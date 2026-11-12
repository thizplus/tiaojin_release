import { memo, useCallback, useMemo } from "react";
import { DataTable } from "@/components/globals/Data-table";
import { type User } from "@/types";
import { getUserColumns } from "./user-columns";

interface UserTableProps {
  data: User[];
  onEdit: (user: User) => void;
  onDelete: (user: User) => void;
  onChangeStatus: (user: User) => void;
  onChangeRole: (user: User) => void;
  onResetPassword: (user: User) => void;
  onAdd: () => void;
  page: number;
  pageSize: number;
  totalPages: number;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
}

const columnLabels = {
  email: "อีเมล",
  company_profile: "บริษัท",
  role: "บทบาท",
  status: "สถานะ",
  created_at: "วันที่สร้าง",
  updated_at: "วันที่อัปเดต",
};

export const UserTable = memo(function UserTable({
  data,
  onEdit,
  onDelete,
  onChangeStatus,
  onChangeRole,
  onResetPassword,
  onAdd,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
}: UserTableProps) {
  const handlers = useMemo(() => ({
    onEdit,
    onDelete,
    onChangeStatus,
    onChangeRole,
    onResetPassword
  }), [onEdit, onDelete, onChangeStatus, onChangeRole, onResetPassword]);

  const columns = useMemo(() => getUserColumns(handlers), [handlers]);

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

  const getRowId = useCallback((row: User) => {
    return row.id.toString();
  }, []);

  return (
    <DataTable
      columns={columns}
      data={data}
      searchPlaceholder="ค้นหาอีเมลหรือชื่อบริษัท..."
      columnToggleText="คอลัมน์"
      showAddButton={true}
      addButtonText="เพิ่มผู้ใช้"
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
      pagination={pagination}
    />
  );
});
