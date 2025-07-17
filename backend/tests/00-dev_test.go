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
