import { useQuery } from '@tanstack/react-query';
import reportService from '@/services/reportService';
import type { AgentReportParams, AgentReportData } from '@/types/report';

interface UseAgentReportReturn {
  data: AgentReportData | null;
  isLoading: boolean;
  isError: boolean;
  error: Error | null;
  refetch: () => void;
}

export const useAgentReport = (
  params: AgentReportParams,
  options: { enabled?: boolean } = {}
): UseAgentReportReturn => {
  const {
    data: reportData,
    isLoading,
    isError,
    error,
    refetch
  } = useQuery({
    queryKey: ['agent-report', params.year, params.month],
    queryFn: async () => {
      const response = await reportService.getAgentReport(params);
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

export default useAgentReport;