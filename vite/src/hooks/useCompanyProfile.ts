import { useCallback } from 'react';
import * as React from 'react';
import {
  useCompanyProfileStore,
  useCompanyProfilesListQuery,
  useSearchCompanyProfilesQuery,
  useCompanyProfileByIdQuery,
  useCompanyProfileByTaxIdQuery,
  useApproveCompanyProfileMutation
} from '@/stores/companyProfileStore';
import type {
  CompanyProfile,
  GetCompanyProfilesRequest,
  SearchCompanyProfilesRequest,
  ApproveCompanyProfileRequest
} from '@/types';

interface UseCompanyProfileReturn {
  // State
  currentCompanyProfile: CompanyProfile | null;
  selectedCompanyProfile: CompanyProfile | null;
  companyProfilesList: CompanyProfile[];
  isLoading: boolean;

  // Pagination data
  pagination: {
    total: number;
    offset: number;
    limit: number;
    has_next: boolean;
  } | null;

  // Actions
  approveCompanyProfile: (data: ApproveCompanyProfileRequest) => Promise<void>;

  // State management
  setCurrentCompanyProfile: (profile: CompanyProfile | null) => void;
  setSelectedCompanyProfile: (profile: CompanyProfile | null) => void;
  setCompanyProfilesList: (profiles: CompanyProfile[]) => void;
  clearCompanyProfileData: () => void;

  // Queries
  getCompanyProfilesList: () => void;
  searchCompanyProfiles: () => void;
  getCompanyProfileById: () => void;
  getCompanyProfileByTaxId: () => void;
  refetchCompanyProfilesList: () => void;
  refetchSearchResults: () => void;
  refetchCompanyProfileById: () => void;
  refetchCompanyProfileByTaxId: () => void;

  // Loading states
  isApproveLoading: boolean;
  isCompanyProfilesListLoading: boolean;
  isSearchLoading: boolean;
  isCompanyProfileByIdLoading: boolean;
  isCompanyProfileByTaxIdLoading: boolean;
}

export const useCompanyProfile = (
  listParams?: GetCompanyProfilesRequest,
  searchParams?: SearchCompanyProfilesRequest,
  companyProfileId?: string,
  taxId?: string,
  options: {
    enableList?: boolean;
    enableSearch?: boolean;
    enableById?: boolean;
    enableByTaxId?: boolean;
  } = {}
): UseCompanyProfileReturn => {
  // Store state
  const {
    currentCompanyProfile,
    selectedCompanyProfile,
    companyProfilesList,
    setCurrentCompanyProfile,
    setSelectedCompanyProfile,
    setCompanyProfilesList,
    clearCompanyProfileData
  } = useCompanyProfileStore();

  // Mutations
  const approveCompanyProfileMutation = useApproveCompanyProfileMutation();

  // Queries
  const {
    data: companyProfilesData,
    refetch: refetchCompanyProfilesList,
    isLoading: isCompanyProfilesListLoading
  } = useCompanyProfilesListQuery(listParams || {}, {
    enabled: options.enableList === true
  });

  const {
    data: searchData,
    refetch: refetchSearchResults,
    isLoading: isSearchLoading
  } = useSearchCompanyProfilesQuery(searchParams || { q: '' }, {
    enabled: options.enableSearch === true && !!searchParams?.q
  });

  const {
    refetch: refetchCompanyProfileById,
    isLoading: isCompanyProfileByIdLoading
  } = useCompanyProfileByIdQuery(companyProfileId || '', {
    enabled: options.enableById !== false && !!companyProfileId
  });

  const {
    refetch: refetchCompanyProfileByTaxId,
    isLoading: isCompanyProfileByTaxIdLoading
  } = useCompanyProfileByTaxIdQuery(taxId || '', {
    enabled: options.enableByTaxId !== false && !!taxId
  });

  // Actions
  const approveCompanyProfile = useCallback(async (data: ApproveCompanyProfileRequest) => {
    await approveCompanyProfileMutation.mutateAsync(data);
  }, [approveCompanyProfileMutation]);

  // Query management
  const getCompanyProfilesList = useCallback(() => {
    refetchCompanyProfilesList();
  }, [refetchCompanyProfilesList]);

  const searchCompanyProfiles = useCallback(() => {
    refetchSearchResults();
  }, [refetchSearchResults]);

  const getCompanyProfileById = useCallback(() => {
    refetchCompanyProfileById();
  }, [refetchCompanyProfileById]);

  const getCompanyProfileByTaxId = useCallback(() => {
    refetchCompanyProfileByTaxId();
  }, [refetchCompanyProfileByTaxId]);

  // Process company profiles data and pagination
  const pagination = React.useMemo(() => {
    const data = companyProfilesData || searchData;
    if (!data?.data?.pagination) return null;

    const paginationMeta = data.data.pagination;
    return {
      total: paginationMeta.total,
      offset: paginationMeta.offset,
      limit: paginationMeta.limit,
      has_next: paginationMeta.has_next,
    };
  }, [companyProfilesData, searchData]);

  // Calculate overall loading state
  const isLoading = approveCompanyProfileMutation.isPending ||
                   isCompanyProfilesListLoading ||
                   isSearchLoading ||
                   isCompanyProfileByIdLoading ||
                   isCompanyProfileByTaxIdLoading;

  return {
    // State
    currentCompanyProfile,
    selectedCompanyProfile,
    companyProfilesList,
    isLoading,

    // Pagination data
    pagination,

    // Actions
    approveCompanyProfile,

    // State management
    setCurrentCompanyProfile,
    setSelectedCompanyProfile,
    setCompanyProfilesList,
    clearCompanyProfileData,

    // Queries
    getCompanyProfilesList,
    searchCompanyProfiles,
    getCompanyProfileById,
    getCompanyProfileByTaxId,
    refetchCompanyProfilesList,
    refetchSearchResults,
    refetchCompanyProfileById,
    refetchCompanyProfileByTaxId,

    // Loading states
    isApproveLoading: approveCompanyProfileMutation.isPending,
    isCompanyProfilesListLoading,
    isSearchLoading,
    isCompanyProfileByIdLoading,
    isCompanyProfileByTaxIdLoading,
  };
};

export default useCompanyProfile;