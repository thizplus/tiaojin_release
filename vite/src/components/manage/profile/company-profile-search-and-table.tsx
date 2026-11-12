import { memo, useCallback } from 'react';
import { useCompanyProfilesData } from '@/hooks/useCompanyProfilesData';
import { ControlledSearchInput } from '@/components/ui/controlled-search-input';
import { CompanyProfileTable } from './company-profile-table';
import type { CompanyProfile } from '@/types';

interface CompanyProfileSearchAndTableProps {
  title?: string;
  description?: string;
  onView: (profile: CompanyProfile) => void;
  onEdit?: (profile: CompanyProfile) => void;
  onApprove: (profile: CompanyProfile) => void;
  selectedRows?: string[];
  onRowSelectionChange?: (selectedRows: string[]) => void;
}

export const CompanyProfileSearchAndTable = memo(function CompanyProfileSearchAndTable({
  title = "จัดการโปรไฟล์บริษัท",
  description = "จัดการข้อมูลโปรไฟล์บริษัท",
  onView,
  onEdit,
  onApprove,
  selectedRows = [],
  onRowSelectionChange
}: CompanyProfileSearchAndTableProps) {

  // ใช้ useCompanyProfilesData hook ที่จัดการทุกอย่าง
  const {
    companyProfiles,
   // pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearchQuery,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
    totalPages,
  } = useCompanyProfilesData();

  // จัดการ search
  const handleSearch = useCallback((query: string) => {
    handleSearchChange(query);
  }, [handleSearchChange]);

  return (
    <div className="space-y-6">
     <div className="flex flex-col items-left gap-4 md:flex-row md:items-center justify-between mb-6 ">
         {/* Header Section */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">{title}</h1>
          <p className="text-muted-foreground mt-1">{description}</p>
        </div>
      </div>

      {/* Search Input */}
      <div className="relative w-full max-w-sm">
        <ControlledSearchInput
          onSearchChange={handleSearch}
          placeholder="ค้นหาชื่อบริษัท..."
          defaultValue={currentSearchQuery}
        />
      </div>
      </div>

      {/* Loading State */}
      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        /* Company Profile Table */
        <CompanyProfileTable
          data={companyProfiles}
          onView={onView}
          onEdit={onEdit}
          onApprove={onApprove}
          page={currentPage - 1} // ปรับให้เริ่มจาก 0 สำหรับ DataTable
          pageSize={currentPageSize}
          totalPages={totalPages}
          setPage={handlePageChange}
          setPageSize={handlePageSizeChange}
          selectedRows={selectedRows}
          onRowSelectionChange={onRowSelectionChange}
        />
      )}
    </div>
  );
});

export default CompanyProfileSearchAndTable;