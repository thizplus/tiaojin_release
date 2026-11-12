import { Percent, Flame, Info } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import type { Period } from '@/types/tours-regular';

interface PeriodCommissionColumnProps {
  commissionItems: Period['commission_items'];
  commissionType: string;
}

export function PeriodCommissionColumn({ commissionItems, commissionType }: PeriodCommissionColumnProps) {
  const getCurrentCommission = (commissionItem: Period['commission_items'][0], commissionType: string) => {
    return commissionType === 'sale' ? commissionItem.sale : commissionItem.regular;
  };

  const getCommissionDisplay = () => {
    const commissions = commissionItems.map(item =>
      getCurrentCommission(item, commissionType).toLocaleString()
    );
    return commissions.join('+');
  };

  const hasMultipleCommissions = commissionItems.length > 1;

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div className="flex items-center space-x-3 cursor-help">
            <div className="flex-shrink-0 p-2 bg-muted/50 rounded-full">
              <Percent className="h-5 w-5 text-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="font-bold text-foreground text-lg leading-tight flex items-center gap-1">
                ฿{getCommissionDisplay()}
                {hasMultipleCommissions && <Info className="h-3 w-3 text-muted-foreground" />}
              </div>
              <div className="text-xs text-muted-foreground flex items-center">
                คอมมิชชั่น
                {hasMultipleCommissions && ` (+${commissionItems.length - 1} คอมอื่น)`}
               
              </div>
              {commissionType === 'sale' && (
                <Badge variant="destructive" className="text-xs mt-1 h-5">
                  <Flame className="h-3 w-3" />
                  พิเศษ
                </Badge>
              )}
            </div>
          </div>
        </TooltipTrigger>
        {hasMultipleCommissions && (
          <TooltipContent className="max-w-sm bg-card text-foreground border-2 rounded-md">
            <div className="space-y-1">
              <div className="font-medium text-sm mb-2 flex items-center justify-between">
                <span>คอมมิชชั่นทั้งหมด:</span>
                {commissionType === 'sale' ? (
                  <Badge variant="outline" className="text-xs text-primary border-primary flex items-center space-x-1">
                    <Flame className="h-3 w-3" />
                    <span>คอมพิเศษ</span>
                  </Badge>
                ) : (
                  <Badge variant="outline" className="text-xs">
                    คอมปกติ
                  </Badge>
                )}
              </div>
              {commissionItems.map((commission, index) => (
                <div key={index} className="flex justify-between items-center px-3 py-1 border rounded-md">
                  <div className="font-medium">คอมมิชชั่น {index + 1}</div>
                  <div className="text-right">
                    <div className="text-sm font-bold">
                      ฿{getCurrentCommission(commission, commissionType).toLocaleString()}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </TooltipContent>
        )}
      </Tooltip>
    </TooltipProvider>
  );
}