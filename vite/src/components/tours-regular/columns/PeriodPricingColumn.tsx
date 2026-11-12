import { DollarSign, Flame, Info } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import type { Period } from '@/types/tours-regular';

interface PeriodPricingColumnProps {
  rateItems: Period['rate_items'];
  rateType: string;
}

export function PeriodPricingColumn({ rateItems, rateType }: PeriodPricingColumnProps) {
  const getCurrentPrice = (rateItem: Period['rate_items'][0], rateType: string) => {
    return rateType === 'sale' ? rateItem.sale : rateItem.regular;
  };

  const hasMultipleRates = rateItems.length > 1;

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div className="flex items-center space-x-3 cursor-help">
            <div className="flex-shrink-0 p-2 bg-muted/50 rounded-full">
              <DollarSign className="h-5 w-5 text-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              {rateItems.slice(0, 1).map((rate) => (
                <div key={rate.id}>
                  <div className="font-bold text-foreground text-lg leading-tight flex items-center gap-1">
                    ฿{getCurrentPrice(rate, rateType).toLocaleString()}
                    {hasMultipleRates && <Info className="h-3 w-3 text-muted-foreground" />}
                  </div>
                  {rateType === 'sale' && rate.regular && rate.regular !== rate.sale && (
                    <div className="text-xs text-muted-foreground line-through">
                      ราคาปกติ ฿{rate.regular.toLocaleString()}
                    </div>
                  )}
                  <div className="text-xs text-muted-foreground truncate">
                    {rate.title}
                    {hasMultipleRates && ` (+${rateItems.length - 1} อัตราอื่น)`}
                  </div>
                </div>
              ))}
              {rateType === 'sale' && (
                <Badge variant="destructive" className="text-xs mt-1 h-5">
                  <Flame className="h-3 w-3" />
                  พิเศษ
                </Badge>
              )}
            </div>
          </div>
        </TooltipTrigger>
        {hasMultipleRates && (
          <TooltipContent className="max-w-sm bg-card text-foreground border-2 rounded-md">
            <div className="space-y-2">
              <div className="font-medium text-sm mb-2 flex items-center justify-between">
                <span>ราคาทั้งหมด:</span>
                {rateType === 'sale' ? (
                  <Badge variant="outline" className="text-xs text-primary border-primary flex items-center space-x-1">
                    <Flame className="h-3 w-3" />
                    <span>พิเศษ</span>
                  </Badge>
                ) : (
                  <Badge variant="outline" className="text-xs">
                    ราคาปกติ
                  </Badge>
                )}
              </div>
              {rateItems.map((rate) => (
                <div key={rate.id} className="flex justify-between items-start space-x-4 py-1">
                  <div className="flex-1">
                    <div className="text-sm font-medium flex items-center gap-2">
                      {rate.title}
                      {rate.exclude && (
                        <Badge variant="secondary" className="text-xs h-5">
                          ไม่ได้ส่วนลด
                        </Badge>
                      )}
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-bold">
                      ฿{getCurrentPrice(rate, rateType).toLocaleString()}
                    </div>
                    {rateType === 'sale' && rate.regular && rate.regular !== rate.sale && (
                      <div className="text-xs text-muted-foreground line-through">
                        ฿{rate.regular.toLocaleString()}
                      </div>
                    )}
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