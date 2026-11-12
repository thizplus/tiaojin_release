import { useQuery } from '@tanstack/react-query';
import cityService from '@/services/cityService';
import type { CityFilter } from '@/types/filter';
import type { City } from '@/types/models/city';

interface UseCitiesReturn {
  cities: City[];
  pagination: {
    total: number;
    has_more: boolean;
    current_page: number;
    total_pages: number;
    page_size: number;
  } | null;
  isLoading: boolean;
  isError: boolean;
  error: Error | null;
  refetch: () => void;
}

export const useCities = (
  filters: CityFilter = {},
  options: { enabled?: boolean } = {}
): UseCitiesReturn => {
  const {
    data: citiesData,
    isLoading,
    isError,
    error,
    refetch
  } = useQuery({
    queryKey: ['public-cities', filters],
    queryFn: async () => {
      const response = await cityService.getPublicCities(filters);
      return response;
    },
    enabled: options.enabled !== false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  const cities = citiesData?.success ? citiesData.data.cities : [];
  const pagination = citiesData?.success ? citiesData.data.pagination : null;

  return {
    cities,
    pagination,
    isLoading,
    isError,
    error: error as Error | null,
    refetch,
  };
};

export default useCities;