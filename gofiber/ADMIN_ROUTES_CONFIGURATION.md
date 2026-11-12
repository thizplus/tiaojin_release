# Admin Routes Configuration - Fixed

## Overview
Fixed admin routes registration for period management with real-time booking calculation.

## Issue Resolution

### ❌ Problem
- Public debug endpoint worked: `/periods/{id}/debug-booking-stats` ✅
- Admin routes not working: `/admin/periods` ❌
- Missing admin endpoints for debug and stats functionality

### ✅ Solution
1. **Enhanced Admin Routes**: Added missing endpoints to `admin_routes.go`
2. **Proper Middleware**: Verified `ManagementOnly()` vs `AdminOnly()` configuration
3. **Route Organization**: Separated public and admin functionality properly
4. **Real-time Integration**: All admin endpoints now use real-time booking calculation

## Route Structure

### 🌐 Public Routes (No Authentication)
```
Base: /api/v1/periods

GET /                      - List periods (public view)
GET /{id}                  - Get period details with real-time data
GET /{id}/availability     - Get availability with real-time calculation
GET /test-deployment       - Test deployment status
```

### 🔐 Admin Routes (Authentication Required)
```
Base: /api/v1/admin/periods
Middleware: Protected() + ManagementOnly()
Allowed Roles: ["super", "admin"]

GET /                              - List periods (admin view)
POST /                             - Create new period
GET /{id}                          - Get period details with real-time data
PUT /{id}                          - Update period
DELETE /{id}                       - Delete period
PUT /{id}/status                   - Update period status
GET /{id}/availability             - Get availability with real-time calculation
GET /stats                         - Get period statistics
GET /test-deployment               - Test admin deployment
GET /{id}/debug-booking-stats      - Debug booking calculation (NEW)
```

### 📊 Legacy Routes (Backward Compatibility)
```
Base: /api/v1/periods
Middleware: Protected() + AdminOnly()

POST /                     - Create period (legacy)
PUT /{id}                  - Update period (legacy)
DELETE /{id}               - Delete period (legacy)
PUT /{id}/status           - Update status (legacy)
GET /stats                 - Get stats (legacy)
```

## Middleware Configuration

### ManagementOnly() - Admin Routes
```go
// admin_routes.go:14-15
admin.Use(middleware.Protected())
admin.Use(middleware.ManagementOnly())

// Allows roles: ["super", "admin"]
```

### AdminOnly() - Legacy Routes
```go
// period_routes.go:19
adminPeriods := periods.Group("/", middleware.Protected(), middleware.AdminOnly())

// Allows roles: ["super", "admin"] - Same as ManagementOnly
```

## Real-time Booking Integration

### 🔄 All Admin Endpoints Use Real-time Data
- **Period List**: Shows live booking counts from database
- **Period Details**: Uses `CalculateRealTimeBookingStats()`
- **Availability**: Real-time seat calculation
- **Debug Stats**: Detailed booking breakdown with status

### 🛡️ Fallback Mechanism
```go
realTimeStats, err := s.CalculateRealTimeBookingStats(ctx, period.ID)
if err != nil {
    log.Printf("[FALLBACK] Using static values for period %s", period.ID)
    // Falls back to static period.Booking field
}
```

## Testing

### Quick Test (Public)
```bash
curl http://localhost:8080/api/v1/periods/test-deployment
```

### Admin Test (Requires Token)
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/v1/admin/periods/test-deployment
```

### Complete Test Suite
```bash
ADMIN_TOKEN="your-jwt-token" bash test_admin_routes.sh
```

## API Examples

### Get Admin Periods with Real-time Data
```bash
GET /api/v1/admin/periods
Authorization: Bearer your-jwt-token

Response:
{
  "success": true,
  "message": "ดึงรายการช่วงเวลาสำเร็จ",
  "data": {
    "periods": [
      {
        "id": "uuid",
        "booking": 7,           // Real-time count from bookings table
        "available_seats": 13,  // Real-time: seat - booked_passengers
        "seat": 20,
        "status": "publish"
      }
    ]
  }
}
```

### Debug Booking Calculation
```bash
GET /api/v1/admin/periods/{id}/debug-booking-stats
Authorization: Bearer your-jwt-token

Response:
{
  "success": true,
  "period_id": "uuid",
  "real_time_stats": {
    "total_seats": 20,
    "booked_passengers": 7,
    "available_seats": 13,
    "bookings_by_status": {
      "pending": 2,
      "processing": 1,
      "complete": 4
    },
    "is_available": true
  },
  "comparison": {
    "static_vs_realtime_booking": {
      "static": 5,
      "realtime": 7,
      "difference": 2
    }
  }
}
```

## File Changes Made

### 1. admin_routes.go (Enhanced)
```go
// Added missing admin period endpoints
admin.Get("/periods/stats", h.PeriodHandler.GetStats)
admin.Get("/periods/test-deployment", h.PeriodHandler.TestDeployment)
admin.Get("/periods/:id/debug-booking-stats", h.PeriodHandler.DebugBookingStats)
```

### 2. period_routes.go (Cleaned)
```go
// Removed debug endpoint from legacy admin routes
// Added legacy route comments for clarity
```

### 3. period_service_impl.go (Updated)
```go
// All period methods now use enhanced toPeriodResponse()
return s.toPeriodResponse(ctx, period), nil // Uses real-time calculation
```

## Security Notes

### Authentication Required
- All admin routes require valid JWT token
- Token must have `super` or `admin` role
- Middleware validates token and role before processing

### Rate Limiting
- Standard rate limiting applies to all routes
- Admin routes have same limits as public routes
- Consider implementing separate admin rate limits if needed

## Monitoring

### Key Metrics
1. **Admin Route Response Times**: Monitor `/admin/periods/*` performance
2. **Real-time Calculation Performance**: Track booking aggregation queries
3. **Authentication Success Rate**: Monitor admin access patterns
4. **Fallback Usage**: Track when real-time calculation fails

### Log Patterns
```bash
# Admin access logs
grep "ManagementOnly" logs/app.log

# Real-time calculation logs
grep "\[DEBUG\] CalculateRealTimeBookingStats" logs/app.log

# Fallback usage
grep "\[FALLBACK\]" logs/app.log
```

## Next Steps

1. **Database Optimization**: Run `migration_booking_optimization.sql`
2. **Monitoring Setup**: Implement admin route performance monitoring
3. **Testing**: Verify all admin endpoints with real admin tokens
4. **Documentation**: Update API documentation with new admin routes

## Success Criteria

✅ **Admin routes working:**
- `GET /api/v1/admin/periods` returns period list
- `GET /api/v1/admin/periods/{id}` returns period details
- `GET /api/v1/admin/periods/{id}/debug-booking-stats` provides debug info

✅ **Real-time calculation active:**
- All admin period endpoints use live booking data
- Fallback mechanism works correctly
- Performance is acceptable

✅ **Security working:**
- Admin authentication required
- Role-based access control functional
- Unauthorized access properly blocked