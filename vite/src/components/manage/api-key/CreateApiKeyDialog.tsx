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
import { useApiKey } from "@/hooks/useApiKey";
import { Loader2, Copy, Eye, EyeOff, CheckCircle2 } from "lucide-react";
import type { CreateApiKeyRequest } from "@/types/request/apiKey";
import { toast } from "sonner";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { cn } from "@/lib/utils";

interface CreateApiKeyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export const CreateApiKeyDialog = memo(function CreateApiKeyDialog({
  open,
  onOpenChange,
  onSuccess
}: CreateApiKeyDialogProps) {
  const [formData, setFormData] = useState<CreateApiKeyRequest>({
    name: '',
    description: '',
    expires_at: null,
    ip_whitelist: [],
    permissions: []
  });

  const [ipInput, setIpInput] = useState('');
  const [permissionInput, setPermissionInput] = useState('');
  const [createdApiKey, setCreatedApiKey] = useState<string | null>(null);
  const [showKey, setShowKey] = useState(true);
  const [copied, setCopied] = useState(false);

  const { createApiKey, isCreateLoading, refetchApiKeysList } = useApiKey();

  // Reset form when dialog closes
  useEffect(() => {
    if (!open) {
      // Small delay to allow dialog close animation
      const timer = setTimeout(() => {
        setFormData({
          name: '',
          description: '',
          expires_at: null,
          ip_whitelist: [],
          permissions: []
        });
        setIpInput('');
        setPermissionInput('');
        setCreatedApiKey(null);
        setShowKey(true);
        setCopied(false);
      }, 200);
      return () => clearTimeout(timer);
    }
  }, [open]);

  // Handle form input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
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

  // Handle copy API key
  const handleCopyApiKey = () => {
    if (createdApiKey) {
      navigator.clipboard.writeText(createdApiKey);
      setCopied(true);
      toast.success("คัดลอก API Key แล้ว");
      setTimeout(() => setCopied(false), 2000);
    }
  };

  // Form validation
  const isFormValid = formData.name.trim() !== '';

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isFormValid) {
      return;
    }

    try {
      const result = await createApiKey(formData);
      console.log('Create API Key Result:', result);

      // Revalidate API keys list to show new data
      refetchApiKeysList();

      // Show the full API key (only shown once)
      if (result?.data?.api_key) {
        console.log('Found api_key:', result.data.api_key);
        setCreatedApiKey(result.data.api_key);
      } else if (result?.api_key) {
        console.log('Found api_key at root:', result.api_key);
        setCreatedApiKey(result.api_key);
      } else {
        console.log('No api_key found in result');
      }
    } catch (error) {
      console.error('Failed to create API key:', error);
    }
  };

  // Handle final close (after showing API key)
  const handleFinalClose = () => {
    console.log('handleFinalClose called');
    onOpenChange(false);
    onSuccess?.();
  };

  // Close dialog on ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && !isCreateLoading && !createdApiKey) {
        onOpenChange(false);
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [onOpenChange, isCreateLoading, createdApiKey]);

  // Show created API key screen
  if (createdApiKey) {
    return (
      <Dialog open={open} onOpenChange={handleFinalClose}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <CheckCircle2 className="h-5 w-5 text-green-500" />
              สร้าง API Key สำเร็จ
            </DialogTitle>
            <DialogDescription>
              คัดลอก API Key นี้และเก็บไว้ในที่ปลอดภัย คุณจะไม่สามารถดู API Key นี้อีกครั้ง
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="rounded-lg bg-muted p-4 space-y-2">
              <div className="flex items-center justify-between">
                <Label>API Key ของคุณ</Label>
                <div className="flex gap-2">
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => setShowKey(!showKey)}
                  >
                    {showKey ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={handleCopyApiKey}
                  >
                    {copied ? <CheckCircle2 className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
                  </Button>
                </div>
              </div>
              <code className="block rounded bg-background px-3 py-2 text-sm font-mono break-all">
                {showKey ? createdApiKey : '••••••••••••••••••••••••••••••••'}
              </code>
            </div>

            <div className="rounded-lg bg-yellow-50 dark:bg-yellow-950/20 p-4 text-sm text-yellow-800 dark:text-yellow-200">
              <strong>⚠️ คำเตือน:</strong> นี่เป็นครั้งเดียวที่คุณจะเห็น API Key นี้
              กรุณาคัดลอกและเก็บไว้ในที่ปลอดภัยก่อนปิดหน้าต่างนี้
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              onClick={handleFinalClose}
            >
              เสร็จสิ้น
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog open={open} onOpenChange={isCreateLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>สร้าง API Key ใหม่</DialogTitle>
          <DialogDescription>
            กรอกข้อมูลเพื่อสร้าง API Key ใหม่สำหรับเชื่อมต่อกับระบบ
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">ชื่อ API Key <span className="text-red-500">*</span></Label>
              <Input
                id="name"
                name="name"
                value={formData.name}
                onChange={handleInputChange}
                placeholder="ชื่อ API Key"
                required
                autoFocus
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">คำอธิบาย</Label>
              <Textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                placeholder="คำอธิบายเกี่ยวกับ API Key นี้"
                disabled={isCreateLoading}
                rows={3}
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
                    disabled={isCreateLoading}
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
                  disabled={isCreateLoading}
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
                  disabled={isCreateLoading}
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
                        disabled={isCreateLoading}
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
                  disabled={isCreateLoading}
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
                  disabled={isCreateLoading}
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
                        disabled={isCreateLoading}
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
              disabled={isCreateLoading}
            >
              ยกเลิก
            </Button>
            <Button
              type="submit"
              disabled={!isFormValid || isCreateLoading}
            >
              {isCreateLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              สร้าง
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
});