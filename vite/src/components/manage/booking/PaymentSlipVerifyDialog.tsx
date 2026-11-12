import  { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Badge } from '@/components/ui/badge';
import { CheckCircle, XCircle, Receipt } from 'lucide-react';
import { getApprovalStatusConfig } from '@/utils/bookingStatus';
import type { PaymentSlipResponse, PaymentSlipVerifyRequest } from '@/types';

interface PaymentSlipVerifyDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  paymentSlip: PaymentSlipResponse | null;
  onVerify: (paymentSlipId: string, data: PaymentSlipVerifyRequest) => Promise<void>;
  onSuccess?: () => void;
  isLoading?: boolean;
}

export function PaymentSlipVerifyDialog({
  isOpen,
  onOpenChange,
  paymentSlip,
  onVerify,
  onSuccess,
  isLoading = false
}: PaymentSlipVerifyDialogProps) {
  const [action, setAction] = useState<'confirm' | 'reject'>('confirm');
  const [note, setNote] = useState('');
  const [actualAmount, setActualAmount] = useState('');
  const [adjustmentReason, setAdjustmentReason] = useState('');

  const handleSubmit = async () => {
    if (!paymentSlip) return;

    const data: PaymentSlipVerifyRequest = {
      status: action,
      note: note || undefined,
    };

    if (action === 'confirm' && actualAmount) {
      data.actual_amount = parseFloat(actualAmount);
      data.adjustment_reason = adjustmentReason || undefined;
    }

    try {
      await onVerify(paymentSlip.id, data);
      await onSuccess?.();
      handleClose();
    } catch (error) {
      console.error('Verify payment slip error:', error);
    }
  };

  const handleClose = () => {
    setAction('confirm');
    setNote('');
    setActualAmount('');
    setAdjustmentReason('');
    onOpenChange(false);
  };

  if (!paymentSlip) return null;

  const statusConfig = getApprovalStatusConfig(paymentSlip.status);

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Receipt className="h-5 w-5" />
            ยืนยันสลิปการชำระเงิน
          </DialogTitle>
          <DialogDescription>
            สลิป #{paymentSlip.slip_id} - จำนวน {paymentSlip.amount?.toLocaleString()} บาท
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
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

          {/* Action Selection */}
          <div className="space-y-3">
            <Label>เลือกการดำเนินการ</Label>
            <div className="flex gap-2">
              <Button
                variant={action === 'confirm' ? 'default' : 'outline'}
                onClick={() => setAction('confirm')}
                className={`flex-1 text-lg ${action === 'confirm' ? 'bg-chart-2 hover:bg-chart-2/90 text-chart-2-foreground' : 'text-chart-2 border-chart-2 hover:bg-chart-2/10'}`}
              >
                <CheckCircle className="h-4 w-4" />
                อนุมัติ
              </Button>
              <Button
                variant={action === 'reject' ? 'destructive' : 'outline'}
                onClick={() => setAction('reject')}
                className="flex-1 text-lg"
              >
                <XCircle className="h-4 w-4" />
                ปฏิเสธ
              </Button>
            </div>
          </div>

          {/* Amount Adjustment (only for confirm) */}
          {action === 'confirm' && (
            <div className="space-y-3 p-3 bg-chart-2/10 rounded-lg border border-chart-2/20">
              <Label className="text-chart-2">ปรับยอดเงิน (ถ้ามี)</Label>
              <div className="space-y-2">
                <div>
                  <Label htmlFor="actualAmount" className="text-xs">จำนวนเงินจริง (บาท)</Label>
                  <Input
                    id="actualAmount"
                    type="number"
                    placeholder={paymentSlip.amount?.toString()}
                    value={actualAmount}
                    onChange={(e) => setActualAmount(e.target.value)}
                    step="0.01"
                  />
                </div>
                {actualAmount && parseFloat(actualAmount) !== paymentSlip.amount && (
                  <div>
                    <Label htmlFor="adjustmentReason" className="text-xs">เหตุผลการปรับยอด</Label>
                    <Input
                      id="adjustmentReason"
                      placeholder="ระบุเหตุผลการปรับยอด"
                      value={adjustmentReason}
                      onChange={(e) => setAdjustmentReason(e.target.value)}
                    />
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Note */}
          <div className="space-y-2">
            <Label htmlFor="note">หมายเหตุ (ไม่บังคับ)</Label>
            <Textarea
              id="note"
              placeholder="เพิ่มหมายเหตุเกี่ยวกับการดำเนินการ..."
              value={note}
              onChange={(e) => setNote(e.target.value)}
              rows={3}
            />
          </div>
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
            onClick={handleSubmit}
            disabled={isLoading}
            variant={action === 'reject' ? 'destructive' : 'default'}
            className={action === 'reject' ? 'bg-destructive hover:bg-destructive/90 text-destructive-foreground' : 'bg-chart-2 hover:bg-chart-2/90 text-chart-2-foreground'}
          >
            {isLoading ? 'กำลังดำเนินการ...' : action === 'confirm' ? 'อนุมัติสลิป' : 'ปฏิเสธสลิป'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}