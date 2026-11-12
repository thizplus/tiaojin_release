import { Banknote } from 'lucide-react';

interface PeriodDepositColumnProps {
  deposit: number;
  depositType: string;
}

export function PeriodDepositColumn({ deposit, depositType }: PeriodDepositColumnProps) {
  const formatDeposit = (deposit: number, depositType: string) => {
    if (depositType === 'percent') {
      return `${deposit}%`;
    }
    return `฿${deposit.toLocaleString()}`;
  };

  return (
    <div className="lg:col-span-1 flex items-center space-x-3">
      <div className="flex-shrink-0 p-2 bg-muted/50 rounded-full">
        <Banknote className="h-5 w-5 text-foreground" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="font-bold text-lg text-foreground leading-tight">
          {formatDeposit(deposit, depositType)}
        </div>
        <div className="text-xs text-muted-foreground">
          {depositType === 'percent' ? 'มัดจำ' : 'มัดจำ / ที่นั่ง'}
        </div>
      </div>
    </div>
  );
}