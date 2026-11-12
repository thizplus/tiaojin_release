import { useQuery } from '@tanstack/react-query';
import bannerService from '@/services/bannerService';
import type { Banner } from '@/types/models/banner';

interface UsePublicBannersReturn {
  banners: Banner[];
  isLoading: boolean;
  isError: boolean;
  error: Error | null;
  refetch: () => void;
}

export const usePublicBanners = (
  options: { enabled?: boolean } = {}
): UsePublicBannersReturn => {
  const {
    data: bannersData,
    isLoading,
    isError,
    error,
    refetch
  } = useQuery({
    queryKey: ['public-banners', 'internal'],
    queryFn: async () => {
      const response = await bannerService.getPublicBanners('internal');
      return response;
    },
    enabled: options.enabled !== false,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });

  const banners = bannersData?.success ? bannersData.data : [];

  return {
    banners,
    isLoading,
    isError,
    error: error as Error | null,
    refetch,
  };
};

export default usePublicBanners;