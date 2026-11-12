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
import type { City, UpdateCityRequest } from "@/types";

interface EditCityDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  city: City | null;
  onSuccess?: () => void;
}

export const EditCityDialog = memo(function EditCityDialog({
  open,
  onOpenChange,
  city,
  onSuccess
}: EditCityDialogProps) {
  const [formData, setFormData] = useState<UpdateCityRequest>({
    title: '',
    detail: '',
    image: '',
    status: 'draft'
  });

  const { updateCity, isUpdateLoading, refetchCitiesList } = useCity();

  // Load city data into form when city changes or dialog opens
  useEffect(() => {
    if (city && open) {
      setFormData({
        title: city.title,
        detail: city.detail || '',
        image: city.image,
        status: city.status as 'publish' | 'draft'
      });
    }
  }, [city, open]);

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
  const isFormValid = formData.title?.trim() !== '' && formData.image?.trim() !== '';

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isFormValid || !city) {
      return;
    }

    try {
      await updateCity(city.id, formData);

      // Revalidate cities list to show updated data
      refetchCitiesList();

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error('Failed to update city:', error);
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

  // Don't render if no city is selected
  if (!city && open) {
    return null;
  }

  return (
    <Dialog open={open} onOpenChange={isUpdateLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>แก้ไขเมือง</DialogTitle>
          <DialogDescription>
            แก้ไขข้อมูลเมือง: {city?.title}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="title">ชื่อเมือง <span className="text-red-500">*</span></Label>
              <Input
                id="title"
                name="title"
                value={formData.title || ''}
                onChange={handleInputChange}
                placeholder="ชื่อเมือง"
                required
                autoFocus
                disabled={isUpdateLoading}
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
                value={formData.image || ''}
                onChange={(url) => {
                  setFormData(prev => ({
                    ...prev,
                    image: url || ''
                  }));
                }}
                variant="dropzone"
                size="md"
                placeholder="เลือกรูปภาพเมืองหรือลากมาวาง"
                disabled={isUpdateLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="detail">รายละเอียด</Label>
              <Textarea
                id="detail"
                name="detail"
                value={formData.detail || ''}
                onChange={handleInputChange}
                placeholder="รายละเอียดเกี่ยวกับเมือง"
                disabled={isUpdateLoading}
                rows={4}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="status">สถานะ</Label>
              <Select
                value={formData.status}
                onValueChange={handleSelectChange}
                disabled={isUpdateLoading}
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