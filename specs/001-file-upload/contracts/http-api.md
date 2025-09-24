# HTTP API Contract

## Base URL
```
http://localhost:8080
```

## Endpoints

### GET /
**Purpose**: Serve HTML upload form
**Response**:
- Status: 200 OK
- Content-Type: text/html
- Body: HTML page with file upload form

### POST /upload
**Purpose**: Handle file upload
**Request**:
- Method: POST
- Content-Type: multipart/form-data
- Body: File data with field name "file"

**Response Success**:
- Status: 200 OK
- Content-Type: text/plain
- Body: "File uploaded successfully: {filename}"

**Response Errors**:
- 400 Bad Request: No file provided
- 413 Payload Too Large: File exceeds size limit
- 500 Internal Server Error: Storage failure

**curl Examples**:
```bash
# Upload with form data
curl -X POST -F "file=@test.txt" http://localhost:8080/upload

# Upload with raw POST (for testing)
curl -X POST \
  -H "Content-Type: application/octet-stream" \
  -H "X-Filename: test.txt" \
  --data-binary @test.txt \
  http://localhost:8080/upload
```

### GET /health
**Purpose**: Health check for Kubernetes probes
**Response**:
- Status: 200 OK
- Content-Type: text/plain
- Body: "OK"

## Configuration

Environment variables:
- `PORT`: Server port (default: 8080)
- `UPLOAD_DIR`: Upload directory (default: /uploads)
- `MAX_SIZE`: Maximum file size in MB (default: 10)

## File Storage

Files are stored in `UPLOAD_DIR` with naming pattern:
```
{timestamp}_{original_filename}
```

Example: `20250923_143022_document.pdf`

## Error Handling

All errors return plain text messages suitable for debugging:
- Clear description of what went wrong
- No sensitive information exposed
- Actionable feedback for troubleshooting