import React, { memo, useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useUser } from "@/hooks/useUser";
import { Loader2 } from "lucide-react";
import type { CreateUserRequest, UserRole, UserStatus } from "@/types";

interface CreateUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export const CreateUserDialog = memo(function CreateUserDialog({
  open,
  onOpenChange,
  onSuccess
}: CreateUserDialogProps) {
  const [formData, setFormData] = useState<CreateUserRequest>({
    email: "",
    password: "",
    role: "sales",
    status: "pending"
  });

  const { createUser, isCreateLoading, refetchUsersList } = useUser();

  // Reset form when dialog closes
  useEffect(() => {
    if (!open) {
      setFormData({
        email: "",
        password: "",
        role: "sales",
        status: "pending"
      });
    }
  }, [open]);

  // Handle form input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Handle role select changes
  const handleRoleChange = (value: string) => {
    setFormData(prev => ({
      ...prev,
      role: value as UserRole
    }));
  };

  // Handle status select changes
  const handleStatusChange = (value: string) => {
    setFormData(prev => ({
      ...prev,
      status: value as UserStatus
    }));
  };

  // Form validation
  const isFormValid = formData.email.trim() !== "" && 
                     formData.password.trim() !== "" && 
                     formData.password.length >= 6;

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isFormValid) {
      return;
    }

    try {
      await createUser(formData);

      // Revalidate users list to show new data
      refetchUsersList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error("Failed to create user:", error);
    }
  };

  // Close dialog on ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === "Escape" && !isCreateLoading) {
        onOpenChange(false);
      }
    };
    window.addEventListener("keydown", handleEsc);
    return () => window.removeEventListener("keydown", handleEsc);
  }, [onOpenChange, isCreateLoading]);

  return (
    <Dialog open={open} onOpenChange={isCreateLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>สร้างผู้ใช้ใหม่</DialogTitle>
          <DialogDescription>
            กรอกข้อมูลเพื่อสร้างผู้ใช้ใหม่ในระบบ
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="email">อีเมล <span className="text-red-500">*</span></Label>
              <Input
                id="email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleInputChange}
                placeholder="user@example.com"
                required
                autoFocus
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="password">รหัสผ่าน <span className="text-red-500">*</span></Label>
              <Input
                id="password"
                name="password"
                type="password"
                value={formData.password}
                onChange={handleInputChange}
                placeholder="รหัสผ่านอย่างน้อย 6 ตัวอักษร"
                required
                minLength={6}
                disabled={isCreateLoading}
              />
              <p className="text-xs text-muted-foreground">
                รหัสผ่านต้องมีความยาวอย่างน้อย 6 ตัวอักษร
              </p>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="role">บทบาท</Label>
              <Select
                value={formData.role}
                onValueChange={handleRoleChange}
                disabled={isCreateLoading}
              >
                <SelectTrigger>
                  <SelectValue placeholder="เลือกบทบาท" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="sales">Sales</SelectItem>
                  <SelectItem value="agent">Agent</SelectItem>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="super">Super Admin</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="status">สถานะ</Label>
              <Select
                value={formData.status}
                onValueChange={handleStatusChange}
                disabled={isCreateLoading}
              >
                <SelectTrigger>
                  <SelectValue placeholder="เลือกสถานะ" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pending">รออนุมัติ</SelectItem>
                  <SelectItem value="active">ใช้งาน</SelectItem>
                  <SelectItem value="inactive">ระงับ</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isCreateLoading}
            >
              ยกเลิก
            </Button>
            <Button
              type="submit"
              disabled={!isFormValid || isCreateLoading}
            >
              {isCreateLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              สร้าง
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
});
