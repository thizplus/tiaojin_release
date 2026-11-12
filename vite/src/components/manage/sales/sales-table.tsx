import { memo, useCallback, useMemo } from "react";
import { DataTable } from "@/components/globals/Data-table";
import { type Sale } from "@/types";
import { getSalesColumns } from "./sales-columns";

interface SalesTableProps {
  data: Sale[];
  onEdit: (sale: Sale) => void;
  onDelete: (sale: Sale) => void;
  onActivate: (sale: Sale) => void;
  onDeactivate: (sale: Sale) => void;
  onAdd: () => void;
  page: number;
  pageSize: number;
  totalPages: number;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
}

const columnLabels = {
  code: "รหัส",
  first_name: "ชื่อ-นามสกุล",
  nick_name: "ชื่อเล่น",
  user_email: "อีเมล",
  active: "สถานะ",
  created_at: "วันที่สร้าง",
};

export const SalesTable = memo(function SalesTable({
  data,
  onEdit,
  onDelete,
  onActivate,
  onDeactivate,
  onAdd,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
}: SalesTableProps) {
  const handlers = useMemo(() => ({
    onEdit,
    onDelete,
    onActivate,
    onDeactivate
  }), [onEdit, onDelete, onActivate, onDeactivate]);

  const columns = useMemo(() => getSalesColumns(handlers), [handlers]);

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

  const getRowId = useCallback((row: Sale) => {
    return row.id.toString();
  }, []);

  return (
    <DataTable
      columns={columns}
      data={data}
      searchPlaceholder="ค้นหารหัส ชื่อ หรือชื่อเล่น..."
      columnToggleText="คอลัมน์"
      showAddButton={true}
      addButtonText="เพิ่มพนักงานขาย"
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
      pagination={pagination}
    />
  );
});