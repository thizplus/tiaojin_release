import { useCallback, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useCompanyProfile } from './useCompanyProfile';
import { useDebounce } from './useDebounce';
import type { GetCompanyProfilesRequest } from '@/types';

interface UseCompanyProfilesDataProps {
  initialPageSize?: number;
}

export function useCompanyProfilesData({ initialPageSize = 25 }: UseCompanyProfilesDataProps = {}) {
  const [searchParams, setSearchParams] = useSearchParams();

  // Read from URL params
  const page = useMemo(() => {
    const p = searchParams.get("page");
    const parsed = p ? parseInt(p, 10) : 1;
    return isNaN(parsed) || parsed < 1 ? 1 : parsed;
  }, [searchParams]);

  const limit = useMemo(() => {
    const l = searchParams.get("limit");
    const parsed = l ? parseInt(l, 10) : initialPageSize;
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
  const filterParams: GetCompanyProfilesRequest = useMemo(() => ({
    search: debouncedSearchQuery || undefined,
    page,
    limit,
  }), [debouncedSearchQuery, page, limit]);

  // ใช้ useCompanyProfile hook
  const {
    companyProfilesList,
    pagination,
    isLoading,
    approveCompanyProfile,
    refetchCompanyProfilesList,
    isApproveLoading,
    isCompanyProfilesListLoading
  } = useCompanyProfile(filterParams, undefined, undefined, undefined, { enableList: true });

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

  // Calculate total pages from pagination
  const totalPages = useMemo(() => {
    if (!pagination) return 1;
    return Math.ceil(pagination.total / pagination.limit);
  }, [pagination]);

  return {
    // Data
    companyProfiles: companyProfilesList,
    pagination,
    isLoading,

    // Current state
    currentPage: page,
    currentPageSize: limit,
    currentSearchQuery: searchQuery,

    // Actions
    approveCompanyProfile,
    refetchCompanyProfilesList,

    // Handlers
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,

    // Loading states
    isApproveLoading,
    isCompanyProfilesListLoading,

    // Calculated values
    totalPages,
  };
}