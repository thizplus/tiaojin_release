# Period-Level Expiration Override

## สรุป Feature

ระบบ **Override การหมดอายุ Invoice ระดับ Period** สำหรับทัวร์ไฟไหม้หรือ Period พิเศษที่ต้องการกำหนดเวลาหมดอายุเอง

---

## ปัญหาที่ต้องการแก้

ปัจจุบันระบบใช้ **Global Expiration Tiers** ที่คำนวณจากจำนวนวันก่อนทัวร์เริ่ม:
- Critical (≤3 วัน) → 30 นาที
- Urgent (≤7 วัน) → 6 ชั่วโมง
- etc.

แต่บาง Period ต้องการ **กำหนดเวลาเอง** เช่น:
- ทัวร์ไฟไหม้ลดราคา 70% → ต้องยืนยันภายใน **15 นาที**
- Flash sale → หมดอายุใน **1 ชั่วโมง**
- Period VIP → ให้เวลา **7 วัน** (นานกว่าปกติ)

---

## Solution: Period Expiration Override

### แนวคิด

เพิ่ม field ใน `periods` table ที่สามารถ **override** การคำนวณ expiration ของ Global Tiers

```
Booking Created
      ↓
Check: Period มี expiration_override?
      ↓
  [YES] → ใช้ค่า override
  [NO]  → ใช้ Global Tiers ตามปกติ
```

---

## Database Changes

### Option A: เพิ่ม Fields ใน periods table (แนะนำ)

```sql
ALTER TABLE periods ADD COLUMN expiration_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE periods ADD COLUMN expiration_minutes INTEGER;
ALTER TABLE periods ADD COLUMN expiration_hours INTEGER;
ALTER TABLE periods ADD COLUMN expiration_days INTEGER;
ALTER TABLE periods ADD COLUMN expiration_message VARCHAR(255);
```

| Field | Type | Description |
|-------|------|-------------|
| `expiration_enabled` | BOOLEAN | เปิด/ปิด override (default: false = ใช้ global tiers) |
| `expiration_minutes` | INTEGER | หมดอายุในกี่นาที (ใช้อันใดอันหนึ่ง) |
| `expiration_hours` | INTEGER | หมดอายุในกี่ชั่วโมง |
| `expiration_days` | INTEGER | หมดอายุในกี่วัน |
| `expiration_message` | VARCHAR | ข้อความแสดง เช่น "ทัวร์ไฟไหม้! ต้องยืนยันภายใน 15 นาที" |

### Option B: JSONB (ยืดหยุ่นกว่า)

```sql
ALTER TABLE periods ADD COLUMN expiration_override JSONB;
```

```json
{
  "enabled": true,
  "expire_minutes": 15,
  "message": "ทัวร์ไฟไหม้! ต้องยืนยันภายใน 15 นาที",
  "notify_before_minutes": [5, 1]
}
```

---

## Model Changes

### periods model

```go
type Period struct {
    // ... existing fields ...

    // Expiration Override
    ExpirationEnabled  bool    `gorm:"default:false" json:"expiration_enabled"`
    ExpirationMinutes  *int    `json:"expiration_minutes,omitempty"`
    ExpirationHours    *int    `json:"expiration_hours,omitempty"`
    ExpirationDays     *int    `json:"expiration_days,omitempty"`
    ExpirationMessage  *string `json:"expiration_message,omitempty"`
}

// HasExpirationOverride checks if period has custom expiration
func (p *Period) HasExpirationOverride() bool {
    return p.ExpirationEnabled && (p.ExpirationMinutes != nil || p.ExpirationHours != nil || p.ExpirationDays != nil)
}

// CalculateExpirationTime calculates expiration from period override
func (p *Period) CalculateExpirationTime(createdAt time.Time) *time.Time {
    if p.ExpirationMinutes != nil {
        t := createdAt.Add(time.Duration(*p.ExpirationMinutes) * time.Minute)
        return &t
    }
    if p.ExpirationHours != nil {
        t := createdAt.Add(time.Duration(*p.ExpirationHours) * time.Hour)
        return &t
    }
    if p.ExpirationDays != nil {
        t := createdAt.Add(time.Duration(*p.ExpirationDays) * 24 * time.Hour)
        return &t
    }
    return nil
}
```

---

## Service Changes

### booking_service_impl.go

```go
func (s *BookingServiceImpl) CalculateBookingExpiration(ctx context.Context, periodID uuid.UUID, createdAt time.Time) (*time.Time, error) {
    // Get period
    period, err := s.periodRepo.GetByID(ctx, periodID)
    if err != nil {
        return nil, err
    }

    // Check for period-level override
    if period.HasExpirationOverride() {
        log.Printf("[DEBUG] Using period expiration override for period %s", periodID)
        return period.CalculateExpirationTime(createdAt), nil
    }

    // Fall back to global tiers
    log.Printf("[DEBUG] Using global expiration tiers for period %s", periodID)
    // ... existing tier calculation logic ...
}
```

---

## API Changes

### PUT /api/v1/periods/:id

เพิ่ม fields ใน request:

```json
{
    "start_date": "2026-04-15",
    "end_date": "2026-04-20",
    "seat": 30,
    "deposit": 5000,

    // NEW: Expiration Override
    "expiration_enabled": true,
    "expiration_minutes": 15,
    "expiration_message": "ทัวร์ไฟไหม้! กรุณายืนยันภายใน 15 นาที"
}
```

### GET /api/v1/periods/:id

Response จะมี fields เพิ่ม:

```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "2026-04-15",
    "seat": 30,

    // Expiration Override Info
    "expiration_enabled": true,
    "expiration_minutes": 15,
    "expiration_message": "ทัวร์ไฟไหม้! กรุณายืนยันภายใน 15 นาที",
    "expiration_display": "15 นาที"
}
```

---

## Frontend Changes

### Period Edit Form

เพิ่ม Section "การหมดอายุ Invoice":

```
┌─────────────────────────────────────────────────────┐
│ การหมดอายุ Invoice                                   │
├─────────────────────────────────────────────────────┤
│ ☐ ใช้ค่าเริ่มต้นจากระบบ (Global Tiers)              │
│ ☑ กำหนดเอง (Override)                               │
│                                                     │
│ หมดอายุใน: [15] [นาที ▼]                            │
│                                                     │
│ ข้อความแสดง:                                        │
│ [ทัวร์ไฟไหม้! กรุณายืนยันภายใน 15 นาที          ]   │
└─────────────────────────────────────────────────────┘
```

### Period List

แสดง badge สำหรับ period ที่มี custom expiration:

```
| Tour Code  | Start Date | Seat | Expiration     |
|------------|------------|------|----------------|
| BCKMG48-MU | 2026-04-15 | 30   | 🔥 15 นาที     |
| BCKMG48-MU | 2026-04-22 | 30   | (ใช้ค่าระบบ)   |
```

---

## Migration

```sql
-- Add expiration override fields to periods
ALTER TABLE periods ADD COLUMN IF NOT EXISTS expiration_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE periods ADD COLUMN IF NOT EXISTS expiration_minutes INTEGER;
ALTER TABLE periods ADD COLUMN IF NOT EXISTS expiration_hours INTEGER;
ALTER TABLE periods ADD COLUMN IF NOT EXISTS expiration_days INTEGER;
ALTER TABLE periods ADD COLUMN IF NOT EXISTS expiration_message VARCHAR(255);

-- Create index for quick lookup
CREATE INDEX IF NOT EXISTS idx_periods_expiration_enabled ON periods(expiration_enabled) WHERE expiration_enabled = true;
```

---

## Priority Flow

```
1. Period Override (ถ้า enabled)
   ↓ (ถ้าไม่มี)
2. Global Expiration Tiers (ตาม days until tour)
   ↓ (ถ้าไม่มี tier match)
3. Default Tier (3 วัน)
```

---

## Use Cases

### Case 1: ทัวร์ไฟไหม้ 70% OFF

```
Period: BCKMG48-MU / 15 เม.ย. 2026
Setting: expiration_minutes = 15
Message: "🔥 ทัวร์ไฟไหม้! ยืนยันภายใน 15 นาที"

→ ทุก booking ใน period นี้ หมดอายุใน 15 นาที
→ ไม่สนใจว่าทัวร์จะเริ่มอีกกี่วัน
```

### Case 2: Flash Sale 6 ชั่วโมง

```
Period: BCKMG48-MU / 22 เม.ย. 2026
Setting: expiration_hours = 6
Message: "Flash Sale! ยืนยันภายใน 6 ชั่วโมง"

→ ทุก booking หมดอายุใน 6 ชั่วโมง
```

### Case 3: VIP Period (เวลานานกว่าปกติ)

```
Period: BCVIP-JAPAN / 1 พ.ค. 2026
Setting: expiration_days = 7
Message: "VIP Package - มีเวลายืนยัน 7 วัน"

→ ทุก booking หมดอายุใน 7 วัน (แม้ว่า global tier จะกำหนดน้อยกว่า)
```

### Case 4: ปกติ (ไม่ override)

```
Period: BCKMG48-MU / 30 เม.ย. 2026
Setting: expiration_enabled = false

→ ใช้ Global Tiers ตามปกติ
→ ถ้าทัวร์เริ่มอีก 26 วัน → Standard tier → 2 วัน
```

---

## สรุปความเป็นไปได้

| หัวข้อ | ความเป็นไปได้ | หมายเหตุ |
|--------|--------------|----------|
| Database | ✅ ง่าย | เพิ่ม 5 columns |
| Backend | ✅ ง่าย | แก้ไข CalculateBookingExpiration |
| Frontend | ✅ ปานกลาง | เพิ่ม form fields ใน Period Edit |
| Migration | ✅ ง่าย | ALTER TABLE |

**Estimated Effort**: 2-4 ชั่วโมง

---

## ขั้นตอนการ Implement

1. **Database Migration** - เพิ่ม columns ใน periods
2. **Model Update** - เพิ่ม fields ใน Period model
3. **DTO Update** - เพิ่ม fields ใน PeriodRequest/Response
4. **Service Update** - แก้ไข CalculateBookingExpiration
5. **Handler Update** - รับ/ส่ง expiration fields
6. **Frontend** - เพิ่ม UI สำหรับตั้งค่า

---

## คำถาม

1. ต้องการ validation อะไรบ้าง? เช่น ต้องระบุอย่างน้อย 1 อย่าง (minutes/hours/days)?
2. ต้องการ notify ก่อนหมดอายุไหม? (เช่น แจ้งเตือน 5 นาทีก่อนหมด)
3. ต้องการ log/audit trail สำหรับการเปลี่ยนแปลง expiration override ไหม?
