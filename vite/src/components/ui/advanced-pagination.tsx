
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";

interface AdvancedPaginationProps {
  currentPage: number;
  totalPages: number;
  pageSize: number;
  total?: number;
  pageSizeOptions?: number[];
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  showPageSizeSelector?: boolean;
  showTotal?: boolean;
  isLoading?: boolean;
}

export function AdvancedPagination({
  currentPage,
  totalPages,
  pageSize,
  total,
  pageSizeOptions = [25, 50, 75, 100],
  onPageChange,
  onPageSizeChange,
  showPageSizeSelector = true,
  showTotal = true,
  isLoading = false
}: AdvancedPaginationProps) {
  const canPreviousPage = currentPage > 1;
  const canNextPage = currentPage < totalPages;

  const handleFirstPage = () => {
    if (canPreviousPage && !isLoading) {
      onPageChange(1);
    }
  };

  const handlePreviousPage = () => {
    if (canPreviousPage && !isLoading) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    if (canNextPage && !isLoading) {
      onPageChange(currentPage + 1);
    }
  };

  const handleLastPage = () => {
    if (canNextPage && !isLoading) {
      onPageChange(totalPages);
    }
  };

  const handlePageSizeChange = (value: string) => {
    const newPageSize = Number(value);
    if (!isLoading) {
      onPageSizeChange(newPageSize);
      // Reset to first page when changing page size
      onPageChange(1);
    }
  };

  // Calculate displayed range
  const startItem = totalPages > 0 ? ((currentPage - 1) * pageSize) + 1 : 0;
  const endItem = Math.min(currentPage * pageSize, total || 0);

  if (totalPages <= 1 && !showPageSizeSelector) {
    return null;
  }

  return (
    <div className="flex flex-col gap-4 py-4 md:flex-row md:items-center md:justify-between">
      <div className="flex flex-col gap-4 items-center md:flex-row md:gap-6">
        {/* Page Size Selector */}
        {showPageSizeSelector && (
          <div className="flex items-center gap-2">
            <Label htmlFor="rows-per-page" className="text-sm font-medium">
              รายการต่อหน้า
            </Label>
            <Select
              value={`${pageSize}`}
              onValueChange={handlePageSizeChange}
              disabled={isLoading}
            >
              <SelectTrigger size="sm" className="w-20" id="rows-per-page">
                <SelectValue />
              </SelectTrigger>
              <SelectContent side="top">
                {pageSizeOptions.map((size) => (
                  <SelectItem key={size} value={`${size}`}>
                    {size}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}

        {/* Current Page Info */}
        <div className="flex items-center justify-center text-sm font-medium whitespace-nowrap">
          หน้า {currentPage} จาก {totalPages}
        </div>

        {/* Navigation Buttons */}
        <div className="flex items-center gap-1 md:gap-2">
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={handleFirstPage}
            disabled={!canPreviousPage || isLoading}
          >
            <span className="sr-only">ไปยังหน้าแรก</span>
            <ChevronsLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={handlePreviousPage}
            disabled={!canPreviousPage || isLoading}
          >
            <span className="sr-only">ไปยังหน้าก่อนหน้า</span>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={handleNextPage}
            disabled={!canNextPage || isLoading}
          >
            <span className="sr-only">ไปยังหน้าถัดไป</span>
            <ChevronRight className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={handleLastPage}
            disabled={!canNextPage || isLoading}
          >
            <span className="sr-only">ไปยังหน้าสุดท้าย</span>
            <ChevronsRight className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Total Items Info */}
      {showTotal && total !== undefined && (
        <div className="text-sm text-muted-foreground">
          แสดง {startItem.toLocaleString()} ถึง {endItem.toLocaleString()} จากทั้งหมด {total.toLocaleString()} รายการ
        </div>
      )}
    </div>
  );
}