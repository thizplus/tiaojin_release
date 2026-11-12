export interface DailyStatusReport {
  date: string
  pending: number
  processing: number
  confirm: number
  cancel: number
  completed: number
  expired: number
  total_seats: number
}

export interface AdminReportData {
  total_bookings: number
  total_seats: number
  total_sales: number
  actual_revenue: number
  daily_status_report: DailyStatusReport[]
}

export interface AdminReportResponse {
  success: boolean
  message: string
  data: AdminReportData
}

export interface AdminReportParams {
  year: number
  month: number
}

// Sales Report Types
export interface SalesReportData {
  total_bookings: number
  total_seats: number
  total_sales: number
  actual_revenue: number
  daily_status_report: DailyStatusReport[]
}

export interface SalesReportResponse {
  success: boolean
  message: string
  data: SalesReportData
}

export interface SalesReportParams {
  year: number
  month: number
}

// Agent Report Types
export interface AgentReportData {
  total_bookings: number
  total_seats: number
  total_sales: number
  actual_revenue: number
  daily_status_report: DailyStatusReport[]
}

export interface AgentReportResponse {
  success: boolean
  message: string
  data: AgentReportData
}

export interface AgentReportParams {
  year: number
  month: number
}