import type { CompanyProfile } from '../models/companyProfile';
import type { PaginationResponse } from './index';

// Company Profile Response Types

export interface GetCompanyProfilesResponse {
  success: boolean;
  message: string;
  data: {
    profiles: CompanyProfile[];
    pagination: PaginationResponse;
  };
}

export interface SearchCompanyProfilesResponse {
  success: boolean;
  message: string;
  data: {
    profiles: CompanyProfile[];
    pagination: PaginationResponse;
  };
}

export interface GetCompanyProfileResponse {
  success: boolean;
  message: string;
  data: CompanyProfile;
}

export interface ApproveCompanyProfileResponse {
  success: boolean;
  message: string;
  data: {
    id: string;
    user_id: string;
    status: string;
    message: string;
    updated_at: string;
  };
}

export interface CompanyProfileErrorResponse {
  success: false;
  message: string;
}