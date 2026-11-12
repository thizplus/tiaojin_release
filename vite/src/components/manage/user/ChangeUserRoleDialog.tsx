import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useUser } from "@/hooks/useUser";
import type { User, UserRole } from "@/types";

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

const changeRoleSchema = z.object({
  role: z.enum(["super", "admin", "agent", "sales"]),
});

type ChangeRoleFormData = z.infer<typeof changeRoleSchema>;

interface ChangeUserRoleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  user: User | null;
  onSuccess?: () => void;
}

const roleLabels: Record<UserRole, string> = {
  super: "Super Admin",
  admin: "Admin",
  agent: "Agent",
  sales: "Sales",
};

const roleVariants: Record<UserRole, "default" | "secondary" | "destructive" | "outline"> = {
  super: "destructive",
  admin: "default",
  agent: "secondary",
  sales: "outline",
};

export function ChangeUserRoleDialog({
  open,
  onOpenChange,
  user,
  onSuccess,
}: ChangeUserRoleDialogProps) {
  const { updateUserRole, isUpdateRoleLoading, refetchUsersList } = useUser();
  const [error, setError] = useState<string>("");

  const form = useForm<ChangeRoleFormData>({
    resolver: zodResolver(changeRoleSchema),
    defaultValues: {
      role: "sales",
    },
  });

  useEffect(() => {
    if (user && open) {
      form.reset({
        role: user.role,
      });
      setError("");
    }
  }, [user, open, form]);

  const onSubmit = async (data: ChangeRoleFormData) => {
    if (!user) return;

    try {
      setError("");
      await updateUserRole(user.id, data);
      refetchUsersList();
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการเปลี่ยนบทบาทผู้ใช้");
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
          <DialogTitle>เปลี่ยนบทบาทผู้ใช้</DialogTitle>
          <DialogDescription>
            เปลี่ยนบทบาทและสิทธิ์การเข้าถึงของผู้ใช้ในระบบ
          </DialogDescription>
        </DialogHeader>

        {user && (
          <div className="space-y-3 py-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">ผู้ใช้:</span>
              <span className="text-sm">{user.email}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">บทบาทปัจจุบัน:</span>
              <Badge variant={roleVariants[user.role]}>
                {roleLabels[user.role]}
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
              name="role"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>บทบาทใหม่</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="เลือกบทบาทใหม่" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {Object.entries(roleLabels).map(([value, label]) => (
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

            <div className="text-sm text-muted-foreground p-3 bg-muted rounded-md">
              <div className="font-medium mb-2">คำอธิบายบทบาท:</div>
              <ul className="space-y-1 text-xs">
                <li><strong>Super Admin:</strong> สิทธิ์เต็มทุกฟังก์ชัน</li>
                <li><strong>Admin:</strong> จัดการผู้ใช้และข้อมูลระบบ</li>
                <li><strong>Agent:</strong> ดูแลลูกค้าและรายการ</li>
                <li><strong>Sales:</strong> ดูข้อมูลและทำรายการขาย</li>
              </ul>
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose}>
                ยกเลิก
              </Button>
              <Button type="submit" disabled={isUpdateRoleLoading}>
                {isUpdateRoleLoading ? "กำลังเปลี่ยน..." : "เปลี่ยนบทบาท"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}