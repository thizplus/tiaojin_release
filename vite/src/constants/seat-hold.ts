import type { SeatHoldStatus } from '@/types/models/seat-hold';

export const SEAT_HOLD_STATUS: Record<SeatHoldStatus, string> = {
  active: 'ยังคงจอง',
  released: 'ปล่อยแล้ว',
  expired: 'หมดอายุ'
} as const;

export const SEAT_HOLD_STATUS_COLORS: Record<SeatHoldStatus, string> = {
  active: 'bg-muted text-muted-foreground',
  released: 'bg-muted text-muted-foreground',
  expired: 'bg-muted text-muted-foreground'
} as const;

export const SEAT_HOLD_CONSTRAINTS = {
  MIN_HELD_SEATS: 1,
  MAX_HELD_SEATS: 50,
  DEFAULT_HOLD_DURATION_HOURS: 24
} as const;

export const SEAT_HOLD_REASONS = [
  'ลูกค้าขอจองล่วงหน้า',
  'รอการชำระเงิน',
  'ประชุมทีมขาย',
  'จองสำหรับลูกค้า VIP',
  'รอการยืนยันจากลูกค้า'
] as const;