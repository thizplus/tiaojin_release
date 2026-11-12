import React, { useState, useEffect } from 'react';
import { type DateRange } from 'react-day-picker';
import { format } from 'date-fns';
import { DateRangePicker } from '@/components/admin/DateRangePicker';
import { SalesSelector } from '@/components/admin/SalesSelector';
import { BookingStatusFilter, type BookingStatus } from '@/components/admin/BookingStatusFilter';
import { BookingGridCard } from '@/components/manage/booking/BookingGridCard';
import { BookingDetailSheet } from '@/components/manage/booking/BookingDetailSheet';
import { FileDownloadSheet } from '@/components/manage/booking/FileDownloadSheet';
import { AdvancedPagination } from '@/components/ui/advanced-pagination';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Loader2, CalendarRange, Download } from 'lucide-react';
import { useAdminBooking } from '@/hooks/useAdminBooking';
import { exportBookingsToCSV } from '@/utils/exportBookingsCSV';
import type { Sale } from '@/types/models/sales';
import type { BookingResponse } from '@/types';

interface BookingAction {
  key: string;
  label: string;
  icon?: React.ReactNode;
  onClick: (booking: BookingResponse) => void;
  variant?: 'default' | 'destructive';
  showFor?: (booking: BookingResponse) => boolean;
  separator?: boolean;
}

export default function ReportsDateRangePage() {
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined);
  const [selectedSales, setSelectedSales] = useState<Sale | null>(null);
  const [selectedStatuses, setSelectedStatuses] = useState<BookingStatus[]>([]);
  const [selectedBooking, setSelectedBooking] = useState<BookingResponse | null>(null);
  const [isDetailSheetOpen, setIsDetailSheetOpen] = useState(false);
  const [isFilesSheetOpen, setIsFilesSheetOpen] = useState(false);
  const [selectedBookingForFiles, setSelectedBookingForFiles] = useState<BookingResponse | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);

  // สร้าง filters object
  const shouldFetch = dateRange?.from && dateRange?.to;
  const filters = shouldFetch ? {
    start_date: format(dateRange.from!, 'yyyy-MM-dd'),
    end_date: format(dateRange.to!, 'yyyy-MM-dd'),
    sales_id: selectedSales?.id,
    status: selectedStatuses.length > 0 ? selectedStatuses.join(',') : undefined,
    page: currentPage,
    limit: pageSize,
  } : {};

  // ดึงข้อมูล bookings
  const {
    bookingsList: bookings,
    pagination,
    isLoading,
    setBookingFilters,
  } = useAdminBooking(
    filters,
    undefined,
    { enableList: !!shouldFetch, enableById: false }
  );

  // อัปเดต filters เมื่อมีการเปลี่ยนแปลง
  useEffect(() => {
    if (shouldFetch) {
      setBookingFilters(filters);
    }
  }, [dateRange, selectedSales, selectedStatuses, currentPage, pageSize, setBookingFilters]);

  // Reset pagination เมื่อเปลี่ยน filters
  useEffect(() => {
    setCurrentPage(1);
  }, [dateRange, selectedSales, selectedStatuses]);

  const handleView = (booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDetailSheetOpen(true);
  };

  const handleCardSelect = (booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDetailSheetOpen(true);
  };

  const handleViewFiles = (booking: BookingResponse) => {
    setSelectedBookingForFiles(booking);
    setIsFilesSheetOpen(true);
  };

  const handlePageChange = React.useCallback((page: number) => {
    setCurrentPage(page);
  }, []);

  const handlePageSizeChange = React.useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    setCurrentPage(1);
  }, []);

  // Admin actions
  const adminActions: BookingAction[] = [];

  const displayBookings = bookings || [];

  // Export to CSV handler
  const handleExportCSV = React.useCallback(() => {
    if (!dateRange?.from || !dateRange?.to) {
      alert('กรุณาเลือกช่วงวันที่ก่อน');
      return;
    }

    if (!bookings || bookings.length === 0) {
      alert('ไม่มีข้อมูลสำหรับการส่งออก');
      return;
    }

    const salesName = selectedSales
      ? ` - ${selectedSales.first_name} ${selectedSales.last_name}`
      : '';

    exportBookingsToCSV(bookings, {
      tourTitle: `รายงานตามช่วงวันที่${salesName}`,
      periodStartDate: format(dateRange.from, 'yyyy-MM-dd'),
      periodEndDate: format(dateRange.to, 'yyyy-MM-dd')
    });
  }, [bookings, dateRange, selectedSales]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <CalendarRange className="h-8 w-8 text-primary" />
        <div>
          <h1 className="text-2xl font-bold">รายงานตามช่วงวันที่</h1>
          <p className="text-muted-foreground">ดูรายการจองตามช่วงวันที่และ Sales</p>
        </div>
      </div>

      {/* Date Range และ Sales Selectors */}
      <div className="space-y-6">
        <DateRangePicker
          dateRange={dateRange}
          onDateRangeChange={setDateRange}
        />
        <SalesSelector
          selectedSales={selectedSales}
          onSalesSelect={setSelectedSales}
        />
      </div>

      {/* Status Filter */}
      {shouldFetch && (
        <BookingStatusFilter
          selectedStatuses={selectedStatuses}
          onStatusChange={setSelectedStatuses}
        />
      )}

      {/* Summary Information */}
      {shouldFetch && pagination && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">จำนวนการจองทั้งหมด</div>
            <div className="text-2xl font-bold">{pagination.total || 0}</div>
          </div>
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">ช่วงวันที่</div>
            <div className="text-lg font-medium">
              {dateRange?.from && dateRange?.to ? (
                <>
                  {format(dateRange.from, 'dd MMM yyyy')} - {format(dateRange.to, 'dd MMM yyyy')}
                </>
              ) : '-'}
            </div>
          </div>
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">Sales</div>
            <div className="text-lg font-medium truncate">
              {selectedSales
                ? `${selectedSales.code} - ${selectedSales.first_name} ${selectedSales.last_name}`
                : 'ทั้งหมด'
              }
            </div>
          </div>
        </div>
      )}

      {/* Booking List */}
      {shouldFetch && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <h2 className="text-lg font-semibold">รายการจอง</h2>
              {pagination && (
                <Badge variant="outline">
                  {pagination.total} รายการ
                </Badge>
              )}
            </div>
            {/* Export CSV Button */}
            <Button
              onClick={handleExportCSV}
              disabled={!bookings || bookings.length === 0 || isLoading}
              variant="outline"
              size="sm"
            >
              <Download className="h-4 w-4 mr-2" />
              Export CSV
            </Button>
          </div>

          {/* Results Count */}
          {pagination && pagination.total !== undefined && pagination.total !== null && (
            <div className="text-sm text-muted-foreground">
              หน้า {pagination.currentPage} จาก {pagination.totalPages} (ทั้งหมด {pagination.total.toLocaleString()} รายการ)
            </div>
          )}

          {/* Loading State */}
          {isLoading ? (
            <div className="flex items-center justify-center p-8">
              <Loader2 className="h-8 w-8 animate-spin" />
              <span className="ml-2">กำลังโหลดข้อมูล...</span>
            </div>
          ) : (
            <div className="space-y-4">
              {displayBookings.length === 0 ? (
                <div className="text-center py-8">
                  <p className="text-muted-foreground">ไม่พบการจองในช่วงวันที่นี้</p>
                </div>
              ) : (
                displayBookings.map((booking) => (
                  <BookingGridCard
                    key={booking.id}
                    booking={booking}
                    onSelect={handleCardSelect}
                    onView={handleView}
                    onViewFiles={handleViewFiles}
                    actions={adminActions}
                    showDefaultActions={true}
                  />
                ))
              )}
            </div>
          )}

          {/* Pagination */}
          {pagination && displayBookings.length > 0 && (
            <div className="mt-6">
              <AdvancedPagination
                currentPage={pagination.currentPage}
                totalPages={pagination.totalPages}
                pageSize={pageSize}
                total={pagination.total}
                onPageChange={handlePageChange}
                onPageSizeChange={handlePageSizeChange}
                showPageSizeSelector={true}
                showTotal={true}
                isLoading={isLoading}
                pageSizeOptions={[25, 50, 75, 100]}
              />
            </div>
          )}
        </div>
      )}

      {/* ถ้ายังไม่ได้เลือกช่วงวันที่ */}
      {!shouldFetch && (
        <div className="text-center py-12">
          <CalendarRange className="h-16 w-16 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-medium text-muted-foreground mb-2">
            เลือกช่วงวันที่เพื่อดูรายการจอง
          </h3>
          <p className="text-sm text-muted-foreground">
            ระบบจะแสดงรายการจองทั้งหมดในช่วงวันที่ที่เลือก
          </p>
        </div>
      )}

      {/* Detail Sheet */}
      <BookingDetailSheet
        isOpen={isDetailSheetOpen}
        onOpenChange={setIsDetailSheetOpen}
        booking={selectedBooking}
        userRole="admin"
      />

      {/* File Download Sheet */}
      <FileDownloadSheet
        isOpen={isFilesSheetOpen}
        onOpenChange={setIsFilesSheetOpen}
        files={selectedBookingForFiles?.period_details?.files || []}
        periodInfo={selectedBookingForFiles?.period_details ? {
          startDate: selectedBookingForFiles.period_details.start_date,
          endDate: selectedBookingForFiles.period_details.end_date,
          days: Math.ceil((new Date(selectedBookingForFiles.period_details.end_date).getTime() - new Date(selectedBookingForFiles.period_details.start_date).getTime()) / (1000 * 60 * 60 * 24)) + 1
        } : undefined}
        tourInfo={selectedBookingForFiles?.tour_details ? {
          code: selectedBookingForFiles.tour_details.code,
          title: selectedBookingForFiles.tour_details.title,
        } : undefined}
      />
    </div>
  );
}