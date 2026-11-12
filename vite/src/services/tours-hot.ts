import apiService from '@/services/apiService';
import { TOUR_API } from '@/constants/api';
import type {
  ToursRegularResponse,
  ToursRegularRequest,
} from '@/types/tours-regular';

export class ToursHotService {
  /**
   * Get hot/special tours list
   */
  async getHotTours(params?: ToursRegularRequest): Promise<ToursRegularResponse> {
    const searchParams = new URLSearchParams();

    if (params?.page !== undefined) {
      searchParams.append('page', params.page.toString());
    }
    if (params?.limit !== undefined) {
      searchParams.append('limit', params.limit.toString());
    }
    if (params?.search) {
      searchParams.append('search', params.search);
    }
    if (params?.status) {
      searchParams.append('status', params.status);
    }
    if (params?.city_id) {
      searchParams.append('city_id', params.city_id);
    }
    if (params?.available_only !== undefined) {
      searchParams.append('available_only', params.available_only.toString());
    }

    const url = searchParams.toString()
      ? `${TOUR_API.SPECIAL}?${searchParams.toString()}`
      : TOUR_API.SPECIAL;

    return await apiService.get<ToursRegularResponse>(url);
  }
}

export const toursHotService = new ToursHotService();