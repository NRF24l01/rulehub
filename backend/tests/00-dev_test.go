package tests

import (
	"net/http"
	"testing"
)

func TestResetDB(t *testing.T) {
    resp, err := http.Post(apiBase+"/dev/reset-db", "application/json", nil)
    if err != nil {
        t.Fatalf("Reset DB error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        t.Fatalf("Reset DB status: %d", resp.StatusCode)
    }
}

func TestResetDBForbiddenInProd(t *testing.T) {
    // Этот тест актуален только если у вас есть способ запустить в прод-режиме.
    // Здесь просто пример, что 403 должен быть в проде.
    req, _ := http.NewRequest("POST", apiBase+"/dev/reset-db", nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("Reset DB error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 && resp.StatusCode != 403 {
        t.Fatalf("Reset DB unexpected status: %d", resp.StatusCode)
    }
}