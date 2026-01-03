# Simple Gateway

**Simple Gateway** adalah API Gateway ringan berbasis Go yang dibuat untuk **self-hosted** dan **internal services**—bukan untuk enterprise-scale systems.  
Project ini lahir dari kebutuhan praktis: ketika nginx config mulai berulang dan sulit dirawat, tapi solusi seperti Kong atau Spring Cloud Gateway terasa terlalu berat.

Fokus utamanya adalah **simplicity, clarity, dan low operational overhead**.

---

## Kenapa Simple Gateway?

Awalnya menggunakan **nginx langsung**, setiap menambah service harus menambah konfigurasi baru:
routing, auth, rate limit, dan aturan lainnya.  
Saat jumlah service bertambah, konfigurasi menjadi panjang, duplikatif, dan sulit dirawat.

Solusi seperti **Kong** atau **Spring Cloud Gateway** memang sangat powerful, tapi sering kali **overkill** untuk kebutuhan self-hosted atau internal tools.

Simple Gateway mengambil pendekatan berbeda:
- satu binary
- satu file konfigurasi
- tanpa framework
- tanpa setup berlapis

---

## Target Use Case

Simple Gateway **bukan** untuk:
- enterprise dengan ribuan API
- IAM kompleks (OAuth2 variants, RBAC enterprise)
- multi-team, multi-tenant gateway
- cloud-native platform besar

Simple Gateway **cocok untuk**:
- self-hosted environments
- internal services
- small–medium scale systems
- developer atau startup kecil
- kebutuhan gateway yang sederhana tapi rapi

---

## Fitur Utama

- Host & path based routing
- Reverse proxy ke backend service
- API Key authentication (per route)
- Multi subdomain dalam satu gateway
- Konfigurasi berbasis YAML
- Environment variable untuk secret
- Bisa dijalankan sebagai systemd service
- Mudah di-versioning (Git-friendly)

---

## Arsitektur Singkat

Client / Browser / API Client  
↓  
Simple Gateway  
↓  
Backend Services (Nginx / API / App)

---

## Struktur Project

simple-gateway/
├── cmd/gateway/main.go  
├── internal/  
├── config/gateway.yaml.example  
├── go.mod  
├── go.sum  
├── README.md  
└── .gitignore  

---

## Konfigurasi (gateway.yaml)

Contoh konfigurasi:

listen: :8080

apis:
- host: server.example.com
  path: /
  forward_to: http://localhost:80
  protect:
    api_key: true

---

## API Key Authentication

API key diambil dari environment variable:

export GATEWAY_API_KEY=super-secret-key

Header request:

X-API-Key: super-secret-key

---

## Menjalankan Gateway

Development:

go run ./cmd/gateway -config config/gateway.yaml.example

Build:

go build -o gateway ./cmd/gateway

---

## systemd Service

Binary: /usr/local/bin/gateway  
Config: /etc/simple-gateway/gateway.yaml

---

## Filosofi Desain

Project ini sengaja **tidak mencoba menyelesaikan semua masalah**.  
Jika kebutuhan sudah enterprise-scale, gunakan tool enterprise.

---

## Roadmap

- Rate limiting
- Basic auth
- Custom error handling
- Graceful shutdown
- Config reload

---

## Lisensi

MIT License

---

## Author

Andreas Nainggolan
