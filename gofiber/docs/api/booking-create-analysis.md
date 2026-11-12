# POST /api/v1/bookings - Agent Booking Creation Analysis

## Endpoint Overview
**URL:** `POST {{base_url}}/bookings`
**Description:** Agent creates booking for customers
**Authentication:** Required - JWT Token (Agent+ access level)
**Middleware:** `middleware.Protected()`, `middleware.BookingAccess()`

## Request Structure

### Headers
```http
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Request Body
```json
{
  "tour_id": "uuid",
  "period_id": "uuid",
  "sales_id": "uuid",
  "customer_name": "Customer Name",
  "customer_phone": "0123456789",
  "note": "Optional booking note",
  "room_types": [
    {
      "id": "room_type_id",
      "quantity": 2
    }
  ],
  "rate_types": [
    {
      "id": "rate_item_id",
      "quantity": 2
    }
  ],
  "additional_services": [
    {
      "id": "service_id",
      "quantity": 1
    }
  ]
}
```

### Request Fields Validation

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `tour_id` | UUID | ✅ | required | Tour ID to book |
| `period_id` | UUID | ✅ | required | Period ID within the tour |
| `sales_id` | UUID | ✅ | required | Sales agent ID for commission |
| `customer_name` | string | ✅ | required, max 200 chars | Customer's full name |
| `customer_phone` | string | ✅ | required | Customer's phone number |
| `note` | string | ❌ | optional | Booking notes |
| `room_types` | array | ❌ | optional | Room type selections |
| `rate_types` | array | ✅ | required, min 1 item | Rate/passenger type selections |
| `additional_services` | array | ❌ | optional | Additional service selections |

### Nested Object Validations

#### Room Type Request
```json
{
  "id": "string", // required - room type ID
  "quantity": 1   // required - min 1
}
```

#### Rate Type Request
```json
{
  "id": "string", // required - rate item ID
  "quantity": 2   // required - min 1
}
```

#### Additional Service Request
```json
{
  "id": "string", // required - service ID
  "quantity": 1   // required - min 1
}
```

---

## Response Structure

### Success Response (201 Created)
```json
{
  "success": true,
  "message": "สร้างการจองเรียบร้อยแล้ว",
  "data": {
    "id": "uuid",
    "invoice_id": "INV-2025-001234",
    "user_id": "uuid", // Agent user ID
    "tour_id": "uuid",
    "period_id": "uuid",
    "sales_id": "uuid",
    "customer_name": "Customer Name",
    "customer_phone": "0123456789",
    "note": "Booking note",
    "sum_room": 7800.0,
    "sum_rate": 24000.0,
    "sum_additional_services": 500.0,
    "sum_qty": 2,
    "sum_commission_qty": 2,
    "sum_net": 32300.0,
    "sum_commission": 1000.0,
    "discount": 0.0,
    "sum_deposit": 10000.0,
    "sum_payment": 0.0,
    "date_deposit": "2025-09-28",
    "date_final": "2025-10-25",
    "status": "pending",
    "expires_at": "2025-09-28T23:59:59+07:00",
    "created_at": "2025-09-26T01:30:00+07:00",
    "updated_at": "2025-09-26T01:30:00+07:00",

    // Related data (populated if included)
    "tour_details": {
      "id": "uuid",
      "title": "Tour Title",
      "code": "TOUR001",
      "thumbnail": "image_url"
    },
    "period_details": {
      "id": "uuid",
      "start_date": "2025-10-01",
      "end_date": "2025-10-05",
      "seat": 25,
      "booking": 5,
      "available_seats": 20
    },
    "sales_details": {
      "id": "uuid",
      "code": "SALES001",
      "first_name": "Sales",
      "last_name": "Agent"
    },
    "room_types": [
      {
        "id": "uuid",
        "title": "Single Room",
        "price": 3900.0,
        "quantity": 2,
        "total": 7800.0
      }
    ],
    "rate_types": [
      {
        "id": "uuid",
        "title": "ผู้ใหญ่",
        "price": 12000.0,
        "quantity": 2,
        "exclude": false,
        "total": 24000.0
      }
    ],
    "additional_services": [
      {
        "id": "uuid",
        "title": "Travel Insurance",
        "price": 150.0,
        "quantity": 2,
        "rate_type": "per_unit",
        "unit_name": "คน",
        "total": 300.0
      }
    ],
    "payment_slips": [],
    "passports": []
  }
}
```

### Response Fields Description

#### Main Booking Fields
| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Booking unique identifier |
| `invoice_id` | string | Auto-generated invoice number |
| `user_id` | UUID | Agent user ID (from JWT token) |
| `tour_id` | UUID | Tour being booked |
| `period_id` | UUID | Period being booked |
| `sales_id` | UUID | Sales agent for commission |
| `customer_name` | string | Customer's name |
| `customer_phone` | string | Customer's phone |
| `note` | string\|null | Booking notes |

#### Financial Summary Fields
| Field | Type | Description |
|-------|------|-------------|
| `sum_room` | number | Total room charges |
| `sum_rate` | number | Total rate/passenger charges |
| `sum_additional_services` | number | Total additional service charges |
| `sum_qty` | number | Total passenger count |
| `sum_commission_qty` | number | Commission-eligible passenger count |
| `sum_net` | number | Total booking amount |
| `sum_commission` | number | Total commission amount |
| `discount` | number | Discount applied |
| `sum_deposit` | number | Required deposit amount |
| `sum_payment` | number | Payment received so far |

#### Date/Status Fields
| Field | Type | Description |
|-------|------|-------------|
| `date_deposit` | string | Deposit due date (YYYY-MM-DD) |
| `date_final` | string | Final payment due date (YYYY-MM-DD) |
| `status` | string | Booking status (pending, processing, confirm, cancel, completed) |
| `expires_at` | string\|null | Booking expiration timestamp (ISO format) |
| `created_at` | string | Creation timestamp (ISO format) |
| `updated_at` | string | Last update timestamp (ISO format) |

#### Related Objects
| Field | Type | Description |
|-------|------|-------------|
| `tour_details` | object\|null | Tour information |
| `period_details` | object\|null | Period information |
| `sales_details` | object\|null | Sales agent information |
| `room_types` | array | Booked room types with quantities/totals |
| `rate_types` | array | Booked rate types with quantities/totals |
| `additional_services` | array | Booked additional services |
| `payment_slips` | array | Payment slip records |
| `passports` | array | Passenger passport information |

---

## Error Responses

### 400 Bad Request - Validation Error
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "tour_id": ["tour_id is required"],
    "customer_name": ["customer_name is required and must be less than 200 characters"],
    "rate_types": ["rate_types is required and must have at least 1 item"]
  }
}
```

### 401 Unauthorized - Invalid/Missing Token
```json
{
  "success": false,
  "message": "Invalid token",
  "error": "Unauthorized"
}
```

### 403 Forbidden - Insufficient Permission
```json
{
  "success": false,
  "message": "Access denied",
  "error": "Insufficient permissions"
}
```

### 400 Bad Request - Business Logic Error
```json
{
  "success": false,
  "message": "สร้างการจองไม่สำเร็จ",
  "error": "Period is fully booked"
}
```

---

## Business Logic & Features

### 🔐 Authentication & Authorization
- **Required:** Valid JWT token
- **Access Level:** Agent+ (agent, admin, super)
- **Middleware:** `BookingAccess()` - ensures agent+ role

### 💰 Automatic Calculations
- **Financial Totals:** Auto-calculates all sum_* fields based on selections
- **Commission:** Calculated based on sales agent and commission structure
- **Deposit:** Calculated based on period's deposit settings
- **Due Dates:** Auto-calculated based on tour dates and payment rules

### ⏰ Booking Expiration
- **Auto-expiry:** Bookings have expiration time for deposit payment
- **Status Management:** Automatically handles pending → expired transitions

### 🔄 Real-time Validation
- **Seat Availability:** Validates against real-time seat availability
- **Period Status:** Ensures period is active and bookable
- **Sales Agent:** Validates sales agent exists and is active

### 📊 Data Relationships
- **Tour/Period:** Links to tour and specific period
- **Sales Agent:** Links for commission tracking
- **Room/Rate/Service:** Detailed breakdown of selected items

---

## Usage Examples

### Basic Booking (Adult passengers only)
```json
{
  "tour_id": "1adde5e4-4031-4b46-8db5-cd7ece185437",
  "period_id": "2c6637f0-83b6-4e3e-8b20-84d201bd4b93",
  "sales_id": "sales-agent-uuid",
  "customer_name": "John Doe",
  "customer_phone": "0812345678",
  "rate_types": [
    {
      "id": "dc69c4cb-e895-498f-95de-b5ed850ce01b", // Adult rate
      "quantity": 2
    }
  ]
}
```

### Full Booking (With rooms and services)
```json
{
  "tour_id": "1adde5e4-4031-4b46-8db5-cd7ece185437",
  "period_id": "7aeb7b57-9d6f-4b8a-96b1-2b6424c8d06b",
  "sales_id": "sales-agent-uuid",
  "customer_name": "Jane Smith",
  "customer_phone": "0987654321",
  "note": "Customer requests vegetarian meals",
  "room_types": [
    {
      "id": "49a5d314-a8aa-4cbd-9e16-90274630b283", // Single room
      "quantity": 1
    }
  ],
  "rate_types": [
    {
      "id": "4b30fa60-154a-4eb7-80b6-4edd34fb99f0", // Adult rate
      "quantity": 1
    }
  ],
  "additional_services": [
    {
      "id": "ecdf606b-417d-46f7-992a-0bfaa8bb6c84", // Car rental
      "quantity": 3
    }
  ]
}
```

---

## TypeScript Interface

```typescript
// Request Types
export interface CreateBookingRequest {
  tour_id: string;
  period_id: string;
  sales_id: string;
  customer_name: string;
  customer_phone: string;
  note?: string;
  room_types?: CreateBookingRoomTypeRequest[];
  rate_types: CreateBookingRateTypeRequest[]; // required, min 1
  additional_services?: CreateBookingAdditionalServiceRequest[];
}

export interface CreateBookingRoomTypeRequest {
  id: string;
  quantity: number; // min 1
}

export interface CreateBookingRateTypeRequest {
  id: string;
  quantity: number; // min 1
}

export interface CreateBookingAdditionalServiceRequest {
  id: string;
  quantity: number; // min 1
}

// Response Types
export interface BookingCreateResponse {
  success: boolean;
  message: string;
  data: BookingResponse;
}

export interface BookingResponse {
  id: string;
  invoice_id: string;
  user_id: string;
  tour_id: string;
  period_id: string;
  sales_id: string;
  customer_name: string;
  customer_phone: string;
  note: string | null;
  sum_room: number;
  sum_rate: number;
  sum_additional_services: number;
  sum_qty: number;
  sum_commission_qty: number;
  sum_net: number;
  sum_commission: number;
  discount: number;
  sum_deposit: number;
  sum_payment: number;
  date_deposit: string; // YYYY-MM-DD
  date_final: string; // YYYY-MM-DD
  status: "pending" | "processing" | "confirm" | "cancel" | "completed";
  expires_at: string | null; // ISO timestamp
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp

  // Relations (optional)
  tour_details?: TourResponse;
  period_details?: PeriodSummaryResponse;
  sales_details?: SalesResponse;
  room_types?: BookingRoomTypeResponse[];
  rate_types?: BookingRateTypeResponse[];
  additional_services?: BookingAdditionalServiceResponse[];
  payment_slips?: PaymentSlipResponse[];
  passports?: PassengerPassportResponse[];
}
```