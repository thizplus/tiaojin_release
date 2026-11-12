import React, { useState, useEffect } from 'react';
import { useAdminBooking } from '@/hooks/useAdminBooking';
import { BookingGridCard } from './BookingGridCard';
import { BookingDetailSheet } from './BookingDetailSheet';
import { FileDownloadSheet } from './FileDownloadSheet';
import { PassportUploadSheet } from './PassportUploadSheet';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { AdvancedPagination } from '@/components/ui/advanced-pagination';
import { IsolatedSearchInput } from '@/components/ui/isolated-search-input';
import { SalesSelect } from '@/components/ui/sales-select';
import { Loader2, Filter, SortAsc, SortDesc } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import adminBookingService from '@/services/adminBookingService';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
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

interface BookingCardListProps {
  title?: string;
  description?: string;
  onView?: (booking: BookingResponse) => void;
  actions?: BookingAction[];
  showDefaultActions?: boolean;
  // Legacy props for backward compatibility
  onDelete?: (booking: BookingResponse) => void;
  onCreatePayment?: (booking: BookingResponse) => void;
}

const BookingCardListContent = React.memo(function BookingCardListContent({
  title = "จัดการการจองทัวร์",
  description = "จัดการข้อมูลการจองทัวร์และการชำระเงิน",
  onView,
  actions = [],
  showDefaultActions = true,
  // Legacy props for backward compatibility
  onDelete,
  onCreatePayment
}: BookingCardListProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired' | 'payment_slip_pending'>('all');
  const [salesFilter, setSalesFilter] = useState<string | undefined>(undefined);
  const [sortBy, setSortBy] = useState<'created_at' | 'sum_net' | 'customer_name'>('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);
  const [selectedBooking, setSelectedBooking] = useState<BookingResponse | null>(null);
  const [isDetailSheetOpen, setIsDetailSheetOpen] = useState(false);
  const [isFilesSheetOpen, setIsFilesSheetOpen] = useState(false);
  const [selectedBookingForFiles, setSelectedBookingForFiles] = useState<BookingResponse | null>(null);
  const [isPassportSheetOpen, setIsPassportSheetOpen] = useState(false);
  const [passportBooking, setPassportBooking] = useState<BookingResponse | null>(null);
  const [isUploadingPassport, setIsUploadingPassport] = useState(false);
  const { toast } = useToast();

  const {
    bookingsList: bookings,
    pagination,
    isLoading,
    refetchBookingsList: refetch,
    setBookingFilters,
    verifyPaymentSlip,
    reversePaymentSlip
  } = useAdminBooking(
    {
      search: searchTerm,
      status: statusFilter === 'all' || statusFilter === 'payment_slip_pending' ? undefined : statusFilter as 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired',
      payment_slip_status: statusFilter === 'payment_slip_pending' ? 'pending' : undefined,
      sales_id: salesFilter,
      sort: sortBy,
      order: sortOrder,
      page: currentPage,
      limit: pageSize,
    },
    undefined,
    { enableList: true, enableStats: false, enableById: false }
  );

  // Update selectedBooking when bookings list updates after refetch
  useEffect(() => {
    if (selectedBooking && bookings && bookings.length > 0 && isDetailSheetOpen) {
      const updatedBooking = bookings.find(b => b.id === selectedBooking.id);
      if (updatedBooking) {
        // Check if payment_slips actually changed
        const oldSlipsStr = JSON.stringify(selectedBooking.payment_slips || []);
        const newSlipsStr = JSON.stringify(updatedBooking.payment_slips || []);

        if (oldSlipsStr !== newSlipsStr) {
          setSelectedBooking(updatedBooking);
        }
      }
    }
  }, [bookings, selectedBooking, isDetailSheetOpen]);

  // Update filters when state changes
  useEffect(() => {
    setBookingFilters({
      search: searchTerm,
      status: statusFilter === 'all' || statusFilter === 'payment_slip_pending' ? undefined : statusFilter as 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired',
      payment_slip_status: statusFilter === 'payment_slip_pending' ? 'pending' : undefined,
      sales_id: salesFilter,
      sort: sortBy,
      order: sortOrder,
      page: currentPage,
      limit: pageSize,
    });
  }, [searchTerm, statusFilter, salesFilter, sortBy, sortOrder, currentPage, pageSize, setBookingFilters]);

  // Reset to first page when filters change (แยก searchTerm ออก)
  useEffect(() => {
    setCurrentPage(1);
  }, [statusFilter, salesFilter, sortBy, sortOrder]);

  // Reset page เมื่อ search query เปลี่ยนเท่านั้น
  useEffect(() => {
    if (searchTerm !== '') {
      setCurrentPage(1);
    }
  }, [searchTerm]);

  const handleCardSelect = (booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDetailSheetOpen(true);
  };

  const handleView = (booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDetailSheetOpen(true);
    onView?.(booking);
  };

  const handleDelete = (booking: BookingResponse) => {
    onDelete?.(booking);
    refetch();
  };

  const handleCreatePayment = (booking: BookingResponse) => {
    onCreatePayment?.(booking);
    refetch();
  };

  const handleViewFiles = (booking: BookingResponse) => {
    setSelectedBookingForFiles(booking);
    setIsFilesSheetOpen(true);
  };

  const handleUploadPassport = (booking: BookingResponse) => {
    setPassportBooking(booking);
    setIsDetailSheetOpen(false);
    setIsPassportSheetOpen(true);
  };

  const handlePassportUpload = async (bookingId: string, file: File) => {
    setIsUploadingPassport(true);
    try {
      await adminBookingService.uploadPassport(bookingId, file);
      toast({
        title: 'อัปโหลดสำเร็จ',
        description: 'อัปโหลดหนังสือเดินทางเรียบร้อยแล้ว',
      });
      refetch();
    } catch (error) {
      toast({
        title: 'เกิดข้อผิดพลาด',
        description: error instanceof Error ? error.message : 'ไม่สามารถอัปโหลดหนังสือเดินทางได้',
        variant: 'destructive',
      });
      throw error;
    } finally {
      setIsUploadingPassport(false);
    }
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
        <h1 className="text-2xl font-bold">{title}</h1>
        <p className="text-muted-foreground">{description}</p>
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
          <Badge
            variant={statusFilter === 'payment_slip_pending' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('payment_slip_pending')}
          >
            สลิปรออนุมัติ
          </Badge>
        </div>

        {/* Sales Filter และ Sort Options */}
        <div className="flex flex-col sm:flex-row gap-4">
          {/* Sales Filter */}
          <div className="w-full sm:w-80">
            <SalesSelect
              value={salesFilter}
              onValueChange={setSalesFilter}
              placeholder="กรองตามพนักงานขาย..."
            />
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
              actions={actions}
              showDefaultActions={showDefaultActions}
              onDelete={handleDelete}
              onCreatePayment={handleCreatePayment}
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
        userRole="admin"
        onDelete={handleDelete}
        onCreatePayment={handleCreatePayment}
        onUploadPassport={handleUploadPassport}
        onVerifyPaymentSlip={verifyPaymentSlip}
        onReversePaymentSlip={reversePaymentSlip}
        onRefetch={refetch}
      />

      {/* Passport Upload Sheet */}
      <PassportUploadSheet
        isOpen={isPassportSheetOpen}
        onOpenChange={setIsPassportSheetOpen}
        booking={passportBooking}
        onUpload={handlePassportUpload}
        isLoading={isUploadingPassport}
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
      />
    </div>
  );
});

export const BookingCardList = BookingCardListContent;