# Simple Gateway

**Simple Gateway** adalah API Gateway ringan berbasis Go yang dirancang untuk **self-hosted**, **internal tools**, dan **indie/startup kecil** yang membutuhkan kontrol routing dan security dasar **tanpa kompleksitas berlebihan** seperti Kong atau Envoy.

Project ini berfokus pada:
- konfigurasi **sederhana (YAML)**
- **satu pintu masuk** (single entry point)
- **mudah dipahami & dioperasikan**
- siap untuk **production awal**

---

## ✨ Fitur Utama (v1)

- Host & path based routing
- Reverse proxy ke backend (HTTP)
- API Key Authentication (per-route)
- Multi subdomain dalam satu gateway
- Konfigurasi berbasis YAML
- Environment variable untuk secret
- Bisa dijalankan sebagai **systemd service**
- Cocok untuk self-hosted / on-prem

---

## 🧱 Arsitektur Singkat

Client / Browser / API Client
↓
Simple Gateway (Go)
↓
Backend Service
(Nginx / API / App)


Semua request **harus melewati gateway**.  
Routing dan security dikontrol di satu tempat.

---

## 📁 Struktur Project
simple-gateway/
├── cmd/
│ └── gateway/
│ └── main.go
├── internal/
│ ├── auth/
│ │ └── apikey.go
│ ├── proxy/
│ ├── router/
│ │ └── router.go
│ └── server/
│ └── http.go
├── config/
│ └── gateway.yaml.example
├── go.mod
├── go.sum
├── README.md
└── .gitignore


---

## ⚙️ Konfigurasi (`gateway.yaml`)

Contoh konfigurasi:

```yaml
listen: :8080

apis:
  - host: server.example.com
    path: /
    forward_to: http://localhost:80
    protect:
      api_key: true

  - host: api.example.com
    path: /
    forward_to: http://localhost:9000
    protect:
      api_key: true

| Field             | Deskripsi                          |
| ----------------- | ---------------------------------- |
| `listen`          | Port gateway                       |
| `host`            | Domain / subdomain yang dicocokkan |
| `path`            | Path rule (`/` berarti semua path) |
| `forward_to`      | Backend tujuan                     |
| `protect.api_key` | Aktifkan API key authentication    |

⚠️ Secret TIDAK disimpan di YAML

🔐 API Key Authentication
Set API key

API key diambil dari environment variable:
export GATEWAY_API_KEY=super-secret-key

Request harus membawa header:
X-API-Key: super-secret-key

| Kondisi           | Response                      |
| ----------------- | ----------------------------- |
| Tidak ada API key | 401 Unauthorized              |
| API key salah     | 401 Unauthorized              |
| API key benar     | Request diteruskan ke backend |

🚀 Menjalankan Gateway
Mode Development
go run ./cmd/gateway -config config/gateway.yaml.example

Build Binary
go build -o gateway ./cmd/gateway

🖥️ Menjalankan sebagai systemd service (Production)
Lokasi standar
Binary: /usr/local/bin/gateway
Config: /etc/simple-gateway/gateway.yaml

Contoh service file
/etc/systemd/system/simple-gateway.service

[Unit]
Description=Simple Gateway
After=network.target

[Service]
ExecStart=/usr/local/bin/gateway -config /etc/simple-gateway/gateway.yaml
Environment=GATEWAY_API_KEY=super-secret-key
Restart=always
User=www-data
Group=www-data

[Install]
WantedBy=multi-user.target

Aktifkan service:
sudo systemctl daemon-reload
sudo systemctl enable simple-gateway
sudo systemctl start simple-gateway

Cek status dan log:
sudo systemctl status simple-gateway
journalctl -u simple-gateway -f

🧪 Testing
Tanpa API key (harus gagal)
curl -H "Host: api.example.com" http://localhost:8080

Dengan API key (harus tembus)
curl -H "Host: api.example.com" \
     -H "X-API-Key: super-secret-key" \
     http://localhost:8080

🆚 Kenapa Simple Gateway?
Simple Gateway bukan pengganti Kong atau Envoy.

Project ini dibuat untuk:
developer yang tidak butuh fitur enterprise
ingin kontrol penuh
setup cepat
mudah dipahami
mudah dimodifikasi

🛣️ Roadmap
Rate limiting (429)
Basic authentication
Error override dari YAML
Graceful shutdown
Reload config tanpa restart
Metrics (Prometheus)
