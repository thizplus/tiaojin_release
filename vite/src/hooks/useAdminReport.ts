import { useQuery } from '@tanstack/react-query';
import reportService from '@/services/reportService';
import type { AdminReportParams, AdminReportData } from '@/types/report';

interface UseAdminReportReturn {
  data: AdminReportData | null;
  isLoading: boolean;
  isError: boolean;
  error: Error | null;
  refetch: () => void;
}

export const useAdminReport = (
  params: AdminReportParams,
  options: { enabled?: boolean } = {}
): UseAdminReportReturn => {
  const {
    data: reportData,
    isLoading,
    isError,
    error,
    refetch
  } = useQuery({
    queryKey: ['admin-report', params.year, params.month],
    queryFn: async () => {
      const response = await reportService.getAdminReport(params);
      return response;
    },
    enabled: options.enabled !== false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  const data = reportData?.success ? reportData.data : null;

  return {
    data,
    isLoading,
    isError,
    error: error as Error | null,
    refetch,
  };
};

export default useAdminReport;