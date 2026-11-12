import { type UseFormReturn, useFieldArray } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Plus,
  Trash2,
  Percent,
  Tag,
  TrendingUp,
  TrendingDown,
  Eye,
  Info
} from 'lucide-react';
import { type PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodCommissionProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodCommission({ form }: PeriodCommissionProps) {
  const {
    fields: commissionFields,
    append: appendCommission,
    remove: removeCommission,
  } = useFieldArray({
    control: form.control,
    name: 'commission_items',
  });

  const commissionType = form.watch('commission_type');

  return (
    <Card className="shadow-sm border-0 bg-card">
      <CardHeader className="border-b border-border/50">
        <CardTitle className="flex items-center space-x-3 text-foreground">
          <div className="p-2 rounded-lg bg-muted">
            <Percent className="h-10 w-10 text-foreground" />
          </div>
          <div>
            <span className="text-lg font-semibold">คอมมิشชั่น</span>
            <p className="text-sm text-muted-foreground font-normal mt-1">
              กำหนดคอมมิชชั่นสำหรับตัวแทนจำหน่าย ({commissionFields.length} รายการ)
            </p>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6 space-y-8">
        {/* Commission Type Selection */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Tag className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">ประเภทคอมมิชชั่นที่ใช้</h3>
            <Badge variant="outline" className="ml-auto">
              <Info className="h-3 w-3 mr-1" />
              Required
            </Badge>
          </div>

          <div className="flex space-x-4">
            <div className="flex items-center space-x-2">
              <input
                type="radio"
                id="commission_regular"
                value="regular"
                {...form.register('commission_type')}
              />
              <label htmlFor="commission_regular" className="text-sm text-foreground">คอมปกติ</label>
            </div>
            <div className="flex items-center space-x-2">
              <input
                type="radio"
                id="commission_sale"
                value="sale"
                {...form.register('commission_type')}
              />
              <label htmlFor="commission_sale" className="text-sm text-foreground">คอมพิเศษ</label>
            </div>
          </div>
        </div>

        <Separator />

        {/* Commission Items */}
        <div className="space-y-6">
          <div className="flex items-center space-x-2 mb-4">
            <Percent className="h-5 w-5 text-muted-foreground" />
            <h3 className="font-semibold text-foreground">คอมมิชชั่นแต่ละระดับ</h3>
          </div>

          <div className="space-y-4">
            {commissionFields.map((field, index) => (
              <div key={field.id} className="space-y-5 p-6 rounded-lg bg-muted/30 border border-border/50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline" className="bg-muted text-foreground border-border/50">
                      คอมมิชชั่นที่ {index + 1}
                    </Badge>
                  </div>
                  {commissionFields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeCommission(index)}
                      className="text-destructive hover:text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>

                {/* Commission Fields - Single Row */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <TrendingUp className="h-4 w-4" />
                      <span>คอมปกติ (บาท)</span>
                    </Label>
                    <Input
                      type="number"
                      min="0"
                      step="0.01"
                      placeholder="850"
                      className="bg-background"
                      {...form.register(`commission_items.${index}.regular`, { valueAsNumber: true })}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <TrendingDown className="h-4 w-4" />
                      <span>คอมพิเศษ (บาท)</span>
                    </Label>
                    <Input
                      type="number"
                      min="0"
                      step="0.01"
                      placeholder="765"
                      className="bg-background"
                      {...form.register(`commission_items.${index}.sale`, { valueAsNumber: true })}
                    />
                  </div>
                </div>

                {/* Commission Preview */}
                <div className="p-4 rounded-lg bg-muted border border-border/50">
                  <div className="flex items-center space-x-2 mb-2">
                    <Eye className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm font-medium text-foreground">คอมมิชชั่นที่จะใช้</span>
                  </div>
                  <div className="flex items-baseline space-x-2">
                    <span className="text-xl font-bold text-foreground">
                      ฿{(form.watch(`commission_items.${index}.${commissionType}`) || 0).toLocaleString()}
                    </span>
                    <span className="text-sm text-muted-foreground">
                      ({commissionType === 'regular' ? 'คอมปกติ' : 'คอมพิเศษ'})
                    </span>
                    {commissionType === 'sale' && form.watch(`commission_items.${index}.regular`) && (
                      <span className="text-sm line-through text-muted-foreground">
                        ฿{form.watch(`commission_items.${index}.regular`).toLocaleString()}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}

            {/* Add New Commission Item */}
            <Button
              type="button"
              variant="outline"
              onClick={() => appendCommission({ regular: 0, sale: 0 })}
              className="w-full"
            >
              <Plus className="h-4 w-4 mr-2" />
              เพิ่มคอมมิชชั่น
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}