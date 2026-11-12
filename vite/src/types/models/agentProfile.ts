

export interface AgentProfile {
  id: string;
  user_id: string;
  name: string;
  tax_id: string;
  phone: string;
  address: string;
  province: string;
  zipcode: string;
  atta?: string | null;
  dbd?: string | null;
  pp20?: string | null;
  logo?: string | null;
  pdfHeaderImage?: string | null;
  created_at: string;
  updated_at: string;
}

// Create Profile Request
export interface CreateAgentProfileRequest {
  name: string;
  tax_id: string;
  phone: string;
  address: string;
  province: string;
  zipcode: string;
  atta?: string;
  dbd?: string;
  pp20?: string;
  logo?: string;
  pdfHeaderImage?: string;
}

export interface UpdateAgentProfileRequest {
  name?: string;
  phone?: string;
  address?: string;
  province?: string;
  zipcode?: string;
  atta?: string;
  dbd?: string;
  pp20?: string;
  logo?: string;
  pdfHeaderImage?: string;
}

export interface ProfileCheckResponse {
  has_profile: boolean;
  profile_id?: string;
  company_name?: string;
}

export interface AgentProfileFilter {
  name?: string;
  tax_id?: string;
  status?: string;
  // For future extensions if needed
}