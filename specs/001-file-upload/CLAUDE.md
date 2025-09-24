# Implementation Guidance for Claude Code

## Overview
This document provides specific guidance for implementing the file upload web application while adhering to the KISS principle and constitution requirements.

## Implementation Order

### 1. Core Server (main.go)
```go
// Structure:
// - Embed HTML template
// - Single main function
// - Three handlers: index, upload, health
// - Configuration from environment
```

Key points:
- Use embed directive for HTML
- No external packages beyond standard library
- Implement streaming for large files
- Clear error messages

### 2. HTML Template (index.html)
```html
<!-- Embedded in main.go -->
<!-- Simple form, no JavaScript -->
<!-- Semantic HTML5 elements -->
```

Requirements:
- Single form element
- POST to /upload
- multipart/form-data encoding
- Basic styling with inline CSS

### 3. Dockerfile
```dockerfile
# Multi-stage build
# FROM golang:alpine AS builder
# FROM alpine:latest for runtime
```

Objectives:
- Static binary compilation
- Minimal final image
- Non-root user
- Health check command

### 4. GitHub Actions (.github/workflows/build.yml)
```yaml
# Triggers: push to main, tags
# Steps: test, build, push to GHCR
# Multi-arch: amd64, arm64
```

Requirements:
- Use GitHub token for auth
- Tag with version and latest
- Build cache optimization
- Security scanning

### 5. Helm Chart (charts/file-upload-web/)
```
charts/file-upload-web/
├── Chart.yaml
├── values.yaml
└── templates/
    ├── deployment.yaml
    ├── service.yaml
    └── pvc.yaml
```

Principles:
- No complex templating
- Sensible defaults
- Optional persistence
- Standard probes

### 6. Documentation (README.md)
Structure:
- Purpose (debugging tool)
- Quick start examples
- curl upload instructions
- Docker/K8s deployment
- Troubleshooting guide

## Code Standards

### Error Handling
```go
if err != nil {
    http.Error(w, "Clear error message", http.StatusCode)
    log.Printf("Context: %v", err)
    return
}
```

### Logging
- Errors only by default
- Request logging optional (DEBUG env var)
- Structured format: timestamp, level, message

### Configuration
```go
port := getEnv("PORT", "8080")
uploadDir := getEnv("UPLOAD_DIR", "/uploads")
maxSize := getEnvInt("MAX_SIZE", 10) * 1024 * 1024
```

## Testing Approach

### Manual Testing Script
```bash
#!/bin/bash
# test.sh
echo "Testing health endpoint..."
curl -f http://localhost:8080/health || exit 1

echo "Testing file upload..."
echo "test content" > test.txt
curl -f -X POST -F "file=@test.txt" http://localhost:8080/upload || exit 1

echo "All tests passed!"
```

### Integration Testing
- Use GitHub Actions to test container
- Verify Helm chart with kind/k3s
- Test both amd64 and arm64 builds

## Common Pitfalls to Avoid

1. **Over-engineering**
   - No logging frameworks
   - No configuration libraries
   - No web frameworks

2. **Security Issues**
   - Always sanitize filenames
   - Set max request size
   - No path traversal

3. **Complexity Creep**
   - Resist adding features
   - Keep it debuggable
   - Maintain single purpose

## File Size Handling

For large files:
```go
// Stream to disk, don't buffer in memory
file, header, err := r.FormFile("file")
defer file.Close()

dst, err := os.Create(filepath)
defer dst.Close()

io.Copy(dst, file) // Streams efficiently
```

## Filename Sanitization

```go
func sanitizeFilename(name string) string {
    // Remove path components
    name = filepath.Base(name)
    // Replace problematic characters
    replacer := strings.NewReplacer(" ", "_", "/", "-")
    return replacer.Replace(name)
}
```

## Success Criteria

Your implementation is complete when:
1. Browser upload works with drag/drop file selection
2. curl upload works as documented
3. Docker image is under 20MB
4. Helm chart deploys successfully
5. README includes working examples
6. No external dependencies in go.mod

## Constitution Compliance Checklist

Before committing:
- [ ] Single file main.go under 200 lines?
- [ ] Zero external dependencies?
- [ ] Clear error messages?
- [ ] Works without configuration?
- [ ] Docker image minimal?
- [ ] Documentation unnecessary (self-explanatory)?

## Debugging Tips

If uploads fail:
1. Check `docker logs` for errors
2. Verify volume mount permissions
3. Test with small file first
4. Use curl -v for verbose output
5. Check disk space in container

Remember: This is a debugging tool. Keep it simple, make errors obvious, and prioritize clarity over cleverness.