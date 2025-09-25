package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestIndexEndpoint(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("GET / returns 200 OK", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET / returns HTML content", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Check Content-Type
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Expected Content-Type to contain 'text/html', got '%s'", contentType)
		}

		// Read body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		html := string(body)

		// Verify HTML structure
		if !strings.Contains(html, "<form") {
			t.Error("Expected HTML to contain a form element")
		}

		if !strings.Contains(html, `action="/upload"`) {
			t.Error("Expected form action to be '/upload'")
		}

		if !strings.Contains(html, `method="POST"`) {
			t.Error("Expected form method to be 'POST'")
		}

		if !strings.Contains(html, `enctype="multipart/form-data"`) {
			t.Error("Expected form enctype to be 'multipart/form-data'")
		}

		if !strings.Contains(html, `type="file"`) {
			t.Error("Expected form to contain a file input")
		}
	})

	t.Run("POST / returns 405 Method Not Allowed", func(t *testing.T) {
		resp, err := http.Post("http://localhost:8080/", "text/plain", strings.NewReader("test"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /nonexistent returns 404 Not Found", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/nonexistent")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

// TestIndexEndpointPerformance tests the performance of the index endpoint
func TestIndexEndpointPerformance(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("response time under 100ms", func(t *testing.T) {
		start := time.Now()
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()

		elapsed := time.Since(start)
		if elapsed > 100*time.Millisecond {
			t.Errorf("Response time too slow: %v", elapsed)
		}
	})
}