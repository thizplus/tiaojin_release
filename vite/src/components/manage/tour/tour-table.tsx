import { memo, useCallback, useMemo } from "react";
import { DataTable } from "@/components/globals/Data-table";
import { type Tour } from "@/types";
import { getTourColumns } from "./tour-columns";

interface TourTableProps {
  data: Tour[];
  onEdit: (tour: Tour) => void;
  onDelete: (tour: Tour) => void;
  onUpdateStatus: (tour: Tour) => void;
  onManagePeriods: (tour: Tour) => void;
  onAdd: () => void;
  page: number;
  pageSize: number;
  totalPages: number;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
}

const columnLabels = {
  code: "รหัสทัวร์",
  title: "ชื่อทัวร์",
  city_title: "จังหวัด",
  thumbnail: "รูปภาพ",
  status: "สถานะ",
  files: "ไฟล์แนบ",
  periods: "จัดการ Periods",
  created_at: "วันที่สร้าง",
};

export const TourTable = memo(function TourTable({
  data,
  onEdit,
  onDelete,
  onUpdateStatus,
  onManagePeriods,
  onAdd,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
}: TourTableProps) {
  const handlers = useMemo(() => ({
    onEdit,
    onDelete,
    onUpdateStatus,
    onManagePeriods
  }), [onEdit, onDelete, onUpdateStatus, onManagePeriods]);

  const columns = useMemo(() => getTourColumns(handlers), [handlers]);

  const handlePageChange = useCallback((newPage: number) => {
    setPage(newPage);
  }, [setPage]);

  const handlePageSizeChange = useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    setPage(0);
  }, [setPageSize, setPage]);

  const pagination = useMemo(() => ({
    pageIndex: page,
    pageSize: pageSize,
    pageCount: totalPages,
    canPreviousPage: page > 0,
    canNextPage: page < totalPages - 1,
    setPageIndex: handlePageChange,
    setPageSize: handlePageSizeChange,
    previousPage: () => handlePageChange(page - 1),
    nextPage: () => handlePageChange(page + 1),
    getPageCount: () => totalPages,
  }), [page, pageSize, totalPages, handlePageChange, handlePageSizeChange]);

  const getRowId = useCallback((row: Tour) => {
    return row.id.toString();
  }, []);

  return (
    <DataTable
      columns={columns}
      data={data}
      searchPlaceholder="ค้นหาชื่อทัวร์หรือรหัสทัวร์..."
      columnToggleText="คอลัมน์"
      showAddButton={true}
      addButtonText="เพิ่มทัวร์"
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
      pagination={pagination}
    />
  );
});