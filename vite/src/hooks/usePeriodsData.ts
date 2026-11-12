import { useState, useCallback, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { usePeriodsListQuery } from '@/stores/periodStore';
import type { PeriodListRequest } from '@/types/request/period';
import type { PeriodWithTourDetails } from '@/types/response/period';

interface UsePeriodsDataOptions {
  tourId?: string;
  autoFetch?: boolean;
  refetchInterval?: number;
}

interface UsePeriodsDataReturn {
  periods: PeriodWithTourDetails[];
  pagination: {
    total: number;
    total_pages: number;
    current_page: number;
    page_size: number;
    has_more: boolean;
    limit: number;
    offset: number;
  } | null;
  loading: boolean;
  currentPage: number;
  currentPageSize: number;
  totalPages: number;
  params: PeriodListRequest;
  filters: PeriodListRequest;
  goToPage: (page: number) => void;
  changePageSize: (pageSize: number) => void;
  refresh: () => void;
  updateParams: (newParams: Partial<PeriodListRequest>) => void;
  setFilters: (filters: Partial<PeriodListRequest>) => void;
}

export const usePeriodsData = (
  initialParams: PeriodListRequest = {},
  options: UsePeriodsDataOptions = {}
): UsePeriodsDataReturn => {
  const { tourId, autoFetch = true } = options;
  const [searchParams, setSearchParams] = useSearchParams();
  const [filters, setFiltersState] = useState<PeriodListRequest>({});

  // Read from URL params
  const currentPage = useMemo(() => {
    const page = searchParams.get("page");
    const parsed = page ? parseInt(page, 10) : 1;
    return isNaN(parsed) || parsed < 1 ? 1 : parsed;
  }, [searchParams]);

  const currentPageSize = useMemo(() => {
    const limit = searchParams.get("limit");
    const parsed = limit ? parseInt(limit, 10) : (initialParams.limit || 25);
    return isNaN(parsed) || parsed < 1 ? 25 : parsed;
  }, [searchParams, initialParams.limit]);

  // Helper to update URL params
  const updateUrlParams = useCallback(
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

  const params: PeriodListRequest = useMemo(() => {
    const baseParams = {
      ...initialParams,
      ...filters,
      limit: currentPageSize,
      page: currentPage,
      ...(tourId && { tour_id: tourId }),
    };
    return baseParams;
  }, [initialParams, filters, currentPage, currentPageSize, tourId]);

  const {
    data,
    isLoading: loading,
    refetch,
  } = usePeriodsListQuery(params, {
    enabled: autoFetch,
  });

  const periods = data?.data?.periods || [];
  const pagination = data?.data?.pagination || {
    total: 0,
    has_more: false,
    current_page: 1,
    total_pages: 1,
    page_size: currentPageSize,
    limit: currentPageSize,
    offset: 0,
  };

  const goToPage = useCallback((page: number) => {
    updateUrlParams({ page }); // Direct page number (1-based)
  }, [updateUrlParams]);

  const changePageSize = useCallback((pageSize: number) => {
    updateUrlParams({ limit: pageSize, page: 1 });
  }, [updateUrlParams]);

  const updateParams = useCallback((newParams: Partial<PeriodListRequest>) => {
    const updates: Record<string, string | number | undefined> = {};
    if (newParams.limit) {
      updates.limit = newParams.limit;
    }
    if (newParams.page !== undefined) {
      updates.page = newParams.page;
    }
    if (Object.keys(updates).length > 0) {
      updateUrlParams(updates);
    }
  }, [updateUrlParams]);

  const refresh = useCallback(() => {
    refetch();
  }, [refetch]);

  const setFilters = useCallback((newFilters: Partial<PeriodListRequest>) => {
    setFiltersState(newFilters);
    updateUrlParams({ page: 1 }); // Reset to first page when filters change
  }, [updateUrlParams]);

  return {
    periods,
    pagination,
    loading,
    currentPage,
    currentPageSize,
    totalPages: pagination.total_pages,
    params,
    filters,
    goToPage,
    changePageSize,
    refresh,
    updateParams,
    setFilters,
  };
};