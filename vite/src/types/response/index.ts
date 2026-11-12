// Response Types Export
export * from './auth';
export * from './user';
export * from './companyProfile';

// Common pagination response type
export interface PaginationResponse {
  total: number;
  offset: number;
  limit: number;
  has_next: boolean;
}