import { useQuery } from '@tanstack/react-query';
import apiService from '@/services/apiService';
import { COMPANY_PROFILE_API } from '@/constants/api';

interface LogosResponse {
  success: boolean;
  data: { logos: string[] };
}

export function useAgentLogos(limit = 20) {
  const { data, isLoading } = useQuery({
    queryKey: ['agent-logos', limit],
    queryFn: () => apiService.get<LogosResponse>(`${COMPANY_PROFILE_API.RANDOM_LOGOS}?limit=${limit}`),
    staleTime: 60 * 1000, // 1 นาที แล้ว refetch ใหม่ (random ใหม่)
  });

  return {
    logos: data?.success ? data.data.logos : [],
    isLoading,
  };
}
