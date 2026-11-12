# Tour Routes API Summary

## Public Routes (No Authentication Required)

### 1. GET /api/v1/tours
**Handler:** `TourHandler.List`
**Description:** Get paginated list of tours (basic information only)

**Request Parameters:**
```
offset (int, optional): Default 0, pagination offset
limit (int, optional): Default 10, max 100, number of items per page
search (string, optional): Search in tour title/code
status (string, optional): Filter by status (publish/draft)
city_id (UUID, optional): Filter by city ID
```

**Example Request:**
```http
GET {{base_url}}/tours?offset=0&limit=20&status=publish&search=bangkok
```

**Response:**
```json
{
  "success": true,
  "message": "ดึงรายการทัวร์สำเร็จ",
  "data": {
    "tours": [
      {
        "id": "uuid",
        "title": "Tour Title",
        "code": "TOUR001",
        "city_id": "uuid",
        "city_title": "City Name",
        "detail": "Tour description",
        "thumbnail": "image_url",
        "file_pdf": "pdf_url",
        "file_word": "doc_url",
        "status": "publish",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": 100,
      "has_more": true,
      "current_page": 1,
      "total_pages": 5,
      "page_size": 20,
      "limit": 20,
      "offset": 0
    }
  }
}
```

---

### 2. ⭐ GET /api/v1/tours/regular
**Handler:** `TourHandler.ListWithPeriods`
**Description:** Get tours with complete periods information including real-time seat availability, pricing, room types, and additional services

**Request Parameters:**
```
offset (int, optional): Default 0, pagination offset
limit (int, optional): Default 10, max 100, number of items per page
search (string, optional): Search in tour title/code
status (string, optional): Default "publish", filter by status
city_id (UUID, optional): Filter by city ID
available_only (boolean, optional): Default false, show only tours with available seats
```

**Example Request:**
```http
GET {{base_url}}/tours/regular?offset=0&limit=10&available_only=true&status=publish
```

**Response:**
```json
{
  "success": true,
  "message": "ดึงรายการทัวร์พร้อม periods สำเร็จ",
  "data": {
    "tours": [
      {
        "id": "uuid",
        "title": "Tour Title",
        "code": "TOUR001",
        "city_id": "uuid",
        "city_title": "City Name",
        "detail": "Tour description",
        "thumbnail": "image_url",
        "file_pdf": "pdf_url",
        "file_word": "doc_url",
        "status": "publish",
        "created_at": "timestamp",
        "updated_at": "timestamp",
        "available_periods": [
          {
            "id": "uuid",
            "tour_id": "uuid",
            "start_date": "2025-09-27",
            "end_date": "2025-09-28",
            "seat": 25,
            "booking": 5,
            "held_seats": 3,
            "available_seats": 17,
            "deposit": 5000,
            "deposit_type": "per_unit",
            "airline": "Thai Airways",
            "airport": "BKK",
            "note": "Special note",
            "sector_1": [{"date": "2025-09-27", "flight": "TG123", "route": "BKK-HKG"}],
            "sector_2": [{"date": "2025-09-28", "flight": "TG124", "route": "HKG-BKK"}],
            "rate_type": "regular",
            "rate_items": [
              {
                "id": "uuid",
                "title": "ผู้ใหญ่",
                "regular": 15000,
                "sale": 12000,
                "exclude": false,
                "sort_order": 0
              }
            ],
            "commission_type": "regular",
            "commission_items": [{"regular": 500, "sale": 500}],
            "room_types": [
              {
                "id": "uuid",
                "title": "single",
                "price": 3900,
                "sort_order": 0
              }
            ],
            "additional_services": [
              {
                "id": "uuid",
                "title": "Travel Insurance",
                "price": 150,
                "rate_type": "per_unit",
                "unit_name": "คน",
                "sort_order": 0
              }
            ],
            "files": [],
            "status": "publish",
            "created_at": "timestamp",
            "updated_at": "timestamp"
          }
        ],
        "min_start_date": "2025-09-27"
      }
    ],
    "pagination": {
      "total": 50,
      "has_more": true,
      "current_page": 1,
      "total_pages": 5,
      "page_size": 10,
      "limit": 10,
      "offset": 0
    }
  }
}
```

---

### 3. GET /api/v1/tours/:id
**Handler:** `TourHandler.GetByID`
**Description:** Get single tour by ID (basic information only)

**Request Parameters:**
```
id (UUID, required): Tour ID in URL path
```

**Example Request:**
```http
GET {{base_url}}/tours/1adde5e4-4031-4b46-8db5-cd7ece185437
```

**Response:**
```json
{
  "success": true,
  "message": "ดึงข้อมูลทัวร์สำเร็จ",
  "data": {
    "id": "uuid",
    "title": "Tour Title",
    "code": "TOUR001",
    "city_id": "uuid",
    "city_title": "City Name",
    "detail": "Tour description",
    "thumbnail": "image_url",
    "file_pdf": "pdf_url",
    "file_word": "doc_url",
    "status": "publish",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

---

### 4. GET /api/v1/tours/:id/periods
**Handler:** `TourHandler.GetPeriods`
**Description:** Get periods for specific tour with pagination

**Request Parameters:**
```
id (UUID, required): Tour ID in URL path
offset (int, optional): Default 0, pagination offset
limit (int, optional): Default 20, max 100, number of periods per page
```

**Example Request:**
```http
GET {{base_url}}/tours/1adde5e4-4031-4b46-8db5-cd7ece185437/periods?offset=0&limit=10
```

**Response:**
```json
{
  "success": true,
  "message": "ดึงช่วงเวลาทัวร์สำเร็จ",
  "data": {
    "periods": [
      {
        "id": "uuid",
        "tour_id": "uuid",
        "start_date": "2025-09-27",
        "end_date": "2025-09-28",
        "seat": 25,
        "booking": 5,
        "held_seats": 3,
        "available_seats": 17,
        "deposit": 5000,
        "deposit_type": "per_unit",
        "status": "publish",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": 10,
      "has_more": false,
      "current_page": 1,
      "total_pages": 1,
      "page_size": 10,
      "limit": 10,
      "offset": 0
    }
  }
}
```

---

## Admin Routes (Require Authentication + Admin Role)

### 5. POST /api/v1/tours
**Handler:** `TourHandler.Create`
**Description:** Create new tour (admin only)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Tour Title",
  "code": "TOUR001",
  "city_id": "uuid",
  "detail": "Tour description",
  "thumbnail": "image_url",
  "status": "draft"
}
```

**Response:**
```json
{
  "success": true,
  "message": "สร้างทัวร์เรียบร้อยแล้ว",
  "data": {
    "id": "uuid",
    "title": "Tour Title",
    "code": "TOUR001",
    "city_id": "uuid",
    "city_title": "City Name",
    "detail": "Tour description",
    "thumbnail": "image_url",
    "file_pdf": null,
    "file_word": null,
    "status": "draft",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

---

### 6. GET /api/v1/tours/stats
**Handler:** `TourHandler.GetStats`
**Description:** Get tour statistics (admin only)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```http
GET {{base_url}}/tours/stats
```

**Response:**
```json
{
  "success": true,
  "message": "ดึงสถิติทัวร์สำเร็จ",
  "data": {
    "total_tours": 150,
    "published_tours": 120,
    "draft_tours": 30,
    "tours_with_periods": 85
  }
}
```

---

### 7. PUT /api/v1/tours/:id
**Handler:** `TourHandler.Update`
**Description:** Update tour information (admin only)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Parameters:**
```
id (UUID, required): Tour ID in URL path
```

**Request Body:**
```json
{
  "title": "Updated Tour Title",
  "city_id": "uuid",
  "detail": "Updated description",
  "thumbnail": "new_image_url",
  "file_pdf": "new_pdf_url",
  "file_word": "new_doc_url",
  "status": "publish"
}
```

**Response:**
```json
{
  "success": true,
  "message": "อัปเดตทัวร์เรียบร้อยแล้ว",
  "data": {
    "id": "uuid",
    "title": "Updated Tour Title",
    "code": "TOUR001",
    "city_id": "uuid",
    "city_title": "City Name",
    "detail": "Updated description",
    "thumbnail": "new_image_url",
    "file_pdf": "new_pdf_url",
    "file_word": "new_doc_url",
    "status": "publish",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

---

### 8. DELETE /api/v1/tours/:id
**Handler:** `TourHandler.Delete`
**Description:** Delete tour (admin only)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Parameters:**
```
id (UUID, required): Tour ID in URL path
```

**Example Request:**
```http
DELETE {{base_url}}/tours/1adde5e4-4031-4b46-8db5-cd7ece185437
```

**Response:**
```json
{
  "success": true,
  "message": "ลบทัวร์เรียบร้อยแล้ว",
  "data": null
}
```

---

### 9. PUT /api/v1/tours/:id/status
**Handler:** `TourHandler.UpdateStatus`
**Description:** Update tour status only (admin only)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Parameters:**
```
id (UUID, required): Tour ID in URL path
```

**Request Body:**
```json
{
  "status": "publish"
}
```

**Response:**
```json
{
  "success": true,
  "message": "อัปเดตสถานะทัวร์เรียบร้อยแล้ว",
  "data": {
    "id": "uuid",
    "title": "Tour Title",
    "code": "TOUR001",
    "city_id": "uuid",
    "city_title": "City Name",
    "detail": "Tour description",
    "thumbnail": "image_url",
    "file_pdf": "pdf_url",
    "file_word": "doc_url",
    "status": "publish",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

---

## Key Features of /tours/regular Route

### ⭐ Special Features:
1. **Real-time Seat Calculation**: Uses same logic as admin system
   - `booking`: Live count of confirmed + pending + processing bookings
   - `held_seats`: Active seat holds in the system
   - `available_seats`: Calculated as `total_seats - booking - held_seats`

2. **Complete Pricing Information**:
   - `rate_items`: Full pricing structure for different passenger types
   - `room_types`: Available room configurations with pricing
   - `additional_services`: Optional services with pricing

3. **Flight Information**: Complete sector details with dates, flights, and routes

4. **Commission Structure**: Full commission breakdown for agents

5. **File Attachments**: Associated files and documents

### Usage Scenarios:
- **Public booking system**: Customers can see complete tour information with real availability
- **Agent portal**: Travel agents get full pricing and commission details
- **Tour comparison**: Complete information for comparing different tours
- **Inventory management**: Real-time availability for booking systems

### Performance Notes:
- Uses preloaded relationships for optimal database queries
- Real-time calculations ensure accurate seat availability
- Efficient pagination for large tour catalogs