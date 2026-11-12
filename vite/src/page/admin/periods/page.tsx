import { useState, useCallback, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Calendar as CalendarComponent } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ArrowLeft, Calendar, Plus, CalendarIcon, X } from "lucide-react";
import { IsolatedSearchInput } from "@/components/ui/isolated-search-input";
import { TourPeriodsTable } from "@/components/manage/period/tour-periods-table";
import { PeriodSheet } from "@/components/manage/period/PeriodSheet";
import { usePeriodsData } from "@/hooks/usePeriodsData";
import { usePeriod } from "@/hooks/usePeriod";
import { useToast } from "@/hooks/use-toast";
import { cn } from "@/lib/utils";
import type { TourPeriod } from "@/types";


export default function PeriodPage() {
  const { tourId } = useParams<{ tourId: string }>();
  const navigate = useNavigate();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create');
  const [selectedPeriod, setSelectedPeriod] = useState<TourPeriod | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'draft' | 'publish'>('all');
  const [dateFrom, setDateFrom] = useState<Date | undefined>(undefined);
  const [dateTo, setDateTo] = useState<Date | undefined>(undefined);
  const { toast } = useToast();
  const { handleDeletePeriod: deletePeriod } = usePeriod();


  const {
    periods,
    loading,
    pagination,
    currentPage,
    currentPageSize,
    totalPages,
    goToPage,
    changePageSize,
    refresh,
    setFilters,
  } = usePeriodsData(
    { limit: 10 },
    { tourId: tourId || undefined, autoFetch: !!tourId }
  );

  // Update filters when search/filter changes
  useEffect(() => {
    const filters: Partial<{
      status: 'draft' | 'publish';
      date_from: string;
      date_to: string;
    }> = {};

    if (statusFilter !== 'all') {
      filters.status = statusFilter;
    }

    if (dateFrom) {
      // Convert Date to YYYY-MM-DD format for API
      filters.date_from = format(dateFrom, 'yyyy-MM-dd');
    }

    if (dateTo) {
      // Convert Date to YYYY-MM-DD format for API
      filters.date_to = format(dateTo, 'yyyy-MM-dd');
    }

    // Note: API doesn't support text search, but we can filter client-side if needed
    setFilters(filters);
  }, [statusFilter, dateFrom, dateTo, setFilters]);

 

  const handleBackToTours = () => {
    navigate("/admin/tours");
  };

  const handleCreatePeriod = () => {
    setDialogMode('create');
    setSelectedPeriod(null);
    setDialogOpen(true);
  };

  const handleEditPeriod = useCallback((period: TourPeriod) => {
    setDialogMode('edit');
    setSelectedPeriod(period);
    setDialogOpen(true);
  }, []);

  const handleDialogSuccess = () => {
    refresh();
  };

  const handleDeletePeriod = async (period: TourPeriod) => {
    if (window.confirm(`คุณต้องการลบช่วงเวลาทัวร์ ${period.start_date} - ${period.end_date} หรือไม่?`)) {
      const result = await deletePeriod(period.id);
      if (result.success) {
        toast({
          title: 'สำเร็จ',
          description: 'ลบช่วงเวลาทัวร์เรียบร้อยแล้ว',
        });
        refresh();
      } else {
        toast({
          variant: 'destructive',
          title: 'เกิดข้อผิดพลาด',
          description: result.error || 'ไม่สามารถลบช่วงเวลาทัวร์ได้',
        });
      }
    }
  };

  if (!tourId) {
    return (
      <div className="text-center py-10">
        <p>Tour ID not found</p>
      </div>
    );
  }

  // Filter periods client-side if search term is provided
  const filteredPeriods = searchTerm
    ? periods.filter((period) => {
        const searchLower = searchTerm.toLowerCase();
        const startDate = new Date(period.start_date).toLocaleDateString('th-TH');
        const endDate = new Date(period.end_date).toLocaleDateString('th-TH');
        const tourCode = period.tour_details?.code?.toLowerCase() || '';
        const tourTitle = period.tour_details?.title?.toLowerCase() || '';

        return (
          startDate.includes(searchLower) ||
          endDate.includes(searchLower) ||
          tourCode.includes(searchLower) ||
          tourTitle.includes(searchLower)
        );
      })
    : periods;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleBackToTours}
              className="hover:bg-muted"
            >
              <ArrowLeft className="h-4 w-4 mr-2" />
              กลับไปหน้าทัวร์
            </Button>
          </div>
          <h1 className="text-3xl font-bold">ช่วงเวลาทัวร์</h1>
        </div>
        <Button onClick={handleCreatePeriod}>
          <Plus className="h-4 w-4 mr-2" />
          เพิ่มช่วงเวลาทัวร์
        </Button>
      </div>

      {/* Search and Filters */}
      <div className="space-y-4">
        {/* Search */}
        <IsolatedSearchInput
          placeholder="ค้นหาด้วยวันที่, รหัสทัวร์, หรือชื่อทัวร์..."
          onSearchChange={setSearchTerm}
        />

        {/* Status Filter */}
        <div className="flex flex-wrap gap-2 items-center">
          <span className="text-sm font-medium text-muted-foreground">สถานะ:</span>
          <Badge
            variant={statusFilter === 'all' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('all')}
          >
            ทั้งหมด {pagination?.total !== undefined ? `(${pagination.total})` : ''}
          </Badge>
          <Badge
            variant={statusFilter === 'publish' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('publish')}
          >
            เผยแพร่
          </Badge>
          <Badge
            variant={statusFilter === 'draft' ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => setStatusFilter('draft')}
          >
            แบบร่าง
          </Badge>
        </div>

        {/* Date Range Filter */}
        <div className="flex flex-wrap gap-4 items-center">
          <span className="text-sm font-medium text-muted-foreground">ช่วงวันที่:</span>
          <div className="flex gap-2 items-center flex-wrap">
            {/* Date From */}
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className={cn(
                    "w-[200px] justify-start text-left font-normal",
                    !dateFrom && "text-muted-foreground"
                  )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {dateFrom ? format(dateFrom, "dd MMM yyyy", { locale: th }) : "เลือกวันที่เริ่มต้น"}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <CalendarComponent
                  mode="single"
                  selected={dateFrom}
                  onSelect={setDateFrom}
                  initialFocus
                  locale={th}
                />
              </PopoverContent>
            </Popover>

            <span className="text-muted-foreground">ถึง</span>

            {/* Date To */}
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className={cn(
                    "w-[200px] justify-start text-left font-normal",
                    !dateTo && "text-muted-foreground"
                  )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {dateTo ? format(dateTo, "dd MMM yyyy", { locale: th }) : "เลือกวันที่สิ้นสุด"}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <CalendarComponent
                  mode="single"
                  selected={dateTo}
                  onSelect={setDateTo}
                  initialFocus
                  disabled={(date) => dateFrom ? date < dateFrom : false}
                  locale={th}
                />
              </PopoverContent>
            </Popover>

            {/* Clear Button */}
            {(dateFrom || dateTo) && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setDateFrom(undefined);
                  setDateTo(undefined);
                }}
              >
                <X className="h-4 w-4 mr-1" />
                ล้างวันที่
              </Button>
            )}
          </div>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Calendar className="h-5 w-5" />
              <span>รายการช่วงเวลาทัวร์</span>
            </div>
            {searchTerm && (
              <span className="text-sm font-normal text-muted-foreground">
                ผลการค้นหา: {filteredPeriods.length} รายการ
              </span>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-10">
              <p>กำลังโหลด...</p>
            </div>
          ) : filteredPeriods.length === 0 ? (
            <div className="text-center py-10">
              <Calendar className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium text-muted-foreground">
                {searchTerm ? 'ไม่พบผลการค้นหา' : 'ไม่มีช่วงเวลาทัวร์'}
              </h3>
              <p className="text-sm text-muted-foreground mt-1 mb-4">
                {searchTerm
                  ? 'ลองเปลี่ยนคำค้นหาหรือปรับตัวกรอง'
                  : 'ยังไม่มีการเพิ่มช่วงเวลาสำหรับทัวร์นี้'}
              </p>
              {!searchTerm && (
                <Button onClick={handleCreatePeriod}>
                  <Plus className="h-4 w-4 mr-2" />
                  เพิ่มช่วงเวลาทัวร์แรก
                </Button>
              )}
            </div>
          ) : (
            <TourPeriodsTable
              data={filteredPeriods}
              page={currentPage - 1}
              pageSize={currentPageSize}
              totalPages={totalPages}
              setPage={(page) => goToPage(page + 1)}
              setPageSize={changePageSize}
              onEdit={handleEditPeriod}
              onDelete={handleDeletePeriod}
            />
          )}
        </CardContent>
      </Card>

      <PeriodSheet
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        tourId={tourId}
        period={selectedPeriod}
        mode={dialogMode}
        onSuccess={handleDialogSuccess}
      />
    </div>
  );
}