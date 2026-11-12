# GET /api/v1/sales/active - Active Sales List for Agents

## Endpoint Overview
**URL:** `GET {{base_url}}/sales/active`
**Description:** Get list of active sales agents (for agent+ access level)
**Authentication:** Required - JWT Token (Any authenticated user)
**Access Level:** Agent+ (any authenticated user can access)

## Purpose
This endpoint was created specifically to allow agents to get a list of active sales personnel when creating bookings, since the main sales list endpoint requires admin privileges.

## Request Structure

### Headers
```http
Authorization: Bearer <jwt_token>
```

### Query Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `offset` | integer | No | 0 | Number of records to skip for pagination |
| `limit` | integer | No | 20 | Maximum number of records to return (1-100) |
| `search` | string | No | - | Search term for sales agent name/code |

### Example Request
```http
GET {{base_url}}/sales/active?offset=0&limit=10&search=john
```

## Response Structure

### Success Response (200 OK)
```json
{
  "success": true,
  "message": "ดึงรายการยอดขายที่ใช้งานสำเร็จ",
  "data": {
    "sales": [
      {
        "id": "uuid",
        "user_id": "uuid",
        "code": "SALES001",
        "first_name": "John",
        "last_name": "Doe",
        "nick_name": "Johnny",
        "active": true,
        "created_at": "2025-09-26T01:30:00+07:00",
        "updated_at": "2025-09-26T01:30:00+07:00"
      },
      {
        "id": "uuid",
        "user_id": "uuid",
        "code": "SALES002",
        "first_name": "Jane",
        "last_name": "Smith",
        "nick_name": "Jane",
        "active": true,
        "created_at": "2025-09-25T10:15:00+07:00",
        "updated_at": "2025-09-25T10:15:00+07:00"
      }
    ],
    "pagination": {
      "total": 25,
      "has_more": true,
      "current_page": 1,
      "total_pages": 3,
      "page_size": 10,
      "limit": 10,
      "offset": 0
    }
  }
}
```

### Response Fields

#### Sales Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Sales agent unique identifier |
| `user_id` | string (UUID) | Associated user account ID |
| `code` | string | Sales agent code (e.g., SALES001) |
| `first_name` | string | First name |
| `last_name` | string | Last name |
| `nick_name` | string | Nickname/display name |
| `active` | boolean | Always true (only active sales returned) |
| `created_at` | string | Creation timestamp (ISO format) |
| `updated_at` | string | Last update timestamp (ISO format) |

#### Pagination Object
| Field | Type | Description |
|-------|------|-------------|
| `total` | number | Total number of active sales agents |
| `has_more` | boolean | Whether there are more pages |
| `current_page` | number | Current page number |
| `total_pages` | number | Total number of pages |
| `page_size` | number | Items per page |
| `limit` | number | Requested limit |
| `offset` | number | Requested offset |

## Error Responses

### 401 Unauthorized - Invalid/Missing Token
```json
{
  "success": false,
  "message": "Invalid token",
  "error": "Unauthorized"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "ดึงรายการยอดขายที่ใช้งานไม่สำเร็จ",
  "error": "Internal server error details"
}
```

## Key Features

### 🔓 **Public Access (Authenticated)**
- Any authenticated user can access (not restricted to admin only)
- Requires valid JWT token but no specific role restriction
- Perfect for agent booking flows

### 🔍 **Active Filter**
- Automatically filters to show only active sales agents (`active: true`)
- No need to specify active parameter - it's hardcoded

### 📊 **Pagination Support**
- Standard pagination with offset/limit
- Maximum limit of 100 records per request

### 🔍 **Search Capability**
- Search through sales agent names and codes
- Case-insensitive search functionality

## Usage Scenarios

### 1. **Booking Creation**
When agents create bookings, they need to select a sales agent for commission tracking:

```javascript
// Get active sales for dropdown
const response = await fetch('/api/v1/sales/active?limit=50');
const { data } = await response.json();
const salesOptions = data.sales.map(sale => ({
  value: sale.id,
  label: `${sale.code} - ${sale.first_name} ${sale.last_name}`
}));
```

### 2. **Sales Agent Search**
Search for specific sales agents during booking process:

```javascript
// Search sales agents
const searchSales = async (searchTerm) => {
  const response = await fetch(`/api/v1/sales/active?search=${searchTerm}&limit=10`);
  return response.json();
};
```

### 3. **Commission Assignment**
List available sales agents for commission assignment:

```javascript
// Get all active sales for commission reports
const activeSales = await fetch('/api/v1/sales/active?limit=100');
```

## Integration Notes

### **Booking Creation Integration**
This endpoint is designed to work seamlessly with the booking creation process:

```json
// Booking creation request uses sales_id from this endpoint
{
  "tour_id": "uuid",
  "period_id": "uuid",
  "sales_id": "uuid", // ← From /sales/active response
  "customer_name": "Customer Name",
  "customer_phone": "0123456789",
  "rate_types": [...]
}
```

### **Frontend Integration**
```typescript
interface ActiveSalesResponse {
  success: boolean;
  message: string;
  data: {
    sales: SalesAgent[];
    pagination: PaginationMeta;
  };
}

interface SalesAgent {
  id: string;
  user_id: string;
  code: string;
  first_name: string;
  last_name: string;
  nick_name: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

## Comparison with Other Sales Endpoints

| Endpoint | Access Level | Purpose | Active Filter |
|----------|-------------|---------|---------------|
| `GET /sales/active` | **Any authenticated user** | **Agent booking workflow** | **Auto-filtered (active only)** |
| `GET /sales/admin/` | Admin/Super only | Full sales management | Manual parameter |
| `GET /sales/my-*` | Sales role only | Own data viewing | N/A |

## Implementation Details

### **Route Definition**
```go
// In sales_routes.go
sales.Get("/active", h.SalesHandler.ListActiveSales) // GET /api/v1/sales/active
```

### **Handler Implementation**
```go
// In SalesHandler
func (h *SalesHandler) ListActiveSales(c *fiber.Ctx) error {
    // Auto-set active filter to true
    req := dto.SalesListRequest{
        Active: &[]bool{true}[0], // Force active: true
        Offset: offset,
        Limit:  limit,
    }

    // Use existing ListSales service method
    sales, err := h.salesService.ListSales(c.Context(), req)
    // ...
}
```

### **Middleware Stack**
1. `middleware.Protected()` - Requires valid JWT token
2. No additional role restrictions - any authenticated user can access

This endpoint provides the perfect solution for agents who need to select active sales personnel during the booking creation process without requiring admin privileges.