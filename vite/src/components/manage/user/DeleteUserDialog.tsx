import { memo } from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { useUser } from "@/hooks/useUser";
import { Loader2 } from "lucide-react";
import type { User } from "@/types";

interface DeleteUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  user: User | null;
  onSuccess?: () => void;
}

export const DeleteUserDialog = memo(function DeleteUserDialog({
  open,
  onOpenChange,
  user,
  onSuccess
}: DeleteUserDialogProps) {
  const { deleteUser, isDeleteLoading, refetchUsersList } = useUser();

  const handleDelete = async () => {
    if (!user) return;

    try {
      await deleteUser(user.id);
      refetchUsersList();
      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error("Failed to delete user:", error);
    }
  };

  if (!user && open) {
    return null;
  }

  return (
    <AlertDialog open={open} onOpenChange={isDeleteLoading ? undefined : onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>ยืนยันการลบผู้ใช้</AlertDialogTitle>
          <AlertDialogDescription>
            คุณต้องการลบผู้ใช้ "{user?.email}" ใช่หรือไม่?
            <br />
            <br />
            <span className="font-semibold text-red-500">
              การกระทำนี้ไม่สามารถย้อนกลับได้ และจะลบข้อมูลผู้ใช้นี้ออกจากระบบอย่างถาวร
            </span>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel asChild>
            <Button variant="outline" disabled={isDeleteLoading}>
              ยกเลิก
            </Button>
          </AlertDialogCancel>
          <AlertDialogAction asChild>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={isDeleteLoading}
            >
              {isDeleteLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              ลบผู้ใช้
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
});
