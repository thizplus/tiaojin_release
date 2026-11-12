//src/services/admin/bannerService.ts

import { API } from '@/constants/api';
import apiService from '@/services/apiService';
import {
  type Banner,
  type BannerFilter,
  type CreateBannerRequest,
  type UpdateBannerRequest,
  type UpdateBannerStatusRequest,
  type UpdateBannerOrderRequest,
  type PaginatedResponse,
  type ApiResponse
} from '@/types';

/**
 * Service สำหรับจัดการข้อมูลแบนเนอร์ (Admin)
 */
const bannerService = {
  /**
   * ดึงรายการแบนเนอร์ทั้งหมด (สำหรับ admin+)
   * GET /api/v1/admin/banners
   * @param params - ตัวกรองและการแบ่งหน้า
   * @returns รายการแบนเนอร์
   */
  getBanners: async (params: BannerFilter = {}): Promise<PaginatedResponse<Banner>> => {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append('search', params.search);
    if (params.status) queryParams.append('status', params.status);
    if (params.position) queryParams.append('position', params.position);
    if (params.page !== undefined) queryParams.append('page', params.page.toString());
    if (params.limit !== undefined) queryParams.append('limit', params.limit.toString());
    if (params.sort) queryParams.append('sort', params.sort);
    if (params.order) queryParams.append('order', params.order);

    const url = queryParams.toString() ? `${API.BANNER.LIST}?${queryParams.toString()}` : API.BANNER.LIST;
    return await apiService.get<PaginatedResponse<Banner>>(url);
  },

  /**
   * สร้างแบนเนอร์ใหม่ (สำหรับ admin+)
   * POST /api/v1/admin/banners
   * @param data - ข้อมูลแบนเนอร์ที่ต้องการสร้าง
   * @returns ข้อมูลแบนเนอร์ที่สร้างแล้ว
   */
  createBanner: async (data: CreateBannerRequest): Promise<ApiResponse<Banner>> => {
    return await apiService.post<ApiResponse<Banner>>(API.BANNER.CREATE, data);
  },

  /**
   * ดึงข้อมูลแบนเนอร์ตาม ID (สำหรับ admin+)
   * GET /api/v1/admin/banners/:id
   * @param id - ID ของแบนเนอร์
   * @returns ข้อมูลแบนเนอร์
   */
  getBannerById: async (id: string): Promise<ApiResponse<Banner>> => {
    return await apiService.get<ApiResponse<Banner>>(API.BANNER.GET_BY_ID(id));
  },

  /**
   * อัปเดตข้อมูลแบนเนอร์ (สำหรับ admin+)
   * PUT /api/v1/admin/banners/:id
   * @param id - ID ของแบนเนอร์
   * @param data - ข้อมูลที่ต้องการอัปเดต
   * @returns ข้อมูลแบนเนอร์ที่อัปเดตแล้ว
   */
  updateBanner: async (id: string, data: UpdateBannerRequest): Promise<ApiResponse<Banner>> => {
    return await apiService.put<ApiResponse<Banner>>(API.BANNER.UPDATE(id), data);
  },

  /**
   * ลบแบนเนอร์ (สำหรับ admin+)
   * DELETE /api/v1/admin/banners/:id
   * @param id - ID ของแบนเนอร์
   * @returns ผลลัพธ์การลบ
   */
  deleteBanner: async (id: string): Promise<ApiResponse<null>> => {
    return await apiService.delete<ApiResponse<null>>(API.BANNER.DELETE(id));
  },

  /**
   * อัปเดตสถานะแบนเนอร์ (สำหรับ admin+)
   * PUT /api/v1/admin/banners/:id/status
   * @param id - ID ของแบนเนอร์
   * @param data - สถานะใหม่
   * @returns ข้อมูลแบนเนอร์ที่อัปเดตแล้ว
   */
  updateBannerStatus: async (id: string, data: UpdateBannerStatusRequest): Promise<ApiResponse<Banner>> => {
    return await apiService.put<ApiResponse<Banner>>(API.BANNER.UPDATE_STATUS(id), data);
  },

  /**
   * อัปเดตลำดับการแสดงผลแบนเนอร์ (สำหรับ admin+)
   * PUT /api/v1/admin/banners/:id/order
   * @param id - ID ของแบนเนอร์
   * @param data - ลำดับใหม่
   * @returns ข้อมูลแบนเนอร์ที่อัปเดตแล้ว
   */
  updateBannerOrder: async (id: string, data: UpdateBannerOrderRequest): Promise<ApiResponse<Banner>> => {
    return await apiService.put<ApiResponse<Banner>>(API.BANNER.UPDATE_ORDER(id), data);
  },
};

export default bannerService;