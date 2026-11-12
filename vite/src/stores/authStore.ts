import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, QueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import authService from '@/services/authService';
import userService from '@/services/userService';
import type {
  UserProfile,
  LoginRequest,
  RegisterRequest,
  ResetPasswordRequest,
  RefreshTokenRequest,
  UpdatePasswordRequest,
  //ApiResponse,
  LoginResponse,
  User,
  TokenValidationResponse
} from '@/types';

const authKeys = {
  profile: () => ['auth', 'profile'] as const,
  validation: () => ['auth', 'validation'] as const,
};

// Global reference สำหรับ cancel queries เมื่อ logout
let globalQueryClient: QueryClient | null = null;

export const setGlobalQueryClient = (client: QueryClient) => {
  globalQueryClient = client;
};

// Force logout พร้อม cancel queries
export const forceLogout = () => {
  // Cancel all queries ก่อน
  if (globalQueryClient) {
    globalQueryClient.cancelQueries();
    globalQueryClient.clear();
  }
  // แล้วค่อย logout
  useAuthStore.getState().logout();
};

interface AuthState {
  user: UserProfile | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoggingOut: boolean; // Flag ป้องกัน infinite loop
  setAuth: (token: string, isAuthenticated: boolean) => void;
  setUser: (user: UserProfile | null) => void;
  logout: () => void;
  clearAuth: () => void;
  setLoggingOut: (value: boolean) => void;
}

// Storage version สำหรับ migration
const CURRENT_STORAGE_VERSION = 1;

// ตรวจสอบว่า user object ถูก format หรือไม่ (ระบบใหม่ใช้ UserProfile)
const isValidUserFormat = (user: unknown): user is UserProfile => {
  if (!user || typeof user !== 'object') return false;
  const u = user as Record<string, unknown>;
  // ระบบใหม่ใช้ created_at, updated_at (snake_case)
  // ระบบเก่าใช้ createdAt, updatedAt (camelCase)
  return 'created_at' in u || !('createdAt' in u);
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoggingOut: false,
      setAuth: (token, isAuthenticated) => set({ token, isAuthenticated, isLoggingOut: false }),
      setUser: (user) => set({ user }),
      logout: () => set({ user: null, token: null, isAuthenticated: false, isLoggingOut: true }),
      clearAuth: () => set({ user: null, token: null, isAuthenticated: false, isLoggingOut: true }),
      setLoggingOut: (value) => set({ isLoggingOut: value }),
    }),
    {
      name: 'auth-storage',
      version: CURRENT_STORAGE_VERSION,
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        isAuthenticated: state.isAuthenticated
        // ไม่ persist isLoggingOut
      }),
      // Migration: เมื่อ version ไม่ตรงหรือ format เก่า ให้ clear ทิ้ง
      migrate: (persistedState, version) => {
        // ถ้า version เก่ากว่า หรือ format ไม่ถูกต้อง ให้ clear
        if (version < CURRENT_STORAGE_VERSION) {
          console.log('[Auth] Migrating from old storage version, clearing...');
          return {
            user: null,
            token: null,
            isAuthenticated: false,
            isLoggingOut: false,
          };
        }

        const state = persistedState as Partial<AuthState>;

        // ตรวจสอบ user format
        if (state.user && !isValidUserFormat(state.user)) {
          console.log('[Auth] Invalid user format detected, clearing...');
          return {
            user: null,
            token: null,
            isAuthenticated: false,
            isLoggingOut: false,
          };
        }

        return state as AuthState;
      },
    }
  )
);

// Clear storage ถ้า format เก่า (เรียกตอน app start)
export const clearOldStorage = () => {
  try {
    const stored = localStorage.getItem('auth-storage');
    if (stored) {
      const parsed = JSON.parse(stored);
      // ถ้าไม่มี version หรือ version เก่า ให้ clear
      if (!parsed.version || parsed.version < CURRENT_STORAGE_VERSION) {
        console.log('[Auth] Old storage detected, clearing...');
        localStorage.removeItem('auth-storage');
        return true;
      }
      // ถ้า user format ไม่ถูกต้อง ให้ clear
      if (parsed.state?.user && !isValidUserFormat(parsed.state.user)) {
        console.log('[Auth] Invalid user format, clearing...');
        localStorage.removeItem('auth-storage');
        return true;
      }
    }
  } catch (e) {
    console.error('[Auth] Error checking storage:', e);
    localStorage.removeItem('auth-storage');
    return true;
  }
  return false;
};

export const useUserProfileQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<UserProfile, Error> => {
  const { token, isAuthenticated, isLoggingOut, setUser } = useAuthStore();

  return useQuery({
    queryKey: authKeys.profile(),
    queryFn: async () => {
      // ถ้ากำลัง logout อยู่ ไม่ต้อง fetch
      if (useAuthStore.getState().isLoggingOut) {
        throw new Error('Logging out');
      }
      if (!token) {
        throw new Error('Token not found for authentication');
      }
      const response = await userService.getProfile();

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch user data');
      }

      // Transform User to UserProfile
      const userProfile: UserProfile = {
        ...response.data,
        created_at: response.data.created_at || new Date().toISOString(),
        updated_at: response.data.updated_at || new Date().toISOString(),
        company_profile: response.data.company_profile || undefined
      };

      setUser(userProfile);
      return userProfile;
    },
    // เพิ่มเงื่อนไข: ต้อง authenticated และ ไม่กำลัง logout
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!token && isAuthenticated && !isLoggingOut,
    retry: (failureCount, error) => {
      // 401 = unauthorized, 404 = user deleted
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const errorStatus = (error as any)?.status;
      const isAuthError = errorStatus === 401 || errorStatus === 404;

      if (isAuthError) {
        forceLogout(); // ใช้ forceLogout เพื่อ cancel queries ด้วย
        return false;
      }
      return failureCount < 2;
    },
    refetchOnWindowFocus: false,
    refetchOnMount: false,
    refetchOnReconnect: false,
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
  });
};

export const useLoginMutation = (): UseMutationResult<
  { token: string; user?: UserProfile },
  Error,
  LoginRequest,
  unknown
> => {
  const queryClient = useQueryClient();
  const { setAuth, setUser } = useAuthStore();

  return useMutation({
    mutationFn: async (data: LoginRequest) => {
      const response = await authService.login(data);
      if (!response.success || !response.data?.access_token) {
        throw new Error(response.message || 'Login failed');
      }
      setAuth(response.data.access_token, true);

      const profileResponse = await userService.getProfile();
      if (profileResponse.success && profileResponse.data) {
        // Transform User to UserProfile
        const userProfile: UserProfile = {
          ...profileResponse.data,
          created_at: profileResponse.data.created_at || new Date().toISOString(),
          updated_at: profileResponse.data.updated_at || new Date().toISOString(),
          company_profile: profileResponse.data.company_profile || undefined
        };

        setUser(userProfile);
        queryClient.setQueryData(authKeys.profile(), userProfile);
        return { token: response.data.access_token, user: userProfile };
      }

      return { token: response.data.access_token };
    }
  });
};

export const useRegisterMutation = (): UseMutationResult<
  User,
  Error,
  RegisterRequest,
  unknown
> => {
  return useMutation({
    mutationFn: async (data: RegisterRequest) => {
      const response = await authService.register(data);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Registration failed');
      }
      return response.data;
    }
  });
};

export const useResetPasswordMutation = (): UseMutationResult<
  void,
  Error,
  ResetPasswordRequest,
  unknown
> => {
  return useMutation({
    mutationFn: async (data: ResetPasswordRequest) => {
      const response = await authService.resetPassword(data);
      if (!response.success) {
        throw new Error(response.message || 'Password reset failed');
      }
    }
  });
};

export const useRefreshTokenMutation = (): UseMutationResult<
  LoginResponse,
  Error,
  RefreshTokenRequest,
  unknown
> => {
  const { setAuth } = useAuthStore();

  return useMutation({
    mutationFn: async (data: RefreshTokenRequest) => {
      const response = await authService.refreshToken(data);
      if (!response.success || !response.data) {
        throw new Error(response.message || 'Token refresh failed');
      }

      setAuth(response.data.access_token, true);
      return response.data;
    }
  });
};

export const useValidateTokenQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<TokenValidationResponse, Error> => {
  const { token, isAuthenticated, isLoggingOut } = useAuthStore();

  return useQuery({
    queryKey: authKeys.validation(),
    queryFn: async () => {
      // ถ้ากำลัง logout อยู่ ไม่ต้อง validate
      if (useAuthStore.getState().isLoggingOut) {
        throw new Error('Logging out');
      }

      const response = await authService.validateToken();
      if (!response.status) {
        forceLogout(); // ใช้ forceLogout เพื่อ cancel queries ด้วย
        throw new Error('Token validation failed');
      }
      return response;
    },
    // เพิ่มเงื่อนไข: ต้อง authenticated และ ไม่กำลัง logout
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!token && isAuthenticated && !isLoggingOut,
    retry: false,
    refetchOnWindowFocus: false,
    refetchOnMount: false,
    refetchOnReconnect: false,
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
  });
};

export const useLogoutMutation = (): UseMutationResult<
  void,
  Error,
  void,
  unknown
> => {
  const queryClient = useQueryClient();
  const { clearAuth } = useAuthStore();

  return useMutation({
    mutationFn: async () => {
      try {
        await authService.logout();
      } catch (error) {
        // Continue with local logout even if server logout fails
        console.warn('Server logout failed:', error);
      }
    },
    onSettled: () => {
      clearAuth();
      queryClient.clear();
    }
  });
};

export const useUpdatePasswordMutation = (): UseMutationResult<
  void,
  Error,
  UpdatePasswordRequest,
  unknown
> => {
  return useMutation({
    mutationFn: async (data: UpdatePasswordRequest) => {
      const response = await authService.updatePassword(data);
      if (!response.success) {
        throw new Error(response.message || 'Password update failed');
      }
    }
  });
};

export const useRevokeAllMutation = (): UseMutationResult<
  void,
  Error,
  void,
  unknown
> => {
  const queryClient = useQueryClient();
  const { clearAuth } = useAuthStore();

  return useMutation({
    mutationFn: async () => {
      const response = await authService.revokeAll();
      if (!response.success) {
        throw new Error(response.message || 'Revoke all sessions failed');
      }
    },
    onSuccess: () => {
      clearAuth();
      queryClient.clear();
    }
  });
};