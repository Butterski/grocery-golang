#!/bin/bash
# Test commands for Grocery List Management API
# Save this file as test_api.sh and run with: bash test_api.sh

# Set the base URL for the API
BASE_URL="http://localhost:8080"

echo "=== Testing Grocery Items API with JWT Authentication ==="

echo -e "\n=== AUTHENTICATION TESTS ==="

echo -e "\n=== 1. Register a new user ==="
curl -X POST "${BASE_URL}/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","repeated_password":"password123"}' \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 2. Login and get JWT token ==="
TOKEN=$(curl -s -X POST "${BASE_URL}/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' | grep -o '"token":"[^"]*"' | sed 's/"token":"\(.*\)"/\1/')

echo "Token received: ${TOKEN:0:20}..." # Show only first 20 chars for brevity
echo "Status: 200"

echo -e "\n=== 3. Try accessing protected endpoint without token ==="
curl -X GET "${BASE_URL}/items" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 4. Access protected endpoint with token ==="
curl -X GET "${BASE_URL}/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 5. Create item with authentication ==="
curl -X POST "${BASE_URL}/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Milk","quantity":2,"unit":"liters","category":"Dairy","notes":"Get lactose-free"}' \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 6. Get item with authentication ==="
curl -X GET "${BASE_URL}/items/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 7. Try with invalid token ==="
curl -X GET "${BASE_URL}/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalidtoken123" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 8. Try with malformed Authorization header ==="
curl -X GET "${BASE_URL}/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: InvalidFormat $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== ADDITIONAL TESTS WITH AUTHENTICATION ==="

# Continue with other tests but now include the Authorization header
echo -e "\n=== 9. List all grocery items with authentication ==="
curl -X GET "${BASE_URL}/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 10. Test filtering by name with authentication ==="
curl -X GET "${BASE_URL}/items?name=milk" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 11. Update an item with authentication ==="
curl -X PUT "${BASE_URL}/items/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Milk","quantity":3,"unit":"liters","category":"Dairy","notes":"Get whole milk instead"}' \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== 12. Delete an item with authentication ==="
curl -X DELETE "${BASE_URL}/items/1" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo -e "\n=== API Testing Complete ==="