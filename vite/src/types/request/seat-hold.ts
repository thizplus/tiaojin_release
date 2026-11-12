import type { SeatHold, PeriodSeatHoldsData } from '@/types/models/seat-hold';
import type { ApiResponse } from '../common';


// Request Types
export interface CreateSeatHoldRequest {
  held_seats: number;
  reason?: string;
  expires_at?: string;
}

export interface ReleaseSeatHoldRequest {
  reason?: string;
}

export interface GetPeriodSeatHoldsRequest {
  include_inactive?: boolean;
}

// Response Types
export interface CreateSeatHoldResponse extends ApiResponse {
  data: SeatHold;
}

export interface GetPeriodSeatHoldsResponse extends ApiResponse {
  data: PeriodSeatHoldsData;
}

export interface ReleaseSeatHoldResponse extends ApiResponse {
  data: SeatHold;
}