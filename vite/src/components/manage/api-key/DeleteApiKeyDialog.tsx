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
import { useApiKey } from "@/hooks/useApiKey";
import { Loader2 } from "lucide-react";
import type { ApiKey } from "@/types/models/apiKey";

interface DeleteApiKeyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  apiKey: ApiKey | null;
  onSuccess?: () => void;
}

export const DeleteApiKeyDialog = memo(function DeleteApiKeyDialog({
  open,
  onOpenChange,
  apiKey,
  onSuccess
}: DeleteApiKeyDialogProps) {
  const { deleteApiKey, isDeleteLoading, refetchApiKeysList } = useApiKey();

  // Handle delete action
  const handleDelete = async () => {
    if (!apiKey) return;

    try {
      await deleteApiKey(apiKey.id);

      // Revalidate API keys list to remove deleted item
      refetchApiKeysList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to delete API key:', error);
    }
  };

  // Don't render if no API key is selected
  if (!apiKey && open) {
    return null;
  }

  return (
    <AlertDialog open={open} onOpenChange={isDeleteLoading ? undefined : onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>ยืนยันการลบ API Key</AlertDialogTitle>
          <AlertDialogDescription>
            คุณต้องการลบ API Key "{apiKey?.name}" ใช่หรือไม่?
            <br />
            <br />
            <span className="font-semibold text-red-500">
              การกระทำนี้ไม่สามารถย้อนกลับได้ และจะทำให้ API Key นี้ไม่สามารถใช้งานได้อีกต่อไป
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
              ลบ API Key
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
});