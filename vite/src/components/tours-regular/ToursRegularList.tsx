import { useState, useMemo, useCallback, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Card, CardContent } from '@/components/ui/card';

import { Button } from '@/components/ui/button';
import { AdvancedPagination } from '@/components/ui/advanced-pagination';

import {
  AlertCircle,

  RefreshCw,

} from 'lucide-react';
import { TourCard } from './TourCard';
import { useToursRegularWithWebSocket } from '@/hooks/useToursRegular';
import { useToursHotWithWebSocket } from '@/hooks/useToursHot';
import type { Tour, Period, ToursRegularRequest } from '@/types/tours-regular';

interface ToursRegularListProps {
  onBookPeriod?: (period: Period, tour: Tour) => void;
  canBook?: boolean;
  userStatus?: string;
  userRole?: string; // 'admin' | 'super' | 'agent' | etc.
  initialParams?: ToursRegularRequest;
  useToursHook?: 'regular' | 'hot'; // เพิ่ม prop เพื่อเลือก hook
}

export function ToursRegularList({
  onBookPeriod,
  canBook = false,
  userStatus = 'pending',
  userRole,
  initialParams = {},
  useToursHook = 'regular' // default ใช้ regular
}: ToursRegularListProps) {
  const [searchParams] = useSearchParams();
  const [statusFilter] = useState<'all' | 'publish' | 'draft'>('all');
  const [availableOnlyFilter] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);

  // Get search text from URL params
  const searchText = useMemo(() => searchParams.get('search') || '', [searchParams]);

  // Build query parameters
  const queryParams = useMemo(() => ({
    ...initialParams,
    search: searchText.trim() || undefined,
    status: statusFilter === 'all' ? undefined : statusFilter,
    available_only: availableOnlyFilter || undefined,
    page: currentPage,
    limit: pageSize,
  }), [initialParams, searchText, statusFilter, availableOnlyFilter, currentPage, pageSize]);

  // Fetch tours data - เลือกใช้ hook ตาม useToursHook prop
  const toursRegularResult = useToursRegularWithWebSocket(
    useToursHook === 'regular' ? queryParams : {},
    { enabled: useToursHook === 'regular' }
  );

  const toursHotResult = useToursHotWithWebSocket(
    useToursHook === 'hot' ? queryParams : {},
    { enabled: useToursHook === 'hot' }
  );

  // เลือกผลลัพธ์ตาม hook ที่ใช้
  const {
    tours,
    pagination,
    isLoading,
    error,
    refetch
  } = useToursHook === 'hot' ? toursHotResult : toursRegularResult;


  const handleRefresh = () => {
    refetch();
  };

  const handlePageChange = useCallback((page: number) => {
    setCurrentPage(page);
  }, []);

  const handlePageSizeChange = useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    setCurrentPage(1); // Reset to first page
  }, []);

  // Reset to first page when search changes
  useEffect(() => {
    setCurrentPage(1);
  }, [searchText]);



  if (error) {
    return (
      <Card>
        <CardContent className="p-6 text-center">
          <AlertCircle className="h-8 w-8 mx-auto mb-2 text-destructive" />
          <p className="text-destructive mb-4">เกิดข้อผิดพลาดในการโหลดข้อมูล</p>
          <Button onClick={handleRefresh} variant="outline">
            <RefreshCw className="h-4 w-4 mr-2" />
            ลองใหม่
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
    

      {/* Tours List */}
      {isLoading ? (
        <Card>
          <CardContent className="p-6 text-center">
            <RefreshCw className="h-8 w-8 mx-auto mb-2 animate-spin text-muted-foreground" />
            <p className="text-muted-foreground">กำลังโหลดข้อมูลทัวร์...</p>
          </CardContent>
        </Card>
      ) : tours.length === 0 ? (
        <Card>
          <CardContent className="p-6 text-center">
            <AlertCircle className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
            <p className="text-muted-foreground">ไม่พบทัวร์ที่ตรงกับเงื่อนไขการค้นหา</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4 ">
          {tours.map((tour) => (
            <TourCard
              key={tour.id}
              tour={tour}
              onBookPeriod={onBookPeriod}
              canBook={canBook}
              userStatus={userStatus}
              userRole={userRole}
            />
          ))}
        </div>
      )}

      {/* Advanced Pagination */}
      {pagination  && (
        <div className="mt-6">
          <AdvancedPagination
            currentPage={pagination.current_page}
            totalPages={pagination.total_pages}
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
  );
}