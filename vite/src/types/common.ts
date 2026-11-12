//types/common.ts

// Base API Response
export interface ApiResponse<T = unknown> {
  success: boolean;
  message: string;
  data: T;
}

// Pagination Meta (ตรงกับ backend response)
export interface PaginationMeta {
  total: number;
  has_more: boolean;
  current_page: number;
  total_pages: number;
  page_size: number;
  limit: number;
  offset: number;
}

// Paginated Response
export interface PaginatedResponse<T> {
  success: boolean;
  message: string;
  data: {
    [key: string]: T[] | PaginationMeta;
    pagination: PaginationMeta;
  };
}

// Pagination Query Parameters
export interface PaginationParams {
  page?: number;
  limit?: number;
}

// Token Types
export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

