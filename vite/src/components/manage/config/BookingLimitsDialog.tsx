import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { AgentSelect } from '@/components/ui/agent-select';
import {
  Plus,
  X,
  Users,
  Settings,
  Save,
  Loader2,
  AlertCircle
} from 'lucide-react';
import { useUpdateBookingLimitsConfig } from '@/stores/configStore';
import { useToast } from '@/hooks/use-toast';
import type { BookingLimitsConfig, BookingLimitsRequest } from '@/types/config';

interface BookingLimitsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialData?: BookingLimitsConfig | null;
}

const DEFAULT_STATUSES = ['pending', 'processing', 'confirm'];

export function BookingLimitsDialog({ open, onOpenChange, initialData }: BookingLimitsDialogProps) {
  const { toast } = useToast();
  const updateMutation = useUpdateBookingLimitsConfig();

  const [formData, setFormData] = useState<BookingLimitsRequest>({
    agent_global_bypasses: [],
    global_settings: {
      count_all_statuses: DEFAULT_STATUSES,
      default_limit_per_period: 10,
      enabled: true,
    },
  });

  const [selectedAgentId, setSelectedAgentId] = useState<string | undefined>(undefined);
  const [newStatus, setNewStatus] = useState('');

  useEffect(() => {
    if (initialData) {
      setFormData({
        agent_global_bypasses: [...initialData.agent_global_bypasses],
        global_settings: {
          count_all_statuses: [...initialData.global_settings.count_all_statuses],
          default_limit_per_period: initialData.global_settings.default_limit_per_period,
          enabled: initialData.global_settings.enabled,
        },
      });
    } else {
      setFormData({
        agent_global_bypasses: [],
        global_settings: {
          count_all_statuses: DEFAULT_STATUSES,
          default_limit_per_period: 10,
          enabled: true,
        },
      });
    }
  }, [initialData, open]);

  const handleAddAgent = () => {
    if (!selectedAgentId) return;

    if (formData.agent_global_bypasses.includes(selectedAgentId)) {
      toast({
        title: "ข้อผิดพลาด",
        description: "Agent นี้มีอยู่แล้ว",
        variant: "destructive",
      });
      return;
    }

    setFormData(prev => ({
      ...prev,
      agent_global_bypasses: [...prev.agent_global_bypasses, selectedAgentId]
    }));
    setSelectedAgentId(undefined);
  };

  const handleRemoveAgent = (agentId: string) => {
    setFormData(prev => ({
      ...prev,
      agent_global_bypasses: prev.agent_global_bypasses.filter(id => id !== agentId)
    }));
  };

  const handleAddStatus = () => {
    if (!newStatus.trim()) return;

    if (formData.global_settings.count_all_statuses.includes(newStatus.trim())) {
      toast({
        title: "ข้อผิดพลาด",
        description: "สถานะนี้มีอยู่แล้ว",
        variant: "destructive",
      });
      return;
    }

    setFormData(prev => ({
      ...prev,
      global_settings: {
        ...prev.global_settings,
        count_all_statuses: [...prev.global_settings.count_all_statuses, newStatus.trim()]
      }
    }));
    setNewStatus('');
  };

  const handleRemoveStatus = (status: string) => {
    setFormData(prev => ({
      ...prev,
      global_settings: {
        ...prev.global_settings,
        count_all_statuses: prev.global_settings.count_all_statuses.filter(s => s !== status)
      }
    }));
  };

  const handleSave = async () => {
    try {
      await updateMutation.mutateAsync(formData);
      toast({
        title: "บันทึกสำเร็จ",
        description: "อัพเดตการตั้งค่าการจำกัดการจองแล้ว",
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

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center">
            <Settings className="h-5 w-5 mr-2" />
            การตั้งค่าการจำกัดการจอง
          </DialogTitle>
          <DialogDescription>
            กำหนดข้อจำกัดการจองและผู้ใช้ที่ได้รับการยกเว้น
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Enable/Disable */}
          <div className="flex items-center justify-between p-4 bg-muted rounded-lg">
            <div className="space-y-1">
              <Label className="text-sm font-medium">เปิด/ปิดระบบการจำกัดการจอง</Label>
              <p className="text-xs text-muted-foreground">
                เมื่อเปิดใช้งาน ระบบจะจำกัดจำนวนการจองตามที่กำหนด
              </p>
            </div>
            <Switch
              checked={formData.global_settings.enabled}
              onCheckedChange={(checked) =>
                setFormData(prev => ({
                  ...prev,
                  global_settings: {
                    ...prev.global_settings,
                    enabled: checked
                  }
                }))
              }
            />
          </div>

          <Separator />

          {/* Global Settings */}
          <div className="space-y-4">
            <h4 className="font-medium flex items-center">
              <Settings className="h-4 w-4 mr-2" />
              การตั้งค่าทั่วไป
            </h4>

            <div className="grid grid-cols-1 gap-4">
              <div className="space-y-2">
                <Label htmlFor="default_limit">จำนวนจำกัดต่อรอบ (ค่าเริ่มต้น)</Label>
                <Input
                  id="default_limit"
                  type="number"
                  min="1"
                  value={formData.global_settings.default_limit_per_period}
                  onChange={(e) =>
                    setFormData(prev => ({
                      ...prev,
                      global_settings: {
                        ...prev.global_settings,
                        default_limit_per_period: parseInt(e.target.value) || 1
                      }
                    }))
                  }
                />
                <p className="text-xs text-muted-foreground">
                  จำนวนการจองสูงสุดที่อนุญาตต่อรอบทัวร์
                </p>
              </div>

              <div className="space-y-2">
                <Label>สถานะการจองที่นับรวม</Label>
                <div className="flex flex-wrap gap-2 p-3 border rounded-lg min-h-[60px]">
                  {formData.global_settings.count_all_statuses.map((status) => (
                    <Badge key={status} variant="secondary" className="flex items-center gap-1">
                      {status}
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-4 w-4 p-0 ml-1"
                        onClick={() => handleRemoveStatus(status)}
                      >
                        <X className="h-3 w-3" />
                      </Button>
                    </Badge>
                  ))}
                </div>
                <div className="flex gap-2">
                  <Input
                    placeholder="เพิ่มสถานะใหม่..."
                    value={newStatus}
                    onChange={(e) => setNewStatus(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && handleAddStatus()}
                  />
                  <Button onClick={handleAddStatus} size="sm">
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  สถานะการจองที่จะนับรวมในการคำนวณข้อจำกัด
                </p>
              </div>
            </div>
          </div>

          <Separator />

          {/* Agent Bypasses */}
          <div className="space-y-4">
            <h4 className="font-medium flex items-center">
              <Users className="h-4 w-4 mr-2" />
              ผู้ใช้ที่ได้รับการยกเว้น
              <Badge variant="secondary" className="ml-2 text-xs">
                {formData.agent_global_bypasses.length}
              </Badge>
            </h4>

            <div className="space-y-3">
              <div className="space-y-2 max-h-40 overflow-y-auto">
                {formData.agent_global_bypasses.length === 0 ? (
                  <div className="text-center py-4 text-muted-foreground border-2 border-dashed rounded-lg">
                    <Users className="h-6 w-6 mx-auto mb-2 opacity-50" />
                    <p className="text-sm">ไม่มีผู้ใช้ที่ได้รับการยกเว้น</p>
                  </div>
                ) : (
                  formData.agent_global_bypasses.map((agentId) => (
                    <div
                      key={agentId}
                      className="flex items-center justify-between p-2 bg-muted rounded border"
                    >
                      <div className="flex items-center space-x-2">
                        <Users className="h-4 w-4 text-muted-foreground" />
                        <span className="font-mono text-sm">{agentId}</span>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRemoveAgent(agentId)}
                        className="h-6 w-6 p-0 text-destructive hover:text-destructive"
                      >
                        <X className="h-4 w-4" />
                      </Button>
                    </div>
                  ))
                )}
              </div>

              <div className="flex gap-2">
                <div className="flex-1">
                  <AgentSelect
                    value={selectedAgentId}
                    onValueChange={setSelectedAgentId}
                    placeholder="เลือก Agent..."
                  />
                </div>
                <Button
                  onClick={handleAddAgent}
                  size="sm"
                  disabled={!selectedAgentId}
                >
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              <p className="text-xs text-muted-foreground">
                ผู้ใช้เหล่านี้จะไม่ถูกจำกัดจำนวนการจอง
              </p>
            </div>
          </div>

          {/* Warning */}
          {!formData.global_settings.enabled && (
            <div className="flex items-start space-x-2 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
              <AlertCircle className="h-4 w-4 text-yellow-600 mt-0.5" />
              <div className="text-sm text-yellow-800">
                <p className="font-medium">ระบบการจำกัดการจองถูกปิดใช้งาน</p>
                <p>การตั้งค่าจะไม่มีผลจนกว่าจะเปิดใช้งานระบบ</p>
              </div>
            </div>
          )}
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={updateMutation.isPending}
          >
            ยกเลิก
          </Button>
          <Button
            onClick={handleSave}
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending && (
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
            )}
            <Save className="h-4 w-4 mr-2" />
            บันทึก
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}