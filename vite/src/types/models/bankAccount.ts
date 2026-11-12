export interface CompanyBankAccount {
  id: string;
  bankCode: string;
  bankName: string;
  accountNumber: string;
  accountName: string;
  branchName?: string;
  isActive: boolean;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateBankAccountRequest {
  bankCode: string;
  bankName: string;
  accountNumber: string;
  accountName: string;
  branchName?: string;
  isActive?: boolean;
  sortOrder?: number;
}

export interface UpdateBankAccountRequest {
  bankCode?: string;
  bankName?: string;
  accountNumber?: string;
  accountName?: string;
  branchName?: string;
  isActive?: boolean;
  sortOrder?: number;
}
