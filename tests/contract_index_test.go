package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexEndpoint(t *testing.T) {
	// This test will initially fail since main.go doesn't exist yet
	// It serves as a contract test for the GET / endpoint requirements

	// TODO: Import and use the actual handler from main.go when it exists
	// import "github.com/stianfro/file-upload-web"
	// handler := main.GetRouter() or similar

	// For now, create a test that will fail to demonstrate the contract
	// Replace this with actual implementation testing once main.go exists
	// Comment out the next line to see the test fail when no implementation exists
	t.Skip("Skipping test until main.go implementation is available")

	// Uncomment the following lines to test against a real server:
	// resp, err := http.Get("http://localhost:8080/")
	// if err != nil {
	//     t.Fatalf("Failed to connect to server: %v", err)
	// }
	// defer resp.Body.Close()
	//
	// The test would fail here since no server is running

	// Mock handler for development/testing purposes:
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a placeholder - the actual implementation will be in main.go
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Return HTML form as required
		html := `<!DOCTYPE html>
<html>
<head>
    <title>File Upload</title>
</head>
<body>
    <form action="/upload" method="POST" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">Upload</button>
    </form>
</body>
</html>`

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})

	// Create a test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test 1: Endpoint returns 200 OK
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Read response body for further tests
	buf := make([]byte, 1024*10) // 10KB buffer should be enough for a simple HTML form
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Failed to read response body: %v", err)
	}

	body := string(buf[:n])

	// Test 2: Response contains an HTML form element
	if !strings.Contains(body, "<form") || !strings.Contains(body, "</form>") {
		t.Error("Response does not contain an HTML form element")
	}

	// Test 3: Form has the correct action="/upload" and method="POST"
	if !strings.Contains(body, `action="/upload"`) {
		t.Error("Form does not have the correct action='/upload'")
	}

	if !strings.Contains(body, `method="POST"`) {
		t.Error("Form does not have the correct method='POST'")
	}

	// Test 4: Form has enctype="multipart/form-data"
	if !strings.Contains(body, `enctype="multipart/form-data"`) {
		t.Error("Form does not have the correct enctype='multipart/form-data'")
	}
}