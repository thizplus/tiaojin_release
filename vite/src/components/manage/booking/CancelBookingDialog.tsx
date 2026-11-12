
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { AlertTriangle } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useToast } from '@/hooks/use-toast';
import { getBookingStatusConfig } from '@/utils/bookingStatus';
import type { BookingResponse } from '@/types';
import { useState } from 'react';

interface CancelBookingDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  booking: BookingResponse | null;
  onCancel: (id: string, reason?: string) => Promise<void>;
  isLoading?: boolean;
}

export function CancelBookingDialog({
  isOpen,
  onOpenChange,
  booking,
  onCancel,
  isLoading = false
}: CancelBookingDialogProps) {
  const [reason, setReason] = useState<string>('');
  const { toast } = useToast();

  const handleSubmit = async () => {
    if (!booking) {
      toast({
        title: 'เกิดข้อผิดพลาด',
        description: 'ไม่พบข้อมูลการจอง',
        variant: 'destructive',
      });
      return;
    }

    try {
      await onCancel(booking.id, reason || 'ยกเลิกโดยเอเจนต์');
      toast({
        title: 'ยกเลิกสำเร็จ',
        description: 'ยกเลิกการจองเรียบร้อยแล้ว',
      });
      handleClose();
    } catch (error) {
      toast({
        title: 'เกิดข้อผิดพลาด',
        description: error instanceof Error ? error.message : 'ไม่สามารถยกเลิกการจองได้',
        variant: 'destructive',
      });
    }
  };

  const handleClose = () => {
    setReason('');
    onOpenChange(false);
  };


  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-destructive">
            <AlertTriangle className="h-5 w-5" />
            ยืนยันการยกเลิกการจอง
          </DialogTitle>
          <DialogDescription>
            คุณต้องการยกเลิกการจอง {booking?.invoice_id} หรือไม่?
          </DialogDescription>
        </DialogHeader>

        {booking && (
          <div className="space-y-4">
            {/* Booking Info */}
            <div className="bg-muted rounded-lg p-4 space-y-2">
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-medium">{booking.customer_name}</p>
                  <p className="text-sm text-muted-foreground">{booking.customer_phone}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-muted-foreground">สถานะปัจจุบัน</p>
                  {(() => {
                    const statusConfig = getBookingStatusConfig(booking.status);
                    return (
                      <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium ${statusConfig.bgColor} ${statusConfig.textColor}`}>
                        {statusConfig.icon}
                        {statusConfig.label}
                      </span>
                    );
                  })()}
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4 pt-2 border-t border-border">
                <div>
                  <p className="text-sm text-muted-foreground">ทัวร์</p>
                  <p className="font-medium">{booking.tour_details?.title || 'N/A'}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">ยอดเงิน</p>
                  <p className="font-medium">{(booking.grand_total || booking.sum_net)?.toLocaleString()} บาท</p>
                </div>
              </div>
            </div>

            {/* Warning */}
            <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
              <div className="flex">
                <div className="flex-shrink-0">
                  <AlertTriangle className="h-5 w-5 text-destructive" />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-destructive">
                    ข้อมูลสำคัญ
                  </h3>
                  <div className="mt-2 text-sm text-destructive/80">
                    <ul className="list-disc pl-5 space-y-1">
                      <li>การยกเลิกจะไม่สามารถเปลี่ยนแปลงได้อีก</li>
                      <li>ระบบจะปรับสถานะเป็น "ยกเลิก" ทันที</li>
                      <li>หากมีการชำระเงินแล้ว ต้องติดต่อ Admin เพื่อคืนเงิน</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            {/* Reason Input */}
            <div className="space-y-2">
              <Label htmlFor="cancel-reason">เหตุผลในการยกเลิก (ไม่บังคับ)</Label>
              <Textarea
                id="cancel-reason"
                placeholder="ระบุเหตุผลในการยกเลิกการจอง..."
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                rows={3}
              />
            </div>
          </div>
        )}

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={handleClose}
            disabled={isLoading}
          >
            ไม่ยกเลิก
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleSubmit}
            disabled={isLoading}
            className="min-w-[100px]"
          >
            {isLoading ? 'กำลังยกเลิก...' : 'ยืนยันยกเลิก'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}