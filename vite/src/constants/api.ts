const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const AUTH_API = {
  // Public endpoints
  REGISTER: '/auth/register',
  LOGIN: '/auth/login',
  REFRESH: '/auth/refresh',
  VALIDATE: '/auth/validate',
  RESET_PASSWORD: '/auth/reset-password',

  // Protected endpoints
  LOGOUT: '/auth/logout',
  UPDATE_PASSWORD: '/auth/update-password',
  REVOKE_ALL: '/auth/revoke-all',
}

export const USER_API = {
  // Public endpoints
  BY_EMAIL: '/users/by-email',

  // Protected endpoints (all users)
  PROFILE: '/users/profile',

  // Admin+ endpoints
  LIST: '/admin/users',
  CREATE: '/admin/users',
  STATS: '/admin/users/stats',
  GET_BY_ID: (id: string) => `/admin/users/${id}`,
  UPDATE_BY_ID: (id: string) => `/admin/users/${id}`,
  UPDATE_STATUS: (id: string) => `/admin/users/${id}/status`,
  UPDATE_ROLE: (id: string) => `/admin/users/${id}/role`,
  DELETE: (id: string) => `/admin/users/${id}`,

  // Legacy endpoints (keep for compatibility)
  APPROVE: (id: string) => `/users/${id}/approve`,
  SUSPEND: (id: string) => `/users/${id}/suspend`,
  REACTIVATE: (id: string) => `/users/${id}/reactivate`,
  CHANGE_ROLE: (id: string) => `/users/${id}/change-role`,
}

export const CITY_API = {
  // Public endpoints
  PUBLIC_LIST: '/cities',                                   // GET /api/v1/cities
  PUBLIC_GET_BY_ID: (id: string) => `/cities/${id}`,      // GET /api/v1/cities/:id

  // Admin endpoints (ตาม backend routes)
  LIST: '/admin/cities',                                    // GET /api/v1/admin/cities
  CREATE: '/admin/cities',                                  // POST /api/v1/admin/cities
  UPDATE: (id: string) => `/admin/cities/${id}`,           // PUT /api/v1/admin/cities/:id
  DELETE: (id: string) => `/admin/cities/${id}`,           // DELETE /api/v1/admin/cities/:id
  UPDATE_STATUS: (id: string) => `/admin/cities/${id}/status`, // PUT /api/v1/admin/cities/:id/status

  GET_BY_ID: (id: string) => `/admin/cities/${id}`,        // GET /api/v1/admin/cities/:id
}

export const BANNER_API = {
  // Public endpoints
  PUBLIC_LIST: '/banners',                                  // GET /api/v1/banners

  // Admin endpoints (ตาม backend routes)
  LIST: '/admin/banners',                                   // GET /api/v1/admin/banners
  CREATE: '/admin/banners',                                 // POST /api/v1/admin/banners
  UPDATE: (id: string) => `/admin/banners/${id}`,          // PUT /api/v1/admin/banners/:id
  DELETE: (id: string) => `/admin/banners/${id}`,          // DELETE /api/v1/admin/banners/:id
  UPDATE_STATUS: (id: string) => `/admin/banners/${id}/status`, // PUT /api/v1/admin/banners/:id/status
  UPDATE_ORDER: (id: string) => `/admin/banners/${id}/order`,   // PUT /api/v1/admin/banners/:id/order

  GET_BY_ID: (id: string) => `/admin/banners/${id}`,       // GET /api/v1/admin/banners/:id
}

export const SALES_API = {
  // Admin endpoints (ตาม backend routes)
  LIST: '/admin/sales',                                      // GET /api/v1/admin/sales
  CREATE: '/admin/sales',                                    // POST /api/v1/admin/sales
  UPDATE: (id: string) => `/admin/sales/${id}`,             // PUT /api/v1/admin/sales/:id
  DELETE: (id: string) => `/admin/sales/${id}`,             // DELETE /api/v1/admin/sales/:id
  ACTIVATE: (id: string) => `/admin/sales/${id}/activate`,  // PATCH /api/v1/admin/sales/:id/activate
  DEACTIVATE: (id: string) => `/admin/sales/${id}/deactivate`, // PATCH /api/v1/admin/sales/:id/deactivate

  GET_BY_ID: (id: string) => `/admin/sales/${id}`,          // GET /api/v1/admin/sales/:id

  //public endpoints
  GET_SALE_ACTIVE: '/sales/active',                                      // GET /api/v1/sales/active
}

export const TOUR_API = {
  // Admin endpoints (ตาม backend routes)
  LIST: '/admin/tours',                                      // GET /api/v1/admin/tours
  CREATE: '/admin/tours',                                    // POST /api/v1/admin/tours
  UPDATE: (id: string) => `/admin/tours/${id}`,             // PUT /api/v1/admin/tours/:id
  DELETE: (id: string) => `/admin/tours/${id}`,             // DELETE /api/v1/admin/tours/:id
  UPDATE_STATUS: (id: string) => `/admin/tours/${id}/status`, // PATCH /api/v1/admin/tours/:id/status

  GET_BY_ID: (id: string) => `/admin/tours/${id}`,          // GET /api/v1/admin/tours/:id

  // Tour periods endpoints
  GET_PERIODS: (id: string) => `/admin/tours/${id}/periods`, // GET /api/v1/admin/tours/:id/periods

  // Public tours endpoints
  REGULAR: '/tours/regular',                                 // GET /api/v1/tours/regular
  SPECIAL: '/tours/special',                                 // GET /api/v1/tours/special
}

export const PERIOD_API = {
  // Admin period endpoints (ตาม backend routes)
  LIST: '/admin/periods',                                      // GET /api/v1/admin/periods
  CREATE: '/admin/periods',                                    // POST /api/v1/admin/periods
  UPDATE: (id: string) => `/admin/periods/${id}`,             // PUT /api/v1/admin/periods/:id
  DELETE: (id: string) => `/admin/periods/${id}`,             // DELETE /api/v1/admin/periods/:id
  UPDATE_STATUS: (id: string) => `/admin/periods/${id}/status`, // PUT /api/v1/admin/periods/:id/status
  CHECK_AVAILABILITY: (id: string) => `/admin/periods/${id}/availability`, // GET /api/v1/admin/periods/:id/availability

  GET_BY_ID: (id: string) => `/admin/periods/${id}`,          // GET /api/v1/admin/periods/:id

  // Seat Hold endpoints
  HOLD_SEATS: (id: string) => `/admin/periods/${id}/hold-seats`,        // POST /api/v1/admin/periods/:id/hold-seats
  GET_SEAT_HOLDS: (id: string) => `/admin/periods/${id}/seat-holds`,    // GET /api/v1/admin/periods/:id/seat-holds
}

export const SEAT_HOLD_API = {
  // Seat Hold management endpoints
  RELEASE: (holdId: string) => `/admin/seat-holds/${holdId}/release`,   // DELETE /api/v1/admin/seat-holds/:hold_id/release
}

export const AGENT_PROFILE_API = {
  // Agent profile endpoints
  PROFILE: '/agent/profile',                           // GET, PUT /api/v1/agent/profile
  CREATE_PROFILE: '/agent/profile',                   // POST /api/v1/agent/profile
  CHECK_PROFILE: '/agent/profile/check',              // GET /api/v1/agent/profile/check
}

export const COMPANY_PROFILE_API = {
  // Admin company profile endpoints
  LIST: '/admin/companies',                              // GET /api/v1/admin/companies
  SEARCH: '/admin/companies/search',                     // GET /api/v1/admin/companies/search
  GET_BY_ID: (id: string) => `/admin/companies/${id}`,  // GET /api/v1/admin/companies/:id
  GET_BY_TAX_ID: (taxId: string) => `/admin/companies/tax-id/${taxId}`, // GET /api/v1/admin/companies/tax-id/:tax_id
  UPDATE: (id: string) => `/admin/companies/${id}`,           // PUT /api/v1/admin/companies/:id
  APPROVE: (id: string) => `/admin/companies/${id}/approve`, // PATCH /api/v1/admin/companies/:id/approve
  // Public
  RANDOM_LOGOS: '/company-profiles/logos',               // GET /api/v1/company-profiles/logos
}

export const BOOKING_API = {
  // Booking endpoints (Agent+ access level)
  CREATE: '/bookings',                                   // POST /api/v1/bookings
  LIST: '/bookings',                                     // GET /api/v1/bookings
  GET_BY_ID: (id: string) => `/bookings/${id}`,         // GET /api/v1/bookings/:id
  UPDATE: (id: string) => `/bookings/${id}`,            // PUT /api/v1/bookings/:id
  DELETE: (id: string) => `/bookings/${id}`,            // DELETE /api/v1/bookings/:id
  UPDATE_STATUS: (id: string) => `/bookings/${id}/status`, // PUT /api/v1/bookings/:id/status
}

export const ADMIN_BOOKING_API = {
  // Admin booking endpoints
  LIST: '/admin/bookings',                               // GET /api/v1/admin/bookings
  STATS: '/admin/bookings/stats',                        // GET /api/v1/admin/bookings/stats
  REPORT: '/bookings/report',                            // GET /api/v1/bookings/report
  GET_BY_ID: (id: string) => `/admin/bookings/${id}`,   // GET /api/v1/admin/bookings/:id
  UPDATE_STATUS: (id: string) => `/admin/bookings/${id}/status`, // PUT /api/v1/admin/bookings/:id/status
  CREATE_PAYMENT: (bookingId: string) => `/admin/bookings/${bookingId}/payments`, // POST /api/v1/admin/bookings/:bookingId/payments
  UPLOAD_PASSPORT: (bookingId: string) => `/admin/bookings/${bookingId}/passport`, // POST /api/v1/admin/bookings/:bookingId/passport

  // Payment slip verification endpoints
  VERIFY_PAYMENT_SLIP: (paymentSlipId: string) => `/admin/payment-slips/${paymentSlipId}/verify`, // POST /api/v1/admin/payment-slips/:paymentSlipId/verify
  REVERSE_PAYMENT_SLIP: (paymentSlipId: string) => `/admin/payment-slips/${paymentSlipId}/reverse`, // POST /api/v1/admin/payment-slips/:paymentSlipId/reverse
}

export const AGENT_BOOKING_API = {
  // Agent booking endpoints (Agent access to own bookings)
  LIST: '/agent/bookings',                               // GET /api/v1/agent/bookings
  GET_BY_ID: (id: string) => `/agent/bookings/${id}`,   // GET /api/v1/agent/bookings/:id

  // Agent booking actions
  CANCEL: (id: string) => `/agent/bookings/${id}`,      // DELETE /api/v1/agent/bookings/:id
  UPLOAD_PAYMENT_SLIP: (id: string) => `/agent/bookings/${id}/payment-slip`, // POST /api/v1/agent/bookings/:id/payment-slip
  UPLOAD_PASSPORT: (id: string) => `/agent/bookings/${id}/passport`,         // POST /api/v1/agent/bookings/:id/passport
  GET_DOCUMENTS: (id: string) => `/agent/bookings/${id}/documents`,          // GET /api/v1/agent/bookings/:id/documents
}

export const SALE_BOOKING_API = {
  // Sale booking endpoints (Sale access to own bookings)
  LIST: '/sales/bookings',                               // GET /api/v1/sales/bookings
  REPORT: '/bookings/sales-report',                      // GET /api/v1/bookings/sales-report
  GET_BY_ID: (id: string) => `/sales/bookings/${id}`,   // GET /api/v1/sales/bookings/:id

  // Sale booking actions (for assigned customers)
  UPLOAD_PAYMENT_SLIP: (id: string) => `/sales/bookings/${id}/payment-slip`, // POST /api/v1/sales/bookings/:id/payment-slip
  UPLOAD_PASSPORT: (id: string) => `/sales/bookings/${id}/passport`,         // POST /api/v1/sales/bookings/:id/passport
}

export const AGENT_REPORT_API = {
  // Agent report endpoints
  REPORT: '/bookings/agent-report',                      // GET /api/v1/bookings/agent-report
}

export const API_KEY_API = {
  // API Key endpoints
  LIST: '/api-keys',                                     // GET /api/v1/api-keys
  CREATE: '/api-keys',                                   // POST /api/v1/api-keys
  GET_BY_ID: (id: string) => `/api-keys/${id}`,        // GET /api/v1/api-keys/:id
  UPDATE: (id: string) => `/api-keys/${id}`,           // PUT /api/v1/api-keys/:id
  DELETE: (id: string) => `/api-keys/${id}`,           // DELETE /api/v1/api-keys/:id
  VALIDATE: '/api-keys/validate',                       // POST /api/v1/api-keys/validate
}

export const FILE_API = {
  UPLOAD: '/files/upload',
}

export const CONFIG_API = {
  // Base info endpoint
  INFO: '/config',                                               // GET /api/v1/config

  // Expiration tiers management
  EXPIRATION_TIERS: '/config/expiration/tiers',                 // GET, PUT /api/v1/config/expiration/tiers
  VALIDATE_TIERS: '/config/expiration/tiers/validate',          // POST /api/v1/config/expiration/tiers/validate
  CALCULATE_TIER: '/config/expiration/tiers/calculate',         // POST /api/v1/config/expiration/tiers/calculate
  EXPIRATION_STATS: '/config/expiration/stats',                // GET /api/v1/config/expiration/stats
  PROCESS_EXPIRED: '/config/expiration/process-expired',       // POST /api/v1/config/expiration/process-expired

  // Agent booking limits management
  AGENT_BOOKING_LIMITS: '/config/agent-booking-limits',        // GET, PUT /api/v1/config/agent-booking-limits

  // General config management
  CONFIGS: '/config/configs',                                   // GET, POST /api/v1/config/configs
  CONFIG_BY_ID: (id: string) => `/config/configs/${id}`,      // GET, PUT, DELETE /api/v1/config/configs/:id
}

export const BANK_ACCOUNT_API = {
  LIST: '/bank-accounts/',                                        // GET /api/v1/bank-accounts/ (admin)
  ACTIVE: '/bank-accounts/active',                                // GET /api/v1/bank-accounts/active (public)
  CREATE: '/bank-accounts/',                                      // POST /api/v1/bank-accounts/
  GET_BY_ID: (id: string) => `/bank-accounts/${id}`,             // GET /api/v1/bank-accounts/:id
  UPDATE: (id: string) => `/bank-accounts/${id}`,                // PUT /api/v1/bank-accounts/:id
  DELETE: (id: string) => `/bank-accounts/${id}`,                // DELETE /api/v1/bank-accounts/:id
}

export const API = {
    AUTH: AUTH_API,
    USER: USER_API,
    CITY: CITY_API,
    BANNER: BANNER_API,
    TOUR: TOUR_API,
    PERIOD: PERIOD_API,
    SEAT_HOLD: SEAT_HOLD_API,
    SALES: SALES_API,
    AGENT_PROFILE: AGENT_PROFILE_API,
    COMPANY_PROFILE: COMPANY_PROFILE_API,
    BOOKING: BOOKING_API,
    ADMIN_BOOKING: ADMIN_BOOKING_API,
    AGENT_BOOKING: AGENT_BOOKING_API,
    SALE_BOOKING: SALE_BOOKING_API,
    AGENT_REPORT: AGENT_REPORT_API,
    API_KEY: API_KEY_API,
    BANK_ACCOUNT: BANK_ACCOUNT_API,

    // File API
    FILE: FILE_API,

    // Config API
    CONFIG: CONFIG_API,

    BASE_URL: API_BASE_URL
  };
  
  export default API;