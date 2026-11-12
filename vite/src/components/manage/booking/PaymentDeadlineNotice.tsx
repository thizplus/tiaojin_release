import { useState, useEffect } from 'react';
import { Clock, AlertTriangle } from 'lucide-react';

interface PaymentDeadlineNoticeProps {
  expiresAt: string | null;
}

export function PaymentDeadlineNotice({ expiresAt }: PaymentDeadlineNoticeProps) {
  const [timeLeft, setTimeLeft] = useState<string>('');

  useEffect(() => {
    if (!expiresAt) return;

    const updateTime = () => {
      const now = new Date();
      const expiry = new Date(expiresAt);

      if (expiry <= now) {
        setTimeLeft('หมดเวลาแล้ว');
        return;
      }

      const difference = expiry.getTime() - now.getTime();
      const days = Math.floor(difference / (1000 * 60 * 60 * 24));
      const hours = Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
      const minutes = Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((difference % (1000 * 60)) / 1000);

      setTimeLeft(`${days}วัน ${hours}ชั่วโมง ${minutes}นาที ${seconds}วินาที`);
    };

    updateTime(); // Initial call
    const interval = setInterval(updateTime, 1000);

    return () => clearInterval(interval);
  }, [expiresAt]);

  if (!expiresAt || !timeLeft) return null;

  const isExpired = timeLeft === 'หมดเวลาแล้ว';

  return (
    <div className={`flex items-center gap-2 text-sm ${isExpired ? 'text-destructive' : 'text-chart-4'}`}>
      {isExpired ? <AlertTriangle className="h-4 w-4" /> : <Clock className="h-4 w-4" />}
      <span>
        {isExpired ? 'หมดเวลาชำระเงินแล้ว' : `กรุณาชำระเงินภายใน ${timeLeft}`}
      </span>
    </div>
  );
}