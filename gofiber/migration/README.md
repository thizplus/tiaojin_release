# MongoDB to PostgreSQL Migration Tool

This tool migrates data from the old MongoDB-based tour booking system to the new PostgreSQL-based system.

## Prerequisites

1. Go 1.21 or later installed
2. Access to MongoDB (old system database)
3. Access to PostgreSQL (new system database)

## Setup

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Update `.env` with your actual credentials:
   - `MONGODB_URI`: Your MongoDB connection string from the old system
   - PostgreSQL credentials matching your new system configuration

## Usage

### Run Full Migration

```bash
cd migration
go mod tidy
go run .
```

This will:
1. Connect to both databases
2. Clear existing PostgreSQL data (mock data)
3. Migrate all collections in the correct order
4. Verify the migration counts

### Commands

```bash
# Run full migration
go run .

# Clear PostgreSQL data only
go run . clear

# Run verification only
go run . verify

# Show help
go run . help
```

## Migration Order

The migration runs in this order to respect foreign key dependencies:

1. **Cities** (tbl_city -> cities)
2. **Banners** (tbl_banner -> banners)
3. **Sales** (tbl_sales -> sales)
4. **Users** (tbl_users -> users + company_profiles)
5. **Tours** (tbl_tours -> tours)
6. **Periods** (tbl_period -> periods)
7. **Bookings** (tbl_booking -> bookings + child tables)

## Data Transformations

### Email Normalization
- Emails are converted to lowercase
- Leading/trailing spaces are removed
- Internal spaces are removed
- Duplicate emails after normalization are skipped

### Status Mappings
- "active", "publish" -> "publish"
- "inactive", "draft", "" -> "draft"

### Date Handling
- Supports multiple date formats (ISO 8601, MongoDB $date, etc.)
- Falls back to current date if parsing fails

### ID Mapping
- MongoDB ObjectIDs are converted to PostgreSQL UUIDs
- A mapping table is maintained during migration for foreign key relationships

## Child Tables Created

For bookings, the following child tables are populated:
- `booking_room_types` - Room type selections
- `booking_rate_types` - Rate type selections with pricing
- `payment_slips` - Payment slip records
- `passenger_passports` - Passenger passport information
- `booking_additional_services` - Additional services (if any)

## Verification

After migration, a verification report shows:
- Record counts in MongoDB vs PostgreSQL
- Any discrepancies found
- Child table record counts

## Troubleshooting

### "User not found" warnings
This means a user referenced in a booking was not migrated (possibly due to duplicate email). Check the email normalization output.

### "Period not found" warnings
The tour or period referenced doesn't exist in the database. Check if the tour was migrated successfully.

### Connection errors
- Verify MongoDB URI is correct and accessible
- Verify PostgreSQL credentials match your setup
- Check network connectivity

## Safety Notes

- This tool only READS from MongoDB (no modifications)
- PostgreSQL data will be OVERWRITTEN (existing data cleared first)
- Run on a backup/test database first to verify results
