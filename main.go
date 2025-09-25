package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed index.html
var indexHTML string

func main() {
	port := getEnv("PORT", "8080")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")
	maxSizeStr := getEnv("MAX_SIZE", "10")

	maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil {
		maxSize = 10
	}
	maxSizeBytes := maxSize * 1024 * 1024

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Setup routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler(uploadDir, maxSizeBytes))
	http.HandleFunc("/health", healthHandler)

	log.Printf("Server starting on port %s", port)
	log.Printf("Upload directory: %s", uploadDir)
	log.Printf("Max file size: %d MB", maxSize)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func uploadHandler(uploadDir string, maxSizeBytes int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form with size limit
		r.Body = http.MaxBytesReader(w, r.Body, maxSizeBytes)
		if err := r.ParseMultipartForm(maxSizeBytes); err != nil {
			if err.Error() == "http: request body too large" {
				http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "No file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create timestamp first
		timestamp := time.Now().Format("20060102_150405")

		// Sanitize filename - account for timestamp in length limit
		// Timestamp format "20060102_150405_" is 16 chars
		filename := sanitizeFilenameWithMaxLen(header.Filename, 255-16)

		// Create final filename
		finalName := fmt.Sprintf("%s_%s", timestamp, filename)
		filepath := filepath.Join(uploadDir, finalName)

		// Create destination file
		dst, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer dst.Close()

		// Stream file to disk
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			log.Printf("Failed to write file: %v", err)
			return
		}

		log.Printf("File uploaded: %s", finalName)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "File uploaded successfully: %s", finalName)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func sanitizeFilenameWithMaxLen(name string, maxLen int) string {
	// Remove path components
	name = filepath.Base(name)

	// Replace problematic characters
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "-",
		"\\", "-",
		"..", ".",
	)
	name = replacer.Replace(name)

	// Ensure filename is not empty
	if name == "" || name == "." {
		name = "unnamed"
	}

	// Limit length based on provided max length
	if len(name) > maxLen {
		ext := filepath.Ext(name)
		if len(ext) > maxLen {
			// If extension itself is too long, truncate everything
			name = name[:maxLen]
		} else {
			// Keep extension if reasonable, truncate base
			maxBase := maxLen - len(ext)
			if maxBase > 0 {
				base := name[:len(name)-len(ext)]
				if len(base) > maxBase {
					base = base[:maxBase]
				}
				name = base + ext
			} else {
				name = name[:maxLen]
			}
		}
	}

	return name
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}