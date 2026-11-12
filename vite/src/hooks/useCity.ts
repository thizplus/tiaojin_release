import { useCallback } from 'react';
import * as React from 'react';
import {
  useCityStore,
  useCitiesListQuery,
  useCityByIdQuery,
  useCreateCityMutation,
  useUpdateCityMutation,
  useDeleteCityMutation,
  useUpdateCityStatusMutation
} from '@/stores/cityStore';
import type {
  City,
  CityFilter,
  CreateCityRequest,
  UpdateCityRequest,
  UpdateCityStatusRequest
} from '@/types';

interface UseCityReturn {
  // State
  currentCity: City | null;
  selectedCity: City | null;
  citiesList: City[];
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
  createCity: (data: CreateCityRequest) => Promise<City>;
  updateCity: (id: string, data: UpdateCityRequest) => Promise<City>;
  deleteCity: (id: string) => Promise<null>;
  updateCityStatus: (id: string, data: UpdateCityStatusRequest) => Promise<City>;

  // State management
  setCurrentCity: (city: City | null) => void;
  setSelectedCity: (city: City | null) => void;
  setCitiesList: (cities: City[]) => void;
  clearCityData: () => void;

  // Queries
  getCitiesList: () => void;
  getCityById: () => void;
  refetchCitiesList: () => void;
  refetchCityById: () => void;

  // Loading states
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  isDeleteLoading: boolean;
  isUpdateStatusLoading: boolean;
  isCitiesListLoading: boolean;
  isCityByIdLoading: boolean;
}

export const useCity = (
  listParams?: CityFilter,
  cityId?: string,
  options: { enableList?: boolean; enableById?: boolean } = {}
): UseCityReturn => {
  // Store state
  const {
    currentCity,
    selectedCity,
    citiesList,
    setCurrentCity,
    setSelectedCity,
    setCitiesList,
    clearCityData
  } = useCityStore();

  // Mutations
  const createCityMutation = useCreateCityMutation();
  const updateCityMutation = useUpdateCityMutation();
  const deleteCityMutation = useDeleteCityMutation();
  const updateCityStatusMutation = useUpdateCityStatusMutation();

  // Queries
  const {
    data: citiesData,
    refetch: refetchCitiesList,
    isLoading: isCitiesListLoading
  } = useCitiesListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchCityById,
    isLoading: isCityByIdLoading
  } = useCityByIdQuery(cityId || '', {
    enabled: options.enableById !== false && !!cityId
  });

  // Actions
  const createCity = useCallback(async (data: CreateCityRequest) => {
    const result = await createCityMutation.mutateAsync(data);
    return result;
  }, [createCityMutation]);

  const updateCity = useCallback(async (id: string, data: UpdateCityRequest) => {
    const result = await updateCityMutation.mutateAsync({ id, data });
    return result;
  }, [updateCityMutation]);

  const deleteCity = useCallback(async (id: string) => {
    const result = await deleteCityMutation.mutateAsync(id);
    return result;
  }, [deleteCityMutation]);

  const updateCityStatus = useCallback(async (id: string, data: UpdateCityStatusRequest) => {
    const result = await updateCityStatusMutation.mutateAsync({ id, data });
    return result;
  }, [updateCityStatusMutation]);

  // Query management
  const getCitiesList = useCallback(() => {
    refetchCitiesList();
  }, [refetchCitiesList]);

  const getCityById = useCallback(() => {
    refetchCityById();
  }, [refetchCityById]);

  // Process cities data and pagination
  const pagination = React.useMemo(() => {
    if (!citiesData?.data?.pagination) return null;

    const paginationMeta = citiesData.data.pagination;
    return {
      total: paginationMeta.total,
      totalPages: paginationMeta.total_pages,
      currentPage: paginationMeta.current_page,
      pageSize: paginationMeta.page_size,
      hasMore: paginationMeta.has_more,
    };
  }, [citiesData]);

  // Calculate overall loading state
  const isLoading = createCityMutation.isPending ||
                   updateCityMutation.isPending ||
                   deleteCityMutation.isPending ||
                   updateCityStatusMutation.isPending ||
                   isCitiesListLoading ||
                   isCityByIdLoading;

  return {
    // State
    currentCity,
    selectedCity,
    citiesList,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    createCity,
    updateCity,
    deleteCity,
    updateCityStatus,

    // State management
    setCurrentCity,
    setSelectedCity,
    setCitiesList,
    clearCityData,

    // Queries
    getCitiesList,
    getCityById,
    refetchCitiesList,
    refetchCityById,

    // Loading states
    isCreateLoading: createCityMutation.isPending,
    isUpdateLoading: updateCityMutation.isPending,
    isDeleteLoading: deleteCityMutation.isPending,
    isUpdateStatusLoading: updateCityStatusMutation.isPending,
    isCitiesListLoading,
    isCityByIdLoading,
  };
};

export default useCity;