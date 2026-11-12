export type UserRole = 'super' | 'admin' | 'agent' | 'sales';
export type UserStatus = 'pending' | 'active' | 'inactive';

import type { CompanyProfile } from './companyProfile';

export interface User {
  id: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  created_at: string;
  updated_at: string;
  company_profile?: CompanyProfile | null;
}

export interface UserStats {
  total_users: number;
  active_users: number;
  pending_users: number;
  inactive_users: number;
}