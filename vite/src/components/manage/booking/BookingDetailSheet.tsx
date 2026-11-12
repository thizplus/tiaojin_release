import { useState, useEffect } from 'react';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { getStorageUrl } from '@/lib/utils';

import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Calendar,
  User,
  MapPin,
  DollarSign,
  CreditCard,
  Receipt,
  Users,
  Trash,
  FileText,
  Bed,
  Plus,
  Upload,
  CheckCircle,
  Building,
  Mail,
  UserCheck,
  RotateCcw
} from 'lucide-react';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import { getBookingStatusConfig, getApprovalStatusConfig } from '@/utils/bookingStatus';
import { PaymentSlipVerifyDialog } from './PaymentSlipVerifyDialog';
import { PaymentSlipReverseDialog } from './PaymentSlipReverseDialog';
import { FileDownloadSheet } from './FileDownloadSheet';
import { PaymentDeadlineNotice } from './PaymentDeadlineNotice';
import PdfInvoiceNew from '@/components/globals/pdfInvoiceNew';
import type { BookingResponse, PaymentSlipResponse, PaymentSlipVerifyRequest, PaymentSlipReverseRequest } from '@/types';

interface BookingDetailSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  booking: BookingResponse | null;

  // User role to determine which actions to show
  userRole?: 'admin' | 'agent' | 'sale';

  // Admin actions
  onDelete?: (booking: BookingResponse) => void;
  onCreatePayment?: (booking: BookingResponse) => void;

  // Agent actions
  onUploadPaymentSlip?: (booking: BookingResponse) => void;
  onUploadPassport?: (booking: BookingResponse) => void;
  onCancelBooking?: (booking: BookingResponse) => void;

  // Admin payment slip actions
  onVerifyPaymentSlip?: (paymentSlipId: string, data: PaymentSlipVerifyRequest) => Promise<void>;
  onReversePaymentSlip?: (paymentSlipId: string, data: PaymentSlipReverseRequest) => Promise<void>;

  // Refetch callback
  onRefetch?: () => void;
}


const isBookingExpired = (booking: BookingResponse): boolean => {
  if (!booking.expires_at) return false;
  return new Date(booking.expires_at) < new Date();
};


export function BookingDetailSheet({
  isOpen,
  onOpenChange,
  booking,
  userRole = 'admin',
  onDelete,
  onCreatePayment,
  onUploadPaymentSlip,
  onUploadPassport,
  onCancelBooking,
  onVerifyPaymentSlip,
  onReversePaymentSlip,
  onRefetch
}: BookingDetailSheetProps) {
  const [isVerifyDialogOpen, setIsVerifyDialogOpen] = useState(false);
  const [isReverseDialogOpen, setIsReverseDialogOpen] = useState(false);
  const [selectedPaymentSlip, setSelectedPaymentSlip] = useState<PaymentSlipResponse | null>(null);
  const [isFilesSheetOpen, setIsFilesSheetOpen] = useState(false);
  const [currentBooking, setCurrentBooking] = useState<BookingResponse | null>(booking);

  // Update currentBooking when booking prop changes
  useEffect(() => {
    if (booking) {
      setCurrentBooking(booking);
    }
  }, [booking]);

  // Handle refetch with updated booking data
  const handleRefetch = async () => {
    if (!currentBooking?.id) return;

    try {
      // Call parent refetch first to update the list
      if (onRefetch) {
        await onRefetch();
      }

      // Parent will update booking prop via BookingCardList useEffect
      // The useEffect will automatically update currentBooking
    } catch (error) {
      console.error('Failed to refetch booking:', error);
    }
  };

  if (!currentBooking) return null;

  const statusConfig = getBookingStatusConfig(currentBooking.status);

  const formatDateRange = () => {
    if (!currentBooking.period_details) return { dateRange: 'ไม่มีข้อมูลรอบ', days: 0 };

    const start = new Date(currentBooking.period_details.start_date);
    const end = new Date(currentBooking.period_details.end_date);
    const diffDays = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;

    return {
      dateRange: `${format(start, 'dd MMM', { locale: th })} - ${format(end, 'dd MMM yyyy', { locale: th })}`,
      days: diffDays
    };
  };

  const { dateRange, days } = formatDateRange();

  // Payment slip handlers
  const handleVerifyPaymentSlip = (slip: PaymentSlipResponse) => {
    setSelectedPaymentSlip(slip);
    setIsVerifyDialogOpen(true);
  };

 

  const handleReversePaymentSlip = (slip: PaymentSlipResponse) => {
    setSelectedPaymentSlip(slip);
    setIsReverseDialogOpen(true);
  };

  const handleVerifyDialogClose = () => {
    setIsVerifyDialogOpen(false);
    setSelectedPaymentSlip(null);
  };

  const handleReverseDialogClose = () => {
    setIsReverseDialogOpen(false);
    setSelectedPaymentSlip(null);
  };

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg overflow-y-auto p-6">
        <SheetHeader className="p-0">
          <SheetTitle className="text-left pr-6">
            รายละเอียดการจอง
          </SheetTitle>
          <SheetDescription className="text-left">
            เลขที่ใบจอง: {currentBooking.invoice_id}
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* สถานะ */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                {statusConfig.iconLarge}
              </div>
              <div className="leading-none">
                <span>สถานะการจอง</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {statusConfig.label}
                </div>
                <div className="text-xs text-muted-foreground">
                  {format(new Date(currentBooking.created_at), "d MMM yyyy HH:mm", { locale: th })}
                </div>
              </div>
            </h3>

            {/* แสดงข้อความ deadline สำหรับ pending bookings */}
            {currentBooking.status === 'pending' && currentBooking.expires_at && (
              <div className="mt-2">
                <PaymentDeadlineNotice expiresAt={currentBooking.expires_at} />
              </div>
            )}
          </div>

          {/* ข้อมูลลูกค้า */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <User className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>ข้อมูลลูกค้า</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {currentBooking.customer_name}
                </div>
                <div className="text-xs text-muted-foreground">
                  {currentBooking.customer_phone}
                </div>
              </div>
            </h3>
          </div>

          {/* ข้อมูลผู้ทำรายการจอง */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Building className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>ผู้ทำรายการจอง</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {currentBooking.company_profile?.name || 'ไม่มีข้อมูลบริษัท'}
                </div>
                <div className="text-xs text-muted-foreground flex items-center gap-1">
                  <Mail className="h-3 w-3" />
                  {currentBooking.user_details?.email || 'ไม่มีข้อมูลอีเมล'}
                </div>
              </div>
            </h3>
          </div>

          {/* ข้อมูลเซลล์ */}
          {currentBooking.sales_details && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <UserCheck className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>เซลล์ผู้ขาย</span>
                  <div className="text-lg font-bold text-foreground leading-none">
                    {currentBooking.sales_details.first_name} {currentBooking.sales_details.last_name}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    รหัส: {currentBooking.sales_details.code} | ชื่อเล่น: {currentBooking.sales_details.nick_name}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* ข้อมูลทัวร์ */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <MapPin className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>ข้อมูลทัวร์</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {currentBooking.tour_details?.code || 'ไม่มีข้อมูลทัวร์'}
                </div>
                {currentBooking.tour_details?.title && (
                  <div className="text-xs text-muted-foreground">
                    {currentBooking.tour_details.title}
                  </div>
                )}
              </div>
            </h3>
          </div>

          {/* วันเดินทาง */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Calendar className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>วันเดินทาง</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {dateRange}
                </div>
                <div className="text-xs text-muted-foreground">
                  {days} วัน
                </div>
              </div>
            </h3>
          </div>

          {/* ไฟล์แนับจาก Period */}
          {currentBooking.period_details?.files && currentBooking.period_details.files.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <FileText className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="flex-1 leading-none">
                  <span>เอกสารประกอบ</span>
                  <div className="text-lg font-bold text-foreground leading-none">
                    {currentBooking.period_details.files.length} ไฟล์
                  </div>
                  <div className="flex items-center gap-2 mt-2">
               

                    <Button
      variant="ghost"
      size="sm"
      onClick={() => setIsFilesSheetOpen(true)}
      className={`cursor-pointer bg-muted`}
    >
      <FileText className={`h-4 w-4 mr-1`} />
      ดูไฟล์ทั้งหมด
    </Button>


                    {currentBooking && (
                      <PdfInvoiceNew
                        booking={currentBooking}
                        variant="ghost"
                        size="sm"
                        style="row"
                      />
                    )}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* จำนวนผู้เดินทาง */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Users className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>จำนวนผู้เดินทาง</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  {currentBooking.sum_qty} คน
                </div>
              </div>
            </h3>
          </div>

          {/* ประเภทห้องพัก */}
          {currentBooking.room_types && currentBooking.room_types.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <Bed className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>ประเภทห้องพัก</span>
                  <div className="text-xs text-muted-foreground space-y-1 mt-1">
                    {currentBooking.room_types.map((room, index) => (
                      <div key={index} className="border rounded p-2">
                        <div className="flex justify-between">
                          <span>{room.title}</span>
                          <span className="font-medium">฿{Number(room.price).toLocaleString()}</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                          จำนวน {room.quantity} ห้อง = ฿{Number(room.total).toLocaleString()}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* ประเภทราคา */}
          {currentBooking.rate_types && currentBooking.rate_types.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <Users className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>ประเภทราคา</span>
                  <div className="text-xs text-muted-foreground space-y-1 mt-1">
                    {currentBooking.rate_types.map((rate, index) => (
                      <div key={index} className="border rounded p-2">
                        <div className="flex justify-between">
                          <span>{rate.title}</span>
                          <span className="font-medium">฿{Number(rate.price).toLocaleString()}</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                          จำนวน {rate.quantity} คน = ฿{Number(rate.total).toLocaleString()}
                          {rate.exclude && <span className="ml-2 text-orange-500">(ไม่ได้ส่วนลด)</span>}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* บริการเพิ่มเติม */}
          {currentBooking.additional_services && currentBooking.additional_services.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <Plus className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>บริการเพิ่มเติม</span>
                  <div className="text-xs text-muted-foreground space-y-1 mt-1">
                    {currentBooking.additional_services.map((service, index) => (
                      <div key={index} className="border rounded p-2">
                        <div className="flex justify-between">
                          <span>{service.title}</span>
                          <span className="font-medium">฿{Number(service.price).toLocaleString()}</span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                          จำนวน {service.quantity} {service.unit_name} = ฿{Number(service.total).toLocaleString()}
                          <span className="ml-2">({service.rate_type === 'per_unit' ? 'ต่อหน่วย' : 'เหมาจ่าย'})</span>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* ข้อมูลการเงิน */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <DollarSign className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>สรุปการเงิน</span>
                <div className="text-lg font-bold text-foreground leading-none">
                  ฿{Number(currentBooking.grand_total || currentBooking.sum_net).toLocaleString()}
                </div>
                <div className="text-xs text-muted-foreground space-y-1">
                  <div>ห้องพัก: ฿{Number(currentBooking.sum_room || 0).toLocaleString()}</div>
                  <div>ราคาทัวร์: ฿{Number(currentBooking.sum_rate || 0).toLocaleString()}</div>
                  <div>บริการเพิ่มเติม: ฿{Number(currentBooking.sum_additional_services || 0).toLocaleString()}</div>
                  <div>ส่วนลด: -฿{Number(currentBooking.discount || 0).toLocaleString()}</div>
                  <div className="border-t pt-1 mt-2">
                    <div className="font-medium text-orange-600">
                      {currentBooking.period_details?.deposit_type === 'percent'
                        ? `มัดจำ ${currentBooking.period_details.deposit}% (฿${Number(currentBooking.sum_deposit || 0).toLocaleString()})`
                        : currentBooking.period_details?.deposit_type === 'per_unit'
                        ? `มัดจำ ฿${Number(currentBooking.period_details.deposit || 0).toLocaleString()}/คน (฿${Number(currentBooking.sum_deposit || 0).toLocaleString()})`
                        : `มัดจำ: ฿${Number(currentBooking.sum_deposit || 0).toLocaleString()}`
                      }
                    </div>
                    <div className="text-xs">
                      กำหนดมัดจำ: {currentBooking.date_deposit ? format(new Date(currentBooking.date_deposit), "d MMM yyyy", { locale: th }) : '-'}
                    </div>
                  </div>
                  <div>ชำระแล้ว: ฿{Number(currentBooking.sum_payment).toLocaleString()}</div>
                  <div>คงเหลือ: ฿{Number((currentBooking.grand_total || currentBooking.sum_net) - currentBooking.sum_payment).toLocaleString()}</div>
                </div>
              </div>
            </h3>
          </div>

          {/* หมายเหตุ */}
          {currentBooking.note && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <FileText className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>หมายเหตุ</span>
                  <div className="text-xs text-muted-foreground mt-1">
                    {currentBooking.note}
                  </div>
                </div>
              </h3>
            </div>
          )}

          {/* สลิปการชำระเงิน */}
          {(currentBooking.payment_slips && currentBooking.payment_slips.length > 0) && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <Receipt className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>สลิปการชำระเงิน</span>
                </div>
              </h3>
              <div className="text-xs text-muted-foreground space-y-1 mt-1">
                {currentBooking.payment_slips.map((slip) => (
                  <div key={`${slip.id}-${slip.status}`} className="border rounded p-2 space-y-2">
                    <div className="flex items-start justify-between">
                      <div className="flex items-start space-x-3">
                        {slip.image ? (
                          <img
                            src={getStorageUrl(slip.image)}
                            alt="Payment slip thumbnail"
                            className="w-16 h-16 object-cover rounded border cursor-pointer hover:opacity-80 flex-shrink-0"
                            onClick={() => window.open(getStorageUrl(slip.image), '_blank')}
                          />
                        ) : (
                          <div className="w-16 h-16 bg-muted/50 rounded border flex items-center justify-center flex-shrink-0">
                            <Receipt className="h-6 w-6 text-muted-foreground" />
                          </div>
                        )}
                        <div className="flex-1">
                          <p className="font-medium text-sm text-foreground">สลิป #{slip.slip_id}</p>
                          <p className="text-xs text-muted-foreground">
                            จำนวน: {slip.amount?.toLocaleString()} บาท
                          </p>
                            <p className="text-xs text-muted-foreground">
                            วิธีชำระ: สลิปการโอน
                          </p>
                          <p className="text-xs text-muted-foreground">
                            อัปโหลดเมื่อ: {new Date(slip.created_at).toLocaleDateString('th-TH', {
                              year: 'numeric',
                              month: 'long',
                              day: 'numeric',
                              hour: '2-digit',
                              minute: '2-digit'
                            })}
                          </p>
                          {/* SlipOK Verification Info */}
                          {slip.slip_verified && (
                            <div className="mt-1 p-1.5 rounded bg-green-50 dark:bg-green-950/30 border border-green-200 dark:border-green-800 space-y-0.5">
                              <p className="text-xs font-medium text-green-700 dark:text-green-400">
                                ✅ สลิปผ่านการตรวจสอบ{slip.status === 'confirm' && ' (ยืนยันอัตโนมัติ)'}
                              </p>
                              {slip.payment_date && (
                                <p className="text-xs text-green-600 dark:text-green-500">
                                  วันที่โอน: {new Date(slip.payment_date).toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                                </p>
                              )}
                              {slip.sending_bank && slip.receiving_bank && (
                                <p className="text-xs text-green-600 dark:text-green-500">
                                  {slip.sending_bank} → {slip.receiving_bank}
                                </p>
                              )}
                              {slip.sender_name && (
                                <p className="text-xs text-green-600 dark:text-green-500">ผู้โอน: {slip.sender_name}</p>
                              )}
                              {slip.trans_ref && (
                                <p className="text-xs text-green-600 dark:text-green-500">Ref: {slip.trans_ref}</p>
                              )}
                            </div>
                          )}
                          {slip.slip_verified === false && slip.slip_verify_error && (
                            <p className="text-xs text-destructive mt-1">⚠️ {slip.slip_verify_error}</p>
                          )}
                        </div>
                      </div>
                      <div className="flex flex-col items-end space-y-2">
                        <div className="flex items-center space-x-2">
                          {(() => {
                            const statusConfig = getApprovalStatusConfig(slip.status);
                            return (
                              <>
                                {statusConfig.icon}
                                <Badge className={`${statusConfig.bgColor} ${statusConfig.textColor}`}>
                                  {statusConfig.label}
                                </Badge>
                              </>
                            );
                          })()}
                        </div>

                        {/* Admin payment slip actions */}
                        {userRole === 'admin' && (onVerifyPaymentSlip || onReversePaymentSlip) && (
                          <div className="flex gap-1">
                            {slip.status === 'pending' && onVerifyPaymentSlip && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleVerifyPaymentSlip(slip)}
                                className="h-7 px-2 text-xs"
                              >
                                <CheckCircle className="h-3 w-3 mr-1" />
                                ตรวจสอบสลิป
                              </Button>
                            )}
                            {( slip.status === 'confirm') && onReversePaymentSlip && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleReversePaymentSlip(slip)}
                                className="h-7 px-2 text-xs text-destructive border-destructive hover:bg-destructive hover:text-destructive-foreground"
                              >
                                <RotateCcw className="h-3 w-3 mr-1" />
                                ยกเลิก
                              </Button>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* หนังสือเดินทาง */}
          {currentBooking.passports && currentBooking.passports.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium flex items-center gap-2">
                <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                  <User className="h-10 w-10 text-muted-foreground" />
                </div>
                <div className="leading-none">
                  <span>หนังสือเดินทาง</span>
                </div>
              </h3>
              <div className="text-xs text-muted-foreground space-y-1 mt-1">
                {currentBooking.passports.map((passport, index) => (
                  <div key={passport.id} className="border rounded p-2 space-y-2">
                    <div className="flex items-start space-x-3">
                      {passport.image ? (
                        <img
                          src={getStorageUrl(passport.image)}
                          alt="Passport thumbnail"
                          className="w-16 h-16 object-cover rounded border cursor-pointer hover:opacity-80 flex-shrink-0"
                          onClick={() => window.open(getStorageUrl(passport.image), '_blank')}
                        />
                      ) : (
                        <div className="w-16 h-16 bg-muted/50 rounded border flex items-center justify-center flex-shrink-0">
                          <User className="h-6 w-6 text-muted-foreground" />
                        </div>
                      )}
                      <div className="flex-1">
                        <p className="font-medium text-sm text-foreground">
                          {passport.title_th && passport.name_th && passport.surname_th
                            ? `${passport.title_th}${passport.name_th} ${passport.surname_th}`
                            : passport.title_en && passport.name_en && passport.surname_en
                            ? `${passport.title_en} ${passport.name_en} ${passport.surname_en}`
                            : `หนังสือเดินทาง #${index + 1}`
                          }
                        </p>
                        <p className="text-xs text-muted-foreground">
                          เลขที่: {passport.passport_number || 'ไม่ระบุ'}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          อัปโหลดเมื่อ: {new Date(passport.created_at).toLocaleDateString('th-TH', {
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* การดำเนินการ */}
          <div className="pt-4 border-t space-y-3">
            {userRole === 'admin' ? (
              // Admin Actions
              <>
                {/* Processing: payment + passport + cancel */}
                {currentBooking.status === "processing" && (
                  <>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onCreatePayment?.(currentBooking)}
                    >
                      <CreditCard className="h-8 w-8 " />
                      เพิ่มการชำระเงิน
                    </Button>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPassport?.(currentBooking)}
                    >
                      <FileText className="h-4 w-4" />
                      แนบหนังสือเดินทาง
                    </Button>
                    <Button
                      variant="destructive"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onDelete?.(currentBooking)}
                    >
                      <Trash className="h-4 w-4" />
                      ยกเลิกการจอง
                    </Button>
                  </>
                )}

                {/* Completed: passport + cancel */}
                {currentBooking.status === "completed" && (
                  <>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPassport?.(currentBooking)}
                    >
                      <FileText className="h-4 w-4" />
                      แนบหนังสือเดินทาง
                    </Button>
                    <Button
                      variant="destructive"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onDelete?.(currentBooking)}
                    >
                      <Trash className="h-4 w-4" />
                      ยกเลิกการจอง
                    </Button>
                  </>
                )}

                {/* Confirm: payment + passport + cancel (ไม่เช็ค expired) */}
                {currentBooking.status === "confirm" && (
                  <>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onCreatePayment?.(currentBooking)}
                    >
                      <CreditCard className="h-4 w-4" />
                      เพิ่มการชำระเงิน
                    </Button>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPassport?.(currentBooking)}
                    >
                      <FileText className="h-4 w-4" />
                      แนบหนังสือเดินทาง
                    </Button>
                    <Button
                      variant="destructive"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onDelete?.(currentBooking)}
                    >
                      <Trash className="h-4 w-4" />
                      ยกเลิกการจอง
                    </Button>
                  </>
                )}

                {/* Pending: payment + passport + cancel (เช็ค expired) */}
                {currentBooking.status === "pending" && !isBookingExpired(currentBooking) && (
                  <>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onCreatePayment?.(currentBooking)}
                    >
                      <CreditCard className="h-4 w-4" />
                      เพิ่มการชำระเงิน
                    </Button>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPassport?.(currentBooking)}
                    >
                      <FileText className="h-4 w-4" />
                      แนบหนังสือเดินทาง
                    </Button>
                    <Button
                      variant="destructive"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onDelete?.(currentBooking)}
                    >
                      <Trash className="h-4 w-4" />
                      ยกเลิกการจอง
                    </Button>
                  </>
                )}
              </>
            ) : userRole === 'agent' ? (
              // Agent Actions
              <>
                {/* Agent can upload payment slip and passport for active bookings */}
                {(currentBooking.status === "pending" || currentBooking.status === "processing" || currentBooking.status === "confirm") && (
                  <>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPaymentSlip?.(currentBooking)}
                    >
                      <Upload className="h-4 w-4" />
                      แจ้งชำระเงิน
                    </Button>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPassport?.(currentBooking)}
                    >
                      <FileText className="h-4 w-4" />
                      แนบหนังสือเดินทาง
                    </Button>
                  </>
                )}

                {/* Agent can cancel their own booking */}
                {(currentBooking.status === "pending" || currentBooking.status === "processing") && (
                  <Button
                    variant="destructive"
                    className="w-full text-lg font-bold cursor-pointer"
                    onClick={() => onCancelBooking?.(currentBooking)}
                  >
                    <Trash className="h-4 w-4" />
                    ยกเลิกการจอง
                  </Button>
                )}
              </>
            ) : userRole === 'sale' ? (
              // Sale Actions - สามารถแนบสลิปและ passport ให้ลูกค้าได้
              <>
                {(currentBooking.status === "pending" || currentBooking.status === "processing" || currentBooking.status === "confirm") && (
                  <>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPaymentSlip?.(currentBooking)}
                    >
                      <Upload className="h-4 w-4" />
                      แจ้งชำระเงินแทนลูกค้า
                    </Button>
                    <Button
                      variant="outline"
                      className="w-full text-lg font-bold cursor-pointer"
                      onClick={() => onUploadPassport?.(currentBooking)}
                    >
                      <FileText className="h-4 w-4" />
                      แนบหนังสือเดินทางแทนลูกค้า
                    </Button>
                  </>
                )}
              </>
            ) : null}

            {/* Cancel, Expired: ไม่มีปุ่ม actions */}
          </div>
        </div>
      </SheetContent>

      {/* Payment Slip Verify Dialog */}
      <PaymentSlipVerifyDialog
        isOpen={isVerifyDialogOpen}
        onOpenChange={handleVerifyDialogClose}
        paymentSlip={selectedPaymentSlip}
        onVerify={onVerifyPaymentSlip || (() => Promise.resolve())}
        onSuccess={handleRefetch}
        isLoading={false}
      />

      {/* Payment Slip Reverse Dialog */}
      <PaymentSlipReverseDialog
        isOpen={isReverseDialogOpen}
        onOpenChange={handleReverseDialogClose}
        paymentSlip={selectedPaymentSlip}
        onReverse={onReversePaymentSlip || (() => Promise.resolve())}
        onSuccess={handleRefetch}
        isLoading={false}
      />

      {/* File Download Sheet */}
      <FileDownloadSheet
        isOpen={isFilesSheetOpen}
        onOpenChange={setIsFilesSheetOpen}
        files={currentBooking.period_details?.files || []}
        periodInfo={currentBooking.period_details ? {
          startDate: currentBooking.period_details.start_date,
          endDate: currentBooking.period_details.end_date,
          days: Math.ceil((new Date(currentBooking.period_details.end_date).getTime() - new Date(currentBooking.period_details.start_date).getTime()) / (1000 * 60 * 60 * 24)) + 1
        } : undefined}
        tourInfo={currentBooking.tour_details ? {
          code: currentBooking.tour_details.code,
          title: currentBooking.tour_details.title,
        } : undefined}
        isFullyPaid={userRole === 'agent'
          ? currentBooking.sum_payment >= (currentBooking.grand_total || currentBooking.sum_net)
          : undefined}
      />
    </Sheet>
  );
}