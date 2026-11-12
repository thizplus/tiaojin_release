//src/services/public/bannerService.ts

import { API } from '@/constants/api';
import apiService from '@/services/apiService';
import type { Banner, BannerFilter, PaginatedResponse } from '@/types';

/**
 * Service สำหรับดึงข้อมูลแบนเนอร์แบบ Public (ไม่ต้อง auth)
 */
const publicBannerService = {
  /**
   * ดึงรายการแบนเนอร์ที่ active (Public)
   * GET /api/v1/banners
   * @param params - ตัวกรอง
   * @returns รายการแบนเนอร์
   */
  getBanners: async (params: BannerFilter = {}): Promise<PaginatedResponse<Banner>> => {
    const queryParams = new URLSearchParams();

    if (params.is_active !== undefined) queryParams.append('is_active', params.is_active.toString());
    if (params.position) queryParams.append('position', params.position);
    if (params.limit !== undefined) queryParams.append('limit', params.limit.toString());
    if (params.page !== undefined) queryParams.append('page', params.page.toString());

    const url = queryParams.toString()
      ? `${API.BANNER.PUBLIC_LIST}?${queryParams.toString()}`
      : API.BANNER.PUBLIC_LIST;

    return await apiService.get<PaginatedResponse<Banner>>(url);
  },
};

export default publicBannerService;
