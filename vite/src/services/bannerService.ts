// src/services/bannerService.ts
import { API } from '@/constants/api';
import apiService from './apiService';
import type { PublicBannersResponse } from '@/types/response/banner';

/**
 * Service สำหรับจัดการข้อมูล Banner (Public)
 */
const bannerService = {
  /**
   * ดึงรายการแบนเนอร์ที่ publish แล้ว (Public endpoint)
   * @returns รายการแบนเนอร์
   */
  getPublicBanners: async (position: string = 'internal'): Promise<PublicBannersResponse> => {
    return await apiService.get<PublicBannersResponse>(`${API.BANNER.PUBLIC_LIST}?position=${position}`);
  },
};

export default bannerService;