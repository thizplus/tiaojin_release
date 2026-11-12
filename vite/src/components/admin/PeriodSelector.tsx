import { useState, useEffect } from 'react';
import { Calendar, Search, Users, AlertCircle } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2 } from 'lucide-react';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import { usePeriodSearch } from '@/hooks/usePeriodSearch';
import type { TourPeriod } from '@/types/models/period';

interface PeriodSelectorProps {
  tourId: string | null;
  selectedPeriod: TourPeriod | null;
  onPeriodSelect: (period: TourPeriod | null) => void;
  /** อนุญาตให้เลือก period ที่เต็มแล้ว (สำหรับหน้า reports) */
  allowFullyBooked?: boolean;
}

export function PeriodSelector({ tourId, selectedPeriod, onPeriodSelect, allowFullyBooked = false }: PeriodSelectorProps) {
  const [isOpen, setIsOpen] = useState(false);
 

  const {
    periods,
    searchTerm,
    setSearchTerm,
    isLoading,
    clearSearch,
  } = usePeriodSearch({
    tourId,
    limit: 20,
    debounceMs: 300,
    autoSearch: true,
  });

  // Reset selected period when tourId changes
  useEffect(() => {
    if (!tourId) {
      onPeriodSelect(null);
    }
  }, [tourId, onPeriodSelect]);

  const handlePeriodSelect = (period: TourPeriod) => {
    onPeriodSelect(period);
    setIsOpen(false);
    clearSearch();
  };

  const handleClear = () => {
    onPeriodSelect(null);
    clearSearch();
  };

  const formatDateRange = (startDate: string, endDate: string) => {
    const start = new Date(startDate);
    const end = new Date(endDate);
    const days = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;

    return {
      range: `${format(start, 'dd MMM', { locale: th })} - ${format(end, 'dd MMM yyyy', { locale: th })}`,
      days
    };
  };


  return (
    <div className="space-y-2">
      <label className="text-sm font-medium">เลือกรอบทัวร์</label>
      <div className="relative">
        <div className="flex gap-2">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={tourId ? "ค้นหารอบทัวร์ด้วยวันที่..." : "กรุณาเลือกทัวร์ก่อน"}
              value={selectedPeriod ?
                (() => {
                  const { range } = formatDateRange(selectedPeriod.start_date, selectedPeriod.end_date);
                  return range;
                })() : searchTerm
              }
              onChange={(e) => {
                if (!selectedPeriod && tourId) {
                  setSearchTerm(e.target.value);
                  setIsOpen(true);
                }
              }}
              onFocus={() => tourId && setIsOpen(true)}
              className="pl-10"
              readOnly={!!selectedPeriod || !tourId}
              disabled={!tourId}
            />
          </div>
          {selectedPeriod && (
            <Button variant="outline" onClick={handleClear}>
              เคลียร์
            </Button>
          )}
        </div>

        {/* Dropdown รายการ periods */}
        {isOpen && !selectedPeriod && tourId && (
          <div className="absolute top-full left-0 right-0 z-10 mt-1 bg-background border rounded-md shadow-lg">
            <ScrollArea className="max-h-80">
              {isLoading ? (
                <div className="flex items-center justify-center p-4">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span className="ml-2 text-sm text-muted-foreground">กำลังโหลด...</span>
                </div>
              ) : periods.length === 0 ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  {searchTerm ? 'ไม่พบรอบทัวร์ที่ค้นหา' : 'ไม่มีรอบทัวร์'}
                </div>
              ) : (
                <div className="p-2">
                  {periods.map((period) => {
                    const { range, days } = formatDateRange(period.start_date, period.end_date);
                    const availableSeats = period.available_seats || 0;
                    const totalSeats = period.seat || 0;
                    const isFullyBooked = availableSeats === 0;

                    return (
                      <div
                        key={period.id}
                        className={`flex items-center gap-3 p-3 hover:bg-muted cursor-pointer rounded ${
                          isFullyBooked && !allowFullyBooked ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                        onClick={() => (allowFullyBooked || !isFullyBooked) && handlePeriodSelect(period)}
                      >
                        <Calendar className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1">
                            <span className="text-sm font-medium">{range}</span>
                            <Badge variant="outline" className="text-xs">
                              {days} วัน
                            </Badge>
                          </div>
                          <div className="flex items-center gap-3 text-xs text-muted-foreground">
                            <div className="flex items-center gap-1">
                              <Users className="h-3 w-3" />
                              <span>
                                ที่นั่ง: {totalSeats - period.booking || 0}/{totalSeats}
                              </span>
                            </div>
                            <div className="flex items-center gap-1">
                              <span>จอง: {period.booking || 0}</span>
                            </div>
                            {isFullyBooked && (
                              <div className="flex items-center gap-1 text-red-500">
                                <AlertCircle className="h-3 w-3" />
                                <span>เต็ม</span>
                              </div>
                            )}
                          </div>
                          
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </ScrollArea>
          </div>
        )}
      </div>
    </div>
  );
}