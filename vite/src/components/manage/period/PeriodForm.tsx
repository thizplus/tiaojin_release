import { useState } from 'react';
import { useForm, type Resolver } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Calendar, Plane, DollarSign, Users, Bed, FileText } from 'lucide-react';
import { periodFormSchema, type PeriodFormData } from './schemas/period-form-schemas';
import { PeriodBasicInfo } from './sections/PeriodBasicInfo';
import { PeriodFlightSectors } from './sections/PeriodFlightSectors';
import { PeriodRateItems } from './sections/PeriodRateItems';
import { PeriodCommission } from './sections/PeriodCommission';
import { PeriodRoomTypes } from './sections/PeriodRoomTypes';
import { PeriodAdditionalServices } from './sections/PeriodAdditionalServices';
import { PeriodFiles } from './sections/PeriodFiles';
import type { TourPeriod } from '@/types/models/period';
import type { CreatePeriodRequest, UpdatePeriodRequest } from '@/types/request/period';

interface PeriodFormProps {
  tourId: string;
  period?: TourPeriod | null;
  onSubmit: (data: CreatePeriodRequest | UpdatePeriodRequest) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

// Create a type-safe resolver
const createFormResolver = (): Resolver<PeriodFormData> => {
  return zodResolver(periodFormSchema) as Resolver<PeriodFormData>;
};

export function PeriodForm({ tourId, period, onSubmit, onCancel, loading = false }: PeriodFormProps) {
  const isEdit = !!period;
  const [currentSection, setCurrentSection] = useState('basic');

  const sections = [
    { id: 'basic', label: 'ข้อมูลพื้นฐาน', icon: Calendar },
    { id: 'flights', label: 'เที่ยวบิน', icon: Plane },
    { id: 'pricing', label: 'ราคาและคอมมิชชั่น', icon: DollarSign },
    { id: 'rooms', label: 'ประเภทห้องพัก', icon: Bed },
    { id: 'services', label: 'บริการเสริม', icon: Users },
    ...(isEdit ? [{ id: 'files', label: 'ไฟล์แนบ', icon: FileText }] : []),
  ];

  const form = useForm<PeriodFormData>({
    resolver: createFormResolver(),
    defaultValues: {
      tour_id: tourId,
      start_date: period?.start_date?.split('T')[0] || '',
      end_date: period?.end_date?.split('T')[0] || '',
      seat: period?.seat || 0,
      deposit: period?.deposit || 0,
      deposit_type: period?.deposit_type || 'per_unit',
    //  single_stay: period?.single_stay || 0,
      airline: period?.airline || '',
      airport: period?.airport || '',
      note: period?.note || '',
      sector_1: period?.sector_1 || [{ date: '', flight: '', route: '' }],
      sector_2: period?.sector_2 || [{ date: '', flight: '', route: '' }],
      rate_type: period?.rate_type || 'regular',
      rate_items: period?.rate_items || [
        { title: 'ผู้ใหญ่', regular: 0, sale: 0, exclude: false }
      ],
      commission_type: period?.commission_type || 'regular',
      commission_items: period?.commission_items || [
        { regular: 0, sale: 0 },
        { regular: 0, sale: 0 }
      ],
      room_types: period?.room_types || [
        { title: 'single', price: 3900 },
        { title: 'double', price: 0 },
        { title: 'twin', price: 0 },
        { title: 'triple', price: 0 }
      ],
      additional_services: period?.additional_services || [],
      files: period?.files || [],
      status: period?.status || 'publish',
      service_fee_per_pax: period?.service_fee_per_pax || 0,

      // Expiration Override - แปลงจาก API format เป็น form format
      expiration_override: period?.expiration_override || false,
      expiration_unit: period?.expiration_minutes ? 'minutes' : period?.expiration_hours ? 'hours' : period?.expiration_days ? 'days' : 'hours',
      expiration_value: period?.expiration_minutes || period?.expiration_hours || period?.expiration_days || undefined,
      expiration_message: period?.expiration_message || '',
    },
  });

  const handleSubmit = async (data: PeriodFormData) => {
    console.log('Form data before submit:', data);
    console.log('Files data:', data.files);

    // แปลง expiration_unit + expiration_value กลับเป็น expiration_minutes/hours/days
    const expirationData: {
      expiration_override?: boolean;
      expiration_minutes?: number;
      expiration_hours?: number;
      expiration_days?: number;
      expiration_message?: string;
    } = {};

    if (data.expiration_override && data.expiration_value) {
      expirationData.expiration_override = true;
      expirationData.expiration_message = data.expiration_message || undefined;

      switch (data.expiration_unit) {
        case 'minutes':
          expirationData.expiration_minutes = data.expiration_value;
          break;
        case 'hours':
          expirationData.expiration_hours = data.expiration_value;
          break;
        case 'days':
          expirationData.expiration_days = data.expiration_value;
          break;
      }
    } else {
      expirationData.expiration_override = false;
    }

    // Remove form-specific fields and add API-specific fields
    const { expiration_unit, expiration_value, ...restData } = data;

    const submitData = {
      ...restData,
      ...expirationData,
      start_date: `${data.start_date}T00:00:00Z`,
      end_date: `${data.end_date}T00:00:00Z`,
    };

    console.log('Submit data:', submitData);
    await onSubmit(submitData);
  };

  return (
    <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">
          {isEdit ? 'แก้ไขช่วงเวลาทัวร์' : 'เพิ่มช่วงเวลาทัวร์'}
        </h2>
        <Badge
          variant="outline"
          className={`text-sm ${form.watch('status') === 'publish' ? 'bg-muted text-chart-2 border-chart-2/50' : 'bg-muted text-muted-foreground'}`}
        >
          <div className={`w-2 h-2 rounded-full ${form.watch('status') === 'publish' ? 'bg-chart-2' : 'bg-muted-foreground'}`}></div>
          {form.watch('status') === 'publish' ? 'เผยแพร่' : 'ฉบับร่าง'}
        </Badge>
      </div>

      {/* Mobile Select Navigation */}
      <div className="block lg:hidden">
        <Select value={currentSection} onValueChange={setCurrentSection}>
          <SelectTrigger className="w-full">
            <SelectValue>
              {(() => {
                const currentSectionData = sections.find(s => s.id === currentSection);
                if (currentSectionData) {
                  const Icon = currentSectionData.icon;
                  return (
                    <div className="flex items-center space-x-2">
                      <Icon className="h-4 w-4" />
                      <span>{currentSectionData.label}</span>
                    </div>
                  );
                }
                return null;
              })()}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {sections.map((section) => {
              const Icon = section.icon;
              return (
                <SelectItem key={section.id} value={section.id}>
                  <div className="flex items-center space-x-2">
                    <Icon className="h-4 w-4" />
                    <span>{section.label}</span>
                  </div>
                </SelectItem>
              );
            })}
          </SelectContent>
        </Select>
      </div>

      {/* Desktop Tabs Navigation */}
      <div className="hidden lg:block">
        <Tabs value={currentSection} onValueChange={setCurrentSection} className="space-y-6">
          <TabsList className={`grid w-full ${isEdit ? 'grid-cols-6' : 'grid-cols-5'}`}>
            {sections.map((section) => {
              const Icon = section.icon;
              return (
                <TabsTrigger key={section.id} value={section.id} className="flex items-center space-x-2">
                  <Icon className="h-4 w-4" />
                  <span>{section.label}</span>
                </TabsTrigger>
              );
            })}
          </TabsList>
        </Tabs>
      </div>

      {/* Content Sections */}
      <div className="space-y-6 mt-6">
        {currentSection === 'basic' && <PeriodBasicInfo form={form} />}
        {currentSection === 'flights' && <PeriodFlightSectors form={form} />}
        {currentSection === 'pricing' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <PeriodRateItems form={form} />
            <PeriodCommission form={form} />
          </div>
        )}
        {currentSection === 'rooms' && <PeriodRoomTypes form={form} />}
        {currentSection === 'services' && <PeriodAdditionalServices form={form} />}
        {currentSection === 'files' && isEdit && <PeriodFiles form={form} />}
      </div>

      <div className="flex justify-end space-x-4 pt-6 border-t">
        <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
          ยกเลิก
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? 'กำลังบันทึก...' : isEdit ? 'อัปเดต' : 'สร้าง'}
        </Button>
      </div>
    </form>
  );
}