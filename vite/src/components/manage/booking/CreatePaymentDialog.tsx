import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { toast } from "sonner";
import { format } from "date-fns";

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
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { FileUpload } from "@/components/ui/file-upload";
import { useAdminBooking } from "@/hooks/useAdminBooking";
import type { BookingResponse } from "@/types";

const paymentMethodOptions = [
  { value: "cash", label: "เงินสด" },
  { value: "bank_transfer", label: "โอนเงิน" },
  { value: "credit_card", label: "บัตรเครดิต" },
  { value: "online", label: "ออนไลน์" },
] as const;

const createPaymentSchema = z.object({
  amount: z.number().min(1, "จำนวนเงินต้องมากกว่า 0"),
  payment_method: z.enum(["cash", "bank_transfer", "credit_card", "online"]),
  reference_number: z.string().optional(),
  note: z.string().optional(),
  payment_date: z.string().optional(),
  slip_image: z.string().optional(),
});

type CreatePaymentFormData = z.infer<typeof createPaymentSchema>;

interface CreatePaymentDialogProps {
  booking: BookingResponse | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function CreatePaymentDialog({
  booking,
  open,
  onOpenChange,
  onSuccess,
}: CreatePaymentDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { createPayment } = useAdminBooking();

  const remainingAmount = booking ? (booking.grand_total || booking.sum_net) - booking.sum_payment : 0;

  const form = useForm<CreatePaymentFormData>({
    resolver: zodResolver(createPaymentSchema),
    defaultValues: {
      amount: remainingAmount,
      payment_method: "bank_transfer",
      reference_number: "",
      note: "",
      payment_date: format(new Date(), "yyyy-MM-dd"),
      slip_image: "",
    },
  });

  const onSubmit = async (data: CreatePaymentFormData) => {
    if (!booking) return;

    setIsSubmitting(true);
    try {
      await createPayment(booking.id, data);

      toast.success("เพิ่มการชำระเงินสำเร็จ");
      onSuccess?.();
      onOpenChange(false);
      form.reset();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "เกิดข้อผิดพลาดในการเพิ่มการชำระเงิน"
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
      <DialogContent className="sm:max-w-md max-h-[90vh] flex flex-col">
        <DialogHeader className="flex-shrink-0">
          <DialogTitle>เพิ่มการชำระเงิน</DialogTitle>
          <DialogDescription>
            เพิ่มการชำระเงินสำหรับการจองหมายเลข: {booking.invoice_id}
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto py-4">
          <div className="mb-4 p-4 bg-muted rounded-lg">
            <div className="grid grid-cols-1 gap-2 text-sm">
              <div>
                <span className="font-medium">ลูกค้า:</span> {booking.customer_name}
              </div>
              <div className="grid grid-cols-3 gap-2">
                <div>
                  <span className="font-medium">ยอดรวม:</span>
                  <div className="text-lg font-semibold">
                    {(booking.grand_total || booking.sum_net).toLocaleString()} ฿
                  </div>
                </div>
                <div>
                  <span className="font-medium">ชำระแล้ว:</span>
                  <div className="text-lg font-semibold text-chart-2">
                    {booking.sum_payment.toLocaleString()} ฿
                  </div>
                </div>
                <div>
                  <span className="font-medium">คงเหลือ:</span>
                  <div className="text-lg font-semibold text-chart-4">
                    {remainingAmount.toLocaleString()} ฿
                  </div>
                </div>
              </div>
            </div>
          </div>

          <Form {...form}>
            <form id="payment-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>จำนวนเงิน (บาท)</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.01"
                        min="0.01"
                        max={remainingAmount}
                        {...field}
                        onChange={(e) => field.onChange(parseFloat(e.target.value) || 0)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="payment_method"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>วิธีการชำระเงิน</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="เลือกวิธีการชำระเงิน" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {paymentMethodOptions.map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
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
                name="reference_number"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>หมายเลขอ้างอิง (ไม่บังคับ)</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        placeholder="เลขที่ธุรกรรม, เลขที่ใบเสร็จ หรือหมายเลขอ้างอิงอื่นๆ"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="payment_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>วันที่ชำระเงิน</FormLabel>
                    <FormControl>
                      <Input
                        type="date"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="slip_image"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>สลิปโอนเงิน (ไม่บังคับ)</FormLabel>
                    <FormControl>
                      <FileUpload
                        accept="image/*"
                        maxSize={5 * 1024 * 1024} // 5MB
                        placeholder="เลือกรูปสลิปโอนเงินหรือลากมาวาง"
                        value={field.value}
                        onChange={field.onChange}
                        variant="dropzone"
                        size="sm"
                        showPreview={true}
                        structuredOptions={{
                          category: 'tours',
                          file_type: 'document',
                          entity_id: booking?.id
                        }}
                      />
                    </FormControl>
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
                        placeholder="หมายเหตุเพิ่มเติมสำหรับการชำระเงิน..."
                        rows={3}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

            </form>
          </Form>
        </div>

        <DialogFooter className="flex-shrink-0">
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
            form="payment-form"
            disabled={isSubmitting}
          >
            {isSubmitting ? "กำลังบันทึก..." : "บันทึก"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}