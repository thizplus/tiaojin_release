import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import saleBookingService from '@/services/saleBookingService';
import type {
  AdminBookingListResponse,
  AdminBookingListFilters,
  BookingResponse,
  ApiResponse,
  AgentUploadPaymentSlipRequest,
  AgentUploadPaymentSlipResponse,
  AgentUploadPassportRequest,
  AgentUploadPassportResponse,
} from '@/types';

const saleBookingKeys = {
  list: (filters?: AdminBookingListFilters) => ['sale', 'bookings', 'list', filters] as const,
  detail: (id: string) => ['sale', 'bookings', 'detail', id] as const,
};

interface SaleBookingState {
  selectedBooking: BookingResponse | null;
  bookingFilters: AdminBookingListFilters;
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: AdminBookingListFilters) => void;
  clearBookingState: () => void;
}

export const useSaleBookingStore = create<SaleBookingState>()(
  persist(
    (set) => ({
      selectedBooking: null,
      bookingFilters: {},
      setSelectedBooking: (booking) => set({ selectedBooking: booking }),
      setBookingFilters: (filters) => set({ bookingFilters: filters }),
      clearBookingState: () => set({ selectedBooking: null, bookingFilters: {} }),
    }),
    {
      name: 'sale-booking-storage',
      partialize: (state) => ({
        bookingFilters: state.bookingFilters,
      }),
    }
  )
);

// Query Hooks
export const useSaleBookingListQuery = (
  filters?: AdminBookingListFilters,
  options: { enabled?: boolean } = {}
): UseQueryResult<AdminBookingListResponse, Error> => {
  return useQuery({
    queryKey: saleBookingKeys.list(filters),
    queryFn: () => saleBookingService.list(filters),
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

export const useSaleBookingDetailQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<BookingResponse>, Error> => {
  const { setSelectedBooking } = useSaleBookingStore();

  return useQuery({
    queryKey: saleBookingKeys.detail(id),
    queryFn: async () => {
      const response = await saleBookingService.getById(id);
      if (response.success && response.data) {
        setSelectedBooking(response.data);
      }
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Mutation Hooks
export const useUploadPaymentSlipMutation = (): UseMutationResult<
  AgentUploadPaymentSlipResponse,
  Error,
  { id: string; data: AgentUploadPaymentSlipRequest }
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) =>
      saleBookingService.uploadPaymentSlip(id, data.payment_slip, data.amount),
    onSuccess: () => {
      // Invalidate list queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['sale', 'bookings'] });
    },
  });
};

export const useUploadPassportMutation = (): UseMutationResult<
  AgentUploadPassportResponse,
  Error,
  { id: string; data: AgentUploadPassportRequest }
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) =>
      saleBookingService.uploadPassport(id, data.passport),
    onSuccess: () => {
      // Invalidate list queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['sale', 'bookings'] });
    },
  });
};