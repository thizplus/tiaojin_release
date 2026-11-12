import { SALES_API as API } from '@/constants/api';
import apiService from '../apiService';
import type {
  Sale,
  SaleFilter,
  CreateSaleRequest,
  UpdateSaleRequest,
  SalesActiveResponse,

  ApiResponse,
  PaginationMeta
} from '@/types';

const salesService = {
  async getSales(params: SaleFilter = {}): Promise<ApiResponse<{ sales: Sale[], pagination: PaginationMeta }>> {
    const queryParams = new URLSearchParams();

    if (params.search) queryParams.append("search", params.search);
    if (params.active !== undefined) queryParams.append("active", params.active.toString());
    if (params.page !== undefined) queryParams.append("page", params.page.toString());
    if (params.limit !== undefined) queryParams.append("limit", params.limit.toString());
    if (params.sort) queryParams.append("sort", params.sort);
    if (params.order) queryParams.append("order", params.order);

    const url = queryParams.toString() ? `${API.LIST}?${queryParams.toString()}` : API.LIST;
    return await apiService.get<ApiResponse<{ sales: Sale[], pagination: PaginationMeta }>>(url);
  },

  async getSaleById(id: string): Promise<ApiResponse<Sale>> {
    return await apiService.get<ApiResponse<Sale>>(API.GET_BY_ID(id));
  },

  async createSale(data: CreateSaleRequest): Promise<ApiResponse<Sale>> {
    return await apiService.post<ApiResponse<Sale>>(API.CREATE, data);
  },

  async updateSale(id: string, data: UpdateSaleRequest): Promise<ApiResponse<Sale>> {
    return await apiService.put<ApiResponse<Sale>>(API.UPDATE(id), data);
  },

  async deleteSale(id: string): Promise<ApiResponse<null>> {
    return await apiService.delete<ApiResponse<null>>(API.DELETE(id));
  },

  async activateSale(id: string): Promise<ApiResponse<Sale>> {
    return await apiService.patch<ApiResponse<Sale>>(API.ACTIVATE(id));
  },

  async deactivateSale(id: string): Promise<ApiResponse<Sale>> {
    return await apiService.patch<ApiResponse<Sale>>(API.DEACTIVATE(id));
  },

  async getSalesActive(): Promise<ApiResponse<SalesActiveResponse>> {
    return await apiService.get<ApiResponse<SalesActiveResponse>>(API.GET_SALE_ACTIVE);
  }
};

export default salesService;