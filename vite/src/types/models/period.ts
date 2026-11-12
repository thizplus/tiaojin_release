export type PeriodStatus = 'draft' | 'publish';
export type DepositType = 'per_unit' | 'percent';
export type RateType = 'regular' | 'sale';
export type CommissionType = 'regular' | 'sale';
export type ServiceRateType = 'per_unit' | 'per_booking';

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
  rate_type: ServiceRateType;
  unit_name: string;
  sort_order: number;
}

export interface PeriodFile {
  name: string;
  type: string;
  url: string;
}

export interface TourPeriod {
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
  deposit_type: DepositType;
//  single_stay: number;
  airline: string;
  airport: string;
  note?: string;
  sector_1: FlightSector[];
  sector_2: FlightSector[];
  rate_type: RateType;
  rate_items: RateItem[];
  commission_type: CommissionType;
  commission_items: CommissionItem[];
  room_types: RoomType[];
  additional_services: AdditionalService[];
  files: PeriodFile[];
  status: PeriodStatus;
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

export interface TourPeriodsFilter {
  page?: number;
  offset?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
  status?: PeriodStatus;
  search?: string;
  start_date_from?: string;
  start_date_to?: string;
  end_date_from?: string;
  end_date_to?: string;
}

export interface TourPeriodsResponse {
  periods: TourPeriod[];
  pagination: {
    total: number;
    has_more: boolean;
    current_page: number;
    total_pages: number;
    page_size: number;
    limit: number;
    offset: number;
  };
}

export interface CreateTourPeriodRequest {
  tour_id: string;
  start_date: string;
  end_date: string;
  seat: number;
  deposit: number;
  deposit_type: DepositType;
  //single_stay: number;
  airline: string;
  airport: string;
  note?: string;
  sector_1: FlightSector[];
  sector_2: FlightSector[];
  rate_type: RateType;
  rate_items: RateItem[];
  commission_type: CommissionType;
  commission_items: CommissionItem[];
  room_types: RoomType[];
  additional_services: AdditionalService[];
  files?: PeriodFile[];
  status?: PeriodStatus;
  service_fee_per_pax?: number;

  // Expiration Override
  expiration_override?: boolean;
  expiration_minutes?: number;
  expiration_hours?: number;
  expiration_days?: number;
  expiration_message?: string;
}

export interface UpdateTourPeriodRequest {
  start_date?: string;
  end_date?: string;
  seat?: number;
  deposit?: number;
  deposit_type?: DepositType;
  single_stay?: number;
  airline?: string;
  airport?: string;
  note?: string;
  sector_1?: FlightSector[];
  sector_2?: FlightSector[];
  rate_type?: RateType;
  rate_items?: RateItem[];
  commission_type?: CommissionType;
  commission_items?: CommissionItem[];
  room_types?: RoomType[];
  additional_services?: AdditionalService[];
  files?: PeriodFile[];
  status?: PeriodStatus;
  service_fee_per_pax?: number;

  // Expiration Override
  expiration_override?: boolean;
  expiration_minutes?: number;
  expiration_hours?: number;
  expiration_days?: number;
  expiration_message?: string;
}