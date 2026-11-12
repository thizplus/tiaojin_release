# TIAOJIN - รายการแก้ไข (Change Log / TODO)

สร้างเมื่อ: 2026-09-06
อัปเดตล่าสุด: 2026-09-06 (ตรวจสอบ code จริงแล้ว)
สถานะ: รอดำเนินการ

---

## สารบัญ

1. [Banner หน้าทัวร์ (Agent)](#1-banner-หน้าทัวร์-agent)
2. [ระบบสร้าง Banner หน้าแรก](#2-ระบบสร้าง-banner-หน้าแรก)
3. [แสดงโปรไฟไหม้ใต้ช่องค้นหา](#3-แสดงโปรไฟไหม้ใต้ช่องค้นหา)
4. [WL แสดงใน Period List กรณีจ่ายเงินไม่ครบ](#4-wl-แสดงใน-period-list-กรณีจ่ายเงินไม่ครบ)
5. [แก้ไขระบบ หัก ณ ที่จ่าย 3% / VAT 7%](#5-แก้ไขระบบ-หัก-ณ-ที่จ่าย-3--vat-7)
6. [ใบนัดหมาย แสดงเมื่อชำระเงินครบเท่านั้น](#6-ใบนัดหมาย-แสดงเมื่อชำระเงินครบเท่านั้น)
7. [PDF ดาวน์โหลดโปรแกรม แสดง Header ของ Agent](#7-pdf-ดาวน์โหลดโปรแกรม-แสดง-header-ของ-agent)
8. [ปุ่มคัดลอกโปรแกรมขาย (Sale Script)](#8-ปุ่มคัดลอกโปรแกรมขาย-sale-script)
9. [ปุ่มคัดลอกรายละเอียด](#9-ปุ่มคัดลอกรายละเอียด)
10. [แสดงเลขที่ Invoice ในรายงาน](#10-แสดงเลขที่-invoice-ในรายงาน)

---

## 1. Banner หน้าทัวร์ (Agent)

### สถานะ: ต้องทำ (ง่าย — FE 1 ไฟล์)

### สถานะปัจจุบัน (ตรวจ code แล้ว)
- `BannerCarousel` component มีอยู่แล้วที่ `vite/src/components/agent/dashboard/BannerCarousel.tsx`
- Hook มี 2 ตัว:
  - `usePublicBanners` (`vite/src/hooks/usePublicBanners.ts`) — ใช้ใน Agent Dashboard
  - `usePublicBanner` (`vite/src/hooks/usePublicBanner.ts`) — ใช้ใน Homepage
- ใช้งานอยู่แล้ว 2 ที่: Agent Dashboard (`page/agent/dashboard/page.tsx`) + Homepage (`page/globals/home/page.tsx`)
- หน้า Agent Tours (`page/agent/tours/page.tsx`) มีแค่ `CityFilter`, `TourSearchForm`, `ToursRegularList` — ไม่มี Banner

### สิ่งที่ต้องทำ
เพิ่ม `BannerCarousel` ในหน้ารายการทัวร์ของ Agent

### ไฟล์ที่ต้องแก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/page/agent/tours/page.tsx` | เพิ่ม `usePublicBanners()` หรือ `usePublicBanner()` + `<BannerCarousel>` ระหว่าง Search กับ Tour List |

### Layout
```
หน้า Agent Tours ปัจจุบัน:
  [City Filter]
  [Tour Search]
  [Tour List]

หน้า Agent Tours ใหม่:
  [City Filter]
  [Tour Search]
  [BannerCarousel]     <-- เพิ่มตรงนี้
  [Tour List]
```

---

## 2. ระบบสร้าง Banner หน้าแรก

### สถานะ: มีครบแล้ว — ไม่ต้องทำ

### ตรวจ code แล้ว
- Admin จัดการ Banner: `vite/src/page/admin/banners/page.tsx` (CRUD ครบ: สร้าง/แก้/ลบ)
- Homepage แสดง Banner: `vite/src/page/globals/home/page.tsx` — ใช้ `usePublicBanner({ is_active: true, limit: 10 })` + `<BannerCarousel>` ในส่วน "โปรโมชั่นพิเศษ" (แสดงเมื่อมี banner)
- Backend: `gofiber/interfaces/api/handlers/banner_handler.go` — handler ครบ

### TODO (optional)
- Drag & drop จัดลำดับ Banner — ยังไม่มี (ถ้าต้องการ)

---

## 3. แสดงโปรไฟไหม้ใต้ช่องค้นหา

### สถานะ: ต้องทำ (ง่าย — FE 1 ไฟล์, BE มีแล้ว)

### สถานะปัจจุบัน (ตรวจ code แล้ว)

**Backend — พร้อมใช้แล้ว:**
- Route `GET /api/v1/tours/special` เป็น **public route** อยู่แล้ว (ไม่ต้อง auth)
- อยู่ใน `gofiber/interfaces/api/routes/tour_routes.go` line 16

**Frontend — Agent side มีครบ:**
- Service: `vite/src/services/tours-hot.ts` — `ToursHotService.getHotTours()`
- Hook: `vite/src/hooks/useToursHot.ts` — ครบทั้ง search, filter, city, WebSocket, summary
- Store: `vite/src/stores/toursHotStore.ts`
- Agent page: `vite/src/page/agent/tours-hot/page.tsx` — route `/agent/tours-hot`

**Frontend — Homepage ยังไม่มี:**
- `vite/src/page/globals/home/page.tsx` ไม่มี reference ถึง `useToursHot` หรือ hot tours เลย

### สิ่งที่ต้องทำ
เพิ่มส่วนแสดงทัวร์ไฟไหม้ใน Homepage โดยใช้ hook/service ที่มีอยู่แล้ว

### ไฟล์ที่ต้องแก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/page/globals/home/page.tsx` | เพิ่มส่วน Hot Tours ใต้ Hero search ก่อน Banner โดยใช้ `useToursHot` hook ที่มีอยู่แล้ว |

> Backend ไม่ต้องแก้ — route public มีอยู่แล้ว

### Layout
```
หน้าแรกปัจจุบัน:
  [Hero + Search Bar]
  [Banner โปรโมชั่นพิเศษ]
  [ทัวร์ทั้งหมด (Regular)]

หน้าแรกใหม่:
  [Hero + Search Bar]
  [โปรไฟไหม้ / Hot Tours]   <-- เพิ่มตรงนี้ (แสดง 3-6 รายการ)
  [Banner โปรโมชั่นพิเศษ]
  [ทัวร์ทั้งหมด (Regular)]
```

### ข้อควรพิจารณา
- แสดง limit 3-6 รายการ พร้อมปุ่ม "ดูทั้งหมด" (ถ้า login แล้ว ไปหน้า `/agent/tours-hot`)

---

## 4. WL แสดงใน Period List กรณีจ่ายเงินไม่ครบ

### สถานะ: ต้องทำ (ปานกลาง — BE 2 ไฟล์ + FE 2 ไฟล์) + ต้อง confirm definition กับลูกค้า

### สถานะปัจจุบัน (ตรวจ code แล้ว)

**ไม่มีระบบ WL ใดๆ ในระบบ:**
- Booking มี 6 สถานะ: `pending`, `processing`, `confirm`, `completed`, `cancel`, `expired` — ไม่มี `waiting_list`
- `IsPaid()` มีอยู่แล้ว (`gofiber/domain/models/booking.go` line 104) — เช็ค `SumPayment >= GrandTotal`
- `RemainingBalance()` มีอยู่แล้ว (line 114)
- `PeriodResponse` มีแค่: `seat`, `booking`, `held_seats`, `available_seats` — ไม่มี WL count
- `CalculateRealTimeBookingStats()` ใน `period_service_impl.go` นับ booking status `["pending", "processing", "confirm"]` — ไม่แยก WL
- Frontend `tour-periods-columns.tsx` แสดงแค่: ทั้งหมด / จอง / จำนอง / ว่าง
- Frontend `TourPeriod` type (`vite/src/types/models/period.ts`) ไม่มี `waiting_list_count`

### นิยาม WL (ต้องยืนยันกับลูกค้า)
- Booking ที่ status = `pending`, `processing`, หรือ `confirm` แต่ `RemainingBalance() > 0` (จ่ายยังไม่ครบ)?
- หรือนิยามอื่น?

### สิ่งที่ต้องทำ

**Backend:**

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `gofiber/domain/dto/period.go` | เพิ่ม field `WaitingListCount int json:"waitingListCount"` ใน `PeriodResponse` |
| `gofiber/application/serviceimpl/period_service_impl.go` | เพิ่มคำนวณ WL ใน `CalculateRealTimeBookingStats()` + อัปเดต `toPeriodResponse` / `toPeriodResponseWithSeparateTables` mappers |

**Frontend:**

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/types/models/period.ts` | เพิ่ม `waiting_list_count: number` ใน `TourPeriod` interface |
| `vite/src/components/manage/period/tour-periods-columns.tsx` | เพิ่มแสดง WL count ในคอลัมน์ "ที่นั่ง" (badge สีส้ม "WL: 3") |

### Layout
```
คอลัมน์ที่นั่งปัจจุบัน:
  ทั้งหมด: 30
  จอง: 20 | จำนอง: 5 | ว่าง: 5

คอลัมน์ที่นั่งใหม่:
  ทั้งหมด: 30
  จอง: 20 | จำนอง: 5 | ว่าง: 5
  WL: 3 (badge สีส้ม)    <-- Booking ที่ยังจ่ายไม่ครบ
```

---

## 5. แก้ไขระบบ หัก ณ ที่จ่าย 3% / VAT 7%

### สถานะ: ต้องทำ (ปานกลาง — FE 3 ไฟล์ + BE 1 ไฟล์ optional)

### สถานะปัจจุบัน (ตรวจ code แล้ว)

**Frontend — `BookingSheet.tsx`:**
- มี 2 toggles: `includeVAT` + `includeWHT` (lines 144-148)
- VAT = include model: `serviceFeeTotal x vatRate / (100 + vatRate)` — VAT ไม่บวกเพิ่ม แค่แสดงยอด
- WHT = `serviceFeeTotal x 3%` — **คิดจาก serviceFeeTotal (ไม่ใช่ netBeforeTax)**
- `grandTotal = netBeforeTax - whtAmount`

**Backend — `booking_service_impl.go`:**
- สูตรเดียวกับ frontend ทุกประการ (lines 187-205)
- Fields ที่เก็บ: `IncludeWHT`, `WHTRate`, `WHTAmount`, `IncludeVAT`, `VATRate`, `VATAmount`, `NetBeforeTax`, `GrandTotal`

**PDF — `pdfInvoiceNew.tsx`:**
- แสดง VAT (ถ้า `booking.include_vat`) + WHT (ถ้า `booking.include_wht`) ถูกต้อง

**BookingDetailSheet.tsx:**
- แสดงแค่ `grand_total` — ไม่แสดง VAT/WHT breakdown (ไม่มี conditional rows)

### สิ่งที่ต้องทำ

**A) ลบ Switch "หัก ณ ที่จ่าย 3%" ออกจาก UI**

**B) เมื่อเลือก VAT 7% — auto-fill หมายเหตุ**

### ตัวอย่างการคำนวณ
```
ค่าบริการ/คน = 500 บาท (รวม VAT), จำนวน 10 คน
serviceFeeTotal       = 500 x 10 = 5,000 บาท
vatAmount             = 5,000 x 7 / 107 = 327.10 บาท
serviceFeeExcludeVAT  = 5,000 - 327.10 = 4,672.90 บาท
whtForNote            = 4,672.90 x 3% = 140.19 บาท

หมายเหตุ: "กรุณาทำการหัก ณ ที่จ่าย 140.19 บาท"
```

### ไฟล์ที่ต้องแก้ไข

**Frontend:**

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/components/tours-regular/BookingSheet.tsx` | 1) ลบ Switch WHT + state `includeWHT` 2) เมื่อ VAT ON: คำนวณ whtForNote แล้ว auto-fill ช่องหมายเหตุ 3) เมื่อ VAT OFF: ลบข้อความจากหมายเหตุ 4) แก้ `grandTotal` ไม่หัก WHT (ส่ง `include_wht: false` เสมอ) 5) ลบ WHT ในส่วนสรุปราคา |
| `vite/src/components/globals/pdfInvoiceNew.tsx` | Booking ใหม่: ไม่แสดง WHT / Booking เก่าที่มี `include_wht: true`: ยังแสดง WHT (backward compatible — เช็คจาก `booking.include_wht`) |
| `vite/src/components/manage/booking/BookingDetailSheet.tsx` | ปัจจุบันไม่แสดง VAT/WHT breakdown อยู่แล้ว — ไม่ต้องแก้ (หรือเพิ่ม VAT row ถ้าต้องการ) |

**Backend (optional):**

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `gofiber/application/serviceimpl/booking_service_impl.go` | (optional) force `include_wht = false` เพื่อป้องกัน API call ตรง |

### ข้อควรระวัง
- Booking เก่าที่มี `include_wht = true` ยังต้องแสดงข้อมูลถูกต้อง (backward compatible)
- WHT คิดจาก `serviceFeeTotal` ไม่ใช่ `netBeforeTax` (ตาม code จริง)

---

## 6. ใบนัดหมาย แสดงเมื่อชำระเงินครบเท่านั้น

### สถานะ: ต้องทำ (ปานกลาง — FE 3 ไฟล์)

### สถานะปัจจุบัน (ตรวจ code แล้ว)

**Admin side — จัดการไฟล์:**
- `PeriodFilesDialog.tsx` — Admin ใช้อัปโหลด/ลบไฟล์ต่อ Period (ไม่มีเงื่อนไข — ถูกแล้ว Admin ต้องเห็นเสมอ)
- ใช้ใน `tour-periods-table.tsx` เท่านั้น

**Agent side — ดู/ดาวน์โหลด:**
- `FileDownloadSheet.tsx` — Agent ดูไฟล์ใบนัดหมาย
- เงื่อนไขปัจจุบัน: `canDownload = hasConfirmedSlip && hasPassport` — **เช็คแค่มี slip + passport ไม่ได้เช็ค full payment**
- เรียกจาก `page/agent/bookings/page.tsx` ส่ง `documentStatus: { hasConfirmedSlip, hasPassport }`
- `BookingGridCard.tsx` มีปุ่ม "ดูไฟล์" แสดงเมื่อ `period_details?.files?.length > 0` — ไม่มี payment gate

**ข้อมูลที่มีอยู่แล้ว:**
- `BookingResponse` มี `sum_payment` และ `grand_total` — สามารถเช็ค `sum_payment >= grand_total` ได้เลย

### สิ่งที่ต้องทำ
เปลี่ยนเงื่อนไขจาก `hasConfirmedSlip && hasPassport` เป็นเช็ค **ชำระเงินครบ** (`sum_payment >= grand_total`)

### ไฟล์ที่ต้องแก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/components/manage/booking/FileDownloadSheet.tsx` | เพิ่ม/แก้เงื่อนไข `canDownload` ให้เช็ค `isFullyPaid` (sum_payment >= grand_total) แทนหรือเพิ่มเติมจาก slip+passport |
| `vite/src/page/agent/bookings/page.tsx` | ส่ง prop `isFullyPaid: booking.sum_payment >= booking.grand_total` ไปให้ `FileDownloadSheet` |
| `vite/src/components/manage/booking/BookingGridCard.tsx` | เพิ่ม payment gate ที่ปุ่ม "ดูไฟล์" — ซ่อนหรือ disable ถ้ายังจ่ายไม่ครบ |

> หมายเหตุ: `PeriodFilesDialog.tsx` เป็น Admin side ไม่ต้องแก้ — Admin ต้องเห็นเสมอ

---

## 7. PDF ดาวน์โหลดโปรแกรม แสดง Header ของ Agent

### สถานะ: ต้องทำ (ปานกลาง — FE 1-2 ไฟล์)

### สถานะปัจจุบัน (ตรวจ code แล้ว)

**PDF generators มี 2 ตัว:**

1. `vite/src/utils/pdfGenerator.ts` — ใช้ `pdf-lib` + `@pdf-lib/fontkit`
   - Header hardcode: "BOOKING CHINA (บุ๊คกิ้ง ไช่น่า)" + Tax ID
   - **ไม่มี component ไหนใช้อยู่ (อาจเป็น dead code)**

2. `vite/src/components/globals/pdfInvoiceNew.tsx` — ใช้ `pdf-lib` + `@pdf-lib/fontkit`
   - Header hardcode: "BOOKING CHINA (บุ๊คกิ้ง ไช่น่า)" / "บริษัท เที่ยวจีน จำกัด" + ที่อยู่ + Tax ID + เบอร์โทร + บัญชีธนาคาร
   - **ตัวนี้ใช้งานจริง** — เรียกจาก `BookingDetailSheet.tsx` + `BookingGridCard.tsx`

**Agent Profile มีครบ:**
- Type: `vite/src/types/models/agentProfile.ts` — `AgentProfile` มี `name`, `tax_id`, `phone`, `address`, `province`, `zipcode`, `logo`
- Hook: `vite/src/hooks/agent/useAgentProfile.ts` — `useAgentProfile()` พร้อมใช้
- Store: `vite/src/stores/agent/agentProfileStore.ts`

**หมายเหตุ:** `BookingResponse` มี `company_profile` (agent company data) อยู่แล้ว แต่ใช้แสดงใน customer info section (ซ้าย) ไม่ใช่ header (บน)

### สิ่งที่ต้องทำ
แก้ `pdfInvoiceNew.tsx` ให้ใช้ Agent company profile เป็น header แทน operator header

### ไฟล์ที่ต้องแก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/components/globals/pdfInvoiceNew.tsx` | รับ agent profile data -> วาดโลโก้ + ชื่อบริษัท Agent ที่ header แทน hardcode "BOOKING CHINA" |

### Layout
```
PDF ปัจจุบัน:
  [Logo BOOKING CHINA]
  [ชื่อบริษัท + Tax ID]
  [Invoice content...]

PDF ใหม่:
  [Logo Agent Company]         <-- ใช้ข้อมูลจาก booking.company_profile หรือ AgentProfile
  [ชื่อบริษัท Agent + ที่อยู่]
  [Invoice content...]
```

### ข้อควรพิจารณา
- ต้องตรวจสอบว่า "ดาวน์โหลดโปรแกรม" หมายถึง PDF Invoice หรือ PDF โปรแกรมทัวร์ (itinerary)?
- Agent logo ต้อง fetch เป็น image blob เพื่อ embed ใน pdf-lib
- `pdfGenerator.ts` อาจเป็น dead code — ตรวจสอบว่ายังใช้อยู่ไหม ถ้าไม่ใช้อาจลบได้

---

## 8. ปุ่มคัดลอกโปรแกรมขาย (Sale Script)

### สถานะ: ต้องสร้างใหม่ทั้งหมด (FE 2 ไฟล์)

### สถานะปัจจุบัน (ตรวจ code แล้ว)
- ไม่มี Sale Script, clipboard copy, หรือ tour text generation ใดๆ ใน tour components
- Clipboard pattern มีตัวอย่างอยู่ใน:
  - `vite/src/components/manage/profile/company-profile-columns.tsx` — `copyToClipboard()` + `navigator.clipboard.writeText()`
  - `vite/src/components/manage/profile/ViewCompanyProfileDialog.tsx`
  - `vite/src/components/manage/api-key/CreateApiKeyDialog.tsx`

### ตำแหน่งที่เหมาะสม
- `TourCard.tsx` — tour card ที่ Agent เห็น (มี tour code, title, city, files, periods)
- `PeriodDetailSheet.tsx` — period detail sheet (มีข้อมูลครบ: dates, pricing, seats, commission)

### ไฟล์ที่ต้องสร้าง/แก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/utils/saleScript.ts` | **(สร้างใหม่)** ฟังก์ชัน `generateSaleScript(tour, period)` สร้าง text template |
| `vite/src/components/tours-regular/TourCard.tsx` หรือ `PeriodDetailSheet.tsx` | เพิ่มปุ่ม "คัดลอกโปรแกรมขาย" ใช้ `navigator.clipboard.writeText()` |

### ตัวอย่าง Sale Script Template
```
ทัวร์: {tour.title}
รหัส: {tour.code}
เดินทาง: {period.start_date} - {period.end_date}
ราคาเริ่มต้น: {lowest_rate} บาท/ท่าน
ที่นั่งเหลือ: {available_seats} ที่นั่ง
สายการบิน: {period.airline}

สนใจติดต่อ: {agent.company_name}
โทร: {agent.phone}
```

### ข้อควรพิจารณา
- Template ควรอยู่ที่ไหน? (hardcode, config, หรือ admin กำหนดเอง?)
- ต้องถาม user ว่าต้องการข้อมูลอะไรบ้างใน sale script
- Data ที่ต้องการ (`Tour`, `Period` types) มีครบใน `vite/src/types/tours-regular.ts`

---

## 9. ปุ่มคัดลอกรายละเอียด

### สถานะ: ต้องสร้างใหม่ทั้งหมด (FE 2 ไฟล์)

### สถานะปัจจุบัน (ตรวจ code แล้ว)
- ไม่มีปุ่มคัดลอกรายละเอียดทัวร์/Booking ที่ไหนเลย
- Clipboard pattern เดียวกับ Item 8

### ตำแหน่งที่เหมาะสม
- `PeriodDetailSheet.tsx` — มีข้อมูลครบที่สุด (dates, pricing, rooms, commission, additional services, seats)

### ความแตกต่างจาก Sale Script (Item 8)
- **Sale Script** = สั้น กระชับ สำหรับส่งลูกค้า
- **Copy Details** = ยาว รายละเอียดครบ (ราคาทุก rate, ห้องทุกประเภท, บริการเสริม)

### ไฟล์ที่ต้องสร้าง/แก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/utils/copyDetails.ts` | **(สร้างใหม่)** ฟังก์ชัน `generateTourDetails(tour, period)` สร้าง text รายละเอียดครบ |
| `vite/src/components/tours-regular/PeriodDetailSheet.tsx` | เพิ่มปุ่ม "คัดลอกรายละเอียด" (icon Copy) |

### ตัวอย่าง Copy Details Template
```
{tour.title} ({tour.code})
วันที่: {period.start_date} - {period.end_date} ({days} วัน)
สายการบิน: {period.airline} ({period.airport})

ราคา:
- {rate.name}: {rate.price} บาท
- ...

ห้องพัก:
- {room.name}: {room.price} บาท
- ...

บริการเสริม:
- {service.name}: {service.price} บาท
- ...

มัดจำ: {deposit} บาท
ที่นั่งเหลือ: {available} / {total}
```

---

## 10. แสดงเลขที่ Invoice ในรายงาน

### สถานะ: ต้องทำแค่ CSV Export (ง่าย — FE 1 ไฟล์)

### สถานะปัจจุบัน (ตรวจ code แล้ว)

**Backend — InvoiceID มีครบ:**
- Model: `Booking.InvoiceID` (`gofiber/domain/models/booking.go` line 23) — `gorm:"uniqueIndex;not null;size:20"`
- DTO: `BookingResponse.InvoiceID` (`gofiber/domain/dto/booking.go` line 11) — `json:"invoice_id"`

**Frontend — หน้ารายงานแสดง invoice_id แล้ว:**
- `BookingGridCard.tsx` แสดง `booking.invoice_id` อยู่แล้ว (line 141 desktop, line 188 mobile)
- `admin/reports-tours` ใช้ `BookingGridCard` → เห็น invoice_id แล้ว
- `admin/reports-daterange` ใช้ `BookingGridCard` → เห็น invoice_id แล้ว
- `agent/reports` เป็น chart dashboard (แสดง summary: total bookings, seats, sales) — ไม่มี booking list (by design)

**CSV Export — ยังขาด:**
- `vite/src/utils/exportBookingsCSV.ts` มี 26 columns แต่ **ไม่มี `invoice_id`**

### สิ่งที่ต้องทำ
เพิ่ม column `invoice_id` ใน CSV export เท่านั้น

### ไฟล์ที่ต้องแก้ไข

| ไฟล์ | การเปลี่ยนแปลง |
|------|---------------|
| `vite/src/utils/exportBookingsCSV.ts` | เพิ่ม header `'เลขที่ใบจอง'` + เพิ่ม `booking.invoice_id` ใน row data |

---

## สรุปภาพรวม (ตรวจ code จริงแล้ว)

| # | รายการ | สถานะ | ต้องทำ | ไฟล์ที่แก้ |
|---|--------|-------|--------|-----------|
| 1 | Banner หน้าทัวร์ | ขาด FE | ง่าย | FE 1 ไฟล์ |
| 2 | Banner หน้าแรก | **มีครบแล้ว** | ไม่ต้องทำ | — |
| 3 | โปรไฟไหม้ใต้ค้นหา | ขาด FE (BE มีแล้ว) | ง่าย | FE 1 ไฟล์ |
| 4 | WL ใน Period List | ไม่มีเลย | ปานกลาง + confirm definition | BE 2 + FE 2 ไฟล์ |
| 5 | แก้ WHT/VAT | มี WHT toggle ต้องลบ | ปานกลาง | FE 2-3 + BE 1 (opt) |
| 6 | ใบนัดหมาย เฉพาะจ่ายครบ | เงื่อนไขผิด (เช็ค slip+passport ไม่ใช่ full payment) | ปานกลาง | FE 3 ไฟล์ |
| 7 | PDF Agent Header | Hardcode operator | ปานกลาง | FE 1-2 ไฟล์ |
| 8 | ปุ่มคัดลอก Sale Script | ไม่มีเลย | ปานกลาง | FE 2 ไฟล์ (1 ใหม่) |
| 9 | ปุ่มคัดลอกรายละเอียด | ไม่มีเลย | ปานกลาง | FE 2 ไฟล์ (1 ใหม่) |
| 10 | Invoice ในรายงาน | หน้าเว็บมีแล้ว ขาดแค่ CSV | ง่าย | FE 1 ไฟล์ |
