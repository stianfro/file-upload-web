package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestUploadEndpoint tests the POST /upload endpoint
func TestUploadEndpoint(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("successful file upload returns 200 OK", func(t *testing.T) {
		// Create a multipart form with a test file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add a test file
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		testContent := "This is a test file content"
		_, err = part.Write([]byte(testContent))
		if err != nil {
			t.Fatalf("Failed to write test content: %v", err)
		}

		writer.Close()

		// Make HTTP request
		resp, err := http.Post(
			"http://localhost:8080/upload",
			writer.FormDataContentType(),
			body,
		)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Verify response
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Verify response message
		message := string(respBody)
		if !strings.Contains(message, "File uploaded successfully") {
			t.Errorf("Expected success message, got: %s", message)
		}

		if !strings.Contains(message, "test.txt") {
			t.Errorf("Expected message to contain filename 'test.txt', got: %s", message)
		}
	})

	t.Run("no file provided returns 400 Bad Request", func(t *testing.T) {
		// Create empty multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		// Make HTTP request
		resp, err := http.Post(
			"http://localhost:8080/upload",
			writer.FormDataContentType(),
			body,
		)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Verify response
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Verify error message
		message := string(respBody)
		if !strings.Contains(message, "No file provided") {
			t.Errorf("Expected 'No file provided' error, got: %s", message)
		}
	})

	t.Run("file too large returns 413 Request Entity Too Large", func(t *testing.T) {
		// Create a multipart form with a large file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add a large test file (simulate a file that exceeds the limit)
		part, err := writer.CreateFormFile("file", "large_file.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		// Create content that would be too large (assuming 10MB limit)
		largeContent := make([]byte, 11*1024*1024) // 11MB
		for i := range largeContent {
			largeContent[i] = 'A'
		}

		_, err = part.Write(largeContent)
		if err != nil {
			t.Fatalf("Failed to write large content: %v", err)
		}

		writer.Close()

		// Make HTTP request
		resp, err := http.Post(
			"http://localhost:8080/upload",
			writer.FormDataContentType(),
			body,
		)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Verify response
		if resp.StatusCode != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status 413, got %d", resp.StatusCode)
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Verify error message
		message := string(respBody)
		if !strings.Contains(message, "File too large") {
			t.Errorf("Expected 'File too large' error, got: %s", message)
		}
	})

	t.Run("multipart form data content type validation", func(t *testing.T) {
		// Test with non-multipart content type
		body := strings.NewReader("not multipart data")

		// Make HTTP request with wrong content type
		req, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Verify response - should return 400 for invalid content type
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid content type, got %d", resp.StatusCode)
		}
	})
}

// TestUploadFilenameWithSpecialCharacters tests filename sanitization
func TestUploadFilenameWithSpecialCharacters(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	specialFilenames := []string{
		"file with spaces.txt",
		"file/with/slashes.txt",
		"file\\with\\backslashes.txt",
		"../../../etc/passwd",
		"file@special#chars$.txt",
	}

	for _, filename := range specialFilenames {
		t.Run(fmt.Sprintf("upload file: %s", filename), func(t *testing.T) {
			// Create multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", filename)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}

			_, err = io.WriteString(part, "Test content")
			if err != nil {
				t.Fatalf("Failed to write content: %v", err)
			}

			writer.Close()

			// Make HTTP request
			resp, err := http.Post(
				"http://localhost:8080/upload",
				writer.FormDataContentType(),
				body,
			)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Should succeed with sanitized filename
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for file '%s', got %d", filename, resp.StatusCode)
			}

			// Read response
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Verify upload succeeded
			message := string(respBody)
			if !strings.Contains(message, "File uploaded successfully") {
				t.Errorf("Expected success message for file '%s', got: %s", filename, message)
			}
		})
	}
}