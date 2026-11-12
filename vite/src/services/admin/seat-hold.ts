import apiService from '@/services/apiService';


import { PERIOD_API, SEAT_HOLD_API } from '@/constants/api';
import type {
  CreateSeatHoldRequest,
  CreateSeatHoldResponse,
  GetPeriodSeatHoldsRequest,
  GetPeriodSeatHoldsResponse,
  ReleaseSeatHoldRequest,
  ReleaseSeatHoldResponse,
} from '@/types/request/seat-hold';

class SeatHoldService {
  async createSeatHold(
    periodId: string,
    data: CreateSeatHoldRequest
  ): Promise<CreateSeatHoldResponse> {
    return apiService.post<CreateSeatHoldResponse>(
      PERIOD_API.HOLD_SEATS(periodId),
      data
    );
  }

  async getPeriodSeatHolds(
    periodId: string,
    params?: GetPeriodSeatHoldsRequest
  ): Promise<GetPeriodSeatHoldsResponse> {
    const searchParams = new URLSearchParams();

    if (params?.include_inactive) {
      searchParams.append('include_inactive', 'true');
    }

    const url = `${PERIOD_API.GET_SEAT_HOLDS(periodId)}${
      searchParams.toString() ? `?${searchParams.toString()}` : ''
    }`;

    return apiService.get<GetPeriodSeatHoldsResponse>(url);
  }

  async releaseSeatHold(
    holdId: string,
    data?: ReleaseSeatHoldRequest
  ): Promise<ReleaseSeatHoldResponse> {
    // Support both body and query parameter for reason
    if (data?.reason) {
      return apiService.delete<ReleaseSeatHoldResponse>(
        SEAT_HOLD_API.RELEASE(holdId),
        { data }
      );
    }

    return apiService.delete<ReleaseSeatHoldResponse>(
      SEAT_HOLD_API.RELEASE(holdId)
    );
  }

  // Helper method สำหรับ release ผ่าน query parameter
  async releaseSeatHoldWithQueryParam(
    holdId: string,
    reason?: string
  ): Promise<ReleaseSeatHoldResponse> {
    const url = reason
      ? `${SEAT_HOLD_API.RELEASE(holdId)}?reason=${encodeURIComponent(reason)}`
      : SEAT_HOLD_API.RELEASE(holdId);

    return apiService.delete<ReleaseSeatHoldResponse>(url);
  }
}

export const seatHoldService = new SeatHoldService();