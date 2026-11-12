# Migration Tool - Change Log

## Version 1.4.0 (2026-03-04)

### สรุปการแก้ไข

#### 1. เพิ่ม Auto-Seed สำหรับ System Configurations

**ปัญหา**: เมื่อ clear database และ sync ใหม่ ไม่มี expiration tiers config ทำให้ระบบ expiration ไม่ทำงาน

**ไฟล์ที่เพิ่ม/แก้ไข**:
- `seed_configs.go` - ไฟล์ใหม่สำหรับ seed default configs
- `sync_full.go` - เรียก seed functions ก่อน sync data

**การทำงาน**:
- ตรวจสอบว่า config มีอยู่หรือไม่ก่อน INSERT (ไม่ overwrite ค่าที่มีอยู่)
- Seed ทั้ง `booking.expiration_tiers` และ `booking.agent_period_limits`
- Seed default job `backup-booking-expiry-cleanup`

**Default Configs ที่ถูก Seed**:

| Config Key | Description |
|------------|-------------|
| `booking.expiration_tiers` | 5 tiers: Critical (30min), Urgent (6h), Normal (24h), Standard (2d), Default (3d) |
| `booking.agent_period_limits` | Agent booking limits per period (disabled by default) |

**Default Job ที่ถูก Seed**:

| Job Name | Cron | Description |
|----------|------|-------------|
| `backup-booking-expiry-cleanup` | `0 */1 * * *` | Backup cleanup ทุกชั่วโมง |

**วิธีใช้**:
```bash
# Clear และ Sync ใหม่ - configs จะถูก seed อัตโนมัติ
go run . clear
go run . sync-full
```

---

#### 2. แก้ไข Config Handler ให้เรียก BookingService จริง

**ปัญหา**: `GetExpirationStats` และ `ProcessExpiredBookings` ใน config_handler.go เป็น mock response

**ไฟล์ที่แก้ไข**:
- `interfaces/api/handlers/config_handler.go` - เพิ่ม bookingService dependency
- `interfaces/api/handlers/handlers.go` - ส่ง bookingService ให้ ConfigHandler

**การแก้ไข**:
```go
// Before: Mock response
response := dto.ExpirationStatsResponse{
    TotalBookingsWithExpiration: 0,
    // ... hardcoded values
}

// After: Call actual service
response, err := h.bookingService.GetBookingExpirationStats(c.Context())
```

---

## Version 1.3.0 (2026-03-04)

### สรุปการแก้ไข

#### 1. แก้ไขปัญหา rate_items และ room_types ไม่แสดงบน Frontend (CRITICAL FIX)

**ปัญหา**: ระบบใหม่ใช้ **separate tables** (`period_rate_items`, `period_room_types`) แทน JSONB columns
แต่ migration เดิมใส่ข้อมูลลง JSONB columns ทำให้ Frontend ไม่แสดงข้อมูล

**สาเหตุ**:
- Service layer ใช้ `toPeriodResponseWithSeparateTables()` ที่อ่านจาก `period_rate_items` และ `period_room_types` tables
- JSONB columns (`rate_items`, `room_types`) ถูกตั้งค่าเป็น empty arrays โดย design

**ไฟล์ที่แก้ไข**:
- `migrate_periods.go` - เพิ่ม `insertPeriodRateItems()` และ `insertPeriodRoomTypes()` functions
- `sync_full.go` - ใช้ empty JSONB และ INSERT เข้า separate tables

**การแก้ไข**:

1. เพิ่ม function สำหรับ INSERT เข้า `period_rate_items`:
```go
func insertPeriodRateItems(ctx context.Context, periodID uuid.UUID, items []MongoRateItem, createdAt time.Time) error {
    query := `INSERT INTO period_rate_items (id, period_id, title, regular, sale, exclude, sort_order, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
    for i, item := range items {
        // INSERT each rate item with sort_order
    }
}
```

2. เพิ่ม function สำหรับ INSERT เข้า `period_room_types`:
```go
func insertPeriodRoomTypes(ctx context.Context, periodID uuid.UUID, singleStay float64, createdAt time.Time) error {
    // INSERT 4 standard room types: TWIN, DOUBLE, TRIPLE, SINGLE
    // SINGLE gets price from single_stay
}
```

3. เปลี่ยน JSONB columns เป็น empty arrays:
```go
"[]", // rate_items - empty JSONB, data goes to period_rate_items table
"[]", // room_types - empty JSONB, data goes to period_room_types table
```

**ผลลัพธ์**:
- rate_items → INSERT เข้า `period_rate_items` table
- room_types → INSERT เข้า `period_room_types` table (4 standard types)
- JSONB columns เป็น empty arrays (ตาม design ของระบบใหม่)

**Tables ที่ถูกใช้**:
| Data | Table | Fields |
|------|-------|--------|
| Rate Items | `period_rate_items` | id, period_id, title, regular, sale, exclude, sort_order |
| Room Types | `period_room_types` | id, period_id, title, price, sort_order |

---

## Version 1.2.0 (2026-03-04)

### สรุปการแก้ไข

#### 1. แก้ไขปัญหา room_types ว่างเปล่า (SUPERSEDED by v1.3.0)

**หมายเหตุ**: การแก้ไขนี้ถูกแทนที่ด้วย v1.3.0 ที่ใช้ separate tables แทน JSONB

~~**ปัญหา**: MongoDB มี `single_stay` (ราคาพักเดี่ยว) แต่ไม่ได้ migrate ไปยัง `room_types` ใน PostgreSQL~~

---

#### 2. แก้ไขปัญหา deposit_type ผิด

**ปัญหา**: ใช้ค่า `"fixed"` แต่ Backend validation รองรับแค่ `"percent"` หรือ `"per_unit"`

**ไฟล์ที่แก้ไข**:
- `migrate_periods.go`
- `sync_full.go`

**การแก้ไข**:
```go
// Before: depositType = "fixed"    ❌
// After:  depositType = "per_unit" ✅
depositType := "percent"
if period.Deposit > 100 {
    depositType = "per_unit"
}
```

---

#### 3. แก้ไขปัญหา Sector Date Format ผิด

**ปัญหา**: MongoDB เก็บ date เป็น `"DD/MM/YYYY"` แต่ Frontend ต้องการ `"YYYY-MM-DD"`

**ไฟล์ที่แก้ไข**:
- `migrate_periods.go` - เพิ่ม `convertDateFormat()` function

**การแก้ไข**:
```go
// convertDateFormat converts date from DD/MM/YYYY to YYYY-MM-DD
func convertDateFormat(date string) string {
    parts := strings.Split(date, "/")
    if len(parts) == 3 {
        return fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])
    }
    return date
}
```

**ผลลัพธ์**:
- MongoDB `"02/06/2026"` → PostgreSQL `"2026-06-02"`

---

#### 4. แก้ไขปัญหา Rate Items exclude ไม่ถูก migrate

**ปัญหา**: `exclude` field ถูก hardcode เป็น `false` แทนที่จะอ่านจาก MongoDB

**ไฟล์ที่แก้ไข**:
- `migrate_periods.go`

**การแก้ไข**:
```go
// Before
type MongoRateItem struct {
    Title   interface{} `bson:"title"`
    Regular interface{} `bson:"regular"`
    Sale    interface{} `bson:"sale"`
}
// Exclude: false  ← hardcode

// After
type MongoRateItem struct {
    Title   interface{} `bson:"title"`
    Regular interface{} `bson:"regular"`
    Sale    interface{} `bson:"sale"`
    Exclude bool        `bson:"exclude"`  // ← เพิ่ม field
}
// Exclude: item.Exclude  ← อ่านจาก MongoDB
```

---

## Version 1.1.0 (2026-03-03)

### สรุปการแก้ไข

#### 1. แก้ไขปัญหา Timezone - วันที่ผิด

**ปัญหา**: MongoDB เก็บวันที่ในรูปแบบ UTC (เช่น `2026-03-26T17:00:00.000Z`)
แต่ความจริงคือวันที่ 27 มีนาคม ตามเวลาไทย (UTC+7)

**ไฟล์ที่แก้ไข**:
- `sync_full.go` - เพิ่ม timezone conversion ใน `upsertPeriodFull()`
- `migrate_periods.go` - เพิ่ม timezone conversion ใน `MigratePeriods()`

**การแก้ไข**:
```go
// Convert to Thailand timezone (UTC+7) before extracting date
thailandTZ := time.FixedZone("Asia/Bangkok", 7*60*60)
startDateTH := startDate.In(thailandTZ)
endDateTH := endDate.In(thailandTZ)

// Normalize dates to DATE only
startDateOnly := time.Date(startDateTH.Year(), startDateTH.Month(), startDateTH.Day(), 0, 0, 0, 0, time.UTC)
endDateOnly := time.Date(endDateTH.Year(), endDateTH.Month(), endDateTH.Day(), 0, 0, 0, 0, time.UTC)
```

**ผลลัพธ์**:
- MongoDB `2026-03-26T17:00:00Z` → PostgreSQL `2026-03-27` (ถูกต้อง)

---

#### 2. แก้ไขปัญหา Period Duplicate

**ปัญหา**: `FixPeriodTours` ใน `fix_all.go` ย้าย periods จากทัวร์อื่นมาเนื่องจาก match ด้วย start_date/end_date อย่างเดียว

**ไฟล์ที่แก้ไข**:
- `fix_all.go` - ลบ `FixPeriodTours` ออกจาก fix list
- `fix_period_tours.go` - เพิ่ม timezone conversion (แต่ไม่ใช้แล้ว)

**การแก้ไข**:
```go
// REMOVED: {"Period-Tour Relationships", FixPeriodTours} - handled by sync
```

**เหตุผล**: `upsertPeriodFull` ใน sync-full ตั้ง tour_id ถูกต้องอยู่แล้วจาก MongoDB mapping

---

#### 3. แก้ไขปัญหา Files ไม่ถูก Migrate

**ปัญหา**: ไฟล์แนบ (files) จาก MongoDB ไม่ถูก migrate มาที่ PostgreSQL

**ไฟล์ที่แก้ไข**:
- `migrate_periods.go` - เพิ่ม `MongoFile` struct และ `convertFilesToJSON()` function
- `sync_full.go` - ใช้ `filesJSON` แทน `"[]"`

**การแก้ไข**:

1. เพิ่ม struct ใน `migrate_periods.go`:
```go
type MongoFile struct {
    Name string `bson:"name"`
    Type string `bson:"type"`
    URL  string `bson:"url"`
}
```

2. เพิ่ม field ใน `MongoPeriod`:
```go
Files []MongoFile `bson:"files"`
```

3. เพิ่ม function convert:
```go
func convertFilesToJSON(files []MongoFile) string {
    // ... convert to JSON
}
```

4. ใช้ใน INSERT:
```go
filesJSON, // files from MongoDB (แทน "[]")
```

---

### คำสั่งที่แนะนำ

```bash
# Clear และ Sync ใหม่ทั้งหมด
go run . clear
go run . sync-full

# ตรวจสอบ Period
go run . check-periods TOUR_CODE
go run . debug-period TOUR_CODE
```

---

### ข้อมูลที่ถูก Migrate

| Field | MongoDB | PostgreSQL | Status |
|-------|---------|------------|--------|
| start_date | UTC timestamp | DATE (Thailand) | ✅ Fixed |
| end_date | UTC timestamp | DATE (Thailand) | ✅ Fixed |
| sector_1 | Array | JSONB | ✅ OK |
| sector_2 | Array | JSONB | ✅ OK |
| deposit | Number | DECIMAL | ✅ OK |
| deposit_type | - | VARCHAR | ✅ Calculated |
| rate.item | Array | `period_rate_items` table | ✅ Fixed (v1.3.0) |
| commission | Object | JSONB | ✅ OK |
| files | Array | JSONB | ✅ Fixed |
| single_stay | Number | `period_room_types` table | ✅ Fixed (v1.3.0) - 4 standard types |
| additional_services | - | JSONB (empty) | ✅ Empty (not in MongoDB) |
| booking | Number | - | ⚠️ Calculated from bookings table |

---

### Architecture: JSONB vs Separate Tables

ระบบใหม่ใช้ **separate tables** สำหรับ rate_items และ room_types:

| JSONB Column (periods table) | Separate Table | Status |
|------------------------------|----------------|--------|
| `rate_items` (empty) | `period_rate_items` | ✅ Data here |
| `room_types` (empty) | `period_room_types` | ✅ Data here |
| `commission_items` | - | ✅ JSONB only |
| `additional_services` | `period_additional_services` | ⚠️ Empty (not in MongoDB) |

**เหตุผล**: Service layer ใช้ `toPeriodResponseWithSeparateTables()` ที่อ่านจาก separate tables

---

### หมายเหตุ

- `single_stay` และ `booking` ไม่ถูก migrate เพราะ:
  - `single_stay` - ระบบใหม่ใช้ `period_room_types` table แทน
  - `booking` - ระบบใหม่คำนวณจาก bookings table แบบ real-time

- JSONB columns (`rate_items`, `room_types`) ใน periods table จะเป็น empty arrays ตาม design ของระบบใหม่
