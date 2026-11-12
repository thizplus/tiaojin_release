import { type UseFormReturn } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { DatePickerField } from '@/components/ui/date-picker-field';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Calendar,
  MapPin,
  Plane,
  Users,
  CreditCard,
  FileText,
  Settings,
  Clock,
  Flame,
  Receipt
} from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { type PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodBasicInfoProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodBasicInfo({ form }: PeriodBasicInfoProps) {
  return (
    <Card className="shadow-sm border-0 bg-card">
      <CardHeader className=" border-b border-border/50">
        <CardTitle className="flex items-center space-x-3 text-foreground">
          <div className="p-2 rounded-lg bg-muted">
            <Calendar className="h-10 w-10 text-foreground" />
          </div>
          <div>
            <span className="text-lg font-semibold">ข้อมูลพื้นฐาน</span>
            <p className="text-sm text-muted-foreground font-normal mt-1">
              กำหนดช่วงเวลาและข้อมูลสำคัญของทัวร์
            </p>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6 space-y-8">
        {/* Date Range Section */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Calendar className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">ช่วงเวลาทัวร์</h3>
            <Badge variant="outline" className="ml-auto">
              <Calendar className="h-3 w-3 mr-1" />
              Required
            </Badge>
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
            <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Calendar className="h-5 w-5" />
                      <span>วันที่เริ่มต้น</span>
                    </Label>
              <DatePickerField
                id="start_date"
               // label="วันที่เริ่มต้น"
                placeholder="เลือกวันที่เริ่มต้น"
                value={form.watch('start_date')}
                onChange={(date) => {
                  if (date) {
                    const year = date.getFullYear()
                    const month = String(date.getMonth() + 1).padStart(2, '0')
                    const day = String(date.getDate()).padStart(2, '0')
                    form.setValue('start_date', `${year}-${month}-${day}`);
                  } else {
                    form.setValue('start_date', '');
                  }
                }}
                required
                error={form.formState.errors.start_date?.message}
              />
            </div>
            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
            <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Calendar className="h-5 w-5" />
                      <span>วันที่สิ้นสุด</span>
                    </Label>
              <DatePickerField
                id="end_date"
               // label="วันที่สิ้นสุด"
                placeholder="เลือกวันที่สิ้นสุด"
                value={form.watch('end_date')}
                onChange={(date) => {
                  if (date) {
                    const year = date.getFullYear()
                    const month = String(date.getMonth() + 1).padStart(2, '0')
                    const day = String(date.getDate()).padStart(2, '0')
                    form.setValue('end_date', `${year}-${month}-${day}`);
                  } else {
                    form.setValue('end_date', '');
                  }
                }}
                required
                error={form.formState.errors.end_date?.message}
              />
            </div>
            
            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <Plane className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="airline" className="text-sm font-medium text-foreground">สายการบิน</Label>
                <Badge variant="secondary" className="ml-auto">
                  Required
                </Badge>
              </div>
              <Input
                id="airline"
                placeholder="Thai Airways"
                className="bg-background"
                {...form.register('airline')}
              />
              {form.formState.errors.airline && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.airline.message}
                </p>
              )}
            </div>

            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <MapPin className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="airport" className="text-sm font-medium text-foreground">สนามบิน</Label>
                <Badge variant="secondary" className="ml-auto">
                  Required
                </Badge>
              </div>
              <Input
                id="airport"
                placeholder="BKK"
                className="bg-background"
                {...form.register('airport')}
              />
              {form.formState.errors.airport && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.airport.message}
                </p>
              )}
            </div>
          </div>
        </div>

        <Separator />

        {/* Capacity & Deposit Section */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Users className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">ที่นั่งและเงินมัดจำ</h3>
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <Users className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="seat" className="text-sm font-medium text-foreground">จำนวนที่นั่ง</Label>
                <Badge variant="secondary" className="ml-auto">
                  Required
                </Badge>
              </div>
              <Input
                id="seat"
                type="number"
                min="1"
                placeholder="25"
                className="bg-background"
                {...form.register('seat', { valueAsNumber: true })}
              />
              {form.formState.errors.seat && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.seat.message}
                </p>
              )}
            </div>

            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <CreditCard className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="deposit" className="text-sm font-medium text-foreground">เงินมัดจำ</Label>
                <Badge variant="secondary" className="ml-auto">
                  Required
                </Badge>
              </div>
              <Input
                id="deposit"
                type="number"
                min="0"
                step="0.01"
                placeholder="10"
                className="bg-background"
                {...form.register('deposit', { valueAsNumber: true })}
              />
              {form.formState.errors.deposit && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.deposit.message}
                </p>
              )}
            </div>

            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <Settings className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="deposit_type" className="text-sm font-medium text-foreground">ประเภทมัดจำ</Label>
              </div>
              <Select
                onValueChange={(value) => form.setValue('deposit_type', value as 'per_unit' | 'percent')}
                value={form.watch('deposit_type')}
              >
                <SelectTrigger className="bg-background">
                  <SelectValue placeholder="เลือกประเภทมัดจำ" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="per_unit">ต่อคน (บาท)</SelectItem>
                  <SelectItem value="percent">เปอร์เซ็นต์ (%)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <Receipt className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="service_fee_per_pax" className="text-sm font-medium text-foreground">ค่าบริการ/คน</Label>
                <Badge variant="outline" className="ml-auto">
                  Optional
                </Badge>
              </div>
              <Input
                id="service_fee_per_pax"
                type="number"
                min="0"
                step="0.01"
                placeholder="500"
                className="bg-background"
                {...form.register('service_fee_per_pax', { valueAsNumber: true })}
              />
              <p className="text-xs text-muted-foreground">
                บาท/คน (รวม VAT) — ใช้คำนวณภาษี VAT 7%
              </p>
              {form.formState.errors.service_fee_per_pax && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.service_fee_per_pax.message}
                </p>
              )}
            </div>
          </div>
        </div>

        <Separator />

     

        {/* Additional Information Section */}
        <div className="space-y-6">
          <div className="flex items-center space-x-2 mb-4">
            <FileText className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">ข้อมูลเพิ่มเติม</h3>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <FileText className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="note" className="text-sm font-medium text-foreground">หมายเหตุ</Label>
                <Badge variant="outline" className="ml-auto">
                  Optional
                </Badge>
              </div>
              <Textarea
                id="note"
                placeholder="หมายเหตุเพิ่มเติมหรือข้อมูลสำคัญของช่วงเวลาทัวร์นี้"
                rows={4}
                className="bg-background resize-none"
                {...form.register('note')}
              />
            </div>

            <div className="space-y-3 p-4 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center space-x-2">
                <Settings className="h-5 w-5 text-muted-foreground" />
                <Label htmlFor="status" className="text-sm font-medium text-foreground">สถานะการเผยแพร่</Label>
              </div>
              <Select
                onValueChange={(value) => form.setValue('status', value as 'draft' | 'publish')}
                value={form.watch('status')}
              >
                <SelectTrigger className="bg-background">
                  <SelectValue placeholder="เลือกสถานะ" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="draft">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 rounded-full bg-muted-foreground"></div>
                      <span>ฉบับร่าง</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="publish">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 rounded-full bg-chart-2"></div>
                      <span>เผยแพร่</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground mt-1">
                {form.watch('status') === 'publish'
                  ? 'ช่วงเวลาทัวร์นี้จะแสดงบนเว็บไซต์'
                  : 'ช่วงเวลาทัวร์นี้ยังไม่แสดงบนเว็บไซต์'
                }
              </p>
            </div>
          </div>
        </div>

        <Separator />

        {/* Expiration Override Section - ทัวร์ไฟไหม้ */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Clock className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">การหมดอายุ Invoice</h3>
            <Badge variant="outline" className="ml-auto">
              Optional
            </Badge>
          </div>

          <div className="p-4 rounded-lg bg-muted/30 border border-border/50 space-y-4">
            {/* Toggle Switch */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Flame className="h-5 w-5 text-orange-500" />
                <div>
                  <Label htmlFor="expiration_override" className="text-sm font-medium text-foreground">
                    กำหนดเวลาหมดอายุเอง
                  </Label>
                  <p className="text-xs text-muted-foreground">
                    สำหรับทัวร์ไฟไหม้หรือ Flash Sale ที่ต้องการกำหนดเวลาพิเศษ
                  </p>
                </div>
              </div>
              <Switch
                id="expiration_override"
                checked={form.watch('expiration_override') || false}
                onCheckedChange={(checked) => form.setValue('expiration_override', checked)}
              />
            </div>

            {/* Expiration Settings - แสดงเมื่อเปิด override */}
            {form.watch('expiration_override') && (
              <div className="pt-4 border-t border-border/50 space-y-4">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  {/* เวลาหมดอายุ */}
                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground">หมดอายุใน</Label>
                    <div className="flex space-x-2">
                      <Input
                        type="number"
                        min="1"
                        placeholder="เช่น 15"
                        className="bg-background w-24"
                        value={form.watch('expiration_value') || ''}
                        onChange={(e) => form.setValue('expiration_value', e.target.value ? parseInt(e.target.value) : undefined)}
                      />
                      <Select
                        value={form.watch('expiration_unit') || 'hours'}
                        onValueChange={(value) => form.setValue('expiration_unit', value as 'minutes' | 'hours' | 'days')}
                      >
                        <SelectTrigger className="bg-background w-28">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="minutes">นาที</SelectItem>
                          <SelectItem value="hours">ชั่วโมง</SelectItem>
                          <SelectItem value="days">วัน</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>

                  {/* ข้อความแสดง */}
                  <div className="space-y-2">
                    <Label htmlFor="expiration_message" className="text-sm font-medium text-foreground">
                      ข้อความแสดง (ไม่บังคับ)
                    </Label>
                    <Input
                      id="expiration_message"
                      placeholder="เช่น 🔥 ทัวร์ไฟไหม้! ยืนยันภายใน 15 นาที"
                      className="bg-background"
                      {...form.register('expiration_message')}
                    />
                  </div>
                </div>

                <p className="text-xs text-muted-foreground bg-orange-500/10 p-2 rounded">
                  💡 เมื่อมีการ Booking ใน Period นี้ Invoice จะหมดอายุตามเวลาที่กำหนด แทนที่จะใช้ค่าเริ่มต้นของระบบ
                </p>
              </div>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}