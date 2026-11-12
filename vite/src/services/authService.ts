// src/services/authService.ts
import { API } from '@/constants/api';
import apiService from './apiService';

import type {
  RegisterRequest,
  LoginRequest,
  RefreshTokenRequest,
  ResetPasswordRequest,
  UpdatePasswordRequest,
  ApiResponse,
  LoginResponse,
  TokenValidationResponse,
  User
} from '@/types';

/**
 * Service สำหรับจัดการการยืนยันตัวตน
 */
const authService = {
  /**
   * ลงทะเบียนผู้ใช้ใหม่
   * @param data - ข้อมูลการลงทะเบียน
   * @returns ข้อมูลการลงทะเบียนและผู้ใช้
   */
  register: async (data: RegisterRequest): Promise<ApiResponse<User>> => {
    return await apiService.post<ApiResponse<User>>(API.AUTH.REGISTER, data);
  },

  /**
   * เข้าสู่ระบบ
   * @param data - ข้อมูลการเข้าสู่ระบบ
   * @returns ข้อมูลการเข้าสู่ระบบและ token
   */
  login: async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    return await apiService.post<ApiResponse<LoginResponse>>(API.AUTH.LOGIN, data);
  },

  /**
   * รีเซ็ตรหัสผ่าน
   * @param data - ข้อมูลสำหรับการรีเซ็ตรหัสผ่าน
   * @returns ผลลัพธ์จากการรีเซ็ตรหัสผ่าน
   */
  resetPassword: async (data: ResetPasswordRequest): Promise<ApiResponse<void>> => {
    return await apiService.post<ApiResponse<void>>(API.AUTH.RESET_PASSWORD, data);
  },

  /**
   * รีเฟรช token
   * @param data - refresh token
   * @returns token ใหม่
   */
  refreshToken: async (data: RefreshTokenRequest): Promise<ApiResponse<LoginResponse>> => {
    return await apiService.post<ApiResponse<LoginResponse>>(API.AUTH.REFRESH, data);
  },

  /**
   * ตรวจสอบ token
   * @returns ผลลัพธ์จากการตรวจสอบ token
   */
  validateToken: async (): Promise<TokenValidationResponse> => {
    return await apiService.post<TokenValidationResponse>(API.AUTH.VALIDATE);
  },

  /**
   * ออกจากระบบ
   * @returns ผลลัพธ์จากการออกจากระบบ
   */
  logout: async (): Promise<ApiResponse<void>> => {
    return await apiService.post<ApiResponse<void>>(API.AUTH.LOGOUT);
  },

  /**
   * เปลี่ยนรหัสผ่าน
   * @param data - ข้อมูลสำหรับการเปลี่ยนรหัสผ่าน
   * @returns ผลลัพธ์จากการเปลี่ยนรหัสผ่าน
   */
  updatePassword: async (data: UpdatePasswordRequest): Promise<ApiResponse<void>> => {
    return await apiService.post<ApiResponse<void>>(API.AUTH.UPDATE_PASSWORD, data);
  },

  /**
   * ยกเลิกการเข้าสู่ระบบทั้งหมด
   * @returns ผลลัพธ์จากการยกเลิกการเข้าสู่ระบบทั้งหมด
   */
  revokeAll: async (): Promise<ApiResponse<void>> => {
    return await apiService.post<ApiResponse<void>>(API.AUTH.REVOKE_ALL);
  }
};

export default authService;