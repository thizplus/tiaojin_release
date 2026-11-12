import { memo, useCallback } from "react";
import { useSalesData } from "@/hooks/useSalesData";
import { ControlledSearchInput } from "@/components/ui/controlled-search-input";
import { SalesTable } from "./sales-table";
import type { Sale } from "@/types";

interface SalesSearchAndTableProps {
  title?: string;
  description?: string;
  onEdit: (sale: Sale) => void;
  onDelete: (sale: Sale) => void;
  onActivate: (sale: Sale) => void;
  onDeactivate: (sale: Sale) => void;
  onAdd: () => void;
}

export const SalesSearchAndTable = memo(function SalesSearchAndTable({
  title = "จัดการพนักงานขาย",
  description = "จัดการข้อมูลพนักงานขายและสถานะการใช้งาน",
  onEdit,
  onDelete,
  onActivate,
  onDeactivate,
  onAdd,
}: SalesSearchAndTableProps) {

  const {
    sales,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearch,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
  } = useSalesData();

  const handleSearch = useCallback((query: string) => {
    handleSearchChange(query);
  }, [handleSearchChange]);

  return (
    <div className="space-y-6">
     <div className="flex flex-col items-left gap-4 md:flex-row md:items-center justify-between mb-6 ">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">{title}</h1>
          <p className="text-muted-foreground mt-1">{description}</p>
        </div>
      </div>

      <div className="relative w-full max-w-sm">
        <ControlledSearchInput
          onSearchChange={handleSearch}
          placeholder="ค้นหารหัส ชื่อ หรือชื่อเล่น..."
          defaultValue={currentSearch}
        />
      </div>
      </div>

      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        <SalesTable
          data={sales}
          onEdit={onEdit}
          onDelete={onDelete}
          onActivate={onActivate}
          onDeactivate={onDeactivate}
          onAdd={onAdd}
          page={currentPage - 1}
          pageSize={currentPageSize}
          totalPages={pagination?.totalPages || 1}
          setPage={handlePageChange}
          setPageSize={handlePageSizeChange}
        />
      )}
    </div>
  );
});

export default SalesSearchAndTable;