import  { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Badge } from '@/components/ui/badge';
import { RotateCcw, AlertTriangle } from 'lucide-react';
import { getApprovalStatusConfig } from '@/utils/bookingStatus';
import { getStorageUrl } from '@/lib/utils';
import type { PaymentSlipResponse, PaymentSlipReverseRequest } from '@/types';

interface PaymentSlipReverseDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  paymentSlip: PaymentSlipResponse | null;
  onReverse: (paymentSlipId: string, data: PaymentSlipReverseRequest) => Promise<void>;
  onSuccess?: () => void;
  isLoading?: boolean;
}

export function PaymentSlipReverseDialog({
  isOpen,
  onOpenChange,
  paymentSlip,
  onReverse,
  onSuccess,
  isLoading = false
}: PaymentSlipReverseDialogProps) {
  const [reversalReason, setReversalReason] = useState('');

  const handleSubmit = async () => {
    if (!paymentSlip || reversalReason.length < 10) return;

    try {
      await onReverse(paymentSlip.id, { reversal_reason: reversalReason });
      await onSuccess?.();
      handleClose();
    } catch (error) {
      console.error('Reverse payment slip error:', error);
    }
  };

  const handleClose = () => {
    setReversalReason('');
    onOpenChange(false);
  };

  if (!paymentSlip) return null;

  const statusConfig = getApprovalStatusConfig(paymentSlip.status);
  const isValidReason = reversalReason.length >= 10 && reversalReason.length <= 500;

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-destructive">
            <RotateCcw className="h-5 w-5" />
            ยกเลิกสลิปการชำระเงิน
          </DialogTitle>
          <DialogDescription>
            สลิป #{paymentSlip.slip_id} - จำนวน {paymentSlip.amount?.toLocaleString()} บาท
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Warning */}
          <div className="flex items-start gap-3 p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
            <AlertTriangle className="h-5 w-5 text-destructive flex-shrink-0 mt-0.5" />
            <div className="text-sm">
              <p className="font-medium text-destructive">คำเตือน</p>
              <p className="text-muted-foreground mt-1">
                การยกเลิกสลิปจะทำให้ยอดชำระเงินในการจองลดลง และอาจเปลี่ยนสถานะการจองกลับเป็น "รอดำเนินการ"
              </p>
            </div>
          </div>

          {/* Current Status */}
          <div className="flex items-center justify-between p-3 bg-muted/50 rounded-lg">
            <span className="text-sm font-medium">สถานะปัจจุบัน:</span>
            <div className="flex items-center space-x-2">
              {statusConfig.icon}
              <Badge className={`${statusConfig.bgColor} ${statusConfig.textColor}`}>
                {statusConfig.label}
              </Badge>
            </div>
          </div>

          {/* Reversal Reason */}
          <div className="space-y-2">
            <Label htmlFor="reversalReason">
              เหตุผลการยกเลิก <span className="text-destructive">*</span>
            </Label>
            <Textarea
              id="reversalReason"
              placeholder="ระบุเหตุผลการยกเลิกสลิป (อย่างน้อย 10 ตัวอักษร)"
              value={reversalReason}
              onChange={(e) => setReversalReason(e.target.value)}
              rows={4}
              className={!isValidReason && reversalReason.length > 0 ? 'border-destructive' : ''}
            />
            <div className="flex justify-between text-xs text-muted-foreground">
              <span>
                {reversalReason.length < 10 && reversalReason.length > 0 && (
                  <span className="text-destructive">ต้องมีอย่างน้อย 10 ตัวอักษร</span>
                )}
                {reversalReason.length > 500 && (
                  <span className="text-destructive">ไม่เกิน 500 ตัวอักษร</span>
                )}
              </span>
              <span className={reversalReason.length > 500 ? 'text-destructive' : ''}>
                {reversalReason.length}/500
              </span>
            </div>
          </div>

          {/* Payment Slip Preview */}
          {paymentSlip.image && (
            <div className="space-y-2">
              <Label>ตัวอย่างสลิป</Label>
              <div className="border rounded-lg p-2">
                <img
                  src={getStorageUrl(paymentSlip.image)}
                  alt="Payment slip"
                  className="w-full h-32 object-cover rounded cursor-pointer hover:opacity-80"
                  onClick={() => window.open(getStorageUrl(paymentSlip.image), '_blank')}
                />
              </div>
            </div>
          )}
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={handleClose}
            disabled={isLoading}
          >
            ยกเลิก
          </Button>
          <Button
            variant="destructive"
            onClick={handleSubmit}
            disabled={!isValidReason || isLoading}
          >
            {isLoading ? 'กำลังยกเลิก...' : 'ยกเลิกสลิป'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}