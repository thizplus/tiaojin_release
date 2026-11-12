import { memo, useCallback } from 'react';
import { useBannersData } from '@/hooks/useBannersData';
import { ControlledSearchInput } from '@/components/ui/controlled-search-input';
import { BannerTable } from './banner-table';
import type { Banner, BannerPosition } from '@/types';

interface BannerSearchAndTableProps {
  title?: string;
  description?: string;
  position?: BannerPosition;
  onEdit: (banner: Banner) => void;
  onDelete: (banner: Banner) => void;
  onAdd: () => void;
  onBulkDelete?: (selectedIds: string[]) => void;
  selectedRows?: string[];
  onRowSelectionChange?: (selectedRows: string[]) => void;
}

export const BannerSearchAndTable = memo(function BannerSearchAndTable({
  title,
  description,
  position,
  onEdit,
  onDelete,
  onAdd,
  onBulkDelete,
  selectedRows = [],
  onRowSelectionChange
}: BannerSearchAndTableProps) {

  // ใช้ useBannersData hook ที่จัดการทุกอย่าง
  const {
    banners,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearchQuery,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
  } = useBannersData({ position });

  // จัดการ search
  const handleSearch = useCallback((query: string) => {
    handleSearchChange(query);
  }, [handleSearchChange]);

  return (
    <div className="space-y-6">
     <div className="flex flex-col items-left gap-4 md:flex-row md:items-center justify-between mb-6 ">
         {/* Header Section */}
      {(title || description) && (
        <div className="flex justify-between items-center">
          <div>
            {title && <h1 className="text-2xl font-bold">{title}</h1>}
            {description && <p className="text-muted-foreground mt-1">{description}</p>}
          </div>
        </div>
      )}

      {/* Search Input */}
      <div className="relative w-full max-w-sm">
        <ControlledSearchInput
          onSearchChange={handleSearch}
          placeholder="ค้นหาชื่อแบนเนอร์..."
          defaultValue={currentSearchQuery}
        />
      </div>
      </div>

      {/* Loading State */}
      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        /* Banner Table */
        <BannerTable
          data={banners}
          onEdit={onEdit}
          onDelete={onDelete}
          onBulkDelete={onBulkDelete}
          onAdd={onAdd}
          page={currentPage - 1} // ปรับให้เริ่มจาก 0 สำหรับ DataTable
          pageSize={currentPageSize}
          totalPages={pagination?.totalPages || 1}
          setPage={handlePageChange}
          setPageSize={handlePageSizeChange}
          selectedRows={selectedRows}
          onRowSelectionChange={onRowSelectionChange}
          searchQuery="" // ไม่ต้องส่ง เพราะ search แยกออกมาแล้ว
          onSearchChange={() => {}} // ไม่ต้องใช้
        />
      )}
    </div>
  );
});

export default BannerSearchAndTable;