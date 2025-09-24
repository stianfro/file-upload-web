# Research: File Upload Web Application

## Technology Analysis

### Server Implementation
**Selected: Go with net/http**
- Single static binary deployment
- Zero runtime dependencies
- Built-in HTTP server
- Excellent Docker container size (~10MB)
- Native multipart form handling

**Alternatives Considered:**
- Python http.server: Requires Python runtime, larger image
- Node.js: Heavy runtime, requires npm dependencies
- Nginx + CGI: More complex configuration

### Frontend Approach
**Selected: Pure HTML5**
- Native file input element
- Standard form POST
- No JavaScript required
- Works in all browsers
- Zero build process

### Container Strategy
**Selected: Multi-stage Docker build**
```dockerfile
# Build stage: golang:alpine
# Runtime stage: scratch or alpine
```
- Minimal attack surface
- ~10MB final image size
- No shell access in production

### CI/CD Pipeline
**Selected: GitHub Actions**
- Native to repository
- Free for public repos
- Direct GHCR integration
- Supports multi-arch builds

### Kubernetes Deployment
**Selected: Helm 3**
- OCI registry support
- Simple values override
- Standard K8s resources only
- No CRDs required

## Technical Decisions

### File Storage
- Local filesystem in `/uploads`
- Mounted as volume in K8s
- No database required
- Files stored with timestamp prefix

### Size Limits
- Default: 10MB per file
- Configurable via MAX_SIZE env var
- Enforced at server level
- Clear error messages

### Security Considerations
- No file execution permissions
- Sanitized filenames
- MIME type validation
- Rate limiting (optional)
- No directory traversal

### Monitoring
- `/health` endpoint for probes
- Basic request logging
- Upload success/failure metrics
- No external dependencies

## Implementation Constraints

### Must Have
1. HTML form upload interface
2. curl POST support
3. Docker image < 20MB
4. Helm chart with basic options
5. GitHub Actions automation

### Nice to Have (Future)
1. Upload progress indication
2. File listing endpoint
3. Prometheus metrics
4. Rate limiting

### Out of Scope
1. Authentication/authorization
2. File retrieval/download
3. Virus scanning
4. S3/cloud storage
5. WebSocket progress

## Risk Mitigation

### Disk Space
- **Risk**: Container fills up
- **Mitigation**: Volume mount, size limits

### Large Files
- **Risk**: Memory exhaustion
- **Mitigation**: Streaming to disk, size cap

### Concurrent Uploads
- **Risk**: Resource contention
- **Mitigation**: Go routines, connection limits

### Path Traversal
- **Risk**: Writing outside upload dir
- **Mitigation**: Filename sanitization

## Validation Approach

### Local Testing
```bash
# HTML form test
curl -X POST -F "file=@test.txt" http://localhost:8080/upload

# Direct POST test
curl -X POST --data-binary @test.txt http://localhost:8080/upload
```

### Docker Testing
```bash
docker build -t upload-test .
docker run -p 8080:8080 -v /tmp/uploads:/uploads upload-test
```

### Kubernetes Testing
```bash
helm install test-upload ./charts/file-upload-web
kubectl port-forward svc/file-upload-web 8080:80
```

## Performance Targets

- Startup time: < 1 second
- Upload latency: < 500ms for 1MB file
- Memory usage: < 50MB baseline
- Concurrent uploads: 10 simultaneous
- Container size: < 20MB

## Conclusion

The selected stack (Go + pure HTML + Docker + Helm) provides the simplest possible solution while meeting all requirements. No external dependencies, minimal attack surface, and straightforward debugging capabilities align perfectly with the KISS principle.