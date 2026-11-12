import { type UseFormReturn, useFieldArray } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Plus, Trash2, Users, Tag, DollarSign, Settings } from 'lucide-react';
import { type PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodAdditionalServicesProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodAdditionalServices({ form }: PeriodAdditionalServicesProps) {
  const {
    fields: serviceFields,
    append: appendService,
    remove: removeService,
  } = useFieldArray({
    control: form.control,
    name: 'additional_services',
  });

  return (
    <Card className="shadow-sm border-0 bg-card">
      <CardHeader className="border-b border-border/50">
        <CardTitle className="flex items-center space-x-3 text-foreground">
          <div className="p-2 rounded-lg bg-muted">
            <Users className="h-10 w-10 text-foreground" />
          </div>
          <div>
            <span className="text-lg font-semibold">บริการเสริม</span>
            <p className="text-sm text-muted-foreground font-normal mt-1">
              กำหนดบริการเสริมและราคาเพิ่มเติม ({serviceFields.length} บริการ)
            </p>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6 space-y-6">
        <div className="space-y-4">
          {serviceFields.map((field, index) => (
            <div key={field.id} className="space-y-4 p-6 rounded-lg bg-muted/30 border border-border/50">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Badge variant="outline" className="bg-muted text-foreground border-border/50">
                    บริการที่ {index + 1}
                  </Badge>
                </div>
                {serviceFields.length > 1 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => removeService(index)}
                    className="text-destructive hover:text-destructive hover:bg-destructive/10"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </div>

              {/* Service Fields - Single Row */}
              <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                    <Tag className="h-4 w-4" />
                    <span>ชื่อบริการ</span>
                  </Label>
                  <Input
                    placeholder="ประกันการเดินทาง"
                    className="bg-background"
                    {...form.register(`additional_services.${index}.title`)}
                  />
                </div>

                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                    <DollarSign className="h-4 w-4" />
                    <span>ราคา (บาท)</span>
                  </Label>
                  <Input
                    type="number"
                    min="0"
                    step="0.01"
                    placeholder="500"
                    className="bg-background"
                    {...form.register(`additional_services.${index}.price`, { valueAsNumber: true })}
                  />
                </div>

                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                    <Settings className="h-4 w-4" />
                    <span>ประเภทราคา</span>
                  </Label>
                  <Select
                    onValueChange={(value) => form.setValue(`additional_services.${index}.rate_type`, value as 'per_unit' | 'per_booking')}
                    value={form.watch(`additional_services.${index}.rate_type`)}
                  >
                    <SelectTrigger className="bg-background">
                      <SelectValue placeholder="เลือกประเภท" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="per_unit">ต่อหน่วย</SelectItem>
                      <SelectItem value="per_booking">ต่อการจอง</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                    <Tag className="h-4 w-4" />
                    <span>หน่วย</span>
                  </Label>
                  <Input
                    placeholder="คน, ใบ, ครั้ง"
                    className="bg-background"
                    {...form.register(`additional_services.${index}.unit_name`)}
                  />
                </div>
              </div>
            </div>
          ))}

          <Button
            type="button"
            variant="outline"
            onClick={() => appendService({ title: '', price: 0, rate_type: 'per_unit', unit_name: '' })}
            className="w-full"
          >
            <Plus className="h-4 w-4 mr-2" />
            เพิ่มบริการเสริม
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}