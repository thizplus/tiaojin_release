import { useCallback } from "react";
import * as React from "react";
import {
  useTourStore,
  useToursListQuery,
  useTourByIdQuery,
  useCreateTourMutation,
  useUpdateTourMutation,
  useDeleteTourMutation,
  useUpdateTourStatusMutation
} from "@/stores/tourStore";
import type {
  Tour,
  TourFilter,
  CreateTourRequest,
  UpdateTourRequest,
  UpdateTourStatusRequest
} from "@/types";

interface UseToursReturn {
  // State
  selectedTour: Tour | null;
  toursList: Tour[];
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
  createTour: (data: CreateTourRequest) => Promise<Tour>;
  updateTour: (id: string, data: UpdateTourRequest) => Promise<Tour>;
  deleteTour: (id: string) => Promise<null>;
  updateTourStatus: (id: string, data: UpdateTourStatusRequest) => Promise<Tour>;

  // State management
  setSelectedTour: (tour: Tour | null) => void;
  setToursList: (tours: Tour[]) => void;
  clearTourData: () => void;

  // Queries
  getToursList: (params?: TourFilter, options?: { enabled?: boolean }) => void;
  getTourById: (id: string, options?: { enabled?: boolean }) => void;
  refetchToursList: () => void;
  refetchTourById: () => void;

  // Loading states
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  isDeleteLoading: boolean;
  isUpdateStatusLoading: boolean;
  isToursListLoading: boolean;
  isTourByIdLoading: boolean;
}

export const useTours = (
  listParams?: TourFilter,
  tourId?: string,
  options: { enableList?: boolean; enableById?: boolean } = {}
): UseToursReturn => {
  // Store state
  const {
    selectedTour,
    toursList,
    setSelectedTour,
    setToursList,
    clearTourData
  } = useTourStore();

  // Mutations
  const createTourMutation = useCreateTourMutation();
  const updateTourMutation = useUpdateTourMutation();
  const deleteTourMutation = useDeleteTourMutation();
  const updateTourStatusMutation = useUpdateTourStatusMutation();

  // Queries - ต้อง explicit enable เพื่อป้องกัน duplicate queries
  const {
    data: toursData,
    refetch: refetchToursList,
    isLoading: isToursListLoading
  } = useToursListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchTourById,
    isLoading: isTourByIdLoading
  } = useTourByIdQuery(tourId || "", {
    enabled: options.enableById !== false && !!tourId
  });

  // Actions
  const createTour = useCallback(async (data: CreateTourRequest) => {
    const result = await createTourMutation.mutateAsync(data);
    return result;
  }, [createTourMutation]);

  const updateTour = useCallback(async (id: string, data: UpdateTourRequest) => {
    const result = await updateTourMutation.mutateAsync({ id, data });
    return result;
  }, [updateTourMutation]);

  const deleteTour = useCallback(async (id: string) => {
    const result = await deleteTourMutation.mutateAsync(id);
    return result;
  }, [deleteTourMutation]);

  const updateTourStatus = useCallback(async (id: string, data: UpdateTourStatusRequest) => {
    const result = await updateTourStatusMutation.mutateAsync({ id, data });
    return result;
  }, [updateTourStatusMutation]);

  // Query management
  const getToursList = useCallback(() => {
    refetchToursList();
  }, [refetchToursList]);

  const getTourById = useCallback(() => {
    refetchTourById();
  }, [refetchTourById]);

  // Process tours data and pagination
  const pagination = React.useMemo(() => {
    if (!toursData?.data?.pagination) return null;

    const paginationMeta = toursData.data.pagination;
    return {
      total: paginationMeta.total,
      totalPages: paginationMeta.total_pages,
      currentPage: paginationMeta.current_page,
      pageSize: paginationMeta.page_size || paginationMeta.limit,
      hasMore: paginationMeta.has_more,
    };
  }, [toursData]);

  // Calculate overall loading state
  const isLoading = createTourMutation.isPending ||
                   updateTourMutation.isPending ||
                   deleteTourMutation.isPending ||
                   updateTourStatusMutation.isPending ||
                   isToursListLoading ||
                   isTourByIdLoading;

  return {
    // State
    selectedTour,
    toursList,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    createTour,
    updateTour,
    deleteTour,
    updateTourStatus,

    // State management
    setSelectedTour,
    setToursList,
    clearTourData,

    // Queries
    getToursList,
    getTourById,
    refetchToursList,
    refetchTourById,

    // Loading states
    isCreateLoading: createTourMutation.isPending,
    isUpdateLoading: updateTourMutation.isPending,
    isDeleteLoading: deleteTourMutation.isPending,
    isUpdateStatusLoading: updateTourStatusMutation.isPending,
    isToursListLoading,
    isTourByIdLoading,
  };
};

export default useTours;