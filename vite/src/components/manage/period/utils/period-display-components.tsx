import { Badge } from '@/components/ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { Flame } from 'lucide-react';
import type { TourPeriod } from '@/types/models/period';

// Component สำหรับแสดงคอมมิชชั่นทั้งหมดใน Tooltip
export const AllCommissionDisplay = ({ period }: { period: TourPeriod }) => {
  if (!period.commission_items || period.commission_items.length === 0) {
    return <span className="text-muted-foreground text-sm">-</span>;
  }

  // เลือกคอมมิชชั่นตาม commission_type
  const getCurrentCommission = (commissionItem: (typeof period.commission_items)[0]) => {
    return period.commission_type === 'sale' ? commissionItem.sale : commissionItem.regular;
  };

  const getCommissionTypeLabel = () => {
    return period.commission_type === 'sale' ? 'คอมพิเศษ' : 'คอมปกติ';
  };

  const primaryCommission = period.commission_items[0];

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div className="space-y-1 cursor-pointer hover:bg-muted/50 p-1 rounded">
            <div className="text-sm font-medium">
              ฿{getCurrentCommission(primaryCommission).toLocaleString()}
            </div>
            {period.commission_type === 'sale' ? (
              <div className="text-xs text-primary flex items-center space-x-1">
                <Flame className="h-3 w-3" />
                <span>{getCommissionTypeLabel()}</span>
              </div>
            ) : (
              <div className="text-xs text-foreground">
                {getCommissionTypeLabel()}
              </div>
            )}
            {period.commission_items.length > 1 && (
              <div className="text-xs text-foreground">
                +{period.commission_items.length - 1} คอมอื่น
              </div>
            )}
          </div>
        </TooltipTrigger>
        <TooltipContent className="max-w-sm bg-card text-foreground border-2 rounded-md">
          <div className="space-y-2">
            <div className="font-medium text-sm mb-2 flex items-center justify-between">
              <span>คอมมิชชั่นทั้งหมด:</span>
              {period.commission_type === 'sale' ? (
                <Badge variant="outline" className="text-xs text-primary border-primary flex items-center space-x-1">
                  <Flame className="h-3 w-3" />
                  <span>{getCommissionTypeLabel()}</span>
                </Badge>
              ) : (
                <Badge variant="outline" className="text-xs">
                  {getCommissionTypeLabel()}
                </Badge>
              )}
            </div>
            {period.commission_items.map((commission, index) => (
              <div key={index} className="flex justify-between items-center p-3 border rounded-md">
                <div className="font-medium">คอมมิชชั่น {index + 1}</div>
                <div className="text-right">
                  <div className="text-sm font-bold">
                    ฿{getCurrentCommission(commission).toLocaleString()}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

// Component สำหรับแสดงราคาทั้งหมดใน Tooltip
export const AllRatesDisplay = ({ period }: { period: TourPeriod }) => {
  // หาราคาหลัก (ไม่ exclude และ sort_order ต่ำสุด)
  const primaryRate = period.rate_items
    .filter(item => !item.exclude)
    .sort((a, b) => a.sort_order - b.sort_order)[0] || period.rate_items[0];

  if (!primaryRate) return <span className="text-muted-foreground text-sm">-</span>;

  // เลือกราคาตาม rate_type
  const getCurrentPrice = (rateItem: typeof primaryRate) => {
    return period.rate_type === 'sale' ? rateItem.sale : rateItem.regular;
  };

  const getRateTypeLabel = () => {
    return period.rate_type === 'sale' ? 'ราคาพิเศษ' : 'ราคาปกติ';
  };

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div className="space-y-1 cursor-pointer hover:bg-muted/50 p-1 rounded">
            <div className="text-sm font-medium">
              ฿{getCurrentPrice(primaryRate).toLocaleString()}
            </div>
            <div className="text-xs text-muted-foreground">
              {primaryRate.title}
            </div>
            {period.rate_type === 'sale' ? (
              <div className="text-xs text-primary flex items-center space-x-1">
                <Flame className="h-3 w-3" />
                <span>{getRateTypeLabel()}</span>
              </div>
            ) : (
              <div className="text-xs text-foreground">
                {getRateTypeLabel()}
              </div>
            )}
            {period.rate_items.length > 1 && (
              <div className="text-xs text-foreground">
                +{period.rate_items.length - 1} อัตราอื่น
              </div>
            )}
          </div>
        </TooltipTrigger>
        <TooltipContent className="max-w-sm bg-card text-foreground border-2 rounded-md">
          <div className="space-y-2">
            <div className="font-medium text-sm mb-2 flex items-center justify-between">
              <span>ราคาทั้งหมด:</span>
              {period.rate_type === 'sale' ? (
                <Badge variant="outline" className="text-xs text-primary border-primary flex items-center space-x-1">
                  <Flame className="h-3 w-3" />
                  <span>{getRateTypeLabel()}</span>
                </Badge>
              ) : (
                <Badge variant="outline" className="text-xs">
                  {getRateTypeLabel()}
                </Badge>
              )}
            </div>
            {period.rate_items.map((rate) => (
              <div key={rate.id} className="flex justify-between items-start space-x-4 py-1">
                <div className="flex-1">
                  <div className="text-sm font-medium">{rate.title}</div>
                  {rate.exclude && (
                    <Badge variant="outline" className="text-xs mt-1 bg-muted text-muted-foreground">
                      ไม่ได้ส่วนลด
                    </Badge>
                  )}
                </div>
                <div className="text-right">
                  <div className="text-sm font-bold">
                    ฿{getCurrentPrice(rate).toLocaleString()}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};