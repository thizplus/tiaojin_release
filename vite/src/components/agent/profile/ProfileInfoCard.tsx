import { memo, useRef } from "react";
import { Building2, Edit, FileText, ExternalLink, Upload, Trash2, ImageIcon, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { getStorageUrl } from "@/lib/utils";
import type { AgentProfile } from "@/types";

interface ProfileInfoCardProps {
  profile: AgentProfile;
  onEdit: () => void;
  onUploadPdfHeader?: (file: File) => void;
  onDeletePdfHeader?: () => void;
  isUploadingPdfHeader?: boolean;
  className?: string;
}

export const ProfileInfoCard = memo(function ProfileInfoCard({
  profile,
  onEdit,
  onUploadPdfHeader,
  onDeletePdfHeader,
  isUploadingPdfHeader = false,
  className = ""
}: ProfileInfoCardProps) {
  const pdfHeaderInputRef = useRef<HTMLInputElement>(null);
  // Embedded labels to remove dependency on constants
  const labels = {
    tax_id: 'เลขประจำตัวผู้เสียภาษี',
    phone: 'เบอร์โทรศัพท์',
    province: 'จังหวัด',
    address: 'ที่อยู่',
    atta: 'ใบอนุญาต ATTA',
    dbd: 'ทะเบียนพาณิชย์อิเล็กทรอนิกส์ (DBD)',
    pp20: 'ใบทะเบียนภาษีมูลค่าเพิ่ม (PP20)'
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Mobile-First Header */}
      <div className="bg-card border border-border rounded-lg p-4">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex items-start gap-3">
            {/* Logo */}
            {profile.logo ? (
              <img
                src={profile.logo}
                alt="โลโก้บริษัท"
                className="w-24 h-24 sm:w-22 sm:h-22 rounded-lg object-cover border border-border flex-shrink-0"
              />
            ) : (
              <div className="w-12 h-12 sm:w-16 sm:h-16 rounded-lg bg-muted border border-border flex items-center justify-center flex-shrink-0">
                <Building2 className="w-6 h-6 sm:w-8 sm:h-8 text-muted-foreground" />
              </div>
            )}

            {/* Company Info */}
            <div className="min-w-0 flex-1">
              <h2 className="text-lg sm:text-xl font-bold text-foreground truncate">
                {profile.name}
              </h2>
              <p className="text-sm text-muted-foreground">
                บริษัททัวร์
              </p>
              <div className="mt-2 inline-block">
                <span className="text-xs px-2 py-1 bg-muted text-muted-foreground rounded border border-border">
                  {labels.tax_id}: {profile.tax_id}
                </span>
              </div>
            </div>
          </div>

          {/* Edit Button */}
          <Button
            onClick={onEdit}
            variant="outline"
            size="sm"
            className="self-start sm:self-center"
          >
            <Edit className="w-4 h-4 mr-2" />
            แก้ไข
          </Button>
        </div>
      </div>

      {/* Contact Info */}
      <div className="bg-card border border-border rounded-lg p-4">
        <h3 className="text-base font-semibold text-foreground mb-3">
          ข้อมูลติดต่อ
        </h3>

        <div className="space-y-3">
          {profile.phone && (
            <div>
              <label className="text-sm text-muted-foreground">
                {labels.phone}
              </label>
              <p className="text-sm font-medium text-foreground">
                {profile.phone}
              </p>
            </div>
          )}

          {profile.province && (
            <div>
              <label className="text-sm text-muted-foreground">
                {labels.province}
              </label>
              <p className="text-sm font-medium text-foreground">
                {profile.province}
              </p>
            </div>
          )}

          {profile.address && (
            <div>
              <label className="text-sm text-muted-foreground">
                {labels.address}
              </label>
              <p className="text-sm font-medium text-foreground leading-relaxed">
                {profile.address}, {profile.province} {profile.zipcode}
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Documents */}
      {(profile.atta || profile.dbd || profile.pp20) && (
        <div className="bg-card border border-border rounded-lg p-4">
          <h3 className="text-base font-semibold text-foreground mb-3">
            เอกสารใบอนุญาต
          </h3>

          <div className="space-y-2">
            {profile.atta && (
              <div className="flex items-center justify-between p-3 bg-muted border border-border rounded-lg">
                <div className="flex items-center gap-2">
                  <FileText className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                  <p className="text-sm font-medium text-foreground truncate">
                    {labels.atta}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => window.open(profile.atta || '', '_blank')}
                  className="flex-shrink-0"
                >
                  <ExternalLink className="h-3 w-3" />
                </Button>
              </div>
            )}

            {profile.dbd && (
              <div className="flex items-center justify-between p-3 bg-muted border border-border rounded-lg">
                <div className="flex items-center gap-2">
                  <FileText className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                  <p className="text-sm font-medium text-foreground truncate">
                    {labels.dbd}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => window.open(profile.dbd || '', '_blank')}
                  className="flex-shrink-0"
                >
                  <ExternalLink className="h-3 w-3" />
                </Button>
              </div>
            )}

            {profile.pp20 && (
              <div className="flex items-center justify-between p-3 bg-muted border border-border rounded-lg">
                <div className="flex items-center gap-2">
                  <FileText className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                  <p className="text-sm font-medium text-foreground truncate">
                    {labels.pp20}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => window.open(profile.pp20 || '', '_blank')}
                  className="flex-shrink-0"
                >
                  <ExternalLink className="h-3 w-3" />
                </Button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* PDF Header Image */}
      <div className="bg-card border border-border rounded-lg p-4">
        <h3 className="text-base font-semibold text-foreground mb-3">
          ภาพหัว PDF
        </h3>
        <p className="text-sm text-muted-foreground mb-3">
          ภาพที่จะแสดงเป็นหัวกระดาษในหน้าปกโปรแกรมทัวร์ (แทนที่ header ที่สร้างอัตโนมัติ)
        </p>

        {profile.pdfHeaderImage ? (
          <div className="space-y-3">
            <img
              src={getStorageUrl(profile.pdfHeaderImage)}
              alt="ภาพหัว PDF"
              className="w-full max-h-48 object-contain rounded-lg border border-border bg-muted"
            />
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => pdfHeaderInputRef.current?.click()}
                disabled={isUploadingPdfHeader}
              >
                {isUploadingPdfHeader ? (
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                ) : (
                  <Upload className="w-4 h-4 mr-2" />
                )}
                เปลี่ยนภาพ
              </Button>
              <Button
                variant="destructive"
                size="sm"
                onClick={onDeletePdfHeader}
                disabled={isUploadingPdfHeader}
              >
                <Trash2 className="w-4 h-4 mr-2" />
                ลบภาพ
              </Button>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center gap-3 p-6 border-2 border-dashed border-border rounded-lg">
            <ImageIcon className="w-10 h-10 text-muted-foreground" />
            <p className="text-sm text-muted-foreground">ยังไม่มีภาพหัว PDF</p>
            <Button
              variant="outline"
              size="sm"
              onClick={() => pdfHeaderInputRef.current?.click()}
              disabled={isUploadingPdfHeader}
            >
              {isUploadingPdfHeader ? (
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              ) : (
                <Upload className="w-4 h-4 mr-2" />
              )}
              อัพโหลดภาพ
            </Button>
          </div>
        )}

        {/* Hidden file input */}
        <input
          ref={pdfHeaderInputRef}
          type="file"
          accept="image/*"
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file && onUploadPdfHeader) {
              onUploadPdfHeader(file);
            }
            // Reset input value เพื่อให้เลือกไฟล์เดิมซ้ำได้
            e.target.value = '';
          }}
        />
      </div>
    </div>
  );
});