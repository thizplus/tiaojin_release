import { useCallback } from "react";
import * as React from "react";
import {
  useSalesStore,
  useSalesListQuery,
  useSaleByIdQuery,
  useCreateSaleMutation,
  useUpdateSaleMutation,
  useDeleteSaleMutation,
  useActivateSaleMutation,
  useDeactivateSaleMutation,
  useSalesActiveQuery
} from "@/stores/salesStore";
import type {
  Sale,
  SaleFilter,
  CreateSaleRequest,
  UpdateSaleRequest,
 // PaginatedResponse
} from "@/types";

interface UseSalesReturn {
  // State
  selectedSale: Sale | null;
  salesList: Sale[];
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
  createSale: (data: CreateSaleRequest) => Promise<Sale>;
  updateSale: (id: string, data: UpdateSaleRequest) => Promise<Sale>;
  deleteSale: (id: string) => Promise<null>;
  activateSale: (id: string) => Promise<Sale>;
  deactivateSale: (id: string) => Promise<Sale>;

  // State management
  setSelectedSale: (sale: Sale | null) => void;
  setSalesList: (sales: Sale[]) => void;
  clearSalesData: () => void;

  // Queries
  getSalesList: (params?: SaleFilter, options?: { enabled?: boolean }) => void;
  getSaleById: (id: string, options?: { enabled?: boolean }) => void;
  getSalesActive: (options?: { enabled?: boolean }) => void;
  refetchSalesList: () => void;
  refetchSaleById: () => void;
  refetchSalesActive: () => void;

  // Loading states
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  isDeleteLoading: boolean;
  isActivateLoading: boolean;
  isDeactivateLoading: boolean;
  isSalesListLoading: boolean;
  isSaleByIdLoading: boolean;
  isSalesActiveLoading: boolean;
}

export const useSales = (
  listParams?: SaleFilter,
  saleId?: string,
  options: { enableList?: boolean; enableById?: boolean; enableActive?: boolean } = {}
): UseSalesReturn => {
  // Store state
  const {
    selectedSale,
    salesList,
    setSelectedSale,
    setSalesList,
    clearSalesData
  } = useSalesStore();

  // Mutations
  const createSaleMutation = useCreateSaleMutation();
  const updateSaleMutation = useUpdateSaleMutation();
  const deleteSaleMutation = useDeleteSaleMutation();
  const activateSaleMutation = useActivateSaleMutation();
  const deactivateSaleMutation = useDeactivateSaleMutation();

  // Queries
  const {
    data: salesData,
    refetch: refetchSalesList,
    isLoading: isSalesListLoading
  } = useSalesListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchSaleById,
    isLoading: isSaleByIdLoading
  } = useSaleByIdQuery(saleId || "", {
    enabled: options.enableById !== false && !!saleId
  });

  const {
    refetch: refetchSalesActive,
    isLoading: isSalesActiveLoading
  } = useSalesActiveQuery({
    enabled: options.enableActive === true
  });

  // Actions
  const createSale = useCallback(async (data: CreateSaleRequest) => {
    const result = await createSaleMutation.mutateAsync(data);
    return result;
  }, [createSaleMutation]);

  const updateSale = useCallback(async (id: string, data: UpdateSaleRequest) => {
    const result = await updateSaleMutation.mutateAsync({ id, data });
    return result;
  }, [updateSaleMutation]);

  const deleteSale = useCallback(async (id: string) => {
    const result = await deleteSaleMutation.mutateAsync(id);
    return result;
  }, [deleteSaleMutation]);

  const activateSale = useCallback(async (id: string) => {
    const result = await activateSaleMutation.mutateAsync(id);
    return result;
  }, [activateSaleMutation]);

  const deactivateSale = useCallback(async (id: string) => {
    const result = await deactivateSaleMutation.mutateAsync(id);
    return result;
  }, [deactivateSaleMutation]);

  // Query management
  const getSalesList = useCallback(() => {
    refetchSalesList();
  }, [refetchSalesList]);

  const getSaleById = useCallback(() => {
    refetchSaleById();
  }, [refetchSaleById]);

  const getSalesActive = useCallback(() => {
    refetchSalesActive();
  }, [refetchSalesActive]);

  // Process sales data and pagination
  const pagination = React.useMemo(() => {
    if (!salesData?.data?.pagination) return null;

    const paginationMeta = salesData.data.pagination;
    return {
      total: paginationMeta.total,
      totalPages: paginationMeta.total_pages,
      currentPage: paginationMeta.current_page,
      pageSize: paginationMeta.page_size || paginationMeta.limit,
      hasMore: paginationMeta.has_more,
    };
  }, [salesData]);

  // Calculate overall loading state
  const isLoading = createSaleMutation.isPending ||
                   updateSaleMutation.isPending ||
                   deleteSaleMutation.isPending ||
                   activateSaleMutation.isPending ||
                   deactivateSaleMutation.isPending ||
                   isSalesListLoading ||
                   isSaleByIdLoading ||
                   isSalesActiveLoading;

  return {
    // State
    selectedSale,
    salesList,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    createSale,
    updateSale,
    deleteSale,
    activateSale,
    deactivateSale,

    // State management
    setSelectedSale,
    setSalesList,
    clearSalesData,

    // Queries
    getSalesList,
    getSaleById,
    getSalesActive,
    refetchSalesList,
    refetchSaleById,
    refetchSalesActive,

    // Loading states
    isCreateLoading: createSaleMutation.isPending,
    isUpdateLoading: updateSaleMutation.isPending,
    isDeleteLoading: deleteSaleMutation.isPending,
    isActivateLoading: activateSaleMutation.isPending,
    isDeactivateLoading: deactivateSaleMutation.isPending,
    isSalesListLoading,
    isSaleByIdLoading,
    isSalesActiveLoading,
  };
};

export default useSales;