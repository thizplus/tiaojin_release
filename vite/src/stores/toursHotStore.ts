import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult
} from '@tanstack/react-query';
import { toursHotService } from '@/services/tours-hot';
import type {
  Tour,
  ToursRegularResponse,
  ToursRegularRequest,
  PaginationMeta
} from '@/types/tours-regular';

// Query Keys
const toursHotKeys = {
  all: ['tours-hot'] as const,
  lists: () => [...toursHotKeys.all, 'list'] as const,
  list: (params?: ToursRegularRequest) => [...toursHotKeys.lists(), params] as const,
};

// Zustand Store Interface
interface ToursHotState {
  currentTour: Tour | null;
  selectedTour: Tour | null;
  toursList: Tour[];
  paginationMeta: PaginationMeta | null;
  setCurrentTour: (tour: Tour | null) => void;
  setSelectedTour: (tour: Tour | null) => void;
  setToursList: (tours: Tour[]) => void;
  setPaginationMeta: (meta: PaginationMeta | null) => void;
  clearToursData: () => void;
}

// Zustand Store
export const useToursHotStore = create<ToursHotState>()(
  persist(
    (set) => ({
      currentTour: null,
      selectedTour: null,
      toursList: [],
      paginationMeta: null,
      setCurrentTour: (tour) => set({ currentTour: tour }),
      setSelectedTour: (tour) => set({ selectedTour: tour }),
      setToursList: (tours) => set({ toursList: tours }),
      setPaginationMeta: (meta) => set({ paginationMeta: meta }),
      clearToursData: () => set({
        currentTour: null,
        selectedTour: null,
        toursList: [],
        paginationMeta: null,
      }),
    }),
    {
      name: 'tours-hot-storage',
      partialize: (state) => ({
        selectedTour: state.selectedTour,
      }),
    }
  )
);

// React Query Hook
export const useToursHotQuery = (
  params: ToursRegularRequest = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<ToursRegularResponse, Error> => {
  const { setToursList, setPaginationMeta } = useToursHotStore();

  return useQuery({
    queryKey: toursHotKeys.list(params),
    queryFn: async () => {
      const response = await toursHotService.getHotTours(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch hot tours list");
      }

      const tours = response.data.tours || [];
      const pagination = response.data.pagination || null;

      setToursList(tours);
      setPaginationMeta(pagination);

      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
    staleTime: 1000 * 60, // 1 minute
    gcTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    refetchInterval: false,
  });
};

// Hook สำหรับ refresh tours data
export const useRefreshToursHot = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await queryClient.invalidateQueries({
        queryKey: toursHotKeys.lists(),
      });
    },
  });
};

export { toursHotKeys };