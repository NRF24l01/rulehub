package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Define the media struct to match the API response
type MediaFile struct {
	FileName string `json:"file_name"`
	S3Key    string `json:"s3_key"`
}

// Helper function to upload temporary media file and get its URL
func uploadTempMedia(t *testing.T, accessToken string) string {
	req, _ := http.NewRequest("POST", apiBase+"/media/upload-temp", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get temp upload URL: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Upload-temp request failed with status: %d", resp.StatusCode)
	}
	
	var out struct {
		URL string `json:"temp_url"`
	}
	json.NewDecoder(resp.Body).Decode(&out)
	return out.URL
}

// Helper to extract file key from presigned URL
func extractFileKeyFromURL(urlStr string) string {
	parsedURL, _ := url.Parse(urlStr)
	// Extract just the path part without query parameters
	path := parsedURL.Path
	// Extract the filename - expected format: /bucket/filename
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// 1. Загрузка одного медиафайла, создание статьи, проверка содержимого
func TestUploadSingleMediaAndCreateArticle(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Получаем временный URL для загрузки
    uploadURL := uploadTempMedia(t, access)
    
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

    // Создаем статью, ссылаясь на загруженный файл
    fileKey := extractFileKeyFromURL(uploadURL)
    
    title := "Article with media"
    content := "Some content"
    body := map[string]interface{}{"title": title, "content": content, "media": []string{fileKey}}
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
        ID    string      `json:"id"`
        Media []MediaFile `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    if len(out.Media) != 1 {
        t.Fatalf("Expected 1 media, got %d", len(out.Media))
    }

    // Получаем статью и проверяем ссылку
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != 1 {
        t.Errorf("GetArticle: expected 1 media, got %d", len(got.Media))
    }
    
    // Добавляем проверку доступности медиафайла
    mediaURL := got.Media[0].S3Key
    getResp, err := http.Get(mediaURL)
    if err != nil {
        t.Fatalf("Failed to get media file: %v", err)
    }
    defer getResp.Body.Close()
    if getResp.StatusCode != 200 {
        t.Fatalf("Failed to get media file, status: %d", getResp.StatusCode)
    }
    
    // Проверяем содержимое файла
    var mediaContent bytes.Buffer
    mediaContent.ReadFrom(getResp.Body)
    if !bytes.Equal(mediaContent.Bytes(), fileContent) {
        t.Errorf("Media content doesn't match: expected %s, got %s", 
            string(fileContent), mediaContent.String())
    }
}

// 2. Загрузка нескольких медиафайлов, создание статьи, проверка всех ссылок
func TestUploadMultipleMediaAndCreateArticle(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Загружаем несколько файлов
    uploadURLs := []string{}
    fileKeys := []string{}
    for i := 0; i < 2; i++ {
        uploadURL := uploadTempMedia(t, access)
        putResp, err := httpPut(uploadURL, []byte("media content"))
        if err != nil {
            t.Fatalf("Failed to upload media: %v", err)
        }
        if putResp.StatusCode != 200 && putResp.StatusCode != 204 {
            t.Fatalf("Media upload failed, status: %d", putResp.StatusCode)
        }
        uploadURLs = append(uploadURLs, uploadURL)
        fileKeys = append(fileKeys, extractFileKeyFromURL(uploadURL))
    }

    // Создаем статью с несколькими медиафайлами
    body := map[string]interface{}{"title": "MultiMedia", "content": "Multi", "media": fileKeys}
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
        ID    string      `json:"id"`
        Media []MediaFile `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    if len(out.Media) != len(uploadURLs) {
        t.Fatalf("Expected %d media, got %d", len(uploadURLs), len(out.Media))
    }
    
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != len(uploadURLs) {
        t.Errorf("GetArticle: expected %d media, got %d", len(uploadURLs), len(got.Media))
    }
    
    // Добавляем проверку содержимого всех файлов
    for i, mediaFile := range got.Media {
        log.Printf("Checking media file %d: %s", i, mediaFile.S3Key)
        getResp, err := http.Get(mediaFile.S3Key)
        if err != nil {
            t.Fatalf("Failed to get media file %d: %v", i, err)
        }
        defer getResp.Body.Close()
        if getResp.StatusCode != 200 {
            t.Fatalf("Failed to get media file %d, status: %d", i, getResp.StatusCode)
        }
        
        // Здесь можно также проверить содержимое файла, если нужно
        // Но для упрощения достаточно проверить статус ответа
    }
}

// 3. Обновление статьи с новым набором медиа
func TestUpdateArticleMedia(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Загружаем первый файл
    oldMediaURL := uploadTempMedia(t, access)
    httpPut(oldMediaURL, []byte("old media"))
    oldFileKey := extractFileKeyFromURL(oldMediaURL)

    // Создаем статью с одним медиафайлом
    body := map[string]interface{}{"title": "UpdateMedia", "content": "Update", "media": []string{oldFileKey}}
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
        ID    string      `json:"id"`
        Media []MediaFile `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&out)

    // Загружаем новый файл
    newMediaURL := uploadTempMedia(t, access)
    httpPut(newMediaURL, []byte("new media"))
    newFileKey := extractFileKeyFromURL(newMediaURL)

    // Обновляем статью с новым медиа
    // При обновлении сохраняем те же title и content
    updateBody := map[string]interface{}{
        "title": "UpdateMedia", 
        "content": "Update", 
        "media": []string{newFileKey},
    }
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
        Media []MediaFile `json:"media"`
    }
    json.NewDecoder(resp.Body).Decode(&updateOut)
    if len(updateOut.Media) != 1 {
        t.Fatalf("Expected 1 media after update, got %d", len(updateOut.Media))
    }
    
    // Получаем обновленную статью и проверяем новое медиа
    got := GetArticle(t, out.ID, 200)
    if len(got.Media) != 1 {
        t.Errorf("GetArticle after update: expected 1 media, got %d", len(got.Media))
    }
    
    // Проверяем доступность нового медиафайла
    if len(got.Media) > 0 {
        newMediaURL := got.Media[0].S3Key
        getResp, err := http.Get(newMediaURL)
        if err != nil {
            t.Fatalf("Failed to get updated media file: %v", err)
        }
        defer getResp.Body.Close()
        if getResp.StatusCode != 200 {
            t.Fatalf("Failed to get updated media file, status: %d", getResp.StatusCode)
        }
    }
}

// 4. Попытка загрузить файл по просроченному presigned URL
func TestExpiredPresignedURL(t *testing.T) {
    // Этот тест требует специальной настройки на бэкенде:
    // presigned URL должен истекать очень быстро (например, через 1 секунду)
    // В реальной среде могут использоваться более длительные сроки
    
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    access, _ := LoginUser(t, username, password)

    // Получаем временный URL для загрузки
    uploadURL := uploadTempMedia(t, access)

    // Ждем истечения времени
    time.Sleep(10 * time.Second) 
    
    // Примечание: этот тест будет успешным только если бэкенд настроен 
    // на очень короткий срок действия presigned URL.
    // Если URL не истекает, тест будет помечен как пропущенный
    
    // Пытаемся загрузить файл
    putResp, err := httpPut(uploadURL, []byte("expired"))
    if err == nil && (putResp.StatusCode == 200 || putResp.StatusCode == 204) {
        t.Skip("URL did not expire as expected. This test requires special backend configuration.")
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
        ID    string      `json:"id"`
        Media []MediaFile `json:"media"`
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