import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import { periodService } from '@/services/admin/periodService';
import type { TourPeriod, PeriodStatus } from '@/types/models/period';
import type { PeriodWithTourDetails } from '@/types/response/period';
import type {
  PeriodListRequest,
  CreatePeriodRequest,
  UpdatePeriodRequest,
} from '@/types/request/period';
import type { ApiResponse, PaginationMeta } from '@/types';

const periodKeys = {
  periods: (params?: PeriodListRequest) => ['period', 'list', JSON.stringify(params || {})] as const,
  periodById: (id: string) => ['period', 'byId', id] as const,
};

interface PeriodState {
  currentPeriod: TourPeriod | null;
  selectedPeriod: TourPeriod | null;
  periodsList: PeriodWithTourDetails[];
  setCurrentPeriod: (period: TourPeriod | null) => void;
  setSelectedPeriod: (period: TourPeriod | null) => void;
  setPeriodsList: (periods: PeriodWithTourDetails[]) => void;
  clearPeriodData: () => void;
}

export const usePeriodStore = create<PeriodState>()(
  persist(
    (set) => ({
      currentPeriod: null,
      selectedPeriod: null,
      periodsList: [],
      setCurrentPeriod: (period) => set({ currentPeriod: period }),
      setSelectedPeriod: (period) => set({ selectedPeriod: period }),
      setPeriodsList: (periods) => set({ periodsList: periods }),
      clearPeriodData: () => set({
        currentPeriod: null,
        selectedPeriod: null,
        periodsList: [],
      }),
    }),
    {
      name: 'period-storage',
      partialize: (state) => ({
        selectedPeriod: state.selectedPeriod,
      }),
    }
  )
);

export const usePeriodsListQuery = (
  params: PeriodListRequest = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<{ periods: PeriodWithTourDetails[], pagination: PaginationMeta }>, Error> => {
  const { setPeriodsList } = usePeriodStore();

  return useQuery({
    queryKey: periodKeys.periods(params),
    queryFn: async () => {
      const response = await periodService.getAll(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch periods list");
      }

      const periods = response.data.periods || [];
      setPeriodsList(periods);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
    staleTime: 0,
    gcTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    refetchInterval: false,
  });
};

export const usePeriodByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<TourPeriod, Error> => {
  const { setCurrentPeriod } = usePeriodStore();

  return useQuery({
    queryKey: periodKeys.periodById(id),
    queryFn: async () => {
      const response = await periodService.getById(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch period data");
      }

      setCurrentPeriod(response.data);
      return response.data;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

export const useCreatePeriodMutation = (): UseMutationResult<
  TourPeriod,
  Error,
  CreatePeriodRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreatePeriodRequest) => {
      const response = await periodService.create(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to create period");
      }

      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['period', 'list'] });
    },
  });
};

export const useUpdatePeriodMutation = (): UseMutationResult<
  TourPeriod,
  Error,
  { id: string; data: UpdatePeriodRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await periodService.update(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update period");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['period', 'byId', id] });
      queryClient.invalidateQueries({ queryKey: ['period', 'list'] });
    },
  });
};

export const useDeletePeriodMutation = (): UseMutationResult<
  null,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await periodService.delete(id);

      if (!response.success) {
        throw new Error(response.message || "Failed to delete period");
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['period', 'byId', id] });
      queryClient.invalidateQueries({ queryKey: ['period', 'list'] });
    },
  });
};

export const useUpdatePeriodStatusMutation = (): UseMutationResult<
  { id: string; status: string; updated_at: string },
  Error,
  { id: string; status: PeriodStatus },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, status }) => {
      const response = await periodService.updateStatus(id, { status });

      if (!response.success) {
        throw new Error(response.message || "Failed to update period status");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['period', 'byId', id] });
      queryClient.invalidateQueries({ queryKey: ['period', 'list'] });
    },
  });
};

export const useCheckPeriodAvailabilityMutation = (): UseMutationResult<
  {
    available: boolean;
    requested_seats: number;
    available_seats: number;
    total_seats: number;
    booked_seats: number;
  },
  Error,
  { id: string; requestedSeats?: number },
  unknown
> => {
  return useMutation({
    mutationFn: async ({ id, requestedSeats }) => {
      const response = await periodService.checkAvailability(id, { requested_seats: requestedSeats });

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to check availability");
      }

      return response.data;
    },
  });
};