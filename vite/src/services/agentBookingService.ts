// src/services/agentBookingService.ts
import { API } from '@/constants/api';
import apiService from './apiService';

import type {
  AdminBookingListResponse,
  AdminBookingListFilters,
  BookingResponse,
  ApiResponse,
  AgentCancelBookingRequest,
  AgentCancelBookingResponse,
  AgentUploadPaymentSlipRequest,
  AgentUploadPaymentSlipResponse,
  AgentUploadPassportRequest,
  AgentUploadPassportResponse,
  AgentBookingDocumentsResponse
} from '@/types';

/**
 * Service สำหรับจัดการการจองของ Agent (ดึงข้อมูลการจองของตัวเอง)
 */
const agentBookingService = {
  /**
   * ดึงรายการการจองของ Agent ตัวเอง (มี pagination และ filters)
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

    const url = queryParams.toString() ? `${API.AGENT_BOOKING.LIST}?${queryParams.toString()}` : API.AGENT_BOOKING.LIST;
    return await apiService.get<AdminBookingListResponse>(url);
  },

  /**
   * ดึงข้อมูลการจองตาม ID สำหรับ Agent (เฉพาะการจองของตัวเอง)
   * @param id - ID ของการจอง
   * @returns ข้อมูลการจอง
   */
  getById: async (id: string): Promise<ApiResponse<BookingResponse>> => {
    return await apiService.get<ApiResponse<BookingResponse>>(API.AGENT_BOOKING.GET_BY_ID(id));
  },

  /**
   * ยกเลิกการจองสำหรับ Agent (เฉพาะการจองของตัวเอง)
   * @param id - ID ของการจอง
   * @param data - ข้อมูลการยกเลิก (เหตุผล)
   * @returns ผลลัพธ์การยกเลิก
   */
  cancel: async (id: string, data: AgentCancelBookingRequest = {}): Promise<AgentCancelBookingResponse> => {
    return await apiService.deleteWithData<AgentCancelBookingResponse>(API.AGENT_BOOKING.CANCEL(id), data);
  },

  /**
   * อัปโหลดสลิปการโอนเงินสำหรับการจอง
   * @param id - ID ของการจอง
   * @param data - ข้อมูลสลิปและจำนวนเงิน
   * @returns ข้อมูลสลิปที่อัปโหลด
   */
  uploadPaymentSlip: async (id: string, data: AgentUploadPaymentSlipRequest): Promise<AgentUploadPaymentSlipResponse> => {
    const formData = new FormData();
    formData.append('payment_slip', data.payment_slip);
    formData.append('amount', data.amount.toString());

    return await apiService.postFile<AgentUploadPaymentSlipResponse>(API.AGENT_BOOKING.UPLOAD_PAYMENT_SLIP(id), formData);
  },

  /**
   * อัปโหลดหนังสือเดินทางสำหรับการจอง
   * @param id - ID ของการจอง
   * @param data - ข้อมูลหนังสือเดินทาง
   * @returns ข้อมูลหนังสือเดินทางที่อัปโหลด
   */
  uploadPassport: async (id: string, data: AgentUploadPassportRequest): Promise<AgentUploadPassportResponse> => {
    const formData = new FormData();
    formData.append('passport', data.passport);

    return await apiService.postFile<AgentUploadPassportResponse>(API.AGENT_BOOKING.UPLOAD_PASSPORT(id), formData);
  },

  /**
   * ดึงเอกสารการจองทั้งหมด (สลิป + หนังสือเดินทาง)
   * @param id - ID ของการจอง
   * @returns เอกสารการจองทั้งหมด
   */
  getDocuments: async (id: string): Promise<AgentBookingDocumentsResponse> => {
    return await apiService.get<AgentBookingDocumentsResponse>(API.AGENT_BOOKING.GET_DOCUMENTS(id));
  }
};

export default agentBookingService;