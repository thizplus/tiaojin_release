import { useCallback, useMemo } from "react";
import { useSearchParams } from "react-router-dom";
import { useSales } from "./useSales";
import type { Sale, SaleFilter } from "@/types";

interface UseSalesDataReturn {
  sales: Sale[];
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
  currentSearch: string;
  handleSearchChange: (query: string) => void;
  handlePageChange: (page: number) => void;
  handlePageSizeChange: (pageSize: number) => void;
}

export const useSalesData = (): UseSalesDataReturn => {
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

  const salesParams: SaleFilter = useMemo(() => ({
    search: searchQuery || undefined,
    page: currentPage,
    limit: currentPageSize,
    sort: "created_at",
    order: "desc"
  }), [searchQuery, currentPage, currentPageSize]);

  const {
    salesList: sales,
    pagination,
    isLoading,
  } = useSales(salesParams, undefined, { enableList: true });

  const handleSearchChange = useCallback((query: string) => {
    updateParams({ search: query || undefined, page: 1 });
  }, [updateParams]);

  const handlePageChange = useCallback((page: number) => {
    updateParams({ page: page + 1 }); // Convert 0-based to 1-based
  }, [updateParams]);

  const handlePageSizeChange = useCallback((pageSize: number) => {
    updateParams({ limit: pageSize, page: 1 });
  }, [updateParams]);

  return {
    sales,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearch: searchQuery,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
  };
};