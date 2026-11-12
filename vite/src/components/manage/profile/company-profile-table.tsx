import { memo, useCallback, useMemo } from 'react';
import { DataTable } from "@/components/globals/Data-table";
import { type CompanyProfile } from "@/types";
import { getCompanyProfileColumns } from './company-profile-columns';

interface CompanyProfileTableProps {
  data: CompanyProfile[];
  onView: (profile: CompanyProfile) => void;
  onEdit?: (profile: CompanyProfile) => void;
  onApprove: (profile: CompanyProfile) => void;
  onAdd?: () => void;
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
  name: "ชื่อบริษัท",
  logo: "โลโก้",
  tax_id: "เลขประจำตัวผู้เสียภาษี",
  phone: "เบอร์โทรศัพท์",
  province: "จังหวัด",
  atta: "ATTA",
  dbd: "DBD",
  pp20: "PP20",
  created_at: "วันที่ส่งคำขอ",
  updated_at: "วันที่อัปเดต",
};

export const CompanyProfileTable = memo(function CompanyProfileTable({
  data,
  onView,
  onEdit,
  onApprove,
  onAdd,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
  selectedRows = [],
  //onRowSelectionChange,
}: CompanyProfileTableProps) {
  const handlers = useMemo(() => ({
    onView,
    onEdit,
    onApprove
  }), [onView, onEdit, onApprove]);

  const columns = useMemo(() => getCompanyProfileColumns(handlers), [handlers]);

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

  const getRowId = useCallback((row: CompanyProfile) => {
    return row.id.toString();
  }, []);

  console.log("CompanyProfileTable rendering", {
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
      searchPlaceholder="ค้นหาชื่อบริษัท..."
      columnToggleText="คอลัมน์"
      showAddButton={false}
      addButtonText=""
      onAdd={onAdd}
      columnLabels={columnLabels}
      getRowId={getRowId}
      pagination={pagination}
    />
  );
});