package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRegisterValidation(t *testing.T) {
    ResetDB(t)
    // Успешная регистрация
    username, password := UniqueUser()
    t.Logf("Registering user: username=%s, password=%s", username, password)
    RegisterUser(t, username, password)

    // Короткий username
    body := map[string]string{"username": "ab", "password": "password123"}
    b, _ := json.Marshal(body)
    resp, _ := http.Post(apiBase+"/auth/register", "application/json", bytes.NewReader(b))
    if resp.StatusCode != 400 {
        t.Errorf("Short username: expected 400, got %d", resp.StatusCode)
    }

    // Короткий пароль
    body = map[string]string{"username": "validuser", "password": "123"}
    b, _ = json.Marshal(body)
    resp, _ = http.Post(apiBase+"/auth/register", "application/json", bytes.NewReader(b))
    if resp.StatusCode != 400 {
        t.Errorf("Short password: expected 400, got %d", resp.StatusCode)
    }

    // Невалидные символы
    body = map[string]string{"username": "bad user", "password": "password123"}
    b, _ = json.Marshal(body)
    resp, _ = http.Post(apiBase+"/auth/register", "application/json", bytes.NewReader(b))
    if resp.StatusCode != 400 {
        t.Errorf("Invalid chars: expected 400, got %d", resp.StatusCode)
    }

    // Повторная регистрация
    body = map[string]string{"username": username, "password": password}
    b, _ = json.Marshal(body)
    resp, _ = http.Post(apiBase+"/auth/register", "application/json", bytes.NewReader(b))
    if resp.StatusCode != 409 {
        t.Errorf("Duplicate user: expected 409, got %d", resp.StatusCode)
    }
}

func TestLoginValidation(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)

    // Успешный вход
    access, refresh := LoginUser(t, username, password)
    if len(access) < 16 || len(refresh) < 16 {
        t.Errorf("Tokens too short")
    }

    // Неверный пароль
    _, _ = LoginUserExpectStatus(t, username, "wrongpass123", 401)

    // Неверный пользователь
    _, _ = LoginUserExpectStatus(t, "nouser", "password123", 401)
}

func LoginUserExpectStatus(t *testing.T, username, password string, wantStatus int) (string, string) {
    t.Logf("Logging in user: username=%s, password=%s", username, password)
    body := map[string]string{"username": username, "password": password}
    b, _ := json.Marshal(body)
    resp, err := http.Post(apiBase+"/auth/login", "application/json", bytes.NewReader(b))
    if err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != wantStatus {
        t.Fatalf("Login: expected %d, got %d", wantStatus, resp.StatusCode)
    }
    return "", ""
}

func TestRefreshToken(t *testing.T) {
    ResetDB(t)
    username, password := UniqueUser()
    RegisterUser(t, username, password)
    _, refresh := LoginUser(t, username, password)

    // Валидный refresh через cookie
    req, err := http.NewRequest("POST", apiBase+"/auth/refresh", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.AddCookie(&http.Cookie{
        Name:  "refresh_token",
        Value: refresh,
        Path:  "/",
        HttpOnly: true,
    })
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Refresh failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        t.Fatalf("Refresh: expected 200, got %d", resp.StatusCode)
    }
    var out struct {
        AccessToken string `json:"access_token"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    if len(out.AccessToken) < 16 {
        t.Errorf("Refresh: access_token too short")
    }

    // Невалидный refresh через cookie
    req, err = http.NewRequest("POST", apiBase+"/auth/refresh", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.AddCookie(&http.Cookie{
        Name:  "refresh_token",
        Value: "badtoken",
        Path:  "/",
        HttpOnly: true,
    })
    resp, err = client.Do(req)
    if err != nil {
        t.Fatalf("Refresh failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 401 {
        t.Errorf("Refresh: expected 401, got %d", resp.StatusCode)
    }
}