import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult,
  type UseMutationResult
} from '@tanstack/react-query';
import configService from '@/services/config';
import type {
  ExpirationTiersData,
  ExpirationTiersRequest,
  ExpirationValidationResult,
  TierCalculationRequest,
  TierCalculationResult,
  ExpirationStats,
  ProcessExpiredRequest,
  ProcessExpiredResult,
  ConfigListRequest,
  ConfigListResponse,
  SystemConfig,
  CreateConfigRequest,
  UpdateConfigRequest,
  ConfigApiInfo,
  BookingLimitsConfig,
  BookingLimitsRequest,
} from '@/types/config';

// ==================== Query Keys ====================
export const configKeys = {
  all: ['config'] as const,
  info: () => [...configKeys.all, 'info'] as const,
  expirationTiers: () => [...configKeys.all, 'expiration-tiers'] as const,
  expirationStats: () => [...configKeys.all, 'expiration-stats'] as const,
  bookingLimits: () => [...configKeys.all, 'booking-limits'] as const,
  configs: () => [...configKeys.all, 'configs'] as const,
  configList: (params?: ConfigListRequest) => [...configKeys.configs(), 'list', params] as const,
  config: (id: string) => [...configKeys.configs(), 'detail', id] as const,
  configsByType: (type: string) => [...configKeys.configs(), 'by-type', type] as const,
  configByKey: (key: string) => [...configKeys.configs(), 'by-key', key] as const,
};

// ==================== Zustand Store ====================
interface ConfigState {
  // Current data
  currentConfig: SystemConfig | null;
  selectedConfigs: string[];
  expirationTiers: ExpirationTiersData | null;
  bookingLimitsConfig: BookingLimitsConfig | null;

  // UI state
  isExpirationTiersDialogOpen: boolean;
  isConfigDialogOpen: boolean;
  configDialogMode: 'create' | 'edit' | 'view';

  // Actions
  setCurrentConfig: (config: SystemConfig | null) => void;
  setSelectedConfigs: (ids: string[]) => void;
  setExpirationTiers: (data: ExpirationTiersData | null) => void;
  setBookingLimitsConfig: (data: BookingLimitsConfig | null) => void;
  setExpirationTiersDialogOpen: (open: boolean) => void;
  setConfigDialogOpen: (open: boolean, mode?: 'create' | 'edit' | 'view') => void;
  clearConfigData: () => void;
}

export const useConfigStore = create<ConfigState>()(
  persist(
    (set) => ({
      // Initial state
      currentConfig: null,
      selectedConfigs: [],
      expirationTiers: null,
      bookingLimitsConfig: null,
      isExpirationTiersDialogOpen: false,
      isConfigDialogOpen: false,
      configDialogMode: 'create',

      // Actions
      setCurrentConfig: (config) => set({ currentConfig: config }),
      setSelectedConfigs: (ids) => set({ selectedConfigs: ids }),
      setExpirationTiers: (data) => set({ expirationTiers: data }),
      setBookingLimitsConfig: (data) => set({ bookingLimitsConfig: data }),
      setExpirationTiersDialogOpen: (open) => set({ isExpirationTiersDialogOpen: open }),
      setConfigDialogOpen: (open, mode = 'create') => set({
        isConfigDialogOpen: open,
        configDialogMode: mode
      }),
      clearConfigData: () => set({
        currentConfig: null,
        selectedConfigs: [],
        expirationTiers: null,
        bookingLimitsConfig: null,
        isExpirationTiersDialogOpen: false,
        isConfigDialogOpen: false,
        configDialogMode: 'create',
      }),
    }),
    {
      name: 'config-storage',
      partialize: (state) => ({
        selectedConfigs: state.selectedConfigs,
      }),
    }
  )
);

// ==================== React Query Hooks ====================

// Config Info
export const useConfigInfo = (): UseQueryResult<ConfigApiInfo, Error> => {
  return useQuery({
    queryKey: configKeys.info(),
    queryFn: async () => {
      const response = await configService.getConfigInfo();
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to fetch config info');
      }
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  });
};

// Expiration Tiers
export const useExpirationTiers = (): UseQueryResult<ExpirationTiersData, Error> => {
  const { setExpirationTiers } = useConfigStore();

  return useQuery({
    queryKey: configKeys.expirationTiers(),
    queryFn: async () => {
      const response = await configService.getExpirationTiers();
      console.log('🔍 Raw API response:', response); // Debug log

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to fetch expiration tiers');
      }

      console.log('✅ Processed tiers data:', response.data); // Debug log
      setExpirationTiers(response.data);
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useUpdateExpirationTiers = (): UseMutationResult<void, Error, ExpirationTiersRequest> => {
  const queryClient = useQueryClient();
  //const { setExpirationTiers } = useConfigStore();

  return useMutation({
    mutationFn: async (data: ExpirationTiersRequest) => {
      const response = await configService.updateExpirationTiers(data);
      if (!response.success) {
        throw new Error(response.message || 'Failed to update expiration tiers');
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.expirationTiers() });
      queryClient.invalidateQueries({ queryKey: configKeys.expirationStats() });
    },
  });
};

export const useValidateExpirationTiers = (): UseMutationResult<ExpirationValidationResult, Error, ExpirationTiersRequest> => {
  return useMutation({
    mutationFn: async (data: ExpirationTiersRequest) => {
      const response = await configService.validateExpirationTiers(data);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to validate expiration tiers');
      }
      return response.data;
    },
  });
};

export const useCalculateTier = (): UseMutationResult<TierCalculationResult, Error, TierCalculationRequest> => {
  return useMutation({
    mutationFn: async (data: TierCalculationRequest) => {
      const response = await configService.calculateTier(data);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to calculate tier');
      }
      return response.data;
    },
  });
};

// Expiration Stats
export const useExpirationStats = (): UseQueryResult<ExpirationStats, Error> => {
  return useQuery({
    queryKey: configKeys.expirationStats(),
    queryFn: async () => {
      const response = await configService.getExpirationStats();
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to fetch expiration stats');
      }
      return response.data;
    },
    staleTime: 1 * 60 * 1000, // 1 minute
    gcTime: 3 * 60 * 1000, // 3 minutes
    refetchInterval: 2 * 60 * 1000, // Auto-refresh every 2 minutes
  });
};

export const useProcessExpiredBookings = (): UseMutationResult<ProcessExpiredResult, Error, ProcessExpiredRequest> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: ProcessExpiredRequest) => {
      const response = await configService.processExpiredBookings(data);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to process expired bookings');
      }
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.expirationStats() });
    },
  });
};

// General Config Management
export const useConfigs = (params: ConfigListRequest = {}): UseQueryResult<ConfigListResponse, Error> => {
  return useQuery({
    queryKey: configKeys.configList(params),
    queryFn: async () => {
      const response = await configService.getConfigs(params);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to fetch configs');
      }
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useConfig = (id: string): UseQueryResult<SystemConfig, Error> => {
  const { setCurrentConfig } = useConfigStore();

  return useQuery({
    queryKey: configKeys.config(id),
    queryFn: async () => {
      const response = await configService.getConfigById(id);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Failed to fetch config');
      }

      setCurrentConfig(response.data);
      return response.data;
    },
    enabled: !!id,
    staleTime: 1 * 60 * 1000, // 1 minute
    gcTime: 3 * 60 * 1000, // 3 minutes
  });
};

export const useCreateConfig = (): UseMutationResult<string | null, Error, CreateConfigRequest> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateConfigRequest) => {
      const response = await configService.createConfig(data);
      if (!response.success) {
        throw new Error(response.message || 'Failed to create config');
      }
      // API อาจไม่ return data กลับมา แต่ success = true ก็ถือว่าสำเร็จ
      return response.data?.id || null;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.configs() });
    },
  });
};

export const useUpdateConfig = (): UseMutationResult<void, Error, { id: string; data: UpdateConfigRequest }> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await configService.updateConfig(id, data);
      if (!response.success) {
        throw new Error(response.message || 'Failed to update config');
      }
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: configKeys.config(id) });
      queryClient.invalidateQueries({ queryKey: configKeys.configs() });
    },
  });
};

export const useDeleteConfig = (): UseMutationResult<void, Error, string> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await configService.deleteConfig(id);
      if (!response.success) {
        throw new Error(response.message || 'Failed to delete config');
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.configs() });
    },
  });
};

// Booking Limits Management
export const useBookingLimitsConfig = (): UseQueryResult<BookingLimitsConfig | null, Error> => {
  const { setBookingLimitsConfig } = useConfigStore();

  return useQuery({
    queryKey: configKeys.bookingLimits(),
    queryFn: async () => {
      const data = await configService.getBookingLimitsConfig();
      setBookingLimitsConfig(data);
      return data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useUpdateBookingLimitsConfig = (): UseMutationResult<void, Error, BookingLimitsRequest> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: BookingLimitsRequest) => {
      const response = await configService.updateBookingLimitsConfig(data);
      if (!response.success) {
        throw new Error(response.message || 'Failed to update booking limits config');
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.bookingLimits() });
      queryClient.invalidateQueries({ queryKey: configKeys.configs() });
    },
  });
};

// Helper Hooks
export const useConfigsByType = (configType: string): UseQueryResult<SystemConfig[], Error> => {
  return useQuery({
    queryKey: configKeys.configsByType(configType),
    queryFn: async () => {
      return configService.getConfigsByType(configType);
    },
    enabled: !!configType,
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useConfigByKey = (configKey: string): UseQueryResult<SystemConfig | null, Error> => {
  return useQuery({
    queryKey: configKeys.configByKey(configKey),
    queryFn: async () => {
      return configService.getConfigByKey(configKey);
    },
    enabled: !!configKey,
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

export default useConfigStore;