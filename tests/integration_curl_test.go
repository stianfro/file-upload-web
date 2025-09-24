package tests

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	serverURL = "http://localhost:8080/upload"
	testDir   = "/tmp/file_upload_test"
)

func TestMain(m *testing.M) {
	// Setup test directory and files
	setupTestEnvironment()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestEnvironment()

	os.Exit(code)
}

func setupTestEnvironment() {
	// Create test directory
	os.MkdirAll(testDir, 0755)

	// Create test files
	createTestFile("simple.txt", "Hello, World!")
	createTestFile("file with spaces.txt", "File with spaces in name")
	createTestFile("special-chars_@#$.txt", "File with special characters")
	createTestFile("large.txt", strings.Repeat("Large file content\n", 100))
	createTestFile("binary.dat", string([]byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}))
}

func cleanupTestEnvironment() {
	os.RemoveAll(testDir)
}

func createTestFile(filename, content string) {
	filepath := filepath.Join(testDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		panic(fmt.Sprintf("Failed to create test file %s: %v", filename, err))
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		panic(fmt.Sprintf("Failed to write to test file %s: %v", filename, err))
	}
}

func isServerRunning() bool {
	cmd := exec.Command("curl", "-s", "-f", "-o", "/dev/null", serverURL)
	err := cmd.Run()
	return err == nil
}

func waitForServer(timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		if isServerRunning() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not available after %v", timeout)
}

func TestBasicFileUpload(t *testing.T) {
	if err := waitForServer(5 * time.Second); err != nil {
		t.Skip("Server not running:", err)
	}

	testFile := filepath.Join(testDir, "simple.txt")

	// Test basic file upload using curl
	cmd := exec.Command("curl",
		"-X", "POST",
		"-F", fmt.Sprintf("file=@%s", testFile),
		"-w", "%{http_code}",
		"-s", "-o", "/dev/null",
		serverURL,
	)

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("curl command failed: %v", err)
	}

	statusCode := strings.TrimSpace(string(output))
	if statusCode != "200" && statusCode != "201" {
		t.Errorf("Expected status code 200 or 201, got: %s", statusCode)
	}

	t.Logf("Basic file upload test completed with status: %s", statusCode)
}

func TestMultipleFilesSequential(t *testing.T) {
	if err := waitForServer(5 * time.Second); err != nil {
		t.Skip("Server not running:", err)
	}

	testFiles := []string{"simple.txt", "large.txt", "binary.dat"}

	for i, filename := range testFiles {
		t.Run(fmt.Sprintf("Upload_%d_%s", i+1, filename), func(t *testing.T) {
			testFile := filepath.Join(testDir, filename)

			cmd := exec.Command("curl",
				"-X", "POST",
				"-F", fmt.Sprintf("file=@%s", testFile),
				"-w", "%{http_code}",
				"-s", "-o", "/dev/null",
				serverURL,
			)

			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("curl command failed for %s: %v", filename, err)
			}

			statusCode := strings.TrimSpace(string(output))
			if statusCode != "200" && statusCode != "201" {
				t.Errorf("Expected status code 200 or 201 for %s, got: %s", filename, statusCode)
			}

			t.Logf("Sequential upload %d (%s) completed with status: %s", i+1, filename, statusCode)
		})
	}
}

func TestFilenameSpecialCharacters(t *testing.T) {
	if err := waitForServer(5 * time.Second); err != nil {
		t.Skip("Server not running:", err)
	}

	testCases := []struct {
		name        string
		filename    string
		description string
	}{
		{
			name:        "SpacesInFilename",
			filename:    "file with spaces.txt",
			description: "File with spaces in filename",
		},
		{
			name:        "SpecialCharacters",
			filename:    "special-chars_@#$.txt",
			description: "File with special characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFile := filepath.Join(testDir, tc.filename)

			// Verify test file exists
			if _, err := os.Stat(testFile); os.IsNotExist(err) {
				t.Fatalf("Test file does not exist: %s", testFile)
			}

			cmd := exec.Command("curl",
				"-X", "POST",
				"-F", fmt.Sprintf("file=@%s", testFile),
				"-w", "%{http_code}",
				"-s", "-o", "/dev/null",
				serverURL,
			)

			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("curl command failed for %s: %v", tc.filename, err)
			}

			statusCode := strings.TrimSpace(string(output))
			if statusCode != "200" && statusCode != "201" {
				t.Errorf("Expected status code 200 or 201 for %s, got: %s", tc.filename, statusCode)
			}

			t.Logf("Special character filename test (%s) completed with status: %s", tc.description, statusCode)
		})
	}
}

func TestConcurrentUploads(t *testing.T) {
	if err := waitForServer(5 * time.Second); err != nil {
		t.Skip("Server not running:", err)
	}

	const numConcurrent = 5
	testFiles := []string{"simple.txt", "large.txt", "binary.dat", "file with spaces.txt", "special-chars_@#$.txt"}

	var wg sync.WaitGroup
	results := make(chan struct {
		filename   string
		statusCode string
		err        error
	}, numConcurrent)

	// Launch concurrent uploads
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			filename := testFiles[index%len(testFiles)]
			testFile := filepath.Join(testDir, filename)

			cmd := exec.Command("curl",
				"-X", "POST",
				"-F", fmt.Sprintf("file=@%s", testFile),
				"-w", "%{http_code}",
				"-s", "-o", "/dev/null",
				serverURL,
			)

			output, err := cmd.Output()
			statusCode := ""
			if err == nil {
				statusCode = strings.TrimSpace(string(output))
			}

			results <- struct {
				filename   string
				statusCode string
				err        error
			}{filename, statusCode, err}
		}(i)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and validate results
	successCount := 0
	for result := range results {
		if result.err != nil {
			t.Errorf("Concurrent upload failed for %s: %v", result.filename, result.err)
			continue
		}

		if result.statusCode != "200" && result.statusCode != "201" {
			t.Errorf("Expected status code 200 or 201 for concurrent upload of %s, got: %s", result.filename, result.statusCode)
		} else {
			successCount++
		}

		t.Logf("Concurrent upload of %s completed with status: %s", result.filename, result.statusCode)
	}

	if successCount == 0 {
		t.Fatal("No concurrent uploads succeeded")
	}

	t.Logf("Concurrent uploads completed: %d/%d successful", successCount, numConcurrent)
}

func TestUploadWithAdditionalFormData(t *testing.T) {
	if err := waitForServer(5 * time.Second); err != nil {
		t.Skip("Server not running:", err)
	}

	testFile := filepath.Join(testDir, "simple.txt")

	// Test file upload with additional form fields
	cmd := exec.Command("curl",
		"-X", "POST",
		"-F", fmt.Sprintf("file=@%s", testFile),
		"-F", "description=Test file upload",
		"-F", "category=test",
		"-F", "tags=integration,curl,test",
		"-w", "%{http_code}",
		"-s", "-o", "/dev/null",
		serverURL,
	)

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("curl command failed: %v", err)
	}

	statusCode := strings.TrimSpace(string(output))
	if statusCode != "200" && statusCode != "201" {
		t.Errorf("Expected status code 200 or 201, got: %s", statusCode)
	}

	t.Logf("Upload with additional form data completed with status: %s", statusCode)
}

func TestUploadNonExistentFile(t *testing.T) {
	if err := waitForServer(5 * time.Second); err != nil {
		t.Skip("Server not running:", err)
	}

	nonExistentFile := filepath.Join(testDir, "does_not_exist.txt")

	// Test upload of non-existent file (should fail)
	cmd := exec.Command("curl",
		"-X", "POST",
		"-F", fmt.Sprintf("file=@%s", nonExistentFile),
		"-w", "%{http_code}",
		"-s", "-o", "/dev/null",
		serverURL,
	)

	output, err := cmd.Output()
	// We expect this to fail at the curl level or return an error status
	if err == nil {
		statusCode := strings.TrimSpace(string(output))
		t.Logf("Upload of non-existent file returned status: %s", statusCode)
		// Server should handle this gracefully, either 400 or 404
		if statusCode == "200" || statusCode == "201" {
			t.Error("Upload of non-existent file should not succeed")
		}
	} else {
		t.Logf("Upload of non-existent file failed as expected: %v", err)
	}
}

// Benchmark functions for performance testing
func BenchmarkSingleFileUpload(b *testing.B) {
	if err := waitForServer(5 * time.Second); err != nil {
		b.Skip("Server not running:", err)
	}

	testFile := filepath.Join(testDir, "simple.txt")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("curl",
			"-X", "POST",
			"-F", fmt.Sprintf("file=@%s", testFile),
			"-s", "-o", "/dev/null",
			serverURL,
		)

		err := cmd.Run()
		if err != nil {
			b.Fatalf("curl command failed: %v", err)
		}
	}
}

func BenchmarkLargeFileUpload(b *testing.B) {
	if err := waitForServer(5 * time.Second); err != nil {
		b.Skip("Server not running:", err)
	}

	testFile := filepath.Join(testDir, "large.txt")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("curl",
			"-X", "POST",
			"-F", fmt.Sprintf("file=@%s", testFile),
			"-s", "-o", "/dev/null",
			serverURL,
		)

		err := cmd.Run()
		if err != nil {
			b.Fatalf("curl command failed: %v", err)
		}
	}
}