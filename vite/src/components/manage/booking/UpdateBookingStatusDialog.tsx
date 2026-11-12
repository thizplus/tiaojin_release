import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { toast } from "sonner";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useAdminBooking } from "@/hooks/useAdminBooking";
import type { BookingResponse } from "@/types";

const statusOptions = [
  { value: "pending", label: "รอดำเนินการ", color: "text-yellow-600" },
  { value: "processing", label: "กำลังดำเนินการ", color: "text-blue-600" },
  { value: "confirm", label: "ยืนยันแล้ว", color: "text-green-600" },
  { value: "cancel", label: "ยกเลิก", color: "text-red-600" },
  { value: "completed", label: "เสร็จสิ้น", color: "text-emerald-600" },
] as const;

const updateStatusSchema = z.object({
  status: z.enum(["pending", "processing", "confirm", "cancel", "completed"]),
  note: z.string().optional(),
});

type UpdateStatusFormData = z.infer<typeof updateStatusSchema>;
type BookingStatus = "pending" | "processing" | "confirm" | "cancel" | "completed";

const isValidBookingStatus = (status: string): status is BookingStatus => {
  return ["pending", "processing", "confirm", "cancel", "completed"].includes(status);
};

interface UpdateBookingStatusDialogProps {
  booking: BookingResponse | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function UpdateBookingStatusDialog({
  booking,
  open,
  onOpenChange,
  onSuccess,
}: UpdateBookingStatusDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { updateBookingStatus } = useAdminBooking();

  const form = useForm<UpdateStatusFormData>({
    resolver: zodResolver(updateStatusSchema),
    defaultValues: {
      status: booking?.status && isValidBookingStatus(booking.status) ? booking.status : "pending",
      note: "",
    },
  });

  const onSubmit = async (data: UpdateStatusFormData) => {
    if (!booking) return;

    setIsSubmitting(true);
    try {
      await updateBookingStatus(booking.id, data);

      toast.success("อัปเดตสถานะการจองสำเร็จ");
      onSuccess?.();
      onOpenChange(false);
      form.reset();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "เกิดข้อผิดพลาดในการอัปเดตสถานะ"
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset();
    }
    onOpenChange(newOpen);
  };

  if (!booking) return null;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>เปลี่ยนสถานะการจอง</DialogTitle>
          <DialogDescription>
            เปลี่ยนสถานะการจองสำหรับหมายเลข: {booking.invoice_id}
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <div className="mb-4 p-4 bg-muted rounded-lg">
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div>
                <span className="font-medium">ลูกค้า:</span> {booking.customer_name}
              </div>
              <div>
                <span className="font-medium">สถานะปัจจุบัน:</span>
                <Badge variant="outline" className="ml-2">
                  {statusOptions.find(s => s.value === booking.status)?.label || booking.status}
                </Badge>
              </div>
            </div>
          </div>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="status"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>สถานะใหม่</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="เลือกสถานะ" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {statusOptions.map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            <span className={option.color}>{option.label}</span>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="note"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>หมายเหตุ (ไม่บังคับ)</FormLabel>
                    <FormControl>
                      <Textarea
                        {...field}
                        placeholder="เพิ่มหมายเหตุสำหรับการเปลี่ยนสถานะ..."
                        rows={3}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

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
                  type="submit"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "กำลังบันทึก..." : "บันทึก"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </div>
      </DialogContent>
    </Dialog>
  );
}