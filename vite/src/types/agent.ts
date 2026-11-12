import type {
  AgentProfile,
  CreateAgentProfileRequest,
  UpdateAgentProfileRequest,
  ProfileCheckResponse
} from "./models/agentProfile";

// API Response Types
export interface ProfileCheckApiResponse {
  success: boolean;
  message: string;
  data: ProfileCheckResponse;
}

export interface ProfileApiResponse {
  success: boolean;
  message: string;
  data: AgentProfile;
}

export interface ProfileCreateApiResponse {
  success: boolean;
  message: string;
  data: AgentProfile;
}

export interface ProfileUpdateApiResponse {
  success: boolean;
  message: string;
  data: AgentProfile;
}

// Error Response
export interface ProfileErrorResponse {
  success: false;
  message: string;
  error: string;
  errors?: Record<string, string>;
}

// File Upload Types for Agent Profile
export type AgentFileType = 'logo' | 'atta' | 'dbd' | 'pp20';

export type AgentFileCategory = 'agent_profiles' | 'agent_licenses';

export interface AgentFileConfig {
  category: AgentFileCategory;
  fileType: AgentFileType;
  accept: string;
  maxSize: number;
  allowedTypes: string[];
  description: string;
}

export type AgentProfileFileMapping = {
  [K in AgentFileType]: AgentFileCategory;
};

// Validation Types
export interface ValidationRule {
  required?: boolean;
  maxLength?: number;
  pattern?: RegExp;
  message?: string;
}

export interface ProfileValidationConfig {
  name: ValidationRule;
  tax_id: ValidationRule;
  phone: ValidationRule;
  address: ValidationRule;
  province: ValidationRule;
  zipcode: ValidationRule;
}

// Profile State Types
export type ProfileState = 'loading' | 'no_profile' | 'has_profile' | 'creating' | 'editing' | 'error';

export type ProfileSection = 'overview' | 'basic_info' | 'contact' | 'address' | 'documents' | 'logo';

// Hook Options Types
export interface UseAgentProfileOptions {
  enableCheck?: boolean;
  enableProfile?: boolean;
  autoFetchProfile?: boolean;
  showToast?: boolean;
}

export interface UseAgentFileUploadOptions {
  onSuccess?: (result: { url: string; fileName: string }) => void;
  onError?: (error: Error) => void;
  showToast?: boolean;
}

// Component Props Types
export interface ProfileStatusCardProps {
  profileStatus: ProfileCheckResponse;
  onCreateProfile?: () => void;
  isLoading?: boolean;
  className?: string;
}

export interface ProfileInfoCardProps {
  profile: AgentProfile;
  onEdit: () => void;
  className?: string;
}

export interface CreateProfileFormProps {
  onSubmit: (data: CreateAgentProfileRequest) => Promise<void>;
  onUploadFile?: (file: File, category: string, fileType: string) => Promise<{ url: string; fileName: string }>;
  isLoading?: boolean;
  isUploadingFile?: boolean;
  className?: string;
}

export interface EditProfileDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  profile: AgentProfile | null;
  onSubmit: (data: UpdateAgentProfileRequest) => Promise<void>;
  onUploadFile?: (file: File, category: string, fileType: string) => Promise<{ url: string; fileName: string }>;
  isLoading?: boolean;
  isUploadingFile?: boolean;
}

// Utility Types
export interface AgentProfileState {
  profile: AgentProfile | null;
  profileStatus: ProfileCheckResponse | null;
  isInitialized: boolean;
  lastUpdated: string | null;
}

export interface ProfileFormData {
  name: string;
  tax_id?: string; // Only for create
  phone: string;
  address: string;
  province: string;
  zipcode: string;
  atta?: string;
  dbd?: string;
  pp20?: string;
  logo?: string;
}

// API Service Types
export interface ProfileServiceValidation {
  isValid: boolean;
  errors: Record<string, string>;
}

export interface FileValidation {
  isValid: boolean;
  error?: string;
}