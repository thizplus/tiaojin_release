# Tours Regular API

## Get Tours with Periods

**Endpoint:** `GET {{base_url}}/tours/regular`

**Description:** Retrieves tours with their complete period information, including real-time seat availability, pricing details, room types, and additional services.

### Request Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `offset` | integer | No | 0 | Number of records to skip for pagination |
| `limit` | integer | No | 10 | Maximum number of records to return (1-100) |
| `search` | string | No | - | Search term for tour title or code |
| `status` | string | No | "publish" | Tour status filter (`publish`, `draft`) |
| `city_id` | string (UUID) | No | - | Filter tours by city ID |
| `available_only` | boolean | No | false | Show only tours with available seats |

### Example Request

```http
GET {{base_url}}/tours/regular?limit=2&available_only=true&status=publish
```

### Response

#### Success Response (200 OK)

```json
{
  "success": true,
  "message": "ดึงรายการทัวร์พร้อม periods สำเร็จ",
  "data": {
    "tours": [
      {
        "id": "1adde5e4-4031-4b46-8db5-cd7ece185437",
        "title": "WHISPERS OF WINTER CHONGQING CHENGDU ฉงชิ่ง เฉิงตู สี่ดรุณีซวงเฉียวโกว 6 วัน 5 คืน",
        "code": "BCCKG12-HU",
        "city_id": "8e777264-01b1-4e4b-91a1-c3d337eaa71a",
        "city_title": "Bangkok",
        "detail": "*เที่ยวครบทุกวัน *นั่งรถไฟความเร็วสูง 1 เที่ยว *อุทยานซวงเฉียวโกว (รวมรถเวียน) *สะพานหนานเฉียวธารน้ำสีฟ้า *เมนูพิเศษ!สุกี้เห็ด",
        "thumbnail": "https://storage-humtung.b-cdn.net/uploads/tours/image/8e7e4dd6-5d88-42bf-bbf9-e32d58222739.jpg",
        "file_pdf": "https://storage-humtung.b-cdn.net/uploads/tours/pdf/13a6fd7c-aff7-4e71-abeb-3064f63e0307.pdf",
        "file_word": "https://storage-humtung.b-cdn.net/uploads/tours/document/fdadc247-7552-476b-902f-8c607d81a1ed.docx",
        "status": "publish",
        "created_at": "2025-09-24T14:07:14+07:00",
        "updated_at": "2025-09-25T08:36:22+07:00",
        "available_periods": [
          {
            "id": "2c6637f0-83b6-4e3e-8b20-84d201bd4b93",
            "tour_id": "1adde5e4-4031-4b46-8db5-cd7ece185437",
            "start_date": "2025-09-27",
            "end_date": "2025-09-28",
            "seat": 1,
            "booking": 0,
            "held_seats": 0,
            "available_seats": 1,
            "deposit": 5000,
            "deposit_type": "per_unit",
            "airline": "BKK",
            "airport": "SUWAN",
            "note": "TESTER",
            "sector_1": [
              {
                "date": "2025-09-26",
                "flight": "A",
                "route": "DS"
              }
            ],
            "sector_2": [
              {
                "date": "2025-09-27",
                "flight": "A",
                "route": "SD"
              }
            ],
            "rate_type": "regular",
            "rate_items": [
              {
                "id": "dc69c4cb-e895-498f-95de-b5ed850ce01b",
                "title": "ผู้ใหญ่",
                "regular": 15000,
                "sale": 12000,
                "exclude": false,
                "sort_order": 0
              }
            ],
            "commission_type": "regular",
            "commission_items": [
              {
                "regular": 500,
                "sale": 500
              },
              {
                "regular": 200,
                "sale": 300
              }
            ],
            "room_types": [
              {
                "id": "b37ca1ae-27dc-42d8-acdf-3ac181a8361c",
                "title": "single",
                "price": 3900,
                "sort_order": 0
              },
              {
                "id": "dd442caa-a308-46a0-858d-33c7c14a97fa",
                "title": "double",
                "price": 0,
                "sort_order": 1
              },
              {
                "id": "e66bf7cb-bfdd-4a91-ae48-c5a38ad40b4f",
                "title": "twin",
                "price": 0,
                "sort_order": 2
              },
              {
                "id": "cc0b7d8b-31ec-4ddf-8002-cef451b9847a",
                "title": "triple",
                "price": 0,
                "sort_order": 3
              }
            ],
            "additional_services": [],
            "files": [],
            "status": "publish",
            "created_at": "2025-09-25T13:42:42+07:00",
            "updated_at": "2025-09-26T00:35:21+07:00"
          }
        ],
        "min_start_date": "2025-09-27"
      }
    ],
    "pagination": {
      "total": 1,
      "has_more": false,
      "current_page": 1,
      "total_pages": 1,
      "page_size": 2,
      "limit": 0,
      "offset": 0
    }
  }
}
```

### Response Fields

#### Tour Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Tour unique identifier |
| `title` | string | Tour title |
| `code` | string | Tour code |
| `city_id` | string (UUID) | City ID where tour starts |
| `city_title` | string | City name |
| `detail` | string | Tour description |
| `thumbnail` | string | Tour image URL |
| `file_pdf` | string | PDF file URL |
| `file_word` | string | Word document URL |
| `status` | string | Tour status (`publish`, `draft`) |
| `created_at` | string | Creation timestamp |
| `updated_at` | string | Last update timestamp |
| `available_periods` | array | Array of available periods |
| `min_start_date` | string | Earliest period start date |

#### Period Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Period unique identifier |
| `tour_id` | string (UUID) | Tour ID this period belongs to |
| `start_date` | string | Period start date (YYYY-MM-DD) |
| `end_date` | string | Period end date (YYYY-MM-DD) |
| `seat` | integer | Total available seats |
| `booking` | integer | Currently booked passengers (real-time) |
| `held_seats` | integer | Temporarily held seats (real-time) |
| `available_seats` | integer | Available seats (calculated: seat - booking - held_seats) |
| `deposit` | number | Deposit amount |
| `deposit_type` | string | Deposit calculation type (`per_unit`, `percent`) |
| `airline` | string | Airline name |
| `airport` | string | Airport code |
| `note` | string | Additional notes |
| `sector_1` | array | Outbound flight sectors |
| `sector_2` | array | Return flight sectors |
| `rate_type` | string | Pricing type (`regular`, `sale`) |
| `rate_items` | array | Pricing structure for different passenger types |
| `commission_type` | string | Commission structure type |
| `commission_items` | array | Commission rates |
| `room_types` | array | Available room types and pricing |
| `additional_services` | array | Available additional services |
| `files` | array | Associated files |
| `status` | string | Period status |
| `created_at` | string | Creation timestamp |
| `updated_at` | string | Last update timestamp |

#### Rate Item Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Rate item unique identifier |
| `title` | string | Rate item title (e.g., "ผู้ใหญ่", "เด็ก") |
| `regular` | number | Regular price |
| `sale` | number | Sale price |
| `exclude` | boolean | Whether this rate is excluded from calculations |
| `sort_order` | integer | Display order |

#### Room Type Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Room type unique identifier |
| `title` | string | Room type name |
| `price` | number | Additional price for this room type |
| `sort_order` | integer | Display order |

#### Additional Service Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Service unique identifier |
| `title` | string | Service name |
| `price` | number | Service price |
| `rate_type` | string | Pricing calculation (`per_unit`, `per_booking`) |
| `unit_name` | string | Unit description |
| `sort_order` | integer | Display order |

### Error Responses

#### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request parameters",
  "error": "Validation failed"
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "message": "ดึงรายการทัวร์พร้อม periods ไม่สำเร็จ",
  "error": "Internal server error"
}
```

### Usage Notes

1. **Real-time Data**: The `booking`, `held_seats`, and `available_seats` fields are calculated in real-time using active bookings and seat holds.

2. **Filtering**: Use `available_only=true` to show only periods with available seats (where `available_seats > 0`).

3. **Pagination**: Use `offset` and `limit` parameters for pagination. Maximum `limit` is 100.

4. **Search**: The `search` parameter searches through tour titles and codes.

5. **Complete Data**: This endpoint provides the same level of detail as the admin periods endpoint, including pricing, room types, and additional services.