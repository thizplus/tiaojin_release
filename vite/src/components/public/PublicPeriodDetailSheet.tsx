import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Calendar,
  Plane,
  Users,
  DollarSign,
  Banknote,
  MapPin,
  Flame,
  Plus,
  LogIn
} from 'lucide-react';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import type { Tour, Period } from '@/types/tours-regular';

interface PublicPeriodDetailSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  period: Period | null;
  tour: Tour | null;
}

export function PublicPeriodDetailSheet({
  isOpen,
  onOpenChange,
  period,
  tour,
}: PublicPeriodDetailSheetProps) {
  if (!period || !tour) return null;

  const getCurrentPrice = (rateItem: Period['rate_items'][0], rateType: string) => {
    return rateType === 'sale' ? rateItem.sale : rateItem.regular;
  };

  const formatDeposit = (deposit: number, depositType: string) => {
    if (depositType === 'percent') {
      return `${deposit}%`;
    }
    return `฿${deposit.toLocaleString()}`;
  };

  const formatDateRange = () => {
    const start = new Date(period.start_date);
    const end = new Date(period.end_date);
    const diffDays = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;
    return {
      dateRange: `${format(start, 'dd MMM', { locale: th })} - ${format(end, 'dd MMM yyyy', { locale: th })}`,
      days: diffDays
    };
  };

  const { dateRange, days } = formatDateRange();

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg overflow-y-auto p-6">
        <SheetHeader className="p-0">
          <SheetTitle className="text-left pr-6">
            {tour.title}
          </SheetTitle>
          <SheetDescription className="text-left">
            รหัสทัวร์: {tour.code}
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* วันเดินทาง */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Calendar className="h-10 w-10 text-muted-foreground" />
              </div>

              <div className="leading-none">
                <span>วันเดินทาง</span>
                <div className="text-lg font-bold text-foreground leading-none">{dateRange}</div>
                <div className="text-xs text-muted-foreground">{days} วัน</div>
              </div>
            </h3>
          </div>

          {/* สายการบิน */}
          {(period.airline || period.airport) && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <Plane className="h-10 w-10 text-muted-foreground" />
                </div>

                <div className="leading-none">
                  <span>สายการบิน</span>
                  <div className="text-lg font-bold text-foreground leading-none">{period.airline || 'ไม่ระบุ'}</div>
                  {period.airport && (
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <MapPin className="h-3 w-3" />
                      {period.airport}
                    </div>
                  )}
                </div>
              </h3>
            </div>
          )}

          {/* ที่นั่ง */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Users className="h-10 w-10 text-muted-foreground" />
              </div>

              <div className="leading-none">
                <span>ที่นั่ง</span>
                <div className="text-lg font-bold text-foreground leading-none">ว่าง {period.available_seats} ที่นั่ง</div>
                <div className="text-xs text-muted-foreground space-y-1">
                  <div>ทั้งหมด: {period.seat} ที่นั่ง</div>
                </div>
              </div>
            </h3>
          </div>

          {/* ราคา */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <DollarSign className="h-10 w-10 text-muted-foreground" />
              </div>

              <div className="leading-none">
                <span>ราคา</span>
                <div className="space-y-2">
                  {period.rate_items.map((rate) => (
                    <div key={rate.id}>
                      <div className="text-xs text-muted-foreground flex items-center gap-2">
                        {rate.title}
                        {rate.exclude && (
                          <Badge variant="secondary" className="text-xs h-4">
                            ไม่ได้ส่วนลด
                          </Badge>
                        )}
                      </div>
                      <div className="text-lg font-bold text-foreground leading-none">
                        ฿{getCurrentPrice(rate, period.rate_type).toLocaleString()}
                      </div>
                      {period.rate_type === 'sale' && rate.regular && rate.regular !== rate.sale && (
                        <div className="text-xs text-muted-foreground line-through">
                          ราคาปกติ ฿{rate.regular.toLocaleString()}
                        </div>
                      )}
                    </div>
                  ))}
                  {period.rate_type === 'sale' && (
                    <Badge variant="destructive" className="text-xs mt-1 h-5">
                      <Flame className="h-3 w-3" />
                      พิเศษ
                    </Badge>
                  )}
                </div>
              </div>
            </h3>
          </div>

          {/* ไม่แสดง Commission สำหรับ Public */}

          {/* มัดจำ */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Banknote className="h-10 w-10 text-muted-foreground" />
              </div>

              <div className="leading-none">
                <span>มัดจำ</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {formatDeposit(period.deposit, period.deposit_type)}
                </div>
                <div className="text-xs text-muted-foreground">
                  {period.deposit_type === 'percent' ? 'มัดจำ' : 'มัดจำ / ที่นั่ง'}
                </div>
              </div>
            </h3>
          </div>

          {/* บริการเพิ่มเติม */}
          {period.additional_services && period.additional_services.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <Plus className="h-10 w-10 text-muted-foreground" />
                </div>

                <div className="leading-none">
                  <span>บริการเพิ่มเติม</span>
                  <div className="space-y-2">
                    {period.additional_services.map((service, index) => (
                      <div key={service.id || index}>
                        <div className="text-xs text-muted-foreground">{service.title || `บริการ ${index + 1}`}</div>
                        <div className="text-lg font-bold text-foreground leading-none">
                          ฿{service.price?.toLocaleString() || '0'}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {service.rate_type === 'per_unit' ? `ต่อ / ${service.unit_name || 'หน่วย'}` : service.unit_name || 'เหมาจ่าย'}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* ปุ่ม Login เพื่อจอง */}
          <div className="pt-4 border-t">
            <Button
              size="lg"
              className="w-full text-lg font-bold"
              asChild
            >
              <a href="/login">
                <LogIn className="h-5 w-5 mr-2" />
                เข้าสู่ระบบเพื่อจองทัวร์
              </a>
            </Button>
            <p className="text-center text-sm text-muted-foreground mt-3">
              ยังไม่มีบัญชี?{' '}
              <a href="/register" className="text-primary hover:underline">
                ลงทะเบียนที่นี่
              </a>
            </p>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
