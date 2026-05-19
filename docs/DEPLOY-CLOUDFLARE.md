# Deploy Gin API ke Cloudflare Containers

Arsitektur: **Worker** (`src/index.ts`) → **Container** (Docker / Gin :8080) → **Supabase PostgreSQL**.

Nama project di Dashboard dan `wrangler.toml` harus **sama**: `rest-api-go`.

---

## PENTING — Keamanan (baca dulu)

Jika Anda pernah menaruh secret sebagai **plain text Variables** di Dashboard (bukan **Encrypted Secrets**), nilai tersebut bisa **terbaca di log build**. Lakukan segera:

1. **Rotate** semua credential yang pernah dimasukkan sebagai plain var:
   - `JWT_SECRET` (ganti string baru)
   - `SUPABASE_DB_PASSWORD` (reset di Supabase → Database)
   - `SUPABASE_ANON_KEY` (rotate di Supabase → API, jika perlu)
   - **Build token** Cloudflare (revoke token lama, buat baru)
2. **Hapus** dari Worker **Variables** (plain):
   - `CLOUDFLARE_API_TOKEN` — **tidak boleh** di Worker vars
   - `CLOUDFLARE_ACCOUNT_ID` — tidak perlu di Worker runtime
   - Semua secret aplikasi yang sempat di-set sebagai plain text
3. Pasang ulang hanya sebagai **Encrypted Secrets** (lihat bawah).

| Variable | Di mana | Tipe |
|----------|---------|------|
| `CLOUDFLARE_API_TOKEN` | **Build** settings saja | Encrypted |
| `CLOUDFLARE_ACCOUNT_ID` | **Build** settings saja | Plain OK |
| `JWT_SECRET`, `SUPABASE_*` | **Worker** → Variables and Secrets | **Secret** (encrypted) |
| `GIN_MODE`, `SUPABASE_REGION`, dll. | `wrangler.toml` `[vars]` | Non-secret |

---

## Prasyarat

- Akun Cloudflare dengan akses **Workers + Containers**
- Repo terhubung ke **Workers Builds**
- Build environment punya **Docker** (untuk build image)

---

## 1. API Token untuk Build (fix `Unauthorized` setelah Docker build)

Log tipikal: Worker ter-upload, image Docker **berhasil di-build**, lalu gagal:

```text
✘ [ERROR] Unauthorized
```

Itu biasanya gagal **push image** ke registry Cloudflare — token build kurang permission.

### Buat token baru

1. [API Tokens](https://dash.cloudflare.com/profile/api-tokens) → **Create Token**
2. **Custom token** dengan permission minimal:

| Resource | Permission |
|----------|------------|
| Account → **Workers Scripts** | Edit |
| Account → **Workers Containers** | Edit |
| Account → **Account Settings** | Read |

   (Template "Edit Cloudflare Workers" kadang kurang untuk Containers — pakai custom jika masih `Unauthorized`.)

3. Salin token (sekali tampil).

### Pasang hanya di Build

**Workers & Pages** → **rest-api-go** → **Settings** → **Build** → **Environment variables**:

| Name | Value | Encrypt |
|------|-------|---------|
| `CLOUDFLARE_API_TOKEN` | token baru | Ya |
| `CLOUDFLARE_ACCOUNT_ID` | Account ID (32 char) | Tidak |

**Jangan** tambahkan `CLOUDFLARE_API_TOKEN` di Worker **Variables and Secrets** untuk runtime.

### Build command

```bash
npm install && npx wrangler deploy
```

---

## 2. Secret aplikasi (Worker → Container)

Secret harus ada di Worker agar `src/index.ts` bisa meneruskannya ke container Gin.

### Via Dashboard (disarankan setelah leak)

**rest-api-go** → **Settings** → **Variables and Secrets** → **Add**:

| Name | Type |
|------|------|
| `JWT_SECRET` | Secret |
| `SUPABASE_PROJECT_ID` | Secret |
| `SUPABASE_URL` | Secret |
| `SUPABASE_ANON_KEY` | Secret |
| `SUPABASE_DB_PASSWORD` | Secret |

### Via CLI (lokal, sudah `wrangler login`)

```bash
npx wrangler secret put JWT_SECRET --name rest-api-go
npx wrangler secret put SUPABASE_PROJECT_ID --name rest-api-go
npx wrangler secret put SUPABASE_URL --name rest-api-go
npx wrangler secret put SUPABASE_ANON_KEY --name rest-api-go
npx wrangler secret put SUPABASE_DB_PASSWORD --name rest-api-go
```

Non-secret tetap dari `wrangler.toml`: `GIN_MODE`, `SUPABASE_REGION`, `SUPABASE_POOLER`, `SUPABASE_DB_MODE`.

---

## 3. Deploy lokal (opsional)

```bash
npm install
npx wrangler login
docker info          # wajib untuk build image
npx wrangler deploy
```

Deploy pertama bisa **3–5 menit** (provision container).

URL: `https://rest-api-go.<subdomain>.workers.dev`

```bash
curl https://rest-api-go.<subdomain>.workers.dev/health
```

---

## Troubleshooting

| Gejala | Penyebab | Solusi |
|--------|----------|--------|
| Worker name mismatch | `wrangler.toml` ≠ Dashboard | Pakai `name = "rest-api-go"` (sudah diset) |
| `Unauthorized` setelah Docker build | Token tanpa **Workers Containers Edit** | Token custom + pasang di **Build** env |
| Build token deleted/rolled | Token di-revoke | Buat token baru di Build settings |
| API jalan tapi DB/auth gagal | Secret Worker hilang setelah deploy override | Set ulang **Encrypted Secrets** di Dashboard |
| Secret di log build | Plain **Variables** bukan **Secrets** | Rotate + pindah ke Secret + hapus plain vars |

---

## Perintah berguna

```bash
npx wrangler containers list
npx wrangler containers images list
npx wrangler tail
```

## Catatan teknis

- Gin hanya jalan di **Containers**, bukan Workers runtime biasa.
- Image: `linux/amd64` (lihat `Dockerfile`).
- Supabase pooler: `SUPABASE_POOLER=aws-1`, region `ap-southeast-1`.
- `AUTO_MIGRATE` tidak dipakai di production; schema di Supabase sudah ada.
