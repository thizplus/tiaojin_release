import React, { useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { AlertCircle } from 'lucide-react';
import type { Tour, Period } from '@/types/tours-regular';
import {
  PeriodDateColumn,
  PeriodAirlineColumn,
  PeriodPricingColumn,
  PeriodSeatsColumn,
  PeriodDepositColumn,
} from '@/components/tours-regular/columns';
import { PublicPeriodDetailSheet } from './PublicPeriodDetailSheet';

interface PublicPeriodsAccordionProps {
  tour: Tour;
}

export function PublicPeriodsAccordion({ tour }: PublicPeriodsAccordionProps) {
  const [selectedPeriod, setSelectedPeriod] = useState<Period | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);

  const handleRowClick = (period: Period, event: React.MouseEvent) => {
    const target = event.target as HTMLElement;
    if (target.closest('button') || target.closest('[data-no-sheet]')) {
      return;
    }

    setSelectedPeriod(period);
    setIsSheetOpen(true);
  };

  if (!tour.available_periods.length) {
    return (
      <div className="p-4 text-center">
        <AlertCircle className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
        <p className="text-muted-foreground">ไม่มีช่วงเวลาทัวร์ที่พร้อมใช้งาน</p>
      </div>
    );
  }

  return (
    <>
      <div className="space-y-4">
        {tour.available_periods.map((period) => (
          <Card
            key={period.id}
            className="p-0 overflow-hidden hover:shadow-md transition-all duration-300 border-2 hover:border-primary/20 cursor-pointer"
            onClick={(e) => handleRowClick(period, e)}
          >
            <CardContent className="p-0">
              {/* ไม่แสดง Commission Column สำหรับ public */}
              <div className="min-[320px]:space-y-2 sm:space-y-0 sm:block md:grid md:grid-cols-[120px_20%_12%_12%_12%_1fr] py-4 gap-2 px-4">
                <PeriodDateColumn
                  startDate={period.start_date}
                  endDate={period.end_date}
                  createdAt={period.created_at}
                  status={period.status}
                  tourTitle={tour.title}
                />

                <PeriodPricingColumn
                  rateItems={period.rate_items}
                  rateType={period.rate_type}
                />

                <PeriodAirlineColumn
                  airline={period.airline || undefined}
                  airport={period.airport || undefined}
                />

                <PeriodSeatsColumn
                  availableSeats={period.available_seats}
                  totalSeats={period.seat}
                  hasWaitingList={period.has_waiting_list}
                />

                <PeriodDepositColumn
                  deposit={period.deposit}
                  depositType={period.deposit_type}
                />

                {/* ปุ่ม Login เพื่อจอง */}
                <div className="flex items-center justify-end">
                  <a
                    href="/login"
                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2"
                    data-no-sheet
                  >
                    เข้าสู่ระบบเพื่อจอง
                  </a>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Period Detail Sheet - ไม่แสดง commission */}
      <PublicPeriodDetailSheet
        isOpen={isSheetOpen}
        onOpenChange={setIsSheetOpen}
        period={selectedPeriod}
        tour={tour}
      />
    </>
  );
}
