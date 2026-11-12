// Booking Types based on API analysis

import type { CompanyProfileResponse } from "./response";

// Request Types
export interface CreateBookingRequest {
  tour_id: string;
  period_id: string;
  sales_id: string;
  customer_name: string;
  customer_phone: string;
  note?: string;
  room_types?: CreateBookingRoomTypeRequest[];
  rate_types: CreateBookingRateTypeRequest[]; // required, min 1
  additional_services?: CreateBookingAdditionalServiceRequest[];

  // VAT and Withholding Tax options (Agent เลือกตอนจอง)
  include_vat?: boolean; // รวม VAT 7% (default: false)
  include_wht?: boolean; // หัก ณ ที่จ่าย 3% (default: false)
}

export interface CreateBookingRoomTypeRequest {
  id: string;
  quantity: number; // min 1
}

export interface CreateBookingRateTypeRequest {
  id: string;
  quantity: number; // min 1
}

export interface CreateBookingAdditionalServiceRequest {
  id: string;
  quantity: number; // min 1
}

// Response Types
export interface BookingCreateResponse {
  success: boolean;
  message: string;
  data: BookingResponse;
}

export interface BookingResponse {
  id: string;
  invoice_id: string;
  user_id: string;
  tour_id: string;
  period_id: string;
  sales_id: string;
  customer_name: string;
  customer_phone: string;
  note: string | null;
  sum_room: number;
  sum_rate: number;
  sum_additional_services: number;
  sum_qty: number;
  sum_commission_qty: number;
  sum_net: number;
  sum_commission: number;
  discount: number;
  sum_deposit: number;
  sum_payment: number;

  // VAT fields (ภาษีมูลค่าเพิ่ม)
  include_vat: boolean;
  vat_rate: number;
  vat_amount: number;

  // Withholding Tax fields (หัก ณ ที่จ่าย)
  include_wht: boolean;
  wht_rate: number;
  wht_amount: number;

  // Calculated totals
  net_before_tax: number;
  grand_total: number;

  date_deposit: string; // YYYY-MM-DD
  date_final: string; // YYYY-MM-DD
  status: "pending" | "processing" | "confirm" | "cancel" | "completed";
  expires_at: string | null; // ISO timestamp
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp

  // Relations (optional)
  tour_details?: TourResponse;
  period_details?: PeriodSummaryResponse;
  sales_details?: SalesResponse;
  company_profile?: CompanyProfileResponse;
  user_details?: UserDetailsResponse;
  room_types?: BookingRoomTypeResponse[];
  rate_types?: BookingRateTypeResponse[];
  additional_services?: BookingAdditionalServiceResponse[];
  payment_slips?: PaymentSlipResponse[];
  passports?: PassengerPassportResponse[];
}

// Related Response Types
export interface TourResponse {
  id: string;
  title: string;
  code: string;
  city_id: string;
  detail: string;
  thumbnail: string;
  file_pdf: string;
  file_word: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface PeriodSummaryResponse {
  id: string;
  start_date: string;
  end_date: string;
  seat: number;
  booking: number;
  held_seats: number;
  available_seats: number;
  deposit: number;
  deposit_type?: string;
  files?: PeriodFile[];
}

interface PeriodFile {
  name: string;
  type: string;
  url: string;
}

export interface SalesResponse {
  id: string;
  user_id: string;
  code: string;
  first_name: string;
  last_name: string;
  nick_name: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BookingRoomTypeResponse {
  id: string;
  title: string;
  price: number;
  quantity: number;
  total: number;
}

export interface BookingRateTypeResponse {
  id: string;
  title: string;
  price: number;
  quantity: number;
  exclude: boolean;
  total: number;
}

export interface BookingAdditionalServiceResponse {
  id: string;
  title: string;
  price: number;
  quantity: number;
  rate_type: string;
  unit_name: string;
  total: number;
}

export interface PaymentSlipResponse {
  id: string;
  booking_id: string;
  slip_id: string;
  image: string;
  amount: number;
  status: "pending" | "confirm" | "reject" | "reversed";
  created_at: string;

  // SlipOK verification fields
  payment_date?: string;
  sending_bank?: string;
  receiving_bank?: string;
  receiver_account?: string;
  receiver_name?: string;
  sender_name?: string;
  trans_ref?: string;
  slip_verified?: boolean;
  slip_verify_error?: string;
}

export interface PassengerPassportResponse {
  id: string;
  booking_id: string;
  passport_number: string;
  image: string;
  status: "pending" | "approved" | "rejected";
  title_th?: string;
  name_th?: string;
  surname_th?: string;
  title_en?: string;
  name_en?: string;
  surname_en?: string;
  gender?: string;
  nationality?: string;
  date_of_birth?: string;
  id_number?: string;
  place_of_birth?: string;
  date_of_expiry?: string;
  created_at: string;
}



export interface UserDetailsResponse {
  id: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
  updated_at: string;
}

// Error Response Types
export interface BookingValidationError {
  success: false;
  message: string;
  errors: {
    [key: string]: string[];
  };
}

export interface BookingBusinessLogicError {
  success: false;
  message: string;
  error: string;
}

// Admin Booking Types
export interface AdminBookingListResponse {
  success: boolean;
  message: string;
  data: {
    bookings: BookingResponse[];
    pagination: {
      total: number;
      current_page: number;
      page_size: number;
      total_pages: number;
    };
  };
}

export interface AdminBookingListFilters {
  status?: string; // รองรับ single หรือ multiple statuses (comma-separated): "pending" หรือ "pending,processing,confirm"
  payment_slip_status?: "pending" | "confirm" | "reject" | "reversed";
  user_id?: string;
  tour_id?: string;
  period_id?: string;
  sales_id?: string;
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
  search?: string;
  page?: number;
  limit?: number;
  sort?: string;
  order?: "asc" | "desc";
}

export interface AdminBookingStatsResponse {
  success: boolean;
  message: string;
  data: {
    total_bookings: number;
    total_revenue: number;
    total_commission: number;
    total_deposit: number;
    total_payment: number;
    status_counts: {
      pending: number;
      processing: number;
      confirm: number;
      cancel: number;
      completed: number;
    };
    monthly_stats: {
      month: string; // YYYY-MM
      bookings: number;
      revenue: number;
      commission: number;
    }[];
  };
}

export interface AdminBookingStatusUpdateRequest {
  status: "pending" | "processing" | "confirm" | "cancel" | "completed";
  note?: string;
}

export interface AdminBookingStatusUpdateResponse {
  success: boolean;
  message: string;
  data: BookingResponse;
}

export interface AdminCreatePaymentRequest {
  amount: number;
  payment_method: "cash" | "bank_transfer" | "credit_card" | "online";
  reference_number?: string;
  note?: string;
  payment_date?: string; // YYYY-MM-DD
  slip_image?: string;
}

export interface AdminCreatePaymentResponse {
  success: boolean;
  message: string;
  data: {
    id: string;
    booking_id: string;
    amount: number;
    payment_method: string;
    reference_number: string | null;
    note: string | null;
    payment_date: string;
    created_at: string;
    updated_at: string;
  };
}

// Agent Booking Action Types
export interface AgentCancelBookingRequest {
  reason?: string;
}

export interface AgentCancelBookingResponse {
  success: boolean;
  message: string;
  data: {
    booking_id: string;
    status: "cancel";
    reason: string;
  };
}

export interface AgentUploadPaymentSlipRequest {
  payment_slip: File;
  amount: number;
}

export interface AgentUploadPaymentSlipResponse {
  success: boolean;
  message: string;
  data: {
    payment_slip: PaymentSlipResponse;
  };
}

export interface AgentUploadPassportRequest {
  passport: File;
}

export interface AgentUploadPassportResponse {
  success: boolean;
  message: string;
  data: {
    passport: PassengerPassportResponse;
  };
}

// Document types for agent booking documents API
export interface BookingDocumentResponse {
  id: string;
  filename: string;
  type: "payment_slip" | "passport";
  uploaded_at: string;
  file_url: string;
  file_size: number;
  status: string;
  amount?: number; // Only for payment_slip type
}

export interface AgentBookingDocumentsResponse {
  success: boolean;
  message: string;
  data: {
    booking_id: string;
    documents: BookingDocumentResponse[];
  };
}

// Payment Slip Verification Types
export interface PaymentSlipVerifyRequest {
  status: "confirm" | "reject";
  note?: string;
  actual_amount?: number;
  adjustment_reason?: string;
}

export interface PaymentSlipReverseRequest {
  reversal_reason: string; // 10-500 characters
}

export interface PaymentSlipVerifyResponse {
  success: boolean;
  message: string;
  data: {
    payment_slip: {
      id: string;
      booking_id: string;
      slip_id: string;
      amount: number;
      original_amount?: number;
      adjustment_reason?: string;
      status: "confirm" | "reject";
      payment_method: string;
      is_admin_created: boolean;
      is_reversed: boolean;
      created_at: string;
      updated_at: string;
    };
    booking_update: {
      booking_id: string;
      status: string;
      sum_payment: number;
      sum_deposit: number;
      sum_net: number;
      payment_ratio: number;
      deposit_ratio: number;
      fully_paid: boolean;
      deposit_satisfied: boolean;
      remaining_amount: number;
      updated_at: string;
    };
    amount_adjustment?: {
      original_amount: number;
      adjusted_amount: number;
      adjustment_reason: string;
      was_adjusted: boolean;
    };
  };
}

export interface PaymentSlipReverseResponse {
  success: boolean;
  message: string;
  data: {
    payment_slip: {
      id: string;
      booking_id: string;
      slip_id: string;
      amount: number;
      status: "reversed";
      is_reversed: true;
      reversed_at: string;
      reversed_by: string;
      reversal_reason: string;
      created_at: string;
      updated_at: string;
    };
    booking_update: {
      booking_id: string;
      status: string;
      sum_payment: number;
      sum_deposit: number;
      sum_net: number;
      remaining_amount: number;
      updated_at: string;
    };
  };
}