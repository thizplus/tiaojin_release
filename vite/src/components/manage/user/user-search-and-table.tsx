import { memo, useCallback } from "react";
import { useUsersData } from "@/hooks/useUsersData";
import { ControlledSearchInput } from "@/components/ui/controlled-search-input";
import { UserTable } from "./user-table";
import type { User } from "@/types";

interface UserSearchAndTableProps {
  title?: string;
  description?: string;
  onEdit: (user: User) => void;
  onDelete: (user: User) => void;
  onChangeStatus: (user: User) => void;
  onChangeRole: (user: User) => void;
  onResetPassword: (user: User) => void;
  onAdd: () => void;
}

export const UserSearchAndTable = memo(function UserSearchAndTable({
  title = "จัดการผู้ใช้",
  description = "จัดการข้อมูลผู้ใช้และสิทธิ์การเข้าถึงระบบ",
  onEdit,
  onDelete,
  onChangeStatus,
  onChangeRole,
  onResetPassword,
  onAdd,
}: UserSearchAndTableProps) {

  const {
    users,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    currentSearchQuery,
    handleSearchChange,
    handlePageChange,
    handlePageSizeChange,
  } = useUsersData();

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
          placeholder="ค้นหาอีเมลหรือชื่อบริษัท..."
          defaultValue={currentSearchQuery}
        />
      </div>
      </div>

      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        <UserTable
          data={users}
          onEdit={onEdit}
          onDelete={onDelete}
          onChangeStatus={onChangeStatus}
          onChangeRole={onChangeRole}
          onResetPassword={onResetPassword}
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

export default UserSearchAndTable;
