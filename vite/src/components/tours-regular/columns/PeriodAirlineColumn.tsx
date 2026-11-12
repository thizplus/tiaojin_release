import { Plane } from 'lucide-react';

interface PeriodAirlineColumnProps {
  airline?: string;
  airport?: string;
}

export function PeriodAirlineColumn({ airline, airport }: PeriodAirlineColumnProps) {
  return (
    <div className="flex items-center space-x-3">
      <div className="flex-shrink-0 p-2 bg-muted/50 rounded-full">
        <Plane className="h-5 w-5 text-foreground" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="font-bold text-lg text-foreground leading-tight">
          {airline || 'ไม่ระบุสายการบิน'}
        </div>
        <div className="text-xs text-muted-foreground truncate">
          {airport || 'สนามบิน'}
        </div>
      </div>
    </div>
  );
}