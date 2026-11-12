import React, { memo, useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { useApiKey } from "@/hooks/useApiKey";
import { Loader2 } from "lucide-react";
import type { ApiKey } from "@/types/models/apiKey";
import type { UpdateApiKeyRequest } from "@/types/request/apiKey";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { cn } from "@/lib/utils";

interface EditApiKeyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  apiKey: ApiKey | null;
  onSuccess?: () => void;
}

export const EditApiKeyDialog = memo(function EditApiKeyDialog({
  open,
  onOpenChange,
  apiKey,
  onSuccess
}: EditApiKeyDialogProps) {
  const [formData, setFormData] = useState<UpdateApiKeyRequest>({
    name: '',
    description: '',
    is_active: true,
    expires_at: null,
    ip_whitelist: [],
    permissions: []
  });

  const [ipInput, setIpInput] = useState('');
  const [permissionInput, setPermissionInput] = useState('');

  const { updateApiKey, isUpdateLoading, refetchApiKeysList } = useApiKey();

  // Load API key data into form when apiKey changes or dialog opens
  useEffect(() => {
    if (apiKey && open) {
      setFormData({
        name: apiKey.name,
        description: apiKey.description || '',
        is_active: apiKey.is_active,
        expires_at: apiKey.expires_at,
        ip_whitelist: apiKey.ip_whitelist || [],
        permissions: apiKey.permissions || []
      });
    }
  }, [apiKey, open]);

  // Handle form input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Handle switch change
  const handleSwitchChange = (checked: boolean) => {
    setFormData(prev => ({
      ...prev,
      is_active: checked
    }));
  };

  // Handle IP whitelist
  const handleAddIp = () => {
    if (ipInput.trim() && !formData.ip_whitelist?.includes(ipInput.trim())) {
      setFormData(prev => ({
        ...prev,
        ip_whitelist: [...(prev.ip_whitelist || []), ipInput.trim()]
      }));
      setIpInput('');
    }
  };

  const handleRemoveIp = (ip: string) => {
    setFormData(prev => ({
      ...prev,
      ip_whitelist: prev.ip_whitelist?.filter(i => i !== ip) || []
    }));
  };

  // Handle permissions
  const handleAddPermission = () => {
    if (permissionInput.trim() && !formData.permissions?.includes(permissionInput.trim())) {
      setFormData(prev => ({
        ...prev,
        permissions: [...(prev.permissions || []), permissionInput.trim()]
      }));
      setPermissionInput('');
    }
  };

  const handleRemovePermission = (permission: string) => {
    setFormData(prev => ({
      ...prev,
      permissions: prev.permissions?.filter(p => p !== permission) || []
    }));
  };

  // Form validation
  const isFormValid = formData.name?.trim() !== '';

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isFormValid || !apiKey) {
      return;
    }

    try {
      await updateApiKey(apiKey.id, formData);

      // Revalidate API keys list to show updated data
      refetchApiKeysList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to update API key:', error);
    }
  };

  // Close dialog on ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && !isUpdateLoading) {
        onOpenChange(false);
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [onOpenChange, isUpdateLoading]);

  // Don't render if no API key is selected
  if (!apiKey && open) {
    return null;
  }

  return (
    <Dialog open={open} onOpenChange={isUpdateLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>แก้ไข API Key</DialogTitle>
          <DialogDescription>
            แก้ไขข้อมูล API Key: {apiKey?.name}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">ชื่อ API Key <span className="text-red-500">*</span></Label>
              <Input
                id="name"
                name="name"
                value={formData.name || ''}
                onChange={handleInputChange}
                placeholder="ชื่อ API Key"
                required
                autoFocus
                disabled={isUpdateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">คำอธิบาย</Label>
              <Textarea
                id="description"
                name="description"
                value={formData.description || ''}
                onChange={handleInputChange}
                placeholder="คำอธิบายเกี่ยวกับ API Key นี้"
                disabled={isUpdateLoading}
                rows={3}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="is_active">สถานะการใช้งาน</Label>
                <div className="text-sm text-muted-foreground">
                  เปิด/ปิดการใช้งาน API Key นี้
                </div>
              </div>
              <Switch
                id="is_active"
                checked={formData.is_active}
                onCheckedChange={handleSwitchChange}
                disabled={isUpdateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label>วันหมดอายุ (ไม่บังคับ)</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className={cn(
                      "justify-start text-left font-normal",
                      !formData.expires_at && "text-muted-foreground"
                    )}
                    disabled={isUpdateLoading}
                  >
                    {formData.expires_at ? (
                      format(new Date(formData.expires_at), "d MMM yyyy", { locale: th })
                    ) : (
                      "เลือกวันหมดอายุ"
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={formData.expires_at ? new Date(formData.expires_at) : undefined}
                    onSelect={(date) => {
                      setFormData(prev => ({
                        ...prev,
                        expires_at: date ? date.toISOString() : null
                      }));
                    }}
                    disabled={(date) => date < new Date()}
                    initialFocus
                  />
                  {formData.expires_at && (
                    <div className="p-3 border-t">
                      <Button
                        variant="ghost"
                        className="w-full"
                        onClick={() => {
                          setFormData(prev => ({
                            ...prev,
                            expires_at: null
                          }));
                        }}
                      >
                        ล้างวันที่
                      </Button>
                    </div>
                  )}
                </PopoverContent>
              </Popover>
            </div>

            <div className="grid gap-2">
              <Label>IP Whitelist (ไม่บังคับ)</Label>
              <div className="flex gap-2">
                <Input
                  value={ipInput}
                  onChange={(e) => setIpInput(e.target.value)}
                  placeholder="เช่น 192.168.1.1"
                  disabled={isUpdateLoading}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      handleAddIp();
                    }
                  }}
                />
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleAddIp}
                  disabled={isUpdateLoading}
                >
                  เพิ่ม
                </Button>
              </div>
              {formData.ip_whitelist && formData.ip_whitelist.length > 0 && (
                <div className="flex flex-wrap gap-2 mt-2">
                  {formData.ip_whitelist.map((ip, index) => (
                    <div key={index} className="flex items-center gap-1 bg-muted px-2 py-1 rounded text-sm">
                      <span>{ip}</span>
                      <button
                        type="button"
                        onClick={() => handleRemoveIp(ip)}
                        className="text-muted-foreground hover:text-foreground"
                        disabled={isUpdateLoading}
                      >
                        ×
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="grid gap-2">
              <Label>สิทธิ์ (ไม่บังคับ)</Label>
              <div className="flex gap-2">
                <Input
                  value={permissionInput}
                  onChange={(e) => setPermissionInput(e.target.value)}
                  placeholder="เช่น read, write, admin"
                  disabled={isUpdateLoading}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      handleAddPermission();
                    }
                  }}
                />
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleAddPermission}
                  disabled={isUpdateLoading}
                >
                  เพิ่ม
                </Button>
              </div>
              {formData.permissions && formData.permissions.length > 0 && (
                <div className="flex flex-wrap gap-2 mt-2">
                  {formData.permissions.map((permission, index) => (
                    <div key={index} className="flex items-center gap-1 bg-muted px-2 py-1 rounded text-sm">
                      <span>{permission}</span>
                      <button
                        type="button"
                        onClick={() => handleRemovePermission(permission)}
                        className="text-muted-foreground hover:text-foreground"
                        disabled={isUpdateLoading}
                      >
                        ×
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isUpdateLoading}
            >
              ยกเลิก
            </Button>
            <Button
              type="submit"
              disabled={!isFormValid || isUpdateLoading}
            >
              {isUpdateLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              บันทึก
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
});