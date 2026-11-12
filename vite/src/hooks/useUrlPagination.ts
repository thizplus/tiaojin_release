import { useCallback, useMemo } from "react";
import { useSearchParams } from "react-router-dom";

interface UseUrlPaginationOptions {
  defaultPage?: number;
  defaultPageSize?: number;
  defaultSearch?: string;
}

interface UseUrlPaginationReturn {
  page: number;
  pageSize: number;
  search: string;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  setSearch: (search: string) => void;
  resetToFirstPage: () => void;
}

/**
 * Hook สำหรับ sync pagination state กับ URL search params
 * เมื่อ refresh หน้า state จะไม่หาย
 */
export const useUrlPagination = (
  options: UseUrlPaginationOptions = {}
): UseUrlPaginationReturn => {
  const {
    defaultPage = 1,
    defaultPageSize = 25,
    defaultSearch = "",
  } = options;

  const [searchParams, setSearchParams] = useSearchParams();

  // Read from URL
  const page = useMemo(() => {
    const pageParam = searchParams.get("page");
    const parsed = pageParam ? parseInt(pageParam, 10) : defaultPage;
    return isNaN(parsed) || parsed < 1 ? defaultPage : parsed;
  }, [searchParams, defaultPage]);

  const pageSize = useMemo(() => {
    const sizeParam = searchParams.get("limit");
    const parsed = sizeParam ? parseInt(sizeParam, 10) : defaultPageSize;
    return isNaN(parsed) || parsed < 1 ? defaultPageSize : parsed;
  }, [searchParams, defaultPageSize]);

  const search = useMemo(() => {
    return searchParams.get("search") || defaultSearch;
  }, [searchParams, defaultSearch]);

  // Update URL helper
  const updateParams = useCallback(
    (updates: Record<string, string | number | undefined>) => {
      setSearchParams((prev) => {
        const newParams = new URLSearchParams(prev);
        Object.entries(updates).forEach(([key, value]) => {
          if (value === undefined || value === "" || value === 1) {
            // Remove default values from URL to keep it clean
            if (key === "page" && value === 1) {
              newParams.delete(key);
            } else if (value === undefined || value === "") {
              newParams.delete(key);
            } else {
              newParams.set(key, String(value));
            }
          } else {
            newParams.set(key, String(value));
          }
        });
        return newParams;
      }, { replace: true }); // Use replace to avoid creating history entries
    },
    [setSearchParams]
  );

  const setPage = useCallback(
    (newPage: number) => {
      updateParams({ page: newPage });
    },
    [updateParams]
  );

  const setPageSize = useCallback(
    (newPageSize: number) => {
      // Reset to page 1 when changing page size
      updateParams({ limit: newPageSize, page: 1 });
    },
    [updateParams]
  );

  const setSearch = useCallback(
    (newSearch: string) => {
      // Reset to page 1 when searching
      updateParams({ search: newSearch || undefined, page: 1 });
    },
    [updateParams]
  );

  const resetToFirstPage = useCallback(() => {
    updateParams({ page: 1 });
  }, [updateParams]);

  return {
    page,
    pageSize,
    search,
    setPage,
    setPageSize,
    setSearch,
    resetToFirstPage,
  };
};
