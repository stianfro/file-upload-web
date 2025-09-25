package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

// TestBrowserUploadFlow simulates a full browser upload flow
func TestBrowserUploadFlow(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("complete browser upload flow", func(t *testing.T) {
		// Step 1: GET / to retrieve the form
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Fatalf("Failed to GET /: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// Read the HTML content
		htmlContent, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read HTML content: %v", err)
		}

		// Step 2: Parse the form to extract action and method
		actionRegex := regexp.MustCompile(`action="([^"]+)"`)
		methodRegex := regexp.MustCompile(`method="([^"]+)"`)

		actionMatches := actionRegex.FindStringSubmatch(string(htmlContent))
		if len(actionMatches) < 2 {
			t.Fatal("Could not find form action")
		}
		formAction := actionMatches[1]

		methodMatches := methodRegex.FindStringSubmatch(string(htmlContent))
		if len(methodMatches) < 2 {
			t.Fatal("Could not find form method")
		}
		formMethod := methodMatches[1]

		// Verify form attributes
		if formAction != "/upload" {
			t.Errorf("Expected form action '/upload', got '%s'", formAction)
		}

		if !strings.EqualFold(formMethod, "POST") {
			t.Errorf("Expected form method 'POST', got '%s'", formMethod)
		}

		// Step 3: POST a file using multipart/form-data
		var uploadBody bytes.Buffer
		writer := multipart.NewWriter(&uploadBody)

		// Add file field
		fileWriter, err := writer.CreateFormFile("file", "test-browser.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		fileContent := "Test content from browser simulation"
		_, err = fileWriter.Write([]byte(fileContent))
		if err != nil {
			t.Fatalf("Failed to write file content: %v", err)
		}

		writer.Close()

		// Make the upload request
		uploadURL := "http://localhost:8080" + formAction
		req, err := http.NewRequest(strings.ToUpper(formMethod), uploadURL, &uploadBody)
		if err != nil {
			t.Fatalf("Failed to create upload request: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		uploadResp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to upload file: %v", err)
		}
		defer uploadResp.Body.Close()

		// Step 4: Verify successful upload response
		if uploadResp.StatusCode != http.StatusOK {
			t.Errorf("Expected upload status 200, got %d", uploadResp.StatusCode)
		}

		responseBody, err := io.ReadAll(uploadResp.Body)
		if err != nil {
			t.Fatalf("Failed to read upload response: %v", err)
		}

		responseText := string(responseBody)
		if !strings.Contains(responseText, "File uploaded successfully") {
			t.Errorf("Expected success message in response, got: %s", responseText)
		}

		if !strings.Contains(responseText, "test-browser.txt") {
			t.Errorf("Expected filename in response, got: %s", responseText)
		}
	})
}

// TestBrowserEdgeCases tests edge cases for browser uploads
func TestBrowserEdgeCases(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("upload with empty filename", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		// Create form file with empty filename
		fileWriter, err := writer.CreateFormFile("file", "")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		_, err = fileWriter.Write([]byte("content"))
		if err != nil {
			t.Fatalf("Failed to write content: %v", err)
		}
		writer.Close()

		resp, err := http.Post(
			"http://localhost:8080/upload",
			writer.FormDataContentType(),
			&body,
		)
		if err != nil {
			t.Fatalf("Failed to upload: %v", err)
		}
		defer resp.Body.Close()

		// Empty filename is a bad request - we need a filename
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for empty filename, got %d", resp.StatusCode)
		}
	})

	t.Run("upload with very long filename", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		// Create a very long filename
		longName := strings.Repeat("a", 300) + ".txt"
		fileWriter, err := writer.CreateFormFile("file", longName)
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		_, err = fileWriter.Write([]byte("content"))
		if err != nil {
			t.Fatalf("Failed to write content: %v", err)
		}
		writer.Close()

		resp, err := http.Post(
			"http://localhost:8080/upload",
			writer.FormDataContentType(),
			&body,
		)
		if err != nil {
			t.Fatalf("Failed to upload: %v", err)
		}
		defer resp.Body.Close()

		// Should succeed with truncated filename
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}