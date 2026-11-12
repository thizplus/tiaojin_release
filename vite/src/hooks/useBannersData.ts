import { useCallback, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useBanner } from './useBanner';
import { useDebounce } from './useDebounce';
import type { BannerFilter, BannerPosition } from '@/types';

interface UseBannersDataProps {
  initialPageSize?: number;
  position?: BannerPosition;
}

export function useBannersData({ initialPageSize = 25, position }: UseBannersDataProps = {}) {
  const [searchParams, setSearchParams] = useSearchParams();

  // Read from URL params
  const page = useMemo(() => {
    const p = searchParams.get("page");
    const parsed = p ? parseInt(p, 10) : 1;
    return isNaN(parsed) || parsed < 1 ? 1 : parsed;
  }, [searchParams]);

  const pageSize = useMemo(() => {
    const limit = searchParams.get("limit");
    const parsed = limit ? parseInt(limit, 10) : initialPageSize;
    return isNaN(parsed) || parsed < 1 ? initialPageSize : parsed;
  }, [searchParams, initialPageSize]);

  const searchQuery = useMemo(() => {
    return searchParams.get("search") || "";
  }, [searchParams]);

  // Helper to update URL params
  const updateParams = useCallback(
    (updates: Record<string, string | number | undefined>) => {
      setSearchParams((prev) => {
        const newParams = new URLSearchParams(prev);
        Object.entries(updates).forEach(([key, value]) => {
          if (value === undefined || value === "" || (key === "page" && value === 1)) {
            newParams.delete(key);
          } else {
            newParams.set(key, String(value));
          }
        });
        return newParams;
      }, { replace: true });
    },
    [setSearchParams]
  );

  // ใช้ debounce สำหรับ search query
  const debouncedSearchQuery = useDebounce(searchQuery, 500);

  // Memoize filter params เพื่อป้องกัน unnecessary rerenders
  const filterParams: BannerFilter = useMemo(() => ({
    search: debouncedSearchQuery,
    position,
    page,
    limit: pageSize,
    sort: 'display_order',
    order: 'asc'
  }), [debouncedSearchQuery, position, page, pageSize]);

  // ใช้ useBanner hook
  const {
    bannersList,
    pagination,
    isLoading,
    createBanner,
    updateBanner,
    deleteBanner,
    updateBannerStatus,
    updateBannerOrder,
    refetchBannersList,
    isCreateLoading,
    isUpdateLoading,
    isDeleteLoading,
    isUpdateOrderLoading
  } = useBanner(filterParams, undefined, { enableList: true });

  // Handlers
  const handleSearchChange = useCallback((query: string) => {
    updateParams({ search: query || undefined, page: 1 });
  }, [updateParams]);

  const handlePageChange = useCallback((newPage: number) => {
    updateParams({ page: newPage + 1 }); // DataTable uses 0-based, API uses 1-based
  }, [updateParams]);

  const handlePageSizeChange = useCallback((newPageSize: number) => {
    updateParams({ limit: newPageSize, page: 1 });
  }, [updateParams]);

  return {
    // Data
    banners: bannersList,
    pagination,
    isLoading,

    // Current state
    currentPage: page,
    currentPageSize: pageSize,
    currentSearchQuery: searchQuery,

    // Actions
    createBanner,
    updateBanner,
    deleteBanner,
    updateBannerStatus,
    updateBannerOrder,
    refetchBannersList,

    // Handlers
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,

    // Loading states
    isCreateLoading,
    isUpdateLoading,
    isDeleteLoading,
    isUpdateOrderLoading,
  };
}