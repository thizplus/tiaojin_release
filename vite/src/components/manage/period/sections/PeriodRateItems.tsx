import { type UseFormReturn, useFieldArray } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Plus,
  Trash2,
  DollarSign,
  Tag,
  TrendingUp,
  TrendingDown,
  Eye,

  Info
} from 'lucide-react';
import { type PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodRateItemsProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodRateItems({ form }: PeriodRateItemsProps) {
  const {
    fields: rateFields,
    append: appendRate,
    remove: removeRate,
  } = useFieldArray({
    control: form.control,
    name: 'rate_items',
  });

  const rateType = form.watch('rate_type');

  return (
    <Card className="shadow-sm border-0 bg-card">
      <CardHeader className="border-b border-border/50">
        <CardTitle className="flex items-center space-x-3 text-foreground">
          <div className="p-2 rounded-lg bg-muted">
            <DollarSign className="h-10 w-10 text-foreground" />
          </div>
          <div>
            <span className="text-lg font-semibold">รายการราคา</span>
            <p className="text-sm text-muted-foreground font-normal mt-1">
              กำหนดราคาแต่ละประเภท ({rateFields.length} รายการ)
            </p>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6 space-y-8">
        {/* Rate Type Selection */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Tag className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">ประเภทราคาที่แสดง</h3>
            <Badge variant="outline" className="ml-auto">
              <Info className="h-3 w-3 mr-1" />
              Required
            </Badge>
          </div>

          <div className="flex space-x-4">
            <div className="flex items-center space-x-2">
              <input
                type="radio"
                id="rate_regular"
                value="regular"
                {...form.register('rate_type')}
              />
              <label htmlFor="rate_regular" className="text-sm text-foreground">ราคาปกติ</label>
            </div>
            <div className="flex items-center space-x-2">
              <input
                type="radio"
                id="rate_sale"
                value="sale"
                {...form.register('rate_type')}
              />
              <label htmlFor="rate_sale" className="text-sm text-foreground">ราคาพิเศษ</label>
            </div>
          </div>
        </div>

        <Separator />

        {/* Rate Items */}
        <div className="space-y-6">
          <div className="flex items-center space-x-2 mb-4">
            <DollarSign className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">รายการราคาแต่ละประเภท</h3>
          </div>

          <div className="space-y-4">
            {rateFields.map((field, index) => (
              <div key={field.id} className="space-y-5 p-6 rounded-lg bg-muted/30 border border-border/50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline" className="bg-muted text-foreground border-border/50">
                  
                      รายการที่ {index + 1}
                    </Badge>
                  </div>
                  {rateFields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeRate(index)}
                      className="text-destructive hover:text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>

                {/* Rate Item Fields - Single Row */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Tag className="h-4 w-4" />
                      <span>ชื่อรายการ</span>
                    </Label>
                    <Input
                      placeholder="ผู้ใหญ่, เด็กมีเตียง"
                      className="bg-background"
                      {...form.register(`rate_items.${index}.title`)}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <TrendingUp className="h-4 w-4" />
                      <span>ราคาปกติ (บาท)</span>
                    </Label>
                    <Input
                      type="number"
                      min="0"
                      step="0.01"
                      placeholder="15000"
                      className="bg-background"
                      {...form.register(`rate_items.${index}.regular`, { valueAsNumber: true })}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <TrendingDown className="h-4 w-4" />
                      <span>ราคาพิเศษ (บาท)</span>
                    </Label>
                    <Input
                      type="number"
                      min="0"
                      step="0.01"
                      placeholder="12000"
                      className="bg-background"
                      {...form.register(`rate_items.${index}.sale`, { valueAsNumber: true })}
                    />
                  </div>
                </div>

                {/* Price Preview */}
                <div className="p-4 rounded-lg bg-muted border border-border/50">
                  <div className="flex items-center space-x-2 mb-2">
                    <Eye className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm font-medium text-foreground">ราคาที่จะแสดงบนเว็บไซต์</span>
                  </div>
                  <div className="flex items-baseline space-x-2">
                    <span className="text-xl font-bold text-foreground">
                      ฿{(form.watch(`rate_items.${index}.${rateType}`) || 0).toLocaleString()}
                    </span>
                    <span className="text-sm text-muted-foreground">
                      ({rateType === 'regular' ? 'ราคาปกติ' : 'ราคาพิเศษ'})
                    </span>
                    {rateType === 'sale' && form.watch(`rate_items.${index}.regular`) && (
                      <span className="text-sm line-through text-muted-foreground">
                        ฿{form.watch(`rate_items.${index}.regular`).toLocaleString()}
                      </span>
                    )}
                  </div>
                </div>

                {/* Exclude Option */}
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id={`rate_exclude_${index}`}
                    checked={form.watch(`rate_items.${index}.exclude`) || false}
                    onCheckedChange={(checked) => {
                      form.setValue(`rate_items.${index}.exclude`, !!checked);
                    }}
                  />
                  <Label htmlFor={`rate_exclude_${index}`} className="text-sm text-foreground">
                  ไม่ได้ส่วนลด
                  </Label>
                </div>
              </div>
            ))}

            {/* Add New Rate Item */}
            <Button
              type="button"
              variant="outline"
              onClick={() => appendRate({ title: '', regular: 0, sale: 0, exclude: false })}
              className="w-full"
            >
              <Plus className="h-4 w-4 mr-2" />
              เพิ่มรายการราคา
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}