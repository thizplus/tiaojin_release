import React, { useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { AlertCircle } from 'lucide-react';
import type { Tour, Period } from '@/types/tours-regular';
import {
  PeriodDateColumn,
  PeriodAirlineColumn,
  PeriodPricingColumn,
  PeriodSeatsColumn,
  PeriodCommissionColumn,
  PeriodDepositColumn,
  PeriodActionColumn
} from './columns';
import { PeriodDetailSheet } from './PeriodDetailSheet';
import { BookingSheet } from './BookingSheet';


interface PeriodsAccordionProps {
  tour: Tour;
  onBookPeriod?: (period: Period, tour: Tour) => void;
  canBook?: boolean;
  userStatus?: string;
  userRole?: string; // 'admin' | 'super' | 'agent' | etc.
}

export function PeriodsAccordion({
  tour,
  onBookPeriod,
  canBook = false,
  userStatus = 'pending',
  userRole
}: PeriodsAccordionProps) {
  const [selectedPeriod, setSelectedPeriod] = useState<Period | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [bookingPeriod, setBookingPeriod] = useState<Period | null>(null);
  const [isBookingSheetOpen, setIsBookingSheetOpen] = useState(false);

  const handleRowClick = (period: Period, event: React.MouseEvent) => {
    // ไม่เปิด sheet ถ้าคลิกปุ่ม booking หรือ element ที่มี data-no-sheet
    const target = event.target as HTMLElement;
    if (target.closest('button') || target.closest('[data-no-sheet]')) {
      return;
    }

    setSelectedPeriod(period);
    setIsSheetOpen(true);
  };

  const handleBookPeriod = (period: Period) => {
    setBookingPeriod(period);
    setIsBookingSheetOpen(true);
  };

  const handleBookingSuccess = () => {
    // Optionally call the parent callback
    if (onBookPeriod && bookingPeriod) {
      onBookPeriod(bookingPeriod, tour);
    }

    // Reset booking state
    setBookingPeriod(null);
    setIsBookingSheetOpen(false);
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
              <div className="min-[320px]:space-y-2 sm:space-y-0 sm:block md:grid md:grid-cols-[120px_14%_14%_10%_15%_10%_1fr]  py-4 gap-2 px-4">
              <PeriodDateColumn
                startDate={period.start_date}
                endDate={period.end_date}
                createdAt={period.created_at}
                status={period.status}
                tourTitle={tour.title}
                expirationOverride={period.expiration_override}
                expirationDisplay={period.expiration_display}
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

              <PeriodCommissionColumn
                commissionItems={period.commission_items}
                commissionType={period.commission_type}
              />

              <PeriodDepositColumn
                deposit={period.deposit}
                depositType={period.deposit_type}
              />

              <PeriodActionColumn
                period={period}
                tour={tour}
                userStatus={userStatus}
                canBook={canBook}
                onBookPeriod={(period) => handleBookPeriod(period)}
              />
            </div>
          </CardContent>
        </Card>
      ))}
      </div>

      {/* Period Detail Sheet */}
      <PeriodDetailSheet
        isOpen={isSheetOpen}
        onOpenChange={setIsSheetOpen}
        period={selectedPeriod}
        tour={tour}
        onBookPeriod={(period) => handleBookPeriod(period)}
        canBook={canBook}
        userStatus={userStatus}
      />

      {/* Booking Sheet */}
      <BookingSheet
        isOpen={isBookingSheetOpen}
        onOpenChange={setIsBookingSheetOpen}
        period={bookingPeriod}
        tour={tour}
        onBookingSuccess={handleBookingSuccess}
        userRole={userRole}
      />
    </>
  );
}