package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

// Add a named struct type for articles
type Article struct {
	UUID    string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
    Media   []struct {
        FileName string `json:"file_name"`
        S3Key    string `json:"s3_key"`
    } `json:"media"`
}

func TestArticleCRUD(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Создание статьи (валидно)
    title := "Test Article"
    content := "This is the content"
    articleUUID := CreateArticle(t, access, title, content, 201)

    // Создание статьи (короткий title)
    CreateArticle(t, access, "ab", content, 400)

    // Создание статьи (пустой контент)
    CreateArticle(t, access, "Valid Title", "", 400)

    // Получение статьи
    got := GetArticle(t, articleUUID, 200)
    if got.Title != title || got.Content != content || got.Author != username {
        t.Errorf("GetArticle: wrong data")
    }

    // Получение несуществующей статьи
    GetArticle(t, "00000000-0000-0000-0000-000000000000", 404)

    // Изменение статьи (валидно)
    newContent := "Updated content"
    UpdateArticle(t, access, articleUUID, title, newContent, 200)

    // Проверка изменения
    got = GetArticle(t, articleUUID, 200)
    if got.Content != newContent {
        t.Errorf("UpdateArticle: content not updated")
    }

    // Изменение статьи (короткий title)
    UpdateArticle(t, access, articleUUID, "ab", newContent, 400)

    // Изменение статьи (пустой контент)
    UpdateArticle(t, access, articleUUID, title, "", 400)
}

func CreateArticle(t *testing.T, access, title, content string, wantStatus int) string {
    body := map[string]string{"title": title, "content": content}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != wantStatus {
        t.Fatalf("CreateArticle: expected %d, got %d", wantStatus, resp.StatusCode)
    }
    if wantStatus == 201 {
        var out Article
        json.NewDecoder(resp.Body).Decode(&out)
        return out.UUID
    }
    return "not 201"
}

func GetArticle(t *testing.T, uuid string, wantStatus int) Article {
    resp, err := http.Get(fmt.Sprintf("%s/articles/%s", apiBase, uuid))
    if err != nil {
        t.Fatalf("GetArticle failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != wantStatus {
        t.Fatalf("GetArticle: expected %d, got %d", wantStatus, resp.StatusCode)
    }
    var out Article
    if wantStatus == 200 {
        json.NewDecoder(resp.Body).Decode(&out)
    }
    return out
}

func UpdateArticle(t *testing.T, access string, uuid string, title, content string, wantStatus int) {
    body := map[string]string{"title": title, "content": content}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/articles/%s", apiBase, uuid), bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("UpdateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != wantStatus {
        t.Fatalf("UpdateArticle: expected %d, got %d", wantStatus, resp.StatusCode)
    }
}

func TestArticleMediaPresignedURLs(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Создание статьи с медиа
    title := "Media Article"
    content := "Article with media"
    mediaFiles := []string{"file1.jpg", "file2.png"}
    body := map[string]interface{}{"title": title, "content": content, "media": mediaFiles}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle with media failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 201 {
        t.Fatalf("CreateArticle with media: expected 201, got %d", resp.StatusCode)
    }
    var out struct {
        ID      string `json:"id"`
        Media   []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    log.Printf("Body: %+v", out)
    if len(out.Media) != len(mediaFiles) {
        t.Errorf("Expected %d media, got %d", len(mediaFiles), len(out.Media))
    }
    for i, m := range out.Media {
        if m.FileName != mediaFiles[i] {
            t.Errorf("Media filename mismatch: expected %s, got %s", mediaFiles[i], m.FileName)
        }
        if !isS3URL(m.S3Key) {
            t.Errorf("Media S3Key is not a valid S3 URL: %s", m.S3Key)
        }
    }

    // Получение статьи и проверка media
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != len(mediaFiles) {
        t.Errorf("GetArticle: expected %d media, got %d", len(mediaFiles), len(got.Media))
    }
    for _, m := range got.Media {
        if !isS3URL(m.S3Key) {
            t.Errorf("GetArticle: S3Key is not a valid S3 URL: %s", m.S3Key)
        }
    }

    // Обновление статьи с новыми медиа
    newMedia := []string{"file3.gif"}
    updateBody := map[string]interface{}{"media": newMedia}
    b, _ = json.Marshal(updateBody)
    req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/articles/%s", apiBase, out.ID), bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("UpdateArticle with media failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        t.Fatalf("UpdateArticle with media: expected 200, got %d", resp.StatusCode)
    }
    var updateOut struct {
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&updateOut)
    if len(updateOut.Media) != len(newMedia) {
        t.Errorf("UpdateArticle: expected %d media, got %d", len(newMedia), len(updateOut.Media))
    }
    for i, m := range updateOut.Media {
        if m.FileName != newMedia[i] {
            t.Errorf("UpdateArticle: media filename mismatch: expected %s, got %s", newMedia[i], m.FileName)
        }
        if !isS3URL(m.S3Key) {
            t.Errorf("UpdateArticle: S3Key is not a valid S3 URL: %s", m.S3Key)
        }
    }
}

// Проверка, что строка похожа на S3 presigned URL
func isS3URL(url string) bool {
    return len(url) > 10 && (url[:4] == "http" || url[:5] == "https")
}