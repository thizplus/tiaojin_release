import type { TourPeriod } from '@/types/models/period';

export interface TourDetails {
  id: string;
  title: string;
  code: string;
}

export interface PeriodWithTourDetails extends TourPeriod {
  tour_details?: TourDetails;
}

export interface PeriodListResponse {
  success: boolean;
  message: string;
  data: {
    periods: PeriodWithTourDetails[];
    pagination: {
      total: number;
      has_more: boolean;
      current_page: number;
      total_pages: number;
      page_size: number;
      limit: number;
      offset: number;
    };
  };
}

export interface PeriodResponse {
  success: boolean;
  message: string;
  data: TourPeriod;
}

export interface CreatePeriodResponse {
  success: boolean;
  message: string;
  data: TourPeriod;
}

export interface UpdatePeriodResponse {
  success: boolean;
  message: string;
  data: TourPeriod;
}

export interface DeletePeriodResponse {
  success: boolean;
  message: string;
  data: null;
}

export interface UpdateStatusResponse {
  success: boolean;
  message: string;
  data: {
    id: string;
    status: string;
    updated_at: string;
  };
}

export interface AvailabilityResponse {
  success: boolean;
  message: string;
  data: {
    available: boolean;
    requested_seats: number;
    available_seats: number;
    total_seats: number;
    booked_seats: number;
  };
}