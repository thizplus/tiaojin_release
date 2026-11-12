import { useCallback } from "react";
import * as React from "react";
import {
  useUserStore,
  useUsersListQuery,
  useUserByIdQuery,
  useUserStatsQuery,
  useCreateUserMutation,
  useUpdateUserMutation,
  useDeleteUserMutation,
  useUpdateUserStatusMutation,
  useUpdateUserRoleMutation
} from "@/stores/userStore";
import type {
  User,
  UserStats,
  UserFilter,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateUserStatusRequest,
  UpdateUserRoleRequest
} from "@/types";

interface UseUserReturn {
  // State
  currentUser: User | null;
  selectedUser: User | null;
  usersList: User[];
  userStats: UserStats | null;
  isLoading: boolean;

  // Pagination data
  pagination: {
    total: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
    hasMore: boolean;
  } | null;

  // Actions
  createUser: (data: CreateUserRequest) => Promise<User>;
  updateUser: (id: string, data: UpdateUserRequest) => Promise<User>;
  deleteUser: (id: string) => Promise<null>;
  updateUserStatus: (id: string, data: UpdateUserStatusRequest) => Promise<User>;
  updateUserRole: (id: string, data: UpdateUserRoleRequest) => Promise<{ id: string; new_role: string; changed_at: string }>;

  // State management
  setCurrentUser: (user: User | null) => void;
  setSelectedUser: (user: User | null) => void;
  setUsersList: (users: User[]) => void;
  setUserStats: (stats: UserStats | null) => void;
  clearUserData: () => void;

  // Queries
  getUsersList: (params?: UserFilter, options?: { enabled?: boolean }) => void;
  getUserById: (id: string, options?: { enabled?: boolean }) => void;
  getUserStats: (options?: { enabled?: boolean }) => void;
  refetchUsersList: () => void;
  refetchUserById: () => void;
  refetchUserStats: () => void;

  // Loading states
  isCreateLoading: boolean;
  isUpdateLoading: boolean;
  isDeleteLoading: boolean;
  isUpdateStatusLoading: boolean;
  isUpdateRoleLoading: boolean;
  isUsersListLoading: boolean;
  isUserByIdLoading: boolean;
  isUserStatsLoading: boolean;
}

export const useUser = (
  listParams?: UserFilter,
  userId?: string,
  options: { enableList?: boolean; enableById?: boolean; enableStats?: boolean } = {}
): UseUserReturn => {
  // Store state
  const {
    currentUser,
    selectedUser,
    usersList,
    userStats,
    setCurrentUser,
    setSelectedUser,
    setUsersList,
    setUserStats,
    clearUserData
  } = useUserStore();

  // Mutations
  const createUserMutation = useCreateUserMutation();
  const updateUserMutation = useUpdateUserMutation();
  const deleteUserMutation = useDeleteUserMutation();
  const updateUserStatusMutation = useUpdateUserStatusMutation();
  const updateUserRoleMutation = useUpdateUserRoleMutation();

  // Queries
  const {
    data: usersData,
    refetch: refetchUsersList,
    isLoading: isUsersListLoading
  } = useUsersListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    refetch: refetchUserById,
    isLoading: isUserByIdLoading
  } = useUserByIdQuery(userId || "", {
    enabled: options.enableById !== false && !!userId
  });

  const {
    refetch: refetchUserStats,
    isLoading: isUserStatsLoading
  } = useUserStatsQuery({
    enabled: options.enableStats === true
  });

  // Actions
  const createUser = useCallback(async (data: CreateUserRequest) => {
    const result = await createUserMutation.mutateAsync(data);
    return result;
  }, [createUserMutation]);

  const updateUser = useCallback(async (id: string, data: UpdateUserRequest) => {
    const result = await updateUserMutation.mutateAsync({ id, data });
    return result;
  }, [updateUserMutation]);

  const deleteUser = useCallback(async (id: string) => {
    const result = await deleteUserMutation.mutateAsync(id);
    return result;
  }, [deleteUserMutation]);

  const updateUserStatus = useCallback(async (id: string, data: UpdateUserStatusRequest) => {
    const result = await updateUserStatusMutation.mutateAsync({ id, data });
    return result;
  }, [updateUserStatusMutation]);

  const updateUserRole = useCallback(async (id: string, data: UpdateUserRoleRequest) => {
    const result = await updateUserRoleMutation.mutateAsync({ id, data });
    return result;
  }, [updateUserRoleMutation]);

  // Query management
  const getUsersList = useCallback(() => {
    refetchUsersList();
  }, [refetchUsersList]);

  const getUserById = useCallback(() => {
    refetchUserById();
  }, [refetchUserById]);

  const getUserStats = useCallback(() => {
    refetchUserStats();
  }, [refetchUserStats]);

  // Process users data and pagination
  const pagination = React.useMemo(() => {
    if (!usersData?.data?.pagination) return null;

    const paginationMeta = usersData.data.pagination;
    return {
      total: paginationMeta.total,
      totalPages: paginationMeta.total_pages,
      currentPage: paginationMeta.current_page,
      pageSize: paginationMeta.page_size,
      hasMore: paginationMeta.has_more,
    };
  }, [usersData]);

  // Calculate overall loading state
  const isLoading = createUserMutation.isPending ||
                   updateUserMutation.isPending ||
                   deleteUserMutation.isPending ||
                   updateUserStatusMutation.isPending ||
                   updateUserRoleMutation.isPending ||
                   isUsersListLoading ||
                   isUserByIdLoading ||
                   isUserStatsLoading;

  return {
    // State
    currentUser,
    selectedUser,
    usersList,
    userStats,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    createUser,
    updateUser,
    deleteUser,
    updateUserStatus,
    updateUserRole,

    // State management
    setCurrentUser,
    setSelectedUser,
    setUsersList,
    setUserStats,
    clearUserData,

    // Queries
    getUsersList,
    getUserById,
    getUserStats,
    refetchUsersList,
    refetchUserById,
    refetchUserStats,

    // Loading states
    isCreateLoading: createUserMutation.isPending,
    isUpdateLoading: updateUserMutation.isPending,
    isDeleteLoading: deleteUserMutation.isPending,
    isUpdateStatusLoading: updateUserStatusMutation.isPending,
    isUpdateRoleLoading: updateUserRoleMutation.isPending,
    isUsersListLoading,
    isUserByIdLoading,
    isUserStatsLoading,
  };
};

export default useUser;
