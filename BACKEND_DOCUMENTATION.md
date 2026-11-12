# TIAOJIN Backend Documentation (GoFiber)

## Overview

ระบบ Backend สำหรับ Tiaojin Tour Management System พัฒนาด้วย Go Fiber Framework ใช้ Clean Architecture แบบ 4 ชั้น

**Tech Stack:**
- Go Fiber v2 (Web Framework)
- PostgreSQL (Database)
- GORM (ORM)
- Redis (Caching & Session)
- JWT (Authentication)
- Bunny CDN (File Storage)

---

## 1. Project Structure (โครงสร้างโปรเจค)

```
gofiber/
├── cmd/api/                      # Application Entry Point
│   └── main.go                   # DI initialization, Server startup, Graceful shutdown
│
├── domain/                       # Business Logic Layer
│   ├── models/                   # Data Models (22 files)
│   ├── repositories/             # Repository Interfaces
│   ├── services/                 # Service Interfaces (19 files)
│   └── dto/                      # Data Transfer Objects
│
├── application/                  # Application Layer
│   └── serviceimpl/              # Service Implementations (17 files)
│
├── infrastructure/               # Infrastructure Layer
│   ├── postgres/                 # PostgreSQL Repositories
│   ├── redis/                    # Redis Client
│   ├── storage/                  # Bunny CDN Storage
│   ├── services/                 # External Services (Vision API, LLM)
│   └── websocket/                # WebSocket Manager
│
├── interfaces/api/               # Interface/HTTP Layer
│   ├── handlers/                 # HTTP Handlers (22 files)
│   ├── routes/                   # Route Files (24 files)
│   ├── middleware/               # Middleware (8 files)
│   └── websocket/                # WebSocket Handlers
│
├── pkg/                          # Shared Packages
│   ├── config/                   # Configuration Management
│   ├── di/                       # Dependency Injection Container
│   ├── scheduler/                # Event-based Job Scheduler
│   ├── utils/                    # JWT, Password Hashing, Validators
│   ├── errors/                   # Custom Error Handling
│   └── credentials/              # Credentials Management
│
└── docs/                         # Documentation
```

---

## 2. Database Models (ตารางในฐานข้อมูล)

ระบบมี **22 Database Models** หลัก:

### 2.1 User Management

| Model | Description | Fields |
|-------|-------------|--------|
| **User** | ผู้ใช้ระบบ | id, email, password, role, status, created_at, updated_at |
| **Sales** | บุคลากรขาย | id, name, phone, email, is_active, commission_rate |
| **APIKey** | API Keys สำหรับ Public Access | id, key, user_id, name, is_active, expires_at |

### 2.2 Tour Management

| Model | Description | Fields |
|-------|-------------|--------|
| **Tour** | ทัวร์/แพ็คเก็จ | id, title, code, city_id, detail, thumbnail, status |
| **Period** | ช่วงเวลาของทัวร์ | id, tour_id, start_date, end_date, total_seats, price |
| **City** | เมือง/สถานที่ | id, name, description, thumbnail |
| **PeriodRateItem** | รายการอัตราราคา | id, period_id, rate_type, price |
| **PeriodRoomType** | ประเภทห้อง | id, period_id, room_type, price, capacity |
| **PeriodAdditionalService** | บริการเพิ่มเติม | id, period_id, service_name, price |

### 2.3 Booking Management

| Model | Description | Fields |
|-------|-------------|--------|
| **Booking** | การจองทัวร์ | id, period_id, user_id, sales_id, status, total_price, passengers |
| **BookingRateType** | ประเภทอัตราการจอง | id, booking_id, rate_type, quantity, price |
| **BookingRoomType** | ประเภทห้อง | id, booking_id, room_type, quantity |
| **BookingAdditionalService** | บริการเพิ่มเติม | id, booking_id, service_name, quantity |
| **SeatHold** | การจองที่นั่งชั่วคราว | id, period_id, admin_id, seats_held, expires_at |
| **PaymentSlip** | ใบเสร็จชำระเงิน | id, booking_id, file_url, amount, verified |
| **PassengerPassport** | หนังสือเดินทาง | id, booking_id, file_url, passport_number |

### 2.4 System Configuration

| Model | Description | Fields |
|-------|-------------|--------|
| **SystemConfig** | การกำหนดค่าระบบ | id, key, value (JSONB), description |
| **CompanyProfile** | โปรไฟล์บริษัท | id, user_id, company_name, tax_id, address |
| **Banner** | แบนเนอร์ | id, title, image_url, link, is_active, order |

### 2.5 Booking Status Flow

```
pending → processing → confirm → completed
                    ↘ cancel
                    ↘ expired
```

---

## 3. API Endpoints (Routes)

ระบบมี **150+ API Endpoints** แบ่งตามหมวดหมู่:

### 3.1 Authentication Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | ลงทะเบียนผู้ใช้ใหม่ | Public |
| POST | `/api/v1/auth/login` | เข้าสู่ระบบ | Public |
| POST | `/api/v1/auth/refresh` | รีเฟรช Token | Public |
| POST | `/api/v1/auth/validate` | ยืนยัน Token | Public |
| POST | `/api/v1/auth/reset-password` | รีเซตรหัสผ่าน | Public |
| POST | `/api/v1/auth/logout` | ออกจากระบบ | Protected |
| POST | `/api/v1/auth/update-password` | เปลี่ยนรหัสผ่าน | Protected |
| POST | `/api/v1/auth/revoke-all` | ยกเลิก Token ทั้งหมด | Protected |

### 3.2 User Management Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/users/by-email` | ค้นหา User ตามอีเมล | Public |
| GET | `/api/v1/users/profile` | ดูโปรไฟล์ตัวเอง | Protected |
| PUT | `/api/v1/users/profile` | อัปเดตโปรไฟล์ | Protected |
| GET | `/api/v1/users/` | ดูรายการผู้ใช้ทั้งหมด | Admin+ |
| GET | `/api/v1/users/stats` | สถิติผู้ใช้ | Admin+ |
| GET | `/api/v1/users/:id` | ดูรายละเอียดผู้ใช้ | Admin+ |
| PUT | `/api/v1/users/:id` | อัปเดตผู้ใช้ | Admin+ |
| PUT | `/api/v1/users/:id/status` | เปลี่ยนสถานะผู้ใช้ | Admin+ |
| POST | `/api/v1/users/:id/approve` | อนุมัติผู้ใช้ | Admin+ |
| POST | `/api/v1/users/:id/suspend` | ระงับผู้ใช้ | Admin+ |
| DELETE | `/api/v1/users/:id` | ลบผู้ใช้ | Super |
| POST | `/api/v1/users/:id/change-role` | เปลี่ยนบทบาท | Super |

### 3.3 Tour Management Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/tours/` | ดูรายการทัวร์ | Public |
| GET | `/api/v1/tours/regular` | ดูทัวร์ปกติพร้อมช่วงเวลา | Public |
| GET | `/api/v1/tours/special` | ดูทัวร์พิเศษ | Public |
| GET | `/api/v1/tours/:id` | ดูรายละเอียดทัวร์ | Public |
| GET | `/api/v1/tours/:id/periods` | ดูช่วงเวลาของทัวร์ | Public |
| POST | `/api/v1/tours/` | สร้างทัวร์ใหม่ | Admin+ |
| GET | `/api/v1/tours/stats` | สถิติทัวร์ | Admin+ |
| PUT | `/api/v1/tours/:id` | อัปเดตทัวร์ | Admin+ |
| DELETE | `/api/v1/tours/:id` | ลบทัวร์ | Admin+ |
| PATCH | `/api/v1/tours/:id/status` | เปลี่ยนสถานะทัวร์ | Admin+ |

### 3.4 Period Management Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/periods/` | ดูช่วงเวลาทั้งหมด | Public |
| GET | `/api/v1/periods/:id` | ดูรายละเอียดช่วงเวลา | Public |
| GET | `/api/v1/periods/:id/availability` | ดูความพร้อมที่นั่ง | Public |
| POST | `/api/v1/periods/` | สร้างช่วงเวลา | Admin+ |
| PUT | `/api/v1/periods/:id` | อัปเดตช่วงเวลา | Admin+ |
| DELETE | `/api/v1/periods/:id` | ลบช่วงเวลา | Admin+ |
| PUT | `/api/v1/periods/:id/status` | เปลี่ยนสถานะ | Admin+ |
| GET | `/api/v1/periods/stats` | สถิติช่วงเวลา | Admin+ |

### 3.5 Booking Management Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/bookings/stats` | สถิติการจอง | Admin+ |
| GET | `/api/v1/bookings/report` | รายงานการจอง | Admin+ |
| GET | `/api/v1/bookings/commission` | ข้อมูลค่าคอมมิชชัน | Sales+ |
| GET | `/api/v1/bookings/sales-report` | รายงานการขาย | Sales+ |
| GET | `/api/v1/bookings/agent-report` | รายงานเอเจนต์ | Agent+ |
| POST | `/api/v1/bookings/` | สร้างการจอง | Agent+ |
| GET | `/api/v1/bookings/` | ดูรายการจอง | Agent+ |
| GET | `/api/v1/bookings/:id` | ดูรายละเอียดจอง | Agent+ |
| PUT | `/api/v1/bookings/:id` | อัปเดตจอง | Agent+ |
| PUT | `/api/v1/bookings/:id/status` | เปลี่ยนสถานะ | Admin+ |
| DELETE | `/api/v1/bookings/:id` | ยกเลิกจอง | Agent+ |

### 3.6 Agent Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/agent/profile` | ดูโปรไฟล์ | Agent |
| PUT | `/api/v1/agent/profile` | อัปเดตโปรไฟล์ | Agent |
| POST | `/api/v1/agent/profile` | สร้างโปรไฟล์ | Agent |
| GET | `/api/v1/agent/profile/check` | ตรวจสอบโปรไฟล์ | Agent |
| GET | `/api/v1/agent/tours` | ดูทัวร์ที่จองได้ | Agent |
| GET | `/api/v1/agent/bookings` | ดูการจองของตัวเอง | Agent |
| POST | `/api/v1/agent/bookings` | สร้างการจอง | Agent |
| DELETE | `/api/v1/agent/bookings/:id` | ยกเลิกการจอง | Agent |
| POST | `/api/v1/agent/bookings/:id/payment-slip` | อัปโหลดใบเสร็จ | Agent |
| POST | `/api/v1/agent/bookings/:id/passport` | อัปโหลด Passport | Agent |
| GET | `/api/v1/agent/bookings/:id/documents` | ดูเอกสาร | Agent |

### 3.7 Sales Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/sales/active` | ดูรายชื่อ Sales | Agent+ |
| GET | `/api/v1/sales/admin/` | ดูรายชื่อทั้งหมด | Admin+ |
| POST | `/api/v1/sales/admin/` | สร้าง Sales | Admin+ |
| GET | `/api/v1/sales/admin/:id` | ดูรายละเอียด | Admin+ |
| PUT | `/api/v1/sales/admin/:id` | อัปเดต | Admin+ |
| DELETE | `/api/v1/sales/admin/:id` | ลบ | Admin+ |
| PATCH | `/api/v1/sales/admin/:id/activate` | เปิดใช้งาน | Admin+ |
| PATCH | `/api/v1/sales/admin/:id/deactivate` | ปิดใช้งาน | Admin+ |
| GET | `/api/v1/sales/bookings` | ดูการจองลูกค้า | Sales |
| GET | `/api/v1/sales/my-performance` | ดูประสิทธิภาพ | Sales |
| GET | `/api/v1/sales/my-commission` | ดูค่าคอมมิชชัน | Sales |

### 3.8 Seat Hold Routes (Admin Only)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/admin/periods/:id/hold-seats` | จองที่นั่งชั่วคราว | Admin+ |
| GET | `/api/v1/admin/periods/:id/seat-holds` | ดูการจองที่นั่ง | Admin+ |
| GET | `/api/v1/admin/periods/:id/seat-availability` | ดูความพร้อม | Admin+ |
| GET | `/api/v1/admin/seat-holds` | ดูทั้งหมด | Admin+ |
| GET | `/api/v1/admin/seat-holds/:hold_id` | ดูรายละเอียด | Admin+ |
| PATCH | `/api/v1/admin/seat-holds/:hold_id` | อัปเดต | Admin+ |
| DELETE | `/api/v1/admin/seat-holds/:hold_id/release` | ปล่อยคืน | Admin+ |
| POST | `/api/v1/admin/seat-holds/:hold_id/extend` | ขยายเวลา | Admin+ |
| DELETE | `/api/v1/admin/seat-holds/:hold_id` | ลบ | Super |
| POST | `/api/v1/admin/seat-holds/:hold_id/convert-to-booking` | แปลงเป็นจอง | Admin+ |

### 3.9 Public API Routes (API Key Auth)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/public/tours` | ดูทัวร์ | API Key |
| GET | `/api/v1/public/tours/:id` | ดูรายละเอียดทัวร์ | API Key |

### 3.10 API Key Management Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/api-keys/` | สร้าง API Key | Agent+ |
| GET | `/api/v1/api-keys/` | ดูรายการ | Agent+ |
| GET | `/api/v1/api-keys/:id` | ดูรายละเอียด | Agent+ |
| PUT | `/api/v1/api-keys/:id` | อัปเดต | Agent+ |
| DELETE | `/api/v1/api-keys/:id` | ยกเลิก | Agent+ |
| POST | `/api/v1/api-keys/validate` | ยืนยัน | Public |

### 3.11 System Configuration Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/config/` | ดูข้อมูลกำหนดค่า | Admin+ |
| GET | `/api/v1/config/expiration/tiers` | ดูระดับหมดอายุ | Admin+ |
| PUT | `/api/v1/config/expiration/tiers` | อัปเดตระดับ | Super |
| POST | `/api/v1/config/expiration/process-expired` | ประมวลผลหมดอายุ | Super |
| GET | `/api/v1/config/agent-booking-limits/` | ดูขีดจำกัด | Admin+ |
| PUT | `/api/v1/config/agent-booking-limits/` | อัปเดตขีดจำกัด | Super |

### 3.12 Health & Debug Routes

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | ตรวจสอบสุขภาพระบบ | Public |
| GET | `/` | ข้อมูลต้อนรับ | Public |
| GET | `/ws` | WebSocket Connection | Protected |

---

## 4. Middleware

ระบบมี **8 Middleware Files**:

| Middleware | Description |
|-----------|-------------|
| **Protected()** | ตรวจสอบ JWT Token |
| **RequireRole(roles...)** | ตรวจสอบบทบาท |
| **SuperAdminOnly()** | อนุญาตเฉพาะ Super Admin |
| **AdminOnly()** | อนุญาตเฉพาะ Admin+ |
| **ManagementOnly()** | อนุญาต Admin/Sales Management |
| **BookingAccess()** | สำหรับการจอง |
| **APIKeyAuth()** | ตรวจสอบ API Key |
| **Security()** | Helmet Security Headers |
| **BlockSuspiciousRequests()** | บล็อก SQL Injection/XSS |
| **CorsMiddleware()** | CORS Headers |
| **LoggerMiddleware()** | Request Logging |
| **ErrorHandler()** | Error Handling |
| **RateLimitMiddleware()** | Rate Limiting |

---

## 5. User Roles & Permissions

| Role | Description | Permissions |
|------|-------------|-------------|
| **super** | เจ้าของระบบ | Full Access ทุกระบบ |
| **admin** | ผู้ดูแลระบบ | จัดการผู้ใช้, ทัวร์, การจอง, รายงาน |
| **agent** | ตัวแทนขาย | จองทัวร์, ดูประสิทธิภาพ, API Keys |
| **sales** | พนักงานขาย | ติดตามค่าคอมมิชชัน, ดูการจอง |

---

## 6. Services (Business Logic)

มี **17 Service Implementations**:

| Service | Description |
|---------|-------------|
| AuthServiceImpl | การยืนยันตัวตน |
| UserServiceImpl | การจัดการผู้ใช้ |
| CompanyServiceImpl | โปรไฟล์บริษัท |
| CityServiceImpl | การจัดการเมือง |
| BannerServiceImpl | แบนเนอร์ |
| SalesServiceImpl | บุคลากรขาย |
| TourServiceImpl | ทัวร์ |
| PeriodServiceImpl | ช่วงเวลา |
| BookingServiceImpl | การจอง |
| TaskServiceImpl | งาน |
| FileServiceImpl | ไฟล์ |
| JobServiceImpl | ตำแหน่ง |
| NotificationServiceImpl | การแจ้งเตือน |
| APIKeyServiceImpl | API Key |
| SeatHoldServiceImpl | จองที่นั่งชั่วคราว |
| ConfigServiceImpl | การกำหนดค่า |
| LLMServiceImpl | Vision & LLM Services |

---

## 7. External Integrations

| Service | Description |
|---------|-------------|
| **PostgreSQL** | Primary Database (GORM ORM) |
| **Redis** | Caching & Session Management |
| **Bunny CDN** | File Storage Integration |
| **Vision API** | OCR for Passport/Receipt scanning |
| **LLM Service** | AI/ML services |
| **go-co-op/gocron** | Event-based Cron Job scheduler |

---

## 8. API Response Format

```json
{
  "success": true,
  "message": "Success message",
  "data": { /* Response data */ },
  "error": null,
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## 9. Key Features

1. **JWT Authentication** - Token-based authentication with refresh
2. **Role-Based Access Control (RBAC)** - การควบคุมการเข้าถึงตามบทบาท
3. **API Key Authentication** - สำหรับ Public API Access
4. **Seat Hold Management** - การจองที่นั่งชั่วคราว (Admin-only)
5. **Booking Status Workflow** - pending → processing → confirm → completed
6. **Commission & Sales Tracking** - ติดตามค่าคอมมิชชัน
7. **Payment Slip Verification** - ยืนยันใบเสร็จ (OCR)
8. **Document Upload** - อัปโหลด Passport, Payment Slip
9. **WebSocket Support** - Real-time seat updates
10. **Multi-tenant Support** - Company Profile Management
11. **System Configuration** - JSONB Config Store
12. **Analytics & Reporting** - รายงานการจอง, การขาย
13. **OCR Vision API** - สแกนเอกสารอัตโนมัติ
14. **Air Hot Reload** - Hot reload for development
15. **Docker Support** - Containerization ready

---

## 10. File Statistics

| Category | Count |
|----------|-------|
| Total Go Files | 150+ |
| Total Handlers | 22 |
| Total Routes | 24 |
| Total API Endpoints | 150+ |
| Total Middleware | 8 |
| Total Database Models | 22 |
| Total Services | 17 |

---

## 11. Important File Paths

| File | Path |
|------|------|
| Main Entry | `gofiber/cmd/api/main.go` |
| Routes Hub | `gofiber/interfaces/api/routes/routes.go` |
| Handlers | `gofiber/interfaces/api/handlers/` |
| Models | `gofiber/domain/models/` |
| Services | `gofiber/application/serviceimpl/` |
| DI Container | `gofiber/pkg/di/container.go` |
| Middleware | `gofiber/interfaces/api/middleware/` |

---

## 12. Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=tiaojin

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Bunny CDN
BUNNY_STORAGE_ZONE=your-zone
BUNNY_STORAGE_KEY=your-key

# Server
PORT=8080
ENV=development
```
