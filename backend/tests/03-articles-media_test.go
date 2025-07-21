package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

// 1. Загрузка одного медиафайла, создание статьи, проверка содержимого
func TestUploadSingleMediaAndCreateArticle(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Получаем presigned URL через создание статьи
    mediaFile := "testfile1.jpg"
    title := "Article with media"
    content := "Some content"
    body := map[string]interface{}{"title": title, "content": content, "media": []string{mediaFile}}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 201 {
        t.Fatalf("CreateArticle: expected 201, got %d", resp.StatusCode)
    }
    var out struct {
        ID    string `json:"id"`
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    if len(out.Media) != 1 {
        t.Fatalf("Expected 1 media, got %d", len(out.Media))
    }
    uploadURL := out.Media[0].S3Key

    // Загружаем файл по presigned URL
    fileContent := []byte("hello world media")
    putResp, err := httpPut(uploadURL, fileContent)
	log.Printf("Upload URL: %s", uploadURL)
    if err != nil {
        t.Fatalf("Failed to upload media: %v", err)
    }
    if putResp.StatusCode != 200 && putResp.StatusCode != 204 {
        t.Fatalf("Media upload failed, status: %d", putResp.StatusCode)
    }

    // Получаем статью и проверяем ссылку
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != 1 {
        t.Errorf("GetArticle: expected 1 media, got %d", len(got.Media))
    }
}

// 2. Загрузка нескольких медиафайлов, создание статьи, проверка всех ссылок
func TestUploadMultipleMediaAndCreateArticle(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    mediaFiles := []string{"fileA.jpg", "fileB.png"}
    body := map[string]interface{}{"title": "MultiMedia", "content": "Multi", "media": mediaFiles}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 201 {
        t.Fatalf("CreateArticle: expected 201, got %d", resp.StatusCode)
    }
    var out struct {
        ID    string `json:"id"`
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    if len(out.Media) != len(mediaFiles) {
        t.Fatalf("Expected %d media, got %d", len(mediaFiles), len(out.Media))
    }
    for _, m := range out.Media {
        putResp, err := httpPut(m.S3Key, []byte("media content"))
        if err != nil {
            t.Fatalf("Failed to upload media: %v", err)
        }
        if putResp.StatusCode != 200 && putResp.StatusCode != 204 {
            t.Fatalf("Media upload failed, status: %d", putResp.StatusCode)
        }
    }
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != len(mediaFiles) {
        t.Errorf("GetArticle: expected %d media, got %d", len(mediaFiles), len(got.Media))
    }
}

// 3. Обновление статьи с новым набором медиа
func TestUpdateArticleMedia(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    mediaFiles := []string{"old1.jpg"}
    body := map[string]interface{}{"title": "UpdateMedia", "content": "Update", "media": mediaFiles}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    var out struct {
        ID    string `json:"id"`
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    oldMediaURL := out.Media[0].S3Key
    httpPut(oldMediaURL, []byte("old media"))

    // Обновляем статью с новым медиа
    newMedia := []string{"new1.png"}
    updateBody := map[string]interface{}{"media": newMedia}
    b, _ = json.Marshal(updateBody)
    req, _ = http.NewRequest("PUT", apiBase+"/articles/"+out.ID, bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("UpdateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    var updateOut struct {
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&updateOut)
    newMediaURL := updateOut.Media[0].S3Key
    httpPut(newMediaURL, []byte("new media"))
}

// 4. Попытка загрузить файл по просроченному presigned URL
func TestExpiredPresignedURL(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    mediaFile := "expire.jpg"
    body := map[string]interface{}{"title": "Expire", "content": "Expire", "media": []string{mediaFile}}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    var out struct {
        ID    string `json:"id"`
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    uploadURL := out.Media[0].S3Key

    // Ждем истечения времени (эмулируем, если возможно)
    time.Sleep(10 * time.Second) // если presigned на 1 час, не сработает, но можно уменьшить TTL в бэке для теста

    // Пытаемся загрузить файл
    putResp, err := httpPut(uploadURL, []byte("expired"))
    if err == nil && (putResp.StatusCode == 200 || putResp.StatusCode == 204) {
        t.Errorf("Upload should fail for expired URL, got status %d", putResp.StatusCode)
    }
}

// 5. Получение статьи без медиа
func TestGetArticleWithoutMedia(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    body := map[string]interface{}{"title": "NoMedia", "content": "NoMedia", "media": []string{}}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", apiBase+"/articles/", bytes.NewReader(b))
    req.Header.Set("Authorization", "Bearer "+access)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("CreateArticle failed: %v", err)
    }
    defer resp.Body.Close()
    var out struct {
        ID    string `json:"id"`
        Media []struct {
            FileName string `json:"file_name"`
            S3Key    string `json:"s3_key"`
        } `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != 0 {
        t.Errorf("Expected no media, got %d", len(got.Media))
    }
}

// Вспомогательная функция для PUT-запроса
func httpPut(url string, data []byte) (*http.Response, error) {
    req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/octet-stream")
    client := &http.Client{}
    return client.Do(req)
}