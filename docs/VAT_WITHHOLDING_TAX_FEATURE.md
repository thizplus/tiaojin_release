# Feature: VAT และ หัก ณ ที่จ่าย สำหรับ Booking

## สถานะ: ✅ เสร็จสมบูรณ์

---

## 1. ภาพรวมระบบ

### 1.1 การคำนวณเงินในระบบ
```
sumRoom                = ผลรวมราคาห้องพัก
sumRate                = ผลรวมราคาผู้โดยสาร (ทุกคนจ่าย รวม excluded)
sumAdditionalServices  = ผลรวมบริการเพิ่มเติม
─────────────────────────────────────────────────────
subtotal               = sumRoom + sumRate + sumAdditionalServices
commission (discount)  = ส่วนลด Agent (เฉพาะ non-excluded passengers)
─────────────────────────────────────────────────────
netBeforeTax (ยอดก่อนภาษี) = subtotal - commission
vatAmount              = netBeforeTax × 0.07 (ถ้าเลือก VAT)
whtAmount              = netBeforeTax × 0.03 (ถ้าเลือกหัก ณ ที่จ่าย)
─────────────────────────────────────────────────────
grandTotal (ยอดที่ต้องชำระ) = netBeforeTax + vatAmount - whtAmount
sumDeposit             = grandTotal × deposit_percent หรือ deposit_per_unit × qty
```

### 1.2 สรุปการตั้งค่า (ตามที่ User กำหนด)

| ประเด็น | คำตอบ |
|--------|-------|
| VAT คิดจากยอดอะไร | ยอดหลังหักส่วนลด (netBeforeTax) |
| WHT คิดจากยอดอะไร | ยอดก่อน VAT (netBeforeTax) |
| ใครเลือก VAT/WHT | Agent เลือกตอนจอง |
| มัดจำคิดจากยอดอะไร | ยอดที่ต้องชำระ (grandTotal) |
| Default ค่า VAT/WHT | ปิดทั้งคู่ (false) |

---

## 2. ตัวอย่างการคำนวณ

### 2.1 ตัวอย่างตัวเลข

| รายการ | ไม่มี Tax | +VAT 7% | +WHT 3% | +VAT +WHT |
|--------|-----------|---------|---------|-----------|
| ยอดรวมก่อนหัก | 100,000 | 100,000 | 100,000 | 100,000 |
| ส่วนลด Agent | -5,000 | -5,000 | -5,000 | -5,000 |
| **ยอดก่อนภาษี** | **95,000** | **95,000** | **95,000** | **95,000** |
| VAT 7% | - | +6,650 | - | +6,650 |
| หัก ณ ที่จ่าย 3% | - | - | -2,850 | -2,850 |
| **ยอดที่ต้องชำระ** | **95,000** | **101,650** | **92,150** | **98,800** |

### 2.2 สูตรคำนวณ

```
grandTotal = netBeforeTax + (netBeforeTax × 0.07 × includeVAT) - (netBeforeTax × 0.03 × includeWHT)
```

---

## 3. การเปลี่ยนแปลงที่ทำเสร็จแล้ว

### 3.1 Backend (Go)

#### A. Booking Model
```go
// domain/models/booking.go

// VAT fields (ภาษีมูลค่าเพิ่ม)
IncludeVAT bool    `gorm:"not null;default:false;column:include_vat" json:"include_vat"`
VATRate    float64 `gorm:"type:decimal(5,2);not null;default:7.00;column:vat_rate" json:"vat_rate"`
VATAmount  float64 `gorm:"type:decimal(10,2);not null;default:0;column:vat_amount" json:"vat_amount"`

// Withholding Tax fields (หัก ณ ที่จ่าย)
IncludeWHT bool    `gorm:"not null;default:false;column:include_wht" json:"include_wht"`
WHTRate    float64 `gorm:"type:decimal(5,2);not null;default:3.00;column:wht_rate" json:"wht_rate"`
WHTAmount  float64 `gorm:"type:decimal(10,2);not null;default:0;column:wht_amount" json:"wht_amount"`

// Calculated totals
NetBeforeTax float64 `gorm:"type:decimal(10,2);not null;default:0;column:net_before_tax" json:"net_before_tax"`
GrandTotal   float64 `gorm:"type:decimal(10,2);not null;default:0;column:grand_total" json:"grand_total"`
```

#### B. Booking DTO
```go
// domain/dto/booking.go

// CreateBookingRequest
IncludeVAT bool `json:"include_vat"` // รวม VAT 7% (default: false)
IncludeWHT bool `json:"include_wht"` // หัก ณ ที่จ่าย 3% (default: false)

// BookingResponse
IncludeVAT   bool    `json:"include_vat"`
VATRate      float64 `json:"vat_rate"`
VATAmount    float64 `json:"vat_amount"`
IncludeWHT   bool    `json:"include_wht"`
WHTRate      float64 `json:"wht_rate"`
WHTAmount    float64 `json:"wht_amount"`
NetBeforeTax float64 `json:"net_before_tax"`
GrandTotal   float64 `json:"grand_total"`
```

#### C. Booking Service - Calculation Logic
```go
// application/serviceimpl/booking_service_impl.go

// Calculate VAT and WHT based on request options
var vatRate, vatAmount, whtRate, whtAmount float64
vatRate = 7.00  // Default VAT rate
whtRate = 3.00  // Default WHT rate

if req.IncludeVAT {
    vatAmount = netBeforeTax * (vatRate / 100)
}
if req.IncludeWHT {
    whtAmount = netBeforeTax * (whtRate / 100)
}

grandTotal := netBeforeTax + vatAmount - whtAmount

// เงินมัดจำคิดจาก grandTotal (ยอดที่ต้องชำระ)
if period.DepositType == "percent" {
    depositAmount = grandTotal * (float64(*period.Deposit) / 100)
} else if period.DepositType == "fixed" {
    depositAmount = float64(*period.Deposit) * float64(sumQty)
}
```

### 3.2 Frontend (React)

#### A. TypeScript Types
```typescript
// types/booking.ts

interface BookingResponse {
  // ... existing fields ...

  // VAT fields (ภาษีมูลค่าเพิ่ม)
  include_vat: boolean;
  vat_rate: number;
  vat_amount: number;

  // Withholding Tax fields (หัก ณ ที่จ่าย)
  include_wht: boolean;
  wht_rate: number;
  wht_amount: number;

  // Calculated totals
  net_before_tax: number;
  grand_total: number;
}

interface CreateBookingRequest {
  // ... existing fields ...
  include_vat?: boolean;
  include_wht?: boolean;
}
```

#### B. Booking Form (BookingSheet.tsx)
- เพิ่ม state: `includeVAT`, `includeWHT`
- เพิ่ม Switch toggles สำหรับเลือก VAT/WHT
- เพิ่มฟังก์ชัน: `calculateVATAmount()`, `calculateWHTAmount()`, `calculateGrandTotal()`
- อัพเดต price summary แสดง VAT, WHT, และ Grand Total
- อัพเดต handleSubmit ส่ง `include_vat`, `include_wht` ไปใน request

#### C. Invoice PDF (pdfInvoiceNew.tsx)
- อัพเดต Summary section แสดง:
  - ยอดก่อนภาษี (net_before_tax)
  - VAT 7% (ถ้าเลือก)
  - หัก ณ ที่จ่าย 3% (ถ้าเลือก)
  - ยอดที่ต้องชำระ (grand_total)
- อัพเดต "จำนวนเงินตัวอักษร" ใช้ grand_total
- อัพเดต payment schedule ใช้ grand_total

---

## 4. Files ที่แก้ไข

### Backend:
| File | Changes |
|------|---------|
| `domain/models/booking.go` | เพิ่ม VAT/WHT fields, อัพเดต IsPaid(), RemainingBalance() |
| `domain/dto/booking.go` | เพิ่ม VAT/WHT fields ใน Request/Response |
| `application/serviceimpl/booking_service_impl.go` | เพิ่ม calculation logic, อัพเดต toBookingResponse() |

### Frontend:
| File | Changes |
|------|---------|
| `types/booking.ts` | เพิ่ม VAT/WHT TypeScript types |
| `components/tours-regular/BookingSheet.tsx` | เพิ่ม toggles, calculation functions, และ UI |
| `components/globals/pdfInvoiceNew.tsx` | อัพเดต Invoice PDF แสดง VAT/WHT |

---

## 5. Database Migration

GORM จะ auto-migrate fields ใหม่เมื่อ start server (fields จะถูกเพิ่มอัตโนมัติ):

```sql
-- Fields ใหม่ที่ถูกเพิ่มใน bookings table
include_vat BOOLEAN DEFAULT FALSE
vat_rate DECIMAL(5,2) DEFAULT 7.00
vat_amount DECIMAL(10,2) DEFAULT 0

include_wht BOOLEAN DEFAULT FALSE
wht_rate DECIMAL(5,2) DEFAULT 3.00
wht_amount DECIMAL(10,2) DEFAULT 0

net_before_tax DECIMAL(10,2) DEFAULT 0
grand_total DECIMAL(10,2) DEFAULT 0
```

---

## 6. การทดสอบ

### 6.1 Test Scenarios

1. **ไม่เลือก VAT, ไม่เลือก WHT (default)**
   - grandTotal = netBeforeTax

2. **เลือก VAT เท่านั้น**
   - grandTotal = netBeforeTax + (netBeforeTax × 0.07)

3. **เลือก WHT เท่านั้น**
   - grandTotal = netBeforeTax - (netBeforeTax × 0.03)

4. **เลือกทั้ง VAT และ WHT**
   - grandTotal = netBeforeTax + (netBeforeTax × 0.07) - (netBeforeTax × 0.03)

### 6.2 Invoice PDF Test
- ตรวจสอบว่า VAT/WHT แสดงถูกต้องใน PDF
- ตรวจสอบจำนวนเงินตัวอักษร
- ตรวจสอบยอดมัดจำและยอดส่วนที่เหลือ

---

**สร้างเมื่อ:** 2026-03-03
**อัพเดตล่าสุด:** 2026-03-03
**สถานะ:** ✅ เสร็จสมบูรณ์
