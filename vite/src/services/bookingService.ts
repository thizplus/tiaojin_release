// src/services/bookingService.ts
import { API } from '@/constants/api';
import apiService from './apiService';

import type {
  CreateBookingRequest,
  BookingCreateResponse,
  BookingResponse,
  ApiResponse
} from '@/types';

/**
 * Service สำหรับจัดการการจองทัวร์
 */
const bookingService = {
  /**
   * สร้างการจองใหม่ (Agent+ access level)
   * @param data - ข้อมูลการจอง
   * @returns ข้อมูลการจองที่สร้างขึ้น
   */
  create: async (data: CreateBookingRequest): Promise<BookingCreateResponse> => {
    return await apiService.post<BookingCreateResponse>(API.BOOKING.CREATE, data);
  },

  /**
   * ดึงรายการการจองทั้งหมด
   * @param params - พารามิเตอร์สำหรับการกรองและการแบ่งหน้า
   * @returns รายการการจอง
   */
  list: async (params?: Record<string, unknown>): Promise<ApiResponse<BookingResponse[]>> => {
    return await apiService.get<ApiResponse<BookingResponse[]>>(API.BOOKING.LIST, { params });
  },

  /**
   * ดึงข้อมูลการจองตาม ID
   * @param id - ID ของการจอง
   * @returns ข้อมูลการจอง
   */
  getById: async (id: string): Promise<ApiResponse<BookingResponse>> => {
    return await apiService.get<ApiResponse<BookingResponse>>(API.BOOKING.GET_BY_ID(id));
  },

  /**
   * อัพเดตข้อมูลการจอง
   * @param id - ID ของการจอง
   * @param data - ข้อมูลที่ต้องการอัพเดต
   * @returns ข้อมูลการจองที่อัพเดตแล้ว
   */
  update: async (id: string, data: Partial<CreateBookingRequest>): Promise<ApiResponse<BookingResponse>> => {
    return await apiService.put<ApiResponse<BookingResponse>>(API.BOOKING.UPDATE(id), data);
  },

  /**
   * ลบการจอง
   * @param id - ID ของการจอง
   * @returns ผลลัพธ์จากการลบ
   */
  delete: async (id: string): Promise<ApiResponse<void>> => {
    return await apiService.delete<ApiResponse<void>>(API.BOOKING.DELETE(id));
  },

  /**
   * อัพเดตสถานะการจอง
   * @param id - ID ของการจอง
   * @param status - สถานะใหม่
   * @returns ข้อมูลการจองที่อัพเดตแล้ว
   */
  updateStatus: async (id: string, status: 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed'): Promise<ApiResponse<BookingResponse>> => {
    return await apiService.put<ApiResponse<BookingResponse>>(API.BOOKING.UPDATE_STATUS(id), { status });
  }
};

export default bookingService;