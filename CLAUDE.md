# TIAOJIN - Tour Booking Management System

ระบบจองทัวร์แบบ B2B สำหรับจัดการแพ็กเกจทัวร์ การจอง การชำระเงิน และ Commission

---

# Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go Fiber v2.52 + GORM (PostgreSQL) + Redis |
| **Frontend** | React 19 + TypeScript + Vite 7 + Tailwind v4 + shadcn/ui |
| **State** | Zustand 5 (18 stores) + TanStack React Query 5 |
| **Storage** | S3-compatible (iDrive e2) |
| **Real-time** | WebSocket (Fiber) |
| **AI/OCR** | Google Cloud Vision (Passport OCR), LLM (optional) |
| **Payment** | SlipOK API (ตรวจสลิปธนาคารไทย) |
| **Scheduler** | gocron (booking expiry, seat hold release) |

---

# Project Structure

```
_main-system/
├── gofiber/                    # Backend (Go Fiber)
│   ├── cmd/api/main.go         # Entry point
│   ├── domain/
│   │   ├── models/             # 23 Entity structs
│   │   ├── dto/                # 21 Request/Response DTOs + mappers.go
│   │   ├── repositories/       # 23 Repository interfaces
│   │   ├── services/           # 20 Service interfaces
│   │   └── ports/              # 3 Adapter interfaces (Storage, SlipOK, Email)
│   ├── application/serviceimpl/  # 18 Service implementations
│   ├── infrastructure/
│   │   ├── postgres/           # 24 Repository implementations + database.go
│   │   ├── redis/              # Redis client
│   │   ├── storage/            # S3 adapter
│   │   ├── email/              # SMTP adapter
│   │   ├── slipverification/   # SlipOK adapter
│   │   ├── services/           # Vision & LLM service impl
│   │   └── websocket/          # WebSocket manager
│   ├── interfaces/api/
│   │   ├── handlers/           # 23 HTTP handlers
│   │   ├── routes/             # 25 Route definitions
│   │   ├── middleware/         # RequestID, Logger, CORS, Auth, APIKey, RateLimit
│   │   └── websocket/         # WebSocket handlers
│   └── pkg/
│       ├── config/             # ENV loader
│       ├── di/container.go     # Manual DI container
│       ├── scheduler/          # gocron wrapper
│       └── utils/              # Response, JWT, Validator
├── vite/                       # Frontend (React)
│   └── src/
│       ├── routes/             # Route definitions
│       ├── page/               # Pages (globals, admin, agent, sale, super)
│       ├── components/         # UI (shadcn/ui 40+), auth, manage, admin, agent
│       ├── layouts/            # Sidebar, nav, root layout
│       ├── hooks/              # 45+ custom hooks
│       ├── stores/             # 18 Zustand stores
│       ├── services/           # 18+ API services (Axios wrapper)
│       ├── constants/          # API endpoints, enums
│       ├── types/              # TypeScript interfaces (models, request, response)
│       └── theme/              # Dark/Light mode
└── docs/                       # Documentation
```

---

# Architecture

```
Frontend (React :5173)
    |  /api proxy
    v
Handler (HTTP logic, validation, auth check)
    |
Service Interface (domain/services/)
    |
Service Implementation (application/serviceimpl/)
    |
Repository Interface (domain/repositories/)
    |
Repository Implementation (infrastructure/postgres/)
    |
PostgreSQL + Redis

Port/Adapter (domain/ports/) --> S3, SlipOK, SMTP
```

---

# User Roles

| Role | สิทธิ์ |
|------|--------|
| `super` | System settings + ทุกอย่างของ Admin |
| `admin` | จัดการทัวร์, Booking, ผู้ใช้, ยืนยันชำระเงิน, รายงาน |
| `agent` | ดูทัวร์, สร้าง Booking, อัพโหลดเอกสาร, รายงานตัวเอง |
| `sales` | ดู Booking ของตัวเอง, Commission |

---

# Core Domains

### Tour & Period
- `Tour` - แพ็กเกจทัวร์ (code, title, city, files)
- `Period` - ช่วงเวลา + ราคา + ที่นั่ง
- `PeriodRateItem` - รายการราคา (ผู้ใหญ่/เด็ก/ทารก)
- `PeriodRoomType` - ประเภทห้อง + ราคา
- `PeriodAdditionalService` - บริการเสริม

### Booking & Payment
- `Booking` - การจอง (status: pending -> processing -> confirm -> completed/cancel/expired)
- `BookingRoomType`, `BookingRateType`, `BookingAdditionalService` - รายการที่เลือก
- `PaymentSlip` - สลิปชำระเงิน (verify ผ่าน SlipOK)
- `PassengerPassport` - พาสปอร์ต (OCR ผ่าน Google Vision)

### User & Company
- `User` - บัญชีผู้ใช้ (4 roles)
- `CompanyProfile` - ข้อมูลบริษัท
- `CompanyBankAccount` - บัญชีธนาคาร

### Other
- `SeatHold` - จองที่นั่งชั่วคราว (admin)
- `Sales` - พนักงานขาย + Commission
- `Job` - Cron jobs (expiry, cleanup)
- `SystemConfig` - System settings
- `APIKey` - API keys สำหรับ 3rd-party
- `Banner`, `City`, `File`

---

# Business Flow

```
1. Admin สร้าง Tour -> เพิ่ม Period (ราคา, ห้อง, บริการเสริม, ที่นั่ง)
2. Agent เลือก Tour/Period -> สร้าง Booking (คำนวณราคาอัตโนมัติ + ตั้ง expiry)
3. ลูกค้าชำระเงิน -> Agent อัพโหลดสลิป
4. Admin ตรวจสอบสลิป (SlipOK API) -> Confirm
5. Agent อัพโหลด Passport -> OCR ดึงข้อมูลอัตโนมัติ
6. Booking สมบูรณ์ (completed)
```

---

# Booking Calculation

```
subtotal       = sumRoom + sumRate + sumAdditionalServices
netBeforeTax   = subtotal - commission (discount)
vatAmount      = netBeforeTax x 7% (ถ้าเลือก VAT)
whtAmount      = netBeforeTax x 3% (ถ้าเลือกหัก ณ ที่จ่าย)
grandTotal     = netBeforeTax + vatAmount - whtAmount
```

---

# Code Conventions

### Backend (Go)
- **Handlers**: HTTP logic, validation, call services, return DTO (ห้าม return model ตรง)
- **Services**: Business logic, no HTTP dependency
- **Repositories**: Database access only, ห้าม JOIN ข้าม module
- **DTOs**: Request/Response with `camelCase` JSON tags, ใช้ mappers.go แปลง
- **Ports**: Interface สำหรับ external dependencies (S3, SlipOK, SMTP)
- **Naming**: PascalCase (struct/method), camelCase (variable/JSON tag), snake_case (DB column - GORM auto)

### Frontend (React)
- **Data Flow**: Page -> Hook -> Store -> Service -> apiClient -> API
- **State**: Zustand (feature state) + React Query (server data) + useState (UI)
- **Services**: ใช้ apiService wrapper (Axios + interceptors + auth)
- **Components**: shadcn/ui, semantic CSS classes (ห้าม hardcode colors)
- **Types**: แยก models/, request/, response/

---

# How to Implement Features

### Backend
1. Create Model -> `domain/models/feature.go`
2. Create Repository Interface -> `domain/repositories/feature.go`
3. Create Service Interface -> `domain/services/feature.go`
4. Implement Repository -> `infrastructure/postgres/feature_repository_impl.go`
5. Implement Service -> `application/serviceimpl/feature_service_impl.go`
6. Create DTOs -> `domain/dto/feature.go` (with mapper)
7. Create Handler -> `interfaces/api/handlers/feature_handler.go`
8. Create Routes -> `interfaces/api/routes/feature_routes.go`
9. Register in DI -> `pkg/di/container.go`

### Frontend
1. Create Types -> `types/feature.ts`
2. Create Service -> `services/featureService.ts`
3. Create Store -> `stores/featureStore.ts`
4. Create Hook -> `hooks/useFeature.ts`
5. Create Page -> `page/admin/feature-page.tsx`
6. Add Route -> `routes/index.tsx`
7. Add API constant -> `constants/api.ts`

---

# Key Files

### Backend
- Entry: `gofiber/cmd/api/main.go`
- DI Container: `gofiber/pkg/di/container.go`
- Config: `gofiber/pkg/config/config.go`
- Database: `gofiber/infrastructure/postgres/database.go`
- Routes Setup: `gofiber/interfaces/api/routes/routes.go`
- Response Utils: `gofiber/pkg/utils/response.go`

### Frontend
- Entry: `vite/src/main.tsx`
- App: `vite/src/App.tsx`
- Routes: `vite/src/routes/index.tsx`
- API Client: `vite/src/services/apiService.ts`
- Auth Store: `vite/src/stores/authStore.ts`
- API Constants: `vite/src/constants/api.ts`

---

# Dev Config

```
Backend:   localhost:8080  (Go Fiber, air hot-reload)
Frontend:  localhost:5173  (Vite, proxy /api -> :8080)
Database:  PostgreSQL
Cache:     Redis
```

---

# Documentation

- `docs/SYSTEM_OVERVIEW.md` - ภาพรวมระบบ
- `docs/BACKEND_STRUCTURE.md` - โครงสร้าง Backend ละเอียด
- `docs/FRONTEND_STRUCTURE.md` - โครงสร้าง Frontend ละเอียด
- `docs/VAT_WITHHOLDING_TAX_FEATURE.md` - ระบบ VAT และหัก ณ ที่จ่าย