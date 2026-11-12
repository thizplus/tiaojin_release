import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import companyProfileService from '@/services/admin/companyProfileService';
import type {
  CompanyProfile,
  GetCompanyProfilesRequest,
  SearchCompanyProfilesRequest,

  ApproveCompanyProfileRequest,
  GetCompanyProfilesResponse,
  SearchCompanyProfilesResponse,
  GetCompanyProfileResponse,
  ApproveCompanyProfileResponse,
} from '@/types';

const companyProfileKeys = {
  companyProfiles: (params?: GetCompanyProfilesRequest) => ['companyProfile', 'list', params] as const,
  searchCompanyProfiles: (params?: SearchCompanyProfilesRequest) => ['companyProfile', 'search', params] as const,
  companyProfileById: (id: string) => ['companyProfile', 'byId', id] as const,
  companyProfileByTaxId: (taxId: string) => ['companyProfile', 'byTaxId', taxId] as const,
};

interface CompanyProfileState {
  currentCompanyProfile: CompanyProfile | null;
  selectedCompanyProfile: CompanyProfile | null;
  companyProfilesList: CompanyProfile[];
  setCurrentCompanyProfile: (profile: CompanyProfile | null) => void;
  setSelectedCompanyProfile: (profile: CompanyProfile | null) => void;
  setCompanyProfilesList: (profiles: CompanyProfile[]) => void;
  clearCompanyProfileData: () => void;
}

export const useCompanyProfileStore = create<CompanyProfileState>()(
  persist(
    (set) => ({
      currentCompanyProfile: null,
      selectedCompanyProfile: null,
      companyProfilesList: [],
      setCurrentCompanyProfile: (profile) => set({ currentCompanyProfile: profile }),
      setSelectedCompanyProfile: (profile) => set({ selectedCompanyProfile: profile }),
      setCompanyProfilesList: (profiles) => set({ companyProfilesList: profiles }),
      clearCompanyProfileData: () => set({
        currentCompanyProfile: null,
        selectedCompanyProfile: null,
        companyProfilesList: [],
      }),
    }),
    {
      name: 'company-profile-storage',
      partialize: (state) => ({
        selectedCompanyProfile: state.selectedCompanyProfile,
      }),
    }
  )
);

// Admin+ Queries

/**
 * ดึงรายการโปรไฟล์บริษัทที่รอการอนุมัติ (สำหรับ admin+)
 */
export const useCompanyProfilesListQuery = (
  params: GetCompanyProfilesRequest = {},
  options: { enabled?: boolean } = {}
): UseQueryResult<GetCompanyProfilesResponse, Error> => {
  const { setCompanyProfilesList } = useCompanyProfileStore();

  return useQuery({
    queryKey: companyProfileKeys.companyProfiles(params),
    queryFn: async () => {
      const response = await companyProfileService.getCompanyProfiles(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch company profiles list');
      }

      const profiles = response.data.profiles || [];
      setCompanyProfilesList(Array.isArray(profiles) ? profiles : []);
      return response;
    },
    enabled: options.enabled !== undefined ? options.enabled : true,
  });
};

/**
 * ค้นหาโปรไฟล์บริษัท (สำหรับ admin+)
 */
export const useSearchCompanyProfilesQuery = (
  params: SearchCompanyProfilesRequest,
  options: { enabled?: boolean } = {}
): UseQueryResult<SearchCompanyProfilesResponse, Error> => {
  const { setCompanyProfilesList } = useCompanyProfileStore();

  return useQuery({
    queryKey: companyProfileKeys.searchCompanyProfiles(params),
    queryFn: async () => {
      const response = await companyProfileService.searchCompanyProfiles(params);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to search company profiles');
      }

      const profiles = response.data.profiles || [];
      setCompanyProfilesList(Array.isArray(profiles) ? profiles : []);
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!params.q,
  });
};

/**
 * ดึงโปรไฟล์บริษัทตาม ID (สำหรับ admin+)
 */
export const useCompanyProfileByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<GetCompanyProfileResponse, Error> => {
  const { setSelectedCompanyProfile } = useCompanyProfileStore();

  return useQuery({
    queryKey: companyProfileKeys.companyProfileById(id),
    queryFn: async () => {
      const response = await companyProfileService.getCompanyProfileById({ id });

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch company profile data');
      }

      setSelectedCompanyProfile(response.data);
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!id,
  });
};

/**
 * ดึงโปรไฟล์บริษัทตามเลขประจำตัวผู้เสียภาษี (สำหรับ admin+)
 */
export const useCompanyProfileByTaxIdQuery = (
  taxId: string,
  options: { enabled?: boolean } = {}
): UseQueryResult<GetCompanyProfileResponse, Error> => {
  const { setSelectedCompanyProfile } = useCompanyProfileStore();

  return useQuery({
    queryKey: companyProfileKeys.companyProfileByTaxId(taxId),
    queryFn: async () => {
      const response = await companyProfileService.getCompanyProfileByTaxId({ tax_id: taxId });

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Unable to fetch company profile data');
      }

      setSelectedCompanyProfile(response.data);
      return response;
    },
    enabled: (options.enabled !== undefined ? options.enabled : true) && !!taxId,
  });
};

// Admin+ Mutations

/**
 * อนุมัติโปรไฟล์บริษัท (สำหรับ admin+)
 */
export const useApproveCompanyProfileMutation = (): UseMutationResult<
  ApproveCompanyProfileResponse,
  Error,
  ApproveCompanyProfileRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: ApproveCompanyProfileRequest) => {
      const response = await companyProfileService.approveCompanyProfile(data);

      if (!response.success) {
        throw new Error(response.message || 'Failed to approve company profile');
      }

      return response;
    },
    onSuccess: (_, { id }) => {
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: companyProfileKeys.companyProfileById(id) });
      queryClient.invalidateQueries({ queryKey: companyProfileKeys.companyProfiles() });
      queryClient.invalidateQueries({ queryKey: companyProfileKeys.searchCompanyProfiles() });
    },
  });
};