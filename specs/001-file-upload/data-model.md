# Data Model

## File Upload Entity

### Stored File
```
Structure:
- Filename: string (sanitized original name)
- Timestamp: ISO 8601 format
- Size: bytes (integer)
- MIME Type: string (from Content-Type header)
- Storage Path: absolute path on filesystem

Example:
{
  "filename": "test.pdf",
  "timestamp": "2025-09-23T14:30:22Z",
  "size": 102400,
  "mime_type": "application/pdf",
  "storage_path": "/uploads/20250923_143022_test.pdf"
}
```

### Upload Request
```
Structure:
- File: binary data
- Filename: string (from form field or header)
- Content-Type: MIME type
- Content-Length: size in bytes

Validation:
- Filename: alphanumeric, dots, dashes, underscores only
- Size: must not exceed MAX_SIZE
- Type: no restrictions (accept all)
```

### Upload Response
```
Success:
- Status: 200
- Message: confirmation with filename
- Timestamp: when stored

Error:
- Status: 4xx or 5xx
- Error: descriptive message
- Details: troubleshooting hints
```

## Storage Schema

### Directory Structure
```
/uploads/
├── 20250923_143022_document.pdf
├── 20250923_143145_image.jpg
├── 20250923_143234_data.csv
└── .gitkeep
```

### Naming Convention
- Prefix: YYYYMMDD_HHMMSS
- Separator: underscore
- Suffix: sanitized original filename
- Purpose: Avoid collisions, maintain order

### File Metadata
No database or metadata files. Information derived from:
- Filesystem attributes (size, modified time)
- Filename parsing (timestamp, original name)
- MIME type detection (optional, from upload)

## Constraints

### Size Limits
- Default: 10 MB per file
- Configurable via MAX_SIZE
- Enforced before writing to disk

### Filename Sanitization
- Remove path components (../, ./, /)
- Replace spaces with underscores
- Remove special characters except .-_
- Truncate to 255 characters max

### Concurrency
- Timestamp precision ensures uniqueness
- No locking required (append-only)
- Each upload independent transaction

## Non-Persistent Data

The application maintains no state between requests:
- No upload history
- No user sessions
- No file indices
- Pure request/response model

This ensures:
- Horizontal scalability
- Simple disaster recovery
- Minimal memory footprint
- Easy debugging