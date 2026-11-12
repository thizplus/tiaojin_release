// System Configuration Types

export interface ExpirationTier {
  name: string;
  description: string;
  threshold_days: number | null;
  expire_minutes?: number | null;
  expire_hours?: number | null;
  expire_days?: number | null;
  sort_order: number;
  example_message?: string;
}

export interface ExpirationTiersData {
  tiers: ExpirationTier[];
  metadata: {
    version: string;
    last_updated: string;
    updated_by: string;
  };
}

export interface ExpirationTiersRequest {
  tiers: Omit<ExpirationTier, 'example_message'>[];
}

export interface ExpirationValidationResult {
  is_valid: boolean;
  errors: string[];
  warnings: string[];
  total_tiers: number;
  default_tier: string | null;
}

export interface TierCalculationRequest {
  tour_start_date: string;
  created_at: string;
}

export interface TierCalculationResult {
  applicable_tier: ExpirationTier;
  days_until_tour: number;
  expiration_time: string;
  time_until_expiry: string;
  is_default_tier: boolean;
  calculation_details: {
    tour_start_date: string;
    booking_created_at: string;
    current_time: string;
    days_until_tour: number;
    tested_tiers: Array<{
      tier_name: string;
      threshold_days: number | null;
      matches: boolean;
      reason: string;
    }>;
    selected_tier_name: string;
    expiration_method: 'hours' | 'days';
  };
}

export interface ExpirationStats {
  total_bookings_with_expiration: number;
  expired_bookings: number;
  near_expiring_bookings: number;
  bookings_by_tier: Record<string, number>;
  average_expiration_time: string;
  expiration_trends: Array<{
    date: string;
    created: number;
    expired: number;
    near_expiring: number;
  }>;
}

export interface ProcessExpiredRequest {
  dry_run?: boolean;
  limit?: number;
}

export interface ProcessedBooking {
  booking_id: string;
  invoice_id: string;
  customer_name: string;
  expires_at: string;
  processed_at: string;
  action: string;
  reason: string;
  success: boolean;
  error: string | null;
}

export interface ProcessExpiredResult {
  processed_at: string;
  dry_run: boolean;
  total_checked: number;
  expired_found: number;
  processed_count: number;
  failed_count: number;
  processed_bookings: ProcessedBooking[];
  errors: string[];
}

export interface SystemConfig {
  id: string;
  config_key: string;
  config_value: any;
  config_type: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ConfigListRequest {
  config_type?: string;
  limit?: number;
  offset?: number;
}

export interface ConfigListResponse {
  configs: SystemConfig[];
  pagination: {
    total: number;
    has_more: boolean;
    current_page: number;
    total_pages: number;
    page_size: number;
    limit: number;
    offset: number;
  };
}

export interface CreateConfigRequest {
  config_key: string;
  config_value: any;
  config_type: string;
  description: string;
  is_active: boolean;
}

export interface UpdateConfigRequest {
  config_value: any;
  description: string;
  is_active: boolean;
}

export interface CreateConfigResponse {
  id: string;
}

// API Response Types
export interface ConfigApiResponse<T = any> {
  success: boolean;
  message: string;
  data: T | null;
}

export interface ConfigApiInfo {
  message: string;
  version: string;
  expiration_tiers: string;
  configs: string;
  stats: string;
}

// Booking Limits Configuration Types
export interface BookingLimitsConfig {
  agent_global_bypasses: string[]; // Array of user UUIDs that bypass limits
  global_settings: {
    count_all_statuses: string[]; // Array of booking statuses to count
    default_limit_per_period: number; // Default limit per period
    enabled: boolean; // Whether booking limits are enabled
  };
}

export interface BookingLimitsRequest {
  agent_global_bypasses: string[];
  global_settings: {
    count_all_statuses: string[];
    default_limit_per_period: number;
    enabled: boolean;
  };
}

// Config Types for different modules
export type ConfigType =
  | 'expiration'
  | 'booking_limits'
  | 'payment'
  | 'notification'
  | 'security'
  | 'integration'
  | 'general';

// Error types
export interface ConfigError {
  field?: string;
  message: string;
  code?: string;
}