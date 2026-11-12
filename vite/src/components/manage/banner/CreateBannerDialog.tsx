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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useBanner } from "@/hooks/useBanner";
import { Loader2 } from "lucide-react";
import { FileUpload } from "@/components/ui/file-upload";
import type { CreateBannerRequest } from "@/types";

interface CreateBannerDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export const CreateBannerDialog = memo(function CreateBannerDialog({
  open,
  onOpenChange,
  onSuccess
}: CreateBannerDialogProps) {
  const [formData, setFormData] = useState<CreateBannerRequest>({
    title: '',
    image: '',
    link_url: undefined,
    alt: undefined,
    display_order: 0,
    position: 'internal',
    status: 'publish'
  });

  const { createBanner, isCreateLoading, refetchBannersList } = useBanner();

  // Reset form when dialog closes
  useEffect(() => {
    if (!open) {
      setFormData({
        title: '',
        image: '',
        link_url: undefined,
        alt: undefined,
        display_order: 0,
        position: 'internal',
        status: 'publish'
      });
    }
  }, [open]);

  // Handle form input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'display_order' ? parseInt(value) || 0 : value
    }));
  };

  // Handle select changes
  const handleSelectChange = (name: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Form validation
  const isFormValid = formData.title.trim() !== '' && formData.image.trim() !== '';

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isFormValid) {
      return;
    }

    try {
      // Prepare data, remove empty optional fields
      const dataToSubmit: CreateBannerRequest = {
        ...formData,
        link_url: formData.link_url?.trim() || undefined,
        alt: formData.alt?.trim() || undefined,
      };

      await createBanner(dataToSubmit);

      // Revalidate banners list to show new data
      refetchBannersList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to create banner:', error);
    }
  };

  // Close dialog on ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && !isCreateLoading) {
        onOpenChange(false);
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [onOpenChange, isCreateLoading]);

  return (
    <Dialog open={open} onOpenChange={isCreateLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>สร้างแบนเนอร์ใหม่</DialogTitle>
          <DialogDescription>
            กรอกข้อมูลเพื่อสร้างแบนเนอร์ใหม่ในระบบ
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="title">ชื่อแบนเนอร์ <span className="text-red-500">*</span></Label>
              <Input
                id="title"
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                placeholder="ชื่อแบนเนอร์"
                required
                autoFocus
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label>รูปภาพ <span className="text-red-500">*</span></Label>
              <FileUpload
                structuredOptions={{
                  category: 'banners',
                  file_type: 'image'
                }}
                accept="image/*"
                maxSize={5 * 1024 * 1024} // 5MB
                allowedTypes={['image/jpeg', 'image/png', 'image/webp']}
                value={formData.image}
                onChange={(url) => {
                  setFormData(prev => ({
                    ...prev,
                    image: url || ''
                  }));
                }}
                variant="dropzone"
                size="md"
                placeholder="เลือกรูปภาพแบนเนอร์หรือลากมาวาง"
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="link_url">ลิงก์ (URL) <span className="text-muted-foreground text-xs font-normal">(ไม่บังคับ)</span></Label>
              <Input
                id="link_url"
                name="link_url"
                type="url"
                value={formData.link_url || ''}
                onChange={handleInputChange}
                placeholder="https://example.com"
                disabled={isCreateLoading}
              />
              <p className="text-xs text-muted-foreground">
                URL ที่จะเปิดเมื่อคลิกที่แบนเนอร์ (หากไม่ระบุจะไม่สามารถคลิกได้)
              </p>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="alt">Alt Text</Label>
              <Input
                id="alt"
                name="alt"
                value={formData.alt}
                onChange={handleInputChange}
                placeholder="คำอธิบายรูปภาพสำหรับผู้พิการทางสายตา"
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="display_order">ลำดับการแสดงผล</Label>
              <Input
                id="display_order"
                name="display_order"
                type="number"
                min="0"
                value={formData.display_order}
                onChange={handleInputChange}
                placeholder="0"
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="position">ตำแหน่งแสดงผล</Label>
              <Select
                value={formData.position || 'internal'}
                onValueChange={(v) => handleSelectChange('position', v)}
                disabled={isCreateLoading}
              >
                <SelectTrigger>
                  <SelectValue placeholder="เลือกตำแหน่ง" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="internal">ข้างใน (Dashboard / Tours)</SelectItem>
                  <SelectItem value="public">หน้าแรก (Public)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="status">สถานะ</Label>
              <Select
                value={formData.status}
                onValueChange={(v) => handleSelectChange('status', v)}
                disabled={isCreateLoading}
              >
                <SelectTrigger>
                  <SelectValue placeholder="เลือกสถานะ" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="publish">เผยแพร่</SelectItem>
                  <SelectItem value="draft">ร่าง</SelectItem>
                </SelectContent>
              </Select>
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