import { useState } from "react";
import { toast } from "sonner";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useAdminBooking } from "@/hooks/useAdminBooking";
import type { BookingResponse } from "@/types";

interface DeleteBookingDialogProps {
  booking: BookingResponse | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function DeleteBookingDialog({
  booking,
  open,
  onOpenChange,
  onSuccess,
}: DeleteBookingDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { deleteBooking } = useAdminBooking();

  const onSubmit = async () => {
    if (!booking) return;

    setIsSubmitting(true);
    try {
      await deleteBooking(booking.id);

      toast.success("ยกเลิกการจองสำเร็จ");
      onSuccess?.();
      onOpenChange(false);
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "เกิดข้อผิดพลาดในการยกเลิกการจอง"
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!booking) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>ยกเลิกการจอง</DialogTitle>
          <DialogDescription>
            คุณแน่ใจหรือไม่ที่ต้องการยกเลิกการจองนี้? การดำเนินการนี้ไม่สามารถยกเลิกได้
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <div className="p-4 bg-muted rounded-lg">
            <div className="space-y-2 text-sm">
              <div>
                <span className="font-medium">หมายเลขใบจอง:</span> {booking.invoice_id}
              </div>
              <div>
                <span className="font-medium">ลูกค้า:</span> {booking.customer_name}
              </div>
              <div>
                <span className="font-medium">เบอร์โทร:</span> {booking.customer_phone}
              </div>
              <div>
                <span className="font-medium">ยอดรวม:</span> {booking.sum_net.toLocaleString()} ฿
              </div>
              <div>
                <span className="font-medium">สถานะ:</span>
                <Badge variant="outline" className="ml-2">
                  {booking.status}
                </Badge>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isSubmitting}
          >
            ยกเลิก
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={onSubmit}
            disabled={isSubmitting}
          >
            {isSubmitting ? "กำลังยกเลิก..." : "ยืนยันการยกเลิก"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}