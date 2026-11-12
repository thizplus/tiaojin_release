# Service Fee & VAT Analysis - สถานะปัจจุบันและแผนแก้ไข

## สถานะ: วิเคราะห์แล้ว - รอดำเนินการ

---

## 1. ระบบปัจจุบัน (Current State)

### 1.1 ค่าบริการ (Service Fee)

| รายการ | ค่าปัจจุบัน |
|--------|------------|
| ที่เก็บค่า | ENV variable `SERVICE_FEE_PER_PAX` หรือ default 500 บาท |
| ขอบเขต | **ทั้งระบบ** (ทุก Period ใช้ค่าเดียวกัน) |
| ต่อ Period | ไม่มี - Period model ไม่มี field นี้ |

**ปัญหา**: ค่าบริการเป็น system-wide constant ไม่สามารถกำหนดค่าบริการแยกแต่ละ Period ได้

### 1.2 การคำนวณ VAT - Backend

**File**: `gofiber/application/serviceimpl/booking_service_impl.go` (lines 177-198)

```
serviceFeePerPax = 500 (จาก ENV)
serviceFeeTotal  = serviceFeePerPax x totalQty

VAT = serviceFeeTotal x 7 / 107  (แยก VAT ออกจากราคารวม VAT)
```

**ตัวอย่าง**: ค่าบริการ 500 บาท/คน x 10 คน = 5,000 บาท
- VAT = 5,000 x 7 / 107 = **327.10 บาท**
- ค่าบริการก่อน VAT = 5,000 - 327.10 = 4,672.90 บาท

**สรุป Backend**: VAT คิดจาก **ค่าบริการเท่านั้น** (include VAT - แยก VAT ออกจากราคา) ถูกต้อง

### 1.3 การคำนวณ VAT - Frontend

**File**: `vite/src/components/tours-regular/BookingSheet.tsx` (lines 188-193)

```typescript
const calculateVATAmount = () => {
  if (!includeVAT) return 0;
  const netBeforeTax = calculateNetTotal();
  return netBeforeTax * (VAT_RATE / 100);  // VAT คิดจาก netBeforeTax ทั้งยอด!
};
```

**ตัวอย่าง**: ยอดก่อนภาษี 100,000 บาท
- VAT = 100,000 x 7% = **7,000 บาท** (ผิด!)

**สรุป Frontend**: VAT คิดจาก **ยอดทั้งหมด** (netBeforeTax) ไม่ใช่ค่าบริการ

### 1.4 Frontend vs Backend Mismatch

| รายการ | Backend | Frontend |
|--------|---------|----------|
| VAT คิดจาก | ServiceFeeTotal (ค่าบริการ) | NetBeforeTax (ยอดทั้งหมด) |
| สูตร VAT | `fee x 7 / 107` (include VAT) | `net x 7 / 100` (exclude VAT) |
| ค่าบริการ | จาก ENV (500 บาท) | ไม่มีข้อมูล |
| ผลลัพธ์ | ถูกต้อง | **ผิด - แสดงผลไม่ตรงกับ backend** |

**ปัญหาหลัก**:
1. Frontend แสดง VAT ผิด (คำนวณจากยอดทั้งหมด แต่ backend คำนวณจากค่าบริการ)
2. ตัวเลขที่ Agent เห็นตอนจอง ไม่ตรงกับที่ backend คำนวณจริง
3. ค่าบริการเป็น system-wide ไม่ยืดหยุ่นพอ (บาง Period อาจค่าบริการต่างกัน)

---

## 2. สิ่งที่ต้องการ (Requirement)

1. แต่ละ Period มี **ช่องกรอกค่าบริการ** (เช่น 500 บาท/คน)
2. VAT 7% คิดจาก **ค่าบริการเท่านั้น** ไม่ใช่ยอดทั้งหมด
3. เป็นแบบ **include VAT** (ราคาค่าบริการรวม VAT แล้ว ต้องแยก VAT ออก)
4. Frontend แสดงผลตรงกับ Backend

**สูตร**:
```
serviceFeeTotal = serviceFeePerPax x totalQty
vatAmount       = serviceFeeTotal x 7 / 107  (แยก VAT จากราคารวม)
grandTotal      = netBeforeTax + vatAmount - whtAmount
```

---

## 3. สิ่งที่ต้องแก้ไข (Changes Required)

### 3.1 Backend

#### A. Period Model - เพิ่ม field ค่าบริการ
**File**: `gofiber/domain/models/period.go`
```go
// เพิ่ม field ใน Period struct
ServiceFeePerPax float64 `gorm:"type:decimal(10,2);not null;default:0;column:service_fee_per_pax" json:"service_fee_per_pax"`
```

#### B. Period DTO - เพิ่ม field ใน request/response
**File**: `gofiber/domain/dto/period.go`
- CreatePeriodRequest: เพิ่ม `ServiceFeePerPax`
- UpdatePeriodRequest: เพิ่ม `ServiceFeePerPax` (pointer)
- PeriodResponse: เพิ่ม `ServiceFeePerPax`

#### C. Period Service - เพิ่มการบันทึกค่าบริการ
**File**: `gofiber/application/serviceimpl/period_service_impl.go`
- Create: เซ็ตค่า `ServiceFeePerPax` จาก request
- Update: อัพเดต `ServiceFeePerPax` ถ้ามีส่งมา

#### D. Booking Service - ใช้ค่าบริการจาก Period แทน ENV
**File**: `gofiber/application/serviceimpl/booking_service_impl.go`
- เปลี่ยนจาก `getServiceFeePerPax()` (ENV) เป็น `period.ServiceFeePerPax`
- ถ้า Period ไม่ได้ตั้งค่า (= 0) ให้ fallback ใช้ค่า default จาก ENV

### 3.2 Frontend

#### E. Period Form - เพิ่มช่อง input ค่าบริการ
**File**: `vite/src/components/manage/period/sections/PeriodBasicInfo.tsx`
- เพิ่ม input field "ค่าบริการ/คน (บาท, รวม VAT)"
- Default value: 500

#### F. Period Schema - เพิ่ม validation
**File**: `vite/src/components/manage/period/schemas/period-form-schemas.ts`
- เพิ่ม `service_fee_per_pax: z.number().min(0)`

#### G. Period Types - เพิ่ม type
- เพิ่ม `service_fee_per_pax` ใน Period type/interface

#### H. Booking Sheet - แก้ไขการคำนวณ VAT (สำคัญที่สุด)
**File**: `vite/src/components/tours-regular/BookingSheet.tsx`

แก้จาก:
```typescript
// ผิด - คิดจากยอดทั้งหมด
const calculateVATAmount = () => {
  if (!includeVAT) return 0;
  const netBeforeTax = calculateNetTotal();
  return netBeforeTax * (VAT_RATE / 100);
};
```

เป็น:
```typescript
// ถูก - คิดจากค่าบริการเท่านั้น (include VAT)
const calculateVATAmount = () => {
  if (!includeVAT) return 0;
  const totalPassengers = getTotalPassengers();
  const serviceFeePerPax = period.service_fee_per_pax || 500;
  const serviceFeeTotal = serviceFeePerPax * totalPassengers;
  return serviceFeeTotal * VAT_RATE / (100 + VAT_RATE); // แยก VAT จากราคารวม
};
```

---

## 4. สรุปไฟล์ที่ต้องแก้ไข

### Backend (5 files)
| File | การเปลี่ยนแปลง |
|------|---------------|
| `domain/models/period.go` | เพิ่ม `ServiceFeePerPax` field |
| `domain/dto/period.go` | เพิ่ม field ใน Create/Update/Response DTOs |
| `application/serviceimpl/period_service_impl.go` | เพิ่มการจัดการ service fee ใน create/update |
| `application/serviceimpl/booking_service_impl.go` | ใช้ `period.ServiceFeePerPax` แทน ENV |
| `domain/dto/mappers.go` | เพิ่ม service_fee_per_pax ใน period mapper (ถ้ามี) |

### Frontend (4-5 files)
| File | การเปลี่ยนแปลง |
|------|---------------|
| `components/manage/period/sections/PeriodBasicInfo.tsx` | เพิ่มช่อง input ค่าบริการ |
| `components/manage/period/schemas/period-form-schemas.ts` | เพิ่ม validation |
| `types/` (period types) | เพิ่ม `service_fee_per_pax` type |
| `components/tours-regular/BookingSheet.tsx` | แก้ VAT calculation ให้ตรง backend |
| `components/globals/pdfInvoiceNew.tsx` | อัพเดต Invoice PDF (ถ้าจำเป็น) |

### Database
- GORM auto-migrate: เพิ่ม column `service_fee_per_pax DECIMAL(10,2) DEFAULT 0` ใน `periods` table

---

## 5. ตัวอย่างการคำนวณที่ถูกต้อง (หลังแก้)

**สมมติ**: Period ตั้งค่าบริการ = 500 บาท/คน, จอง 10 คน

```
ยอดก่อนหักส่วนลด (subtotal)     = 190,000
ส่วนลด Agent (commission)        =  -5,000
ยอดก่อนภาษี (netBeforeTax)       = 185,000

ค่าบริการรวม                     = 500 x 10 = 5,000 (รวมอยู่ในราคาแล้ว)
VAT 7% (แยกจากค่าบริการ)         = 5,000 x 7 / 107 = 327.10

หัก ณ ที่จ่าย 3%                  = 185,000 x 3% = 5,550

ยอดที่ต้องชำระ (grandTotal)       = 185,000 + 327.10 - 5,550 = 179,777.10
```

**หมายเหตุ**: VAT คิดจากค่าบริการ 5,000 เท่านั้น ไม่ใช่ยอด 185,000 ทั้งจำนวน

---

**สร้างเมื่อ**: 2026-09-06
**อัพเดตล่าสุด**: 2026-09-06
**สถานะ**: ดำเนินการเสร็จสมบูรณ์
