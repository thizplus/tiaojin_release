import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import bookingService from '@/services/bookingService';
import type {
  CreateBookingRequest,
  BookingCreateResponse,
  BookingResponse,
  ApiResponse
} from '@/types';

const bookingKeys = {
  list: (params?: Record<string, unknown>) => ['bookings', 'list', params] as const,
  detail: (id: string) => ['bookings', 'detail', id] as const,
};

interface BookingState {
  selectedBooking: BookingResponse | null;
  bookingFilters: Record<string, unknown>;
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: Record<string, unknown>) => void;
  clearBookingState: () => void;
}

export const useBookingStore = create<BookingState>()(
  persist(
    (set) => ({
      selectedBooking: null,
      bookingFilters: {},
      setSelectedBooking: (booking) => set({ selectedBooking: booking }),
      setBookingFilters: (filters) => set({ bookingFilters: filters }),
      clearBookingState: () => set({ selectedBooking: null, bookingFilters: {} }),
    }),
    {
      name: 'booking-storage',
      partialize: (state) => ({
        bookingFilters: state.bookingFilters,
      }),
    }
  )
);

// Query Hooks
export const useBookingListQuery = (
  params?: Record<string, unknown>,
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<BookingResponse[]>, Error> => {
  return useQuery({
    queryKey: bookingKeys.list(params),
    queryFn: () => bookingService.list(params),
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

export const useBookingDetailQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<BookingResponse>, Error> => {
  const { setSelectedBooking } = useBookingStore();

  return useQuery({
    queryKey: bookingKeys.detail(id),
    queryFn: async () => {
      const response = await bookingService.getById(id);
      if (response.success && response.data) {
        setSelectedBooking(response.data);
      }
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Mutation Hooks
export const useCreateBookingMutation = (): UseMutationResult<
  BookingCreateResponse,
  Error,
  CreateBookingRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateBookingRequest) => bookingService.create(data),
    onSuccess: (response) => {
      // Invalidate and refetch booking list queries
      queryClient.invalidateQueries({ queryKey: ['bookings', 'list'] });

      // If successful, set the created booking data in cache
      if (response.success && response.data) {
        queryClient.setQueryData(
          bookingKeys.detail(response.data.id),
          { success: true, data: response.data }
        );
      }
    },
  });
};

export const useUpdateBookingMutation = (): UseMutationResult<
  ApiResponse<BookingResponse>,
  Error,
  { id: string; data: Partial<CreateBookingRequest> },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => bookingService.update(id, data),
    onSuccess: (response, { id }) => {
      // Invalidate and refetch booking list queries
      queryClient.invalidateQueries({ queryKey: ['bookings', 'list'] });

      // Update the specific booking detail query
      if (response.success && response.data) {
        queryClient.setQueryData(bookingKeys.detail(id), response);
      }
    },
  });
};

export const useUpdateBookingStatusMutation = (): UseMutationResult<
  ApiResponse<BookingResponse>,
  Error,
  { id: string; status: 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, status }) => bookingService.updateStatus(id, status),
    onSuccess: (response, { id }) => {
      // Invalidate and refetch booking list queries
      queryClient.invalidateQueries({ queryKey: ['bookings', 'list'] });

      // Update the specific booking detail query
      if (response.success && response.data) {
        queryClient.setQueryData(bookingKeys.detail(id), response);
      }
    },
  });
};

export const useDeleteBookingMutation = (): UseMutationResult<
  ApiResponse<void>,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();
  const { selectedBooking, setSelectedBooking } = useBookingStore();

  return useMutation({
    mutationFn: (id: string) => bookingService.delete(id),
    onSuccess: (response, id) => {
      if (response.success) {
        // Remove from query cache
        queryClient.removeQueries({ queryKey: bookingKeys.detail(id) });

        // Invalidate list queries
        queryClient.invalidateQueries({ queryKey: ['bookings', 'list'] });

        // Clear selected booking if it was the deleted one
        if (selectedBooking?.id === id) {
          setSelectedBooking(null);
        }
      }
    },
  });
};