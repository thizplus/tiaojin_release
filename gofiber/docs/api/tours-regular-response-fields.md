# GET /tours/regular - Response Fields Summary

## Root Response Structure
```typescript
interface ToursRegularResponse {
  success: boolean;
  message: string;
  data: ToursRegularData;
}
```

## Main Data Structure
```typescript
interface ToursRegularData {
  tours: Tour[];
  pagination: PaginationMeta;
}
```

## Tour Object
```typescript
interface Tour {
  id: string; // UUID
  title: string;
  code: string;
  city_id: string | null; // UUID
  city_title: string | null;
  detail: string | null;
  thumbnail: string | null; // URL
  file_pdf: string | null; // URL
  file_word: string | null; // URL
  status: string; // "publish" | "draft"
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp
  available_periods: Period[];
  min_start_date: string | null; // YYYY-MM-DD format
}
```

## Period Object
```typescript
interface Period {
  id: string; // UUID
  tour_id: string; // UUID
  start_date: string; // YYYY-MM-DD
  end_date: string; // YYYY-MM-DD
  seat: number; // Total seats
  booking: number; // Booked passengers (real-time)
  held_seats: number; // Held seats (real-time)
  available_seats: number; // Available seats (calculated)
  deposit: number;
  deposit_type: string; // "per_unit" | "percent"
  airline: string | null;
  airport: string | null;
  note: string | null;
  sector_1: FlightSector[];
  sector_2: FlightSector[];
  rate_type: string; // "regular" | "sale"
  rate_items: RateItem[];
  commission_type: string; // "regular" | "sale"
  commission_items: CommissionItem[];
  room_types: RoomType[];
  additional_services: AdditionalService[];
  files: FileAttachment[];
  status: string; // "publish" | "draft"
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp
}
```

## Flight Sector Object
```typescript
interface FlightSector {
  date: string; // YYYY-MM-DD
  flight: string;
  route: string;
}
```

## Rate Item Object
```typescript
interface RateItem {
  id: string; // UUID
  title: string; // e.g., "ผู้ใหญ่", "เด็ก"
  regular: number; // Regular price
  sale: number; // Sale price
  exclude: boolean; // Whether excluded from calculations
  sort_order: number;
}
```

## Commission Item Object
```typescript
interface CommissionItem {
  regular: number; // Regular commission
  sale: number; // Sale commission
}
```

## Room Type Object
```typescript
interface RoomType {
  id: string; // UUID
  title: string; // e.g., "single", "double", "twin"
  price: number; // Additional price for this room type
  sort_order: number;
}
```

## Additional Service Object
```typescript
interface AdditionalService {
  id: string; // UUID
  title: string; // Service name
  price: number;
  rate_type: string; // "per_unit" | "per_booking"
  unit_name: string; // e.g., "คน", "วัน", "เหมา"
  sort_order: number;
}
```

## File Attachment Object
```typescript
interface FileAttachment {
  id: string; // UUID
  name: string;
  url: string;
  type: string;
  size: number;
}
```

## Pagination Meta Object
```typescript
interface PaginationMeta {
  total: number; // Total number of tours
  has_more: boolean; // Whether there are more pages
  current_page: number;
  total_pages: number;
  page_size: number; // Items per page
  limit: number; // Requested limit
  offset: number; // Requested offset
}
```

---

## Complete TypeScript Interface File

```typescript
// tours-regular.types.ts

export interface ToursRegularResponse {
  success: boolean;
  message: string;
  data: ToursRegularData;
}

export interface ToursRegularData {
  tours: Tour[];
  pagination: PaginationMeta;
}

export interface Tour {
  id: string;
  title: string;
  code: string;
  city_id: string | null;
  city_title: string | null;
  detail: string | null;
  thumbnail: string | null;
  file_pdf: string | null;
  file_word: string | null;
  status: "publish" | "draft";
  created_at: string;
  updated_at: string;
  available_periods: Period[];
  min_start_date: string | null;
}

export interface Period {
  id: string;
  tour_id: string;
  start_date: string;
  end_date: string;
  seat: number;
  booking: number;
  held_seats: number;
  available_seats: number;
  deposit: number;
  deposit_type: "per_unit" | "percent";
  airline: string | null;
  airport: string | null;
  note: string | null;
  sector_1: FlightSector[];
  sector_2: FlightSector[];
  rate_type: "regular" | "sale";
  rate_items: RateItem[];
  commission_type: "regular" | "sale";
  commission_items: CommissionItem[];
  room_types: RoomType[];
  additional_services: AdditionalService[];
  files: FileAttachment[];
  status: "publish" | "draft";
  created_at: string;
  updated_at: string;
}

export interface FlightSector {
  date: string;
  flight: string;
  route: string;
}

export interface RateItem {
  id: string;
  title: string;
  regular: number;
  sale: number;
  exclude: boolean;
  sort_order: number;
}

export interface CommissionItem {
  regular: number;
  sale: number;
}

export interface RoomType {
  id: string;
  title: string;
  price: number;
  sort_order: number;
}

export interface AdditionalService {
  id: string;
  title: string;
  price: number;
  rate_type: "per_unit" | "per_booking";
  unit_name: string;
  sort_order: number;
}

export interface FileAttachment {
  id: string;
  name: string;
  url: string;
  type: string;
  size: number;
}

export interface PaginationMeta {
  total: number;
  has_more: boolean;
  current_page: number;
  total_pages: number;
  page_size: number;
  limit: number;
  offset: number;
}

// Request Parameters Interface
export interface ToursRegularRequest {
  offset?: number;
  limit?: number;
  search?: string;
  status?: "publish" | "draft";
  city_id?: string;
  available_only?: boolean;
}
```

---

## Key Data Types Summary

### 🔢 **Numeric Fields**
- `seat`, `booking`, `held_seats`, `available_seats` - integers
- `deposit`, `regular`, `sale`, `price` - numbers (can be decimal)
- `sort_order` - integer

### 📝 **String Fields**
- All IDs are UUIDs (string format)
- Dates are in `YYYY-MM-DD` format
- Timestamps are ISO format strings
- URLs are string URLs

### 🔄 **Real-time Calculated Fields**
- `booking` - Live count from database
- `held_seats` - Live count from seat hold system
- `available_seats` - Calculated: `seat - booking - held_seats`

### 📊 **Array Fields**
- `available_periods[]` - Array of Period objects
- `sector_1[]`, `sector_2[]` - Arrays of FlightSector objects
- `rate_items[]`, `room_types[]`, `additional_services[]` - Pricing arrays
- `commission_items[]` - Commission structure array
- `files[]` - File attachment array

### 🔤 **Enum-like Fields**
- `status`: `"publish"` | `"draft"`
- `deposit_type`: `"per_unit"` | `"percent"`
- `rate_type`: `"regular"` | `"sale"`
- `commission_type`: `"regular"` | `"sale"`
- `rate_type` (in AdditionalService): `"per_unit"` | `"per_booking"`