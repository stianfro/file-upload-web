package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stianfro/file-upload-web"
)

// Response structure expected from the upload endpoint
type UploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Filename string `json:"filename,omitempty"`
}

// TestUploadEndpoint tests the POST /upload endpoint
func TestUploadEndpoint(t *testing.T) {
	// This will fail initially since main.go doesn't exist yet
	// but provides the contract specification for the endpoint

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

		// Create HTTP request
		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler (this will fail initially)
		handler := main.CreateUploadHandler()
		handler.ServeHTTP(rr, req)

		// Verify response
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		// Parse response
		var response UploadResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response contains success message with filename
		if !response.Success {
			t.Error("Expected success to be true")
		}

		if response.Filename == "" {
			t.Error("Expected filename to be present in response")
		}

		if !strings.Contains(response.Message, "test.txt") {
			t.Errorf("Expected message to contain filename 'test.txt', got: %s", response.Message)
		}
	})

	t.Run("no file provided returns 400 Bad Request", func(t *testing.T) {
		// Create empty multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		// Create HTTP request
		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		handler := main.CreateUploadHandler()
		handler.ServeHTTP(rr, req)

		// Verify response
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}

		// Parse response
		var response UploadResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify error response
		if response.Success {
			t.Error("Expected success to be false")
		}

		if response.Message == "" {
			t.Error("Expected error message to be present")
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

		// Create HTTP request
		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		handler := main.CreateUploadHandler()
		handler.ServeHTTP(rr, req)

		// Verify response
		if rr.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status 413, got %d", rr.Code)
		}

		// Parse response
		var response UploadResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify error response
		if response.Success {
			t.Error("Expected success to be false")
		}

		if response.Message == "" {
			t.Error("Expected error message to be present")
		}
	})

	t.Run("multipart form data content type validation", func(t *testing.T) {
		// Test with non-multipart content type
		body := strings.NewReader("not multipart data")

		// Create HTTP request with wrong content type
		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		handler := main.CreateUploadHandler()
		handler.ServeHTTP(rr, req)

		// Verify response - should return 400 for invalid content type
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid content type, got %d", rr.Code)
		}
	})
}

// TestUploadEndpointIntegration tests the upload endpoint with an actual server
func TestUploadEndpointIntegration(t *testing.T) {
	t.Run("integration test with test server", func(t *testing.T) {
		// This test will also fail initially but demonstrates how to test
		// the endpoint in a more realistic scenario

		// Create test server
		server := httptest.NewServer(main.CreateUploadHandler())
		defer server.Close()

		// Create multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "integration_test.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		_, err = io.WriteString(part, "Integration test content")
		if err != nil {
			t.Fatalf("Failed to write content: %v", err)
		}

		writer.Close()

		// Make HTTP request
		resp, err := http.Post(
			server.URL+"/upload",
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

		// Read and parse response
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		var response UploadResponse
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// Verify response structure
		if !response.Success {
			t.Error("Expected success to be true")
		}

		if response.Filename == "" {
			t.Error("Expected filename to be present")
		}
	})
}