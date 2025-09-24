# Quick Start Guide

## Local Development

### Prerequisites
- Go 1.21+ OR Docker
- curl (for testing)

### Run Locally

#### Option 1: Go
```bash
# Clone repository
git clone https://github.com/stianfro/file-upload-web.git
cd file-upload-web

# Run server
go run main.go

# Test upload via browser
open http://localhost:8080

# Test upload via curl
curl -X POST -F "file=@test.txt" http://localhost:8080/upload
```

#### Option 2: Docker
```bash
# Build image
docker build -t file-upload-web .

# Run container
docker run -p 8080:8080 -v $(pwd)/uploads:/uploads file-upload-web

# Test upload
curl -X POST -F "file=@test.txt" http://localhost:8080/upload
```

## Production Deployment

### Using Pre-built Image
```bash
# Pull from registry
docker pull ghcr.io/stianfro/file-upload-web:latest

# Run with custom settings
docker run -d \
  -p 80:8080 \
  -e MAX_SIZE=50 \
  -v /data/uploads:/uploads \
  ghcr.io/stianfro/file-upload-web:latest
```

### Kubernetes Deployment

#### Quick Install
```bash
# Add Helm repo
helm repo add file-upload oci://ghcr.io/stianfro
helm repo update

# Install chart
helm install my-upload file-upload/file-upload-web

# Port forward for testing
kubectl port-forward svc/my-upload-file-upload-web 8080:80
```

#### Custom Configuration
```yaml
# values.yaml
image:
  tag: latest

service:
  type: LoadBalancer
  port: 80

config:
  maxSize: "50"
  uploadDir: "/uploads"

persistence:
  enabled: true
  size: 10Gi
```

```bash
helm install my-upload file-upload/file-upload-web -f values.yaml
```

## Testing

### Browser Upload
1. Navigate to http://localhost:8080
2. Click "Choose File"
3. Select a file
4. Click "Upload"
5. See success message

### curl Upload
```bash
# Basic upload
curl -X POST -F "file=@image.jpg" http://localhost:8080/upload

# With progress bar
curl -X POST -F "file=@large.zip" http://localhost:8080/upload --progress-bar

# Multiple files (sequential)
for file in *.txt; do
  curl -X POST -F "file=@$file" http://localhost:8080/upload
done
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| PORT | 8080 | Server listen port |
| UPLOAD_DIR | /uploads | Storage directory |
| MAX_SIZE | 10 | Max file size (MB) |

## Troubleshooting

### Upload Fails: 413 Payload Too Large
- File exceeds MAX_SIZE limit
- Solution: Increase MAX_SIZE or use smaller file

### Upload Fails: 500 Internal Server Error
- Check UPLOAD_DIR exists and is writable
- Check disk space available
- Review container logs

### Cannot Access Service
- Verify port forwarding/exposure
- Check firewall rules
- Confirm service is running: `curl localhost:8080/health`

### Files Not Persisted
- Ensure volume is mounted
- Check UPLOAD_DIR matches mount point
- Verify persistence in Kubernetes deployment

## Security Notes

- No authentication by default
- Sanitizes filenames automatically
- No file execution permissions
- Consider network policies in production
- Add rate limiting if exposed publicly

## Next Steps

1. Add authentication if needed
2. Configure monitoring/alerting
3. Set up backup for uploaded files
4. Implement rate limiting
5. Add virus scanning (optional)