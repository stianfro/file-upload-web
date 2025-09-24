# File Upload Web Application

A simple, lightweight file upload service designed for debugging and troubleshooting file upload functionality across different hosting environments. Built with Go using only standard library dependencies.

## Features

- 🚀 **Simple HTTP server** - No external dependencies
- 📁 **File upload via browser** - HTML form interface
- 🔧 **curl support** - Command-line file uploads
- 🐳 **Docker ready** - < 16MB container image
- ☸️ **Kubernetes ready** - Helm chart included
- 🔒 **Security focused** - Filename sanitization, size limits
- 📊 **Health checks** - Built-in `/health` endpoint

## Quick Start

### Local Development

```bash
# Clone repository
git clone https://github.com/stianfro/file-upload-web.git
cd file-upload-web

# Run directly with Go
go run main.go

# Or build and run
go build .
./file-upload-web

# Access via browser
open http://localhost:8080
```

### Docker

```bash
# Build image
docker build -t file-upload-web .

# Run container
docker run -p 8080:8080 -v $(pwd)/uploads:/uploads file-upload-web

# Or use docker-compose
docker-compose up
```

### Kubernetes with Helm

```bash
# Install from local chart
helm install my-upload ./charts/file-upload-web

# Or from OCI registry (when published)
helm install my-upload oci://ghcr.io/stianfro/file-upload-web --version 0.1.0

# Port forward for testing
kubectl port-forward svc/my-upload-file-upload-web 8080:80
```

## Usage Examples

### Browser Upload
1. Navigate to http://localhost:8080
2. Click "Choose File" and select your file
3. Click "Upload File"
4. See confirmation message

### curl Upload

```bash
# Basic file upload
curl -X POST -F "file=@example.txt" http://localhost:8080/upload

# Upload with progress bar
curl -X POST -F "file=@large-file.zip" http://localhost:8080/upload --progress-bar

# Upload from pipe
echo "test data" | curl -X POST -F "file=@-" http://localhost:8080/upload

# Multiple files (sequential)
for file in *.txt; do
  curl -X POST -F "file=@$file" http://localhost:8080/upload
done
```

### Health Check

```bash
curl http://localhost:8080/health
# Returns: OK
```

## Configuration

Configure via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `UPLOAD_DIR` | `./uploads` | Upload directory path |
| `MAX_SIZE` | `10` | Maximum file size in MB |

### Docker Configuration

```bash
docker run -p 80:8080 \
  -e PORT=8080 \
  -e UPLOAD_DIR=/uploads \
  -e MAX_SIZE=50 \
  -v /local/path:/uploads \
  file-upload-web
```

### Kubernetes Configuration

Edit `values.yaml`:

```yaml
config:
  port: "8080"
  uploadDir: "/uploads"
  maxSize: "50"

persistence:
  enabled: true
  size: 10Gi

service:
  type: LoadBalancer
  port: 80
```

## Troubleshooting

### Upload Fails: 413 Payload Too Large
- **Cause**: File exceeds `MAX_SIZE` limit
- **Solution**: Increase `MAX_SIZE` environment variable or use smaller file

### Upload Fails: 500 Internal Server Error
- **Cause**: Usually permission issues or disk space
- **Solution**:
  ```bash
  # Check upload directory permissions
  ls -la uploads/

  # Check disk space
  df -h

  # Check container logs
  docker logs <container-id>
  ```

### Cannot Access Service
- **Check if service is running**:
  ```bash
  curl localhost:8080/health
  ```
- **Check firewall rules**:
  ```bash
  sudo iptables -L
  ```
- **For Kubernetes**:
  ```bash
  kubectl get pods
  kubectl describe svc my-upload-file-upload-web
  ```

### Files Not Persisted (Docker/K8s)
- **Ensure volume is mounted**:
  ```bash
  # Docker
  docker inspect <container> | grep Mounts -A 10

  # Kubernetes
  kubectl describe pod <pod-name>
  ```
- **Check UPLOAD_DIR matches mount point**
- **Verify PVC is bound (Kubernetes)**:
  ```bash
  kubectl get pvc
  ```

### Filename Issues
- Special characters are automatically sanitized
- Spaces replaced with underscores
- Path traversal attempts blocked

## Testing

```bash
# Run unit tests
go test ./...

# Run manual test suite
./test.sh

# Test Docker image
docker build -t file-upload-web:test .
docker run --rm -p 8080:8080 file-upload-web:test &
sleep 2
curl -f http://localhost:8080/health
```

## Security Notes

- **No authentication** - Do not expose publicly without additional security
- **Filename sanitization** - Automatic removal of dangerous characters
- **Size limits** - Configurable via MAX_SIZE
- **Non-root container** - Runs as user 1000
- **No execution** - Uploaded files have no execute permissions

For production use, consider adding:
- Authentication/authorization
- Rate limiting
- Virus scanning
- Network policies
- TLS/HTTPS

## Development

### Project Structure
```
.
├── main.go              # Application code
├── index.html           # Embedded HTML interface
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Local development
├── go.mod              # Go module definition
├── test.sh             # Manual test suite
├── charts/             # Helm chart
│   └── file-upload-web/
└── .github/            # CI/CD workflows
    └── workflows/
```

### Building

```bash
# Build for current platform
go build -o file-upload-web .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o file-upload-web-linux .

# Build Docker image
docker build -t file-upload-web .
```

## License

MIT

## Contributing

Pull requests welcome! Please keep it simple - this is a debugging tool.

## Support

For issues and questions: https://github.com/stianfro/file-upload-web/issues