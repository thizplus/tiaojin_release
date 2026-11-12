import { type UseFormReturn, useFieldArray } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Plus, Trash2, Bed, Hotel, DollarSign } from 'lucide-react';
import { type PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodRoomTypesProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodRoomTypes({ form }: PeriodRoomTypesProps) {
  const {
    fields: roomFields,
    append: appendRoom,
    remove: removeRoom,
  } = useFieldArray({
    control: form.control,
    name: 'room_types',
  });

  return (
    <Card className="shadow-sm border-0 bg-card">
      <CardHeader className="border-b border-border/50">
        <CardTitle className="flex items-center space-x-3 text-foreground">
          <div className="p-2 rounded-lg bg-muted">
            <Bed className="h-10 w-10 text-foreground" />
          </div>
          <div>
            <span className="text-lg font-semibold">ประเภทห้องพัก</span>
            <p className="text-sm text-muted-foreground font-normal mt-1">
              กำหนดประเภทห้องพักและราคาเพิ่มเติม ({roomFields.length} ประเภท)
            </p>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6 space-y-6">
        <div className="space-y-4">
          {roomFields.map((field, index) => (
            <div key={field.id} className="space-y-4 p-6 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Badge variant="outline" className="bg-muted text-foreground border-border/50">
                    ห้องที่ {index + 1}
                  </Badge>
                </div>
                {roomFields.length > 1 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => removeRoom(index)}
                    className="text-destructive hover:text-destructive hover:bg-destructive/10"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </div>

              {/* Room Fields - Single Row */}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                    <Hotel className="h-4 w-4" />
                    <span>ประเภทห้อง</span>
                  </Label>
                  <Input
                    placeholder="Single, Double, Twin, Triple"
                    className="bg-background"
                    {...form.register(`room_types.${index}.title`)}
                  />
                </div>

                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                    <DollarSign className="h-4 w-4" />
                    <span>ราคาเพิ่มเติม (บาท)</span>
                  </Label>
                  <Input
                    type="number"
                    min="0"
                    step="0.01"
                    placeholder="0"
                    className="bg-background"
                    {...form.register(`room_types.${index}.price`, { valueAsNumber: true })}
                  />
                </div>
              </div>
            </div>
          ))}

          <Button
            type="button"
            variant="outline"
            onClick={() => appendRoom({ title: '', price: 0 })}
            className="w-full"
          >
            <Plus className="h-4 w-4 mr-2" />
            เพิ่มประเภทห้องพัก
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}