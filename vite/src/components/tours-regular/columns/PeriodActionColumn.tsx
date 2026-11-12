import { Clock, AlertCircle, CheckCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Period, Tour } from '@/types/tours-regular';

interface PeriodActionColumnProps {
  period: Period;
  tour: Tour;
  userStatus?: string;
  canBook?: boolean;
  onBookPeriod?: (period: Period, tour: Tour) => void;
}

export function PeriodActionColumn({
  period,
  tour,
  userStatus = 'pending',
  canBook = false,
  onBookPeriod
}: PeriodActionColumnProps) {
  const handleBooking = () => {
    if (onBookPeriod) {
      onBookPeriod(period, tour);
    }
  };

  return (
    <div className="flex items-center justify-center" data-no-sheet>
      {userStatus !== 'approve' ? (
        <Badge variant="secondary" className="text-xs px-3 py-1 font-medium">
          <Clock className="h-3 w-3 mr-1" />
          รออนุมัติ
        </Badge>
      ) : canBook ? (
        period.available_seats <= 0 ? (
        


<Button
size="default"
disabled
onClick={handleBooking}
className="bg-muted-foreground text-foreground text-lg h-8 px-4 font-bold w-full sm:w-auto"
>
<AlertCircle className="h-5 w-5" />
FULL BOOKING
</Button>

        ) : (
          <Button
            size="default"
            onClick={handleBooking}
            className="text-lg h-8 px-4 font-bold w-full sm:w-auto cursor-pointer"
          >
            <CheckCircle className="h-5 w-5" />
            BOOKING NOW
          </Button>
        )
      ) : (
        <Badge
          variant={period.status === 'publish' ? 'default' : 'secondary'}
          className="text-xs px-3 py-1 font-medium"
        >
          {period.status === 'publish' ? (
            <>
              <CheckCircle className="h-3 w-3 mr-1" />
              พร้อม
            </>
          ) : (
            <>
              <AlertCircle className="h-3 w-3 mr-1" />
              ร่าง
            </>
          )}
        </Badge>
      )}
    </div>
  );
}