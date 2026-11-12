import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import cityService from '@/services/admin/cityService';
import type {
  City,
  CityFilter,
  CreateCityRequest,
  UpdateCityRequest,
  UpdateCityStatusRequest,
  PaginatedResponse,
  //ApiResponse
} from '@/types';

const cityKeys = {
  cities: (params?: CityFilter) => ['city', 'list', params] as const,
  cityById: (id: string) => ['city', 'byId', id] as const,
};

interface CityState {
  currentCity: City | null;
  selectedCity: City | null;
  citiesList: City[];
  setCurrentCity: (city: City | null) => void;
  setSelectedCity: (city: City | null) => void;
  setCitiesList: (cities: City[]) => void;
  clearCityData: () => void;
}

export const useCityStore = create<CityState>()(
  persist(
    (set) => ({
      currentCity: null,
      selectedCity: null,
      citiesList: [],
      setCurrentCity: (city) => set({ currentCity: city }),
      setSelectedCity: (city) => set({ selectedCity: city }),
      setCitiesList: (cities) => set({ citiesList: cities }),
      clearCityData: () => set({
        currentCity: null,
        selectedCity: null,
        citiesList: [],
      }),
    }),
    {
      name: 'city-storage',
      partialize: (state) => ({
        selectedCity: state.selectedCity,
      }),
    }
  )
);

// Admin+ Queries

/**
 * ดึงรายการเมืองทั้งหมด (สำหรับ admin+)
 */
export const useCitiesListQuery = (
  params: CityFilter = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<PaginatedResponse<City>, Error> => {
  const { setCitiesList } = useCityStore();

  return useQuery({
    queryKey: cityKeys.cities(params),
    queryFn: async () => {
      const response = await cityService.getCities(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch cities list');
      }

      // สมมติว่า response.data มี structure แบบ { cities: City[], pagination: PaginationMeta }
      const cities = response.data.cities || response.data;
      setCitiesList(Array.isArray(cities) ? cities : []);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

/**
 * ดึงข้อมูลเมืองตาม ID (สำหรับ admin+)
 */
export const useCityByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<City, Error> => {
  const { setSelectedCity } = useCityStore();

  return useQuery({
    queryKey: cityKeys.cityById(id),
    queryFn: async () => {
      const response = await cityService.getCityById(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch city data');
      }

      setSelectedCity(response.data);
      return response.data;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Admin+ Mutations

/**
 * สร้างเมืองใหม่ (สำหรับ admin+)
 */
export const useCreateCityMutation = (): UseMutationResult<
  City,
  Error,
  CreateCityRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateCityRequest) => {
      const response = await cityService.createCity(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to create city');
      }

      return response.data;
    },
    onSuccess: () => {
      // Invalidate cities list to refetch
      queryClient.invalidateQueries({ queryKey: cityKeys.cities() });
    },
  });
};

/**
 * อัปเดตข้อมูลเมือง (สำหรับ admin+)
 */
export const useUpdateCityMutation = (): UseMutationResult<
  City,
  Error,
  { id: string; data: UpdateCityRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await cityService.updateCity(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to update city');
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: cityKeys.cityById(id) });
      queryClient.invalidateQueries({ queryKey: cityKeys.cities() });
    },
  });
};

/**
 * ลบเมือง (สำหรับ admin+)
 */
export const useDeleteCityMutation = (): UseMutationResult<
  null,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await cityService.deleteCity(id);

      if (!response.success) {
        throw new Error(response.message || 'Failed to delete city');
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: cityKeys.cityById(id) });
      queryClient.invalidateQueries({ queryKey: cityKeys.cities() });
    },
  });
};

/**
 * อัปเดตสถานะเมือง (สำหรับ admin+)
 */
export const useUpdateCityStatusMutation = (): UseMutationResult<
  City,
  Error,
  { id: string; data: UpdateCityStatusRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await cityService.updateCityStatus(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to update city status');
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: cityKeys.cityById(id) });
      queryClient.invalidateQueries({ queryKey: cityKeys.cities() });
    },
  });
};