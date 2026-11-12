import { TOUR_API as API } from '@/constants/api';
import apiService from '../apiService';
import type {
  Tour,
  TourFilter,
  CreateTourRequest,
  UpdateTourRequest,
  UpdateTourStatusRequest,
  TourPeriodsResponse,
  TourPeriodsFilter,
  ApiResponse,
  PaginationMeta
} from '@/types';

const tourService = {
  async getTours(params: TourFilter = {}): Promise<ApiResponse<{ tours: Tour[], pagination: PaginationMeta }>> {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append("search", params.search);
    if (params.status) queryParams.append("status", params.status);
    if (params.city_id) queryParams.append("city_id", params.city_id);
    if (params.page !== undefined) queryParams.append("page", params.page.toString());
    if (params.limit !== undefined) queryParams.append("limit", params.limit.toString());
    if (params.sort) queryParams.append("sort", params.sort);
    if (params.order) queryParams.append("order", params.order);
    if (params.with_periods !== undefined) queryParams.append("with_periods", params.with_periods.toString());

    const url = queryParams.toString() ? `${API.LIST}?${queryParams.toString()}` : API.LIST;
    return await apiService.get<ApiResponse<{ tours: Tour[], pagination: PaginationMeta }>>(url);
  },

  async getTourById(id: string): Promise<ApiResponse<Tour>> {
    return await apiService.get<ApiResponse<Tour>>(API.GET_BY_ID(id));
  },

  async createTour(data: CreateTourRequest): Promise<ApiResponse<Tour>> {
    return await apiService.post<ApiResponse<Tour>>(API.CREATE, data);
  },

  async updateTour(id: string, data: UpdateTourRequest): Promise<ApiResponse<Tour>> {
    return await apiService.put<ApiResponse<Tour>>(API.UPDATE(id), data);
  },

  async deleteTour(id: string): Promise<ApiResponse<null>> {
    return await apiService.delete<ApiResponse<null>>(API.DELETE(id));
  },

  async updateTourStatus(id: string, data: UpdateTourStatusRequest): Promise<ApiResponse<Tour>> {
    return await apiService.patch<ApiResponse<Tour>>(API.UPDATE_STATUS(id), data);
  },

  async getTourPeriods(tourId: string, params: TourPeriodsFilter = {}): Promise<ApiResponse<TourPeriodsResponse>> {
    const queryParams = new URLSearchParams();

    if (params.page !== undefined) queryParams.append("page", params.page.toString());
    if (params.limit !== undefined) queryParams.append("limit", params.limit.toString());
    if (params.sort) queryParams.append("sort", params.sort);
    if (params.order) queryParams.append("order", params.order);
    if (params.status) queryParams.append("status", params.status);
    if (params.search) queryParams.append("search", params.search);
    if (params.start_date_from) queryParams.append("start_date_from", params.start_date_from);
    if (params.start_date_to) queryParams.append("start_date_to", params.start_date_to);
    if (params.end_date_from) queryParams.append("end_date_from", params.end_date_from);
    if (params.end_date_to) queryParams.append("end_date_to", params.end_date_to);

    const url = queryParams.toString()
      ? `${API.GET_PERIODS(tourId)}?${queryParams.toString()}`
      : API.GET_PERIODS(tourId);

    return await apiService.get<ApiResponse<TourPeriodsResponse>>(url);
  }
};

export default tourService;