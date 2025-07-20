package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// Add a named struct type for articles
type Article struct {
	UUID    string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
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