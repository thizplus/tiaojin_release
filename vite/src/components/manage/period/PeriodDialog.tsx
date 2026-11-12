import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useToast } from '@/hooks/use-toast';
import { usePeriod } from '@/hooks/usePeriod';
import { PeriodForm } from './PeriodForm';
import type { TourPeriod } from '@/types/models/period';
import type { CreatePeriodRequest, UpdatePeriodRequest } from '@/types/request/period';

interface PeriodDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  tourId: string;
  period?: TourPeriod | null;
  mode: 'create' | 'edit';
  onSuccess?: () => void;
}

export function PeriodDialog({
  open,
  onOpenChange,
  tourId,
  period,
  mode,
  onSuccess,
}: PeriodDialogProps) {
  const [loading, setLoading] = useState(false);
  const { toast } = useToast();
  const { handleCreatePeriod, handleUpdatePeriod } = usePeriod();

  const handleSubmit = async (data: CreatePeriodRequest | UpdatePeriodRequest) => {
    setLoading(true);
    try {
      let result;

      if (mode === 'create') {
        result = await handleCreatePeriod(data as CreatePeriodRequest);
      } else if (period?.id) {
        result = await handleUpdatePeriod(period.id, data as UpdatePeriodRequest);
      }

      if (result?.success) {
        toast({
          title: 'สำเร็จ',
          description: mode === 'create'
            ? 'สร้างช่วงเวลาทัวร์เรียบร้อยแล้ว'
            : 'อัปเดตช่วงเวลาทัวร์เรียบร้อยแล้ว',
        });
        onSuccess?.();
        onOpenChange(false);
      } else {
        toast({
          variant: 'destructive',
          title: 'เกิดข้อผิดพลาด',
          description: result?.error || 'ไม่สามารถดำเนินการได้',
        });
      }
    } catch  {
      toast({
        variant: 'destructive',
        title: 'เกิดข้อผิดพลาด',
        description: 'ไม่สามารถดำเนินการได้ กรุณาลองอีกครั้ง',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-6xl max-h-[90vh] overflow-hidden">
        <DialogHeader>
          <DialogTitle>
            {mode === 'create' ? 'เพิ่มช่วงเวลาทัวร์' : 'แก้ไขช่วงเวลาทัวร์'}
          </DialogTitle>
        </DialogHeader>

        <ScrollArea className="h-[calc(90vh-8rem)] pr-4">
          <PeriodForm
            tourId={tourId}
            period={period}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
            loading={loading}
          />
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
}