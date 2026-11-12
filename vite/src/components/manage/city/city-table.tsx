import { memo, useCallback, useMemo } from 'react';
import { DataTable } from "@/components/globals/Data-table";
import { type City } from "@/types";
import { getCityColumns } from './city-columns';
//import { type RowSelectionState } from "@tanstack/react-table";

interface CityTableProps {
  data: City[];
  onEdit: (city: City) => void;
  onDelete: (city: City) => void;
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
  title: "ชื่อเมือง",
  image: "รูปภาพ",
  detail: "รายละเอียด",
  status: "สถานะ",
  created_at: "วันที่สร้าง",
  updated_at: "วันที่อัปเดต",
};

export const CityTable = memo(function CityTable({
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
 

}: CityTableProps) {
  const handlers = useMemo(() => ({
    onEdit,
    onDelete
  }), [onEdit, onDelete]);

  const columns = useMemo(() => getCityColumns(handlers), [handlers]);

 

  

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

  const getRowId = useCallback((row: City) => {
    return row.id.toString();
  }, []);

  console.log("CityTable rendering", {
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
      searchPlaceholder="ค้นหาชื่อเมือง..."
      columnToggleText="คอลัมน์"
      showAddButton={true}
      addButtonText="เพิ่มเมือง"
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
     // enableRowSelection={true}
     // selectedRowIds={selectedRowIds}
     // onRowSelectionChange={handleRowSelectionChange}
      pagination={pagination}

  
      //showBulkActions={true}
      //bulkActionText="ลบ"
      //onBulkAction={handleBulkDelete}
    />
  );
});