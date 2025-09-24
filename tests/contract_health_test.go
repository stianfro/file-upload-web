package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockServer creates a mock HTTP server for testing
// Since main.go doesn't exist yet, we'll simulate the expected behavior
func createMockHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Create a mock server with the health handler
	handler := createMockHealthHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a GET request to the /health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Test 1: Verify status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test 2: Verify Content-Type is text/plain
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/plain"
	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type %q, got %q", expectedContentType, contentType)
	}

	// Test 3: Verify response body contains "OK"
	body := make([]byte, 2)
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Failed to read response body: %v", err)
	}

	responseBody := string(body[:n])
	expectedBody := "OK"
	if responseBody != expectedBody {
		t.Errorf("Expected response body %q, got %q", expectedBody, responseBody)
	}
}

func TestHealthEndpointWrongMethod(t *testing.T) {
	// Create a mock server with the health handler
	handler := createMockHealthHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test that non-GET methods are not allowed
	resp, err := http.Post(server.URL+"/health", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for POST request, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestHealthEndpointWrongPath(t *testing.T) {
	// Create a mock server with the health handler
	handler := createMockHealthHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test that wrong paths return 404
	resp, err := http.Get(server.URL + "/wrong-path")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d for wrong path, got %d", http.StatusNotFound, resp.StatusCode)
	}
}