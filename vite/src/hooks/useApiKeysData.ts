import { useCallback, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useApiKey } from './useApiKey';
import { useDebounce } from './useDebounce';
import type { ApiKeyListRequest } from '@/types/request/apiKey';

interface UseApiKeysDataProps {
  initialPageSize?: number;
}

export function useApiKeysData({ initialPageSize = 25 }: UseApiKeysDataProps = {}) {
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
  const filterParams: ApiKeyListRequest = useMemo(() => ({
    search: debouncedSearchQuery,
    page,
    limit: pageSize,
  }), [debouncedSearchQuery, page, pageSize]);

  // ใช้ useApiKey hook
  const {
    apiKeys,
    isLoading,
    createApiKey,
    updateApiKey,
    deleteApiKey,
    refetchApiKeysList,
    isCreateLoading,
    isUpdateLoading,
    isDeleteLoading
  } = useApiKey(filterParams, undefined, { enableList: true });

  // Calculate pagination from store data
  const pagination = useMemo(() => ({
    total: apiKeys.length, // This should come from API response
    totalPages: Math.ceil(apiKeys.length / pageSize),
    currentPage: page,
  }), [apiKeys.length, page, pageSize]);

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
    apiKeys,
    pagination,
    isLoading,

    // Current state
    currentPage: page,
    currentPageSize: pageSize,
    currentSearchQuery: searchQuery,

    // Actions
    createApiKey,
    updateApiKey,
    deleteApiKey,
    refetchApiKeysList,

    // Handlers
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,

    // Loading states
    isCreateLoading,
    isUpdateLoading,
    isDeleteLoading,
  };
}