//src/services/admin/companyProfileService.ts

import { API } from '@/constants/api';
import apiService from '@/services/apiService';
import type { CompanyProfile, ApiResponse } from '@/types';
import type {
  GetCompanyProfilesRequest,
  SearchCompanyProfilesRequest,
  GetCompanyProfileByIdRequest,
  GetCompanyProfileByTaxIdRequest,
  ApproveCompanyProfileRequest,
  GetCompanyProfilesResponse,
  SearchCompanyProfilesResponse,
  GetCompanyProfileResponse,
  ApproveCompanyProfileResponse,
} from '@/types';

/**
 * Service สำหรับจัดการโปรไฟล์บริษัท (Admin)
 */
const companyProfileService = {
  /**
   * ดึงรายการโปรไฟล์บริษัทที่รอการอนุมัติ (สำหรับ admin+)
   * GET /api/v1/admin/companies
   * @param params - ตัวกรองและการแบ่งหน้า
   * @returns รายการโปรไฟล์บริษัทที่รอการอนุมัติ
   */
  getCompanyProfiles: async (params: GetCompanyProfilesRequest = {}): Promise<GetCompanyProfilesResponse> => {
    const queryParams = new URLSearchParams();

    if (params.page !== undefined) queryParams.append('page', params.page.toString());
    if (params.limit !== undefined) queryParams.append('limit', params.limit.toString());
    if (params.search) queryParams.append('search', params.search);

    const url = queryParams.toString() ? `${API.COMPANY_PROFILE.LIST}?${queryParams.toString()}` : API.COMPANY_PROFILE.LIST;
    return await apiService.get<GetCompanyProfilesResponse>(url);
  },

  /**
   * ค้นหาโปรไฟล์บริษัท (สำหรับ admin+)
   * GET /api/v1/admin/companies/search
   * @param params - พารามิเตอร์การค้นหา
   * @returns รายการโปรไฟล์บริษัทที่ตรงกับการค้นหา
   */
  searchCompanyProfiles: async (params: SearchCompanyProfilesRequest): Promise<SearchCompanyProfilesResponse> => {
    const queryParams = new URLSearchParams();

    queryParams.append('q', params.q);
    if (params.page !== undefined) queryParams.append('page', params.page.toString());
    if (params.limit !== undefined) queryParams.append('limit', params.limit.toString());

    const url = `${API.COMPANY_PROFILE.SEARCH}?${queryParams.toString()}`;
    return await apiService.get<SearchCompanyProfilesResponse>(url);
  },

  /**
   * ดึงโปรไฟล์บริษัทตาม ID (สำหรับ admin+)
   * GET /api/v1/admin/companies/:id
   * @param params - พารามิเตอร์ที่มี company profile ID
   * @returns ข้อมูลโปรไฟล์บริษัท
   */
  getCompanyProfileById: async (params: GetCompanyProfileByIdRequest): Promise<GetCompanyProfileResponse> => {
    return await apiService.get<GetCompanyProfileResponse>(API.COMPANY_PROFILE.GET_BY_ID(params.id));
  },

  /**
   * ดึงโปรไฟล์บริษัทตามเลขประจำตัวผู้เสียภาษี (สำหรับ admin+)
   * GET /api/v1/admin/companies/tax-id/:tax_id
   * @param params - พารามิเตอร์ที่มีเลขประจำตัวผู้เสียภาษี
   * @returns ข้อมูลโปรไฟล์บริษัท
   */
  getCompanyProfileByTaxId: async (params: GetCompanyProfileByTaxIdRequest): Promise<GetCompanyProfileResponse> => {
    return await apiService.get<GetCompanyProfileResponse>(API.COMPANY_PROFILE.GET_BY_TAX_ID(params.tax_id));
  },

  /**
   * อนุมัติโปรไฟล์บริษัท (สำหรับ admin+)
   * PATCH /api/v1/admin/companies/:id/approve
   * @param params - พารามิเตอร์ที่มี company profile ID
   * @returns ผลลัพธ์การอนุมัติ
   */
  updateCompanyProfile: async (id: string, data: Record<string, unknown>): Promise<ApiResponse<CompanyProfile>> => {
    return await apiService.put<ApiResponse<CompanyProfile>>(API.COMPANY_PROFILE.UPDATE(id), data);
  },

  approveCompanyProfile: async (params: ApproveCompanyProfileRequest): Promise<ApproveCompanyProfileResponse> => {
    return await apiService.patch<ApproveCompanyProfileResponse>(API.COMPANY_PROFILE.APPROVE(params.id), {});
  },
};

export default companyProfileService;