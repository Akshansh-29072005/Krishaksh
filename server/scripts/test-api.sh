#!/bin/bash
# Test script for API endpoints after deployment

echo "🧪 Testing Krisho API endpoints..."

API_URL="https://api-krisho.aarcsx.com"

echo "1. Testing health endpoint..."
curl -s -o /dev/null -w "Status: %{http_code}\n" "$API_URL/health"

echo "2. Testing ready endpoint..."
curl -s -o /dev/null -w "Status: %{http_code}\n" "$API_URL/ready"

echo "3. Testing weather endpoint..."
curl -s "$API_URL/api/v1/weather?lat=28.6139&lon=77.2090" | head -c 200
echo "..."

echo "4. Testing metrics endpoint..."
curl -s "$API_URL/metrics" | head -10

echo "✅ API testing complete!"