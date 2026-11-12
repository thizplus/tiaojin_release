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
  has_waiting_list: boolean;
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
  service_fee_per_pax: number;
  created_at: string;
  updated_at: string;

  // Expiration Override - สำหรับทัวร์ไฟไหม้
  expiration_override: boolean;
  expiration_minutes?: number;
  expiration_hours?: number;
  expiration_days?: number;
  expiration_message?: string;
  expiration_display?: string; // Human-readable duration (e.g., "15 นาที")
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
  page?: number;
  limit?: number;
  search?: string;
  status?: "publish" | "draft";
  city_id?: string;
  available_only?: boolean;
}