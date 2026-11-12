import { API } from '@/constants/api';
import apiService from './apiService';
import type {
  ConfigApiResponse,
  ConfigApiInfo,
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
  CreateConfigResponse,
  UpdateConfigRequest,
  BookingLimitsConfig,
  BookingLimitsRequest,
} from '@/types/config';

/**
 * Service สำหรับจัดการการตั้งค่าระบบ
 */
const configService = {
  // ==================== Base Info ====================

  /**
   * ดึงข้อมูลพื้นฐานของ Config API
   * @returns ข้อมูล API info
   */
  getConfigInfo: async (): Promise<ConfigApiResponse<ConfigApiInfo>> => {
    return await apiService.get<ConfigApiResponse<ConfigApiInfo>>(API.CONFIG.INFO);
  },

  // ==================== Expiration Tiers ====================

  /**
   * ดึงการตั้งค่า expiration tiers
   * @returns ข้อมูล expiration tiers
   */
  getExpirationTiers: async (): Promise<ConfigApiResponse<ExpirationTiersData>> => {
    return await apiService.get<ConfigApiResponse<ExpirationTiersData>>(API.CONFIG.EXPIRATION_TIERS);
  },

  /**
   * อัพเดตการตั้งค่า expiration tiers
   * @param data - ข้อมูล tiers ใหม่
   * @returns ผลลัพธ์การอัพเดต
   */
  updateExpirationTiers: async (data: ExpirationTiersRequest): Promise<ConfigApiResponse<null>> => {
    return await apiService.put<ConfigApiResponse<null>>(API.CONFIG.EXPIRATION_TIERS, data);
  },

  /**
   * ตรวจสอบความถูกต้องของการตั้งค่า expiration tiers
   * @param data - ข้อมูล tiers ที่ต้องการตรวจสอบ
   * @returns ผลการตรวจสอบ
   */
  validateExpirationTiers: async (data: ExpirationTiersRequest): Promise<ConfigApiResponse<ExpirationValidationResult>> => {
    return await apiService.post<ConfigApiResponse<ExpirationValidationResult>>(API.CONFIG.VALIDATE_TIERS, data);
  },

  /**
   * คำนวณ tier ที่เหมาะสมสำหรับการจอง
   * @param data - ข้อมูลวันที่สำหรับคำนวณ
   * @returns ผลการคำนวณ tier
   */
  calculateTier: async (data: TierCalculationRequest): Promise<ConfigApiResponse<TierCalculationResult>> => {
    return await apiService.post<ConfigApiResponse<TierCalculationResult>>(API.CONFIG.CALCULATE_TIER, data);
  },

  /**
   * ดึงสถิติการหมดอายุ
   * @returns สถิติการหมดอายุ
   */
  getExpirationStats: async (): Promise<ConfigApiResponse<ExpirationStats>> => {
    return await apiService.get<ConfigApiResponse<ExpirationStats>>(API.CONFIG.EXPIRATION_STATS);
  },

  /**
   * ประมวลผลการจองที่หมดอายุ
   * @param data - พารามิเตอร์สำหรับการประมวลผล
   * @returns ผลการประมวลผล
   */
  processExpiredBookings: async (data: ProcessExpiredRequest = {}): Promise<ConfigApiResponse<ProcessExpiredResult>> => {
    return await apiService.post<ConfigApiResponse<ProcessExpiredResult>>(API.CONFIG.PROCESS_EXPIRED, data);
  },

  // ==================== General Config Management ====================

  /**
   * ดึงรายการการตั้งค่า
   * @param params - พารามิเตอร์สำหรับการกรองและแบ่งหน้า
   * @returns รายการการตั้งค่า
   */
  getConfigs: async (params: ConfigListRequest = {}): Promise<ConfigApiResponse<ConfigListResponse>> => {
    return await apiService.get<ConfigApiResponse<ConfigListResponse>>(API.CONFIG.CONFIGS, params);
  },

  /**
   * สร้างการตั้งค่าใหม่
   * @param data - ข้อมูลการตั้งค่าใหม่
   * @returns ข้อมูลการตั้งค่าที่สร้างขึ้น
   */
  createConfig: async (data: CreateConfigRequest): Promise<ConfigApiResponse<CreateConfigResponse>> => {
    return await apiService.post<ConfigApiResponse<CreateConfigResponse>>(API.CONFIG.CONFIGS, data);
  },

  /**
   * ดึงการตั้งค่าตาม ID
   * @param id - ID ของการตั้งค่า
   * @returns ข้อมูลการตั้งค่า
   */
  getConfigById: async (id: string): Promise<ConfigApiResponse<SystemConfig>> => {
    return await apiService.get<ConfigApiResponse<SystemConfig>>(API.CONFIG.CONFIG_BY_ID(id));
  },

  /**
   * อัพเดตการตั้งค่า
   * @param id - ID ของการตั้งค่า
   * @param data - ข้อมูลที่ต้องการอัพเดต
   * @returns ผลลัพธ์การอัพเดต
   */
  updateConfig: async (id: string, data: UpdateConfigRequest): Promise<ConfigApiResponse<null>> => {
    return await apiService.put<ConfigApiResponse<null>>(API.CONFIG.CONFIG_BY_ID(id), data);
  },

  /**
   * ลบการตั้งค่า
   * @param id - ID ของการตั้งค่า
   * @returns ผลลัพธ์การลบ
   */
  deleteConfig: async (id: string): Promise<ConfigApiResponse<null>> => {
    return await apiService.delete<ConfigApiResponse<null>>(API.CONFIG.CONFIG_BY_ID(id));
  },

  // ==================== Booking Limits Management ====================

  /**
   * ดึงการตั้งค่า booking limits
   * @returns ข้อมูล booking limits config
   */
  getBookingLimitsConfig: async (): Promise<BookingLimitsConfig | null> => {
    try {
      console.log('🔍 getBookingLimitsConfig: Fetching from dedicated endpoint...');
      const response = await apiService.get<ConfigApiResponse<BookingLimitsConfig>>(API.CONFIG.AGENT_BOOKING_LIMITS);
      console.log('🔍 getBookingLimitsConfig: API response:', response);

      if (!response.success || !response.data) {
        console.log('🔍 getBookingLimitsConfig: No config found or API error, returning null');
        return null;
      }

      console.log('🔍 getBookingLimitsConfig: Processed result:', response.data);
      return response.data;
    } catch (error) {
      console.error('❌ Error getting booking limits config:', error);
      return null;
    }
  },

  /**
   * อัพเดตการตั้งค่า booking limits
   * @param data - ข้อมูล booking limits ใหม่
   * @returns ผลลัพธ์การอัพเดต
   */
  updateBookingLimitsConfig: async (data: BookingLimitsRequest): Promise<ConfigApiResponse<null>> => {
    console.log('🔍 updateBookingLimitsConfig: Updating with data:', data);
    return await apiService.put<ConfigApiResponse<null>>(API.CONFIG.AGENT_BOOKING_LIMITS, data);
  },

  // ==================== Helper Methods ====================

  /**
   * ดึงการตั้งค่าตาม key
   * @param configKey - key ของการตั้งค่า
   * @returns ข้อมูลการตั้งค่า หรือ null ถ้าไม่พบ
   */
  getConfigByKey: async (configKey: string): Promise<SystemConfig | null> => {
    try {
      const response = await configService.getConfigs({ limit: 100 });
      if (!response.success || !response.data) {
        return null;
      }

      const config = response.data.configs.find(c => c.config_key === configKey);
      return config || null;
    } catch (error) {
      console.error('Error getting config by key:', error);
      return null;
    }
  },

  /**
   * ดึงการตั้งค่าตามประเภท
   * @param configType - ประเภทของการตั้งค่า
   * @returns รายการการตั้งค่า
   */
  getConfigsByType: async (configType: string): Promise<SystemConfig[]> => {
    try {
      const response = await configService.getConfigs({ config_type: configType, limit: 100 });
      if (!response.success || !response.data) {
        return [];
      }

      return response.data.configs;
    } catch (error) {
      console.error('Error getting configs by type:', error);
      return [];
    }
  },

  /**
   * ประมวลผลการจองที่หมดอายุแบบปลอดภัย (ทดสอบก่อน)
   * @param limit - จำนวนรายการที่ต้องการประมวลผล
   * @returns ผลการประมวลผลทั้งทดสอบและจริง
   */
  processExpiredBookingsSafe: async (limit?: number): Promise<{
    dryRun: ProcessExpiredResult;
    actual?: ProcessExpiredResult;
  }> => {
    // ทำการทดสอบก่อน
    const dryRunResponse = await configService.processExpiredBookings({
      dry_run: true,
      limit,
    });

    if (!dryRunResponse.success || !dryRunResponse.data) {
      throw new Error(dryRunResponse.message || 'Dry run failed');
    }

    const result: {
      dryRun: ProcessExpiredResult;
      actual?: ProcessExpiredResult;
    } = {
      dryRun: dryRunResponse.data,
    };

    // ถ้าพบการจองที่หมดอายุ ให้ประมวลผลจริง
    if (dryRunResponse.data.expired_found > 0) {
      const actualResponse = await configService.processExpiredBookings({
        dry_run: false,
        limit,
      });

      if (actualResponse.success && actualResponse.data) {
        result.actual = actualResponse.data;
      }
    }

    return result;
  }
};

export default configService;