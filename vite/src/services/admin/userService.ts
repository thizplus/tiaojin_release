//src/services/admin/userService.ts

import { API } from "@/constants/api";
import apiService from "@/services/apiService";
import {
  type User,
  type UserStats,
  type UserFilter,
  type CreateUserRequest,
  type UpdateUserRequest,
  type UpdateUserStatusRequest,
  type UpdateUserRoleRequest,
  type PaginatedResponse,
  type ApiResponse
} from "@/types";

/**
 * Service สำหรับจัดการข้อมูลผู้ใช้ (Admin)
 */
const userService = {
  /**
   * ดึงรายการผู้ใช้ทั้งหมด (สำหรับ admin+)
   * GET /api/v1/admin/users
   * @param params - ตัวกรองและการแบ่งหน้า
   * @returns รายการผู้ใช้
   */
  getUsers: async (params: UserFilter = {}): Promise<PaginatedResponse<User>> => {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append("search", params.search);
    if (params.role) queryParams.append("role", params.role);
    if (params.status) queryParams.append("status", params.status);
    if (params.page !== undefined) queryParams.append("page", params.page.toString());
    if (params.limit !== undefined) queryParams.append("limit", params.limit.toString());
    if (params.sort) queryParams.append("sort", params.sort);
    if (params.order) queryParams.append("order", params.order);

    const url = queryParams.toString() ? `${API.USER.LIST}?${queryParams.toString()}` : API.USER.LIST;
    return await apiService.get<PaginatedResponse<User>>(url);
  },

  /**
   * ดึงสถิติผู้ใช้ (สำหรับ admin+)
   * GET /api/v1/admin/users/stats
   * @returns สถิติผู้ใช้
   */
  getUserStats: async (): Promise<ApiResponse<UserStats>> => {
    return await apiService.get<ApiResponse<UserStats>>(API.USER.STATS);
  },

  /**
   * ดึงข้อมูลผู้ใช้ตาม ID (สำหรับ admin+)
   * GET /api/v1/admin/users/:id
   * @param id - ID ของผู้ใช้
   * @returns ข้อมูลผู้ใช้
   */
  getUserById: async (id: string): Promise<ApiResponse<User>> => {
    return await apiService.get<ApiResponse<User>>(API.USER.GET_BY_ID(id));
  },

  /**
   * สร้างผู้ใช้ใหม่ (สำหรับ admin+) - อาจใช้ auth register หรือสร้างแยก
   * POST /api/v1/admin/users (ถ้ามี endpoint นี้)
   * @param data - ข้อมูลผู้ใช้ที่ต้องการสร้าง
   * @returns ข้อมูลผู้ใช้ที่สร้างแล้ว
   */
  createUser: async (data: CreateUserRequest): Promise<ApiResponse<User>> => {
    // อาจต้องใช้ AUTH.REGISTER หรือสร้าง endpoint ใหม่
    return await apiService.post<ApiResponse<User>>("/admin/users", data);
  },

  /**
   * อัปเดตข้อมูลผู้ใช้ (สำหรับ admin+)
   * PUT /api/v1/admin/users/:id
   * @param id - ID ของผู้ใช้
   * @param data - ข้อมูลที่ต้องการอัปเดต
   * @returns ข้อมูลผู้ใช้ที่อัปเดตแล้ว
   */
  updateUser: async (id: string, data: UpdateUserRequest): Promise<ApiResponse<User>> => {
    return await apiService.put<ApiResponse<User>>(API.USER.UPDATE_BY_ID(id), data);
  },

  /**
   * ลบผู้ใช้ (สำหรับ admin+)
   * DELETE /api/v1/admin/users/:id
   * @param id - ID ของผู้ใช้
   * @returns ผลลัพธ์การลบ
   */
  deleteUser: async (id: string): Promise<ApiResponse<null>> => {
    return await apiService.delete<ApiResponse<null>>(API.USER.DELETE(id));
  },

  /**
   * อัปเดตสถานะผู้ใช้ (สำหรับ admin+)
   * PATCH /api/v1/admin/users/:id/status
   * @param id - ID ของผู้ใช้
   * @param data - สถานะใหม่
   * @returns ข้อมูลผู้ใช้ที่อัปเดตแล้ว
   */
  updateUserStatus: async (id: string, data: UpdateUserStatusRequest): Promise<ApiResponse<User>> => {
    return await apiService.patch<ApiResponse<User>>(API.USER.UPDATE_STATUS(id), data);
  },

  /**
   * อัปเดตบทบาทผู้ใช้ (สำหรับ admin+)
   * PATCH /api/v1/admin/users/:id/role
   * @param id - ID ของผู้ใช้
   * @param data - บทบาทใหม่
   * @returns ข้อมูลการเปลี่ยนแปลง
   */
  updateUserRole: async (id: string, data: UpdateUserRoleRequest): Promise<ApiResponse<{ id: string; new_role: string; changed_at: string }>> => {
    return await apiService.patch<ApiResponse<{ id: string; new_role: string; changed_at: string }>>(API.USER.UPDATE_ROLE(id), data);
  },
};

export default userService;
