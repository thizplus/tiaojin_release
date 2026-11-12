import { create } from "zustand";
import { persist } from "zustand/middleware";
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from "@tanstack/react-query";
import agentProfileService from "@/services/agent/agentProfileService";
import { AGENT_PROFILE_MESSAGES } from "@/constants/agent";
import type {
  AgentProfile,
  CreateAgentProfileRequest,
  UpdateAgentProfileRequest,
  ProfileCheckResponse,
  //ApiResponse
} from "@/types";

// Query Keys
const agentProfileKeys = {
  all: ['agentProfile'] as const,
  profile: ['agentProfile', 'profile'] as const,
  check: ['agentProfile', 'check'] as const,
  uploads: ['agentProfile', 'uploads'] as const,
} as const;

// Store State Interface
interface AgentProfileState {
  profile: AgentProfile | null;
  profileStatus: ProfileCheckResponse | null;
  isInitialized: boolean;
  lastUpdated: string | null;

  // Actions
  setProfile: (profile: AgentProfile | null) => void;
  setProfileStatus: (status: ProfileCheckResponse | null) => void;
  setInitialized: (initialized: boolean) => void;
  updateLastUpdated: () => void;
  clearProfileData: () => void;
}

// Zustand Store
export const useAgentProfileStore = create<AgentProfileState>()(
  persist(
    (set) => ({
      profile: null,
      profileStatus: null,
      isInitialized: false,
      lastUpdated: null,

      setProfile: (profile) => set({
        profile,
        lastUpdated: new Date().toISOString()
      }),

      setProfileStatus: (status) => set({
        profileStatus: status,
        lastUpdated: new Date().toISOString()
      }),

      setInitialized: (initialized) => set({ isInitialized: initialized }),

      updateLastUpdated: () => set({
        lastUpdated: new Date().toISOString()
      }),

      clearProfileData: () => set({
        profile: null,
        profileStatus: null,
        isInitialized: false,
        lastUpdated: null,
      }),
    }),
    {
      name: "agent-profile-storage",
      partialize: (state) => ({
        profile: state.profile,
        profileStatus: state.profileStatus,
        lastUpdated: state.lastUpdated,
      }),
    }
  )
);

// Query Hooks
export const useProfileCheckQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<ProfileCheckResponse, Error> => {
  const { setProfileStatus, setInitialized } = useAgentProfileStore();

  return useQuery({
    queryKey: agentProfileKeys.check,
    queryFn: async () => {
      try {
        const response = await agentProfileService.checkProfile();

        if (!response.success || !response.data) {
          throw new Error(response.message || "Unable to check profile status");
        }

        setProfileStatus(response.data);
        setInitialized(true);
        return response.data;
      } catch (error) {
        setInitialized(true);
        throw error;
      }
    },
    enabled: options.enabled !== false,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    refetchOnWindowFocus: false,
    retry: 2,
  });
};

export const useAgentProfileQuery = (
  options: { enabled?: boolean } = {}
): UseQueryResult<AgentProfile, Error> => {
  const { setProfile } = useAgentProfileStore();

  return useQuery({
    queryKey: agentProfileKeys.profile,
    queryFn: async () => {
      const response = await agentProfileService.getProfile();

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to fetch profile data");
      }

      setProfile(response.data);
      return response.data;
    },
    enabled: options.enabled !== false,
    staleTime: 2 * 60 * 1000, // 2 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    refetchOnWindowFocus: false,
    retry: 2,
  });
};

// Mutation Hooks
export const useCreateProfileMutation = (
  options: {
    onSuccess?: (profile: AgentProfile) => void;
    onError?: (error: Error) => void;
  } = {}
): UseMutationResult<AgentProfile, Error, CreateAgentProfileRequest, unknown> => {
  const queryClient = useQueryClient();
  const { setProfile, setProfileStatus } = useAgentProfileStore();

  return useMutation({
    mutationFn: async (data: CreateAgentProfileRequest) => {
      // Validate data first
      const validation = agentProfileService.validateProfileData(data);
      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors).join(', ');
        throw new Error(errorMessage);
      }

      const response = await agentProfileService.createProfile(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || AGENT_PROFILE_MESSAGES.CREATE_FAILED);
      }

      return response.data;
    },
    onSuccess: (profile) => {
      setProfile(profile);
      setProfileStatus({
        has_profile: true,
        profile_id: profile.id,
        company_name: profile.name
      });

      // Invalidate and refetch queries
      queryClient.invalidateQueries({ queryKey: agentProfileKeys.profile });
      queryClient.invalidateQueries({ queryKey: agentProfileKeys.check });

      options.onSuccess?.(profile);
    },
    onError: (error) => {
      options.onError?.(error);
    },
  });
};

export const useUpdateProfileMutation = (
  options: {
    onSuccess?: (profile: AgentProfile) => void;
    onError?: (error: Error) => void;
  } = {}
): UseMutationResult<AgentProfile, Error, UpdateAgentProfileRequest, unknown> => {
  const queryClient = useQueryClient();
  const { setProfile } = useAgentProfileStore();

  return useMutation({
    mutationFn: async (data: UpdateAgentProfileRequest) => {
      // Validate data first
      const validation = agentProfileService.validateProfileData(data);
      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors).join(', ');
        throw new Error(errorMessage);
      }

      const response = await agentProfileService.updateProfile(data);

      if (!response.success || !response.data) {
        throw new Error(response.message || AGENT_PROFILE_MESSAGES.UPDATE_FAILED);
      }

      return response.data;
    },
    onSuccess: (profile) => {
      setProfile(profile);

      // Invalidate and refetch queries
      queryClient.invalidateQueries({ queryKey: agentProfileKeys.profile });
      queryClient.invalidateQueries({ queryKey: agentProfileKeys.check });

      options.onSuccess?.(profile);
    },
    onError: (error) => {
      options.onError?.(error);
    },
  });
};

// File Upload Mutations
export const useUploadFileMutation = (
  fileType: 'logo' | 'atta' | 'dbd' | 'pp20',
  options: {
    onSuccess?: (result: { url: string; fileName: string }) => void;
    onError?: (error: Error) => void;
  } = {}
): UseMutationResult<{ url: string; fileName: string }, Error, File, unknown> => {
  const queryClient = useQueryClient();

  const categoryMap = {
    logo: 'agent_profiles' as const,
    atta: 'agent_licenses' as const,
    dbd: 'agent_licenses' as const,
    pp20: 'agent_licenses' as const,
  };

  return useMutation({
    mutationFn: async (file: File) => {
      // Validate file first
      const validation = agentProfileService.validateFile(file, fileType);
      if (!validation.isValid) {
        throw new Error(validation.error);
      }

      const category = categoryMap[fileType];
      const response = await agentProfileService.uploadFile(file, category, fileType);

      if (!response.success || !response.data) {
        throw new Error(response.message || AGENT_PROFILE_MESSAGES.FILE_UPLOAD_FAILED);
      }

      return response.data;
    },
    onSuccess: (result) => {
      // Invalidate profile query to refetch with new file URL
      queryClient.invalidateQueries({ queryKey: agentProfileKeys.profile });

      options.onSuccess?.(result);
    },
    onError: (error) => {
      options.onError?.(error);
    },
  });
};

// Legacy upload hooks for backward compatibility
export const useUploadLogoMutation = (
  options: {
    onSuccess?: (result: { url: string; fileName: string }) => void;
    onError?: (error: Error) => void;
  } = {}
) => useUploadFileMutation('logo', options);

export const useUploadAttaMutation = (
  options: {
    onSuccess?: (result: { url: string; fileName: string }) => void;
    onError?: (error: Error) => void;
  } = {}
) => useUploadFileMutation('atta', options);

export const useUploadDbdMutation = (
  options: {
    onSuccess?: (result: { url: string; fileName: string }) => void;
    onError?: (error: Error) => void;
  } = {}
) => useUploadFileMutation('dbd', options);

export const useUploadPp20Mutation = (
  options: {
    onSuccess?: (result: { url: string; fileName: string }) => void;
    onError?: (error: Error) => void;
  } = {}
) => useUploadFileMutation('pp20', options);

// Utility Hooks
export const useProfileState = () => {
  const store = useAgentProfileStore();
  return {
    profile: store.profile,
    profileStatus: store.profileStatus,
    isInitialized: store.isInitialized,
    lastUpdated: store.lastUpdated,
    hasProfile: store.profileStatus?.has_profile ?? false,
    companyName: store.profile?.name || store.profileStatus?.company_name,
  };
};

// Clear all profile data (useful for logout)
export const useClearAllProfileData = () => {
  const queryClient = useQueryClient();
  const { clearProfileData } = useAgentProfileStore();

  return () => {
    clearProfileData();
    queryClient.removeQueries({ queryKey: agentProfileKeys.all });
  };
};

export default agentProfileKeys;