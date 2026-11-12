import React, { useState, useEffect } from 'react';
import { TourSelector } from '@/components/admin/TourSelector';
import { PeriodSelector } from '@/components/admin/PeriodSelector';
import { BookingStatusFilter, type BookingStatus } from '@/components/admin/BookingStatusFilter';
import { BookingGridCard } from '@/components/manage/booking/BookingGridCard';
import { BookingDetailSheet } from '@/components/manage/booking/BookingDetailSheet';
import { FileDownloadSheet } from '@/components/manage/booking/FileDownloadSheet';
import { AdvancedPagination } from '@/components/ui/advanced-pagination';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Loader2, FileBarChart, Download } from 'lucide-react';
import { useAdminBooking } from '@/hooks/useAdminBooking';
import { exportBookingsToCSV } from '@/utils/exportBookingsCSV';
import type { TourPeriod } from '@/types/models/period';
import type { Tour } from '@/types/models/tour';
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

export default function ReportsToursPage() {
  const [selectedTour, setSelectedTour] = useState<Tour | null>(null);
  const [selectedPeriod, setSelectedPeriod] = useState<TourPeriod | null>(null);
  const [selectedStatuses, setSelectedStatuses] = useState<BookingStatus[]>([]);
  const [selectedBooking, setSelectedBooking] = useState<BookingResponse | null>(null);
  const [isDetailSheetOpen, setIsDetailSheetOpen] = useState(false);
  const [isFilesSheetOpen, setIsFilesSheetOpen] = useState(false);
  const [selectedBookingForFiles, setSelectedBookingForFiles] = useState<BookingResponse | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);

  // ดึงข้อมูล bookings เมื่อมี period ที่เลือก
  const {
    bookingsList: bookings,
    pagination,
    isLoading,
    setBookingFilters,
  } = useAdminBooking(
    selectedPeriod ? {
      period_id: selectedPeriod.id,
      status: selectedStatuses.length > 0 ? selectedStatuses.join(',') : undefined,
      page: currentPage,
      limit: pageSize,
    } : {},
    undefined,
    { enableList: !!selectedPeriod, enableById: false }
  );

  // อัปเดต filters เมื่อมีการเปลี่ยนแปลง
  useEffect(() => {
    if (selectedPeriod) {
      setBookingFilters({
        period_id: selectedPeriod.id,
        status: selectedStatuses.length > 0 ? selectedStatuses.join(',') : undefined,
        page: currentPage,
        limit: pageSize,
      });
    }
  }, [selectedPeriod, selectedStatuses, currentPage, pageSize, setBookingFilters]);

  // Reset เมื่อเปลี่ยนทัวร์
  useEffect(() => {
    setSelectedPeriod(null);
    setCurrentPage(1);
  }, [selectedTour]);

  // Reset เมื่อเปลี่ยน period
  useEffect(() => {
    setCurrentPage(1);
  }, [selectedPeriod]);

  // Reset pagination เมื่อเปลี่ยน status filter
  useEffect(() => {
    if (selectedPeriod) {
      setCurrentPage(1);
    }
  }, [selectedStatuses]);

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

  // Admin actions - เหมือน admin booking page
  const adminActions: BookingAction[] = [
    // จะใช้ default actions ของ BookingGridCard
  ];

  const displayBookings = bookings || [];

  // Export to CSV handler
  const handleExportCSV = React.useCallback(() => {
    if (!selectedPeriod || !selectedTour) {
      alert('กรุณาเลือกทัวร์และรอบทัวร์ก่อน');
      return;
    }

    if (!bookings || bookings.length === 0) {
      alert('ไม่มีข้อมูลสำหรับการส่งออก');
      return;
    }

    exportBookingsToCSV(bookings, {
      tourTitle: selectedTour.title,
      periodStartDate: selectedPeriod.start_date,
      periodEndDate: selectedPeriod.end_date
    });
  }, [bookings, selectedTour, selectedPeriod]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <FileBarChart className="h-8 w-8 text-primary" />
        <div>
          <h1 className="text-2xl font-bold">รายงานตามทัวร์</h1>
          <p className="text-muted-foreground">ดูรายการจองตามทัวร์และรอบที่เลือก</p>
        </div>
      </div>

      {/* Tour และ Period Selectors */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <TourSelector
          selectedTour={selectedTour}
          onTourSelect={setSelectedTour}
        />
        <PeriodSelector
          tourId={selectedTour?.id || null}
          selectedPeriod={selectedPeriod}
          onPeriodSelect={setSelectedPeriod}
          allowFullyBooked={true}
        />
      </div>

      {/* Status Filter */}
      {selectedPeriod && (
        <BookingStatusFilter
          selectedStatuses={selectedStatuses}
          onStatusChange={setSelectedStatuses}
        />
      )}

      {/* Summary Information */}
      {selectedPeriod && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">ที่นั่งทั้งหมด</div>
            <div className="text-2xl font-bold">{selectedPeriod.seat}</div>
          </div>
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">จองแล้ว</div>
            <div className="text-2xl font-bold text-chart-2">{selectedPeriod.booking || 0}</div>
          </div>
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">ที่นั่งว่าง</div>
            <div className="text-2xl font-bold text-foreground">{selectedPeriod.available_seats || 0}</div>
          </div>
          <div className="p-4 border rounded-lg bg-card">
            <div className="text-sm text-muted-foreground">อัตราการจอง</div>
            <div className="text-2xl font-bold text-chart-2">
              {selectedPeriod.seat > 0 ? Math.round(((selectedPeriod.booking || 0) / selectedPeriod.seat) * 100) : 0}%
            </div>
          </div>
        </div>
      )}

      {/* Booking List */}
      {selectedPeriod && (
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
                  <p className="text-muted-foreground">ไม่พบการจองในรอบทัวร์นี้</p>
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
                    showDefaultActions={true} // แสดง admin default actions
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

      {/* ถ้ายังไม่ได้เลือกรอบทัวร์ */}
      {!selectedPeriod && (
        <div className="text-center py-12">
          <FileBarChart className="h-16 w-16 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-medium text-muted-foreground mb-2">
            เลือกทัวร์และรอบทัวร์เพื่อดูรายการจอง
          </h3>
          <p className="text-sm text-muted-foreground">
            ระบบจะแสดงรายการจองทั้งหมดในรอบทัวร์ที่เลือก
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