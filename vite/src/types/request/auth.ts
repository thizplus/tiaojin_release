// Auth Request Types
import type { UserRole } from '../models/user';
import type { CompanyProfile } from '../models/companyProfile';

// Authentication Requests
export interface RegisterRequest {
  email: string;
  password: string;
  role?: UserRole;
  company_profile?: CompanyProfile;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ResetPasswordRequest {
  email: string;
}

export interface UpdatePasswordRequest {
  id: string;
  old_password: string;
  new_password: string;
}

// Additional Auth Requests
export interface ForgotPasswordRequest {
  email: string;
}

export interface VerifyEmailRequest {
  token: string;
  email: string;
}

export interface ResendVerificationRequest {
  email: string;
}