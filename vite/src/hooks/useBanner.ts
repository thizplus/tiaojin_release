import { useCallback } from 'react';
import * as React from 'react';
import {
  useBannerStore,
  useBannersListQuery,
  useBannerByIdQuery,
  useCreateBannerMutation,
  useUpdateBannerMutation,
  useDeleteBannerMutation,
  useUpdateBannerStatusMutation,
  useUpdateBannerOrderMutation
} from '@/stores/bannerStore';
import type {
  Banner,
  BannerFilter,
  CreateBannerRequest,
  UpdateBannerRequest,
  UpdateBannerStatusRequest,
  UpdateBannerOrderRequest,

} from '@/types';

interface UseBannerReturn {
  // State
  currentBanner: Banner | null;
  selectedBanner: Banner | null;
  bannersList: Banner[];
  isLoading: boolean;

  // Pagination data
  pagination: {
    total: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
    hasMore: boolean;
  } | null;

  // Actions
  createBanner: (data: CreateBannerRequest) => Promise<Banner>;
  updateBanner: (id: string, data: UpdateBannerRequest) => Promise<Banner>;
  deleteBanner: (id: string) => Promise<null>;
  updateBannerStatus: (id: string, data: UpdateBannerStatusRequest) => Promise<Banner>;
  updateBannerOrder: (id: string, data: UpdateBannerOrderRequest) => Promise<Banner>;

  // State management
  setCurrentBanner: (banner: Banner | null) => void;
  setSelectedBanner: (banner: Banner | null) => void;
  setBannersList: (banners: Banner[]) => void;
  clearBannerData: () => void;

  // Queries
  getBannersList: () => void;
  getBannerById: () => void;
  refetchBannersList: () => void;
  refetchBannerById: () => void;

  // Loading states
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  isDeleteLoading: boolean;
  isUpdateStatusLoading: boolean;
  isUpdateOrderLoading: boolean;
  isBannersListLoading: boolean;
  isBannerByIdLoading: boolean;
}

export const useBanner = (
  listParams?: BannerFilter,
  bannerId?: string,
  options: { enableList?: boolean; enableById?: boolean } = {}
): UseBannerReturn => {
  // Store state
  const {
    currentBanner,
    selectedBanner,
    bannersList,
    setCurrentBanner,
    setSelectedBanner,
    setBannersList,
    clearBannerData
  } = useBannerStore();

  // Mutations
  const createBannerMutation = useCreateBannerMutation();
  const updateBannerMutation = useUpdateBannerMutation();
  const deleteBannerMutation = useDeleteBannerMutation();
  const updateBannerStatusMutation = useUpdateBannerStatusMutation();
  const updateBannerOrderMutation = useUpdateBannerOrderMutation();

  // Queries
  const {
    data: bannersData,
    refetch: refetchBannersList,
    isLoading: isBannersListLoading
  } = useBannersListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchBannerById,
    isLoading: isBannerByIdLoading
  } = useBannerByIdQuery(bannerId || '', {
    enabled: options.enableById !== false && !!bannerId
  });

  // Actions
  const createBanner = useCallback(async (data: CreateBannerRequest) => {
    const result = await createBannerMutation.mutateAsync(data);
    return result;
  }, [createBannerMutation]);

  const updateBanner = useCallback(async (id: string, data: UpdateBannerRequest) => {
    const result = await updateBannerMutation.mutateAsync({ id, data });
    return result;
  }, [updateBannerMutation]);

  const deleteBanner = useCallback(async (id: string) => {
    const result = await deleteBannerMutation.mutateAsync(id);
    return result;
  }, [deleteBannerMutation]);

  const updateBannerStatus = useCallback(async (id: string, data: UpdateBannerStatusRequest) => {
    const result = await updateBannerStatusMutation.mutateAsync({ id, data });
    return result;
  }, [updateBannerStatusMutation]);

  const updateBannerOrder = useCallback(async (id: string, data: UpdateBannerOrderRequest) => {
    const result = await updateBannerOrderMutation.mutateAsync({ id, data });
    return result;
  }, [updateBannerOrderMutation]);

  // Query management
  const getBannersList = useCallback(() => {
    refetchBannersList();
  }, [refetchBannersList]);

  const getBannerById = useCallback(() => {
    refetchBannerById();
  }, [refetchBannerById]);

  // Process banners data and pagination
  const pagination = React.useMemo(() => {
    if (!bannersData?.data?.pagination) return null;

    const paginationMeta = bannersData.data.pagination;
    return {
      total: paginationMeta.total,
      totalPages: paginationMeta.total_pages,
      currentPage: paginationMeta.current_page,
      pageSize: paginationMeta.page_size,
      hasMore: paginationMeta.has_more,
    };
  }, [bannersData]);

  // Calculate overall loading state
  const isLoading = createBannerMutation.isPending ||
                   updateBannerMutation.isPending ||
                   deleteBannerMutation.isPending ||
                   updateBannerStatusMutation.isPending ||
                   updateBannerOrderMutation.isPending ||
                   isBannersListLoading ||
                   isBannerByIdLoading;

  return {
    // State
    currentBanner,
    selectedBanner,
    bannersList,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    createBanner,
    updateBanner,
    deleteBanner,
    updateBannerStatus,
    updateBannerOrder,

    // State management
    setCurrentBanner,
    setSelectedBanner,
    setBannersList,
    clearBannerData,

    // Queries
    getBannersList,
    getBannerById,
    refetchBannersList,
    refetchBannerById,

    // Loading states
    isCreateLoading: createBannerMutation.isPending,
    isUpdateLoading: updateBannerMutation.isPending,
    isDeleteLoading: deleteBannerMutation.isPending,
    isUpdateStatusLoading: updateBannerStatusMutation.isPending,
    isUpdateOrderLoading: updateBannerOrderMutation.isPending,
    isBannersListLoading,
    isBannerByIdLoading,
  };
};

export default useBanner;