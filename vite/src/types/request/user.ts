import type { UserRole, UserStatus } from "../models/user";

// User Request Types
export interface CreateUserRequest {
  email: string;
  password: string;
  role: UserRole;
  status?: UserStatus;
}

export interface UpdateUserRequest {
  email?: string;
  role?: UserRole;
  status?: UserStatus;
}

export interface UpdateUserStatusRequest {
  status: UserStatus;
}

export interface UpdateUserRoleRequest {
  role: UserRole;
}
