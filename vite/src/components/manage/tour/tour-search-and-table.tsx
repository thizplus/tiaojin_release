import { memo, useCallback } from "react";
import { useToursData } from "@/hooks/useToursData";
import { ControlledSearchInput } from "@/components/ui/controlled-search-input";
import { TourTable } from "./tour-table";
import type { Tour } from "@/types";

interface TourSearchAndTableProps {
  title?: string;
  description?: string;
  onEdit: (tour: Tour) => void;
  onDelete: (tour: Tour) => void;
  onUpdateStatus: (tour: Tour) => void;
  onManagePeriods: (tour: Tour) => void;
  onAdd: () => void;
}

export const TourSearchAndTable = memo(function TourSearchAndTable({
  title = "จัดการทัวร์",
  description = "จัดการข้อมูลทัวร์และสถานะการเผยแพร่",
  onEdit,
  onDelete,
  onUpdateStatus,
  onManagePeriods,
  onAdd,
}: TourSearchAndTableProps) {

  const {
    tours,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearch,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
  } = useToursData();

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
          placeholder="ค้นหาชื่อทัวร์หรือรหัสทัวร์..."
          defaultValue={currentSearch}
        />
      </div>
      </div>

      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        <TourTable
          data={tours}
          onEdit={onEdit}
          onDelete={onDelete}
          onUpdateStatus={onUpdateStatus}
          onManagePeriods={onManagePeriods}
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

export default TourSearchAndTable;