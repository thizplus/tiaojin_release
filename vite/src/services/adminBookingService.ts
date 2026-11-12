// src/services/adminBookingService.ts
import { API } from '@/constants/api';
import apiService from './apiService';

import type {
  AdminBookingListResponse,
  AdminBookingListFilters,
  AdminBookingStatsResponse,
  BookingResponse,
  AdminBookingStatusUpdateRequest,
  AdminBookingStatusUpdateResponse,
  AdminCreatePaymentRequest,
  AdminCreatePaymentResponse,
  PaymentSlipVerifyRequest,
  PaymentSlipVerifyResponse,
  PaymentSlipReverseRequest,
  PaymentSlipReverseResponse,
  ApiResponse
} from '@/types';

/**
 * Service สำหรับจัดการการจองของ Admin
 */
const adminBookingService = {
  /**
   * ดึงรายการการจองทั้งหมดสำหรับ Admin (มี pagination และ filters)
   * @param filters - พารามิเตอร์สำหรับการกรองและการแบ่งหน้า
   * @returns รายการการจอง
   */
  list: async (filters: AdminBookingListFilters = {}): Promise<AdminBookingListResponse> => {
    const queryParams = new URLSearchParams();

    if (filters.search) queryParams.append('search', filters.search);
    if (filters.status) queryParams.append('status', filters.status);
    if (filters.payment_slip_status) queryParams.append('payment_slip_status', filters.payment_slip_status);
    if (filters.user_id) queryParams.append('user_id', filters.user_id);
    if (filters.tour_id) queryParams.append('tour_id', filters.tour_id);
    if (filters.period_id) queryParams.append('period_id', filters.period_id);
    if (filters.sales_id) queryParams.append('sales_id', filters.sales_id);
    if (filters.start_date) queryParams.append('start_date', filters.start_date);
    if (filters.end_date) queryParams.append('end_date', filters.end_date);
    if (filters.page !== undefined) queryParams.append('page', filters.page.toString());
    if (filters.limit !== undefined) queryParams.append('limit', filters.limit.toString());
    if (filters.sort) queryParams.append('sort', filters.sort);
    if (filters.order) queryParams.append('order', filters.order);

    const url = queryParams.toString() ? `${API.ADMIN_BOOKING.LIST}?${queryParams.toString()}` : API.ADMIN_BOOKING.LIST;
    return await apiService.get<AdminBookingListResponse>(url);
  },

  /**
   * ดึงสถิติการจองสำหรับ Admin
   * @returns สถิติการจอง
   */
  getStats: async (): Promise<AdminBookingStatsResponse> => {
    return await apiService.get<AdminBookingStatsResponse>(API.ADMIN_BOOKING.STATS);
  },

  /**
   * ดึงข้อมูลการจองตาม ID สำหรับ Admin
   * @param id - ID ของการจอง
   * @returns ข้อมูลการจอง
   */
  getById: async (id: string): Promise<ApiResponse<BookingResponse>> => {
    return await apiService.get<ApiResponse<BookingResponse>>(API.ADMIN_BOOKING.GET_BY_ID(id));
  },

  /**
   * อัพเดตสถานะการจองสำหรับ Admin
   * @param id - ID ของการจอง
   * @param data - ข้อมูลสถานะใหม่
   * @returns ข้อมูลการจองที่อัพเดตแล้ว
   */
  updateStatus: async (id: string, data: AdminBookingStatusUpdateRequest): Promise<AdminBookingStatusUpdateResponse> => {
    return await apiService.put<AdminBookingStatusUpdateResponse>(API.ADMIN_BOOKING.UPDATE_STATUS(id), data);
  },

  /**
   * สร้างการชำระเงินใหม่สำหรับการจอง
   * @param bookingId - ID ของการจอง
   * @param data - ข้อมูลการชำระเงิน
   * @returns ข้อมูลการชำระเงินที่สร้างขึ้น
   */
  createPayment: async (bookingId: string, data: AdminCreatePaymentRequest): Promise<AdminCreatePaymentResponse> => {
    return await apiService.post<AdminCreatePaymentResponse>(API.ADMIN_BOOKING.CREATE_PAYMENT(bookingId), data);
  },

  /**
   * ยกเลิกการจอง (Delete)
   * @param id - ID ของการจอง
   * @returns ผลลัพธ์จากการยกเลิก
   */
  delete: async (id: string): Promise<ApiResponse<void>> => {
    return await apiService.delete<ApiResponse<void>>(API.BOOKING.DELETE(id));
  },

  /**
   * ยืนยันหรือปฏิเสธ payment slip
   * @param paymentSlipId - ID ของ payment slip
   * @param data - ข้อมูลการยืนยัน
   * @returns ผลลัพธ์จากการยืนยัน
   */
  verifyPaymentSlip: async (paymentSlipId: string, data: PaymentSlipVerifyRequest): Promise<PaymentSlipVerifyResponse> => {
    return await apiService.post<PaymentSlipVerifyResponse>(API.ADMIN_BOOKING.VERIFY_PAYMENT_SLIP(paymentSlipId), data);
  },

  /**
   * ยกเลิก payment slip ที่อนุมัติแล้ว
   * @param paymentSlipId - ID ของ payment slip
   * @param data - ข้อมูลการยกเลิก
   * @returns ผลลัพธ์จากการยกเลิก
   */
  reversePaymentSlip: async (paymentSlipId: string, data: PaymentSlipReverseRequest): Promise<PaymentSlipReverseResponse> => {
    return await apiService.post<PaymentSlipReverseResponse>(API.ADMIN_BOOKING.REVERSE_PAYMENT_SLIP(paymentSlipId), data);
  },

  /**
   * อัปโหลด passport สำหรับการจอง (Admin)
   * @param bookingId - ID ของการจอง
   * @param file - ไฟล์ passport
   * @returns ผลลัพธ์จากการอัปโหลด
   */
  uploadPassport: async (bookingId: string, file: File): Promise<ApiResponse<unknown>> => {
    const formData = new FormData();
    formData.append('passport', file);
    return await apiService.upload<ApiResponse<unknown>>(API.ADMIN_BOOKING.UPLOAD_PASSPORT(bookingId), formData);
  }
};

export default adminBookingService;