import { memo, useCallback, useMemo } from 'react';
import { DataTable } from "@/components/globals/Data-table";
import type { ApiKey } from "@/types/models/apiKey";
import { getApiKeyColumns } from './api-key-columns';

interface ApiKeyTableProps {
  data: ApiKey[];
  onEdit: (apiKey: ApiKey) => void;
  onDelete: (apiKey: ApiKey) => void;
  onAdd: () => void;
  // Pagination props
  page: number;
  pageSize: number;
  totalPages: number;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  // Row selection props
  selectedRows?: string[];
  onRowSelectionChange?: (selectedRows: string[]) => void;
  // Search props
  searchQuery?: string;
  onSearchChange?: (query: string) => void;
}

const columnLabels = {
  name: "ชื่อ API Key",
  key_preview: "API Key",
  is_active: "สถานะ",
  permissions: "สิทธิ์",
  last_used_at: "ใช้งานล่าสุด",
  expires_at: "วันหมดอายุ",
  created_at: "วันที่สร้าง",
};

export const ApiKeyTable = memo(function ApiKeyTable({
  data,
  onEdit,
  onDelete,
  onAdd,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
 
}: ApiKeyTableProps) {
  const handlers = useMemo(() => ({
    onEdit,
    onDelete
  }), [onEdit, onDelete]);

  const columns = useMemo(() => getApiKeyColumns(handlers), [handlers]);

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



  const getRowId = useCallback((row: ApiKey) => {
    return row.id.toString();
  }, []);

  return (
    <DataTable
      columns={columns}
      data={data}
      searchPlaceholder="ค้นหา API Key..."
      columnToggleText="คอลัมน์"
      showAddButton={true}
      addButtonText="เพิ่ม API Key"
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
      pagination={pagination}
    />
  );
});