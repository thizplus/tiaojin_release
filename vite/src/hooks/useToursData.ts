import { useCallback, useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import { useTours } from "./useTours";
import type { TourFilter, Tour, TourStatus } from "@/types";

interface UseToursDataReturn {
  tours: Tour[];
  pagination: {
    total: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
    hasMore: boolean;
  } | null;
  isLoading: boolean;
  currentPage: number;
  currentPageSize: number;
  handleSearchChange: (query: string) => void;
  handlePageChange: (page: number) => void;
  handlePageSizeChange: (pageSize: number) => void;
  handleCityFilter: (cityId: string) => void;
  handleStatusFilter: (status: string) => void;
  currentCityFilter: string;
  currentStatusFilter: string;
  currentSearch: string;
}

export const useToursData = (): UseToursDataReturn => {
  const [searchParams, setSearchParams] = useSearchParams();

  // Read from URL params
  const currentPage = useMemo(() => {
    const page = searchParams.get("page");
    const parsed = page ? parseInt(page, 10) : 1;
    return isNaN(parsed) || parsed < 1 ? 1 : parsed;
  }, [searchParams]);

  const currentPageSize = useMemo(() => {
    const limit = searchParams.get("limit");
    const parsed = limit ? parseInt(limit, 10) : 25;
    return isNaN(parsed) || parsed < 1 ? 25 : parsed;
  }, [searchParams]);

  const searchQuery = useMemo(() => {
    return searchParams.get("search") || "";
  }, [searchParams]);

  const currentCityFilter = useMemo(() => {
    return searchParams.get("city") || "";
  }, [searchParams]);

  const currentStatusFilter = useMemo(() => {
    return searchParams.get("status") || "";
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

  const toursParams: TourFilter = useMemo(() => {
    const params: TourFilter = {
      search: searchQuery.trim() || undefined,
      city_id: currentCityFilter || undefined,
      status: currentStatusFilter as TourStatus || undefined,
      page: currentPage,
      limit: currentPageSize,
      sort: "created_at",
      order: "desc" as const
    };
    return params;
  }, [searchQuery, currentCityFilter, currentStatusFilter, currentPage, currentPageSize]);

  const {
    toursList: tours,
    pagination,
    isLoading,
  } = useTours(toursParams, undefined, { enableList: true, enableById: false });

  const handleSearchChange = useCallback((query: string) => {
    updateParams({ search: query || undefined, page: 1 });
  }, [updateParams]);

  const handlePageChange = useCallback((page: number) => {
    updateParams({ page: page + 1 }); // Convert 0-based to 1-based
  }, [updateParams]);

  const handlePageSizeChange = useCallback((pageSize: number) => {
    updateParams({ limit: pageSize, page: 1 });
  }, [updateParams]);

  const handleCityFilter = useCallback((cityId: string) => {
    updateParams({ city: cityId || undefined, page: 1 });
  }, [updateParams]);

  const handleStatusFilter = useCallback((status: string) => {
    updateParams({ status: status || undefined, page: 1 });
  }, [updateParams]);

  return {
    tours,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentCityFilter,
    currentStatusFilter,
    currentSearch: searchQuery,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
    handleCityFilter,
    handleStatusFilter,
  };
};