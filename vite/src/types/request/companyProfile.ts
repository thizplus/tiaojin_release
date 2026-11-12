// Company Profile Request Types

export interface GetCompanyProfilesRequest {
  page?: number;
  limit?: number;
  search?: string;
}

export interface SearchCompanyProfilesRequest {
  q: string;
  page?: number;
  limit?: number;
}

export interface GetCompanyProfileByIdRequest {
  id: string;
}

export interface GetCompanyProfileByTaxIdRequest {
  tax_id: string;
}

export interface ApproveCompanyProfileRequest {
  id: string;
}