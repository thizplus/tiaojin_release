// src/services/saleBookingService.ts
import { API } from '@/constants/api';
import apiService from './apiService';

import type {
  AdminBookingListResponse,
  AdminBookingListFilters,
  BookingResponse,
  ApiResponse,
  AgentUploadPaymentSlipResponse,
  AgentUploadPassportResponse,
} from '@/types';

/**
 * Service สำหรับจัดการการจองของ Sale (ดึงข้อมูลการจองของตัวเองแบบอ่านอย่างเดียว)
 */
const saleBookingService = {
  /**
   * ดึงรายการการจองของ Sale ตัวเอง (มี pagination และ filters) - อ่านอย่างเดียว
   * @param filters - พารามิเตอร์สำหรับการกรองและการแบ่งหน้า
   * @returns รายการการจอง
   */
  list: async (filters: AdminBookingListFilters = {}): Promise<AdminBookingListResponse> => {
    const queryParams = new URLSearchParams();

    if (filters.search) queryParams.append('search', filters.search);
    if (filters.status) queryParams.append('status', filters.status);
    if (filters.tour_id) queryParams.append('tour_id', filters.tour_id);
    if (filters.period_id) queryParams.append('period_id', filters.period_id);
    if (filters.sales_id) queryParams.append('sales_id', filters.sales_id);
    if (filters.start_date) queryParams.append('start_date', filters.start_date);
    if (filters.end_date) queryParams.append('end_date', filters.end_date);
    if (filters.page !== undefined) queryParams.append('page', filters.page.toString());
    if (filters.limit !== undefined) queryParams.append('limit', filters.limit.toString());
    if (filters.sort) queryParams.append('sort', filters.sort);
    if (filters.order) queryParams.append('order', filters.order);

    const url = queryParams.toString() ? `${API.SALE_BOOKING.LIST}?${queryParams.toString()}` : API.SALE_BOOKING.LIST;
    return await apiService.get<AdminBookingListResponse>(url);
  },

  /**
   * ดึงข้อมูลการจองตาม ID สำหรับ Sale (เฉพาะการจองของตัวเอง) - อ่านอย่างเดียว
   * @param id - ID ของการจอง
   * @returns ข้อมูลการจอง
   */
  getById: async (id: string): Promise<ApiResponse<BookingResponse>> => {
    return await apiService.get<ApiResponse<BookingResponse>>(API.SALE_BOOKING.GET_BY_ID(id));
  },

  /**
   * อัปโหลดสลิปการชำระเงินสำหรับการจอง (แทนลูกค้า)
   * @param id - ID ของการจอง
   * @param file - ไฟล์สลิป
   * @param amount - จำนวนเงิน
   * @returns ผลลัพธ์จากการอัปโหลด
   */
  uploadPaymentSlip: async (id: string, file: File, amount: number): Promise<AgentUploadPaymentSlipResponse> => {
    const formData = new FormData();
    formData.append('payment_slip', file);
    formData.append('amount', amount.toString());
    return await apiService.upload<AgentUploadPaymentSlipResponse>(API.SALE_BOOKING.UPLOAD_PAYMENT_SLIP(id), formData);
  },

  /**
   * อัปโหลดหนังสือเดินทางสำหรับการจอง (แทนลูกค้า)
   * @param id - ID ของการจอง
   * @param file - ไฟล์หนังสือเดินทาง
   * @returns ผลลัพธ์จากการอัปโหลด
   */
  uploadPassport: async (id: string, file: File): Promise<AgentUploadPassportResponse> => {
    const formData = new FormData();
    formData.append('passport', file);
    return await apiService.upload<AgentUploadPassportResponse>(API.SALE_BOOKING.UPLOAD_PASSPORT(id), formData);
  }
};

export default saleBookingService;