#!/usr/bin/env bash
# Contoh curl untuk REST API (Gin)
# Jalankan server dulu: go run main.go

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
TOKEN=""

echo "==> Health"
curl -s "$BASE_URL/health" | jq .

echo ""
echo "==> Register"
REGISTER_RESP=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }')
echo "$REGISTER_RESP" | jq .
TOKEN=$(echo "$REGISTER_RESP" | jq -r '.token')

# Jika email sudah terdaftar, gunakan login:
# echo ""
# echo "==> Login"
# LOGIN_RESP=$(curl -s -X POST "$BASE_URL/api/auth/login" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "email": "john@example.com",
#     "password": "secret123"
#   }')
# echo "$LOGIN_RESP" | jq .
# TOKEN=$(echo "$LOGIN_RESP" | jq -r '.token')

AUTH_HEADER="Authorization: Bearer $TOKEN"

echo ""
echo "==> Profile"
curl -s "$BASE_URL/api/profile" -H "$AUTH_HEADER" | jq .

echo ""
echo "==> Create item"
CREATE_RESP=$(curl -s -X POST "$BASE_URL/api/items" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d '{
    "title": "Belajar Go",
    "description": "REST API dengan Gin"
  }')
echo "$CREATE_RESP" | jq .
ITEM_ID=$(echo "$CREATE_RESP" | jq -r '.data.id')

echo ""
echo "==> List items"
curl -s "$BASE_URL/api/items" -H "$AUTH_HEADER" | jq .

echo ""
echo "==> Get item by ID"
curl -s "$BASE_URL/api/items/$ITEM_ID" -H "$AUTH_HEADER" | jq .

echo ""
echo "==> Update item"
curl -s -X PUT "$BASE_URL/api/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d '{
    "title": "Belajar Go (updated)",
    "description": "Sudah diperbarui"
  }' | jq .

echo ""
echo "==> Delete item"
curl -s -X DELETE "$BASE_URL/api/items/$ITEM_ID" -H "$AUTH_HEADER" | jq .
