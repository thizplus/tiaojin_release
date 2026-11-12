import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Check, X } from 'lucide-react';
import { cn } from '@/lib/utils';

export type BookingStatus = 'pending' | 'processing' | 'confirm' | 'cancel' | 'completed' | 'expired';

interface BookingStatusOption {
  value: BookingStatus | 'all';
  label: string;
  color: string;
}

interface BookingStatusFilterProps {
  selectedStatuses: BookingStatus[];
  onStatusChange: (statuses: BookingStatus[]) => void;
}

const statusOptions: BookingStatusOption[] = [
  { value: 'all', label: 'ทั้งหมด', color: 'bg-slate-500' },
  { value: 'pending', label: 'รอดำเนินการ', color: 'bg-yellow-500' },
  { value: 'processing', label: 'กำลังดำเนินการ', color: 'bg-blue-500' },
  { value: 'confirm', label: 'ยืนยันแล้ว', color: 'bg-green-500' },
  { value: 'completed', label: 'เสร็จสิ้น', color: 'bg-emerald-500' },
  { value: 'cancel', label: 'ยกเลิก', color: 'bg-red-500' },
  { value: 'expired', label: 'หมดอายุ', color: 'bg-gray-500' },
];

export function BookingStatusFilter({ selectedStatuses, onStatusChange }: BookingStatusFilterProps) {
  const isAllSelected = selectedStatuses.length === 0;

  const handleStatusToggle = (status: BookingStatus | 'all') => {
    if (status === 'all') {
      // เลือก "ทั้งหมด" - clear ทุก status
      onStatusChange([]);
      return;
    }

    // Toggle individual status
    if (selectedStatuses.includes(status)) {
      // ถ้ามีอยู่แล้ว ให้ลบออก
      const newStatuses = selectedStatuses.filter(s => s !== status);
      onStatusChange(newStatuses);
    } else {
      // ถ้ายังไม่มี ให้เพิ่มเข้าไป
      onStatusChange([...selectedStatuses, status]);
    }
  };

  const handleClearAll = () => {
    onStatusChange([]);
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium">กรองตามสถานะ</label>
        {!isAllSelected && (
          <Button
            variant="ghost"
            size="sm"
            onClick={handleClearAll}
            className="h-7 text-xs"
          >
            <X className="h-3 w-3 mr-1" />
            ล้างตัวกรอง
          </Button>
        )}
      </div>

      <div className="flex flex-wrap gap-2">
        {statusOptions.map((option) => {
          const isSelected = option.value === 'all'
            ? isAllSelected
            : selectedStatuses.includes(option.value as BookingStatus);

          return (
            <Badge
              key={option.value}
              variant={isSelected ? 'default' : 'outline'}
              className={cn(
                'cursor-pointer transition-all hover:scale-105',
                'px-3 py-1.5 text-sm',
                isSelected && option.value !== 'all' && option.color,
                isSelected && option.value !== 'all' && 'text-white hover:opacity-90'
              )}
              onClick={() => handleStatusToggle(option.value)}
            >
              {isSelected && (
                <Check className="h-3 w-3 mr-1" />
              )}
              {option.label}
            </Badge>
          );
        })}
      </div>

      {/* แสดงจำนวน status ที่เลือก */}
      {!isAllSelected && (
        <div className="text-xs text-muted-foreground">
          เลือก {selectedStatuses.length} สถานะ
        </div>
      )}
    </div>
  );
}