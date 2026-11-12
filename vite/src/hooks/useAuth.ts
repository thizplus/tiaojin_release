import { useCallback } from 'react';
import {
  useAuthStore,
  useLoginMutation,
  useRegisterMutation,
  useResetPasswordMutation,
  useRefreshTokenMutation,
  useValidateTokenQuery,
  useLogoutMutation,
  useUpdatePasswordMutation,
  useRevokeAllMutation,
  useUserProfileQuery
} from '@/stores/authStore';
import type {
  LoginRequest,
  RegisterRequest,
  ResetPasswordRequest,
  RefreshTokenRequest,
  UpdatePasswordRequest,
  UserProfile
} from '@/types';

interface UseAuthReturn {
  // State
  user: UserProfile | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  login: (data: LoginRequest) => Promise<{ token: string; user?: UserProfile }>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  resetPassword: (data: ResetPasswordRequest) => Promise<void>;
  updatePassword: (data: UpdatePasswordRequest) => Promise<void>;
  refreshToken: (data: RefreshTokenRequest) => Promise<void>;
  revokeAllSessions: () => Promise<void>;

  // Validation
  validateToken: () => void;

  // Profile
  refetchProfile: () => void;

  // Loading states
  isLoginLoading: boolean;
  isRegisterLoading: boolean;
  isLogoutLoading: boolean;
  isResetPasswordLoading: boolean;
  isUpdatePasswordLoading: boolean;
  isRefreshTokenLoading: boolean;
  isRevokeAllLoading: boolean;
  isValidatingToken: boolean;
  isProfileLoading: boolean;
}

export const useAuth = (): UseAuthReturn => {
  // Store state
  const { user, token, isAuthenticated } = useAuthStore();

  // Mutations
  const loginMutation = useLoginMutation();
  const registerMutation = useRegisterMutation();
  const logoutMutation = useLogoutMutation();
  const resetPasswordMutation = useResetPasswordMutation();
  const updatePasswordMutation = useUpdatePasswordMutation();
  const refreshTokenMutation = useRefreshTokenMutation();
  const revokeAllMutation = useRevokeAllMutation();

  // Queries
  const {
    refetch: refetchValidation,
    isLoading: isValidatingToken,
    isError: isValidationError
  } = useValidateTokenQuery({ enabled: !!token && isAuthenticated });

  const {
    refetch: refetchProfile,
    isLoading: isProfileLoading,
    isError: isProfileError
  } = useUserProfileQuery({ enabled: isAuthenticated && !isValidationError });

  // Actions
  const login = useCallback(async (data: LoginRequest) => {
    const result = await loginMutation.mutateAsync(data);
    return result;
  }, [loginMutation]);

  const register = useCallback(async (data: RegisterRequest) => {
    await registerMutation.mutateAsync(data);
  }, [registerMutation]);

  const logout = useCallback(async () => {
    await logoutMutation.mutateAsync();
  }, [logoutMutation]);

  const resetPassword = useCallback(async (data: ResetPasswordRequest) => {
    await resetPasswordMutation.mutateAsync(data);
  }, [resetPasswordMutation]);

  const updatePassword = useCallback(async (data: UpdatePasswordRequest) => {
    await updatePasswordMutation.mutateAsync(data);
  }, [updatePasswordMutation]);

  const refreshToken = useCallback(async (data: RefreshTokenRequest) => {
    await refreshTokenMutation.mutateAsync(data);
  }, [refreshTokenMutation]);

  const revokeAllSessions = useCallback(async () => {
    await revokeAllMutation.mutateAsync();
  }, [revokeAllMutation]);

  const validateToken = useCallback(() => {
    refetchValidation();
  }, [refetchValidation]);

  const refetchUserProfile = useCallback(() => {
    refetchProfile();
  }, [refetchProfile]);

  // Calculate overall loading state
  // ไม่รวม isValidatingToken เมื่อ validation error เพื่อป้องกัน infinite loop
  const isLoading = loginMutation.isPending ||
                   registerMutation.isPending ||
                   logoutMutation.isPending ||
                   resetPasswordMutation.isPending ||
                   updatePasswordMutation.isPending ||
                   refreshTokenMutation.isPending ||
                   revokeAllMutation.isPending ||
                   (isValidatingToken && !isValidationError) ||
                   (isProfileLoading && !isProfileError);

  return {
    // State
    user,
    token,
    isAuthenticated,
    isLoading,

    // Actions
    login,
    register,
    logout,
    resetPassword,
    updatePassword,
    refreshToken,
    revokeAllSessions,

    // Validation
    validateToken,

    // Profile
    refetchProfile: refetchUserProfile,

    // Loading states
    isLoginLoading: loginMutation.isPending,
    isRegisterLoading: registerMutation.isPending,
    isLogoutLoading: logoutMutation.isPending,
    isResetPasswordLoading: resetPasswordMutation.isPending,
    isUpdatePasswordLoading: updatePasswordMutation.isPending,
    isRefreshTokenLoading: refreshTokenMutation.isPending,
    isRevokeAllLoading: revokeAllMutation.isPending,
    isValidatingToken,
    isProfileLoading,
  };
};

export default useAuth;