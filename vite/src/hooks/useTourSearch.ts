import { useState, useEffect, useCallback } from 'react';
import tourService from '@/services/admin/tourService';
import type { Tour } from '@/types';

interface UseTourSearchOptions {
  initialSearch?: string;
  limit?: number;
  debounceMs?: number;
  autoSearch?: boolean;
}

export function useTourSearch(options: UseTourSearchOptions = {}) {
  const {
    initialSearch = '',
    limit = 10,
    debounceMs = 300,
    autoSearch = true,
  } = options;

  const [tours, setTours] = useState<Tour[]>([]);
  const [searchTerm, setSearchTerm] = useState(initialSearch);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const searchTours = useCallback(async (query: string) => {
    if (!query.trim()) {
      setTours([]);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const response = await tourService.getTours({
        search: query,
        limit,
        status: 'publish',
      });

      if (response.success && response.data.tours) {
        setTours(response.data.tours);
      } else {
        setTours([]);
      }
    } catch (err) {
      console.error('Failed to search tours:', err);
      setError('ไม่สามารถค้นหาทัวร์ได้');
      setTours([]);
    } finally {
      setIsLoading(false);
    }
  }, [limit]);

  // Debounced search
  useEffect(() => {
    if (!autoSearch) return;

    const timer = setTimeout(() => {
      if (searchTerm.trim()) {
        searchTours(searchTerm);
      } else {
        setTours([]);
      }
    }, debounceMs);

    return () => clearTimeout(timer);
  }, [searchTerm, searchTours, debounceMs, autoSearch]);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
    setTours([]);
    setError(null);
  }, []);

  return {
    tours,
    searchTerm,
    setSearchTerm,
    isLoading,
    error,
    searchTours,
    clearSearch,
  };
}