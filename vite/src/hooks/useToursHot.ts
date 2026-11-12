import { useMemo, useEffect } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { toursHotService } from '@/services/tours-hot';
import { useToursHotStore as useToursHotStoreOriginal, toursHotKeys } from '@/stores/toursHotStore';
import { useWebSocket } from '@/hooks/useWebSocketHooks';
import type {
  Tour,
  ToursRegularRequest,
  Period,
} from '@/types/tours-regular';

/**
 * Hook สำหรับดึงข้อมูล tours hot พร้อม pagination
 */
export function useToursHot(
  params?: ToursRegularRequest,
  options?: {
    enabled?: boolean;
    refetchInterval?: number;
  }
) {
  const { setToursList, setPaginationMeta } = useToursHotStoreOriginal();

  const query = useQuery({
    queryKey: toursHotKeys.list(params),
    queryFn: async () => {
      const response = await toursHotService.getHotTours(params);

      if (response.success && response.data) {
        setToursList(response.data.tours);
        setPaginationMeta(response.data.pagination);
      }

      return response;
    },
    enabled: options?.enabled !== false,
    refetchInterval: options?.refetchInterval,
    staleTime: 1000 * 60, // 1 minute
    gcTime: 5 * 60 * 1000, // 5 minutes
  });

  return {
    ...query,
    tours: query.data?.data?.tours || [],
    pagination: query.data?.data?.pagination || null,
  };
}

/**
 * Hook สำหรับค้นหา tours ด้วย search text
 */
export function useToursHotSearch(
  searchText: string,
  baseParams?: Omit<ToursRegularRequest, 'search'>,
  options?: { enabled?: boolean; debounceMs?: number }
) {
  const searchParams = useMemo(() => ({
    ...baseParams,
    search: searchText?.trim(),
  }), [searchText, baseParams]);

  return useToursHot(searchParams, {
    enabled: options?.enabled !== false && Boolean(searchText?.trim()),
    refetchInterval: undefined,
  });
}

/**
 * Hook สำหรับ filter tours ตาม status
 */
export function useToursHotByStatus(
  status: 'publish' | 'draft' | 'all',
  baseParams?: Omit<ToursRegularRequest, 'status'>,
  options?: { enabled?: boolean }
) {
  const filterParams = useMemo(() => ({
    ...baseParams,
    ...(status !== 'all' && { status }),
  }), [status, baseParams]);

  return useToursHot(filterParams, options);
}

/**
 * Hook สำหรับ filter tours ที่มี available periods เท่านั้น
 */
export function useAvailableToursHot(
  baseParams?: Omit<ToursRegularRequest, 'available_only'>,
  options?: { enabled?: boolean }
) {
  const availableParams = useMemo(() => ({
    ...baseParams,
    available_only: true,
  }), [baseParams]);

  return useToursHot(availableParams, options);
}

/**
 * Hook สำหรับ filter tours ตาม city
 */
export function useToursHotByCity(
  cityId: string | null,
  baseParams?: Omit<ToursRegularRequest, 'city_id'>,
  options?: { enabled?: boolean }
) {
  const cityParams = useMemo(() => ({
    ...baseParams,
    ...(cityId && { city_id: cityId }),
  }), [cityId, baseParams]);

  return useToursHot(cityParams, {
    ...options,
    enabled: options?.enabled !== false && Boolean(cityId),
  });
}

/**
 * Hook สำหรับการจัดการ store data
 */
export function useToursHotStore() {
  const store = useToursHotStoreOriginal();
  return store;
}

/**
 * Hook สำหรับ invalidate และ refresh tours data
 */
export function useRefreshToursHot() {
  const queryClient = useQueryClient();

  return {
    refresh: () => {
      queryClient.invalidateQueries({
        queryKey: toursHotKeys.lists(),
      });
    },
    refetchActive: () => {
      queryClient.refetchQueries({
        queryKey: toursHotKeys.lists(),
        type: 'active',
      });
    },
    reset: () => {
      queryClient.removeQueries({
        queryKey: toursHotKeys.all,
      });
    }
  };
}

/**
 * Hook สำหรับ auto-refresh tours เมื่อได้รับ WebSocket events
 */
export function useToursHotWithWebSocket(
  params?: ToursRegularRequest,
  options?: {
    enabled?: boolean;
    refetchInterval?: number;
    autoRefreshOnWebSocketEvents?: boolean;
  }
) {
  const queryClient = useQueryClient();
  const { subscribeSeatUpdates } = (useWebSocket() as unknown) as {
    subscribeSeatUpdates: (id: string, callback: (data: unknown) => void) => () => void;
  };

  // ใช้ hook หลักสำหรับดึงข้อมูล
  const toursQuery = useToursHot(params, {
    enabled: options?.enabled,
    refetchInterval: options?.refetchInterval,
  });

  // Subscribe WebSocket events เพื่อ auto-refresh
  useEffect(() => {
    if (options?.autoRefreshOnWebSocketEvents === false) return;

    const unsubscribe = subscribeSeatUpdates(`tours-hot-refresh-${Math.random()}`, () => {
      // Auto-refresh เมื่อมี seat updates
      queryClient.invalidateQueries({
        queryKey: toursHotKeys.lists(),
      });
    });

    return unsubscribe;
  }, [subscribeSeatUpdates, queryClient, options?.autoRefreshOnWebSocketEvents]);

  return toursQuery;
}

/**
 * Utility hook สำหรับหา tour ที่มี available periods
 */
export function useAvailablePeriodsHot(tour: Tour | null): Period[] {
  return useMemo(() => {
    if (!tour?.available_periods) return [];

    return tour.available_periods.filter(period =>
      period.status === 'publish' &&
      period.available_seats > 0
    ).sort((a, b) =>
      new Date(a.start_date).getTime() - new Date(b.start_date).getTime()
    );
  }, [tour]);
}

/**
 * Utility hook สำหรับคำนวณ summary data
 */
export function useToursHotSummary(tours: Tour[]) {
  return useMemo(() => {
    const publishedCount = tours.filter(tour => tour.status === 'publish').length;
    const draftCount = tours.filter(tour => tour.status === 'draft').length;
    const withAvailablePeriodsCount = tours.filter(tour =>
      tour.available_periods.some(period =>
        period.status === 'publish' && period.available_seats > 0
      )
    ).length;

    const totalAvailableSeats = tours.reduce((total, tour) =>
      total + tour.available_periods.reduce((tourTotal, period) =>
        tourTotal + (period.status === 'publish' ? period.available_seats : 0), 0
      ), 0
    );

    return {
      total: tours.length,
      published: publishedCount,
      draft: draftCount,
      withAvailablePeriods: withAvailablePeriodsCount,
      totalAvailableSeats,
    };
  }, [tours]);
}