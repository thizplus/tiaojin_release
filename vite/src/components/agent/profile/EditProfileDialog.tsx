import { useState, useEffect, memo } from "react";
import { Loader2 } from "lucide-react";
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
import { FileUpload } from "@/components/ui/file-upload";
import type { UpdateAgentProfileRequest, AgentProfile } from "@/types";

const LABELS = {
  name: 'ชื่อบริษัท',
  phone: 'เบอร์โทรศัพท์',
  province: 'จังหวัด',
  address: 'ที่อยู่',
  zipcode: 'รหัสไปรษณีย์',
  atta: 'ใบอนุญาต ATTA',
  dbd: 'ทะเบียนพาณิชย์อิเล็กทรอนิกส์ (DBD)',
  pp20: 'ใบทะเบียนภาษีมูลค่าเพิ่ม (PP20)',
  logo: 'โลโก้บริษัท'
};

interface EditProfileDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  profile: AgentProfile | null;
  onSubmit: (data: UpdateAgentProfileRequest) => Promise<void>;
  isLoading?: boolean;
}

export const EditProfileDialog = memo(function EditProfileDialog({
  open,
  onOpenChange,
  profile,
  onSubmit,
  isLoading = false
}: EditProfileDialogProps) {
  const [formData, setFormData] = useState<UpdateAgentProfileRequest>({
    name: '',
    phone: '',
    address: '',
    province: '',
    zipcode: '',
    atta: '',
    dbd: '',
    pp20: '',
    logo: ''
  });

  // Load profile data into form when profile changes or dialog opens
  useEffect(() => {
    if (profile && open) {
      setFormData({
        name: profile.name,
        phone: profile.phone,
        address: profile.address,
        province: profile.province,
        zipcode: profile.zipcode,
        atta: profile.atta || '',
        dbd: profile.dbd || '',
        pp20: profile.pp20 || '',
        logo: profile.logo || ''
      });
    }
  }, [profile, open]);

  // Handle form input changes
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Form validation
  const isFormValid =
    formData.name?.trim() !== '' &&
    formData.phone?.trim() !== '' &&
    formData.address?.trim() !== '' &&
    formData.province?.trim() !== '' &&
    formData.zipcode?.trim() !== '' &&
    /^\d{5}$/.test(formData.zipcode || '');

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isFormValid || !profile) {
      return;
    }

    try {
      await onSubmit(formData);
      onOpenChange(false);
    } catch (error) {
      console.error('Failed to update profile:', error);
    }
  };

  // Close dialog on ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && !isLoading) {
        onOpenChange(false);
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [onOpenChange, isLoading]);

  // Don't render if no profile is selected
  if (!profile && open) {
    return null;
  }

  return (
    <Dialog open={open} onOpenChange={isLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>แก้ไขโปรไฟล์บริษัท</DialogTitle>
          <DialogDescription>
            อัพเดตข้อมูลบริษัททัวร์: {profile?.name}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">{LABELS.name} <span className="text-red-500">*</span></Label>
              <Input
                id="name"
                name="name"
                value={formData.name || ''}
                onChange={handleInputChange}
                placeholder="ชื่อบริษัท"
                required
                autoFocus
                disabled={isLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="phone">{LABELS.phone} <span className="text-red-500">*</span></Label>
              <Input
                id="phone"
                name="phone"
                value={formData.phone || ''}
                onChange={handleInputChange}
                placeholder="เบอร์โทรศัพท์"
                required
                disabled={isLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="address">{LABELS.address} <span className="text-red-500">*</span></Label>
              <Textarea
                id="address"
                name="address"
                value={formData.address || ''}
                onChange={handleInputChange}
                placeholder="ที่อยู่บริษัท"
                required
                disabled={isLoading}
                rows={3}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="province">{LABELS.province} <span className="text-red-500">*</span></Label>
                <Input
                  id="province"
                  name="province"
                  value={formData.province || ''}
                  onChange={handleInputChange}
                  placeholder="จังหวัด"
                  required
                  disabled={isLoading}
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="zipcode">{LABELS.zipcode} <span className="text-red-500">*</span></Label>
                <Input
                  id="zipcode"
                  name="zipcode"
                  value={formData.zipcode || ''}
                  onChange={handleInputChange}
                  placeholder="รหัสไปรษณีย์"
                  required
                  maxLength={5}
                  disabled={isLoading}
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label>{LABELS.atta}</Label>
              <FileUpload
                structuredOptions={{
                  category: 'profiles',
                  file_type: 'document'
                }}
                accept=".pdf,.jpg,.jpeg,.png"
                maxSize={5 * 1024 * 1024} // 5MB
                allowedTypes={['application/pdf', 'image/jpeg', 'image/png']}
                value={formData.atta || ''}
                onChange={(url) => {
                  setFormData(prev => ({
                    ...prev,
                    atta: url || ''
                  }));
                }}
                variant="dropzone"
                size="sm"
                placeholder="เลือกไฟล์ใบอนุญาต ATTA หรือลากมาวาง"
                disabled={isLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label>{LABELS.dbd}</Label>
              <FileUpload
                structuredOptions={{
                  category: 'profiles',
                  file_type: 'document'
                }}
                accept=".pdf,.jpg,.jpeg,.png"
                maxSize={5 * 1024 * 1024} // 5MB
                allowedTypes={['application/pdf', 'image/jpeg', 'image/png']}
                value={formData.dbd || ''}
                onChange={(url) => {
                  setFormData(prev => ({
                    ...prev,
                    dbd: url || ''
                  }));
                }}
                variant="dropzone"
                size="sm"
                placeholder="เลือกไฟล์ใบอนุญาต DBD หรือลากมาวาง"
                disabled={isLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label>{LABELS.pp20}</Label>
              <FileUpload
                structuredOptions={{
                  category: 'profiles',
                  file_type: 'document'
                }}
                accept=".pdf,.jpg,.jpeg,.png"
                maxSize={5 * 1024 * 1024} // 5MB
                allowedTypes={['application/pdf', 'image/jpeg', 'image/png']}
                value={formData.pp20 || ''}
                onChange={(url) => {
                  setFormData(prev => ({
                    ...prev,
                    pp20: url || ''
                  }));
                }}
                variant="dropzone"
                size="sm"
                placeholder="เลือกไฟล์ใบอนุญาต PP.20 หรือลากมาวาง"
                disabled={isLoading}
              />
            </div>

            <div className="grid gap-2">
              <Label>{LABELS.logo}</Label>
              <FileUpload
                structuredOptions={{
                  category: 'profiles',
                  file_type: 'image'
                }}
                accept="image/*"
                maxSize={2 * 1024 * 1024} // 2MB
                allowedTypes={['image/jpeg', 'image/png', 'image/webp']}
                value={formData.logo || ''}
                onChange={(url) => {
                  setFormData(prev => ({
                    ...prev,
                    logo: url || ''
                  }));
                }}
                variant="dropzone"
                size="md"
                placeholder="เลือกโลโก้บริษัทหรือลากมาวาง"
                disabled={isLoading}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              ยกเลิก
            </Button>
            <Button
              type="submit"
              disabled={!isFormValid || isLoading}
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              บันทึก
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
});