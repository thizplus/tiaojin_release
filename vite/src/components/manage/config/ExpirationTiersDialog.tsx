import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Plus,
  Trash2,
  Check,
  X,
  Loader2,
  Play
} from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import {
  useUpdateExpirationTiers,
  useValidateExpirationTiers,
  useCalculateTier
} from '@/stores/configStore';
import type {
  ExpirationTiersData,
  ExpirationTier
} from '@/types/config';

interface ExpirationTiersDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialData?: ExpirationTiersData | null;
}

export function ExpirationTiersDialog({
  open,
  onOpenChange,
  initialData
}: ExpirationTiersDialogProps) {
  const { toast } = useToast();
  const updateMutation = useUpdateExpirationTiers();
  const validateMutation = useValidateExpirationTiers();
  const calculateMutation = useCalculateTier();

  const [tiers, setTiers] = useState<Omit<ExpirationTier, 'example_message'>[]>([]);
  const [validationResult, setValidationResult] = useState<{errors: string[]; warnings: string[]; is_valid: boolean} | null>(null);
  const [testDate, setTestDate] = useState('');
  const [calculationResult, setCalculationResult] = useState<{
    tier?: string;
    expiry_date?: string;
    time_remaining?: string;
    days_to_tour?: number;
    applicable_tier?: { name: string };
    expiration_time?: string;
    time_until_expiry?: string;
    days_until_tour?: number;
  } | null>(null);

  useEffect(() => {
    if (initialData?.tiers) {
      setTiers(initialData.tiers.map(tier => ({
        name: tier.name,
        description: tier.description,
        threshold_days: tier.threshold_days,
        expire_minutes: tier.expire_minutes,
        expire_hours: tier.expire_hours,
        expire_days: tier.expire_days,
        sort_order: tier.sort_order,
      })));
    }
  }, [initialData]);

  const handleAddTier = () => {
    const newSortOrder = Math.max(...tiers.map(t => t.sort_order), 0) + 1;
    setTiers([...tiers, {
      name: '',
      description: '',
      threshold_days: null,
      expire_minutes: null,
      expire_hours: null,
      expire_days: 1,
      sort_order: newSortOrder,
    }]);
  };

  const handleRemoveTier = (index: number) => {
    setTiers(tiers.filter((_, i) => i !== index));
  };

  const handleTierChange = (index: number, field: string, value: string | number | null) => {
    const updatedTiers = [...tiers];
    updatedTiers[index] = {
      ...updatedTiers[index],
      [field]: value === '' ? null : value,
    };

    // If changing expiration type, clear the other fields
    if (field === 'expire_minutes') {
      updatedTiers[index].expire_hours = null;
      updatedTiers[index].expire_days = null;
    } else if (field === 'expire_hours') {
      updatedTiers[index].expire_minutes = null;
      updatedTiers[index].expire_days = null;
    } else if (field === 'expire_days') {
      updatedTiers[index].expire_minutes = null;
      updatedTiers[index].expire_hours = null;
    }

    setTiers(updatedTiers);
  };

  const handleValidate = async () => {
    try {
      const result = await validateMutation.mutateAsync({ tiers });
      setValidationResult(result);

      if (result.is_valid) {
        toast({
          title: "การตรวจสอบสำเร็จ",
          description: "การตั้งค่าถูกต้อง",
        });
      } else {
        toast({
          title: "พบข้อผิดพลาด",
          description: `${result.errors.length} ข้อผิดพลาด, ${result.warnings.length} คำเตือน`,
          variant: "destructive",
        });
      }
    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถตรวจสอบได้",
        variant: "destructive",
      });
    }
  };

  const handleTestCalculation = async () => {
    if (!testDate) {
      toast({
        title: "กรุณาเลือกวันที่",
        description: "เลือกวันที่เริ่มทัวร์เพื่อทดสอบ",
        variant: "destructive",
      });
      return;
    }

    try {
      const result = await calculateMutation.mutateAsync({
        tour_start_date: testDate,
        created_at: new Date().toISOString(),
      });
      setCalculationResult(result);
    } catch {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: "ไม่สามารถทดสอบการคำนวณได้",
        variant: "destructive",
      });
    }
  };

  const handleSave = async () => {
    try {
      await updateMutation.mutateAsync({ tiers });
      toast({
        title: "บันทึกสำเร็จ",
        description: "อัพเดทการตั้งค่าการหมดอายุแล้ว",
      });
      onOpenChange(false);
    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถบันทึกได้",
        variant: "destructive",
      });
    }
  };

  const getTierColor = (sortOrder: number) => {
    switch (sortOrder) {
      case 0: return 'destructive'; // Most urgent (นาที)
      case 1: return 'destructive'; // Very urgent (ชั่วโมง)
      case 2: return 'default';     // Normal (วัน)
      case 3: return 'default';     // Normal (วัน)
      case 4: return 'secondary';   // Default (วัน)
      default: return 'secondary';
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>แก้ไขการตั้งค่าวันหมดอายุการจอง</DialogTitle>
          <DialogDescription>
            กำหนดระยะเวลาหมดอายุการจองตามประเภททัวร์และระยะห่างจากวันเริ่มทัวร์
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Tier Configuration */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h4 className="font-medium">ระดับการหมดอายุ</h4>
              <Button variant="outline" onClick={handleAddTier}>
                <Plus className="h-4 w-4 mr-2" />
                เพิ่มระดับ
              </Button>
            </div>

            {tiers.map((tier, index) => (
              <Card key={index}>
                <CardContent className="p-4 space-y-4">
                  <div className="flex items-center justify-between">
                    <Badge variant={getTierColor(tier.sort_order)}>
                      ระดับ {tier.sort_order}
                    </Badge>
                    {tiers.length > 1 && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleRemoveTier(index)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>ชื่อระดับ</Label>
                      <Input
                        value={tier.name}
                        onChange={(e) => handleTierChange(index, 'name', e.target.value)}
                        placeholder="เช่น urgent, normal"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>ลำดับ</Label>
                      <Input
                        type="number"
                        value={tier.sort_order}
                        onChange={(e) => handleTierChange(index, 'sort_order', parseInt(e.target.value))}
                        min="1"
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label>คำอธิบาย</Label>
                    <Textarea
                      value={tier.description}
                      onChange={(e) => handleTierChange(index, 'description', e.target.value)}
                      placeholder="อธิบายเงื่อนไขของระดับนี้"
                      rows={2}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4 mb-4">
                    <div className="space-y-2">
                      <Label>ระยะห่างจากทัวร์ (วัน)</Label>
                      <Input
                        type="number"
                        value={tier.threshold_days || ''}
                        onChange={(e) => handleTierChange(index, 'threshold_days', e.target.value ? parseInt(e.target.value) : null)}
                        placeholder="เช่น 7, 30 (ว่างสำหรับ default)"
                        min="1"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <Label>หมดอายุ (นาที)</Label>
                      <Input
                        type="number"
                        value={tier.expire_minutes || ''}
                        onChange={(e) => handleTierChange(index, 'expire_minutes', e.target.value ? parseInt(e.target.value) : null)}
                        placeholder="เช่น 1, 5"
                        min="1"
                        disabled={!!tier.expire_hours || !!tier.expire_days}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>หมดอายุ (ชั่วโมง)</Label>
                      <Input
                        type="number"
                        value={tier.expire_hours || ''}
                        onChange={(e) => handleTierChange(index, 'expire_hours', e.target.value ? parseInt(e.target.value) : null)}
                        placeholder="เช่น 2, 6"
                        min="1"
                        disabled={!!tier.expire_minutes || !!tier.expire_days}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>หมดอายุ (วัน)</Label>
                      <Input
                        type="number"
                        value={tier.expire_days || ''}
                        onChange={(e) => handleTierChange(index, 'expire_days', e.target.value ? parseInt(e.target.value) : null)}
                        placeholder="เช่น 1, 3"
                        min="1"
                        disabled={!!tier.expire_minutes || !!tier.expire_hours}
                      />
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          <Separator />

          {/* Validation */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h4 className="font-medium">ตรวจสอบการตั้งค่า</h4>
              <Button
                variant="outline"
                onClick={handleValidate}
                disabled={validateMutation.isPending}
              >
                {validateMutation.isPending ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Check className="h-4 w-4" />
                )}
                ตรวจสอบ
              </Button>
            </div>

            {validationResult && (
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center mb-2">
                    {validationResult.is_valid ? (
                      <Check className="h-5 w-5 text-green-600 mr-2" />
                    ) : (
                      <X className="h-5 w-5 text-red-600 mr-2" />
                    )}
                    <span className={validationResult.is_valid ? 'text-green-600' : 'text-red-600'}>
                      {validationResult.is_valid ? 'การตั้งค่าถูกต้อง' : 'พบข้อผิดพลาด'}
                    </span>
                  </div>

                  {validationResult.errors?.length > 0 && (
                    <div className="mb-2">
                      <p className="text-sm font-medium text-red-600 mb-1">ข้อผิดพลาด:</p>
                      {validationResult.errors.map((error: string, i: number) => (
                        <p key={i} className="text-sm text-red-600">• {error}</p>
                      ))}
                    </div>
                  )}

                  {validationResult.warnings?.length > 0 && (
                    <div>
                      <p className="text-sm font-medium text-yellow-600 mb-1">คำเตือน:</p>
                      {validationResult.warnings.map((warning: string, i: number) => (
                        <p key={i} className="text-sm text-yellow-600">• {warning}</p>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            )}
          </div>

          <Separator />

          {/* Test Calculation */}
          <div className="space-y-4">
            <h4 className="font-medium">ทดสอบการคำนวณ</h4>
            <div className="flex gap-4">
              <div className="flex-1">
                <Label>วันที่เริ่มทัวร์</Label>
                <Input
                  type="datetime-local"
                  value={testDate}
                  onChange={(e) => setTestDate(e.target.value)}
                />
              </div>
              <div className="flex items-end">
                <Button
                  variant="outline"
                  onClick={handleTestCalculation}
                  disabled={calculateMutation.isPending || !testDate}
                >
                  {calculateMutation.isPending ? (
                    <Loader2 className="h-4 w-4 animate-spin mr-2" />
                  ) : (
                    <Play className="h-4 w-4 mr-2" />
                  )}
                  ทดสอบ
                </Button>
              </div>
            </div>

            {calculationResult && (
              <Card>
                <CardContent className="p-4">
                  <div className="space-y-2">
                    <p><strong>ระดับที่ใช้:</strong> {calculationResult.applicable_tier?.name || calculationResult.tier}</p>
                    <p><strong>วันที่หมดอายุ:</strong> {(calculationResult.expiration_time || calculationResult.expiry_date) ? new Date(calculationResult.expiration_time || calculationResult.expiry_date!).toLocaleString('th-TH') : 'ไม่ระบุ'}</p>
                    <p><strong>เวลาที่เหลือ:</strong> {calculationResult.time_until_expiry || calculationResult.time_remaining}</p>
                    <p><strong>วันถึงทัวร์:</strong> {calculationResult.days_until_tour || calculationResult.days_to_tour} วัน</p>
                  </div>
                </CardContent>
              </Card>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            ยกเลิก
          </Button>
          <Button
            onClick={handleSave}
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin mr-2" />
            ) : (
              <Check className="h-4 w-4 mr-2" />
            )}
            บันทึก
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}