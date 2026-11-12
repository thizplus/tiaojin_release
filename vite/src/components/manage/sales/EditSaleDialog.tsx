import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useSales } from "@/hooks/useSales";
import type { Sale, UpdateSaleRequest } from "@/types";

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
import { Switch } from "@/components/ui/switch";
import { UserSelect } from "@/components/ui/user-select";

const editSaleSchema = z.object({
  first_name: z.string().min(1, "กรุณาระบุชื่อ"),
  last_name: z.string().min(1, "กรุณาระบุนามสกุล"),
  nick_name: z.string().min(1, "กรุณาระบุชื่อเล่น"),
  user_id: z.string().optional(),
  active: z.boolean(),
});

type EditSaleFormData = z.infer<typeof editSaleSchema>;

interface EditSaleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sale: Sale | null;
  onSuccess?: () => void;
}

export function EditSaleDialog({
  open,
  onOpenChange,
  sale,
  onSuccess,
}: EditSaleDialogProps) {
  const { updateSale, isUpdateLoading, refetchSalesList } = useSales();
  const [error, setError] = useState<string>("");

  const form = useForm<EditSaleFormData>({
    resolver: zodResolver(editSaleSchema),
    defaultValues: {
      first_name: "",
      last_name: "",
      nick_name: "",
      user_id: "",
      active: true,
    },
  });

  useEffect(() => {
    if (sale && open) {
      form.reset({
        first_name: sale.first_name,
        last_name: sale.last_name,
        nick_name: sale.nick_name,
        user_id: sale.user_id || "",
        active: sale.active,
      });
      setError("");
    }
  }, [sale, open, form]);

  const onSubmit = async (data: EditSaleFormData) => {
    if (!sale) return;

    try {
      setError("");
      const payload: UpdateSaleRequest = {
        first_name: data.first_name,
        last_name: data.last_name,
        nick_name: data.nick_name,
        active: data.active,
        ...(data.user_id && { user_id: data.user_id })
      };

      await updateSale(sale.id, payload);
      refetchSalesList();
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการแก้ไขพนักงานขาย");
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
          <DialogTitle>แก้ไขข้อมูลพนักงานขาย</DialogTitle>
          <DialogDescription>
            แก้ไขข้อมูลพนักงานขายในระบบ กรุณากรอกข้อมูลให้ครบถ้วน
          </DialogDescription>
        </DialogHeader>

        {sale && (
          <div className="py-2">
            <p className="text-sm text-muted-foreground">
              รหัส: <span className="font-mono">{sale.code}</span>
            </p>
          </div>
        )}

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

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

            <FormField
              control={form.control}
              name="active"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">
                      สถานะการใช้งาน
                    </FormLabel>
                    <div className="text-sm text-muted-foreground">
                      เปิดหรือปิดการใช้งานพนักงานขาย
                    </div>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose}>
                ยกเลิก
              </Button>
              <Button type="submit" disabled={isUpdateLoading}>
                {isUpdateLoading ? "กำลังบันทึก..." : "บันทึกการแก้ไข"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}