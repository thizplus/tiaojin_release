import { useState } from "react";
import { useTours } from "@/hooks/useTours";
import type { Tour } from "@/types";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface DeleteTourDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  tour: Tour | null;
  onSuccess?: () => void;
}

export function DeleteTourDialog({
  open,
  onOpenChange,
  tour,
  onSuccess,
}: DeleteTourDialogProps) {
  const { deleteTour, isDeleteLoading, refetchToursList } = useTours();
  const [error, setError] = useState<string>("");

  const handleDelete = async () => {
    if (!tour) return;

    try {
      setError("");
      await deleteTour(tour.id);
      refetchToursList();
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาดในการลบทัวร์");
    }
  };

  const handleClose = () => {
    onOpenChange(false);
    setError("");
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>ยืนยันการลบทัวร์</DialogTitle>
          <DialogDescription>
            คุณต้องการลบทัวร์นี้ออกจากระบบใช่หรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้
          </DialogDescription>
        </DialogHeader>

        {tour && (
          <div className="py-4">
            <div className="space-y-2">
              <p className="text-sm">
                <span className="font-medium">รหัสทัวร์:</span> {tour.code}
              </p>
              <p className="text-sm">
                <span className="font-medium">ชื่อทัวร์:</span> {tour.title}
              </p>
              <p className="text-sm">
                <span className="font-medium">จังหวัด:</span> {tour.city_title}
              </p>
            </div>
          </div>
        )}

        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <DialogFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            ยกเลิก
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={isDeleteLoading}
          >
            {isDeleteLoading ? "กำลังลบ..." : "ลบทัวร์"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}