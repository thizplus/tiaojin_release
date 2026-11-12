import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult
} from '@tanstack/react-query';
import { toursRegularService } from '@/services/tours-regular';
import type {
  Tour,
  ToursRegularResponse,
  ToursRegularRequest,
  PaginationMeta
} from '@/types/tours-regular';

// Query Keys
const toursRegularKeys = {
  all: ['tours-regular'] as const,
  lists: () => [...toursRegularKeys.all, 'list'] as const,
  list: (params?: ToursRegularRequest) => [...toursRegularKeys.lists(), params] as const,
};

// Zustand Store Interface
interface ToursRegularState {
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
export const useToursRegularStore = create<ToursRegularState>()(
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
      name: 'tours-regular-storage',
      partialize: (state) => ({
        selectedTour: state.selectedTour,
      }),
    }
  )
);

// React Query Hook
export const useToursRegularQuery = (
  params: ToursRegularRequest = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<ToursRegularResponse, Error> => {
  const { setToursList, setPaginationMeta } = useToursRegularStore();

  return useQuery({
    queryKey: toursRegularKeys.list(params),
    queryFn: async () => {
      const response = await toursRegularService.getRegularTours(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch tours list");
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
export const useRefreshToursRegular = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await queryClient.invalidateQueries({
        queryKey: toursRegularKeys.lists(),
      });
    },
  });
};

export { toursRegularKeys };