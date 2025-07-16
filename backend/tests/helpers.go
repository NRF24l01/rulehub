package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var apiBase = getAPIBase()

func getAPIBase() string {
    if v := os.Getenv("API_BASE_URL"); v != "" {
        return v
    }
    return "http://localhost:8080/api"
}

func ResetDB(t *testing.T) {
    resp, err := http.Post(apiBase+"/dev/reset-db", "application/json", nil)
    if err != nil {
        t.Fatalf("DB reset failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        t.Fatalf("DB reset failed, status: %d", resp.StatusCode)
    }
}

func RegisterUser(t *testing.T, username, password string) int {
    body := map[string]string{"username": username, "password": password}
    b, _ := json.Marshal(body)
    resp, err := http.Post(apiBase+"/auth/register", "application/json", bytes.NewReader(b))
    if err != nil {
        t.Fatalf("Register failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 201 {
        t.Fatalf("Register failed, status: %d", resp.StatusCode)
    }
    var out struct {
        ID       int    `json:"id"`
        Username string `json:"username"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    return out.ID
}

func LoginUser(t *testing.T, username, password string) (string, string) {
    body := map[string]string{"username": username, "password": password}
    b, _ := json.Marshal(body)
    resp, err := http.Post(apiBase+"/auth/login", "application/json", bytes.NewReader(b))
    if err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        t.Fatalf("Login failed, status: %d", resp.StatusCode)
    }
    var out struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
    }
    json.NewDecoder(resp.Body).Decode(&out)
    return out.AccessToken, out.RefreshToken
}

func UniqueUser() (string, string) {
    return fmt.Sprintf("user_%d", time.Now().UnixNano()), "password123"
}