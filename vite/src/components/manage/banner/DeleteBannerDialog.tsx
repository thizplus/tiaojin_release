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
import { useBanner } from "@/hooks/useBanner";
import { Loader2 } from "lucide-react";
import type { Banner } from "@/types";

interface DeleteBannerDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  banner: Banner | null;
  onSuccess?: () => void;
}

export const DeleteBannerDialog = memo(function DeleteBannerDialog({
  open,
  onOpenChange,
  banner,
  onSuccess
}: DeleteBannerDialogProps) {
  const { deleteBanner, isDeleteLoading, refetchBannersList } = useBanner();

  // Handle delete action
  const handleDelete = async () => {
    if (!banner) return;

    try {
      await deleteBanner(banner.id);

      // Revalidate banners list to remove deleted item
      refetchBannersList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to delete banner:', error);
    }
  };

  // Don't render if no banner is selected
  if (!banner && open) {
    return null;
  }

  return (
    <AlertDialog open={open} onOpenChange={isDeleteLoading ? undefined : onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>ยืนยันการลบแบนเนอร์</AlertDialogTitle>
          <AlertDialogDescription>
            คุณต้องการลบแบนเนอร์ "{banner?.title}" ใช่หรือไม่?
            <br />
            <br />
            <span className="font-semibold text-red-500">
              การกระทำนี้ไม่สามารถย้อนกลับได้ และจะลบข้อมูลแบนเนอร์นี้ออกจากระบบอย่างถาวร
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
              ลบแบนเนอร์
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
});