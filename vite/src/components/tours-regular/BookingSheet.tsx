import { useState, useEffect, useRef } from 'react';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { SalesSelect } from '@/components/ui/sales-select';
import { UserSelect } from '@/components/ui/user-select';
import { Switch } from '@/components/ui/switch';
import {
  Calendar,
  Plane,
  Users,
  DollarSign,
  Banknote,
  MapPin,
  Phone,
  User,
  MessageSquare,
  Minus,
  Plus,
  ShoppingCart,
  AlertCircle,
  Loader2,
  UserStar,
  Receipt
} from 'lucide-react';
import { format } from 'date-fns';
import { th } from 'date-fns/locale';
import { useToast } from '@/hooks/use-toast';
import { useBooking } from '@/hooks/useBooking';
import { useSeatUpdates } from '@/hooks/useWebSocketHooks';
import { currencyToThaiTextFull } from '@/utils/currency';
import { cn } from '@/lib/utils';
import type { Tour, Period } from '@/types/tours-regular';
import type { CreateBookingRequest } from '@/types';

// Extended interface for admin booking
interface AdminCreateBookingRequest extends CreateBookingRequest {
  user_id?: string; // สำหรับ admin/super เลือก user ที่จะทำการจองให้
}

interface BookingSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  period: Period | null;
  tour: Tour | null;
  onBookingSuccess?: (bookingId: string) => void;
  userRole?: string; // 'admin' | 'super' | 'agent' | etc.
}

interface BookingFormData {
  customer_name: string;
  customer_phone: string;
  note: string;
  sales_id: string;
  user_id?: string; // สำหรับ admin/super เลือก user ที่จะทำการจองให้
  room_types: { id: string; quantity: number; title?: string; price?: number }[];
  rate_types: { id: string; quantity: number; title?: string; price?: number; exclude?: boolean }[];
  additional_services: { id: string; quantity: number; title?: string; price?: number; rate_type?: string }[];
}

interface FormErrors {
  customer_name?: boolean;
  customer_phone?: boolean;
  sales_id?: boolean;
  user_id?: boolean; // สำหรับ admin/super validation
  rate_types?: boolean;
  room_types?: boolean;
}

export function BookingSheet({
  isOpen,
  onOpenChange,
  period,
  tour,
  onBookingSuccess,
  userRole
}: BookingSheetProps) {
  const { toast } = useToast();
  const { createBooking, isCreateLoading } = useBooking();

  // Real-time seat updates via WebSocket
  const { seatData } = useSeatUpdates(period?.id);

  // Use real-time seat data if available AND for correct period, fallback to period data
  const currentAvailableSeats = (seatData && seatData.period_id === period?.id)
    ? seatData.available_seats
    : period?.available_seats ?? 0;

  // Refs for scrolling to error fields
  const customerNameRef = useRef<HTMLInputElement>(null);
  const customerPhoneRef = useRef<HTMLInputElement>(null);
  const salesSelectRef = useRef<HTMLDivElement>(null);
  const userSelectRef = useRef<HTMLDivElement>(null);
  const rateTypesRef = useRef<HTMLDivElement>(null);
  const roomTypesRef = useRef<HTMLDivElement>(null);

  const getCurrentPrice = (rateItem: Period['rate_items'][0]) => {
    if (!rateItem || !period) return 0;
    return rateItem[period.rate_type as keyof typeof rateItem] as number;
  };

  // Check if user is admin or super
  const isAdminOrSuper = userRole === 'admin' || userRole === 'super';

  // Debug logging
  console.log('BookingSheet - userRole:', userRole);
  console.log('BookingSheet - isAdminOrSuper:', isAdminOrSuper);

  // TEMPORARY: Force admin mode for testing
  // const isAdminOrSuper = true;

  const initializeFormData = (): BookingFormData => ({
    customer_name: '',
    customer_phone: '',
    note: '',
    sales_id: '',
    user_id: isAdminOrSuper ? '' : undefined, // เฉพาะ admin/super เท่านั้น
    room_types: [],
    rate_types: period?.rate_items ? [{
      id: period.rate_items[0]?.id || '',
      quantity: 1,
      title: period.rate_items[0]?.title,
      price: getCurrentPrice(period.rate_items[0]),
      exclude: period.rate_items[0]?.exclude
    }] : [],
    additional_services: []
  });

  const [formData, setFormData] = useState<BookingFormData>(initializeFormData());
  const [errors, setErrors] = useState<FormErrors>({});

  // VAT option (ใบกำกับภาษี)
  const [includeVAT, setIncludeVAT] = useState(false);
  const VAT_RATE = 7; // 7%
  const WHT_NOTE_PREFIX = 'กรุณาทำการหัก ณ ที่จ่าย';

  // Reset form data when period/tour ID changes (not when seat data updates)
  useEffect(() => {
    if (period && tour && isOpen) {
      setFormData(initializeFormData());
      setIncludeVAT(false); // Reset VAT option
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [period?.id, tour?.id, isOpen]); // Only reset on ID changes, not period object changes

  if (!period || !tour) return null;

  const formatDateRange = () => {
    const start = new Date(period.start_date);
    const end = new Date(period.end_date);
    const diffDays = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;
    return {
      dateRange: `${format(start, 'dd MMM', { locale: th })} - ${format(end, 'dd MMM yyyy', { locale: th })}`,
      days: diffDays
    };
  };

  const calculateTotal = () => {
    const roomTotal = formData.room_types.reduce((sum, room) => sum + (room.price || 0) * room.quantity, 0);
    const rateTotal = formData.rate_types.reduce((sum, rate) => sum + (rate.price || 0) * rate.quantity, 0);
    const serviceTotal = formData.additional_services.reduce((sum, service) => {
      const quantity = service.rate_type === 'per_unit' ? service.quantity : 1;
      return sum + (service.price || 0) * quantity;
    }, 0);
    return roomTotal + rateTotal + serviceTotal;
  };

  const calculateNetTotal = () => {
    const total = calculateTotal();
    const commission = calculateCommission();
    return total - commission; // ยอดรวม - ส่วนลด agent = ราคาที่ลูกค้าจ่าย (ยอดก่อนภาษี)
  };

  // Calculate VAT amount (คิดจากค่าบริการเท่านั้น — include VAT: แยก VAT จากราคารวม)
  const calculateVATAmount = () => {
    if (!includeVAT) return 0;
    const totalPassengers = getTotalPassengers();
    const serviceFeePerPax = period.service_fee_per_pax || 0;
    if (serviceFeePerPax <= 0) return 0;
    const serviceFeeTotal = serviceFeePerPax * totalPassengers;
    // สูตร include VAT: VAT = serviceFeeTotal × vatRate / (100 + vatRate)
    return Math.round(serviceFeeTotal * VAT_RATE / (100 + VAT_RATE) * 100) / 100;
  };

  // คำนวณ WHT 3% สำหรับแสดงในหมายเหตุ (คิดจากค่าบริการก่อน VAT)
  const calculateWHTForNote = () => {
    const totalPassengers = getTotalPassengers();
    const serviceFeePerPax = period.service_fee_per_pax || 0;
    if (serviceFeePerPax <= 0 || totalPassengers <= 0) return 0;
    const serviceFeeTotal = serviceFeePerPax * totalPassengers;
    const vatAmount = serviceFeeTotal * VAT_RATE / (100 + VAT_RATE);
    const serviceFeeExcludeVAT = serviceFeeTotal - vatAmount;
    return Math.round(serviceFeeExcludeVAT * 0.03 * 100) / 100;
  };

  // Calculate grand total — VAT เป็น include (รวมในราคาแล้ว) ไม่ต้องบวกเพิ่ม, ไม่หัก WHT
  const calculateGrandTotal = () => {
    return calculateNetTotal();
  };

  const getTotalPassengers = () => {
    return formData.rate_types.reduce((sum, rate) => sum + rate.quantity, 0);
  };

  const hasInsufficientSeats = () => {
    const totalPassengers = getTotalPassengers();
    return totalPassengers > currentAvailableSeats;
  };

  const getCommissionEligiblePassengers = () => {
    return formData.rate_types
      .filter(rate => !rate.exclude) // เฉพาะที่ exclude = false เท่านั้น
      .reduce((sum, rate) => sum + rate.quantity, 0);
  };

  const calculateCommission = () => {
    const commissionPerPerson = period.commission_items?.reduce((sum, commission) => {
      return sum + (commission[period.commission_type as keyof typeof commission] as number);
    }, 0) || 0;

    return commissionPerPerson * getCommissionEligiblePassengers();
  };

  const handleQuantityChange = (
    type: 'room_types' | 'rate_types' | 'additional_services',
    index: number,
    newQuantity: number
  ) => {
    if (newQuantity < 1) return;

    // Check available seats limit for rate_types (passengers)
    if (type === 'rate_types') {
      const currentPassengers = getTotalPassengers();
      const currentItemQuantity = formData.rate_types[index]?.quantity || 0;
      const passengerDifference = newQuantity - currentItemQuantity;
      const newTotalPassengers = currentPassengers + passengerDifference;

      if (newTotalPassengers > currentAvailableSeats) {
        toast({
          title: "ไม่สามารถจองได้",
          description: `จำนวนผู้โดยสารเกินที่นั่งว่าง (ว่าง ${currentAvailableSeats} ที่นั่ง)`,
          variant: "destructive",
        });
        return;
      }
    }

    setFormData(prev => ({
      ...prev,
      [type]: prev[type].map((item, i) =>
        i === index ? { ...item, quantity: newQuantity } : item
      )
    }));
  };

  const addRoomType = (roomId: string, roomTitle: string, roomPrice: number) => {
    const existing = formData.room_types.find(room => room.id === roomId);
    if (existing) {
      handleQuantityChange('room_types', formData.room_types.indexOf(existing), existing.quantity + 1);
    } else {
      setFormData(prev => ({
        ...prev,
        room_types: [...prev.room_types, { id: roomId, quantity: 1, title: roomTitle, price: roomPrice }]
      }));
      // Clear room_types error when adding
      if (errors.room_types) {
        setErrors(prev => ({ ...prev, room_types: false }));
      }
    }
  };

  const addService = (serviceId: string, serviceTitle: string, servicePrice: number, rateType: string) => {
    const existing = formData.additional_services.find(service => service.id === serviceId);
    if (existing) {
      if (rateType === 'per_unit') {
        handleQuantityChange('additional_services', formData.additional_services.indexOf(existing), existing.quantity + 1);
      }
    } else {
      setFormData(prev => ({
        ...prev,
        additional_services: [...prev.additional_services, {
          id: serviceId,
          quantity: 1,
          title: serviceTitle,
          price: servicePrice,
          rate_type: rateType
        }]
      }));
    }
  };

  const removeItem = (type: 'room_types' | 'additional_services', index: number) => {
    setFormData(prev => ({
      ...prev,
      [type]: prev[type].filter((_, i) => i !== index)
    }));
  };

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (!formData.customer_name.trim()) {
      newErrors.customer_name = true;
    }

    if (!formData.customer_phone.trim()) {
      newErrors.customer_phone = true;
    }

    if (!formData.sales_id) {
      newErrors.sales_id = true;
    }

    // Validate user_id for admin/super
    if (isAdminOrSuper && !formData.user_id) {
      newErrors.user_id = true;
    }

    if (formData.rate_types.length === 0) {
      newErrors.rate_types = true;
    }

    if (formData.room_types.length === 0) {
      newErrors.room_types = true;
    }

    setErrors(newErrors);

    // Scroll to first error field
    if (Object.keys(newErrors).length > 0) {
      if (newErrors.customer_name && customerNameRef.current) {
        customerNameRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
        customerNameRef.current.focus();
      } else if (newErrors.customer_phone && customerPhoneRef.current) {
        customerPhoneRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
        customerPhoneRef.current.focus();
      } else if (newErrors.sales_id && salesSelectRef.current) {
        salesSelectRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
      } else if (newErrors.user_id && userSelectRef.current) {
        userSelectRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
      } else if (newErrors.room_types && roomTypesRef.current) {
        roomTypesRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
      } else if (newErrors.rate_types && rateTypesRef.current) {
        rateTypesRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    }

    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) {
      toast({
        title: "ข้อมูลไม่ครบถ้วน",
        description: "กรุณากรอกข้อมูลที่จำเป็นให้ครบถ้วน",
        variant: "destructive",
      });
      return;
    }

    try {
      const bookingData: AdminCreateBookingRequest = {
        tour_id: tour.id,
        period_id: period.id,
        sales_id: formData.sales_id,
        customer_name: formData.customer_name,
        customer_phone: formData.customer_phone,
        note: formData.note || undefined,
        room_types: formData.room_types.length > 0 ? formData.room_types.map(room => ({
          id: room.id,
          quantity: room.quantity
        })) : undefined,
        rate_types: formData.rate_types.map(rate => ({
          id: rate.id,
          quantity: rate.quantity
        })),
        additional_services: formData.additional_services.length > 0 ? formData.additional_services.map(service => ({
          id: service.id,
          quantity: service.quantity
        })) : undefined,
        // VAT option (ใบกำกับภาษี)
        include_vat: includeVAT,
        include_wht: false,
        // เพิ่ม user_id สำหรับ admin/super
        ...(isAdminOrSuper && formData.user_id && { user_id: formData.user_id })
      };

      // Cast เป็น CreateBookingRequest เพราะ API ยังไม่รองรับ user_id อย่างเป็นทางการ
      const result = await createBooking(bookingData as CreateBookingRequest);

      toast({
        title: "สร้างการจองสำเร็จ",
        description: `เลขที่การจอง: ${result.invoice_id}`,
      });

      onOpenChange(false);
      if (onBookingSuccess) {
        onBookingSuccess(result.id);
      }

      // Reset form
      setFormData(initializeFormData());
      setIncludeVAT(false);

    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถสร้างการจองได้",
        variant: "destructive",
      });
    }
  };

  const { dateRange, days } = formatDateRange();

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-2xl overflow-y-auto p-0">
        <div className="p-6">
          <SheetHeader className="p-0">
            <SheetTitle className="text-left pr-6">
              จองทัวร์: {tour.title}
            </SheetTitle>
            <SheetDescription className="text-left">
              รหัสทัวร์: {tour.code} | {dateRange} ({days} วัน)
            </SheetDescription>
          </SheetHeader>

          <div className="mt-6 space-y-6">
            {/* Period Info */}
            <Card className='gap-0'>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg font-bold">ข้อมูลทัวร์</CardTitle>
              </CardHeader>
              <CardContent className="space-y-1">
                <div className="grid grid-cols-2 gap-1">
                  <div className="flex items-center gap-2">
                    <Calendar className="h-5 w-5 text-muted-foreground" />
                    <span className="text-sm">{dateRange}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Users className="h-5 w-5 text-muted-foreground" />
                    <span className="text-sm">ว่าง {currentAvailableSeats} ที่นั่ง</span>
                  </div>
                  {period.airline && (
                    <div className="flex items-center gap-2">
                      <Plane className="h-5 w-5 text-muted-foreground" />
                      <span className="text-sm">{period.airline}</span>
                    </div>
                  )}
                  {period.airport && (
                    <div className="flex items-center gap-2">
                      <MapPin className="h-5 w-5 text-muted-foreground" />
                      <span className="text-sm">{period.airport}</span>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Customer Info */}
            <Card className='gap-0'>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg font-bold">ข้อมูลลูกค้า</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="grid grid-cols-2 gap-2">
                  <div className="space-y-2">
                    <Label htmlFor="customer_name">
                      <User className="h-5 w-5 inline" />
                      ชื่อ-นามสกุล *
                    </Label>
                    <Input
                      id="customer_name"
                      ref={customerNameRef}
                      value={formData.customer_name}
                      onChange={(e) => {
                        setFormData(prev => ({ ...prev, customer_name: e.target.value }));
                        if (errors.customer_name) {
                          setErrors(prev => ({ ...prev, customer_name: false }));
                        }
                      }}
                      placeholder="กรอกชื่อ-นามสกุลลูกค้า"
                      className={cn(errors.customer_name && "border-red-500 focus-visible:ring-red-500")}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="customer_phone">
                      <Phone className="h-5 w-5 inline" />
                      เบอร์โทรศัพท์ *
                    </Label>
                    <Input
                      id="customer_phone"
                      ref={customerPhoneRef}
                      value={formData.customer_phone}
                      onChange={(e) => {
                        setFormData(prev => ({ ...prev, customer_phone: e.target.value }));
                        if (errors.customer_phone) {
                          setErrors(prev => ({ ...prev, customer_phone: false }));
                        }
                      }}
                      placeholder="กรอกเบอร์โทรศัพท์"
                      className={cn(errors.customer_phone && "border-red-500 focus-visible:ring-red-500")}
                    />
                  </div>
                </div>
                <div className="space-y-2" ref={salesSelectRef}>
                  <Label htmlFor="sales_id">
                    <UserStar className="h-5 w-5 inline" />
                    เซลส์ *
                  </Label>
                  <div className={cn(errors.sales_id && "border border-red-500 rounded-md")}>
                    <SalesSelect
                      value={formData.sales_id}
                      onValueChange={(value) => {
                        setFormData(prev => ({ ...prev, sales_id: value || '' }));
                        if (errors.sales_id) {
                          setErrors(prev => ({ ...prev, sales_id: false }));
                        }
                      }}
                      placeholder="เลือกเซลส์"
                    />
                  </div>
                </div>

                {/* User Select for Admin/Super */}
                {isAdminOrSuper && (
                  <div className="space-y-2" ref={userSelectRef}>
                    <Label htmlFor="user_id">
                      <User className="h-5 w-5 inline" />
                      จองให้ผู้ใช้ *
                    </Label>
                    <div className={cn(errors.user_id && "border border-red-500 rounded-md")}>
                      <UserSelect
                        value={formData.user_id}
                        onValueChange={(value) => {
                          setFormData(prev => ({ ...prev, user_id: value || '' }));
                          if (errors.user_id) {
                            setErrors(prev => ({ ...prev, user_id: false }));
                          }
                        }}
                        placeholder="เลือกผู้ใช้ที่จะทำการจองให้"
                      />
                    </div>
                  </div>
                )}

                <div className="space-y-2">
                  <Label htmlFor="note">
                    <MessageSquare className="h-5 w-5 inline" />
                    หมายเหตุ
                  </Label>
                  <Textarea
                    id="note"
                    value={formData.note}
                    onChange={(e) => setFormData(prev => ({ ...prev, note: e.target.value }))}
                    placeholder="หมายเหตุเพิ่มเติม (ถ้ามี)"
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Rate Types */}
            <Card className={cn('gap-0', errors.rate_types && 'border-red-500')} ref={rateTypesRef}>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg font-bold">ประเภทผู้โดยสาร *</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                {formData.rate_types.map((rate, index) => (
                  <div key={index} className="flex items-center justify-between p-2 border rounded-lg">
                    <div>
                      <div className="font-medium flex items-center gap-2">
                        {rate.title}
                        {rate.exclude && (
                          <Badge variant="secondary" className="text-xs h-4">
                            ไม่ได้ส่วนลด
                          </Badge>
                        )}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        ฿{rate.price?.toLocaleString()} / ที่นั่ง
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Button
                        variant="outline"
                        size="icon"
                        onClick={() => handleQuantityChange('rate_types', index, rate.quantity - 1)}
                        disabled={rate.quantity <= 1}
                      >
                        <Minus className="h-4 w-4" />
                      </Button>
                      <span className="w-8 text-center">{rate.quantity}</span>
                      <Button
                        variant="outline"
                        size="icon"
                        onClick={() => handleQuantityChange('rate_types', index, rate.quantity + 1)}
                      >
                        <Plus className="h-4 w-4" />
                      </Button>
                      {formData.rate_types.length > 1 && (
                        <Button
                          variant="destructive"
                          size="icon"
                          onClick={() => setFormData(prev => ({
                            ...prev,
                            rate_types: prev.rate_types.filter((_, i) => i !== index)
                          }))}
                        >
                          <Minus className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                  </div>
                ))}

                {/* Add more rate types */}
                <div className="grid gap-2">
                  {period.rate_items?.filter(rateItem =>
                    !formData.rate_types.find(selected => selected.id === rateItem.id)
                  ).map((rateItem) => (
                    <Button
                      key={rateItem.id}
                      variant="outline"
                      className="justify-start"
                      onClick={() => {
                        const currentPassengers = getTotalPassengers();
                        if (currentPassengers >= currentAvailableSeats) {
                          toast({
                            title: "ไม่สามารถเพิ่มได้",
                            description: `จำนวนผู้โดยสารเต็มแล้ว (ว่าง ${currentAvailableSeats} ที่นั่ง)`,
                            variant: "destructive",
                          });
                          return;
                        }

                        const price = getCurrentPrice(rateItem);
                        const existing = formData.rate_types.find(rate => rate.id === rateItem.id);
                        if (existing) {
                          handleQuantityChange('rate_types', formData.rate_types.indexOf(existing), existing.quantity + 1);
                        } else {
                          setFormData(prev => ({
                            ...prev,
                            rate_types: [...prev.rate_types, {
                              id: rateItem.id,
                              quantity: 1,
                              title: rateItem.title,
                              price,
                              exclude: rateItem.exclude
                            }]
                          }));
                          // Clear rate_types error when adding
                          if (errors.rate_types) {
                            setErrors(prev => ({ ...prev, rate_types: false }));
                          }
                        }
                      }}
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      <span className="flex items-center gap-2">
                        {rateItem.title} - ฿{getCurrentPrice(rateItem).toLocaleString()}
                        {rateItem.exclude && (
                          <Badge variant="secondary" className="text-xs h-4">
                            ไม่ได้ส่วนลด
                          </Badge>
                        )}
                      </span>
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Room Types */}
            {period.room_types && period.room_types.length > 0 && (
              <Card className={cn('gap-0', errors.room_types && 'border-red-500')} ref={roomTypesRef}>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg font-bold">ประเภทห้องพัก *</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  {formData.room_types.map((room, index) => (
                    <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                      <div>
                        <div className="font-medium">{room.title}</div>
                        <div className="text-sm text-muted-foreground">
                          ฿{room.price?.toLocaleString()} / ห้อง
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Button
                          variant="outline"
                          size="icon"
                          onClick={() => handleQuantityChange('room_types', index, room.quantity - 1)}
                          disabled={room.quantity <= 1}
                        >
                          <Minus className="h-4 w-4" />
                        </Button>
                        <span className="w-8 text-center">{room.quantity}</span>
                        <Button
                          variant="outline"
                          size="icon"
                          onClick={() => handleQuantityChange('room_types', index, room.quantity + 1)}
                        >
                          <Plus className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="destructive"
                          size="icon"
                          onClick={() => removeItem('room_types', index)}
                        >
                          <Minus className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}

                  <div className="grid gap-2">
                    {period.room_types?.filter(roomType =>
                      !formData.room_types.find(selected => selected.id === roomType.id)
                    ).map((roomType) => (
                      <Button
                        key={roomType.id}
                        variant="outline"
                        className="justify-start"
                        onClick={() => addRoomType(roomType.id, roomType.title, roomType.price)}
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        {roomType.title} - ฿{roomType.price.toLocaleString()}
                      </Button>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Additional Services */}
            {period.additional_services && period.additional_services.length > 0 && (
              <Card className='gap-0'>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg font-bold">บริการเพิ่มเติม</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                  {formData.additional_services.map((service, index) => (
                    <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                      <div>
                        <div className="font-medium">{service.title}</div>
                        <div className="text-sm text-muted-foreground">
                          ฿{service.price?.toLocaleString()}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {service.rate_type === 'per_unit' ? (
                          <>
                            <Button
                              variant="outline"
                              size="icon"
                              onClick={() => handleQuantityChange('additional_services', index, service.quantity - 1)}
                              disabled={service.quantity <= 1}
                            >
                              <Minus className="h-4 w-4" />
                            </Button>
                            <span className="w-8 text-center">{service.quantity}</span>
                            <Button
                              variant="outline"
                              size="icon"
                              onClick={() => handleQuantityChange('additional_services', index, service.quantity + 1)}
                            >
                              <Plus className="h-4 w-4" />
                            </Button>
                          </>
                        ) : (
                          <span className="text-sm text-muted-foreground px-3">เหมาจ่าย</span>
                        )}
                        <Button
                          variant="destructive"
                          size="icon"
                          onClick={() => removeItem('additional_services', index)}
                        >
                          <Minus className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}

                  <div className="grid gap-2">
                    {period.additional_services?.filter(service =>
                      !formData.additional_services.find(selected => selected.id === service.id)
                    ).map((service) => (
                      <Button
                        key={service.id}
                        variant="outline"
                        className="justify-start"
                        onClick={() => addService(service.id, service.title, service.price, service.rate_type)}
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        {service.title} - ฿{service.price.toLocaleString()}
                        <span className="ml-auto text-xs text-muted-foreground">
                          {service.rate_type === 'per_unit' ? `ต่อ / ${service.unit_name || 'หน่วย'}` : 'เหมาจ่าย'}
                        </span>
                      </Button>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* ใบกำกับภาษี */}
            <Card className='gap-0'>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg font-bold flex items-center gap-2">
                  <Receipt className="h-5 w-5" />
                  ตัวเลือกภาษี
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label className="text-base">ใบกำกับภาษี (VAT 7%)</Label>
                    <p className="text-sm text-muted-foreground">
                      คือ VAT รวมในค่าบริการ ฿{(period.service_fee_per_pax || 0).toLocaleString()}/คน
                    </p>
                  </div>
                  <Switch
                    checked={includeVAT}
                    onCheckedChange={(checked) => {
                      setIncludeVAT(checked);
                      if (checked) {
                        // Auto-fill note กับ WHT info
                        const wht = calculateWHTForNote();
                        if (wht > 0) {
                          const whtText = `${WHT_NOTE_PREFIX} ฿${wht.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
                          setFormData(prev => ({
                            ...prev,
                            note: prev.note
                              ? (prev.note.includes(WHT_NOTE_PREFIX) ? prev.note.replace(new RegExp(`${WHT_NOTE_PREFIX}.*`), whtText) : `${prev.note}\n${whtText}`)
                              : whtText
                          }));
                        }
                      } else {
                        // ลบ WHT note ออก
                        setFormData(prev => ({
                          ...prev,
                          note: prev.note?.replace(new RegExp(`\\n?${WHT_NOTE_PREFIX}.*`), '').trim() || ''
                        }));
                      }
                    }}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Total */}
            <Card className='gap-0'>
              <CardContent className="pt-4">
                <div className="space-y-2">
                  {/* รายการที่เลือก */}
                  <div className="space-y-2">
                    <h4 className="font-bold text-lg">รายการที่เลือก:</h4>

                    {/* ประเภทผู้โดยสาร */}
                    {formData.rate_types.map((rate, index) => (
                      <div key={index} className="flex justify-between text-sm">
                        <span>
                          {rate.title} x {rate.quantity}
                          {rate.exclude && ' (ไม่ได้ส่วนลด)'}
                        </span>
                        <span>฿{((rate.price || 0) * rate.quantity).toLocaleString()}</span>
                      </div>
                    ))}

                    {/* ห้องพัก */}
                    {formData.room_types.map((room, index) => (
                      <div key={index} className="flex justify-between text-sm">
                        <span>{room.title} x {room.quantity}</span>
                        <span>฿{((room.price || 0) * room.quantity).toLocaleString()}</span>
                      </div>
                    ))}

                    {/* บริการเพิ่มเติม */}
                    {formData.additional_services.map((service, index) => (
                      <div key={index} className="flex justify-between text-sm">
                        <span>
                          {service.title}
                          {service.rate_type === 'per_unit' && ` x ${service.quantity}`}
                          {service.rate_type !== 'per_unit' && ' (เหมาจ่าย)'}
                        </span>
                        <span>฿{((service.price || 0) * (service.rate_type === 'per_unit' ? service.quantity : 1)).toLocaleString()}</span>
                      </div>
                    ))}
                  </div>

                  <Separator />

                  {/* ข้อมูลเพิ่มเติม */}
                  <div className="space-y-1 text-sm text-muted-foreground">
                    <div className="flex justify-between">
                      <span>จำนวนผู้โดยสารทั้งหมด:</span>
                      <span>{getTotalPassengers()} ที่นั่ง</span>
                    </div>
                    <div className="flex justify-between">
                      <span>ผู้โดยสารที่ได้ส่วนลด:</span>
                      <span>{getCommissionEligiblePassengers()} ที่นั่ง</span>
                    </div>
                    {period.commission_items && period.commission_items.length > 0 && (
                      <>
                        <div className="flex justify-between">
                          <span>เรทส่วนลด:</span>
                          <span>
                            {period.commission_items.map((commission, index) => (
                              <span key={index}>
                                ฿{(commission[period.commission_type as keyof typeof commission] as number).toLocaleString()}
                                {index < period.commission_items.length - 1 && ' + '}
                              </span>
                            ))}
                            {period.commission_items.length > 1 && (
                              <span className="ml-1">
                                = ฿{period.commission_items.reduce((sum, commission) =>
                                  sum + (commission[period.commission_type as keyof typeof commission] as number), 0
                                ).toLocaleString()}
                              </span>
                            )}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span>ส่วนลดรวม:</span>
                          <span>
                            ฿{period.commission_items.reduce((sum, commission) =>
                              sum + (commission[period.commission_type as keyof typeof commission] as number), 0
                            ).toLocaleString()} × {getCommissionEligiblePassengers()} ที่นั่ง = ฿{calculateCommission().toLocaleString()}
                          </span>
                        </div>
                      </>
                    )}
                  </div>

                  <Separator />

                  {/* ยอดรวม (ก่อนลดส่วนลด) */}
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <span>ยอดรวม (ก่อนหักส่วนลด):</span>
                    <span>฿{calculateTotal().toLocaleString()}</span>
                  </div>

                  {/* ส่วนลด Agent */}
                  {period.commission_items && period.commission_items.length > 0 && calculateCommission() > 0 && (
                    <div className="flex items-center justify-between text-sm">
                      <span>ส่วนลด:</span>
                      <span>-฿{calculateCommission().toLocaleString()}</span>
                    </div>
                  )}

                  {/* ยอดก่อนภาษี */}
                  <div className="flex items-center justify-between text-sm">
                    <span>ยอดก่อนภาษี:</span>
                    <span>฿{calculateNetTotal().toLocaleString()}</span>
                  </div>

                  {/* ใบกำกับภาษี — แสดงค่าบริการแยก VAT */}
                  {includeVAT && (() => {
                    const serviceFeeTotal = (period.service_fee_per_pax || 0) * getTotalPassengers();
                    const vatAmount = calculateVATAmount();
                    const feeExcludeVAT = Math.round((serviceFeeTotal - vatAmount) * 100) / 100;
                    const whtForNote = calculateWHTForNote();
                    return (
                      <>
                        <div className="flex items-center justify-between text-sm text-muted-foreground">
                          <span>ค่าบริการ:</span>
                          <span>฿{serviceFeeTotal.toLocaleString(undefined, { minimumFractionDigits: 2 })}</span>
                        </div>
                        <div className="flex items-center justify-between text-sm text-muted-foreground">
                          <span>VAT {VAT_RATE}%:</span>
                          <span>-฿{vatAmount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                        </div>
                        <div className="flex items-center justify-between text-sm font-medium">
                          <span>ค่าบริการก่อน VAT:</span>
                          <span>฿{feeExcludeVAT.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                        </div>
                        {whtForNote > 0 && (
                          <div className="text-sm text-orange-600 dark:text-orange-400">
                            * {WHT_NOTE_PREFIX} ฿{whtForNote.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                          </div>
                        )}
                      </>
                    );
                  })()}

                  <Separator />

                  {/* ยอดที่ต้องชำระ (Grand Total) */}
                  <div className="space-y-1">
                    <div className="flex items-center justify-between text-lg font-bold">
                      <div className="flex items-center gap-2">
                        <DollarSign className="h-5 w-5" />
                        ยอดที่ต้องชำระ
                      </div>
                      <div>฿{calculateGrandTotal().toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</div>
                    </div>
                    <div className="text-xs text-muted-foreground text-right">
                      ({currencyToThaiTextFull(calculateGrandTotal())})
                    </div>
                  </div>

                  <Separator />

                  {/* ยอดมัดจำ - คิดจากยอดที่ต้องชำระ (Grand Total) */}
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <div className="flex items-center gap-2">
                      <Banknote className="h-5 w-5" />
                      ยอดเงินมัดจำ ({period.deposit_type === 'percent' ? `${period.deposit}%` : `฿${period.deposit.toLocaleString()}/ที่นั่ง`})
                    </div>
                    <div>
                      ฿{(period.deposit_type === 'percent'
                        ? calculateGrandTotal() * (period.deposit / 100) // คิดจาก Grand Total
                        : period.deposit * (formData.rate_types.reduce((sum, rate) => sum + rate.quantity, 0))
                      ).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Submit */}
            <div className="flex gap-3 ">
              <Button
                variant="outline"
                className="flex-1 text-xl font-bold cursor-pointer"
                onClick={() => onOpenChange(false)}
                disabled={isCreateLoading}
              >
                ยกเลิก
              </Button>
              <Button
                className="flex-1 text-xl font-bold cursor-pointer"
                onClick={handleSubmit}
                disabled={isCreateLoading || hasInsufficientSeats()}
              >
                {isCreateLoading ? (
                  <>
                    <Loader2 className="h-5 w-5 animate-spin" />
                    กำลังสร้างการจอง...
                  </>
                ) : hasInsufficientSeats() ? (
                  <>
                    <AlertCircle className="h-5 w-5" />
                    ที่นั่งไม่เพียงพอ
                  </>
                ) : (
                  <>
                    <ShoppingCart className="h-5 w-5" />
                    ยืนยันการจอง
                  </>
                )}
              </Button>
            </div>

            {/* Seat Availability Warnings */}
            {hasInsufficientSeats() && (
              <div className="flex items-center gap-2 p-3 bg-destructive/10 border border-destructive rounded-lg">
                <AlertCircle className="h-5 w-5 text-destructive" />
                <div className="text-sm">
                  <div className="font-medium text-destructive">ไม่สามารถจองได้</div>
                  <div className="text-destructive/80">
                    เลือกผู้โดยสาร {getTotalPassengers()} ที่นั่ง แต่เหลือที่นั่งว่างเพียง {currentAvailableSeats} ที่นั่ง
                  </div>
                </div>
              </div>
            )}

            {/* General Warning */}
            {!hasInsufficientSeats() && currentAvailableSeats <= 10 && (
              <div className="flex items-center gap-2 p-3 bg-muted border border-muted rounded-lg">
                <AlertCircle className="h-5 w-5 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">
                  เหลือที่นั่งเพียง {currentAvailableSeats} ที่นั่ง
                 
                </span>
              </div>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}