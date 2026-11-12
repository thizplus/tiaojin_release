import { useState, useMemo, useCallback, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AdvancedPagination } from '@/components/ui/advanced-pagination';
import {
  AlertCircle,
  RefreshCw,
} from 'lucide-react';
import { PublicTourCard } from './PublicTourCard';
import { useToursRegular } from '@/hooks/useToursRegular';
import type { ToursRegularRequest } from '@/types/tours-regular';

interface PublicToursRegularListProps {
  initialParams?: ToursRegularRequest;
}

export function PublicToursRegularList({
  initialParams = {},
}: PublicToursRegularListProps) {
  const [searchParams] = useSearchParams();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);

  // Get search text from URL params
  const searchText = useMemo(() => searchParams.get('search') || '', [searchParams]);

  // Build query parameters
  const queryParams = useMemo(() => ({
    ...initialParams,
    search: searchText.trim() || undefined,
    status: 'publish' as const, // Only show published tours
    available_only: true, // Only show tours with available periods
    page: currentPage,
    limit: pageSize,
  }), [initialParams, searchText, currentPage, pageSize]);

  // Fetch tours data
  const {
    tours,
    pagination,
    isLoading,
    error,
    refetch
  } = useToursRegular(queryParams, { enabled: true });

  const handleRefresh = () => {
    refetch();
  };

  const handlePageChange = useCallback((page: number) => {
    setCurrentPage(page);
  }, []);

  const handlePageSizeChange = useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    setCurrentPage(1);
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
        <div className="space-y-4">
          {tours.map((tour) => (
            <PublicTourCard key={tour.id} tour={tour} />
          ))}
        </div>
      )}

      {/* Advanced Pagination */}
      {pagination && (
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
