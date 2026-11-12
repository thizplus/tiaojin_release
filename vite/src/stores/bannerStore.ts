import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import bannerService from '@/services/admin/bannerService';
import type {
  Banner,
  BannerFilter,
  CreateBannerRequest,
  UpdateBannerRequest,
  UpdateBannerStatusRequest,
  UpdateBannerOrderRequest,
  PaginatedResponse,
  //ApiResponse
} from '@/types';

const bannerKeys = {
  banners: (params?: BannerFilter) => ['banner', 'list', params] as const,
  bannerById: (id: string) => ['banner', 'byId', id] as const,
};

interface BannerState {
  currentBanner: Banner | null;
  selectedBanner: Banner | null;
  bannersList: Banner[];
  setCurrentBanner: (banner: Banner | null) => void;
  setSelectedBanner: (banner: Banner | null) => void;
  setBannersList: (banners: Banner[]) => void;
  clearBannerData: () => void;
}

export const useBannerStore = create<BannerState>()(
  persist(
    (set) => ({
      currentBanner: null,
      selectedBanner: null,
      bannersList: [],
      setCurrentBanner: (banner) => set({ currentBanner: banner }),
      setSelectedBanner: (banner) => set({ selectedBanner: banner }),
      setBannersList: (banners) => set({ bannersList: banners }),
      clearBannerData: () => set({
        currentBanner: null,
        selectedBanner: null,
        bannersList: [],
      }),
    }),
    {
      name: 'banner-storage',
      partialize: (state) => ({
        selectedBanner: state.selectedBanner,
      }),
    }
  )
);

// Admin+ Queries

/**
 * ดึงรายการแบนเนอร์ทั้งหมด (สำหรับ admin+)
 */
export const useBannersListQuery = (
  params: BannerFilter = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<PaginatedResponse<Banner>, Error> => {
  const { setBannersList } = useBannerStore();

  return useQuery({
    queryKey: bannerKeys.banners(params),
    queryFn: async () => {
      const response = await bannerService.getBanners(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch banners list');
      }

      // สมมติว่า response.data มี structure แบบ { banners: Banner[], pagination: PaginationMeta }
      const banners = response.data.banners || response.data;
      setBannersList(Array.isArray(banners) ? banners : []);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

/**
 * ดึงข้อมูลแบนเนอร์ตาม ID (สำหรับ admin+)
 */
export const useBannerByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<Banner, Error> => {
  const { setSelectedBanner } = useBannerStore();

  return useQuery({
    queryKey: bannerKeys.bannerById(id),
    queryFn: async () => {
      const response = await bannerService.getBannerById(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch banner data');
      }

      setSelectedBanner(response.data);
      return response.data;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Admin+ Mutations

/**
 * สร้างแบนเนอร์ใหม่ (สำหรับ admin+)
 */
export const useCreateBannerMutation = (): UseMutationResult<
  Banner,
  Error,
  CreateBannerRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateBannerRequest) => {
      const response = await bannerService.createBanner(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to create banner');
      }

      return response.data;
    },
    onSuccess: () => {
      // Invalidate banners list to refetch
      queryClient.invalidateQueries({ queryKey: bannerKeys.banners() });
    },
  });
};

/**
 * อัปเดตข้อมูลแบนเนอร์ (สำหรับ admin+)
 */
export const useUpdateBannerMutation = (): UseMutationResult<
  Banner,
  Error,
  { id: string; data: UpdateBannerRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await bannerService.updateBanner(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to update banner');
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: bannerKeys.bannerById(id) });
      queryClient.invalidateQueries({ queryKey: bannerKeys.banners() });
    },
  });
};

/**
 * ลบแบนเนอร์ (สำหรับ admin+)
 */
export const useDeleteBannerMutation = (): UseMutationResult<
  null,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await bannerService.deleteBanner(id);

      if (!response.success) {
        throw new Error(response.message || 'Failed to delete banner');
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: bannerKeys.bannerById(id) });
      queryClient.invalidateQueries({ queryKey: bannerKeys.banners() });
    },
  });
};

/**
 * อัปเดตสถานะแบนเนอร์ (สำหรับ admin+)
 */
export const useUpdateBannerStatusMutation = (): UseMutationResult<
  Banner,
  Error,
  { id: string; data: UpdateBannerStatusRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await bannerService.updateBannerStatus(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to update banner status');
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: bannerKeys.bannerById(id) });
      queryClient.invalidateQueries({ queryKey: bannerKeys.banners() });
    },
  });
};

/**
 * อัปเดตลำดับการแสดงผลแบนเนอร์ (สำหรับ admin+)
 */
export const useUpdateBannerOrderMutation = (): UseMutationResult<
  Banner,
  Error,
  { id: string; data: UpdateBannerOrderRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await bannerService.updateBannerOrder(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to update banner order');
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: bannerKeys.bannerById(id) });
      queryClient.invalidateQueries({ queryKey: bannerKeys.banners() });
    },
  });
};