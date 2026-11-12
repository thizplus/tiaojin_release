# Backend Structure (Go Fiber)

## Path: `_main-system/gofiber`

---

## Directory Structure

```
gofiber/
├── cmd/api/main.go                          # Entry point (server startup)
│
├── domain/                                  # Business Logic Layer
│   ├── models/                              # 23 Entity structs
│   │   ├── banner.go                        # แบนเนอร์
│   │   ├── booking_additional_service.go    # บริการเสริมใน Booking
│   │   ├── booking_rate_type.go             # ประเภทราคาใน Booking
│   │   ├── booking_room_type.go             # ประเภทห้องใน Booking
│   │   ├── city.go                          # เมือง
│   │   ├── file.go                          # ไฟล์อัพโหลด
│   │   ├── job.go                           # Cron job
│   │   ├── period_additional_service.go     # บริการเสริมของ Period
│   │   ├── period_rate_item.go              # รายการราคาของ Period
│   │   ├── period_room_type.go              # ประเภทห้องของ Period
│   │   ├── sales.go                         # พนักงานขาย
│   │   ├── seat_hold.go                     # การจองที่นั่งชั่วคราว
│   │   ├── task.go                          # Task
│   │   └── tour.go                          # ทัวร์
│   │
│   ├── dto/                                 # 21 Data Transfer Objects
│   │   ├── auth.go                          # Login/Register request/response
│   │   ├── booking.go                       # Booking request/response
│   │   ├── tour.go                          # Tour request/response
│   │   ├── period.go                        # Period request/response
│   │   ├── user.go                          # User request/response
│   │   ├── sales.go                         # Sales response
│   │   ├── company.go                       # Company profile response
│   │   ├── seat_hold.go                     # SeatHold request/response
│   │   ├── api_key.go                       # APIKey request/response
│   │   ├── vision.go                        # Vision OCR response
│   │   ├── file.go, job.go, notification.go, task.go
│   │   └── mappers.go                       # Model → DTO conversion
│   │
│   ├── repositories/                        # 23 Repository Interfaces
│   │   ├── booking_additional_service.go
│   │   ├── booking_rate_type.go
│   │   ├── booking_room_type.go
│   │   ├── city.go
│   │   ├── file_repository.go
│   │   ├── passenger_passport.go
│   │   ├── payment_slip.go
│   │   ├── period_additional_service.go
│   │   ├── period_rate_item.go
│   │   ├── period_room_type.go
│   │   ├── system_config.go
│   │   ├── task_repository.go
│   │   └── ...
│   │
│   ├── services/                            # 20 Service Interfaces
│   │   ├── banner.go
│   │   ├── city_service.go
│   │   ├── file_service.go
│   │   ├── notification.go
│   │   ├── sales.go
│   │   ├── seat_hold_service.go
│   │   ├── task_service.go
│   │   ├── vision.go
│   │   ├── llm.go
│   │   └── ...
│   │
│   └── ports/                               # 3 Port Interfaces (Adapter Pattern)
│       ├── storage_port.go                  # S3 file storage
│       ├── slip_verification_port.go        # SlipOK bank slip verify
│       └── email_port.go                    # SMTP email
│
├── application/serviceimpl/                 # 18 Service Implementations
│   ├── task_service_impl.go
│   ├── notification_service_impl.go
│   └── ...
│
├── infrastructure/                          # External Integrations
│   ├── postgres/                            # 24 Repository Implementations
│   │   ├── database.go                      # DB connection + auto-migration
│   │   ├── task_repository_impl.go
│   │   ├── file_repository_impl.go
│   │   ├── city_repository_impl.go
│   │   ├── booking_room_type_repository_impl.go
│   │   ├── booking_rate_type_repository_impl.go
│   │   ├── booking_additional_service_repository_impl.go
│   │   ├── period_rate_item_repository_impl.go
│   │   ├── period_room_type_repository_impl.go
│   │   ├── period_additional_service_repository_impl.go
│   │   ├── passenger_passport_repository_impl.go
│   │   ├── payment_slip_repository_impl.go
│   │   ├── system_config_repository_impl.go
│   │   └── ...
│   ├── redis/redis.go                       # Redis client (cache/session)
│   ├── storage/                             # S3 adapter (iDrive e2)
│   ├── email/                               # SMTP adapter
│   ├── slipverification/                    # SlipOK adapter
│   ├── services/vision_service_impl.go      # Google Vision OCR
│   └── websocket/websocket.go               # WebSocket manager
│
├── interfaces/api/                          # HTTP Layer
│   ├── handlers/                            # 23 HTTP Handlers
│   │   ├── analytics_handler.go             # Reports & analytics
│   │   └── ...
│   ├── routes/                              # 25 Route Definitions
│   │   ├── auth_routes.go
│   │   ├── user_routes.go
│   │   ├── profile_routes.go
│   │   ├── city_routes.go
│   │   ├── banner_routes.go
│   │   ├── task_routes.go
│   │   ├── file_routes.go
│   │   ├── period_routes.go
│   │   ├── websocket_routes.go
│   │   ├── debug_routes.go
│   │   └── ...
│   ├── middleware/
│   │   ├── logger_middleware.go             # Request/response logging
│   │   ├── error_middleware.go              # Global error handling
│   │   ├── rate_limit_middleware.go         # Rate limiting
│   │   ├── validation_middleware.go         # Input validation
│   │   └── security_middleware.go           # Security headers
│   └── websocket/websocket_handler.go       # WebSocket handlers
│
├── pkg/                                     # Shared Packages
│   ├── config/config.go                     # ENV loader
│   ├── di/container.go                      # Manual DI container
│   ├── scheduler/scheduler.go               # gocron wrapper
│   └── utils/
│       ├── response.go                      # HTTP response helpers
│       ├── jwt.go                           # JWT utilities
│       └── validator.go                     # Validation helpers
│
├── postman/                                 # Postman API collection
└── docs/                                    # Backend-specific docs
```

---

## Data Models

### User & Company
| Model | Fields (key) | หน้าที่ |
|-------|-------------|---------|
| User | id, username, email, password, role, status | บัญชีผู้ใช้ |
| CompanyProfile | user_id, company_name, tax_id, address | ข้อมูลบริษัท |
| CompanyBankAccount | bank_name, account_number, account_name | บัญชีธนาคาร |

### Tour & Period
| Model | Fields (key) | หน้าที่ |
|-------|-------------|---------|
| Tour | id, code, title, city_id, status, files | แพ็กเกจทัวร์ |
| City | id, name, country | เมืองปลายทาง |
| Period | id, tour_id, start/end_date, total_seats, status | ช่วงเวลาทัวร์ |
| PeriodRateItem | period_id, name, regular_price, sale_price, commission | รายการราคา |
| PeriodRoomType | period_id, name, price, quantity | ประเภทห้อง |
| PeriodAdditionalService | period_id, name, price | บริการเสริม |

### Booking & Payment
| Model | Fields (key) | หน้าที่ |
|-------|-------------|---------|
| Booking | id, invoice_id, user_id, tour_id, period_id, sales_id, status, grand_total | การจอง |
| BookingRoomType | booking_id, room_type, quantity, price | ห้องที่เลือก |
| BookingRateType | booking_id, rate_name, quantity, price, commission | ราคาที่เลือก |
| BookingAdditionalService | booking_id, service_name, quantity, price | บริการเสริมที่เลือก |
| PaymentSlip | booking_id, file_url, amount, status, verification_data | สลิปชำระเงิน |
| PassengerPassport | booking_id, name, passport_number, nationality, ocr_data | พาสปอร์ต |

### Other
| Model | หน้าที่ |
|-------|---------|
| SeatHold | จองที่นั่งชั่วคราว (active/released/expired/converted) |
| Sales | พนักงานขาย |
| Job | Cron job definitions |
| SystemConfig | System key-value settings |
| APIKey | API keys สำหรับ 3rd-party |
| Banner | แบนเนอร์โฆษณา |
| File | Metadata ไฟล์อัพโหลด |

---

## API Routes (~150+ endpoints)

### Authentication (Public)
```
POST   /api/v1/auth/login
POST   /api/v1/auth/register
POST   /api/v1/auth/refresh-token
```

### Protected Routes (JWT Required)
```
# User & Profile
GET/PUT    /api/v1/users
GET/PUT    /api/v1/profile
GET/PUT    /api/v1/company

# Tour Management
CRUD       /api/v1/tours
POST       /api/v1/tours/:id/publish
POST       /api/v1/tours/:id/unpublish

# Period Management
CRUD       /api/v1/periods
           /api/v1/periods/:id/rate-items
           /api/v1/periods/:id/room-types
           /api/v1/periods/:id/additional-services

# Booking (Main)
CRUD       /api/v1/bookings
PUT        /api/v1/bookings/:id/status
POST       /api/v1/bookings/:id/payment-slip
POST       /api/v1/bookings/:id/verify-payment
POST       /api/v1/bookings/:id/passport
GET        /api/v1/bookings/reports/*

# Seat Management
GET        /api/v1/seat-availability/:periodId
CRUD       /api/v1/seat-hold

# Other
CRUD       /api/v1/sales
CRUD       /api/v1/cities
CRUD       /api/v1/banners
CRUD       /api/v1/files
GET/POST   /api/v1/storage/presigned-url
CRUD       /api/v1/config (super only)
CRUD       /api/v1/api-keys
GET        /api/v1/public/* (API key auth)
WS         /ws
```

---

## Middleware Chain (ลำดับสำคัญ)
1. `RequestIDMiddleware` - สร้าง request_id
2. `LoggerMiddleware` - Log request/response
3. `CORSMiddleware` - CORS headers
4. `AuthMiddleware` (protected) - ตรวจ JWT
5. `APIKeyMiddleware` (public) - ตรวจ API key

---

## Port/Adapter Pattern

| Port Interface | Adapter | ใช้ทำอะไร |
|---------------|---------|-----------|
| `StoragePort` | S3 (iDrive e2) | อัพโหลด/ลบไฟล์, presigned URL |
| `SlipVerificationPort` | SlipOK API | ตรวจสลิปธนาคารไทย |
| `EmailPort` | SMTP | ส่ง email แจ้งเตือน |

---

## DI Container (`pkg/di/container.go`)

ลำดับ initialization:
1. Config loading (ENV)
2. Infrastructure (DB, Redis, S3, Email, SlipOK)
3. Repositories
4. Services (inject dependencies)
5. Scheduler startup
6. Graceful shutdown

---

## Environment Variables

```env
# App
APP_NAME=TIAOJIN
APP_PORT=3000

# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=xxx
DB_NAME=tiaojin

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Auth
JWT_SECRET=xxx

# Storage (S3)
S3_ENDPOINT=s3.ap-southeast-1.idrivee2.com
S3_REGION=ap-southeast-1
S3_BUCKET=tiaojin
S3_CDN_URL=https://cdn.example.com

# Payment Verification
SLIPOK_API_KEY=xxx
SLIPOK_BRANCH_ID=xxx

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=xxx
SMTP_PASS=xxx

# AI (Optional)
LLM_ENABLED=false
LLM_PROVIDER=gemini
LLM_API_KEY=xxx
GOOGLE_VISION_ENABLED=false
```

---

## DB Connection Config
- Max Open Connections: 25
- Max Idle Connections: 5
- Connection Max Lifetime: 5 min
- Soft deletes via GORM `DeletedAt`
- Auto-migration on startup
