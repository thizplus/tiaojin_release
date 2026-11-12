import React, { useState, useEffect } from 'react';
import { useSaleBooking } from '@/hooks/useSaleBooking';
import { BookingGridCard } from '@/components/manage/booking/BookingGridCard';
import { BookingDetailSheet } from '@/components/manage/booking/BookingDetailSheet';
import { FileDownloadSheet } from '@/components/manage/booking/FileDownloadSheet';
import { PaymentSlipUploadSheet } from '@/components/manage/booking/PaymentSlipUploadSheet';
import { PassportUploadSheet } from '@/components/manage/booking/PassportUploadSheet';
import { Badge } from '@/components/ui/badge';
import { AdvancedPagination } from '@/components/ui/advanced-pagination';
import { IsolatedSearchInput } from '@/components/ui/isolated-search-input';
import { Loader2, Filter, SortAsc, SortDesc, Upload, FileText } from 'lucide-react';
import { Button } from '@/components/ui/button';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { BookingResponse } from '@/types';

// Type สำหรับ action object
interface BookingAction {
  key: string;
  label: string;
  icon: React.ReactNode;
  onClick: (booking: BookingResponse) => void;
  showFor: (booking: BookingResponse) => boolean;
  variant?: 'destructive';
  separator?: boolean;
}

export default function BookingsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired'>('all');
  const [sortBy, setSortBy] = useState<'created_at' | 'sum_net' | 'customer_name'>('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);
  const [selectedBooking, setSelectedBooking] = useState<BookingResponse | null>(null);
  const [isDetailSheetOpen, setIsDetailSheetOpen] = useState(false);
  const [isFilesSheetOpen, setIsFilesSheetOpen] = useState(false);
  const [selectedBookingForFiles, setSelectedBookingForFiles] = useState<BookingResponse | null>(null);
  const [isPaymentSlipSheetOpen, setIsPaymentSlipSheetOpen] = useState(false);
  const [isPassportSheetOpen, setIsPassportSheetOpen] = useState(false);
  const [actionBooking, setActionBooking] = useState<BookingResponse | null>(null);

  const {
    bookingsList: bookings,
    pagination,
    isLoading,
    refetchBookingsList: refetch,
    setBookingFilters,
    uploadPaymentSlip,
    uploadPassport,
    isUploadPaymentSlipLoading,
    isUploadPassportLoading,
  } = useSaleBooking(
    {
      search: searchTerm,
      status: statusFilter === 'all' ? undefined : statusFilter as 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired',
      sort: sortBy,
      order: sortOrder,
      page: currentPage,
      limit: pageSize,
    },
    undefined,
    { enableList: true, enableById: false }
  );

  // Update filters when state changes
  useEffect(() => {
    setBookingFilters({
      search: searchTerm,
      status: statusFilter === 'all' ? undefined : statusFilter as 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired',
      sort: sortBy,
      order: sortOrder,
      page: currentPage,
      limit: pageSize,
    });
  }, [searchTerm, statusFilter, sortBy, sortOrder, currentPage, pageSize, setBookingFilters]);

  // Reset to first page when filters change (แยก searchTerm ออก)
  useEffect(() => {
    setCurrentPage(1);
  }, [statusFilter, sortBy, sortOrder]);

  // Reset page เมื่อ search query เปลี่ยนเท่านั้น
  useEffect(() => {
    if (searchTerm !== '') {
      setCurrentPage(1);
    }
  }, [searchTerm]);

  const handleView = (booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDetailSheetOpen(true);
    console.log('View booking:', booking.id);
  };

  const handleCardSelect = (booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDetailSheetOpen(true);
  };

  const handleViewFiles = (booking: BookingResponse) => {
    setSelectedBookingForFiles(booking);
    setIsFilesSheetOpen(true);
  };

  const handleUploadPaymentSlip = (booking: BookingResponse) => {
    setActionBooking(booking);
    setIsPaymentSlipSheetOpen(true);
  };

  const handleUploadPassport = (booking: BookingResponse) => {
    setActionBooking(booking);
    setIsPassportSheetOpen(true);
  };

  const handlePaymentSlipUpload = async (id: string, file: File, amount: number) => {
    await uploadPaymentSlip(id, { payment_slip: file, amount });
  };

  const handlePassportUpload = async (id: string, file: File) => {
    await uploadPassport(id, { passport: file });
  };

  const handlePageChange = React.useCallback((page: number) => {
    setCurrentPage(page);
  }, []);

  const handlePageSizeChange = React.useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    setCurrentPage(1); // Reset to first page
  }, []);

  const toggleSort = (field: typeof sortBy) => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('desc');
    }
  };

  // Sale-specific actions - แนบสลิป/passport แทนลูกค้า
  const saleActions: BookingAction[] = [
    {
      key: 'upload-payment-slip',
      label: 'แจ้งชำระเงินแทนลูกค้า',
      icon: <Upload className="mr-1 h-3.5 w-3.5" />,
      onClick: handleUploadPaymentSlip,
      showFor: (booking: BookingResponse) => ['pending', 'processing'].includes(booking.status)
    },
    {
      key: 'upload-passport',
      label: 'แนบ Passport แทนลูกค้า',
      icon: <FileText className="mr-1 h-3.5 w-3.5" />,
      onClick: handleUploadPassport,
      showFor: (booking: BookingResponse) => ['pending', 'processing', 'confirm'].includes(booking.status)
    }
  ];

  // Remove client-side filtering since it's handled by the server
  const displayBookings = bookings || [];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin" />
        <span className="ml-2">กำลังโหลดข้อมูล...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold">รายการจองทัวร์ (เซลล์)</h1>
        <p className="text-muted-foreground">จัดการการจองทัวร์และแนบเอกสารแทนลูกค้า</p>
      </div>

      {/* Filters and Search */}
      <div className="space-y-4">
        {/* Search */}
        <IsolatedSearchInput
          placeholder="ค้นหาด้วยชื่อลูกค้า, เบอร์โทร, เลขที่ใบจอง, หรือรหัสทัวร์..."
          onSearchChange={setSearchTerm}
        />

        {/* Status Filter */}
        <div className="flex flex-wrap gap-2">
          <Badge
            variant={statusFilter === 'all' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('all')}
          >
            ทั้งหมด {(pagination?.total !== undefined && pagination?.total !== null) ? `(${pagination.total.toLocaleString()})` : ''}
          </Badge>
          <Badge
            variant={statusFilter === 'pending' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('pending')}
          >
            รอดำเนินการ
          </Badge>
          <Badge
            variant={statusFilter === 'processing' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('processing')}
          >
            กำลังดำเนินการ
          </Badge>
          <Badge
            variant={statusFilter === 'confirm' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('confirm')}
          >
            ยืนยันแล้ว
          </Badge>
          <Badge
            variant={statusFilter === 'completed' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('completed')}
          >
            เสร็จสิ้น
          </Badge>
          <Badge
            variant={statusFilter === 'cancel' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('cancel')}
          >
            ยกเลิก
          </Badge>
          <Badge
            variant={statusFilter === 'expired' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('expired')}
          >
            หมดอายุ
          </Badge>
        </div>

        {/* Sort Options */}
        <div className="flex gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm">
                <Filter className="h-4 w-4 mr-2" />
                เรียงตาม
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => toggleSort('created_at')}>
                <div className="flex items-center justify-between w-full">
                  <span>วันที่จอง</span>
                  {sortBy === 'created_at' && (
                    sortOrder === 'asc' ? <SortAsc className="h-4 w-4" /> : <SortDesc className="h-4 w-4" />
                  )}
                </div>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toggleSort('sum_net')}>
                <div className="flex items-center justify-between w-full">
                  <span>ยอดเงิน</span>
                  {sortBy === 'sum_net' && (
                    sortOrder === 'asc' ? <SortAsc className="h-4 w-4" /> : <SortDesc className="h-4 w-4" />
                  )}
                </div>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toggleSort('customer_name')}>
                <div className="flex items-center justify-between w-full">
                  <span>ชื่อลูกค้า</span>
                  {sortBy === 'customer_name' && (
                    sortOrder === 'asc' ? <SortAsc className="h-4 w-4" /> : <SortDesc className="h-4 w-4" />
                  )}
                </div>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Results Count */}
      {pagination && pagination.total !== undefined && pagination.total !== null && (
        <div className="text-sm text-muted-foreground">
          หน้า {pagination.currentPage} จาก {pagination.totalPages} (ทั้งหมด {pagination.total.toLocaleString()} รายการ)
        </div>
      )}

      {/* Booking Cards */}
      <div className="space-y-4">
        {displayBookings.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-muted-foreground">ไม่พบข้อมูลการจอง</p>
          </div>
        ) : (
          displayBookings.map((booking) => (
            <BookingGridCard
              key={booking.id}
              booking={booking}
              onSelect={handleCardSelect}
              onView={handleView}
              onViewFiles={handleViewFiles}
              requireFullPaymentForFiles
              actions={saleActions}
              showDefaultActions={false}
            />
          ))
        )}
      </div>

      {/* Pagination */}
      {pagination && (
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

      {/* Detail Sheet */}
      <BookingDetailSheet
        isOpen={isDetailSheetOpen}
        onOpenChange={setIsDetailSheetOpen}
        booking={selectedBooking}
        userRole="sale"
        onUploadPaymentSlip={(booking) => {
          setActionBooking(booking);
          setIsDetailSheetOpen(false);
          setIsPaymentSlipSheetOpen(true);
        }}
        onUploadPassport={(booking) => {
          setActionBooking(booking);
          setIsDetailSheetOpen(false);
          setIsPassportSheetOpen(true);
        }}
        onRefetch={refetch}
      />

      {/* Payment Slip Upload Sheet */}
      <PaymentSlipUploadSheet
        isOpen={isPaymentSlipSheetOpen}
        onOpenChange={setIsPaymentSlipSheetOpen}
        booking={actionBooking}
        onUpload={handlePaymentSlipUpload}
        isLoading={isUploadPaymentSlipLoading}
      />

      {/* Passport Upload Sheet */}
      <PassportUploadSheet
        isOpen={isPassportSheetOpen}
        onOpenChange={setIsPassportSheetOpen}
        booking={actionBooking}
        onUpload={handlePassportUpload}
        isLoading={isUploadPassportLoading}
        onRefetch={refetch}
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
        isFullyPaid={selectedBookingForFiles
          ? selectedBookingForFiles.sum_payment >= (selectedBookingForFiles.grand_total || selectedBookingForFiles.sum_net)
          : undefined}
      />
    </div>
  );
}