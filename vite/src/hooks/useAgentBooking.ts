import { useCallback } from 'react';
import * as React from 'react';
import {
  useAgentBookingStore,
  useAgentBookingListQuery,
  useAgentBookingDetailQuery,
  useAgentBookingDocumentsQuery,
  useAgentCancelBookingMutation,
  useAgentUploadPaymentSlipMutation,
  useAgentUploadPassportMutation
} from '@/stores/agentBookingStore';
import type {
  BookingResponse,
  AdminBookingListFilters,
 
  AgentCancelBookingRequest,
  AgentUploadPaymentSlipRequest,
  AgentUploadPassportRequest,
  AgentBookingDocumentsResponse
} from '@/types';

interface UseAgentBookingReturn {
  // State
  selectedBooking: BookingResponse | null;
  bookingsList: BookingResponse[];
  bookingDocuments: AgentBookingDocumentsResponse['data'] | null;
  isLoading: boolean;

  // Pagination data
  pagination: {
    total: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
  } | null;

  // Actions
  cancelBooking: (id: string, data?: AgentCancelBookingRequest) => Promise<void>;
  uploadPaymentSlip: (id: string, data: AgentUploadPaymentSlipRequest) => Promise<void>;
  uploadPassport: (id: string, data: AgentUploadPassportRequest) => Promise<void>;

  // State management
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: AdminBookingListFilters) => void;
  clearBookingState: () => void;

  // Queries
  getBookingsList: () => void;
  getBookingById: () => void;
  getBookingDocuments: () => void;
  refetchBookingsList: () => void;
  refetchBookingById: () => void;
  refetchBookingDocuments: () => void;

  // Loading states
  isBookingsListLoading: boolean;
  isBookingDetailLoading: boolean;
  isBookingDocumentsLoading: boolean;
  isCancelLoading: boolean;
  isUploadPaymentSlipLoading: boolean;
  isUploadPassportLoading: boolean;
}

export const useAgentBooking = (
  listFilters?: AdminBookingListFilters,
  bookingId?: string,
  options: { enableList?: boolean; enableById?: boolean; enableDocuments?: boolean } = {}
): UseAgentBookingReturn => {
  // Store state
  const {
    selectedBooking,
    setSelectedBooking,
    setBookingFilters,
    clearBookingState
  } = useAgentBookingStore();

  // Mutations
  const cancelBookingMutation = useAgentCancelBookingMutation();
  const uploadPaymentSlipMutation = useAgentUploadPaymentSlipMutation();
  const uploadPassportMutation = useAgentUploadPassportMutation();

  // Queries
  const {
    data: bookingsData,
    refetch: refetchBookingsList,
    isLoading: isBookingsListLoading
  } = useAgentBookingListQuery(listFilters || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchBookingById,
    isLoading: isBookingDetailLoading
  } = useAgentBookingDetailQuery(bookingId || '', {
    enabled: options.enableById !== false && !!bookingId
  });

  const {
    data: documentsData,
    refetch: refetchBookingDocuments,
    isLoading: isBookingDocumentsLoading
  } = useAgentBookingDocumentsQuery(bookingId || '', {
    enabled: options.enableDocuments !== false && !!bookingId
  });

  // Actions
  const cancelBooking = useCallback(async (id: string, data?: AgentCancelBookingRequest) => {
    const result = await cancelBookingMutation.mutateAsync({ id, data });
    if (!result.success) {
      throw new Error(result.message || 'การยกเลิกการจองล้มเหลว');
    }
  }, [cancelBookingMutation]);

  const uploadPaymentSlip = useCallback(async (id: string, data: AgentUploadPaymentSlipRequest) => {
    const result = await uploadPaymentSlipMutation.mutateAsync({ id, data });
    if (!result.success) {
      throw new Error(result.message || 'การอัปโหลดสลิปล้มเหลว');
    }
  }, [uploadPaymentSlipMutation]);

  const uploadPassport = useCallback(async (id: string, data: AgentUploadPassportRequest) => {
    const result = await uploadPassportMutation.mutateAsync({ id, data });
    if (!result.success) {
      throw new Error(result.message || 'การอัปโหลดหนังสือเดินทางล้มเหลว');
    }
  }, [uploadPassportMutation]);

  // Query management
  const getBookingsList = useCallback(() => {
    refetchBookingsList();
  }, [refetchBookingsList]);

  const getBookingById = useCallback(() => {
    refetchBookingById();
  }, [refetchBookingById]);

  const getBookingDocuments = useCallback(() => {
    refetchBookingDocuments();
  }, [refetchBookingDocuments]);

  // Process bookings data and pagination
  const bookingsList = React.useMemo(() => {
    if (!bookingsData?.success || !bookingsData.data) return [];
    return bookingsData.data.bookings || [];
  }, [bookingsData]);

  // Extract pagination info
  const pagination = React.useMemo(() => {
    if (!bookingsData?.success || !bookingsData.data?.pagination) return null;

    const paginationMeta = bookingsData.data.pagination;
    return {
      total: paginationMeta.total,
      totalPages: paginationMeta.total_pages,
      currentPage: paginationMeta.current_page,
      pageSize: paginationMeta.page_size,
    };
  }, [bookingsData]);

  // Process documents data
  const bookingDocuments = React.useMemo(() => {
    if (!documentsData?.success || !documentsData.data) return null;
    return documentsData.data;
  }, [documentsData]);

  // Calculate overall loading state
  const isLoading = cancelBookingMutation.isPending ||
                   uploadPaymentSlipMutation.isPending ||
                   uploadPassportMutation.isPending ||
                   isBookingsListLoading ||
                   isBookingDetailLoading ||
                   isBookingDocumentsLoading;

  return {
    // State
    selectedBooking,
    bookingsList,
    bookingDocuments,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    cancelBooking,
    uploadPaymentSlip,
    uploadPassport,

    // State management
    setSelectedBooking,
    setBookingFilters,
    clearBookingState,

    // Queries
    getBookingsList,
    getBookingById,
    getBookingDocuments,
    refetchBookingsList,
    refetchBookingById,
    refetchBookingDocuments,

    // Loading states
    isBookingsListLoading,
    isBookingDetailLoading,
    isBookingDocumentsLoading,
    isCancelLoading: cancelBookingMutation.isPending,
    isUploadPaymentSlipLoading: uploadPaymentSlipMutation.isPending,
    isUploadPassportLoading: uploadPassportMutation.isPending,
  };
};

export default useAgentBooking;