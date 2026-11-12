# TIAOJIN - ระบบจองทัวร์ (System Overview)

## ภาพรวม

TIAOJIN เป็นระบบจองทัวร์แบบ B2B สำหรับจัดการแพ็กเกจทัวร์ การจอง การชำระเงิน และ Commission สำหรับบริษัททัวร์และตัวแทนขาย

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go Fiber v2.52 + GORM (PostgreSQL) + Redis |
| **Frontend** | React 19 + TypeScript + Vite 7 + Tailwind v4 + shadcn/ui |
| **State Management** | Zustand (18 stores) + TanStack React Query |
| **File Storage** | S3-compatible (iDrive e2) |
| **Real-time** | WebSocket (Fiber) |
| **AI/OCR** | Google Cloud Vision (Passport OCR) |
| **LLM** | Gemini / OpenAI / Claude (optional) |
| **Payment Verify** | SlipOK API (ตรวจสลิปธนาคารไทย) |
| **Email** | SMTP |
| **Scheduler** | gocron (event-driven) |

---

## Architecture (Clean Architecture)

```
Frontend (React/Vite :5173)
    |
    | HTTP / WebSocket
    v
Backend (Go Fiber :8080)
    |
    v
Handler Layer        --> HTTP logic, validation, auth
    |
    v
Service Layer        --> Business logic, calculations
    |
    v
Repository Layer     --> Database access (GORM)
    |
    v
PostgreSQL + Redis   --> Data persistence & caching
    |
Port/Adapter         --> External: S3, SlipOK, SMTP, Vision API
```

---

## User Roles (4 Roles)

| Role | สิทธิ์ |
|------|--------|
| **super** | System settings + ทุกอย่างของ Admin |
| **admin** | จัดการทัวร์, Booking, ผู้ใช้, ยืนยันชำระเงิน, รายงาน |
| **agent** | ดูทัวร์, สร้าง Booking, อัพโหลดเอกสาร, รายงานตัวเอง |
| **sales** | ดู Booking ของตัวเอง, Commission |

---

## Core Modules

### 1. Auth & User Management
- JWT Authentication (login, register, refresh token)
- Role-based access control (ProtectedRoute)
- Company profile + Bank account management

### 2. Tour Management
- สร้าง/แก้ไข Tour (code, title, description, city, files)
- จัดการ Period (ช่วงเวลา, ราคา, ที่นั่ง)
- Period Rate Items (ราคาตามประเภท: ผู้ใหญ่/เด็ก/ทารก)
- Period Room Types (ประเภทห้อง + ราคา)
- Period Additional Services (บริการเสริม: ประกัน, ไกด์)

### 3. Booking System (หัวใจของระบบ)
- สร้าง Booking พร้อมคำนวณราคาอัตโนมัติ
- เลือก Room Types, Rate Types, Additional Services
- ระบบหมดอายุอัตโนมัติ (Tiered Expiration)
- Status flow: `pending → processing → confirm → completed / cancel / expired`

### 4. Payment System
- อัพโหลดสลิปการชำระเงิน
- ตรวจสอบสลิปผ่าน SlipOK API (ธนาคารไทย)
- ระบบ Reversal (ยกเลิกการชำระ)
- รองรับ Cash / Transfer

### 5. Passport & OCR
- อัพโหลดรูปพาสปอร์ต
- Google Vision OCR ดึงข้อมูลอัตโนมัติ (ชื่อ, เลขพาสปอร์ต, สัญชาติ)

### 6. Seat Management
- Seat Availability (คำนวณที่นั่งว่าง real-time)
- Seat Hold (Admin จองที่นั่งชั่วคราว + auto-expire)

### 7. Sales & Commission
- จัดการพนักงานขาย
- Commission tracking ต่อ booking/rate type
- รายงาน Commission

### 8. Reporting & Analytics
- Admin reports (revenue, bookings, commission ทั้งหมด)
- Agent reports (booking ของตัวเอง)
- Sales reports (commission ที่ได้รับ)
- Monthly/Yearly breakdown

### 9. System Config
- ตั้งค่า Expiration tiers
- Booking limits
- Bank accounts
- Promotional codes

### 10. API Key Management
- สร้าง API Key สำหรับ 3rd-party integrations
- Public API routes (ใช้ API key auth)

---

## Business Flow หลัก

```
1. Admin สร้าง Tour
       ↓
2. Admin เพิ่ม Period (วันที่, ราคา, ห้อง, บริการเสริม, ที่นั่ง)
       ↓
3. Agent เลือก Tour/Period → สร้าง Booking
       ↓  (ระบบคำนวณราคาอัตโนมัติ + ตั้ง expiry)
4. ลูกค้าชำระเงิน → Agent อัพโหลดสลิป
       ↓
5. Admin ตรวจสอบสลิป (SlipOK API) → Confirm
       ↓
6. Agent อัพโหลด Passport → OCR ดึงข้อมูลอัตโนมัติ
       ↓
7. Booking สมบูรณ์ (completed)
```

---

## การคำนวณเงิน

```
sumRoom                = ผลรวมราคาห้องพัก
sumRate                = ผลรวมราคาผู้โดยสาร
sumAdditionalServices  = ผลรวมบริการเพิ่มเติม
─────────────────────────────────────────
subtotal               = sumRoom + sumRate + sumAdditionalServices
commission (discount)  = ส่วนลด Agent
─────────────────────────────────────────
netBeforeTax           = subtotal - commission
vatAmount              = netBeforeTax x 7% (ถ้าเลือก VAT)
whtAmount              = netBeforeTax x 3% (ถ้าเลือกหัก ณ ที่จ่าย)
─────────────────────────────────────────
grandTotal             = netBeforeTax + vatAmount - whtAmount
sumDeposit             = grandTotal x deposit_percent
```

---

## Seat Availability

```
Available = Period.total_seats - confirmed_bookings(qty) - active_seat_holds
```

---

## Expiration System

- แต่ละ Period มี Expiration Tiers (เช่น 30 วันก่อนออกเดินทาง)
- สร้าง Booking → ตั้ง cron job เฉพาะ booking นั้น
- หมดอายุ → status=expired, คืนที่นั่ง, คืน seat hold
- มี backup cleanup ทุกชั่วโมง

---

## ขนาดโปรเจกต์

| ส่วน | จำนวน |
|------|-------|
| Backend Models | ~23 |
| Service Interfaces | ~20 |
| Repositories | ~23 |
| HTTP Handlers | ~23 |
| API Endpoints | ~150+ |
| Frontend Hooks | ~45 |
| Zustand Stores | ~18 |
| Frontend Services | ~18 |
| shadcn/ui Components | ~40+ |
