# Rewrite Plan: pos-desktop → pos-desktop-tauri

**Stack lama**: Electron + Express.js + Vue 3 + better-sqlite3  
**Stack baru**: Tauri 2 + Go sidecar (HTTP server) + Vue 3 + Go SQLite  
**Referensi proyek lama**: `../pos-desktop/`

---

## Arsitektur Target

```
Tauri shell (Rust — tipis, hanya jadi window manager)
  └── Go sidecar binary (pos-backend)
        ├── HTTP server :3001
        ├── SQLite via modernc.org/sqlite (pure Go, no CGO)
        └── mDNS discovery (grandcat/zeroconf)

Vue 3 renderer
  ├── Axios → localhost:3001 (mode master)
  ├── Axios → master_ip:3001 (mode client)
  └── window.__TAURI__ untuk IPC (config, device id, lisensi)
```

Go binary di-bundle sebagai Tauri **sidecar** — otomatis dijalankan saat app start, dimatikan saat app close.

---

## Step 1 — Setup Go Module

**Tujuan**: Inisialisasi Go project untuk backend.

```bash
mkdir backend
cd backend
go mod init pos-desktop-tauri/backend
```

**Struktur yang dibuat:**
```
backend/
  cmd/server/main.go     # entry point, start HTTP + SQLite
  internal/
    db/db.go             # koneksi SQLite + migration runner
    handler/             # setara controller Express
    model/               # query SQLite
    middleware/          # auth JWT, permission
    migration/           # DDL files
    router/router.go     # definisi semua route
  go.mod
  go.sum
```

**Dependencies yang dibutuhkan:**
```bash
go get github.com/go-chi/chi/v5          # HTTP router (mirip Express)
go get github.com/golang-jwt/jwt/v5      # JWT
go get golang.org/x/crypto               # bcrypt
go get modernc.org/sqlite                # SQLite pure Go (no CGO)
go get github.com/grandcat/zeroconf      # mDNS discovery
```

**Prompt Claude:**
```
Saya sedang membangun Go HTTP backend untuk aplikasi POS desktop (Tauri sidecar).

Stack:
- Chi router (go-chi/chi/v5)
- modernc.org/sqlite (pure Go, no CGO)
- golang-jwt/jwt/v5 untuk auth
- golang.org/x/crypto/bcrypt

Buatkan boilerplate berikut:
1. `backend/cmd/server/main.go` — start HTTP server di port yang bisa dikonfigurasi via env PORT (default 3001), init DB, register routes
2. `backend/internal/db/db.go` — singleton koneksi SQLite, method Init(dbPath string) dan Get() *sql.DB, jalankan migration runner saat Init
3. `backend/internal/router/router.go` — setup Chi router, register semua group route dengan middleware

Standar response JSON:
{ "success": true, "message": "...", "data": ... }
{ "success": false, "message": "..." }

Gunakan helper response di `backend/internal/helper/response.go` dengan fungsi:
- Success(w, data)
- Created(w, data)  
- NoContent(w)
- BadRequest(w, msg)
- NotFound(w, msg)
- ServerError(w, msg)
- Unauthorized(w, msg)
```

---

## Step 2 — Tauri Sidecar Integration

**Tujuan**: Bundle Go binary ke dalam Tauri app, otomatis start/stop.

### 2a. Build Go binary

Tambahkan ke `src-tauri/tauri.conf.json`:
```json
{
  "bundle": {
    "externalBin": ["binaries/pos-backend"]
  }
}
```

Script build Go untuk tiap platform (tambah ke `package.json`):
```json
{
  "scripts": {
    "build:backend:linux": "cd backend && GOOS=linux GOARCH=amd64 go build -o ../src-tauri/binaries/pos-backend-x86_64-unknown-linux-gnu ./cmd/server",
    "build:backend:windows": "cd backend && GOOS=windows GOARCH=amd64 go build -o ../src-tauri/binaries/pos-backend-x86_64-pc-windows-msvc.exe ./cmd/server",
    "build:backend": "npm run build:backend:linux && npm run build:backend:windows"
  }
}
```

### 2b. Tauri spawn sidecar di Rust

Edit `src-tauri/src/lib.rs`:

**Prompt Claude:**
```
Di Tauri 2 (Rust), saya ingin spawn sebuah sidecar binary bernama "pos-backend" saat app start.

Sidecar menerima argument: port (string), db_path (string)
Contoh: pos-backend 3001 /home/user/.local/share/pos/data.db

Kebutuhan:
1. Spawn sidecar saat app setup (tauri::Builder)
2. Kill sidecar saat app close (on_window_event WindowEvent::Destroyed atau app exit)
3. Kirim stdout/stderr sidecar ke log Tauri untuk debugging
4. Handle jika sidecar crash — coba restart sekali

Berikan implementasi di src-tauri/src/lib.rs dengan Tauri 2 API.
```

### 2c. Konfigurasi path DB

DB path ditentukan Rust, dikirim ke Go via argument:
- Linux/macOS: `~/.local/share/pos/data.db`
- Windows: `%APPDATA%\pos\data.db`

Gunakan `tauri::api::path::app_data_dir()` untuk resolve path.

---

## Step 3 — Port Migrations ke Go

**Referensi**: `../pos-desktop/src/backend/migrations/` (17 file JS)

**Pola migration Go:**
```go
// backend/internal/migration/migration.go
type Migration struct {
    Name string
    Up   func(db *sql.DB) error
}
```

**Prompt Claude:**
```
Port semua migration SQLite berikut dari JavaScript ke Go.

Pola Go yang dipakai:
```go
type Migration struct {
    Name string
    Up   func(db *sql.DB) error
}
```

Migration runner: cek tabel `migrations` (buat jika belum ada), jalankan yang belum dijalankan berdasarkan kolom `name`.

Berikut isi migration JS yang perlu diport:
[paste isi file migration JS dari pos-desktop satu per satu]

Catatan:
- SQLite datetime: gunakan `datetime('now')` (single quote)  
- Seed data (admin/admin123) bcrypt hash di Go pakai `golang.org/x/crypto/bcrypt`
- Urutan migration harus sama persis
```

---

## Step 4 — Port Auth

**Referensi**: 
- `../pos-desktop/src/backend/controllers/AuthController.js`
- `../pos-desktop/src/backend/middleware/auth.js`
- `../pos-desktop/src/backend/middleware/can.js`
- `../pos-desktop/src/backend/models/User.js`
- `../pos-desktop/src/backend/config/constants.js` (JWT_SECRET)

**Prompt Claude:**
```
Port sistem auth dari Express.js ke Go (Chi router).

Referensi Express:

AuthController.js:
[paste isi file]

middleware/auth.js:
[paste isi file]

middleware/can.js:
[paste isi file]

models/User.js:
[paste isi file]

Kebutuhan Go:
1. `POST /api/auth/login` — cek username/password (bcrypt), return JWT HS256
2. `GET /api/auth/me` — return user dari token
3. Middleware `AuthMiddleware` — parse Bearer token, inject user ke context
4. Middleware `CanMiddleware(permission string)` — cek permission dari role user

JWT claims: { user_id, username, role_id, permissions []string, exp }
JWT_SECRET dari environment variable JWT_SECRET

Gunakan:
- github.com/go-chi/chi/v5 untuk router
- github.com/golang-jwt/jwt/v5 untuk JWT
- golang.org/x/crypto/bcrypt untuk password

Buat file:
- backend/internal/handler/auth.go
- backend/internal/middleware/auth.go
- backend/internal/middleware/can.go
- backend/internal/model/user.go
```

---

## Step 5 — Port Resource: Kategori & Pemasok

**Referensi**:
- `../pos-desktop/src/backend/controllers/KategoriController.js`
- `../pos-desktop/src/backend/models/Kategori.js`
- `../pos-desktop/src/backend/controllers/PemasokController.js`
- `../pos-desktop/src/backend/models/Pemasok.js`

Simple CRUD dengan pagination. Port dua sekaligus karena strukturnya mirip.

**Prompt Claude:**
```
Port dua resource CRUD dari Express ke Go: Kategori dan Pemasok.

Referensi JS:
[paste KategoriController.js]
[paste Kategori.js model]
[paste PemasokController.js]
[paste Pemasok.js model]

Endpoint yang dibutuhkan:
- GET    /api/kategori          (pagination: page, pageSize, sortBy, sortOrder, search)
- POST   /api/kategori
- PUT    /api/kategori/:id
- DELETE /api/kategori/:id

- GET    /api/pemasok           (pagination: page, pageSize, sortBy, sortOrder, search)
- POST   /api/pemasok
- PUT    /api/pemasok/:id
- DELETE /api/pemasok/:id

Pola pagination response:
{ "success": true, "data": { "data": [...], "total": 123 } }

Whitelist sorting kolom wajib ada di model (jangan pakai sortBy mentah dari request).

Buat file:
- backend/internal/handler/kategori.go
- backend/internal/model/kategori.go
- backend/internal/handler/pemasok.go
- backend/internal/model/pemasok.go
```

---

## Step 6 — Port Resource: Produk (Paling Kompleks)

**Referensi**:
- `../pos-desktop/src/backend/controllers/ProdukController.js`
- `../pos-desktop/src/backend/models/Produk.js`

Model Produk pakai JOIN ke 6 tabel: kategori, pemasok, satuan_produk, barcode_produk, harga_produk, stok_produk, konsinyasi_produk. Soft delete via `deleted_at`.

**Endpoint:**
```
GET    /api/produk              pagination + search (nama/kode/barcode)
POST   /api/produk
GET    /api/produk/trashed      soft-deleted
GET    /api/produk/next-kode    generate kode berikutnya
GET    /api/produk/:id
PUT    /api/produk/:id
DELETE /api/produk/:id          soft delete (set deleted_at)
POST   /api/produk/:id/restore
GET    /api/produk/:id/satuan
POST   /api/produk/:id/satuan
PUT    /api/produk/:id/satuan/:satuanId
DELETE /api/produk/:id/satuan/:satuanId
```

**Prompt Claude:**
```
Port resource Produk dari Express ke Go. Ini yang paling kompleks karena:
1. BASE_SELECT melakukan JOIN ke 6 tabel
2. Ada soft delete (deleted_at)
3. Ada endpoint satuan (sub-resource)
4. Search bisa by nama_produk OR kode_produk OR kode_barcode (kode_barcode ada di tabel barcode_produk via JOIN)
5. Ada next-kode generator

Referensi JS:

ProdukController.js:
[paste isi file]

Produk.js (model):
[paste isi file]

Stack Go: Chi router, modernc.org/sqlite, standar database/sql

Penting:
- Search COUNT query juga harus JOIN ke barcode_produk untuk bisa filter by barcode
- SORT_COLUMNS whitelist wajib
- Soft delete: DELETE set deleted_at = datetime('now'), bukan hapus row
- next-kode: format PREFIX + 6 digit angka (A000001), cari max kode dengan prefix yang sama

Buat file:
- backend/internal/handler/produk.go
- backend/internal/model/produk.go
```

---

## Step 7 — Port Resource: User, Role, Permission, Toko

**Referensi**:
- `../pos-desktop/src/backend/controllers/UserController.js`
- `../pos-desktop/src/backend/controllers/RoleController.js`
- `../pos-desktop/src/backend/controllers/PermissionController.js`
- `../pos-desktop/src/backend/controllers/TokoController.js`
- Model masing-masing

**Prompt Claude:**
```
Port 4 resource dari Express ke Go: User, Role, Permission, Toko.

[paste masing-masing controller dan model JS]

Endpoint:
- GET/POST /api/users [can:pengaturan], PUT/DELETE /api/users/:id
- GET/POST /api/roles [can:pengaturan], PUT /api/roles/:id
- GET /api/permissions [auth]
- GET/PUT /api/toko [auth]

Catatan User:
- Password di-hash dengan bcrypt saat create/update
- Jangan return kolom password di response

Catatan Role:
- Relasi ke permissions via tabel pivot role_permissions
- Response role harus include array permissions

Buat semua handler dan model yang dibutuhkan.
```

---

## Step 8 — Port Vue Renderer

**Tujuan**: Copy seluruh frontend dari pos-desktop, sesuaikan yang perlu.

### 8a. Copy dependencies

Dari `../pos-desktop/package.json`, copy dependencies renderer ke `package.json` tauri:
```bash
# Dependencies yang perlu dicopy:
vue, vue-router, pinia
@tanstack/vue-table
@vueform/multiselect
lucide-vue-next
axios
jsbarcode
tailwindcss, daisyui
```

### 8b. Copy file renderer

```bash
# Copy seluruh struktur renderer
cp -r ../pos-desktop/src/renderer/src/* ./src/

# Yang perlu DIUBAH setelah copy:
# - src/main.js (sesuaikan import jika ada perbedaan)
# - src/services/api.js (baseURL logic sama, tidak perlu ubah)
# - src/router/index.js (guard lisensi perlu update ke Tauri API)
```

### 8c. Yang TIDAK perlu diubah
- Semua `views/` — tidak ada perubahan
- Semua `components/` — tidak ada perubahan  
- Semua `services/*Api.js` — Axios call ke localhost:3001 tetap sama
- `composables/`, `stores/`, `config/` — tidak ada perubahan
- `assets/main.css` — tidak ada perubahan

### 8d. Yang PERLU diubah

**`src/services/api.js`** — baseURL logic sama tapi cara baca config beda:
- Lama: `window.api.config.get('app_mode')` via IPC Electron
- Baru: `window.__TAURI__.core.invoke('get_config', { key: 'app_mode' })`

**`src/router/index.js`** — guard lisensi:
- Lama: `window.api.lisensi.verifyToken(token)`
- Baru: `window.__TAURI__.core.invoke('verify_license_token', { token })`

**Prompt Claude:**
```
Saya sedang port aplikasi Vue 3 dari Electron ke Tauri 2.

Di Electron, IPC dipanggil via preload bridge:
window.api.config.get(key)
window.api.config.set(key, value)
window.api.device.getId()
window.api.lisensi.verifyToken(token)
window.api.server.start()
window.api.discovery.advertise()
window.api.discovery.scan()
window.api.discovery.stopScan()
window.api.discovery.onFound(callback)

Di Tauri 2, IPC menggunakan:
import { invoke } from '@tauri-apps/api/core'
import { listen } from '@tauri-apps/api/event'

Berikut isi file yang perlu diubah:

src/services/api.js:
[paste isi file]

src/router/index.js:
[paste isi file]

Tolong port kedua file tersebut agar pakai Tauri 2 API.
Buat wrapper `src/services/tauriApi.js` yang expose fungsi-fungsi IPC di atas.
```

---

## Step 9 — Port IPC ke Tauri Commands (Rust)

**Tujuan**: Ganti `window.api.*` Electron ke Tauri `invoke()`.

Commands yang dibutuhkan di Rust (`src-tauri/src/lib.rs`):

```rust
// Config
#[tauri::command] get_config(key: String) -> Result<String>
#[tauri::command] set_config(key: String, value: String) -> Result<()>

// Device
#[tauri::command] get_device_id() -> Result<String>

// Lisensi  
#[tauri::command] verify_license_token(token: String) -> Result<serde_json::Value>
#[tauri::command] open_browser(url: String) -> Result<()>
```

Discovery (mDNS) dihandle Go, diexpose via HTTP endpoint ke Vue atau via Tauri event.

**Prompt Claude:**
```
Implementasikan Tauri 2 commands berikut di src-tauri/src/lib.rs:

1. `get_config(key: String)` — baca dari file JSON di app_data_dir/config.json
2. `set_config(key: String, value: serde_json::Value)` — tulis ke config.json
3. `get_device_id()` — generate SHA256 dari MAC address + hostname, simpan permanen ke config.json
4. `verify_license_token(token: String)` — decode JWT RS256 pakai public key (hardcoded), return claims sebagai JSON. Verifikasi device_id di claims cocok dengan device ini.
5. `open_browser(url: String)` — buka URL di browser default

Untuk config.json: path di app_data_dir() / "config.json"
Untuk device_id: hash SHA256(MAC_ADDRESS + HOSTNAME)

Dependencies Rust yang mungkin dibutuhkan:
- serde, serde_json (sudah ada)
- sha2 (untuk SHA256)
- jsonwebtoken (untuk JWT RS256)
- mac_address
- hostname

Berikan implementasi lengkap termasuk Cargo.toml dependencies.
```

---

## Step 10 — Port mDNS Discovery ke Go

**Referensi**: `../pos-desktop/src/main/electron/DiscoveryService.js`

mDNS dihandle Go sidecar, expose via HTTP endpoint:

```
POST /api/internal/discovery/advertise    # mode master: announce ke network
GET  /api/internal/discovery/scan         # mode client: scan, return list master
POST /api/internal/discovery/stop
```

Vue memanggil endpoint ini via Axios biasa (bukan IPC).

**Prompt Claude:**
```
Port DiscoveryService dari Electron/Node.js ke Go menggunakan grandcat/zeroconf.

Referensi JS:
[paste DiscoveryService.js]

Kebutuhan:
1. `POST /api/internal/discovery/advertise` — start mDNS advertise service "_pos._tcp" dengan payload: { name, port, ip }
2. `GET /api/internal/discovery/scan` — scan jaringan untuk service "_pos._tcp", return list dalam 3 detik
3. `POST /api/internal/discovery/stop` — stop advertise/scan

Buat:
- backend/internal/handler/discovery.go
- backend/internal/service/discovery.go (logic mDNS)
```

---

## Step 11 — Port License System ke Go + Rust

**Referensi**: `../pos-desktop/src/main/electron/LisensiService.js`

License verification split:
- **Validasi JWT offline** → Tauri Rust command `verify_license_token` (Step 9)
- **Aktivasi/validasi online** → Go handler yang forward ke license server

```
POST /api/internal/lisensi/aktivasi    # forward ke http://lisansi.test/api/lisensi/aktivasi
POST /api/internal/lisensi/validasi    # forward ke http://lisansi.test/api/lisensi/validasi
POST /api/internal/lisensi/deaktivasi
```

**Prompt Claude:**
```
Port LisensiService dari Electron ke Go. Service ini adalah proxy/wrapper ke license server remote.

Referensi JS:
[paste LisensiService.js]

Go handler di backend/internal/handler/lisensi.go:
1. POST /api/internal/lisensi/aktivasi — forward request ke license server, return response
2. POST /api/internal/lisensi/validasi — sama
3. POST /api/internal/lisensi/deaktivasi — sama

License server base URL: dari env variable LICENSE_SERVER_URL (default: http://lisansi.test/api)

Gunakan net/http standard library untuk HTTP client.
```

---

## Step 12 — Setup Vite + Tailwind + DaisyUI

**Prompt Claude:**
```
Setup Vite config untuk Tauri 2 + Vue 3 + Tailwind CSS 4 + DaisyUI 5.

Requirement:
- Dev server di port 1420 (Tauri default)
- Tailwind CSS 4 via @tailwindcss/vite plugin
- DaisyUI 5 via @plugin di CSS
- @vueform/multiselect
- Path alias: @ → ./src

Berikan:
1. vite.config.js yang lengkap
2. src/assets/main.css dengan import Tailwind + DaisyUI + multiselect overrides
3. package.json scripts untuk dev dan build

Referensi CSS lama (copy dari pos-desktop):
[paste src/assets/main.css dari pos-desktop]
```

---

## Step 13 — Testing Multi-Client

**Checklist sebelum production:**

- [ ] Master start: Go sidecar jalan di port 3001
- [ ] Master advertise mDNS di LAN
- [ ] Client scan dan temukan master
- [ ] Client connect ke master IP:3001
- [ ] Login dari client berhasil
- [ ] CRUD produk dari client tersimpan di master DB
- [ ] 3 client simultan tidak ada race condition di SQLite
- [ ] App close: Go sidecar mati bersih
- [ ] Restart app: sidecar spawn ulang tanpa error
- [ ] Offline mode: license token masih valid (JWT local verification)

---

## Step 14 — Build & Distribusi

```bash
# 1. Build Go binary
npm run build:backend

# 2. Build Tauri app
npm run tauri build
```

Output:
- Linux: `src-tauri/target/release/bundle/appimage/*.AppImage`
- Windows: `src-tauri/target/release/bundle/nsis/*.exe`

**Ukuran estimasi:**
- Installer: ~8–15MB (vs Electron ~150MB)
- RAM idle: ~35MB (vs Electron ~200MB)

---

## Urutan Pengerjaan

```
[ ] Step 1  — Setup Go module + dependencies
[ ] Step 2  — Tauri sidecar integration (spawn/kill Go binary)
[ ] Step 3  — Port migrations ke Go
[ ] Step 4  — Port auth (login, JWT middleware, permission)
[ ] Step 5  — Port Kategori + Pemasok (CRUD sederhana)
[ ] Step 6  — Port Produk (kompleks, JOIN + soft delete + satuan)
[ ] Step 7  — Port User, Role, Permission, Toko
[ ] Step 8  — Copy + adjust Vue renderer
[ ] Step 9  — Port IPC ke Tauri commands (Rust)
[ ] Step 10 — Port mDNS discovery ke Go
[ ] Step 11 — Port license system
[ ] Step 12 — Setup Vite + Tailwind + DaisyUI
[ ] Step 13 — Testing multi-client
[ ] Step 14 — Build + distribusi
```

---

## Catatan

- **SQLite concurrent write**: Go default SQLite tidak support concurrent write. Set `_journal_mode=WAL` dan `_busy_timeout=5000` di connection string untuk multi-client.
  ```go
  db, _ := sql.Open("sqlite", "file:data.db?_journal_mode=WAL&_busy_timeout=5000")
  ```
- **Port Go binary**: nama sidecar binary harus match triple Tauri (`x86_64-unknown-linux-gnu`, `x86_64-pc-windows-msvc`, dll)
- **License server URL**: hardcoded di Step 11, nanti pindah ke config saat production
- **Migration berikutnya**: mulai dari `018_*` (lanjut dari pos-desktop)
