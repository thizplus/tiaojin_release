import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useSales } from "@/hooks/useSales";
import type { CreateSaleRequest } from "@/types";

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
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { UserSelect } from "@/components/ui/user-select";

const createSaleSchema = z.object({
  code: z.string().min(1, "กรุณาระบุรหัสพนักงานขาย"),
  first_name: z.string().min(1, "กรุณาระบุชื่อ"),
  last_name: z.string().min(1, "กรุณาระบุนามสกุล"),
  nick_name: z.string().min(1, "กรุณาระบุชื่อเล่น"),
  user_id: z.string().optional(),
});

type CreateSaleFormData = z.infer<typeof createSaleSchema>;

interface CreateSaleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function CreateSaleDialog({
  open,
  onOpenChange,
  onSuccess,
}: CreateSaleDialogProps) {
  const { createSale, isCreateLoading, refetchSalesList } = useSales();
  const [error, setError] = useState<string>("");

  const form = useForm<CreateSaleFormData>({
    resolver: zodResolver(createSaleSchema),
    defaultValues: {
      code: "",
      first_name: "",
      last_name: "",
      nick_name: "",
      user_id: "",
    },
  });

  useEffect(() => {
    if (!open) {
      form.reset();
      setError("");
    }
  }, [open, form]);

  const onSubmit = async (data: CreateSaleFormData) => {
    try {
      setError("");
      const payload: CreateSaleRequest = {
        code: data.code,
        first_name: data.first_name,
        last_name: data.last_name,
        nick_name: data.nick_name,
        ...(data.user_id && { user_id: data.user_id })
      };

      await createSale(payload);
      refetchSalesList();
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการสร้างพนักงานขาย");
    }
  };

  const handleClose = () => {
    onOpenChange(false);
    form.reset();
    setError("");
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>เพิ่มพนักงานขายใหม่</DialogTitle>
          <DialogDescription>
            เพิ่มพนักงานขายใหม่เข้าสู่ระบบ กรุณากรอกข้อมูลให้ครบถ้วน
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <FormField
              control={form.control}
              name="code"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>รหัสพนักงานขาย *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="เช่น SALES001"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="first_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ชื่อ *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="กรุณาระบุชื่อ"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="last_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>นามสกุล *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="กรุณาระบุนามสกุล"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="nick_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ชื่อเล่น *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="กรุณาระบุชื่อเล่น"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="user_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>เชื่อมโยงกับผู้ใช้ (ถ้ามี)</FormLabel>
                  <FormControl>
                    <UserSelect
                      value={field.value}
                      onValueChange={field.onChange}
                      placeholder="ค้นหาและเลือกผู้ใช้..."
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose}>
                ยกเลิก
              </Button>
              <Button type="submit" disabled={isCreateLoading}>
                {isCreateLoading ? "กำลังสร้าง..." : "สร้างพนักงานขาย"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}