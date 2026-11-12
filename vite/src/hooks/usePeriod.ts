import {
  usePeriodStore,
  useCreatePeriodMutation,
  useUpdatePeriodMutation,
  useDeletePeriodMutation,
  useUpdatePeriodStatusMutation,
  useCheckPeriodAvailabilityMutation
} from '@/stores/periodStore';
import type { PeriodStatus } from '@/types/models/period';
import type {
  CreatePeriodRequest,
  UpdatePeriodRequest,
} from '@/types/request/period';

export const usePeriod = () => {
  const {
    currentPeriod,
    selectedPeriod,
    periodsList,
    setCurrentPeriod,
    setSelectedPeriod,
    clearPeriodData,
  } = usePeriodStore();

  const createPeriodMutation = useCreatePeriodMutation();
  const updatePeriodMutation = useUpdatePeriodMutation();
  const deletePeriodMutation = useDeletePeriodMutation();
  const updateStatusMutation = useUpdatePeriodStatusMutation();
  const checkAvailabilityMutation = useCheckPeriodAvailabilityMutation();

  const handleCreatePeriod = async (data: CreatePeriodRequest) => {
    try {
      const newPeriod = await createPeriodMutation.mutateAsync(data);
      return { success: true, data: newPeriod };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดในการสร้าง period',
      };
    }
  };

  const handleUpdatePeriod = async (id: string, data: UpdatePeriodRequest) => {
    try {
      const updatedPeriod = await updatePeriodMutation.mutateAsync({ id, data });
      return { success: true, data: updatedPeriod };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดในการอัปเดต period',
      };
    }
  };

  const handleDeletePeriod = async (id: string) => {
    try {
      await deletePeriodMutation.mutateAsync(id);
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดในการลบ period',
      };
    }
  };

  const handleUpdateStatus = async (id: string, status: PeriodStatus) => {
    try {
      await updateStatusMutation.mutateAsync({ id, status });
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดในการอัปเดตสถานะ period',
      };
    }
  };

  const handleCheckAvailability = async (id: string, requestedSeats?: number) => {
    try {
      const availability = await checkAvailabilityMutation.mutateAsync({ id, requestedSeats });
      return { success: true, data: availability };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดในการตรวจสอบที่นั่งว่าง',
      };
    }
  };

  const getPeriodsByTourId = (tourId: string) => {
    return periodsList.filter(period => period.tour_id === tourId);
  };

  return {
    // State
    currentPeriod,
    selectedPeriod,
    periodsList,
    loading: createPeriodMutation.isPending || updatePeriodMutation.isPending ||
             deletePeriodMutation.isPending || updateStatusMutation.isPending ||
             checkAvailabilityMutation.isPending,
    error: createPeriodMutation.error?.message || updatePeriodMutation.error?.message ||
           deletePeriodMutation.error?.message || updateStatusMutation.error?.message ||
           checkAvailabilityMutation.error?.message || null,

    // Actions
    handleCreatePeriod,
    handleUpdatePeriod,
    handleDeletePeriod,
    handleUpdateStatus,
    handleCheckAvailability,

    // Store actions
    setCurrentPeriod,
    setSelectedPeriod,
    clearPeriodData,

    // Utility actions
    getPeriodsByTourId,
  };
};