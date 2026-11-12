import { useCallback, useMemo } from "react";
import {
  useAgentProfileStore,
  useProfileCheckQuery,
  useAgentProfileQuery,
  useCreateProfileMutation,
  useUpdateProfileMutation,
  useUploadFileMutation,
  //useProfileState
} from "@/stores/agent/agentProfileStore";
import { useToast } from "@/hooks/use-toast";
import { AGENT_PROFILE_MESSAGES } from "@/constants/agent";
import type {
  AgentProfile,
  CreateAgentProfileRequest,
  UpdateAgentProfileRequest,
  ProfileCheckResponse
} from "@/types";

interface UseAgentProfileOptions {
  enableCheck?: boolean;
  enableProfile?: boolean;
  autoFetchProfile?: boolean;
  showToast?: boolean;
}

interface UseAgentProfileReturn {
  // State
  profile: AgentProfile | null;
  profileStatus: ProfileCheckResponse | null;
  isInitialized: boolean;
  lastUpdated: string | null;

  // Loading states
  isCheckingProfile: boolean;
  isLoadingProfile: boolean;
  isCreatingProfile: boolean;
  isUpdatingProfile: boolean;
  isUploadingFile: boolean;

  // Actions
  createProfile: (data: CreateAgentProfileRequest) => Promise<AgentProfile>;
  updateProfile: (data: UpdateAgentProfileRequest) => Promise<AgentProfile>;
  uploadFile: (file: File, fileType: 'logo' | 'atta' | 'dbd' | 'pp20') => Promise<{ url: string; fileName: string }>;

  // Legacy upload methods
  uploadLogo: (file: File) => Promise<{ url: string; fileName: string }>;
  uploadAtta: (file: File) => Promise<{ url: string; fileName: string }>;
  uploadDbd: (file: File) => Promise<{ url: string; fileName: string }>;
  uploadPp20: (file: File) => Promise<{ url: string; fileName: string }>;

  // State management
  setProfile: (profile: AgentProfile | null) => void;
  setProfileStatus: (status: ProfileCheckResponse | null) => void;
  clearProfileData: () => void;

  // Queries
  refetchProfileCheck: () => void;
  refetchProfile: () => void;

  // Computed
  hasProfile: boolean;
  companyName: string | undefined;
  isLoading: boolean;
  canCreateProfile: boolean;
  canUpdateProfile: boolean;

  // Error states
  checkError: Error | null;
  profileError: Error | null;
}

export const useAgentProfile = (
  options: UseAgentProfileOptions = {}
): UseAgentProfileReturn => {
  const {
    enableCheck = true,
    enableProfile = false,
    autoFetchProfile = true,
    showToast = true
  } = options;

  const { toast } = useToast();

  // Store state
  const {
    profile,
    profileStatus,
    isInitialized,
    lastUpdated,
    setProfile,
    setProfileStatus,
    clearProfileData
  } = useAgentProfileStore();

  // Profile state helper
  //const profileState = useProfileState();

  // Queries
  const {
    data: checkData,
    refetch: refetchProfileCheck,
    isLoading: isCheckingProfile,
    error: checkError
  } = useProfileCheckQuery({
    enabled: enableCheck
  });

  const {
    data: profileData,
    refetch: refetchProfile,
    isLoading: isLoadingProfile,
    error: profileError
  } = useAgentProfileQuery({
    enabled: enableProfile || (autoFetchProfile && (checkData?.has_profile || profileStatus?.has_profile))
  });

  // Mutations with toast notifications
  const createProfileMutation = useCreateProfileMutation({
    onSuccess: () => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.PROFILE_CREATED,
        });
      }
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
    }
  });

  const updateProfileMutation = useUpdateProfileMutation({
    onSuccess: () => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.PROFILE_UPDATED,
        });
      }
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
    }
  });

  // File upload mutations
  const logoUploadMutation = useUploadFileMutation('logo', {
    onSuccess: () => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.FILE_UPLOADED,
        });
      }
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
    }
  });

  const attaUploadMutation = useUploadFileMutation('atta', {
    onSuccess: () => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.FILE_UPLOADED,
        });
      }
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
    }
  });

  const dbdUploadMutation = useUploadFileMutation('dbd', {
    onSuccess: () => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.FILE_UPLOADED,
        });
      }
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
    }
  });

  const pp20UploadMutation = useUploadFileMutation('pp20', {
    onSuccess: () => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.FILE_UPLOADED,
        });
      }
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
    }
  });

  // Actions
  const createProfile = useCallback(async (data: CreateAgentProfileRequest): Promise<AgentProfile> => {
    const result = await createProfileMutation.mutateAsync(data);
    return result;
  }, [createProfileMutation]);

  const updateProfile = useCallback(async (data: UpdateAgentProfileRequest): Promise<AgentProfile> => {
    const result = await updateProfileMutation.mutateAsync(data);
    return result;
  }, [updateProfileMutation]);

  const uploadFile = useCallback(async (
    file: File,
    fileType: 'logo' | 'atta' | 'dbd' | 'pp20'
  ): Promise<{ url: string; fileName: string }> => {
    const mutationMap = {
      logo: logoUploadMutation,
      atta: attaUploadMutation,
      dbd: dbdUploadMutation,
      pp20: pp20UploadMutation
    };

    const mutation = mutationMap[fileType];
    const result = await mutation.mutateAsync(file);
    return result;
  }, [logoUploadMutation, attaUploadMutation, dbdUploadMutation, pp20UploadMutation]);

  // Legacy upload methods
  const uploadLogo = useCallback(async (file: File): Promise<{ url: string; fileName: string }> => {
    return uploadFile(file, 'logo');
  }, [uploadFile]);

  const uploadAtta = useCallback(async (file: File): Promise<{ url: string; fileName: string }> => {
    return uploadFile(file, 'atta');
  }, [uploadFile]);

  const uploadDbd = useCallback(async (file: File): Promise<{ url: string; fileName: string }> => {
    return uploadFile(file, 'dbd');
  }, [uploadFile]);

  const uploadPp20 = useCallback(async (file: File): Promise<{ url: string; fileName: string }> => {
    return uploadFile(file, 'pp20');
  }, [uploadFile]);

  // Computed values
  const currentProfile = useMemo(() => {
    return profile || profileData || null;
  }, [profile, profileData]);

  const currentProfileStatus = useMemo(() => {
    return profileStatus || checkData || null;
  }, [profileStatus, checkData]);

  const hasProfile = useMemo(() => {
    return currentProfileStatus?.has_profile || false;
  }, [currentProfileStatus]);

  const companyName = useMemo(() => {
    return currentProfile?.name || currentProfileStatus?.company_name;
  }, [currentProfile, currentProfileStatus]);

  const isLoading = useMemo(() => {
    return isCheckingProfile || isLoadingProfile;
  }, [isCheckingProfile, isLoadingProfile]);

  const isUploadingFile = useMemo(() => {
    return logoUploadMutation.isPending ||
           attaUploadMutation.isPending ||
           dbdUploadMutation.isPending ||
           pp20UploadMutation.isPending;
  }, [logoUploadMutation.isPending, attaUploadMutation.isPending, dbdUploadMutation.isPending, pp20UploadMutation.isPending]);

  const canCreateProfile = useMemo(() => {
    return isInitialized && !hasProfile && !createProfileMutation.isPending;
  }, [isInitialized, hasProfile, createProfileMutation.isPending]);

  const canUpdateProfile = useMemo(() => {
    return Boolean(hasProfile && currentProfile && !updateProfileMutation.isPending);
  }, [hasProfile, currentProfile, updateProfileMutation.isPending]);

  return {
    // State
    profile: currentProfile,
    profileStatus: currentProfileStatus,
    isInitialized,
    lastUpdated,

    // Loading states
    isCheckingProfile,
    isLoadingProfile,
    isCreatingProfile: createProfileMutation.isPending,
    isUpdatingProfile: updateProfileMutation.isPending,
    isUploadingFile,

    // Actions
    createProfile,
    updateProfile,
    uploadFile,

    // Legacy upload methods
    uploadLogo,
    uploadAtta,
    uploadDbd,
    uploadPp20,

    // State management
    setProfile,
    setProfileStatus,
    clearProfileData,

    // Queries
    refetchProfileCheck,
    refetchProfile,

    // Computed
    hasProfile,
    companyName,
    isLoading,
    canCreateProfile,
    canUpdateProfile,

    // Error states
    checkError,
    profileError,
  };
};

// Additional utility hook for file upload with custom options
export const useAgentFileUpload = (
  fileType: 'logo' | 'atta' | 'dbd' | 'pp20',
  options: {
    onSuccess?: (result: { url: string; fileName: string }) => void;
    onError?: (error: Error) => void;
    showToast?: boolean;
  } = {}
) => {
  const { showToast = true, onSuccess, onError } = options;
  const { toast } = useToast();

  return useUploadFileMutation(fileType, {
    onSuccess: (result) => {
      if (showToast) {
        toast({
          title: "สำเร็จ",
          description: AGENT_PROFILE_MESSAGES.FILE_UPLOADED,
        });
      }
      onSuccess?.(result);
    },
    onError: (error) => {
      if (showToast) {
        toast({
          title: "เกิดข้อผิดพลาด",
          description: error.message,
          variant: "destructive"
        });
      }
      onError?.(error);
    }
  });
};

export default useAgentProfile;