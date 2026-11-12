import { memo, useCallback, useMemo, useState } from "react";
import { DataTable } from "@/components/globals/Data-table";
import { type TourPeriod } from "@/types";
import { getTourPeriodsColumns } from "./tour-periods-columns";
import { PeriodDetailDialog } from "./PeriodDetailDialog";
import { PeriodFilesDialog } from "./PeriodFilesDialog";
import { SeatHoldDialog } from "@/components/manage/seat-hold";

interface TourPeriodsTableProps {
  data: TourPeriod[];
  page: number;
  pageSize: number;
  totalPages: number;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  onEdit?: (period: TourPeriod) => void;
  onDelete?: (period: TourPeriod) => void;
}

const columnLabels = {
  start_date: "ช่วงเวลา",
  seat: "ที่นั่ง",
  airline: "สายการบิน",
  rate_items: "ราคา",
  deposit: "มัดจำ",
  commission_items: "คอมมิชชั่น",
  status: "สถานะ",
};

export const TourPeriodsTable = memo(function TourPeriodsTable({
  data,
  page,
  pageSize,
  totalPages,
  setPage,
  setPageSize,
  onEdit,
  onDelete,
}: TourPeriodsTableProps) {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [seatHoldDialogOpen, setSeatHoldDialogOpen] = useState(false);
  const [filesDialogOpen, setFilesDialogOpen] = useState(false);
  const [selectedPeriod, setSelectedPeriod] = useState<TourPeriod | null>(null);

  const handleEdit = useCallback((period: TourPeriod) => {
    if (onEdit) {
      onEdit(period);
    } else {
      setSelectedPeriod(period);
      setDialogOpen(true);
    }
  }, [onEdit]);

  const handleView = useCallback((period: TourPeriod) => {
    setSelectedPeriod(period);
    setDialogOpen(true);
  }, []);

  const handleManageSeatHold = useCallback((period: TourPeriod) => {
    setSelectedPeriod(period);
    setSeatHoldDialogOpen(true);
  }, []);

  const handleManageFiles = useCallback((period: TourPeriod) => {
    setSelectedPeriod(period);
    setFilesDialogOpen(true);
  }, []);

  const handlers = useMemo(() => ({
    onEdit: handleEdit,
    onView: onEdit ? handleView : undefined,
    onDelete,
    onManageSeatHold: handleManageSeatHold,
    onManageFiles: handleManageFiles,
  }), [handleEdit, handleView, onEdit, onDelete, handleManageSeatHold, handleManageFiles]);

  const columns = useMemo(() => getTourPeriodsColumns(handlers), [handlers]);

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

  const getRowId = useCallback((row: TourPeriod) => {
    return row.id.toString();
  }, []);


  return (
    <>
      <DataTable
        columns={columns}
        data={data}
        searchPlaceholder="ค้นหาช่วงเวลาทัวร์..."
        columnToggleText="คอลัมน์"
        showAddButton={false}
        columnLabels={columnLabels}
        getRowId={getRowId}
        pagination={pagination}
      />

      <PeriodDetailDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        period={selectedPeriod}
      />

      <SeatHoldDialog
        open={seatHoldDialogOpen}
        onOpenChange={setSeatHoldDialogOpen}
        period={selectedPeriod}
      />

      <PeriodFilesDialog
        open={filesDialogOpen}
        onOpenChange={setFilesDialogOpen}
        period={selectedPeriod}
      />
    </>
  );
});