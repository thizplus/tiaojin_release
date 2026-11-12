// src/services/cityService.ts
import { API } from '@/constants/api';
import apiService from './apiService';
import type { CitiesListResponse, CityResponse } from '@/types/response/city';
import type { CityFilter } from '@/types/filter';

/**
 * Service สำหรับจัดการข้อมูลเมือง (Public)
 */
const cityService = {
  /**
   * ดึงรายการเมืองที่ publish แล้ว (Public endpoint)
   * @param filters - พารามิเตอร์สำหรับการกรอง
   * @returns รายการเมือง
   */
  getPublicCities: async (filters: CityFilter = {}): Promise<CitiesListResponse> => {
    const queryParams = new URLSearchParams();

    if (filters.search) queryParams.append('search', filters.search);
    if (filters.status) queryParams.append('status', filters.status);
    if (filters.stat !== undefined) queryParams.append('stat', filters.stat.toString());
    if (filters.page !== undefined) queryParams.append('page', filters.page.toString());
    if (filters.limit !== undefined) queryParams.append('limit', filters.limit.toString());
    if (filters.sort) queryParams.append('sort', filters.sort);
    if (filters.order) queryParams.append('order', filters.order);

    const url = queryParams.toString()
      ? `${API.CITY.PUBLIC_LIST}?${queryParams.toString()}`
      : API.CITY.PUBLIC_LIST;

    return await apiService.get<CitiesListResponse>(url);
  },

  /**
   * ดึงข้อมูลเมืองตาม ID (Public endpoint)
   * @param id - ID ของเมือง
   * @returns ข้อมูลเมือง
   */
  getPublicCityById: async (id: string): Promise<CityResponse> => {
    return await apiService.get<CityResponse>(API.CITY.PUBLIC_GET_BY_ID(id));
  },
};

export default cityService;