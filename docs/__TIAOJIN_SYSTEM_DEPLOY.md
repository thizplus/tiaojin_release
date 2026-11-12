# Tiaojin System — Deploy & Development Guide

## ข้อมูลโปรเจค

| รายการ | ค่า |
|--------|-----|
| ชื่อ | เที่ยวจีนแอพ (Tiaojin System) |
| Frontend | beta.tiaojin.co.th |
| API | api.tiaojin.co.th |
| Server | 5.223.54.168 (alias: `alljs`) |
| Path | `/opt/tiaojin-system/` |
| Stack | Go Fiber + Vite (React) + PostgreSQL 17 + Redis 7 |

---

## SSH เข้า Server

### Server Info
- **IP**: 5.223.54.168
- **Port**: 22
- **User**: root
- **SSH Key**: `id_ed25519_hetzner` (ไฟล์อยู่ที่ `C:\Users\Admin\.ssh\id_ed25519_hetzner`)
- **Public Key**: `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEmhiBXGgGI4OVRK9HgQ+DHuskVOXAEMz1yys6aIsb8z admin@fraudchecker`

### เครื่อง Admin (มี SSH config แล้ว)
```bash
ssh alljs
cd /opt/tiaojin-system
```

### คนอื่นต้องการเข้า (เครื่องใหม่)

#### วิธีที่ 1: ใช้ key เดียวกัน
1. Copy ไฟล์ `id_ed25519_hetzner` (private key) ไปวางที่ `~/.ssh/` ของเครื่องใหม่
2. เพิ่มใน `~/.ssh/config`:
```
Host alljs
    HostName 5.223.54.168
    User root
    IdentityFile ~/.ssh/id_ed25519_hetzner
```
3. ทดสอบ: `ssh alljs`

#### วิธีที่ 2: สร้าง key ใหม่สำหรับคนนั้น (แนะนำ — ปลอดภัยกว่า)
1. สร้าง key ที่เครื่องคนนั้น:
```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519_alljs -N "" -C "ชื่อคน@เครื่อง"
```
2. Copy public key ที่ได้ (ไฟล์ `.pub`) ส่งให้ Admin
3. Admin SSH เข้า server แล้วเพิ่ม key:
```bash
ssh alljs
echo 'วาง public key ของคนนั้นตรงนี้' >> ~/.ssh/authorized_keys
```
4. คนนั้นเพิ่ม SSH config:
```
Host alljs
    HostName 5.223.54.168
    User root
    IdentityFile ~/.ssh/id_ed25519_alljs
```
5. ทดสอบ: `ssh alljs`

---

## โครงสร้างโปรเจค

```
/opt/tiaojin-system/
├── docker-compose.prod.yml          # Docker Compose หลัก
├── gofiber/code/                    # Backend (Go Fiber)
│   ├── cmd/api/                     # Entry point
│   ├── domain/
│   │   ├── models/                  # Entity structs
│   │   ├── dto/                     # Request/Response DTOs
│   │   ├── repositories/           # Repository interfaces
│   │   ├── services/               # Service interfaces
│   │   └── ports/                  # Port interfaces
│   ├── application/serviceimpl/    # Service implementations
│   ├── infrastructure/
│   │   ├── postgres/               # Repository implementations
│   │   ├── redis/                  # Cache
│   │   ├── storage/                # S3/CDN (Cloudflare R2)
│   │   ├── websocket/             # WebSocket
│   │   ├── email/                 # Email service
│   │   ├── slipverification/      # Slip verification
│   │   └── services/              # External services
│   ├── interfaces/api/             # Handlers + Routes + Middleware
│   ├── pkg/
│   │   ├── config/                # Configuration
│   │   ├── di/                    # Dependency injection
│   │   ├── utils/                 # Utilities
│   │   ├── credentials/          # Google service account key
│   │   ├── cdn/                  # CDN helpers
│   │   ├── errors/               # Error handling
│   │   └── scheduler/            # Cron jobs
│   ├── migration/                 # DB migrations
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   └── Makefile
│
└── vite/code/                      # Frontend (React + Vite)
    ├── src/
    │   ├── page/                   # Pages (admin, agent, sale, super, globals)
    │   ├── components/             # UI Components
    │   │   ├── ui/                # shadcn/ui components
    │   │   ├── admin/             # Admin components
    │   │   ├── agent/             # Agent components
    │   │   ├── public/            # Public components
    │   │   ├── manage/            # Management components
    │   │   └── globals/           # Shared components
    │   ├── services/              # API calls (admin, agent, public)
    │   ├── hooks/                 # React hooks
    │   ├── stores/                # Zustand stores
    │   ├── types/                 # TypeScript types (models, request, response)
    │   ├── constants/             # Constants & enums
    │   ├── routes/                # React Router config
    │   ├── contexts/              # React contexts
    │   ├── layouts/               # Layout components
    │   ├── theme/                 # Theme config
    │   ├── lib/                   # Utility libraries
    │   └── utils/                 # Helpers
    ├── Dockerfile
    ├── package.json
    ├── vite.config.ts
    └── tailwind.config.js
```

---

## Containers

| Container | Image | Port | หน้าที่ |
|-----------|-------|------|--------|
| `tiaojin_gofiber` | tiaojin-system-gofiber | 9095 | Backend API |
| `tiaojin_vite` | tiaojin-system-vite | 3004→80 | Frontend |
| `tiaojin_db` | postgres:17-alpine | 5456→5432 | Database |
| `tiaojin_redis` | redis:7-alpine | (internal) | Cache |

---

## Deploy Commands

### แก้ Backend (Go Fiber) แล้ว deploy

```bash
ssh alljs
cd /opt/tiaojin-system

# แก้ไฟล์ใน gofiber/code/...
# แล้ว build + restart เฉพาะ backend
docker compose -f docker-compose.prod.yml up -d --build gofiber
```

### แก้ Frontend (Vite) แล้ว deploy

```bash
ssh alljs
cd /opt/tiaojin-system

# แก้ไฟล์ใน vite/code/src/...
# แล้ว build + restart เฉพาะ frontend
docker compose -f docker-compose.prod.yml up -d --build vite
```

### Deploy ทั้งหมด

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

### Restart โดยไม่ build

```bash
docker compose -f docker-compose.prod.yml restart gofiber   # backend
docker compose -f docker-compose.prod.yml restart vite       # frontend
docker compose -f docker-compose.prod.yml restart            # ทั้งหมด
```

---

## ดู Logs

```bash
# Backend logs (realtime)
docker logs tiaojin_gofiber --tail 50 -f

# Frontend logs
docker logs tiaojin_vite --tail 50 -f

# Database logs
docker logs tiaojin_db --tail 20

# ดูทั้งหมด
docker compose -f docker-compose.prod.yml logs -f
```

---

## Database

```bash
# เข้า psql
docker exec -it tiaojin_db psql -U postgres -d tiaojin_db

# Backup
docker exec tiaojin_db pg_dump -U postgres -d tiaojin_db -F c -f /tmp/backup.dump
docker cp tiaojin_db:/tmp/backup.dump /opt/tiaojin-system/backup_$(date +%Y%m%d).dump

# Restore
docker cp backup.dump tiaojin_db:/tmp/
docker exec tiaojin_db pg_restore -U postgres -d tiaojin_db --clean --if-exists /tmp/backup.dump
```

### DB Connection Info

| รายการ | ค่า |
|--------|-----|
| Host | `postgres` (ใน Docker) / `127.0.0.1` (จาก host) |
| Port | 5432 (ใน Docker) / **5456** (จาก host) |
| User | `postgres` |
| Password | `n147369` |
| Database | `tiaojin_db` |

### SSH Tunnel (เข้า DB จากเครื่อง local)

```bash
ssh -L 5456:127.0.0.1:5456 alljs
# แล้วเปิด pgAdmin/DBeaver เชื่อม localhost:5456
```

---

## Redis

```bash
# เข้า redis-cli
docker exec -it tiaojin_redis redis-cli -a n147369

# ดู keys
KEYS *

# Flush cache
FLUSHDB
```

---

## Environment Variables (Backend)

| ตัวแปร | ค่า | หมายเหตุ |
|--------|-----|----------|
| APP_PORT | 9095 | |
| DB_HOST | postgres | Docker internal |
| DB_PASSWORD | n147369 | |
| DB_NAME | tiaojin_db | |
| REDIS_HOST | redis | Docker internal |
| REDIS_PASSWORD | n147369 | |
| JWT_SECRET | Log2Window$P@ssWord | |
| S3_ENDPOINT | https://fcc0e164ed5f9fcf121a73f8f111ccd1.r2.cloudflarestorage.com | Cloudflare R2 |
| S3_ACCESS_KEY | ec1293bc10ab6bff0c09788a73b839dd | Cloudflare R2 |
| S3_SECRET_KEY | (ดูใน docker-compose.prod.yml) | Cloudflare R2 |
| S3_BUCKET | tiaojin | |
| S3_REGION | auto | |
| S3_PUBLIC_URL | https://pub-1fc850a8ec4e4e1395125c2ce7811e5d.r2.dev | |
| S3_PRIVATE_BUCKET | true | ใช้ presigned URLs |
| LLM_PROVIDER | gemini | |
| LLM_MODEL | gemini-2.5-flash | |
| CORS_ORIGINS | * | เปิดทุก origin (beta) |

**แก้ env**: แก้ใน `docker-compose.prod.yml` → section `gofiber.environment` → แล้ว restart

```bash
docker compose -f docker-compose.prod.yml up -d gofiber
```

---

## Environment Variables (Frontend — Build Args)

| ตัวแปร | ค่า |
|--------|-----|
| VITE_API_URL | https://api.tiaojin.co.th/api/v1 |
| VITE_WS_URL | wss://api.tiaojin.co.th/ws |
| VITE_APP_TITLE | เที่ยวจีนแอพ |
| VITE_APP_VERSION | 1.0.125 |

**แก้ env**: แก้ใน `docker-compose.prod.yml` → section `vite.build.args` → แล้ว **rebuild**

```bash
docker compose -f docker-compose.prod.yml up -d --build vite
```

⚠️ Frontend ต้อง `--build` เสมอเมื่อเปลี่ยน env เพราะ Vite bake env ตอน build

---

## Nginx Config

ไฟล์: `/etc/nginx/sites-available/tiaojin`

```nginx
server {
    listen 80;
    listen 443 ssl;
    server_name beta.tiaojin.co.th;
    ssl_certificate /etc/ssl/certs/tiaojin.crt;
    ssl_certificate_key /etc/ssl/private/tiaojin.key;
    location / {
        proxy_pass http://127.0.0.1:3004;
        ...
    }
}
server {
    listen 80;
    listen 443 ssl;
    server_name api.tiaojin.co.th;
    ssl_certificate /etc/ssl/certs/tiaojin.crt;
    ssl_certificate_key /etc/ssl/private/tiaojin.key;
    location / {
        proxy_pass http://127.0.0.1:9095;
        ...
    }
}
```

SSL: Self-signed cert + **Cloudflare Proxy (Full mode)**

แก้ Nginx แล้วต้อง:
```bash
nginx -t && systemctl reload nginx
```

---

## Health Check

```bash
# API health
curl http://127.0.0.1:9095/health

# Frontend health
curl http://127.0.0.1:3004/health

# ดูสถานะ containers
docker ps --filter name=tiaojin
```

---

## Troubleshooting

### Container ไม่ start
```bash
docker compose -f docker-compose.prod.yml logs gofiber   # ดู error
docker compose -f docker-compose.prod.yml down            # ลบทั้งหมด
docker compose -f docker-compose.prod.yml up -d --build   # สร้างใหม่
```

### DB เชื่อมไม่ได้
```bash
docker ps --filter name=tiaojin_db    # เช็คว่า DB รันอยู่ไหม
docker logs tiaojin_db --tail 20      # ดู error
```

### Frontend แสดงผลเก่า
```bash
# ต้อง rebuild เพราะ Vite เป็น static build
docker compose -f docker-compose.prod.yml up -d --build vite
```

### เปลี่ยน domain
1. แก้ `docker-compose.prod.yml` → vite build args (VITE_API_URL, VITE_WS_URL)
2. แก้ `/etc/nginx/sites-available/tiaojin` → server_name
3. Rebuild + reload:
```bash
docker compose -f docker-compose.prod.yml up -d --build vite
nginx -t && systemctl reload nginx
```

---

## Quick Reference

```bash
# SSH เข้า server
ssh alljs

# เข้า folder
cd /opt/tiaojin-system

# ดูสถานะ
docker ps --filter name=tiaojin

# แก้ backend → deploy
docker compose -f docker-compose.prod.yml up -d --build gofiber

# แก้ frontend → deploy
docker compose -f docker-compose.prod.yml up -d --build vite

# ดู logs
docker logs tiaojin_gofiber --tail 50 -f

# backup DB
docker exec tiaojin_db pg_dump -U postgres -d tiaojin_db -F c -f /tmp/backup.dump

# เข้า DB
docker exec -it tiaojin_db psql -U postgres -d tiaojin_db
```
