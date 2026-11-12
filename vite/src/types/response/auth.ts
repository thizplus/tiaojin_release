// Auth Response Types
import type { User, UserRole, UserStatus } from '../models/user';

// Authentication Responses
export interface LoginResponse {
  user_id: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface RegisterResponse {
  user: User;
  message: string;
}

export interface TokenValidationResponse {
  status: boolean;
  decoded: {
    valid: boolean;
    user_id: string;
    email: string;
    role: UserRole;
    exp?: number;
    iat?: number;
  };
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

// Profile and Company Responses
export interface AuthUserProfile extends User {
  company_profile?: CompanyProfileResponse;
}

export interface CompanyProfileResponse {
  id: string;
  user_id: string;
  name: string;
  tax_id: string;
  phone: string;
  address: string;
  province: string;
  zipcode: string;
  atta: string;
  dbd: string;
  pp20: string;
  logo: string | null;
  pdf_header_image: string | null;
  created_at: string;
  updated_at: string;
}

// API Key Response (if needed)
export interface ApiKeyResponse {
  api_key: string;
  expires_at: string;
  created_at: string;
}