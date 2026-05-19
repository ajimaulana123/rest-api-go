# API Examples (curl / bash)

Base URL: `http://localhost:8080`

## Health

```bash
curl -s http://localhost:8080/health
```

## Register

```bash
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

Simpan token dari response (satu baris, jangan dipotong):

```bash
TOKEN="<paste_token_here>"
# atau dari login/register langsung:
# TOKEN=$(curl -s -X POST .../login ... | jq -r '.token')
```

## Login

```bash
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

## Profile

```bash
# Pakai variabel $TOKEN — jangan tempel JWT panjang di beberapa baris (bisa 400 Bad Request)
curl -s http://localhost:8080/api/profile -H "Authorization: Bearer $TOKEN"
```

## Items CRUD

### List

```bash
curl -s http://localhost:8080/api/items \
  -H "Authorization: Bearer $TOKEN"
```

### Get by ID

```bash
curl -s http://localhost:8080/api/items/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Create

```bash
curl -s -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Belajar Go",
    "description": "REST API dengan Gin"
  }'
```

### Update

```bash
curl -s -X PUT http://localhost:8080/api/items/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Belajar Go (updated)",
    "description": "Sudah diperbarui"
  }'
```

### Delete

```bash
curl -s -X DELETE http://localhost:8080/api/items/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Jalankan semua sekaligus

```bash
chmod +x scripts/api-examples.sh
./scripts/api-examples.sh
```

Opsional: butuh `jq` untuk format JSON (`sudo apt install jq` / `brew install jq`).
