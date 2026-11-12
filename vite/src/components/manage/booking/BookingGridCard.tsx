import React from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Eye,
  Trash,
  CreditCard,
  User,
  DollarSign,
  Receipt,
  Ellipsis,
  FileText,
} from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import type { BookingResponse } from "@/types";
import PdfInvoiceNew from "@/components/globals/pdfInvoiceNew";
import { getBookingStatusConfig } from '@/utils/bookingStatus';

interface BookingAction {
  key: string;
  label: string;
  icon?: React.ReactNode;
  onClick: (booking: BookingResponse) => void;
  variant?: 'default' | 'destructive';
  showFor?: (booking: BookingResponse) => boolean;
  separator?: boolean;
}

interface BookingGridCardProps {
  booking: BookingResponse;
  onView?: (booking: BookingResponse) => void;
  onSelect?: (booking: BookingResponse) => void;
  onViewFiles?: (booking: BookingResponse) => void;
  /** เช็คชำระเงินครบก่อนดูไฟล์ (agent side) */
  requireFullPaymentForFiles?: boolean;
  actions?: BookingAction[];
  showDefaultActions?: boolean;
  // Legacy props for backward compatibility
  onDelete?: (booking: BookingResponse) => void;
  onCreatePayment?: (booking: BookingResponse) => void;
}


const isBookingExpired = (booking: BookingResponse): boolean => {
  if (!booking.expires_at) return false;
  return new Date(booking.expires_at) < new Date();
};

export function BookingGridCard({
  booking,
  onView,
  onSelect,
  onViewFiles,
  requireFullPaymentForFiles = false,
  actions = [],
  showDefaultActions = true,
  // Legacy props for backward compatibility
  onDelete,
  onCreatePayment
}: BookingGridCardProps) {
  const statusConfig = getBookingStatusConfig(booking.status);
  const isFullyPaid = booking.sum_payment >= (booking.grand_total || booking.sum_net);
  const canViewFiles = !requireFullPaymentForFiles || isFullyPaid;

  const handleCardClick = () => {
    if (onSelect) {
      onSelect(booking);
    }
  };

  // Generate default admin actions for backward compatibility
  const defaultAdminActions: BookingAction[] = [
    {
      key: 'add-payment',
      label: 'เพิ่มการชำระเงิน',
      icon: <CreditCard className="mr-1 h-3.5 w-3.5" />,
      onClick: (booking) => onCreatePayment?.(booking),
      showFor: (booking) => {
        const hasPaymentAction = ['processing', 'confirm'].includes(booking.status) ||
          (booking.status === 'pending' && !isBookingExpired(booking));
        return hasPaymentAction && !!onCreatePayment;
      }
    },
    {
      key: 'cancel-booking',
      label: 'ยกเลิกการจอง',
      icon: <Trash className="mr-1 h-3.5 w-3.5 text-destructive" />,
      onClick: (booking) => onDelete?.(booking),
      variant: 'destructive',
      separator: true,
      showFor: (booking) => {
        const canCancel = ['processing', 'confirm', 'completed'].includes(booking.status) ||
          (booking.status === 'pending' && !isBookingExpired(booking));
        return canCancel && !!onDelete;
      }
    }
  ];

  // Combine custom actions with default actions
  const allActions = showDefaultActions
    ? [...actions, ...defaultAdminActions]
    : actions;

  // Filter actions based on showFor condition
  const visibleActions = allActions.filter(action =>
    !action.showFor || action.showFor(booking)
  );



  const formatDateRange = () => {
    if (!booking.period_details) return { dateRange: 'ไม่มีข้อมูล', days: 0 };

    const start = new Date(booking.period_details.start_date);
    const end = new Date(booking.period_details.end_date);
    const diffDays = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;

    return {
      startDay: format(start, 'd', { locale: th }),
      endDay: format(end, 'd', { locale: th }),
      month: format(start, 'MMM yyyy', { locale: th }),
      days: diffDays
    };
  };

  const { startDay, endDay, month } = formatDateRange();

  return (
    <Card className="p-0 hover:shadow-md transition-shadow cursor-pointer" onClick={handleCardClick}>
      <CardContent className="p-0">
        <div className="min-[320px]:space-y-2 sm:space-y-0 sm:block md:grid md:grid-cols-[200px_25%_15%_10%_15%_10%_1fr] py-2 gap-2 px-2">

          {/* เลขที่ใบจอง - แสดงในเดสก์ท็อป */}
          <div className="hidden md:flex items-center">
            <div className="flex flex-col items-start text-left min-w-[100px] px-1 py-1">
              <div className="font-bold text-lg text-foreground">
                {booking.invoice_id}
              </div>
              <Badge className={`${statusConfig.bgColor} ${statusConfig.textColor} text-xs px-1 py-0 h-5 mt-1`}>
                {statusConfig.label}
              </Badge>
              <div className="text-xs text-muted-foreground mt-1">
                {format(new Date(booking.created_at), "d MMM yyyy", { locale: th })}
              </div>
            </div>
          </div>

          {/* ทัวร์และวันเดินทาง - แสดงในเดสก์ท็อป */}
          <div className="hidden md:flex items-center space-x-3">
            <div className="flex-shrink-0">
              <div className="relative flex flex-col items-center text-center min-w-[80px] px-1 py-1 text-foreground rounded-md">
                <div className="flex items-center space-x-1 text-primary font-bold text-sm">
                  <span className="text-foreground text-lg leading-none">{startDay}</span>
                  <span className="text-foreground text-lg leading-none">-</span>
                  <span className="text-foreground text-lg leading-none">{endDay}</span>
                </div>
                <div className="text-xs text-foreground leading-tight">{month}</div>
              </div>
            </div>
            <div className="min-w-0 flex-1">
              <div className="font-bold text-sm text-foreground leading-tight">
                {booking.tour_details?.code || 'N/A'}
              </div>
              <div className="text-xs text-muted-foreground truncate">
                {booking.tour_details?.title?.substring(0, 30) || 'ไม่มีข้อมูล'}...
              </div>
            </div>
          </div>

          {/* ข้อมูลลูกค้าและทัวร์ - แสดงในมือถือ */}
          <div className="block md:hidden">
            <div className="flex items-start space-x-3 p-2 bg-muted/30 rounded-lg">
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between">
                  <div className="flex-1 min-w-0">
                    <h4 className="text-xs font-medium text-foreground line-clamp-3 leading-tight">
                      {booking.tour_details?.title || 'ไม่มีข้อมูลทัวร์'}
                    </h4>
                    <div className="text-xs text-muted-foreground mt-1">
                      {booking.customer_name} • {booking.customer_phone}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {booking.invoice_id}
                    </div>
                  </div>
                  <div className="flex-shrink-0 ml-2">
                    <div className="relative flex flex-col items-center text-center min-w-[90px] px-4 py-1 text-foreground rounded-md">
                      <div className=" z-10">
                        <Badge className={`${statusConfig.bgColor} ${statusConfig.textColor} text-xs px-1 py-0 h-5`}>
                          {statusConfig.label}
                        </Badge>
                      </div>
                      <div className="flex items-center space-x-1 text-primary font-bold text-lg">
                        <span className="text-foreground text-2xl leading-none">{startDay}</span>
                        <span className="text-foreground text-2xl leading-none">-</span>
                        <span className="text-foreground text-2xl leading-none">{endDay}</span>
                      </div>
                      <div className="text-xs text-foreground leading-tight">{month}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* ปุ่มใบแจ้งหนี้และไฟล์ใน mobile */}
            <div className="flex items-center justify-start mt-2 px-2 gap-2">
              {booking.period_details?.files && booking.period_details.files.length > 0 && onViewFiles && canViewFiles && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={(e) => {
                    e.stopPropagation();
                    onViewFiles(booking);
                  }}
                  className="cursor-pointer"
                >
                  <FileText className="h-4 w-4 mr-1" />
                  ไฟล์
                </Button>
              )}
              <PdfInvoiceNew booking={booking} variant="outline" size="sm" style="row" />
            </div>
          </div>

          {/* ข้อมูลลูกค้า */}
          <div className="hidden md:flex items-center space-x-2">
            <div className="flex-shrink-0 p-1 bg-muted/50 rounded-full">
              <User className="h-4 w-4 text-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="font-bold text-foreground text-sm leading-tight">
                {booking.customer_name}
              </div>
              <div className="text-xs text-muted-foreground truncate">
                {booking.customer_phone}
              </div>
            </div>
          </div>

          {/* ยอดเงิน */}
          <div className="flex items-center space-x-2">
            <div className="flex-shrink-0 p-1 bg-muted/50 rounded-full">
              <DollarSign className="h-4 w-4 text-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="font-bold text-foreground text-lg leading-tight">
                ฿{Number(booking.grand_total || booking.sum_net).toLocaleString()}
              </div>
              <div className="text-xs text-muted-foreground truncate">
                {booking.sum_qty} คน
              </div>
            </div>
          </div>


          {/* การชำระเงิน */}
          <div className="flex items-center space-x-2">
            <div className="flex-shrink-0 p-1 bg-muted/50 rounded-full">
              <Receipt className="h-4 w-4 text-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="font-bold text-foreground text-lg leading-tight">
                ฿{Number(booking.sum_payment).toLocaleString()}
              </div>
              <div className="text-xs text-muted-foreground">
                ค้างชำระ ฿{Number((booking.grand_total || booking.sum_net) - booking.sum_payment).toLocaleString()}
              </div>
            </div>
          </div>

          {/* Files & Invoice */}
          <div className="hidden md:flex items-center justify-center gap-1">
            {booking.period_details?.files && booking.period_details.files.length > 0 && onViewFiles && canViewFiles && (
              <Button
                variant="ghost"
                size="sm"
                onClick={(e) => {
                  e.stopPropagation();
                  onViewFiles(booking);
                }}
                className="cursor-pointer flex-col h-full"
              >
                <FileText className="h-4 w-4 mb-1" />
                ไฟล์
              </Button>
            )}
            <PdfInvoiceNew booking={booking} variant="ghost" size="sm" style="column" />
          </div>

          {/* Actions */}
          <div className="flex items-center justify-center">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={(e) => e.stopPropagation()}
                >
                  <Ellipsis className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={(e) => { e.stopPropagation(); onView?.(booking); }}>
                  <Eye className="mr-1 h-3.5 w-3.5" />
                  ดูรายละเอียด
                </DropdownMenuItem>
                {booking.period_details?.files && booking.period_details.files.length > 0 && onViewFiles && canViewFiles && (
                  <DropdownMenuItem
                    onClick={(e) => {
                      e.stopPropagation();
                      onViewFiles(booking);
                    }}
                  >
                    <FileText className="h-4 w-4 mr-1" />
                    ดูไฟล์
                  </DropdownMenuItem>
                )}
                <DropdownMenuItem onClick={(e) => e.stopPropagation()} className="p-0">
                  <div className="w-full">
                    <PdfInvoiceNew booking={booking} variant="ghost" size="sm" />
                  </div>
                </DropdownMenuItem>

                {/* Dynamic Actions */}
                {visibleActions.map((action) => (
                  <React.Fragment key={action.key}>
                    {action.separator && <DropdownMenuSeparator />}
                    <DropdownMenuItem
                      onClick={(e) => {
                        e.stopPropagation();
                        action.onClick(booking);
                      }}
                      className={action.variant === 'destructive' ? 'text-destructive focus:text-destructive' : ''}
                    >
                      {action.icon}
                      {action.label}
                    </DropdownMenuItem>
                  </React.Fragment>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}