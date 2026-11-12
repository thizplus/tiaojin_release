import { useCallback } from 'react';
import * as React from 'react';
import {
  useBookingStore,
  useBookingListQuery,
  useBookingDetailQuery,
  useCreateBookingMutation,
  useUpdateBookingMutation,
  useUpdateBookingStatusMutation,
  useDeleteBookingMutation
} from '@/stores/bookingStore';
import type {
  BookingResponse,
  CreateBookingRequest,

} from '@/types';

interface UseBookingReturn {
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
    hasMore: boolean;
  } | null;

  // Actions
  createBooking: (data: CreateBookingRequest) => Promise<BookingResponse>;
  updateBooking: (id: string, data: Partial<CreateBookingRequest>) => Promise<BookingResponse>;
  updateBookingStatus: (id: string, status: 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed') => Promise<BookingResponse>;
  deleteBooking: (id: string) => Promise<void>;

  // State management
  setSelectedBooking: (booking: BookingResponse | null) => void;
  setBookingFilters: (filters: Record<string, unknown>) => void;
  clearBookingState: () => void;

  // Queries
  getBookingsList: () => void;
  getBookingById: () => void;
  refetchBookingsList: () => void;
  refetchBookingById: () => void;

  // Loading states
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  isUpdateStatusLoading: boolean;
  isDeleteLoading: boolean;
  isBookingsListLoading: boolean;
  isBookingDetailLoading: boolean;
}

export const useBooking = (
  listParams?: Record<string, unknown>,
  bookingId?: string,
  options: { enableList?: boolean; enableById?: boolean } = {}
): UseBookingReturn => {
  // Store state
  const {
    selectedBooking,
    setSelectedBooking,
    setBookingFilters,
    clearBookingState
  } = useBookingStore();

  // Mutations
  const createBookingMutation = useCreateBookingMutation();
  const updateBookingMutation = useUpdateBookingMutation();
  const updateBookingStatusMutation = useUpdateBookingStatusMutation();
  const deleteBookingMutation = useDeleteBookingMutation();

  // Queries
  const {
    data: bookingsData,
    refetch: refetchBookingsList,
    isLoading: isBookingsListLoading
  } = useBookingListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchBookingById,
    isLoading: isBookingDetailLoading
  } = useBookingDetailQuery(bookingId || '', {
    enabled: options.enableById !== false && !!bookingId
  });

  // Actions
  const createBooking = useCallback(async (data: CreateBookingRequest) => {
    const result = await createBookingMutation.mutateAsync(data);
    if (!result.success || !result.data) {
      throw new Error(result.message || 'การสร้างการจองล้มเหลว');
    }
    return result.data;
  }, [createBookingMutation]);

  const updateBooking = useCallback(async (id: string, data: Partial<CreateBookingRequest>) => {
    const result = await updateBookingMutation.mutateAsync({ id, data });
    if (!result.success || !result.data) {
      throw new Error(result.message || 'การอัปเดตการจองล้มเหลว');
    }
    return result.data;
  }, [updateBookingMutation]);

  const updateBookingStatus = useCallback(async (
    id: string,
    status: 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed'
  ) => {
    const result = await updateBookingStatusMutation.mutateAsync({ id, status });
    if (!result.success || !result.data) {
      throw new Error(result.message || 'การอัปเดตสถานะการจองล้มเหลว');
    }
    return result.data;
  }, [updateBookingStatusMutation]);

  const deleteBooking = useCallback(async (id: string) => {
    const result = await deleteBookingMutation.mutateAsync(id);
    if (!result.success) {
      throw new Error(result.message || 'การลบการจองล้มเหลว');
    }
  }, [deleteBookingMutation]);

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
    return Array.isArray(bookingsData.data) ? bookingsData.data : [];
  }, [bookingsData]);

  // If using paginated response, extract pagination info
  const pagination = React.useMemo(() => {
    if (!bookingsData?.data) return null;

    // Check if response has pagination structure
    const paginatedData = bookingsData.data as { pagination?: { total: number; total_pages: number; current_page: number; page_size: number; has_more: boolean } };
    if (paginatedData?.pagination) {
      const paginationMeta = paginatedData.pagination;
      return {
        total: paginationMeta.total,
        totalPages: paginationMeta.total_pages,
        currentPage: paginationMeta.current_page,
        pageSize: paginationMeta.page_size,
        hasMore: paginationMeta.has_more,
      };
    }

    return null;
  }, [bookingsData]);

  // Calculate overall loading state
  const isLoading = createBookingMutation.isPending ||
                   updateBookingMutation.isPending ||
                   updateBookingStatusMutation.isPending ||
                   deleteBookingMutation.isPending ||
                   isBookingsListLoading ||
                   isBookingDetailLoading;

  return {
    // State
    selectedBooking,
    bookingsList,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    createBooking,
    updateBooking,
    updateBookingStatus,
    deleteBooking,

    // State management
    setSelectedBooking,
    setBookingFilters,
    clearBookingState,

    // Queries
    getBookingsList,
    getBookingById,
    refetchBookingsList,
    refetchBookingById,

    // Loading states
    isCreateLoading: createBookingMutation.isPending,
    isUpdateLoading: updateBookingMutation.isPending,
    isUpdateStatusLoading: updateBookingStatusMutation.isPending,
    isDeleteLoading: deleteBookingMutation.isPending,
    isBookingsListLoading,
    isBookingDetailLoading,
  };
};

export default useBooking;