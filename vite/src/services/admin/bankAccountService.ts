import apiService from '@/services/apiService';
import { BANK_ACCOUNT_API } from '@/constants/api';
import type { CompanyBankAccount, CreateBankAccountRequest, UpdateBankAccountRequest } from '@/types/models/bankAccount';

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

class BankAccountService {
  async list(): Promise<ApiResponse<CompanyBankAccount[]>> {
    return apiService.get<ApiResponse<CompanyBankAccount[]>>(BANK_ACCOUNT_API.LIST);
  }

  async listActive(): Promise<ApiResponse<CompanyBankAccount[]>> {
    return apiService.get<ApiResponse<CompanyBankAccount[]>>(BANK_ACCOUNT_API.ACTIVE);
  }

  async getById(id: string): Promise<ApiResponse<CompanyBankAccount>> {
    return apiService.get<ApiResponse<CompanyBankAccount>>(BANK_ACCOUNT_API.GET_BY_ID(id));
  }

  async create(data: CreateBankAccountRequest): Promise<ApiResponse<CompanyBankAccount>> {
    return apiService.post<ApiResponse<CompanyBankAccount>>(BANK_ACCOUNT_API.CREATE, data);
  }

  async update(id: string, data: UpdateBankAccountRequest): Promise<ApiResponse<CompanyBankAccount>> {
    return apiService.put<ApiResponse<CompanyBankAccount>>(BANK_ACCOUNT_API.UPDATE(id), data);
  }

  async delete(id: string): Promise<ApiResponse<null>> {
    return apiService.delete<ApiResponse<null>>(BANK_ACCOUNT_API.DELETE(id));
  }
}

export const bankAccountService = new BankAccountService();
