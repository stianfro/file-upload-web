package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// TestBrowserUploadFlow simulates a complete browser upload flow
func TestBrowserUploadFlow(t *testing.T) {
	// This test will initially fail since main.go doesn't exist yet
	// It tests the complete browser upload flow:
	// 1. GET / to retrieve the form
	// 2. Parse the form to extract action and method
	// 3. POST a file using multipart/form-data
	// 4. Verify successful upload response

	// Start the server (this will fail initially since main.go doesn't exist)
	server := httptest.NewServer(getHandler())
	defer server.Close()

	t.Run("CompleteUploadFlow", func(t *testing.T) {
		// Step 1: GET / to retrieve the form
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("Failed to GET /: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		htmlContent := string(body)

		// Step 2: Parse the form to extract action and method
		formAction, formMethod := parseFormDetails(t, htmlContent)

		// Verify we have a file upload form
		if !strings.Contains(htmlContent, `type="file"`) {
			t.Fatal("Form does not contain a file input")
		}

		// Step 3: POST a file using multipart/form-data
		fileContent := "This is a test file content for upload testing"
		fileName := "test-file.txt"

		// Create multipart form data
		body = &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Create file part
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		// Write file content
		_, err = part.Write([]byte(fileContent))
		if err != nil {
			t.Fatalf("Failed to write file content: %v", err)
		}

		// Close the writer to finalize the form
		err = writer.Close()
		if err != nil {
			t.Fatalf("Failed to close multipart writer: %v", err)
		}

		// Determine upload URL
		uploadURL := server.URL
		if formAction != "" && formAction != "/" {
			uploadURL = server.URL + formAction
		}

		// Create POST request
		req, err := http.NewRequest(strings.ToUpper(formMethod), uploadURL, body)
		if err != nil {
			t.Fatalf("Failed to create POST request: %v", err)
		}

		// Set content type header
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send POST request: %v", err)
		}
		defer resp.Body.Close()

		// Step 4: Verify successful upload response
		if resp.StatusCode != http.StatusOK {
			// Read error response for debugging
			errorBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(errorBody))
		}

		// Read successful response
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read upload response: %v", err)
		}

		responseStr := string(responseBody)

		// Verify response indicates successful upload
		successIndicators := []string{
			"success",
			"uploaded",
			"complete",
			fileName, // The uploaded file name should appear in response
		}

		foundIndicator := false
		for _, indicator := range successIndicators {
			if strings.Contains(strings.ToLower(responseStr), strings.ToLower(indicator)) {
				foundIndicator = true
				break
			}
		}

		if !foundIndicator {
			t.Fatalf("Upload response does not indicate success. Response: %s", responseStr)
		}

		t.Logf("Upload test completed successfully")
		t.Logf("Uploaded file: %s (%d bytes)", fileName, len(fileContent))
		t.Logf("Server response: %s", responseStr)
	})
}

// parseFormDetails extracts form action and method from HTML content
func parseFormDetails(t *testing.T, htmlContent string) (action, method string) {
	// Default values
	action = "/"
	method = "POST"

	// Parse form tag to extract action and method
	formRegex := regexp.MustCompile(`<form[^>]*>`)
	formMatch := formRegex.FindString(htmlContent)

	if formMatch != "" {
		// Extract action attribute
		actionRegex := regexp.MustCompile(`action\s*=\s*["']([^"']*)["']`)
		if actionMatch := actionRegex.FindStringSubmatch(formMatch); len(actionMatch) > 1 {
			action = actionMatch[1]
		}

		// Extract method attribute
		methodRegex := regexp.MustCompile(`method\s*=\s*["']([^"']*)["']`)
		if methodMatch := methodRegex.FindStringSubmatch(formMatch); len(methodMatch) > 1 {
			method = methodMatch[1]
		}
	}

	t.Logf("Parsed form details - Action: %s, Method: %s", action, method)
	return action, method
}

// getHandler returns the HTTP handler for testing
// This function will need to be implemented when main.go exists
func getHandler() http.Handler {
	// TODO: This should return the actual handler from main.go
	// For now, return a mock handler that will cause the test to fail
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a placeholder that will cause the test to fail initially
		http.Error(w, "main.go does not exist yet - implement the actual server", http.StatusNotImplemented)
	})
}

// Helper function to verify file upload functionality
func TestFileUploadValidation(t *testing.T) {
	t.Run("ValidateMultipartFormData", func(t *testing.T) {
		// Test multipart form creation
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Create a test file part
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("Failed to create form file part: %v", err)
		}

		testContent := "test file content"
		_, err = part.Write([]byte(testContent))
		if err != nil {
			t.Fatalf("Failed to write to form file part: %v", err)
		}

		err = writer.Close()
		if err != nil {
			t.Fatalf("Failed to close multipart writer: %v", err)
		}

		// Verify content type is correctly set
		contentType := writer.FormDataContentType()
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Fatalf("Expected multipart/form-data content type, got: %s", contentType)
		}

		// Verify body contains the expected boundary
		bodyStr := body.String()
		if !strings.Contains(bodyStr, "test.txt") {
			t.Fatal("Multipart body does not contain filename")
		}

		if !strings.Contains(bodyStr, testContent) {
			t.Fatal("Multipart body does not contain file content")
		}

		t.Log("Multipart form data validation passed")
	})
}