import { create } from "zustand";
import { persist } from "zustand/middleware";
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from "@tanstack/react-query";
import userService from "@/services/admin/userService";
import type {
  User,
  UserStats,
  UserFilter,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateUserStatusRequest,
  UpdateUserRoleRequest,
  PaginatedResponse,
  //ApiResponse
} from "@/types";

const userKeys = {
  users: (params?: UserFilter) => ["user", "list", params] as const,
  userById: (id: string) => ["user", "byId", id] as const,
  userStats: () => ["user", "stats"] as const,
};

interface UserState {
  currentUser: User | null;
  selectedUser: User | null;
  usersList: User[];
  userStats: UserStats | null;
  setCurrentUser: (user: User | null) => void;
  setSelectedUser: (user: User | null) => void;
  setUsersList: (users: User[]) => void;
  setUserStats: (stats: UserStats | null) => void;
  clearUserData: () => void;
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      currentUser: null,
      selectedUser: null,
      usersList: [],
      userStats: null,
      setCurrentUser: (user) => set({ currentUser: user }),
      setSelectedUser: (user) => set({ selectedUser: user }),
      setUsersList: (users) => set({ usersList: users }),
      setUserStats: (stats) => set({ userStats: stats }),
      clearUserData: () => set({
        currentUser: null,
        selectedUser: null,
        usersList: [],
        userStats: null,
      }),
    }),
    {
      name: "user-storage",
      partialize: (state) => ({
        selectedUser: state.selectedUser,
      }),
    }
  )
);

// Admin+ Queries

/**
 * ดึงรายการผู้ใช้ทั้งหมด (สำหรับ admin+)
 */
export const useUsersListQuery = (
  params: UserFilter = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<PaginatedResponse<User>, Error> => {
  const { setUsersList } = useUserStore();

  return useQuery({
    queryKey: userKeys.users(params),
    queryFn: async () => {
      const response = await userService.getUsers(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch users list");
      }

      // สมมติว่า response.data มี structure แบบ { users: User[], pagination: PaginationMeta }
      const users = response.data.users || response.data;
      setUsersList(Array.isArray(users) ? users : []);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

/**
 * ดึงสถิติผู้ใช้ (สำหรับ admin+)
 */
export const useUserStatsQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<UserStats, Error> => {
  const { setUserStats } = useUserStore();

  return useQuery({
    queryKey: userKeys.userStats(),
    queryFn: async () => {
      const response = await userService.getUserStats();

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch user stats");
      }

      setUserStats(response.data);
      return response.data;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

/**
 * ดึงข้อมูลผู้ใช้ตาม ID (สำหรับ admin+)
 */
export const useUserByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<User, Error> => {
  const { setSelectedUser } = useUserStore();

  return useQuery({
    queryKey: userKeys.userById(id),
    queryFn: async () => {
      const response = await userService.getUserById(id);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch user data");
      }

      setSelectedUser(response.data);
      return response.data;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

// Admin+ Mutations

/**
 * สร้างผู้ใช้ใหม่ (สำหรับ admin+)
 */
export const useCreateUserMutation = (): UseMutationResult<
  User,
  Error,
  CreateUserRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateUserRequest) => {
      const response = await userService.createUser(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to create user");
      }

      return response.data;
    },
    onSuccess: () => {
      // Invalidate users list to refetch
      queryClient.invalidateQueries({ queryKey: userKeys.users() });
      queryClient.invalidateQueries({ queryKey: userKeys.userStats() });
    },
  });
};

/**
 * อัปเดตข้อมูลผู้ใช้ (สำหรับ admin+)
 */
export const useUpdateUserMutation = (): UseMutationResult<
  User,
  Error,
  { id: string; data: UpdateUserRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await userService.updateUser(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update user");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: userKeys.userById(id) });
      queryClient.invalidateQueries({ queryKey: userKeys.users() });
      queryClient.invalidateQueries({ queryKey: userKeys.userStats() });
    },
  });
};

/**
 * ลบผู้ใช้ (สำหรับ admin+)
 */
export const useDeleteUserMutation = (): UseMutationResult<
  null,
  Error,
  string,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await userService.deleteUser(id);

      if (!response.success) {
        throw new Error(response.message || "Failed to delete user");
      }

      return response.data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: userKeys.userById(id) });
      queryClient.invalidateQueries({ queryKey: userKeys.users() });
      queryClient.invalidateQueries({ queryKey: userKeys.userStats() });
    },
  });
};

/**
 * อัปเดตสถานะผู้ใช้ (สำหรับ admin+)
 */
export const useUpdateUserStatusMutation = (): UseMutationResult<
  User,
  Error,
  { id: string; data: UpdateUserStatusRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await userService.updateUserStatus(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update user status");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: userKeys.userById(id) });
      queryClient.invalidateQueries({ queryKey: userKeys.users() });
      queryClient.invalidateQueries({ queryKey: userKeys.userStats() });
    },
  });
};

/**
 * อัปเดตบทบาทผู้ใช้ (สำหรับ admin+)
 */
export const useUpdateUserRoleMutation = (): UseMutationResult<
  { id: string; new_role: string; changed_at: string },
  Error,
  { id: string; data: UpdateUserRoleRequest },
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await userService.updateUserRole(id, data);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Failed to update user role");
      }

      return response.data;
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: userKeys.userById(id) });
      queryClient.invalidateQueries({ queryKey: userKeys.users() });
      queryClient.invalidateQueries({ queryKey: userKeys.userStats() });
    },
  });
};
