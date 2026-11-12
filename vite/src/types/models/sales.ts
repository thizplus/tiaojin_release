export interface Sale {
  id: string;
  user_id: string;
  code: string;
  first_name: string;
  last_name: string;
  nick_name: string;
  active: boolean;
  created_at: string;
  updated_at: string;
  user_email?: string;
  user_role?: string;
}

export interface SaleFilter {
  page?: number;
  limit?: number;
  search?: string;
  active?: boolean;
  sort?: string;
  order?: 'asc' | 'desc';
}

export interface CreateSaleRequest {
  user_id?: string;
  code: string;
  first_name: string;
  last_name: string;
  nick_name: string;
}

export interface UpdateSaleRequest {
  user_id?: string;
  first_name?: string;
  last_name?: string;
  nick_name?: string;
  active?: boolean;
}

export interface SalesActiveResponse {
  sales: Sale[];
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