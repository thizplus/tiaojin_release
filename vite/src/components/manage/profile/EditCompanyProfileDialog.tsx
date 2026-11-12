import { memo, useEffect, useState } from "react";
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
import { Loader2 } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import companyProfileService from "@/services/admin/companyProfileService";
import type { CompanyProfile } from "@/types";

interface EditCompanyProfileDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  profile: CompanyProfile | null;
  onSuccess?: () => void;
}

export const EditCompanyProfileDialog = memo(function EditCompanyProfileDialog({
  open,
  onOpenChange,
  profile,
  onSuccess,
}: EditCompanyProfileDialogProps) {
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    tax_id: "",
    phone: "",
    address: "",
    province: "",
    zipcode: "",
  });

  useEffect(() => {
    if (profile && open) {
      setFormData({
        name: profile.name || "",
        tax_id: profile.tax_id || "",
        phone: profile.phone || "",
        address: profile.address || "",
        province: profile.province || "",
        zipcode: profile.zipcode || "",
      });
    }
  }, [profile, open]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!profile) return;

    setIsLoading(true);
    try {
      await companyProfileService.updateCompanyProfile(profile.id, formData);
      toast({ title: "แก้ไขโปรไฟล์สำเร็จ" });
      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      toast({
        title: "เกิดข้อผิดพลาด",
        description: err instanceof Error ? err.message : "ไม่สามารถแก้ไขโปรไฟล์ได้",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  if (!profile) return null;

  return (
    <Dialog open={open} onOpenChange={isLoading ? undefined : onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>แก้ไขโปรไฟล์บริษัท</DialogTitle>
          <DialogDescription>{profile.name}</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">ชื่อบริษัท</Label>
              <Input id="name" name="name" value={formData.name} onChange={handleChange} disabled={isLoading} />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="tax_id">เลขประจำตัวผู้เสียภาษี</Label>
              <Input id="tax_id" name="tax_id" value={formData.tax_id} onChange={handleChange} maxLength={13} disabled={isLoading} />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="phone">เบอร์โทร</Label>
              <Input id="phone" name="phone" value={formData.phone} onChange={handleChange} disabled={isLoading} />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="address">ที่อยู่</Label>
              <Input id="address" name="address" value={formData.address} onChange={handleChange} disabled={isLoading} />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="province">จังหวัด</Label>
                <Input id="province" name="province" value={formData.province} onChange={handleChange} disabled={isLoading} />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="zipcode">รหัสไปรษณีย์</Label>
                <Input id="zipcode" name="zipcode" value={formData.zipcode} onChange={handleChange} maxLength={5} disabled={isLoading} />
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isLoading}>
              ยกเลิก
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              บันทึก
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
});
