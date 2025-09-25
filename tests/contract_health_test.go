package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	defer resp.Body.Close()

	t.Run("GET /health returns 200 OK", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /health returns OK body", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(body) != "OK" {
			t.Errorf("Expected body to be 'OK', got '%s'", string(body))
		}
	})

	t.Run("GET /health returns text/plain Content-Type", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/plain") {
			t.Errorf("Expected Content-Type to contain 'text/plain', got '%s'", contentType)
		}
	})

	t.Run("POST /health returns 405 Method Not Allowed", func(t *testing.T) {
		resp, err := http.Post("http://localhost:8080/health", "text/plain", strings.NewReader("test"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for POST, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /health returns 405 Method Not Allowed", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/health", strings.NewReader("test"))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for PUT, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /health returns 405 Method Not Allowed", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for DELETE, got %d", resp.StatusCode)
		}
	})
}