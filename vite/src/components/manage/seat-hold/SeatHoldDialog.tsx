import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
import { DatePickerField } from '@/components/ui/date-picker-field';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Lock, Users, Calendar, Clock, Plus } from 'lucide-react';
import { useCreateSeatHold, usePeriodSeatHolds } from '@/hooks/useSeatHold';
import { usePeriodByIdQuery } from '@/stores/periodStore';
import { SEAT_HOLD_CONSTRAINTS } from '@/constants/seat-hold';
import type { TourPeriod } from '@/types';
import { SeatHoldList } from './SeatHoldList';

const createSeatHoldSchema = z.object({
  held_seats: z
    .number()
    .min(SEAT_HOLD_CONSTRAINTS.MIN_HELD_SEATS, `ต้องจองอย่างน้อย ${SEAT_HOLD_CONSTRAINTS.MIN_HELD_SEATS} ที่นั่ง`)
    .max(SEAT_HOLD_CONSTRAINTS.MAX_HELD_SEATS, `ไม่สามารถจองมากกว่า ${SEAT_HOLD_CONSTRAINTS.MAX_HELD_SEATS} ที่นั่งได้`),
  reason: z.string().optional(),
  expires_at: z.string().optional(),
});

type CreateSeatHoldFormData = z.infer<typeof createSeatHoldSchema>;

interface SeatHoldDialogProps {
  period: TourPeriod | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function SeatHoldDialog({ period, open, onOpenChange }: SeatHoldDialogProps) {
  const [activeTab, setActiveTab] = useState<'list' | 'create'>('list');

  // Use real-time period data instead of static prop
  const { data: currentPeriod} = usePeriodByIdQuery(
    period?.id || '',
    { enabled: open && !!period?.id }
  );

  // Use current period data if available, otherwise fallback to prop
  const activePeriod = currentPeriod || period;

  const createMutation = useCreateSeatHold(activePeriod?.id || '');
  const { data: seatHolds, isLoading } = usePeriodSeatHolds(activePeriod?.id || '', {}, { enabled: open && !!activePeriod?.id });

  const form = useForm<CreateSeatHoldFormData>({
    resolver: zodResolver(createSeatHoldSchema),
    defaultValues: {
      held_seats: 1,
      reason: '',
      expires_at: '',
    },
  });

  const onSubmit = async (data: CreateSeatHoldFormData) => {
    if (!activePeriod) return;

    try {
      // Convert expires_at to ISO string if provided
      const submitData = {
        ...data,
        expires_at: data.expires_at ? new Date(data.expires_at).toISOString() : undefined,
      };

      await createMutation.mutateAsync(submitData);
      form.reset();
      setActiveTab('list'); // Switch back to list view
    } catch (error) {
      // Error handling will be done by react-query
      console.error('Failed to create seat hold:', error);
    }
  };

  if (!activePeriod) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
          <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-muted">
              <Lock className="h-6 w-6 text-foreground" />
            </div>
            <div>
              <span className="text-lg">จัดการที่นั่งจำนอง</span>
              <p className="text-sm text-muted-foreground font-normal mt-1">
                ช่วงเวลา: {new Date(activePeriod.start_date).toLocaleDateString('th')} - {new Date(activePeriod.end_date).toLocaleDateString('th')}
              </p>
            </div>
          </DialogTitle>
          <DialogDescription asChild>
            <div className="flex items-center space-x-4 mt-4">
              <div className="flex items-center space-x-2">
                <Users className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">ที่นั่งทั้งหมด: <strong>{activePeriod.seat}</strong></span>
              </div>
              <div className="flex items-center space-x-2">
                <Badge variant="outline" className=" bg-muted text-muted-foreground">
                  ว่าง: {activePeriod.available_seats}
                </Badge>
                <Badge variant="outline" className=" bg-muted text-muted-foreground">
                  จอง: {activePeriod.booking}
                </Badge>
                <Badge variant="outline" className=" bg-muted text-muted-foreground">
                  จำนอง: {activePeriod.held_seats}
                </Badge>
              </div>
            </div>
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'list' | 'create')} className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="list" className="flex items-center space-x-2">
              <Lock className="h-4 w-4" />
              <span>รายการจำนอง ({seatHolds?.data?.total || 0})</span>
            </TabsTrigger>
            <TabsTrigger value="create" className="flex items-center space-x-2">
              <Plus className="h-4 w-4" />
              <span>จำนองที่นั่งใหม่</span>
            </TabsTrigger>
          </TabsList>

          <TabsContent value="list" className="mt-6">
            <SeatHoldList
              periodId={activePeriod.id}
              seatHolds={seatHolds?.data?.seat_holds || []}
              isLoading={isLoading}
            />
          </TabsContent>

          <TabsContent value="create" className="mt-6">
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                  <FormField
                    control={form.control}
                    name="held_seats"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="flex items-center space-x-2">
                          <Users className="h-4 w-4" />
                          <span>จำนวนที่นั่งที่ต้องการจำนอง</span>
                        </FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            min={SEAT_HOLD_CONSTRAINTS.MIN_HELD_SEATS}
                            max={Math.min(activePeriod.available_seats, SEAT_HOLD_CONSTRAINTS.MAX_HELD_SEATS)}
                            placeholder="1"
                            {...field}
                            onChange={(e) => field.onChange(parseInt(e.target.value) || 0)}
                          />
                        </FormControl>
                        <FormMessage />
                        <p className="text-xs text-muted-foreground">
                          ที่นั่งว่าง: {activePeriod.available_seats} | สูงสุด: {Math.min(activePeriod.available_seats, SEAT_HOLD_CONSTRAINTS.MAX_HELD_SEATS)} ที่นั่ง
                        </p>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="expires_at"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="flex items-center space-x-2">
                          <Calendar className="h-4 w-4" />
                          <span>วันหมดอายุ (ไม่บังคับ)</span>
                        </FormLabel>
                        <FormControl>
                          <DatePickerField
                            placeholder="เลือกวันหมดอายุ"
                            value={field.value}
                            onChange={(date) => {
                              if (date) {
                                const year = date.getFullYear();
                                const month = String(date.getMonth() + 1).padStart(2, '0');
                                const day = String(date.getDate()).padStart(2, '0');
                                field.onChange(`${year}-${month}-${day}`);
                              } else {
                                field.onChange('');
                              }
                            }}
                            className="w-full"
                          />
                        </FormControl>
                        <FormMessage />
                        <p className="text-xs text-muted-foreground">
                          หากไม่ระบุ จะหมดอายุตามค่าเริ่มต้น ({SEAT_HOLD_CONSTRAINTS.DEFAULT_HOLD_DURATION_HOURS} ชั่วโมง)
                        </p>
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={form.control}
                  name="reason"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>เหตุผลในการจำนอง</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="ระบุเหตุผลในการจำนองที่นั่ง (ไม่บังคับ)"
                          className="resize-none"
                          rows={3}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <DialogFooter>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => onOpenChange(false)}
                  >
                    ยกเลิก
                  </Button>
                  <Button
                    type="submit"
                    disabled={createMutation.isPending || activePeriod.available_seats === 0}
                  >
                    {createMutation.isPending && <Clock className="mr-2 h-4 w-4 animate-spin" />}
                    จำนองที่นั่ง
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </TabsContent>
          </Tabs>
        </div>
      </DialogContent>
    </Dialog>
  );
}