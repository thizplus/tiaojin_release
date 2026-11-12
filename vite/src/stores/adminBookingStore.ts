import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import adminBookingService from '@/services/adminBookingService';
import type {
  AdminBookingListResponse,
  AdminBookingListFilters,
  AdminBookingStatsResponse,
  BookingResponse,
  AdminBookingStatusUpdateRequest,
  AdminBookingStatusUpdateResponse,
  AdminCreatePaymentRequest,
  AdminCreatePaymentResponse,
  PaymentSlipVerifyRequest,
  PaymentSlipVerifyResponse,
  PaymentSlipReverseRequest,
  PaymentSlipReverseResponse,
  ApiResponse
} from '@/types';

const adminBookingKeys = {
  list: (filters?: AdminBookingListFilters) => ['admin', 'bookings', 'list', filters] as const,
  stats: () => ['admin', 'bookings', 'stats'] as const,
  detail: (id: string) => ['admin', 'bookings', 'detail', id] as const,
};

interface AdminBookingState {
  selectedBooking: BookingResponse | null;
  bookingFilters: AdminBookingListFilters;
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: AdminBookingListFilters) => void;
  clearBookingState: () => void;
}

export const useAdminBookingStore = create<AdminBookingState>()(
  persist(
    (set) => ({
      selectedBooking: null,
      bookingFilters: {},
      setSelectedBooking: (booking) => set({ selectedBooking: booking }),
      setBookingFilters: (filters) => set({ bookingFilters: filters }),
      clearBookingState: () => set({ selectedBooking: null, bookingFilters: {} }),
    }),
    {
      name: 'admin-booking-storage',
      partialize: (state) => ({
        bookingFilters: state.bookingFilters,
      }),
    }
  )
);

// Query Hooks
export const useAdminBookingListQuery = (
  filters?: AdminBookingListFilters,
  options: { enabled?: boolean } = {}
): UseQueryResult<AdminBookingListResponse, Error> => {
  return useQuery({
    queryKey: adminBookingKeys.list(filters),
    queryFn: () => adminBookingService.list(filters),
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

export const useAdminBookingStatsQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<AdminBookingStatsResponse, Error> => {
  return useQuery({
    queryKey: adminBookingKeys.stats(),
    queryFn: () => adminBookingService.getStats(),
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

export const useAdminBookingDetailQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<BookingResponse>, Error> => {
  const { setSelectedBooking } = useAdminBookingStore();

  return useQuery({
    queryKey: adminBookingKeys.detail(id),
    queryFn: async () => {
      const response = await adminBookingService.getById(id);
      if (response.success && response.data) {
        setSelectedBooking(response.data);
      }
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Mutation Hooks
export const useAdminUpdateBookingStatusMutation = (): UseMutationResult<
  AdminBookingStatusUpdateResponse,
  Error,
  { id: string; data: AdminBookingStatusUpdateRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => adminBookingService.updateStatus(id, data),
    onSuccess: (response, { id }) => {
      // Invalidate and refetch admin booking list queries
      queryClient.invalidateQueries({ queryKey: ['admin', 'bookings', 'list'] });

      // Invalidate stats query
      queryClient.invalidateQueries({ queryKey: adminBookingKeys.stats() });

      // Update the specific booking detail query
      if (response.success && response.data) {
        queryClient.setQueryData(adminBookingKeys.detail(id), {
          success: true,
          data: response.data
        });
      }
    },
  });
};

export const useAdminCreatePaymentMutation = (): UseMutationResult<
  AdminCreatePaymentResponse,
  Error,
  { bookingId: string; data: AdminCreatePaymentRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ bookingId, data }) => adminBookingService.createPayment(bookingId, data),
    onSuccess: (_, { bookingId }) => {
      // Invalidate and refetch admin booking list queries
      queryClient.invalidateQueries({ queryKey: ['admin', 'bookings', 'list'] });

      // Invalidate stats query
      queryClient.invalidateQueries({ queryKey: adminBookingKeys.stats() });

      // Invalidate the specific booking detail query to refresh payment data
      queryClient.invalidateQueries({ queryKey: adminBookingKeys.detail(bookingId) });
    },
  });
};

export const useAdminDeleteBookingMutation = (): UseMutationResult<
  ApiResponse<void>,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();
  const { selectedBooking, setSelectedBooking } = useAdminBookingStore();

  return useMutation({
    mutationFn: (id: string) => adminBookingService.delete(id),
    onSuccess: (response, id) => {
      if (response.success) {
        // Remove from query cache
        queryClient.removeQueries({ queryKey: adminBookingKeys.detail(id) });

        // Invalidate list queries
        queryClient.invalidateQueries({ queryKey: ['admin', 'bookings', 'list'] });

        // Invalidate stats query
        queryClient.invalidateQueries({ queryKey: adminBookingKeys.stats() });

        // Clear selected booking if it was the deleted one
        if (selectedBooking?.id === id) {
          setSelectedBooking(null);
        }
      }
    },
  });
};

// Payment Slip Verification Mutations
export const useAdminVerifyPaymentSlipMutation = (): UseMutationResult<
  PaymentSlipVerifyResponse,
  Error,
  { paymentSlipId: string; data: PaymentSlipVerifyRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ paymentSlipId, data }) => adminBookingService.verifyPaymentSlip(paymentSlipId, data),
    onSuccess: (response) => {
      if (response.success) {
        // Invalidate list queries to refresh booking data
        queryClient.invalidateQueries({ queryKey: ['admin', 'bookings', 'list'] });

        // Invalidate stats query
        queryClient.invalidateQueries({ queryKey: adminBookingKeys.stats() });

        // Invalidate specific booking detail if available
        if (response.data.booking_update.booking_id) {
          queryClient.invalidateQueries({
            queryKey: adminBookingKeys.detail(response.data.booking_update.booking_id)
          });
        }
      }
    },
  });
};

export const useAdminReversePaymentSlipMutation = (): UseMutationResult<
  PaymentSlipReverseResponse,
  Error,
  { paymentSlipId: string; data: PaymentSlipReverseRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ paymentSlipId, data }) => adminBookingService.reversePaymentSlip(paymentSlipId, data),
    onSuccess: (response) => {
      if (response.success) {
        // Invalidate list queries to refresh booking data
        queryClient.invalidateQueries({ queryKey: ['admin', 'bookings', 'list'] });

        // Invalidate stats query
        queryClient.invalidateQueries({ queryKey: adminBookingKeys.stats() });

        // Invalidate specific booking detail if available
        if (response.data.booking_update.booking_id) {
          queryClient.invalidateQueries({
            queryKey: adminBookingKeys.detail(response.data.booking_update.booking_id)
          });
        }
      }
    },
  });
};