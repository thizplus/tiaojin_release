import { useQuery } from '@tanstack/react-query';
import publicBannerService from '@/services/public/bannerService';
import type { Banner, BannerFilter } from '@/types';

/**
 * Hook สำหรับดึงข้อมูลแบนเนอร์แบบ Public (ไม่ต้อง auth)
 */
export const usePublicBanner = (params: BannerFilter = {}) => {
  const query = useQuery({
    queryKey: ['public-banners', params],
    queryFn: async () => {
      const response = await publicBannerService.getBanners(params);
      if (!response.success) {
        throw new Error(response.message || 'Failed to fetch banners');
      }
      return response;
    },
  });

  // GetActive API returns flat array: { data: Banner[] }
  const rawData = query.data?.data;
  const banners = Array.isArray(rawData) ? rawData : (rawData?.banners || []);

  return {
    banners: banners as Banner[],
    isLoading: query.isLoading,
    error: query.error,
    refetch: query.refetch,
  };
};

export default usePublicBanner;
