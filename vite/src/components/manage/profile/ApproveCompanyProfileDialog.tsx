import { useState } from "react";
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
import { CheckCircle, AlertTriangle } from "lucide-react";
import { getStorageUrl } from "@/lib/utils";
import type { CompanyProfile } from "@/types";

interface ApproveCompanyProfileDialogProps {
  profile: CompanyProfile | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: (profile: CompanyProfile) => Promise<void>;
  isLoading?: boolean;
}

export function ApproveCompanyProfileDialog({
  profile,
  open,
  onOpenChange,
  onConfirm,
  isLoading = false,
}: ApproveCompanyProfileDialogProps) {
  const [error, setError] = useState<string | null>(null);

  const handleConfirm = async () => {
    if (!profile) return;

    try {
      setError(null);
      await onConfirm(profile);
      onOpenChange(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'เกิดข้อผิดพลาดในการอนุมัติ');
    }
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      setError(null);
    }
    onOpenChange(open);
  };

  if (!profile) return null;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <CheckCircle className="h-5 w-5 text-green-600" />
            ยืนยันการอนุมัติ
          </DialogTitle>
          <DialogDescription>
            คุณต้องการอนุมัติโปรไฟล์บริษัทนี้หรือไม่?
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Company Info */}
          <div className="bg-muted/50 p-4 rounded-lg space-y-2">
            <div className="flex items-center gap-3">
              {profile.logo && (
                <img
                  src={getStorageUrl(profile.logo)}
                  alt="Company Logo"
                  className="h-12 w-12 object-cover rounded-md"
                />
              )}
              <div>
                <p className="font-semibold text-lg">{profile.name}</p>
                <p className="text-sm text-muted-foreground font-mono">
                  {profile.tax_id.replace(/(\d{1})(\d{4})(\d{5})(\d{2})(\d{1})/, "$1-$2-$3-$4-$5")}
                </p>
              </div>
            </div>
          </div>

          {/* Warning Alert */}
          <Alert>
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>
              การอนุมัติจะเปลี่ยนสถานะของผู้ใช้จาก "รอการอนุมัติ" เป็น "ใช้งานได้"
              และไม่สามารถยกเลิกได้ กรุณาตรวจสอบข้อมูลให้ถูกต้องก่อนดำเนินการ
            </AlertDescription>
          </Alert>

          {/* Error Message */}
          {error && (
            <Alert variant="destructive">
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isLoading}
          >
            ยกเลิก
          </Button>
          <Button
            onClick={handleConfirm}
            disabled={isLoading}
            className="bg-green-600 hover:bg-green-700"
          >
            {isLoading ? (
              <>
                <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-b-transparent" />
                กำลังอนุมัติ...
              </>
            ) : (
              <>
                <CheckCircle className="mr-2 h-4 w-4" />
                ยืนยันการอนุมัติ
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}