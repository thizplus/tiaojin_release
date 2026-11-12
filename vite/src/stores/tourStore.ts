import { create } from "zustand";
import { persist } from "zustand/middleware";
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from "@tanstack/react-query";
import tourService from "@/services/admin/tourService";
import type {
  Tour,
  TourFilter,
  CreateTourRequest,
  UpdateTourRequest,
  UpdateTourStatusRequest,
  TourPeriod,
  TourPeriodsResponse,
  TourPeriodsFilter,
  ApiResponse,
  PaginationMeta
} from "@/types";

const tourKeys = {
  tours: (params?: TourFilter) => ["tour", "list", JSON.stringify(params || {})] as const,
  tourById: (id: string) => ["tour", "byId", id] as const,
  tourPeriods: (tourId: string, params?: TourPeriodsFilter) => ["tour", "periods", tourId, JSON.stringify(params || {})] as const,
};

interface TourState {
  selectedTour: Tour | null;
  toursList: Tour[];
  tourPeriods: TourPeriod[];
  setSelectedTour: (tour: Tour | null) => void;
  setToursList: (tours: Tour[]) => void;
  setTourPeriods: (periods: TourPeriod[]) => void;
  clearTourData: () => void;
}

export const useTourStore = create<TourState>()(
  persist(
    (set) => ({
      selectedTour: null,
      toursList: [],
      tourPeriods: [],
      setSelectedTour: (tour) => set({ selectedTour: tour }),
      setToursList: (tours) => set({ toursList: tours }),
      setTourPeriods: (periods) => set({ tourPeriods: periods }),
      clearTourData: () => set({
        selectedTour: null,
        toursList: [],
        tourPeriods: [],
      }),
    }),
    {
      name: "tour-storage",
      partialize: (state) => ({
        selectedTour: state.selectedTour,
      }),
    }
  )
);

export const useToursListQuery = (
  params: TourFilter = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<{ tours: Tour[], pagination: PaginationMeta }>, Error> => {
  const { setToursList } = useTourStore();

  return useQuery({
    queryKey: tourKeys.tours(params),
    queryFn: async () => {
      const response = await tourService.getTours(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch tours list");
      }

      const tours = response.data.tours || [];
      setToursList(tours);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
    staleTime: 0, // Always fresh
    gcTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    refetchInterval: false,
  });
};

export const useTourByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<Tour, Error> => {
  const { setSelectedTour } = useTourStore();

  return useQuery({
    queryKey: tourKeys.tourById(id),
    queryFn: async () => {
      const response = await tourService.getTourById(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch tour data");
      }

      setSelectedTour(response.data);
      return response.data;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

export const useCreateTourMutation = (): UseMutationResult<
  Tour,
  Error,
  CreateTourRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateTourRequest) => {
      const response = await tourService.createTour(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to create tour");
      }

      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tour", "list"] });
    },
  });
};

export const useUpdateTourMutation = (): UseMutationResult<
  Tour,
  Error,
  { id: string; data: UpdateTourRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await tourService.updateTour(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update tour");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["tour", "byId", id] });
      queryClient.invalidateQueries({ queryKey: ["tour", "list"] });
    },
  });
};

export const useDeleteTourMutation = (): UseMutationResult<
  null,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await tourService.deleteTour(id);

      if (!response.success) {
        throw new Error(response.message || "Failed to delete tour");
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["tour", "byId", id] });
      queryClient.invalidateQueries({ queryKey: ["tour", "list"] });
    },
  });
};

export const useUpdateTourStatusMutation = (): UseMutationResult<
  Tour,
  Error,
  { id: string; data: UpdateTourStatusRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await tourService.updateTourStatus(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update tour status");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["tour", "byId", id] });
      queryClient.invalidateQueries({ queryKey: ["tour", "list"] });
    },
  });
};

export const useTourPeriodsQuery = (
  tourId: string,
  params: TourPeriodsFilter = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<TourPeriodsResponse>, Error> => {
  const { setTourPeriods } = useTourStore();

  return useQuery({
    queryKey: tourKeys.tourPeriods(tourId, params),
    queryFn: async () => {
      const response = await tourService.getTourPeriods(tourId, params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch tour periods");
      }

      const periods = response.data.periods || [];
      setTourPeriods(periods);
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!tourId,
    staleTime: 0,
    gcTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    refetchInterval: false,
  });
};