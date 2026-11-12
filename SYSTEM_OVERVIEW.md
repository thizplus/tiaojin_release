# TIAOJIN Tour Management System - Overview

## System Architecture

```
                    ┌─────────────────────────────────────────────┐
                    │              TIAOJIN SYSTEM                  │
                    └─────────────────────────────────────────────┘
                                        │
            ┌───────────────────────────┴───────────────────────────┐
            │                                                       │
    ┌───────▼───────┐                                     ┌────────▼────────┐
    │   FRONTEND    │                                     │    BACKEND      │
    │  (Vite React) │◄────────── REST API ────────────────│   (Go Fiber)    │
    │   Port: 5173  │                                     │   Port: 8080    │
    └───────┬───────┘                                     └────────┬────────┘
            │                                                      │
            │                                              ┌───────┴───────┐
            │                                              │               │
            │                                      ┌───────▼───┐   ┌───────▼───┐
            │                                      │ PostgreSQL│   │   Redis   │
            │                                      │  Database │   │   Cache   │
            │                                      └───────────┘   └───────────┘
            │
            └──────────────── WebSocket ──────────────────┘
                         (Real-time Updates)
```

---

## Tech Stack Summary

| Layer | Technology |
|-------|------------|
| **Frontend** | React 19, TypeScript, Vite, Tailwind CSS |
| **State Management** | Zustand, React Query |
| **UI Components** | shadcn/ui, Radix UI |
| **Backend** | Go Fiber v2 |
| **Database** | PostgreSQL (GORM) |
| **Cache** | Redis |
| **Authentication** | JWT |
| **File Storage** | Bunny CDN |
| **Real-time** | WebSocket |

---

## User Roles & Access

| Role | Description | Access Level |
|------|-------------|--------------|
| **Super Admin** | System Owner | Full Access - All features |
| **Admin** | System Administrator | Tours, Bookings, Users, Reports |
| **Agent** | Travel Agent | Book tours, View own bookings |
| **Sales** | Sales Staff | Track commissions, View bookings |

---

## Core Features

### 1. Authentication System
- JWT Token-based authentication
- Role-based access control (RBAC)
- API Key authentication for public APIs
- Password reset & update

### 2. Tour Management
- Create/Edit/Delete tours
- Tour categories (Regular, Special)
- City-based tour organization
- Tour status management (Draft/Publish)

### 3. Period Management
- Multiple periods per tour
- Seat availability tracking
- Room types & pricing
- Additional services
- Flight sectors

### 4. Booking Management
- Multi-step booking process
- Status workflow: Pending → Processing → Confirm → Completed
- Payment slip upload & verification
- Passport document upload
- Commission tracking

### 5. Seat Hold System (Admin Only)
- Temporary seat reservation
- Auto-expiration
- Convert to booking

### 6. Real-time Updates
- WebSocket for seat availability
- Live booking notifications
- Connection status management

### 7. Reporting & Analytics
- Dashboard statistics
- Booking reports
- Sales performance
- Commission reports

### 8. System Configuration
- Expiration tier settings
- Agent booking limits
- System-wide configurations

---

## API Endpoints Summary

| Category | Count | Base Path |
|----------|-------|-----------|
| Authentication | 8 | `/api/v1/auth/` |
| Users | 12 | `/api/v1/users/` |
| Tours | 10 | `/api/v1/tours/` |
| Periods | 9 | `/api/v1/periods/` |
| Bookings | 11 | `/api/v1/bookings/` |
| Agent | 11 | `/api/v1/agent/` |
| Sales | 11 | `/api/v1/sales/` |
| Admin | 50+ | `/api/v1/admin/` |
| Public API | 2 | `/api/v1/public/` |
| Config | 12 | `/api/v1/config/` |
| API Keys | 6 | `/api/v1/api-keys/` |
| **Total** | **150+** | |

---

## Frontend Pages Summary

| Role | Pages | Features |
|------|-------|----------|
| **Super** | 1 | System settings |
| **Admin** | 14 | Dashboard, Tours, Bookings, Users, Reports |
| **Agent** | 9 | Dashboard, Tours, Bookings, Profile, API Keys |
| **Sales** | 2 | Dashboard, Bookings |
| **Global** | 4 | Login, Register, Profile, Unauthorized |
| **Total** | **30** | |

---

## Database Models

| Category | Models |
|----------|--------|
| **User** | User, Sales, APIKey |
| **Tour** | Tour, Period, City, PeriodRateItem, PeriodRoomType, PeriodAdditionalService |
| **Booking** | Booking, BookingRateType, BookingRoomType, BookingAdditionalService, SeatHold, PaymentSlip, PassengerPassport |
| **System** | SystemConfig, CompanyProfile, Banner |
| **Total** | **22 Models** |

---

## Project Structure

```
_main-system/
├── gofiber/                 # Backend (Go Fiber)
│   ├── cmd/api/             # Entry point
│   ├── domain/              # Business logic
│   ├── application/         # Service implementations
│   ├── infrastructure/      # Database, Redis, Storage
│   ├── interfaces/api/      # HTTP handlers, routes, middleware
│   └── pkg/                 # Shared packages
│
├── vite/                    # Frontend (Vite React)
│   └── src/
│       ├── components/      # React components
│       ├── pages/           # Page components
│       ├── hooks/           # Custom hooks
│       ├── stores/          # Zustand stores
│       ├── services/        # API services
│       ├── types/           # TypeScript types
│       └── layouts/         # Layout components
│
├── BACKEND_DOCUMENTATION.md    # Backend docs
├── FRONTEND_DOCUMENTATION.md   # Frontend docs
└── SYSTEM_OVERVIEW.md          # This file
```

---

## Getting Started

### Backend (GoFiber)

```bash
cd gofiber
cp .env.example .env
# Edit .env with your configuration
go mod download
air  # or go run cmd/api/main.go
```

### Frontend (Vite React)

```bash
cd vite
npm install
npm run dev
```

---

## Environment Variables

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=tiaojin
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret
PORT=8080
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=Tiaojin
```

---

## Documentation Files

| File | Description |
|------|-------------|
| `BACKEND_DOCUMENTATION.md` | รายละเอียด Backend ทั้งหมด |
| `FRONTEND_DOCUMENTATION.md` | รายละเอียด Frontend ทั้งหมด |
| `SYSTEM_OVERVIEW.md` | ภาพรวมของระบบ |

---

## Statistics

| Metric | Backend | Frontend |
|--------|---------|----------|
| Source Files | 150+ | 150+ |
| API Endpoints | 150+ | - |
| Components | - | 60+ |
| Custom Hooks | - | 30+ |
| Database Models | 22 | - |
| Services | 17 | 20+ |
| Stores | - | 16 |
