#!/bin/bash
# test.sh - Manual testing script for file upload web application

echo "=== File Upload Web Application Test Suite ==="
echo ""

# Check if server is running
echo "1. Testing health endpoint..."
curl -f http://localhost:8080/health || { echo "FAIL: Health check failed"; exit 1; }
echo " ✓ Health check passed"
echo ""

# Test file upload
echo "2. Testing file upload..."
echo "test content $(date)" > test.txt
curl -f -X POST -F "file=@test.txt" http://localhost:8080/upload || { echo "FAIL: Upload failed"; exit 1; }
echo " ✓ File upload passed"
echo ""

# Test HTML form endpoint
echo "3. Testing HTML form endpoint..."
curl -f -s http://localhost:8080/ | grep -q "<form" || { echo "FAIL: HTML form not found"; exit 1; }
echo " ✓ HTML form endpoint passed"
echo ""

# Test large file handling
echo "4. Testing large file (creating 5MB test file)..."
dd if=/dev/zero of=large_test.bin bs=1M count=5 2>/dev/null
curl -f -X POST -F "file=@large_test.bin" http://localhost:8080/upload || { echo "FAIL: Large file upload failed"; exit 1; }
echo " ✓ Large file upload passed"
rm -f large_test.bin
echo ""

# Test file with special characters
echo "5. Testing filename with special characters..."
echo "special content" > "test file with spaces.txt"
curl -f -X POST -F "file=@test file with spaces.txt" http://localhost:8080/upload || { echo "FAIL: Special filename upload failed"; exit 1; }
echo " ✓ Special filename upload passed"
rm -f "test file with spaces.txt"
echo ""

# Cleanup
rm -f test.txt

echo "=== All tests passed! ==="