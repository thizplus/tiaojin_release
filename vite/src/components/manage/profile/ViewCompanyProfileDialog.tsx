import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Copy, ExternalLink, FileText, Link, Image as ImageIcon, File } from "lucide-react";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { useToast } from "@/hooks/use-toast";
import { getStorageUrl, isImageFile, isPdfFile } from "@/lib/utils";
import type { CompanyProfile } from "@/types";

// Helper functions
const copyToClipboard = async (text: string, label: string, toast: ReturnType<typeof useToast>['toast']) => {
  try {
    await navigator.clipboard.writeText(text);
    toast({
      title: "คัดลอกแล้ว",
      description: `คัดลอก ${label} เรียบร้อยแล้ว`,
    });
  } catch  {
    toast({
      title: "ไม่สามารถคัดลอกได้",
      description: "เกิดข้อผิดพลาดในการคัดลอกข้อมูล",
      variant: "destructive",
    });
  }
};

// Large Clickable Certificate Badge Component for Dialog
interface LargeCertificateBadgeProps {
  value: string;
  label: string;
}

const LargeCertificateBadge = ({ value, label }: LargeCertificateBadgeProps) => {
  const { toast } = useToast();

  // ตรวจสอบประเภทไฟล์
  const isImage = isImageFile(value);
  const isPdf = isPdfFile(value);
  const isFile = value && (isImage || isPdf || value.includes('/') || value.includes('.'));

  // แปลง path เป็น proxy URL สำหรับ preview
  const proxyUrl = isFile ? getStorageUrl(value) : value;

  // ฟังก์ชันเปิดไฟล์
  const handleOpenFile = () => {
    if (proxyUrl) {
      window.open(proxyUrl, '_blank', 'noopener,noreferrer');
    }
  };

  // ถ้าไม่มี value
  if (!value) {
    return (
      <div className="space-y-2">
        <label className="text-sm font-medium text-muted-foreground">{label}</label>
        <div className="flex items-center gap-3 p-3 border rounded-lg bg-muted/20">
          <FileText className="h-5 w-5 text-muted-foreground" />
          <span className="text-sm text-muted-foreground">ไม่มีข้อมูล</span>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <label className="text-sm font-medium text-muted-foreground">{label}</label>
      <div className="flex items-center gap-3 p-3 border rounded-lg bg-muted/20">
        {/* File/Link Icon */}
        <div className="flex-shrink-0">
          {isImage ? (
            <ImageIcon className="h-5 w-5 text-green-500" />
          ) : isPdf ? (
            <File className="h-5 w-5 text-red-500" />
          ) : isFile ? (
            <FileText className="h-5 w-5 text-blue-500" />
          ) : (
            <Link className="h-5 w-5 text-primary" />
          )}
        </div>

        {/* Preview สำหรับ Image */}
        {isImage && (
          <div className="flex-shrink-0">
            <img
              src={proxyUrl}
              alt={label}
              className="h-12 w-12 object-cover rounded border cursor-pointer hover:opacity-80"
              onClick={handleOpenFile}
            />
          </div>
        )}

        {/* Link/Text Display */}
        <div className="flex-1 min-w-0">
          <Button
            variant="link"
            className="h-auto p-0 text-primary hover:text-primary/80 justify-start text-left truncate w-full"
            onClick={handleOpenFile}
          >
            <span className="truncate">
              {isFile ? `ดู${label}` : value.replace(/^https?:\/\//, '')}
            </span>
          </Button>
        </div>

        {/* Action Buttons */}
        <div className="flex-shrink-0 flex items-center gap-1">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={() => copyToClipboard(value, label, toast)}
                >
                  <Copy className="h-3.5 w-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>คัดลอก {label}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={handleOpenFile}
                >
                  <ExternalLink className="h-3.5 w-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>เปิดดู{label}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </div>
    </div>
  );
};

interface ViewCompanyProfileDialogProps {
  profile: CompanyProfile | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function ViewCompanyProfileDialog({
  profile,
  open,
  onOpenChange,
}: ViewCompanyProfileDialogProps) {
  if (!profile) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[800px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>รายละเอียดโปรไฟล์บริษัท</DialogTitle>
          <DialogDescription>
            ข้อมูลโปรไฟล์บริษัทที่รอการอนุมัติ
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Company Logo */}
          {profile.logo && (
            <div className="flex justify-center">
              <img
                src={getStorageUrl(profile.logo)}
                alt="Company Logo"
                className="h-24 w-24 object-cover rounded-lg border cursor-pointer hover:opacity-80"
                onClick={() => {
                  const logoUrl = getStorageUrl(profile.logo);
                  if (logoUrl) window.open(logoUrl, '_blank');
                }}
              />
            </div>
          )}

          {/* Basic Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                ชื่อบริษัท
              </label>
              <p className="text-base font-semibold">{profile.name}</p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                เลขประจำตัวผู้เสียภาษี
              </label>
              <p className="font-mono text-base">
                {profile.tax_id.replace(/(\d{1})(\d{4})(\d{5})(\d{2})(\d{1})/, "$1-$2-$3-$4-$5")}
              </p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                เบอร์โทรศัพท์
              </label>
              <p className="font-mono text-base">{profile.phone}</p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                จังหวัด
              </label>
              <p className="text-base">{profile.province}</p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                รหัสไปรษณีย์
              </label>
              <p className="font-mono text-base">{profile.zipcode}</p>
            </div>
          </div>

          {/* Address */}
          <div className="space-y-2">
            <label className="text-sm font-medium text-muted-foreground">
              ที่อยู่
            </label>
            <p className="text-base">{profile.address}</p>
          </div>

          {/* Certificates */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">ใบอนุญาต</h3>
            <div className="grid grid-cols-1 gap-4">
              <LargeCertificateBadge value={profile.atta || ''} label="ATTA" />
              <LargeCertificateBadge value={profile.dbd || ''} label="DBD" />
              <LargeCertificateBadge value={profile.pp20 || ''} label="PP20" />
            </div>
          </div>

          {/* Timestamps */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t">
            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                วันที่ส่งคำขอ
              </label>
              <p className="text-base">
                {profile.created_at ? format(new Date(profile.created_at), "d MMMM yyyy เวลา HH:mm น.", { locale: th }) : '-'}
              </p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-muted-foreground">
                วันที่อัปเดตล่าสุด
              </label>
              <p className="text-base">
                {profile.updated_at ? format(new Date(profile.updated_at), "d MMMM yyyy เวลา HH:mm น.", { locale: th }) : '-'}
              </p>
            </div>
          </div>

          {/* User ID for reference */}
          <div className="space-y-2 pt-2 border-t">
            <label className="text-sm font-medium text-muted-foreground">
              User ID
            </label>
            <p className="text-sm text-muted-foreground">{profile.id}</p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}