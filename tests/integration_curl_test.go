package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestCurlUploadScenarios tests curl upload scenarios
func TestCurlUploadScenarios(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("basic file upload with form data", func(t *testing.T) {
		// Simulate: curl -X POST -F "file=@test.txt" http://localhost:8080/upload
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		fileWriter, err := writer.CreateFormFile("file", "curl_test.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		_, err = io.WriteString(fileWriter, "Test content from curl simulation")
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

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		if !strings.Contains(string(respBody), "File uploaded successfully") {
			t.Errorf("Expected success message, got: %s", string(respBody))
		}
	})

	t.Run("multiple files sequential upload", func(t *testing.T) {
		// Simulate uploading multiple files one after another
		files := []struct {
			name    string
			content string
		}{
			{"file1.txt", "Content of file 1"},
			{"file2.txt", "Content of file 2"},
			{"file3.bin", "Binary data simulation"},
		}

		for _, file := range files {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			fileWriter, err := writer.CreateFormFile("file", file.name)
			if err != nil {
				t.Fatalf("Failed to create form file %s: %v", file.name, err)
			}

			_, err = io.WriteString(fileWriter, file.content)
			if err != nil {
				t.Fatalf("Failed to write content for %s: %v", file.name, err)
			}

			writer.Close()

			resp, err := http.Post(
				"http://localhost:8080/upload",
				writer.FormDataContentType(),
				&body,
			)
			if err != nil {
				t.Fatalf("Failed to upload %s: %v", file.name, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Upload of %s failed with status %d", file.name, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response for %s: %v", file.name, err)
			}

			if !strings.Contains(string(respBody), file.name) {
				t.Errorf("Expected filename %s in response, got: %s", file.name, string(respBody))
			}
		}
	})

	t.Run("filename with special characters", func(t *testing.T) {
		specialFilenames := []string{
			"file with spaces.txt",
			"file-with-dashes.txt",
			"file_with_underscores.txt",
			"file.multiple.dots.txt",
		}

		for _, filename := range specialFilenames {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			fileWriter, err := writer.CreateFormFile("file", filename)
			if err != nil {
				t.Fatalf("Failed to create form file %s: %v", filename, err)
			}

			_, err = io.WriteString(fileWriter, "Special filename test")
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
				t.Fatalf("Failed to upload %s: %v", filename, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Upload of %s failed with status %d", filename, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Verify the file was uploaded successfully
			if !strings.Contains(string(respBody), "File uploaded successfully") {
				t.Errorf("Expected success for file %s, got: %s", filename, string(respBody))
			}
		}
	})

	t.Run("concurrent uploads", func(t *testing.T) {
		numUploads := 5
		var wg sync.WaitGroup
		errors := make(chan error, numUploads)

		for i := 0; i < numUploads; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				var body bytes.Buffer
				writer := multipart.NewWriter(&body)

				filename := fmt.Sprintf("concurrent_%d.txt", index)
				fileWriter, err := writer.CreateFormFile("file", filename)
				if err != nil {
					errors <- fmt.Errorf("failed to create form file %s: %v", filename, err)
					return
				}

				content := fmt.Sprintf("Concurrent upload content %d", index)
				_, err = io.WriteString(fileWriter, content)
				if err != nil {
					errors <- fmt.Errorf("failed to write content: %v", err)
					return
				}

				writer.Close()

				resp, err := http.Post(
					"http://localhost:8080/upload",
					writer.FormDataContentType(),
					&body,
				)
				if err != nil {
					errors <- fmt.Errorf("failed to upload %s: %v", filename, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					errors <- fmt.Errorf("upload %s failed with status %d", filename, resp.StatusCode)
					return
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					errors <- fmt.Errorf("failed to read response: %v", err)
					return
				}

				if !strings.Contains(string(respBody), filename) {
					errors <- fmt.Errorf("expected filename %s in response", filename)
					return
				}
			}(i)
		}

		// Wait for all uploads with timeout
		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			// All uploads completed
		case <-time.After(10 * time.Second):
			t.Fatal("Concurrent uploads timed out")
		}

		// Check for errors
		close(errors)
		for err := range errors {
			if err != nil {
				t.Error(err)
			}
		}
	})
}

// TestCurlLargeFileUpload tests uploading larger files
func TestCurlLargeFileUpload(t *testing.T) {
	// Skip if server is not running
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skip("Server not running on localhost:8080, skipping integration tests")
	}
	resp.Body.Close()

	t.Run("upload 5MB file", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		fileWriter, err := writer.CreateFormFile("file", "large_5mb.bin")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		// Create 5MB of data
		largeData := make([]byte, 5*1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		_, err = fileWriter.Write(largeData)
		if err != nil {
			t.Fatalf("Failed to write large data: %v", err)
		}

		writer.Close()

		start := time.Now()
		resp, err := http.Post(
			"http://localhost:8080/upload",
			writer.FormDataContentType(),
			&body,
		)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to upload large file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Check that upload completed in reasonable time (< 5 seconds)
		if elapsed > 5*time.Second {
			t.Errorf("Large file upload took too long: %v", elapsed)
		}
	})
}