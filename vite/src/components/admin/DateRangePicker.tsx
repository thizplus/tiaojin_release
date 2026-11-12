import * as React from "react";
import { CalendarIcon, X } from "lucide-react";
import { type DateRange } from "react-day-picker";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface DateRangePickerProps {
  dateRange: DateRange | undefined;
  onDateRangeChange: (range: DateRange | undefined) => void;
}

export function DateRangePicker({ dateRange, onDateRangeChange }: DateRangePickerProps) {
  const [isOpen, setIsOpen] = React.useState(false);

  const handleClear = () => {
    onDateRangeChange(undefined);
  };

  const formatDateRange = (range: DateRange | undefined) => {
    if (!range?.from) return "เลือกช่วงวันที่";

    if (!range.to) {
      return format(range.from, "d MMM yyyy", { locale: th });
    }

    return `${format(range.from, "d MMM yyyy", { locale: th })} - ${format(range.to, "d MMM yyyy", { locale: th })}`;
  };

  return (
    <div className="space-y-2">
      <label className="text-sm font-medium">เลือกช่วงวันที่</label>
      <div className="flex gap-2">
        <Popover open={isOpen} onOpenChange={setIsOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              className="flex-1 justify-start text-left font-normal"
            >
              <CalendarIcon className="mr-2 h-4 w-4" />
              {formatDateRange(dateRange)}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto overflow-hidden p-0" align="start">
            <Calendar
              mode="range"
              selected={dateRange}
              onSelect={(range) => {
                onDateRangeChange(range);
                // ปิด popover เมื่อเลือกวันที่ครบ
                if (range?.from && range?.to) {
                  setIsOpen(false);
                }
              }}
              captionLayout="dropdown"
              numberOfMonths={2}
              startMonth={new Date(2020, 0)}
              endMonth={new Date(2035, 11)}
            />
          </PopoverContent>
        </Popover>
        {dateRange?.from && (
          <Button variant="outline" onClick={handleClear} className="shrink-0">
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  );
}