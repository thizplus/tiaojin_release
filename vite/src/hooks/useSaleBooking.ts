import { useCallback } from 'react';
import * as React from 'react';
import {
  useSaleBookingStore,
  useSaleBookingListQuery,
  useSaleBookingDetailQuery,
  useUploadPaymentSlipMutation,
  useUploadPassportMutation,
} from '@/stores/saleBookingStore';
import type {
  BookingResponse,
  AdminBookingListFilters,
  AgentUploadPaymentSlipRequest,
  AgentUploadPassportRequest,
} from '@/types';

interface UseSaleBookingReturn {
  // State
  selectedBooking: BookingResponse | null;
  bookingsList: BookingResponse[];
  isLoading: boolean;

  // Pagination data
  pagination: {
    total: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
  } | null;

  // State management
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: AdminBookingListFilters) => void;
  clearBookingState: () => void;

  // Queries
  getBookingsList: () => void;
  getBookingById: () => void;
  refetchBookingsList: () => void;
  refetchBookingById: () => void;

  // Actions (แทนลูกค้า)
  uploadPaymentSlip: (id: string, data: AgentUploadPaymentSlipRequest) => Promise<void>;
  uploadPassport: (id: string, data: AgentUploadPassportRequest) => Promise<void>;

  // Loading states
  isBookingsListLoading: boolean;
  isBookingDetailLoading: boolean;
  isUploadPaymentSlipLoading: boolean;
  isUploadPassportLoading: boolean;
}

export const useSaleBooking = (
  listFilters?: AdminBookingListFilters,
  bookingId?: string,
  options: { enableList?: boolean; enableById?: boolean } = {}
): UseSaleBookingReturn => {
  // Store state
  const {
    selectedBooking,
    setSelectedBooking,
    setBookingFilters,
    clearBookingState
  } = useSaleBookingStore();

  // Queries
  const {
    data: bookingsData,
    refetch: refetchBookingsList,
    isLoading: isBookingsListLoading
  } = useSaleBookingListQuery(listFilters || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchBookingById,
    isLoading: isBookingDetailLoading
  } = useSaleBookingDetailQuery(bookingId || '', {
    enabled: options.enableById !== false && !!bookingId
  });

  // Mutations
  const uploadPaymentSlipMutation = useUploadPaymentSlipMutation();
  const uploadPassportMutation = useUploadPassportMutation();

  // Query management
  const getBookingsList = useCallback(() => {
    refetchBookingsList();
  }, [refetchBookingsList]);

  const getBookingById = useCallback(() => {
    refetchBookingById();
  }, [refetchBookingById]);

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

  // Calculate overall loading state
  const isLoading = isBookingsListLoading || isBookingDetailLoading;

  // Actions
  const uploadPaymentSlip = useCallback(async (id: string, data: AgentUploadPaymentSlipRequest) => {
    await uploadPaymentSlipMutation.mutateAsync({ id, data });
  }, [uploadPaymentSlipMutation]);

  const uploadPassport = useCallback(async (id: string, data: AgentUploadPassportRequest) => {
    await uploadPassportMutation.mutateAsync({ id, data });
  }, [uploadPassportMutation]);

  return {
    // State
    selectedBooking,
    bookingsList,
    isLoading,

    // Pagination data
    pagination,

    // State management
    setSelectedBooking,
    setBookingFilters,
    clearBookingState,

    // Queries
    getBookingsList,
    getBookingById,
    refetchBookingsList,
    refetchBookingById,

    // Actions (แทนลูกค้า)
    uploadPaymentSlip,
    uploadPassport,

    // Loading states
    isBookingsListLoading,
    isBookingDetailLoading,
    isUploadPaymentSlipLoading: uploadPaymentSlipMutation.isPending,
    isUploadPassportLoading: uploadPassportMutation.isPending,
  };
};

export default useSaleBooking;