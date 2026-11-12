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
import { useSales } from "@/hooks/useSales";
import { Loader2 } from "lucide-react";
import type { Sale } from "@/types";

interface DeleteSaleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sale: Sale | null;
  onSuccess?: () => void;
}

export const DeleteSaleDialog = memo(function DeleteSaleDialog({
  open,
  onOpenChange,
  sale,
  onSuccess
}: DeleteSaleDialogProps) {
  const { deleteSale, isDeleteLoading, refetchSalesList } = useSales();

  const handleDelete = async () => {
    if (!sale) return;

    try {
      await deleteSale(sale.id);
      refetchSalesList();
      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error("Failed to delete sale:", error);
    }
  };

  if (!sale && open) {
    return null;
  }

  return (
    <AlertDialog open={open} onOpenChange={isDeleteLoading ? undefined : onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>ยืนยันการลบพนักงานขาย</AlertDialogTitle>
          <AlertDialogDescription>
            คุณต้องการลบพนักงานขาย "{sale?.first_name} {sale?.last_name}"
            รหัส "{sale?.code}" ใช่หรือไม่?
            <br />
            <br />
            <span className="font-semibold text-red-500">
              การกระทำนี้ไม่สามารถย้อนกลับได้ และจะลบข้อมูลพนักงานขายนี้ออกจากระบบอย่างถาวร
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
              ลบพนักงานขาย
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
});