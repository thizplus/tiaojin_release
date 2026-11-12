import { create } from "zustand";
import { persist } from "zustand/middleware";
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from "@tanstack/react-query";
import salesService from "@/services/admin/salesService";
import type {
  Sale,
  SaleFilter,
  CreateSaleRequest,
  UpdateSaleRequest,

  ApiResponse,
  PaginationMeta
} from "@/types";

const salesKeys = {
  sales: (params?: SaleFilter) => ["sale", "list", params] as const,
  saleById: (id: string) => ["sale", "byId", id] as const,
  salesActive: () => ["sale", "active"] as const,
};

interface SalesState {
  selectedSale: Sale | null;
  salesList: Sale[];
  setSelectedSale: (sale: Sale | null) => void;
  setSalesList: (sales: Sale[]) => void;
  clearSalesData: () => void;
}

export const useSalesStore = create<SalesState>()(
  persist(
    (set) => ({
      selectedSale: null,
      salesList: [],
      setSelectedSale: (sale) => set({ selectedSale: sale }),
      setSalesList: (sales) => set({ salesList: sales }),
      clearSalesData: () => set({
        selectedSale: null,
        salesList: [],
      }),
    }),
    {
      name: "sales-storage",
      partialize: (state) => ({
        selectedSale: state.selectedSale,
      }),
    }
  )
);

export const useSalesListQuery = (
  params: SaleFilter = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<{ sales: Sale[], pagination: PaginationMeta }>, Error> => {
  const { setSalesList } = useSalesStore();

  return useQuery({
    queryKey: salesKeys.sales(params),
    queryFn: async () => {
      const response = await salesService.getSales(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch sales list");
      }

      const sales = response.data.sales || [];
      setSalesList(sales);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

export const useSaleByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<Sale, Error> => {
  const { setSelectedSale } = useSalesStore();

  return useQuery({
    queryKey: salesKeys.saleById(id),
    queryFn: async () => {
      const response = await salesService.getSaleById(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch sale data");
      }

      setSelectedSale(response.data);
      return response.data;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

export const useCreateSaleMutation = (): UseMutationResult<
  Sale,
  Error,
  CreateSaleRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateSaleRequest) => {
      const response = await salesService.createSale(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to create sale");
      }

      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: salesKeys.sales() });
    },
  });
};

export const useUpdateSaleMutation = (): UseMutationResult<
  Sale,
  Error,
  { id: string; data: UpdateSaleRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await salesService.updateSale(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update sale");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: salesKeys.saleById(id) });
      queryClient.invalidateQueries({ queryKey: salesKeys.sales() });
    },
  });
};

export const useDeleteSaleMutation = (): UseMutationResult<
  null,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await salesService.deleteSale(id);

      if (!response.success) {
        throw new Error(response.message || "Failed to delete sale");
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: salesKeys.saleById(id) });
      queryClient.invalidateQueries({ queryKey: salesKeys.sales() });
    },
  });
};

export const useActivateSaleMutation = (): UseMutationResult<
  Sale,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await salesService.activateSale(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to activate sale");
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: salesKeys.saleById(id) });
      queryClient.invalidateQueries({ queryKey: salesKeys.sales() });
    },
  });
};

export const useDeactivateSaleMutation = (): UseMutationResult<
  Sale,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await salesService.deactivateSale(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to deactivate sale");
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: salesKeys.saleById(id) });
      queryClient.invalidateQueries({ queryKey: salesKeys.sales() });
    },
  });
};

export const useSalesActiveQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<Sale[], Error> => {
  const { setSalesList } = useSalesStore();

  return useQuery({
    queryKey: salesKeys.salesActive(),
    queryFn: async () => {
      const response = await salesService.getSalesActive();

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch active sales");
      }

      const sales = response.data.sales || [];
      setSalesList(sales);
      return sales;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};