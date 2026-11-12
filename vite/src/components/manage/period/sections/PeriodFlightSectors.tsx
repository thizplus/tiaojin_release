import { type UseFormReturn, useFieldArray } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';

import {
  Plus,
  Trash2,
  Plane,
  Calendar,
  Route,

  ArrowRight
} from 'lucide-react';
import { DatePickerField } from '@/components/ui/date-picker-field';
import { type PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodFlightSectorsProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodFlightSectors({ form }: PeriodFlightSectorsProps) {
  const {
    fields: sector1Fields,
    append: appendSector1,
    remove: removeSector1,
  } = useFieldArray({
    control: form.control,
    name: 'sector_1',
  });

  const {
    fields: sector2Fields,
    append: appendSector2,
    remove: removeSector2,
  } = useFieldArray({
    control: form.control,
    name: 'sector_2',
  });

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Sector 1 */}
      <Card className="shadow-sm border-0 bg-card">
        <CardHeader className=" border-b border-border/50">
          <CardTitle className="flex items-center space-x-3 text-foreground">
            <div className="p-2 rounded-lg bg-muted">
              <Plane className="h-10 w-10 text-foreground" />
            </div>
            <div>
              <span className="text-lg font-semibold">เที่ยวบินไป</span>
              <p className="text-sm text-muted-foreground font-normal mt-1">
                กำหนดเที่ยวบินขาไป ({sector1Fields.length} เที่ยว)
              </p>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent className="p-6 space-y-6">
          <div className="space-y-4">
            {sector1Fields.map((field, index) => (
              <div key={field.id} className="space-y-4 p-5 rounded-lg bg-muted/30 border border-border/50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline" className="bg-muted text-sm text-foreground border-border/50">
                      <Plane className="h-5 w-5 mr-1" />
                      เที่ยวที่ {index + 1}
                    </Badge>
                  </div>
                  {sector1Fields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeSector1(index)}
                      className="text-destructive hover:text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>

                <div className="grid grid-cols-1 gap-4">
                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Calendar className="h-5 w-5" />
                      <span>วันที่</span>
                    </Label>
                    <DatePickerField
                      placeholder="เลือกวันที่บิน"
                      value={form.watch(`sector_1.${index}.date`)}
                      onChange={(date) => {
                        if (date) {
                          const year = date.getFullYear()
                          const month = String(date.getMonth() + 1).padStart(2, '0')
                          const day = String(date.getDate()).padStart(2, '0')
                          form.setValue(`sector_1.${index}.date`, `${year}-${month}-${day}`);
                        } else {
                          form.setValue(`sector_1.${index}.date`, '');
                        }
                      }}
                      className="w-full bg-background"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Plane className="h-5 w-5" />
                      <span>รหัสเที่ยวบิน</span>
                    </Label>
                    <Input
                      placeholder="TG635"
                      className="bg-background"
                      {...form.register(`sector_1.${index}.flight`)}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Route className="h-5 w-5" />
                      <span>เส้นทางการบิน</span>
                    </Label>
                    <div className="flex items-center space-x-2">
                      <Input
                        placeholder="HKG-BKK 08:15-10:45"
                        className="bg-background"
                        {...form.register(`sector_1.${index}.route`)}
                      />
                      <div className="flex-shrink-0 text-muted-foreground">
                        <ArrowRight className="h-4 w-4" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}

            <Button
              type="button"
              variant="outline"
              onClick={() => appendSector1({ date: '', flight: '', route: '' })}
              className="w-full border-dashed border-2 hover:border-chart-2/50 hover:bg-chart-2/5"
            >
              <Plus className="h-4 w-4 mr-2" />
              เพิ่มเที่ยวบิน
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Sector 2 */}
      <Card className="shadow-sm border-0 bg-card">
        <CardHeader className=" border-b border-border/50">
          <CardTitle className="flex items-center space-x-3 text-foreground">
            <div className="p-2 rounded-lg bg-muted">
              <Plane className="h-10 w-10 transform scale-x-[-1] text-foreground" />
            </div>
            <div>
              <span className="text-lg font-semibold">เที่ยวบินกลับ</span>
              <p className="text-sm text-muted-foreground font-normal mt-1">
                กำหนดเที่ยวบินขากลับ ({sector2Fields.length} เที่ยว)
              </p>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent className="p-6 space-y-6">
          <div className="space-y-4">
            {sector2Fields.map((field, index) => (
              <div key={field.id} className="space-y-4 p-5 rounded-lg bg-muted/30 border border-border/50">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline" className="bg-muted text-sm text-foreground border-border/50">
                      <Plane className="h-5 w-5 mr-1 transform scale-x-[-1] text-foreground" />
                      เที่ยวที่ {index + 1}
                    </Badge>
                  </div>
                  {sector2Fields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeSector2(index)}
                      className="text-destructive hover:text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>

                <div className="grid grid-cols-1 gap-4">
                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Calendar className="h-5 w-5" />
                      <span>วันที่</span>
                    </Label>
                    <DatePickerField
                      placeholder="เลือกวันที่บิน"
                      value={form.watch(`sector_2.${index}.date`)}
                      onChange={(date) => {
                        if (date) {
                          const year = date.getFullYear()
                          const month = String(date.getMonth() + 1).padStart(2, '0')
                          const day = String(date.getDate()).padStart(2, '0')
                          form.setValue(`sector_2.${index}.date`, `${year}-${month}-${day}`);
                        } else {
                          form.setValue(`sector_2.${index}.date`, '');
                        }
                      }}
                      className="w-full bg-background"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Plane className="h-5 w-5 transform scale-x-[-1]" />
                      <span>รหัสเที่ยวบิน</span>
                    </Label>
                    <Input
                      placeholder="TG636"
                      className="bg-background"
                      {...form.register(`sector_2.${index}.flight`)}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-medium text-foreground flex items-center space-x-1">
                      <Route className="h-5 w-5" />
                      <span>เส้นทางการบิน</span>
                    </Label>
                    <div className="flex items-center space-x-2">
                   
                      <Input
                        placeholder="BKK-HKG 12:15-16:45"
                        className="bg-background"
                        {...form.register(`sector_2.${index}.route`)}
                      />
                    </div>
                  </div>
                </div>
              </div>
            ))}

            <Button
              type="button"
              variant="outline"
              onClick={() => appendSector2({ date: '', flight: '', route: '' })}
              className="w-full border-dashed border-2 hover:border-chart-1/50 hover:bg-chart-1/5"
            >
              <Plus className="h-4 w-4 mr-2" />
              เพิ่มเที่ยวบิน
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}