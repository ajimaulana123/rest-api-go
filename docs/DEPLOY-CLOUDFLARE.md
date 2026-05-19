# Deploy Gin API ke Cloudflare Containers

Arsitektur: **Worker** menerima request → diteruskan ke **Container** (binary Go/Gin).

## Prasyarat

1. [Docker Desktop](https://docs.docker.com/desktop/) — harus jalan saat deploy
2. [Node.js](https://nodejs.org/) 18+
3. Akun Cloudflare + `wrangler login`

## Setup sekali

```bash
npm install
wrangler login
```

## Environment variables

### Di `wrangler.toml` (non-secret)

Sudah di-set: `GIN_MODE`, `SUPABASE_REGION`, `SUPABASE_POOLER`, `SUPABASE_DB_MODE`

### Secrets (wajib)

```bash
wrangler secret put JWT_SECRET
wrangler secret put SUPABASE_PROJECT_ID
wrangler secret put SUPABASE_URL
wrangler secret put SUPABASE_ANON_KEY
wrangler secret put SUPABASE_DB_PASSWORD
```

Nilai sama seperti file `.env` lokal.

## Deploy

```bash
# Pastikan Docker running
docker info

npm run deploy
```

Deploy pertama **3–5 menit** (build image + provision container).

URL: `https://be-api.<subdomain>.workers.dev`

## Test

```bash
curl https://be-api.<subdomain>.workers.dev/health
```

## Perintah berguna

```bash
npx wrangler containers list
npx wrangler containers images list
npx wrangler tail          # logs Worker
```

## Catatan

- Gin **tidak** jalan di Workers runtime biasa — hanya via **Containers**
- Image harus **linux/amd64** (sudah di Dockerfile)
- Supabase pooler: pakai `aws-1` untuk project ini
- Free tier: container sleep setelah idle (`sleepAfter = 30m`)
