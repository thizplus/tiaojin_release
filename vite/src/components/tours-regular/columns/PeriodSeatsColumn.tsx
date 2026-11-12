import { Users } from 'lucide-react';

interface PeriodSeatsColumnProps {
  availableSeats: number;
  totalSeats: number;
  hasWaitingList?: boolean;
}

export function PeriodSeatsColumn({ availableSeats, totalSeats, hasWaitingList }: PeriodSeatsColumnProps) {
  return (
    <div className="flex items-center space-x-3">
      <div className="flex-shrink-0 p-2 bg-muted/50 rounded-full">
        <Users className="h-5 w-5 text-foreground" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="font-bold text-lg text-foreground leading-tight">
          <span className="">ว่าง {availableSeats}</span>
        </div>
        <div className="text-xs text-muted-foreground">
          {totalSeats} ที่นั่ง
        </div>
        {hasWaitingList && (
          <span className="inline-flex items-center rounded-full bg-orange-100 px-2 py-0.5 text-xs font-medium text-orange-800 dark:bg-orange-900/50 dark:text-orange-300 mt-1">
            WL
          </span>
        )}
      </div>
    </div>
  );
}