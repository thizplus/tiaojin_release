import { useCallback } from 'react';
import * as React from 'react';
import {
  useAdminBookingStore,
  useAdminBookingListQuery,
  useAdminBookingStatsQuery,
  useAdminBookingDetailQuery,
  useAdminUpdateBookingStatusMutation,
  useAdminCreatePaymentMutation,
  useAdminDeleteBookingMutation,
  useAdminVerifyPaymentSlipMutation,
  useAdminReversePaymentSlipMutation
} from '@/stores/adminBookingStore';
import type {
  BookingResponse,
  AdminBookingListFilters,
  AdminBookingStatusUpdateRequest,
  AdminCreatePaymentRequest,
  AdminBookingStatsResponse,
  AdminCreatePaymentResponse,
  PaymentSlipVerifyRequest,
  PaymentSlipReverseRequest
} from '@/types';

interface UseAdminBookingReturn {
  // State
  selectedBooking: BookingResponse | null;
  bookingsList: BookingResponse[];
  bookingsStats: AdminBookingStatsResponse['data'] | null;
  isLoading: boolean;

  // Pagination data
  pagination: {
    total: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
  } | null;

  // Actions
  updateBookingStatus: (id: string, data: AdminBookingStatusUpdateRequest) => Promise<BookingResponse>;
  createPayment: (bookingId: string, data: AdminCreatePaymentRequest) => Promise<AdminCreatePaymentResponse['data']>;
  deleteBooking: (id: string) => Promise<void>;
  verifyPaymentSlip: (paymentSlipId: string, data: PaymentSlipVerifyRequest) => Promise<void>;
  reversePaymentSlip: (paymentSlipId: string, data: PaymentSlipReverseRequest) => Promise<void>;

  // State management
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: AdminBookingListFilters) => void;
  clearBookingState: () => void;

  // Queries
  getBookingsList: () => void;
  getBookingsStats: () => void;
  getBookingById: () => void;
  refetchBookingsList: () => void;
  refetchBookingsStats: () => void;
  refetchBookingById: () => void;

  // Loading states
  isUpdateStatusLoading: boolean;
  isCreatePaymentLoading: boolean;
  isDeleteLoading: boolean;
  isVerifyPaymentSlipLoading: boolean;
  isReversePaymentSlipLoading: boolean;
  isBookingsListLoading: boolean;
  isBookingsStatsLoading: boolean;
  isBookingDetailLoading: boolean;
}

export const useAdminBooking = (
  listFilters?: AdminBookingListFilters,
  bookingId?: string,
  options: { enableList?: boolean; enableStats?: boolean; enableById?: boolean } = {}
): UseAdminBookingReturn => {
  // Store state
  const {
    selectedBooking,
    setSelectedBooking,
    setBookingFilters,
    clearBookingState
  } = useAdminBookingStore();

  // Mutations
  const updateBookingStatusMutation = useAdminUpdateBookingStatusMutation();
  const createPaymentMutation = useAdminCreatePaymentMutation();
  const deleteBookingMutation = useAdminDeleteBookingMutation();
  const verifyPaymentSlipMutation = useAdminVerifyPaymentSlipMutation();
  const reversePaymentSlipMutation = useAdminReversePaymentSlipMutation();

  // Queries
  const {
    data: bookingsData,
    refetch: refetchBookingsList,
    isLoading: isBookingsListLoading
  } = useAdminBookingListQuery(listFilters || {}, {
    enabled: options.enableList === true
  });

  const {
    data: statsData,
    refetch: refetchBookingsStats,
    isLoading: isBookingsStatsLoading
  } = useAdminBookingStatsQuery({
    enabled: options.enableStats === true
  });

  const {
    refetch: refetchBookingById,
    isLoading: isBookingDetailLoading
  } = useAdminBookingDetailQuery(bookingId || '', {
    enabled: options.enableById !== false && !!bookingId
  });

  // Actions
  const updateBookingStatus = useCallback(async (id: string, data: AdminBookingStatusUpdateRequest) => {
    const result = await updateBookingStatusMutation.mutateAsync({ id, data });
    if (!result.success || !result.data) {
      throw new Error(result.message || 'การอัปเดตสถานะการจองล้มเหลว');
    }
    return result.data;
  }, [updateBookingStatusMutation]);

  const createPayment = useCallback(async (bookingId: string, data: AdminCreatePaymentRequest) => {
    const result = await createPaymentMutation.mutateAsync({ bookingId, data });
    if (!result.success || !result.data) {
      throw new Error(result.message || 'การสร้างการชำระเงินล้มเหลว');
    }
    return result.data;
  }, [createPaymentMutation]);

  const deleteBooking = useCallback(async (id: string) => {
    const result = await deleteBookingMutation.mutateAsync(id);
    if (!result.success) {
      throw new Error(result.message || 'การยกเลิกการจองล้มเหลว');
    }
  }, [deleteBookingMutation]);

  const verifyPaymentSlip = useCallback(async (paymentSlipId: string, data: PaymentSlipVerifyRequest) => {
    const result = await verifyPaymentSlipMutation.mutateAsync({ paymentSlipId, data });
    if (!result.success) {
      throw new Error(result.message || 'การยืนยันสลิปการชำระเงินล้มเหลว');
    }
  }, [verifyPaymentSlipMutation]);

  const reversePaymentSlip = useCallback(async (paymentSlipId: string, data: PaymentSlipReverseRequest) => {
    const result = await reversePaymentSlipMutation.mutateAsync({ paymentSlipId, data });
    if (!result.success) {
      throw new Error(result.message || 'การยกเลิกสลิปการชำระเงินล้มเหลว');
    }
  }, [reversePaymentSlipMutation]);

  // Query management
  const getBookingsList = useCallback(() => {
    refetchBookingsList();
  }, [refetchBookingsList]);

  const getBookingsStats = useCallback(() => {
    refetchBookingsStats();
  }, [refetchBookingsStats]);

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

  // Process stats data
  const bookingsStats = React.useMemo(() => {
    if (!statsData?.success || !statsData.data) return null;
    return statsData.data;
  }, [statsData]);

  // Calculate overall loading state
  const isLoading = updateBookingStatusMutation.isPending ||
                   createPaymentMutation.isPending ||
                   deleteBookingMutation.isPending ||
                   verifyPaymentSlipMutation.isPending ||
                   reversePaymentSlipMutation.isPending ||
                   isBookingsListLoading ||
                   isBookingsStatsLoading ||
                   isBookingDetailLoading;

  return {
    // State
    selectedBooking,
    bookingsList,
    bookingsStats,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    updateBookingStatus,
    createPayment,
    deleteBooking,
    verifyPaymentSlip,
    reversePaymentSlip,

    // State management
    setSelectedBooking,
    setBookingFilters,
    clearBookingState,

    // Queries
    getBookingsList,
    getBookingsStats,
    getBookingById,
    refetchBookingsList,
    refetchBookingsStats,
    refetchBookingById,

    // Loading states
    isUpdateStatusLoading: updateBookingStatusMutation.isPending,
    isCreatePaymentLoading: createPaymentMutation.isPending,
    isDeleteLoading: deleteBookingMutation.isPending,
    isVerifyPaymentSlipLoading: verifyPaymentSlipMutation.isPending,
    isReversePaymentSlipLoading: reversePaymentSlipMutation.isPending,
    isBookingsListLoading,
    isBookingsStatsLoading,
    isBookingDetailLoading,
  };
};

export default useAdminBooking;