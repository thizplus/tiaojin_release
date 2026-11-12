import apiService from '@/services/apiService';
import { API } from '@/constants/api';
import type {
  PeriodListRequest,
  CreatePeriodRequest,
  UpdatePeriodRequest,
  UpdatePeriodStatusRequest,
  CheckAvailabilityRequest,
} from '@/types/request/period';
import type {
  PeriodListResponse,
  PeriodResponse,
  CreatePeriodResponse,
  UpdatePeriodResponse,
  DeletePeriodResponse,
  UpdateStatusResponse,
  AvailabilityResponse,
} from '@/types/response/period';

class PeriodService {

  /**
   * ดึงรายการ periods ทั้งหมดพร้อม pagination และ filtering
   */
  async getAll(params: PeriodListRequest = {}): Promise<PeriodListResponse> {
    // Build query string manually for clean URL format
    const queryParams = new URLSearchParams();

    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryParams.append(key, String(value));
      }
    });

    const queryString = queryParams.toString();
    const url = queryString ? `${API.PERIOD.LIST}?${queryString}` : API.PERIOD.LIST;

    return await apiService.get<PeriodListResponse>(url);
  }

  /**
   * ดึงข้อมูล period ตาม ID
   */
  async getById(id: string): Promise<PeriodResponse> {
    return await apiService.get<PeriodResponse>(API.PERIOD.GET_BY_ID(id));
  }

  /**
   * สร้าง period ใหม่
   */
  async create(data: CreatePeriodRequest): Promise<CreatePeriodResponse> {
    return await apiService.post<CreatePeriodResponse>(API.PERIOD.CREATE, data);
  }

  /**
   * อัปเดต period
   */
  async update(id: string, data: UpdatePeriodRequest): Promise<UpdatePeriodResponse> {
    return await apiService.put<UpdatePeriodResponse>(API.PERIOD.UPDATE(id), data);
  }

  /**
   * ลบ period
   */
  async delete(id: string): Promise<DeletePeriodResponse> {
    return await apiService.delete<DeletePeriodResponse>(API.PERIOD.DELETE(id));
  }

  /**
   * อัปเดตสถานะ period
   */
  async updateStatus(id: string, data: UpdatePeriodStatusRequest): Promise<UpdateStatusResponse> {
    return await apiService.put<UpdateStatusResponse>(API.PERIOD.UPDATE_STATUS(id), data);
  }

  /**
   * ตรวจสอบความพร้อมใช้งานของ period
   */
  async checkAvailability(id: string, params: CheckAvailabilityRequest = {}): Promise<AvailabilityResponse> {
    return await apiService.get<AvailabilityResponse>(API.PERIOD.CHECK_AVAILABILITY(id), {
      params,
    });
  }

  /**
   * ดึงรายการ periods ตาม tour_id (helper method)
   */
  async getByTourId(tourId: string, params: Omit<PeriodListRequest, 'tour_id'> = {}): Promise<PeriodListResponse> {
    return this.getAll({ ...params, tour_id: tourId });
  }
}

export const periodService = new PeriodService();