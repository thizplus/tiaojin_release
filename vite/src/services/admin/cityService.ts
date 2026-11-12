//src/services/admin/cityService.ts

import { API } from '@/constants/api';
import apiService from '@/services/apiService';
import {
  type City,
  type CityFilter,
  type CreateCityRequest,
  type UpdateCityRequest,
  type UpdateCityStatusRequest,
  type PaginatedResponse,
  type ApiResponse
} from '@/types';

/**
 * Service สำหรับจัดการข้อมูลเมือง (Admin)
 */
const cityService = {
  /**
   * ดึงรายการเมืองทั้งหมด (สำหรับ admin+)
   * GET /api/v1/admin/cities
   * @param params - ตัวกรองและการแบ่งหน้า
   * @returns รายการเมือง
   */
  getCities: async (params: CityFilter = {}): Promise<PaginatedResponse<City>> => {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append('search', params.search);
    if (params.status) queryParams.append('status', params.status);
    if (params.page !== undefined) queryParams.append('page', params.page.toString());
    if (params.limit !== undefined) queryParams.append('limit', params.limit.toString());
    if (params.sort) queryParams.append('sort', params.sort);
    if (params.order) queryParams.append('order', params.order);

    const url = queryParams.toString() ? `${API.CITY.LIST}?${queryParams.toString()}` : API.CITY.LIST;
    return await apiService.get<PaginatedResponse<City>>(url);
  },

  /**
   * สร้างเมืองใหม่ (สำหรับ admin+)
   * POST /api/v1/admin/cities
   * @param data - ข้อมูลเมืองที่ต้องการสร้าง
   * @returns ข้อมูลเมืองที่สร้างแล้ว
   */
  createCity: async (data: CreateCityRequest): Promise<ApiResponse<City>> => {
    return await apiService.post<ApiResponse<City>>(API.CITY.CREATE, data);
  },

  /**
   * ดึงข้อมูลเมืองตาม ID (สำหรับ admin+)
   * GET /api/v1/admin/cities/:id
   * @param id - ID ของเมือง
   * @returns ข้อมูลเมือง
   */
  getCityById: async (id: string): Promise<ApiResponse<City>> => {
    return await apiService.get<ApiResponse<City>>(API.CITY.GET_BY_ID(id));
  },

  /**
   * อัปเดตข้อมูลเมือง (สำหรับ admin+)
   * PUT /api/v1/admin/cities/:id
   * @param id - ID ของเมือง
   * @param data - ข้อมูลที่ต้องการอัปเดต
   * @returns ข้อมูลเมืองที่อัปเดตแล้ว
   */
  updateCity: async (id: string, data: UpdateCityRequest): Promise<ApiResponse<City>> => {
    return await apiService.put<ApiResponse<City>>(API.CITY.UPDATE(id), data);
  },

  /**
   * ลบเมือง (สำหรับ admin+)
   * DELETE /api/v1/admin/cities/:id
   * @param id - ID ของเมือง
   * @returns ผลลัพธ์การลบ
   */
  deleteCity: async (id: string): Promise<ApiResponse<null>> => {
    return await apiService.delete<ApiResponse<null>>(API.CITY.DELETE(id));
  },

  /**
   * อัปเดตสถานะเมือง (สำหรับ admin+)
   * PUT /api/v1/admin/cities/:id/status
   * @param id - ID ของเมือง
   * @param data - สถานะใหม่
   * @returns ข้อมูลเมืองที่อัปเดตแล้ว
   */
  updateCityStatus: async (id: string, data: UpdateCityStatusRequest): Promise<ApiResponse<City>> => {
    return await apiService.put<ApiResponse<City>>(API.CITY.UPDATE_STATUS(id), data);
  },
};

export default cityService;