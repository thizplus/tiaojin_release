# Migration System - Important Reference

## Quick Commands

```bash
# First time migration (clear all + sync all)
go run . clear
go run . sync-full

# Daily sync (incremental - won't overwrite existing data)
go run . sync-latest

# Sync specific time range
go run . sync-latest --hours 24
go run . sync-latest --since 2025-01-15

# After sync-full, migrate files (images)
go run . files all
```

---

## Data Flow: MongoDB → PostgreSQL

```
MongoDB                          PostgreSQL
─────────────────────────────────────────────────────────────
tbl_city          ───────────►   cities
tbl_banner        ───────────►   banners
tbl_sales         ───────────►   sales
tbl_users         ───────────►   users + company_profiles
tbl_tours         ───────────►   tours
tbl_period        ───────────►   periods
tbl_booking       ───────────►   bookings
                                  ├── booking_room_types
                                  ├── booking_rate_types
                                  ├── payment_slips
                                  └── passenger_passports
```

---

## Relationships Diagram

```
┌─────────┐
│  City   │
└────┬────┘
     │ 1:N
     ▼
┌─────────┐       ┌─────────┐
│  Tour   │◄──────│ Booking │
└────┬────┘  N:1  └────┬────┘
     │ 1:N             │
     ▼                 │
┌─────────┐            │
│ Period  │◄───────────┘ N:1
└─────────┘

┌─────────┐       ┌─────────────────┐
│  User   │──────►│ CompanyProfile  │ 1:1
└────┬────┘       └─────────────────┘
     │ 1:N
     ▼
┌─────────┐
│ Booking │
└────┬────┘
     │ 1:N
     ▼
┌─────────┐       ┌─────────┐
│  Sales  │◄──────│ Booking │ N:1
└─────────┘       └─────────┘
```

---

## Foreign Key Relationships

| Table | Column | References | Required |
|-------|--------|------------|----------|
| tours | city_id | cities.id | Optional |
| periods | tour_id | tours.id | Required |
| bookings | user_id | users.id | Required |
| bookings | tour_id | tours.id | Required |
| bookings | period_id | periods.id | Required |
| bookings | sales_id | sales.id | Required |
| company_profiles | user_id | users.id | Required |
| sales | user_id | users.id | Optional |
| booking_room_types | booking_id | bookings.id | Required |
| booking_rate_types | booking_id | bookings.id | Required |
| payment_slips | booking_id | bookings.id | Required |
| passenger_passports | booking_id | bookings.id | Required |

---

## Sync Behavior - IMPORTANT!

### Sync Mode: UPDATE (Not Skip!)
ข้อมูลจะถูก **UPDATE** ไม่ใช่ SKIP → Status จะ realtime เสมอ

### Fields Updated vs Preserved

| Table | Updated Fields | Protected Fields |
|-------|---------------|------------------|
| cities | detail, status | **img** (ถ้ามีอยู่แล้ว) |
| banners | status | **img** |
| tours | title, city_id, detail, status | **thumbnail, file_pdf, file_word** (ถ้ามีอยู่แล้ว) |
| periods | seat, status | all JSONB fields, file paths |
| bookings | **status only** | all other fields |
| users | role, status | password, email |
| sales | first_name, last_name, nick_name, active | - |

### Sync Logic

```
If record EXISTS in PostgreSQL:
  → UPDATE specific fields (status changes will sync!)
  → PRESERVE file paths (won't overwrite uploaded images)

If record NOT EXISTS:
  → INSERT new record with all data
```

### Example Output (Second Sync)
```
Cities:   22 total, 0 inserted, 22 updated   ← UPDATE not SKIP
Banners:  5 total, 0 inserted, 5 updated
Tours:    324 total, 0 inserted, 324 updated
Bookings: 14776 updated ← Status updates work!
```

---

## MongoDB Field Mappings

### Period (tbl_period)
```
MongoDB                      PostgreSQL
─────────────────────────────────────────
tour (string hex)        →   tour_id (UUID)
start_date (interface{}) →   start_date (date)
end_date (interface{})   →   end_date (date)
seat                     →   seat
deposit                  →   deposit
airline                  →   airline
airport                  →   airport
note                     →   note
sector1[]                →   sector_1 (JSONB)
sector2[]                →   sector_2 (JSONB)
rate.type                →   rate_type
rate.item[]              →   rate_items (JSONB)
commission.type          →   commission_type
commission.item[]        →   commission_items (JSONB)
status                   →   status
```

### Booking (tbl_booking)
```
MongoDB                      PostgreSQL
─────────────────────────────────────────
invoiceId                →   invoice_id (unique)
profileId (ObjectID)     →   user_id (UUID)
tour (ObjectID)          →   tour_id (UUID)
periodCode (ObjectID)    →   period_id (UUID)
salesCode (ObjectID)     →   sales_id (UUID)
cusName                  →   customer_name
cusTel                   →   customer_phone
roomTypes[]              →   booking_room_types table
rateTypes[]              →   booking_rate_types table
slip[]                   →   payment_slips table
passport[]               →   passenger_passports table
```

### Sales (tbl_sales)
```
MongoDB                      PostgreSQL
─────────────────────────────────────────
code                     →   code (unique)
firstName                →   first_name
lastName                 →   last_name
nickName                 →   nick_name
status                   →   active (boolean)
```

---

## Fix Scripts (Run Automatically After Sync)

| Order | Script | Purpose |
|-------|--------|---------|
| 1 | FixTourCities | Match tour.city_id by city title |
| 2 | FixPeriodTours | Match period.tour_id by tour code |
| 3 | FixUserCompany | Create missing company_profiles |
| 4 | FixUserStatus | Fix user status mapping |
| 5 | FixBookingUsers | Match booking.user_id by profile |
| 6 | FixBookingTours | Match booking.tour_id by tour code |
| 7 | FixBookingPeriods | Match booking.period_id by tour + dates |
| 8 | FixBookingSales | Match booking.sales_id by sales code |
| 9 | FixSlipAmounts | Calculate slip amounts from booking.sumPayment |

### Important Note: Payment Slips
MongoDB stores `slip.amount = 0` but `booking.sumPayment` has the correct value.
`FixSlipAmounts` calculates amount per slip = sumPayment / confirmed_slips_count

---

## File Migration

After data migration, run file migration to copy images from DO Spaces to iDrive e2:

```bash
go run . files all        # All files
go run . files companies  # Company logos, ATTA, DBD, PP20
go run . files banners    # Banner images
go run . files cities     # City images
go run . files tours      # Tour thumbnails, PDFs, Word files
go run . files slips      # Payment slip images
go run . files passports  # Passport images
```

---

## Troubleshooting

### Problem: Bookings not syncing
- Check if tours and periods exist first
- Run `go run . sync-full` to ensure correct order

### Problem: Period-Tour relationships broken
- Period.TourID in MongoDB is a **string hex**, not ObjectID
- Fix script handles this automatically

### Problem: Images overwritten after sync
- This should NOT happen - sync preserves existing paths
- Check if PostgreSQL has existing image paths before sync

### Problem: User has no company profile
- FixUserCompany creates missing profiles
- Runs automatically in sync-full

---

## Environment Variables

```env
# MongoDB
MONGODB_URI=mongodb://localhost:27017

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=tiaojin
DB_SSL_MODE=disable

# File Migration (optional)
DO_SPACES_KEY=xxx
DO_SPACES_SECRET=xxx
DO_SPACES_ENDPOINT=xxx
DO_SPACES_BUCKET=xxx
IDRIVE_ACCESS_KEY=xxx
IDRIVE_SECRET_KEY=xxx
IDRIVE_ENDPOINT=xxx
IDRIVE_BUCKET=xxx
```

---

## Last Updated
- Date: 2026-03-02
- Author: Migration System Documentation
