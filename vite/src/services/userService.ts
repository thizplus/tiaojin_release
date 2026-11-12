import { USER_API as API } from '@/constants/api';
import apiService from './apiService';
import type {
  User,
  UserStats,
  UserFilter,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateUserStatusRequest,
  UpdateUserRoleRequest,
  PaginatedResponse,
  ApiResponse
} from '@/types';

const userService = {
  async getProfile(): Promise<ApiResponse<User>> {
    return await apiService.get<ApiResponse<User>>(API.PROFILE);
  },

  async getUsers(params: UserFilter = {}): Promise<ApiResponse<PaginatedResponse<User>>> {
    const queryParams = new URLSearchParams();

    if (params.page) queryParams.append('page', params.page.toString());
    if (params.limit) queryParams.append('limit', params.limit.toString());
    if (params.search) queryParams.append('search', params.search);
    if (params.role) queryParams.append('role', params.role);
    if (params.status) queryParams.append('status', params.status);
    if (params.sort) queryParams.append('sort', params.sort);
    if (params.order) queryParams.append('order', params.order);

    const url = queryParams.toString() ? `${API.LIST}?${queryParams.toString()}` : API.LIST;
    return await apiService.get<ApiResponse<PaginatedResponse<User>>>(url);
  },

  async getUserStats(): Promise<ApiResponse<UserStats>> {
    return await apiService.get<ApiResponse<UserStats>>(API.STATS);
  },

  async getUserById(id: string): Promise<ApiResponse<User>> {
    return await apiService.get<ApiResponse<User>>(API.GET_BY_ID(id));
  },

  async createUser(data: CreateUserRequest): Promise<ApiResponse<User>> {
    return await apiService.post<ApiResponse<User>>(API.CREATE, data);
  },

  async updateUser(id: string, data: UpdateUserRequest): Promise<ApiResponse<User>> {
    return await apiService.put<ApiResponse<User>>(API.UPDATE_BY_ID(id), data);
  },

  async deleteUser(id: string): Promise<ApiResponse<null>> {
    return await apiService.delete<ApiResponse<null>>(API.DELETE(id));
  },

  async updateUserStatus(id: string, data: UpdateUserStatusRequest): Promise<ApiResponse<User>> {
    return await apiService.put<ApiResponse<User>>(API.UPDATE_STATUS(id), data);
  },

  async updateUserRole(id: string, data: UpdateUserRoleRequest): Promise<ApiResponse<{ id: string; new_role: string; changed_at: string }>> {
    return await apiService.put<ApiResponse<{ id: string; new_role: string; changed_at: string }>>(API.UPDATE_ROLE(id), data);
  }
};

export default userService;