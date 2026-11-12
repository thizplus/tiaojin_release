import type {
  PeriodStatus,
  DepositType,
  RateType,
  CommissionType,
  FlightSector,
  RateItem,
  CommissionItem,
  RoomType,
  AdditionalService,
  PeriodFile
} from '@/types/models/period';

export interface PeriodListRequest {
  page?: number;
  limit?: number;
  tour_id?: string;
  status?: PeriodStatus;
  date_from?: string;
  date_to?: string;
  search?: string;
  sort?: string;
  order?: 'asc' | 'desc';
}

export interface CreatePeriodRequest {
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
  rate_items: Omit<RateItem, 'id' | 'sort_order'>[];
  commission_type: CommissionType;
  commission_items: CommissionItem[];
  room_types?: Omit<RoomType, 'id' | 'sort_order'>[];
  additional_services?: Omit<AdditionalService, 'id' | 'sort_order'>[];
  files?: PeriodFile[];
  status?: PeriodStatus;
  service_fee_per_pax?: number;

  // Expiration Override - สำหรับทัวร์ไฟไหม้
  expiration_override?: boolean;
  expiration_minutes?: number;
  expiration_hours?: number;
  expiration_days?: number;
  expiration_message?: string;
}

export interface UpdatePeriodRequest {
  start_date?: string;
  end_date?: string;
  seat?: number;
  deposit?: number;
  deposit_type?: DepositType;
 // single_stay?: number;
  airline?: string;
  airport?: string;
  note?: string;
  sector_1?: FlightSector[];
  sector_2?: FlightSector[];
  rate_type?: RateType;
  rate_items?: Omit<RateItem, 'id' | 'sort_order'>[];
  commission_type?: CommissionType;
  commission_items?: CommissionItem[];
  room_types?: Omit<RoomType, 'id' | 'sort_order'>[];
  additional_services?: Omit<AdditionalService, 'id' | 'sort_order'>[];
  files?: PeriodFile[];
  service_fee_per_pax?: number;

  // Expiration Override - สำหรับทัวร์ไฟไหม้
  expiration_override?: boolean;
  expiration_minutes?: number;
  expiration_hours?: number;
  expiration_days?: number;
  expiration_message?: string;
}

export interface UpdatePeriodStatusRequest {
  status: PeriodStatus;
}

export interface CheckAvailabilityRequest {
  requested_seats?: number;
}