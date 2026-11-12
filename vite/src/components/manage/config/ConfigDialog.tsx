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
import { Switch } from '@/components/ui/switch';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import {
  Save,
 // X,
  Code,
  Loader2,
  AlertTriangle
} from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { useCreateConfig, useUpdateConfig } from '@/stores/configStore';
import type { SystemConfig, ConfigType, CreateConfigRequest, UpdateConfigRequest } from '@/types/config';

const CONFIG_TYPES: { value: ConfigType; label: string }[] = [
  { value: 'expiration', label: 'การหมดอายุ' },
  { value: 'payment', label: 'การชำระเงิน' },
  { value: 'notification', label: 'การแจ้งเตือน' },
  { value: 'security', label: 'ความปลอดภัย' },
  { value: 'integration', label: 'การเชื่อมต่อ' },
  { value: 'general', label: 'ทั่วไป' },
];

interface ConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  mode: 'create' | 'edit' | 'view';
  config?: SystemConfig | null;
}

export function ConfigDialog({
  open,
  onOpenChange,
  mode,
  config
}: ConfigDialogProps) {
  const { toast } = useToast();
  const createMutation = useCreateConfig();
  const updateMutation = useUpdateConfig();

  const [formData, setFormData] = useState({
    config_key: '',
    config_type: 'general' as ConfigType,
    description: '',
    config_value: '',
    is_active: true,
  });

  const [jsonError, setJsonError] = useState<string | null>(null);

  useEffect(() => {
    if (config) {
      setFormData({
        config_key: config.config_key,
        config_type: config.config_type as ConfigType,
        description: config.description,
        config_value: JSON.stringify(config.config_value, null, 2),
        is_active: config.is_active,
      });
    } else {
      setFormData({
        config_key: '',
        config_type: 'general',
        description: '',
        config_value: '{}',
        is_active: true,
      });
    }
    setJsonError(null);
  }, [config, open]);

  const validateJson = (value: string) => {
    try {
      JSON.parse(value);
      setJsonError(null);
      return true;
    } catch (error) {
      setJsonError(error instanceof Error ? error.message : 'Invalid JSON');
      return false;
    }
  };

  const handleConfigValueChange = (value: string) => {
    setFormData({ ...formData, config_value: value });
    validateJson(value);
  };

  const handleSave = async () => {
    if (!formData.config_key.trim()) {
      toast({
        title: "กรุณากรอกข้อมูล",
        description: "กรุณากรอก Config Key",
        variant: "destructive",
      });
      return;
    }

    if (!formData.description.trim()) {
      toast({
        title: "กรุณากรอกข้อมูล",
        description: "กรุณากรอกคำอธิบาย",
        variant: "destructive",
      });
      return;
    }

    if (!validateJson(formData.config_value)) {
      toast({
        title: "ข้อมูล JSON ไม่ถูกต้อง",
        description: jsonError || "กรุณาตรวจสอบรูปแบบ JSON",
        variant: "destructive",
      });
      return;
    }

    try {
      const configValue = JSON.parse(formData.config_value);

      if (mode === 'create') {
        const request: CreateConfigRequest = {
          config_key: formData.config_key,
          config_type: formData.config_type,
          description: formData.description,
          config_value: configValue,
          is_active: formData.is_active,
        };

        await createMutation.mutateAsync(request);
        toast({
          title: "สร้างสำเร็จ",
          description: "เพิ่มการตั้งค่าใหม่แล้ว",
        });
      } else if (mode === 'edit' && config) {
        const request: UpdateConfigRequest = {
          description: formData.description,
          config_value: configValue,
          is_active: formData.is_active,
        };

        await updateMutation.mutateAsync({ id: config.id, data: request });
        toast({
          title: "อัพเดทสำเร็จ",
          description: "แก้ไขการตั้งค่าแล้ว",
        });
      }

      onOpenChange(false);
    } catch (error) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: error instanceof Error ? error.message : "ไม่สามารถบันทึกได้",
        variant: "destructive",
      });
    }
  };

  const isReadOnly = mode === 'view';
  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {mode === 'create' && 'เพิ่มการตั้งค่าใหม่'}
            {mode === 'edit' && 'แก้ไขการตั้งค่า'}
            {mode === 'view' && 'ดูการตั้งค่า'}
          </DialogTitle>
          <DialogDescription>
            {mode === 'create' && 'สร้างการตั้งค่าใหม่สำหรับระบบ'}
            {mode === 'edit' && 'แก้ไขการตั้งค่าที่มีอยู่'}
            {mode === 'view' && 'รายละเอียดการตั้งค่า'}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Basic Info */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="config_key">Config Key *</Label>
              <Input
                id="config_key"
                value={formData.config_key}
                onChange={(e) => setFormData({ ...formData, config_key: e.target.value })}
                placeholder="เช่น payment.stripe_settings"
                disabled={isReadOnly || mode === 'edit'}
                className="font-mono"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="config_type">ประเภท *</Label>
              <Select
                value={formData.config_type}
                onValueChange={(value) => setFormData({ ...formData, config_type: value as ConfigType })}
                disabled={isReadOnly}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CONFIG_TYPES.map((type) => (
                    <SelectItem key={type.value} value={type.value}>
                      {type.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">คำอธิบาย *</Label>
            <Textarea
              id="description"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="อธิบายการใช้งานและวัตถุประสงค์ของการตั้งค่านี้"
              rows={3}
              disabled={isReadOnly}
            />
          </div>

          <div className="flex items-center space-x-2">
            <Switch
              id="is_active"
              checked={formData.is_active}
              onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
              disabled={isReadOnly}
            />
            <Label htmlFor="is_active">เปิดใช้งาน</Label>
          </div>

          <Separator />

          {/* Config Value */}
          <div className="space-y-2">
            <div className="flex items-center space-x-2">
              <Code className="h-4 w-4" />
              <Label htmlFor="config_value">ค่าการตั้งค่า (JSON) *</Label>
            </div>
            <Textarea
              id="config_value"
              value={formData.config_value}
              onChange={(e) => handleConfigValueChange(e.target.value)}
              placeholder='{"example": "value"}'
              rows={12}
              className="font-mono text-sm"
              disabled={isReadOnly}
            />
            {jsonError && (
              <div className="flex items-center space-x-2 text-destructive text-sm">
                <AlertTriangle className="h-4 w-4" />
                <span>{jsonError}</span>
              </div>
            )}
          </div>

          {/* Preview */}
          {mode === 'view' && config && (
            <>
              <Separator />
              <Card>
                <CardContent className="p-4">
                  <div className="space-y-2 text-sm">
                    <div><strong>ID:</strong> {config.id}</div>
                    <div><strong>สร้างเมื่อ:</strong> {new Date(config.created_at).toLocaleString('th-TH')}</div>
                    <div><strong>อัพเดทเมื่อ:</strong> {new Date(config.updated_at).toLocaleString('th-TH')}</div>
                  </div>
                </CardContent>
              </Card>
            </>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isLoading}>
            {mode === 'view' ? 'ปิด' : 'ยกเลิก'}
          </Button>
          {!isReadOnly && (
            <Button
              onClick={handleSave}
              disabled={isLoading || !!jsonError}
            >
              {isLoading ? (
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
              ) : (
                <Save className="h-4 w-4 mr-2" />
              )}
              {mode === 'create' ? 'สร้าง' : 'บันทึก'}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}