import { useQuery } from '@tanstack/react-query';
import reportService from '@/services/reportService';
import type { SalesReportParams, SalesReportData } from '@/types/report';

interface UseSalesReportReturn {
  data: SalesReportData | null;
  isLoading: boolean;
  isError: boolean;
  error: Error | null;
  refetch: () => void;
}

export const useSalesReport = (
  params: SalesReportParams,
  options: { enabled?: boolean } = {}
): UseSalesReportReturn => {
  const {
    data: reportData,
    isLoading,
    isError,
    error,
    refetch
  } = useQuery({
    queryKey: ['sales-report', params.year, params.month],
    queryFn: async () => {
      const response = await reportService.getSalesReport(params);
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

export default useSalesReport;