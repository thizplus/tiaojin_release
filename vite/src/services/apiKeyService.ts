// src/services/apiKeyService.ts
import { API } from '@/constants/api'
import apiService from './apiService'
import type {
  CreateApiKeyRequest,
  UpdateApiKeyRequest,
  ApiKeyListRequest,
  ValidateApiKeyRequest
} from '@/types/request/apiKey'
import type {
  ApiKeyListResponse,
  ApiKeyResponse,
  CreateApiKeyResponse,
  UpdateApiKeyResponse,
  DeleteApiKeyResponse,
  ValidateApiKeyResponse
} from '@/types/response/apiKey'

/**
 * Service สำหรับจัดการ API Keys
 */
const apiKeyService = {
  /**
   * ดึงรายการ API Keys
   * @param params - พารามิเตอร์สำหรับการค้นหาและ pagination
   * @returns รายการ API Keys
   */
  getAll: async (params: ApiKeyListRequest = {}): Promise<ApiKeyListResponse> => {
    return await apiService.get<ApiKeyListResponse>(API.API_KEY.LIST, { params })
  },

  /**
   * ดึงข้อมูล API Key ตาม ID
   * @param id - API Key ID
   * @returns ข้อมูล API Key
   */
  getById: async (id: string): Promise<ApiKeyResponse> => {
    return await apiService.get<ApiKeyResponse>(API.API_KEY.GET_BY_ID(id))
  },

  /**
   * สร้าง API Key ใหม่
   * @param data - ข้อมูลสำหรับสร้าง API Key
   * @returns ข้อมูล API Key ที่สร้างใหม่ (รวม api_key เต็ม)
   */
  create: async (data: CreateApiKeyRequest): Promise<CreateApiKeyResponse> => {
    return await apiService.post<CreateApiKeyResponse>(API.API_KEY.CREATE, data)
  },

  /**
   * แก้ไข API Key
   * @param id - API Key ID
   * @param data - ข้อมูลที่ต้องการแก้ไข
   * @returns ข้อมูล API Key ที่อัปเดตแล้ว
   */
  update: async (id: string, data: UpdateApiKeyRequest): Promise<UpdateApiKeyResponse> => {
    return await apiService.put<UpdateApiKeyResponse>(API.API_KEY.UPDATE(id), data)
  },

  /**
   * ลบ API Key (Soft Delete/Revoke)
   * @param id - API Key ID
   * @returns ผลการลบ
   */
  delete: async (id: string): Promise<DeleteApiKeyResponse> => {
    return await apiService.delete<DeleteApiKeyResponse>(API.API_KEY.DELETE(id))
  },

  /**
   * ตรวจสอบ API Key
   * @param data - API Key ที่ต้องการตรวจสอบ
   * @returns ผลการตรวจสอบ
   */
  validate: async (data: ValidateApiKeyRequest): Promise<ValidateApiKeyResponse> => {
    return await apiService.post<ValidateApiKeyResponse>(API.API_KEY.VALIDATE, data)
  }
}

export default apiKeyService