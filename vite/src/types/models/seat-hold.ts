export type SeatHoldStatus = 'active' | 'released' | 'expired';

export interface SeatHold {
  id: string;
  period_id: string;
  admin_id: string;
  held_seats: number;
  reason?: string;
  expires_at?: string;
  status: SeatHoldStatus;
  created_at: string;
  updated_at: string;
  time_remaining?: string;
  is_expired: boolean;
  can_extend: boolean;
  can_release: boolean;
}

export interface PeriodSeatHoldsData {
  period_id: string;
  seat_holds: SeatHold[];
  total: number;
  include_inactive: boolean;
}