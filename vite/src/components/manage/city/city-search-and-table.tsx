import { memo, useCallback } from 'react';
import { useCitiesData } from '@/hooks/useCitiesData';
import { ControlledSearchInput } from '@/components/ui/controlled-search-input';
import { CityTable } from './city-table';
import type { City } from '@/types';

interface CitySearchAndTableProps {
  title?: string;
  description?: string;
  onEdit: (city: City) => void;
  onDelete: (city: City) => void;
  onAdd: () => void;
  onBulkDelete?: (selectedIds: string[]) => void;
  selectedRows?: string[];
  onRowSelectionChange?: (selectedRows: string[]) => void;
}

export const CitySearchAndTable = memo(function CitySearchAndTable({
  title = "จัดการเมือง",
  description = "จัดการข้อมูลเมืองต่างๆ สำหรับระบบ",
  onEdit,
  onDelete,
  onAdd,
  onBulkDelete,
  selectedRows = [],
  onRowSelectionChange
}: CitySearchAndTableProps) {

  // ใช้ useCitiesData hook ที่จัดการทุกอย่าง
  const {
    cities,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearchQuery,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
  } = useCitiesData();

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
          placeholder="ค้นหาชื่อเมือง..."
          defaultValue={currentSearchQuery}
        />
      </div>
      </div>

      {/* Loading State */}
      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        /* City Table */
        <CityTable
          data={cities}
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

export default CitySearchAndTable;