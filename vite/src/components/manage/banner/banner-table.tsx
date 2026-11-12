import { memo, useCallback, useMemo } from 'react';
import { DataTable } from "@/components/globals/Data-table";
import { type Banner } from "@/types";
import { getBannerColumns } from './banner-columns';

interface BannerTableProps {
  data: Banner[];
  onEdit: (banner: Banner) => void;
  onDelete: (banner: Banner) => void;
  onAdd: () => void;
  onBulkDelete?: (selectedIds: string[]) => void;
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
  title: "ชื่อแบนเนอร์",
  image: "รูปภาพ",
  link_url: "ลิงก์",
  alt: "Alt Text",
  display_order: "ลำดับ",
  status: "สถานะ",
  created_at: "วันที่สร้าง",
  updated_at: "วันที่อัปเดต",
};

export const BannerTable = memo(function BannerTable({
  data,
  onEdit,
  onDelete,
  onAdd,
  //onBulkDelete,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
  selectedRows = [],
  //onRowSelectionChange,

}: BannerTableProps) {
  const handlers = useMemo(() => ({
    onEdit,
    onDelete
  }), [onEdit, onDelete]);

  const columns = useMemo(() => getBannerColumns(handlers), [handlers]);

  const handlePageChange = useCallback((newPage: number) => {
    setPage(newPage);
  }, [setPage]);

  const handlePageSizeChange = useCallback((newPageSize: number) => {
    console.log('Changing page size to:', newPageSize);
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

  const selectedRowIds = useMemo(() => {
    const rowIds: Record<string, boolean> = {};
    selectedRows.forEach(id => {
      rowIds[id] = true;
    });
    return rowIds;
  }, [selectedRows]);

  const getRowId = useCallback((row: Banner) => {
    return row.id.toString();
  }, []);

  console.log("BannerTable rendering", {
    page,
    pageSize,
    totalPages,
    selectedRows,
    selectedRowIds,
    dataLength: data.length,
    sampleRowId: data.length > 0 ? data[0].id : 'no data',
  });

  return (
    <DataTable
      columns={columns}
      data={data}
      searchPlaceholder="ค้นหาชื่อแบนเนอร์..."
      columnToggleText="คอลัมน์"
      showAddButton={true}
      addButtonText="เพิ่มแบนเนอร์"
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
      pagination={pagination}
    />
  );
});