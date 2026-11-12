// src/services/reportService.ts
import { API } from '@/constants/api';
import apiService from './apiService';

import type {
  AdminReportResponse,
  AdminReportParams,
  SalesReportResponse,
  SalesReportParams,
  AgentReportResponse,
  AgentReportParams,
} from '@/types/report';

/**
 * Service สำหรับจัดการรายงาน
 */
const reportService = {
  /**
   * ดึงรายงานการจองรายเดือนสำหรับ Admin
   * @param params - พารามิเตอร์ปีและเดือน
   * @returns ข้อมูลรายงาน
   */
  getAdminReport: async (params: AdminReportParams): Promise<AdminReportResponse> => {
    const queryParams = new URLSearchParams();
    queryParams.append('year', params.year.toString());
    queryParams.append('month', params.month.toString());

    const url = `${API.ADMIN_BOOKING.REPORT}?${queryParams.toString()}`;
    return await apiService.get<AdminReportResponse>(url);
  },

  /**
   * ดึงรายงานการจองรายเดือนสำหรับ Sales
   * @param params - พารามิเตอร์ปีและเดือน
   * @returns ข้อมูลรายงาน
   */
  getSalesReport: async (params: SalesReportParams): Promise<SalesReportResponse> => {
    const queryParams = new URLSearchParams();
    queryParams.append('year', params.year.toString());
    queryParams.append('month', params.month.toString());

    const url = `${API.SALE_BOOKING.REPORT}?${queryParams.toString()}`;
    return await apiService.get<SalesReportResponse>(url);
  },

  /**
   * ดึงรายงานการจองรายเดือนสำหรับ Agent
   * @param params - พารามิเตอร์ปีและเดือน
   * @returns ข้อมูลรายงาน
   */
  getAgentReport: async (params: AgentReportParams): Promise<AgentReportResponse> => {
    const queryParams = new URLSearchParams();
    queryParams.append('year', params.year.toString());
    queryParams.append('month', params.month.toString());

    const url = `${API.AGENT_REPORT.REPORT}?${queryParams.toString()}`;
    return await apiService.get<AgentReportResponse>(url);
  },
};

export default reportService;