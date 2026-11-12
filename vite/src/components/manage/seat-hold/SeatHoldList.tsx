import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  RefreshCw,
  AlertCircle,
} from 'lucide-react';
import { usePeriodSeatHolds } from '@/hooks/useSeatHold';
import { SeatHoldCard } from './SeatHoldCard';
import type { SeatHold } from '@/types/models/seat-hold';
import { useState } from 'react';

interface SeatHoldListProps {
  periodId: string;
  seatHolds: SeatHold[];
  isLoading: boolean;
}

export function SeatHoldList({ periodId, seatHolds, isLoading }: SeatHoldListProps) {
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'released' | 'expired'>('all');
  const [includeInactive, setIncludeInactive] = useState(false);

  const { refetch, isRefetching } = usePeriodSeatHolds(
    periodId,
    { include_inactive: includeInactive },
    { enabled: !!periodId }
  );

  // Filter seat holds by status
  const filteredSeatHolds = seatHolds.filter((seatHold) => {
    if (statusFilter === 'all') return true;
    return seatHold.status === statusFilter;
  });

  

  const totalHeldSeats = seatHolds
    .filter(hold => hold.status === 'active')
    .reduce((sum, hold) => sum + hold.held_seats, 0);

  if (isLoading) {
    return (
      
          <div className="flex items-center justify-center space-x-2 text-muted-foreground">
            <RefreshCw className="h-4 w-4 animate-spin" />
            <span>กำลังโหลดข้อมูลการจำนองที่นั่ง...</span>
          </div>
        
    );
  }

  return (
    <div className="space-y-6">
      {/* Simple Summary */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="text-sm text-muted-foreground">
            รวม: <strong>{seatHolds.length}</strong> รายการ
          </div>
          <div className="text-sm text-muted-foreground">
            ที่นั่งจำนอง: <strong>{totalHeldSeats}</strong> ที่นั่ง
          </div>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => refetch()}
          disabled={isRefetching}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${isRefetching ? 'animate-spin' : ''}`} />
          รีเฟรช
        </Button>
      </div>

      {/* Simple Filter */}
      <div className="flex items-center space-x-4">
        <Select value={statusFilter} onValueChange={(value: 'all' | 'active' | 'released' | 'expired') => setStatusFilter(value)}>
          <SelectTrigger className="w-40">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">ทั้งหมด</SelectItem>
            <SelectItem value="active">กำลังจำนอง</SelectItem>
            <SelectItem value="released">ปล่อยแล้ว</SelectItem>
            <SelectItem value="expired">หมดอายุ</SelectItem>
          </SelectContent>
        </Select>

        <div className="flex items-center space-x-2">
          <input
            type="checkbox"
            id="include-inactive"
            checked={includeInactive}
            onChange={(e) => setIncludeInactive(e.target.checked)}
            className="rounded"
          />
          <label htmlFor="include-inactive" className="text-sm text-muted-foreground">
            แสดงรายการที่ไม่ใช้งาน
          </label>
        </div>
      </div>

      {/* Seat Holds List */}
      {filteredSeatHolds.length === 0 ? (
        <div className="text-center py-12">
          <AlertCircle className="mx-auto h-8 w-8 text-muted-foreground" />
          <h3 className="mt-4 text-lg font-medium text-muted-foreground">
            {seatHolds.length === 0 ? 'ยังไม่มีการจำนองที่นั่ง' : 'ไม่มีรายการที่ตรงกับการกรอง'}
          </h3>
          <p className="mt-2 text-sm text-muted-foreground">
            {seatHolds.length === 0
              ? 'เริ่มต้นจำนองที่นั่งโดยใช้แท็บ "จำนองที่นั่งใหม่"'
              : 'ลองเปลี่ยนตัวกรองหรือรีเฟรชข้อมูล'
            }
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredSeatHolds.map((seatHold) => (
            <SeatHoldCard key={seatHold.id} seatHold={seatHold} />
          ))}
        </div>
      )}
    </div>
  );
}