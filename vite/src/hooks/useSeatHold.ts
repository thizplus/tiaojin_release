import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import type {
  CreateSeatHoldRequest,
  GetPeriodSeatHoldsRequest,
  ReleaseSeatHoldRequest,
} from '@/types/request/seat-hold';
import { seatHoldService } from '@/services/admin/seat-hold';



// Query Keys
export const seatHoldKeys = {
  all: ['seat-holds'] as const,
  lists: () => [...seatHoldKeys.all, 'list'] as const,
  list: (periodId: string, params?: GetPeriodSeatHoldsRequest) =>
    [...seatHoldKeys.lists(), periodId, params] as const,
  details: () => [...seatHoldKeys.all, 'detail'] as const,
  detail: (holdId: string) => [...seatHoldKeys.details(), holdId] as const,
};

// Hook สำหรับ GET period seat holds
export function usePeriodSeatHolds(
  periodId: string,
  params?: GetPeriodSeatHoldsRequest,
  options?: {
    enabled?: boolean;
    refetchInterval?: number;
  }
) {
  return useQuery({
    queryKey: seatHoldKeys.list(periodId, params),
    queryFn: () => seatHoldService.getPeriodSeatHolds(periodId, params),
    enabled: !!periodId && options?.enabled !== false,
    refetchInterval: options?.refetchInterval,
    staleTime: 1000 * 30, // 30 seconds
  });
}

// Hook สำหรับ CREATE seat hold
export function useCreateSeatHold(periodId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSeatHoldRequest) =>
      seatHoldService.createSeatHold(periodId, data),
    onSuccess: () => {
      // Invalidate and refetch period seat holds
      queryClient.invalidateQueries({
        queryKey: seatHoldKeys.lists(),
      });

      // Invalidate period data เพื่อ update available seats
      queryClient.invalidateQueries({
        queryKey: ['period', 'list'],
      });

      // Invalidate specific period if we have the periodId
      queryClient.invalidateQueries({
        queryKey: ['period', 'byId', periodId],
      });
    },
  });
}

// Hook สำหรับ RELEASE seat hold
export function useReleaseSeatHold() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      holdId,
      data,
    }: {
      holdId: string;
      data?: ReleaseSeatHoldRequest;
    }) => seatHoldService.releaseSeatHold(holdId, data),
    onSuccess: () => {
      // Invalidate and refetch all seat holds
      queryClient.invalidateQueries({
        queryKey: seatHoldKeys.all,
      });

      // Invalidate period data เพื่อ update available seats
      queryClient.invalidateQueries({
        queryKey: ['period', 'list'],
      });

      // Invalidate all individual period queries since we don't know which period this hold belongs to
      queryClient.invalidateQueries({
        queryKey: ['period', 'byId'],
      });
    },
  });
}

// Hook สำหรับ RELEASE seat hold ด้วย query parameter
export function useReleaseSeatHoldWithQuery() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      holdId,
      reason,
    }: {
      holdId: string;
      reason?: string;
    }) => seatHoldService.releaseSeatHoldWithQueryParam(holdId, reason),
    onSuccess: () => {
      // Invalidate and refetch all seat holds
      queryClient.invalidateQueries({
        queryKey: seatHoldKeys.all,
      });

      // Invalidate period data เพื่อ update available seats
      queryClient.invalidateQueries({
        queryKey: ['period', 'list'],
      });

      // Invalidate all individual period queries since we don't know which period this hold belongs to
      queryClient.invalidateQueries({
        queryKey: ['period', 'byId'],
      });
    },
  });
}