import { useState, useEffect, useCallback } from 'react';
import { periodService } from '@/services/admin/periodService';
import type { TourPeriod } from '@/types/models/period';

interface UsePeriodSearchOptions {
  tourId: string | null;
  initialSearch?: string;
  limit?: number;
  debounceMs?: number;
  autoSearch?: boolean;
}

export function usePeriodSearch(options: UsePeriodSearchOptions) {
  const {
    tourId,
    initialSearch = '',
    limit = 20,
    debounceMs = 300,
    autoSearch = true,
  } = options;

  const [periods, setPeriods] = useState<TourPeriod[]>([]);
  const [searchTerm, setSearchTerm] = useState(initialSearch);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const searchPeriods = useCallback(async (query?: string) => {
    if (!tourId) {
      setPeriods([]);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const response = await periodService.getByTourId(tourId, {
        search: query?.trim() || undefined,
        limit,
        status: 'publish',
        sort: 'start_date',
        order: 'asc',
      });

      if (response.success && response.data.periods) {
        setPeriods(response.data.periods);
      } else {
        setPeriods([]);
      }
    } catch (err) {
      console.error('Failed to search periods:', err);
      setError('ไม่สามารถค้นหารอบทัวร์ได้');
      setPeriods([]);
    } finally {
      setIsLoading(false);
    }
  }, [tourId, limit]);

  // Debounced search
  useEffect(() => {
    if (!autoSearch || !tourId) {
      if (!tourId) setPeriods([]);
      return;
    }

    const timer = setTimeout(() => {
      searchPeriods(searchTerm);
    }, debounceMs);

    return () => clearTimeout(timer);
  }, [tourId, searchTerm, searchPeriods, debounceMs, autoSearch]);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
    setError(null);
  }, []);

  return {
    periods,
    searchTerm,
    setSearchTerm,
    isLoading,
    error,
    searchPeriods,
    clearSearch,
  };
}