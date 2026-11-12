import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import agentBookingService from '@/services/agentBookingService';
import type {
  AdminBookingListResponse,
  AdminBookingListFilters,
  BookingResponse,
  ApiResponse,
  AgentCancelBookingRequest,
  AgentCancelBookingResponse,
  AgentUploadPaymentSlipRequest,
  AgentUploadPaymentSlipResponse,
  AgentUploadPassportRequest,
  AgentUploadPassportResponse,
  AgentBookingDocumentsResponse
} from '@/types';

const agentBookingKeys = {
  list: (filters?: AdminBookingListFilters) => ['agent', 'bookings', 'list', filters] as const,
  detail: (id: string) => ['agent', 'bookings', 'detail', id] as const,
  documents: (id: string) => ['agent', 'bookings', 'documents', id] as const,
};

interface AgentBookingState {
  selectedBooking: BookingResponse | null;
  bookingFilters: AdminBookingListFilters;
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: AdminBookingListFilters) => void;
  clearBookingState: () => void;
}

export const useAgentBookingStore = create<AgentBookingState>()(
  persist(
    (set) => ({
      selectedBooking: null,
      bookingFilters: {},
      setSelectedBooking: (booking) => set({ selectedBooking: booking }),
      setBookingFilters: (filters) => set({ bookingFilters: filters }),
      clearBookingState: () => set({ selectedBooking: null, bookingFilters: {} }),
    }),
    {
      name: 'agent-booking-storage',
      partialize: (state) => ({
        bookingFilters: state.bookingFilters,
      }),
    }
  )
);

// Query Hooks
export const useAgentBookingListQuery = (
  filters?: AdminBookingListFilters,
  options: { enabled?: boolean } = {}
): UseQueryResult<AdminBookingListResponse, Error> => {
  return useQuery({
    queryKey: agentBookingKeys.list(filters),
    queryFn: () => agentBookingService.list(filters),
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

export const useAgentBookingDetailQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<ApiResponse<BookingResponse>, Error> => {
  const { setSelectedBooking } = useAgentBookingStore();

  return useQuery({
    queryKey: agentBookingKeys.detail(id),
    queryFn: async () => {
      const response = await agentBookingService.getById(id);
      if (response.success && response.data) {
        setSelectedBooking(response.data);
      }
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

export const useAgentBookingDocumentsQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<AgentBookingDocumentsResponse, Error> => {
  return useQuery({
    queryKey: agentBookingKeys.documents(id),
    queryFn: () => agentBookingService.getDocuments(id),
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Mutation Hooks
export const useAgentCancelBookingMutation = (): UseMutationResult<
  AgentCancelBookingResponse,
  Error,
  { id: string; data?: AgentCancelBookingRequest },
  unknown
> => {
  const queryClient = useQueryClient();
  const { selectedBooking, setSelectedBooking } = useAgentBookingStore();

  return useMutation({
    mutationFn: ({ id, data }) => agentBookingService.cancel(id, data),
    onSuccess: (response, { id }) => {
      if (response.success) {
        // Invalidate and refetch agent booking list queries
        queryClient.invalidateQueries({ queryKey: ['agent', 'bookings', 'list'] });

        // Update the specific booking detail query if it exists
        queryClient.invalidateQueries({ queryKey: agentBookingKeys.detail(id) });

        // Clear selected booking if it was the cancelled one
        if (selectedBooking?.id === id) {
          setSelectedBooking(null);
        }
      }
    },
  });
};

export const useAgentUploadPaymentSlipMutation = (): UseMutationResult<
  AgentUploadPaymentSlipResponse,
  Error,
  { id: string; data: AgentUploadPaymentSlipRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => agentBookingService.uploadPaymentSlip(id, data),
    onSuccess: (response, { id }) => {
      if (response.success) {
        // Invalidate booking list to refresh
        queryClient.invalidateQueries({ queryKey: ['agent', 'bookings', 'list'] });

        // Invalidate the specific booking detail to refresh payment slips
        queryClient.invalidateQueries({ queryKey: agentBookingKeys.detail(id) });

        // Invalidate documents query to refresh
        queryClient.invalidateQueries({ queryKey: agentBookingKeys.documents(id) });
      }
    },
  });
};

export const useAgentUploadPassportMutation = (): UseMutationResult<
  AgentUploadPassportResponse,
  Error,
  { id: string; data: AgentUploadPassportRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => agentBookingService.uploadPassport(id, data),
    onSuccess: (response, { id }) => {
      if (response.success) {
        // Invalidate booking list to refresh
        queryClient.invalidateQueries({ queryKey: ['agent', 'bookings', 'list'] });

        // Invalidate the specific booking detail to refresh passports
        queryClient.invalidateQueries({ queryKey: agentBookingKeys.detail(id) });

        // Invalidate documents query to refresh
        queryClient.invalidateQueries({ queryKey: agentBookingKeys.documents(id) });
      }
    },
  });
};