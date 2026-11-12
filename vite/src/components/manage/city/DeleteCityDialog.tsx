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
import { useCity } from "@/hooks/useCity";
import { Loader2 } from "lucide-react";
import type { City } from "@/types";

interface DeleteCityDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  city: City | null;
  onSuccess?: () => void;
}

export const DeleteCityDialog = memo(function DeleteCityDialog({
  open,
  onOpenChange,
  city,
  onSuccess
}: DeleteCityDialogProps) {
  const { deleteCity, isDeleteLoading, refetchCitiesList } = useCity();

  // Handle delete action
  const handleDelete = async () => {
    if (!city) return;

    try {
      await deleteCity(city.id);

      // Revalidate cities list to remove deleted item
      refetchCitiesList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to delete city:', error);
    }
  };

  // Don't render if no city is selected
  if (!city && open) {
    return null;
  }

  return (
    <AlertDialog open={open} onOpenChange={isDeleteLoading ? undefined : onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>ยืนยันการลบเมือง</AlertDialogTitle>
          <AlertDialogDescription>
            คุณต้องการลบเมือง "{city?.title}" ใช่หรือไม่?
            <br />
            <br />
            <span className="font-semibold text-red-500">
              การกระทำนี้ไม่สามารถย้อนกลับได้ และจะลบข้อมูลเมืองนี้ออกจากระบบอย่างถาวร
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
              ลบเมือง
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
});