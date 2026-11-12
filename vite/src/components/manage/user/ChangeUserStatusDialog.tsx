import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useUser } from "@/hooks/useUser";
import type { User, UserStatus } from "@/types";

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
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";

const changeStatusSchema = z.object({
  status: z.enum(["pending", "active", "inactive"]),
});

type ChangeStatusFormData = z.infer<typeof changeStatusSchema>;

interface ChangeUserStatusDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  user: User | null;
  onSuccess?: () => void;
}

const statusLabels: Record<UserStatus, string> = {
  pending: "รอการอนุมัติ",
  active: "ใช้งานได้",
  inactive: "ปิดใช้งาน",
};

const statusVariants: Record<UserStatus, "default" | "secondary" | "destructive"> = {
  pending: "secondary",
  active: "default",
  inactive: "destructive",
};

export function ChangeUserStatusDialog({
  open,
  onOpenChange,
  user,
  onSuccess,
}: ChangeUserStatusDialogProps) {
  const { updateUserStatus, isUpdateStatusLoading, refetchUsersList } = useUser();
  const [error, setError] = useState<string>("");

  const form = useForm<ChangeStatusFormData>({
    resolver: zodResolver(changeStatusSchema),
    defaultValues: {
      status: "pending",
    },
  });

  useEffect(() => {
    if (user && open) {
      form.reset({
        status: user.status,
      });
      setError("");
    }
  }, [user, open, form]);

  const onSubmit = async (data: ChangeStatusFormData) => {
    if (!user) return;

    try {
      setError("");
      await updateUserStatus(user.id, data);
      refetchUsersList();
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการเปลี่ยนสถานะผู้ใช้");
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
          <DialogTitle>เปลี่ยนสถานะผู้ใช้</DialogTitle>
          <DialogDescription>
            เปลี่ยนสถานะการใช้งานของผู้ใช้ในระบบ
          </DialogDescription>
        </DialogHeader>

        {user && (
          <div className="space-y-3 py-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">ผู้ใช้:</span>
              <span className="text-sm">{user.email}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">สถานะปัจจุบัน:</span>
              <Badge variant={statusVariants[user.status]}>
                {statusLabels[user.status]}
              </Badge>
            </div>
            {user.company_profile?.name && (
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">บริษัท:</span>
                <span className="text-sm">{user.company_profile.name}</span>
              </div>
            )}
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
              name="status"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>สถานะใหม่</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="เลือกสถานะใหม่" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {Object.entries(statusLabels).map(([value, label]) => (
                        <SelectItem key={value} value={value}>
                          {label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose}>
                ยกเลิก
              </Button>
              <Button type="submit" disabled={isUpdateStatusLoading}>
                {isUpdateStatusLoading ? "กำลังเปลี่ยน..." : "เปลี่ยนสถานะ"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}