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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useCity } from "@/hooks/useCity";
import { Loader2 } from "lucide-react";
import { FileUpload } from "@/components/ui/file-upload";
import type { CreateCityRequest } from "@/types";

interface CreateCityDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export const CreateCityDialog = memo(function CreateCityDialog({
  open,
  onOpenChange,
  onSuccess
}: CreateCityDialogProps) {
  const [formData, setFormData] = useState<CreateCityRequest>({
    title: '',
    detail: '',
    image: '',
    status: 'publish'
  });

  const { createCity, isCreateLoading, refetchCitiesList } = useCity();

  // Reset form when dialog closes
  useEffect(() => {
    if (!open) {
      setFormData({
        title: '',
        detail: '',
        image: '',
        status: 'publish'
      });
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

  // Handle select changes
  const handleSelectChange = (value: string) => {
    setFormData(prev => ({
      ...prev,
      status: value as 'publish' | 'draft'
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
      await createCity(formData);

      // Revalidate cities list to show new data
      refetchCitiesList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to create city:', error);
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
          <DialogTitle>สร้างเมืองใหม่</DialogTitle>
          <DialogDescription>
            กรอกข้อมูลเพื่อสร้างเมืองใหม่ในระบบ
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="title">ชื่อเมือง <span className="text-red-500">*</span></Label>
              <Input
                id="title"
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                placeholder="ชื่อเมือง"
                required
                autoFocus
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label>รูปภาพ <span className="text-red-500">*</span></Label>
              <FileUpload
                structuredOptions={{
                  category: 'cities',
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
                placeholder="เลือกรูปภาพเมืองหรือลากมาวาง"
                disabled={isCreateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="detail">รายละเอียด</Label>
              <Textarea
                id="detail"
                name="detail"
                value={formData.detail}
                onChange={handleInputChange}
                placeholder="รายละเอียดเกี่ยวกับเมือง"
                disabled={isCreateLoading}
                rows={4}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="status">สถานะ</Label>
              <Select
                value={formData.status}
                onValueChange={handleSelectChange}
                disabled={isCreateLoading}
              >
                <SelectTrigger>
                  <SelectValue placeholder="เลือกสถานะ" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="publish">เผยแพร่</SelectItem>
                  <SelectItem value="draft">ฉบับร่าง</SelectItem>
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