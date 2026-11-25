# Backend dengan Firestore

Backend ini adalah HTTP server Go yang terhubung ke Firestore (project Firebase). Dia membaca/tulis beberapa collection dan juga bisa menampilkan daftar collection yang ada di Firestore.

## Requirement

- Go (versi sesuai `go.mod`)
- Token akses Firebase/GCP via service account JSON
- project ID Firestore yang aktif

## Variabel lingkungan

| Variabel | Default | Keterangan |
|---|---|---|
| `GOOGLE_APPLICATION_CREDENTIALS` | `/home/ninofelino12/service-account.json` | Path ke service account JSON |
| `FIREBASE_PROJECT_ID` | `bizflow-dashboard-vlo72` | Project ID Firestore |
| `FIRESTORE_COLLECTION` | `items` | Collection default (dipakai hanya untuk collection `items`) |
| `PORT` | `8080` | Port HTTP server |

## Menjalankan

```bash
./run.sh
```

atau lakukan manual:

```bash
export PATH="/home/ninofelino12/go/bin:$PATH:$HOME/google-cloud-sdk/bin:$PATH"
export GOOGLE_APPLICATION_CREDENTIALS="${GOOGLE_APPLICATION_CREDENTIALS:-/home/ninofelino12/service-account.json}"
export FIREBASE_PROJECT_ID="${FIREBASE_PROJECT_ID:-bizflow-dashboard-vlo72}"
export FIRESTORE_COLLECTION="${FIRESTORE_COLLECTION:-items}"
export PORT="${PORT:-8080}"

go mod download
go build -o backend .
./backend
```

Setelah jalan, server akan listening di `:PORT`.

## Build di mesin lain

Salin seluruh project ini ke mesin lain, lalu jalankan:

```bash
./build.sh
```

Script ini akan:
1. memeriksa apakah Go terinstall
2. `go mod tidy`
3. `go mod download`
4. `go build -o backend .`

Setelah selesai, binary `backend` siap dijalankan dengan env yang tepat.

Jika ingin build tanpa script, perintah dasarnya:

```bash
go mod tidy
go mod download
go build -o backend .
```

## Endpoint

### Health

- `GET /health`
  - Mengembalikan status dan waktu server.

Contoh:
```bash
curl -fsS http://localhost:8080/health
```

### Collection `items`

- `GET /items`
  - Mengembalikan semua dokumen di collection `items`.
- `GET /items/{id}`
  - Mengembalikan satu dokumen berdasarkan ID.
- `POST /items`
  - Membuat dokumen baru.
  - Body JSON: `{"name": "..."}`
- `DELETE /items/{id}`
  - Menghapus dokumen berdasarkan ID.

Contoh:
```bash
curl -fsS http://localhost:8080/items
curl -fsS http://localhost:8080/items/<ID>
curl -fsS -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name":"contoh"}'
curl -fsS -X DELETE http://localhost:8080/items/<ID>
```

### Collection `device`

- `GET /devices`
  - Mengembalikan semua dokumen di collection `device`.
- `GET /devices/{id}`
  - Mengembalikan satu dokumen berdasarkan ID.
- `POST /devices`
  - Membuat dokumen baru.
  - Body JSON: `{"name": "..."}`
- `DELETE /devices/{id}`
  - Menghapus dokumen berdasarkan ID.

Contoh:
```bash
curl -fsS http://localhost:8080/devices
curl -fsS http://localhost:8080/devices/<ID>
curl -fsS -X POST http://localhost:8080/devices -H "Content-Type: application/json" -d '{"name":"device-1"}'
curl -fsS -X DELETE http://localhost:8080/devices/<ID>
```

### Collection `boottime`

- `GET /boottimes`
  - Mengembalikan semua dokumen di collection `boottime`.
- `GET /boottimes/{id}`
  - Mengembalikan satu dokumen berdasarkan ID.
- `POST /boottimes`
  - Membuat dokumen baru.
  - Body JSON: `{"boot_time": "..."}`
- `DELETE /boottimes/{id}`
  - Menghapus dokumen berdasarkan ID.

Contoh:
```bash
curl -fsS http://localhost:8080/boottimes
curl -fsS http://localhost:8080/boottimes/<ID>
curl -fsS -X POST http://localhost:8080/boottimes -H "Content-Type: application/json" -d '{"boot_time":"2026-11-25T19:00:00+07:00"}'
curl -fsS -X DELETE http://localhost:8080/boottimes/<ID>
```

### Daftar collection Firestore

- `GET /collections`
  - Mengembalikan daftar nama collection yang ada di Firestore project ini.

Contoh:
```bash
curl -fsS http://localhost:8080/collections
```

Output contoh:
```json
["device","items","mic"]
```

## Catatan

- Collection `items`, `device`, dan `boottime` masing-masing disimpan di collection Firestore dengan nama yang sama.
- Endpoint `/collections` menampilkan seluruh collection yang ada di Firestore, bukan hanya yang punya endpoint di backend ini.
- Jika collection belum ada di Firestore, nama itu tidak akan muncul di `/collections`.
