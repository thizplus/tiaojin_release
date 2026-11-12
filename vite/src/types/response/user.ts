// User Response Types
import type { UserRole, UserStatus } from '../models/user';
import type { CompanyProfile } from '../models/companyProfile';

// Base User Model
export interface UserBase {
  id: string;
  email: string;
  username?: string;
  role: UserRole;
  status: UserStatus;
  created_at: string;
  updated_at: string;
}

// Extended User with Profile Info
export interface UserProfile extends UserBase {
  first_name?: string;
  last_name?: string;
  phone?: string;
  company_profile?: CompanyProfile;
}

// User List Item (for admin management)
export interface UserListItem extends UserBase {
  company_profile?: CompanyProfile;
}

// Public User Lookup Response
export interface UserByEmailResponse {
  status: boolean;
  data?: UserProfile;
}

// User Stats Response (Admin+)
export interface UserStatsResponse {
  total_users: number;
  active_users: number;
  pending_users: number;
  inactive_users: number;
}

// User List Response (Admin+)
export interface UsersListResponse {
  users: UserListItem[];
  total: number;
  has_more: boolean;
}

// Profile Update Response
export interface UpdateProfileResponse {
  user_id: string;
  email: string;
  role: UserRole;
  updated_at: string;
}

// User Management Action Responses
export interface UserActionResponse extends UserBase {
  action_performed?: string;
  performed_at?: string;
}

export interface ChangeRoleResponse {
  id: string;
  new_role: UserRole;
  changed_at: string;
}

// Delete User Response
export interface DeleteUserResponse {
  deleted: boolean;
  id: string;
  deleted_at: string;
}