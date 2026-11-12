import  { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { FileUpload } from '@/components/ui/file-upload';
import { Upload, FileText, User } from 'lucide-react';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';

import { Separator } from '@/components/ui/separator';
import { useToast } from '@/hooks/use-toast';
import { getStorageUrl } from '@/lib/utils';
import type { BookingResponse } from '@/types';

interface PassportUploadSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  booking: BookingResponse | null;
  onUpload: (id: string, file: File) => Promise<void>;
  isLoading?: boolean;
  onRefetch?: () => void; // callback เพื่อ reload booking data หลัง upload สำเร็จ
}

export function PassportUploadSheet({
  isOpen,
  onOpenChange,
  booking,
  onUpload,
  isLoading = false,
  onRefetch
}: PassportUploadSheetProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const { toast } = useToast();

  // ใช้ข้อมูล passports จาก booking โดยตรง (รองรับทุก role: Agent, Sales, Admin)
  const passports = booking?.passports || [];

  const handleFileSelect = (file: File | null) => {
    setSelectedFile(file);
  };

  const handleFileError = (error: Error) => {
    toast({
      title: 'เกิดข้อผิดพลาด',
      description: error.message,
      variant: 'destructive',
    });
  };


  const handleSubmit = async () => {
    if (!booking || !selectedFile) {
      toast({
        title: 'ข้อมูลไม่ครบ',
        description: 'กรุณาเลือกไฟล์หนังสือเดินทาง',
        variant: 'destructive',
      });
      return;
    }

    try {
      await onUpload(booking.id, selectedFile);
      toast({
        title: 'อัปโหลดสำเร็จ',
        description: 'อัปโหลดหนังสือเดินทางเรียบร้อยแล้ว',
      });
      // Reload booking data เพื่อดึงข้อมูล passport ใหม่
      onRefetch?.();
      // Clear form but keep sheet open
      setSelectedFile(null);
    } catch (error) {
      toast({
        title: 'เกิดข้อผิดพลาด',
        description: error instanceof Error ? error.message : 'ไม่สามารถอัปโหลดหนังสือเดินทางได้',
        variant: 'destructive',
      });
    }
  };

  const handleClose = () => {
    setSelectedFile(null);
    onOpenChange(false);
  };




  return (
    <Sheet open={isOpen} onOpenChange={handleClose}>
      <SheetContent className="w-full sm:max-w-lg overflow-y-auto p-6">
        <SheetHeader className="p-0">
          <SheetTitle className="text-left pr-6">แนบหนังสือเดินทาง</SheetTitle>
          <SheetDescription className="text-left">
            เลขที่ใบจอง: {booking?.invoice_id}
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* Upload Section */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Upload className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>อัปโหลดหนังสือเดินทางใหม่</span>
              </div>
            </h3>

            {/* File Upload */}
            <div className="space-y-2">
              <Label>หนังสือเดินทาง</Label>
              <FileUpload
                accept=".jpg,.jpeg,.png,.pdf"
                maxSize={5 * 1024 * 1024}
                allowedTypes={['image/jpeg', 'image/jpg', 'image/png', 'application/pdf']}
                placeholder="เลือกไฟล์หนังสือเดินทาง"
                onFileSelect={handleFileSelect}
                onUploadError={handleFileError}
                variant="dropzone"
                size="md"
                showPreview={true}
                uploadType="custom"
                disabled={isLoading}
              />
              <p className="text-xs text-muted-foreground">
                รองรับไฟล์ JPG, PNG, PDF (ขนาดไม่เกิน 5MB)
              </p>
            </div>

            {/* Info */}
            <div className="bg-muted/50 border border-border rounded-lg p-3">
              <div className="flex">
                <div className="flex-shrink-0">
                  <FileText className="h-5 w-5 text-muted-foreground" />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-foreground">
                    ข้อมูลที่ต้องการ
                  </h3>
                  <div className="mt-2 text-sm text-muted-foreground">
                    <p>• หน้าหนังสือเดินทางที่มีรูปภาพและข้อมูลส่วนตัว</p>
                    <p>• ไฟล์ต้องมีความชัดเจน สามารถอ่านข้อมูลได้</p>
                    <p>• ระบบจะแยกข้อมูลจากหนังสือเดินทางอัตโนมัติ</p>
                  </div>
                </div>
              </div>
            </div>

            {/* Submit Button */}
            <Button
              onClick={handleSubmit}
              disabled={!selectedFile || isLoading}
              className="w-full"
            >
              {isLoading ? 'กำลังอัปโหลด...' : 'อัปโหลดหนังสือเดินทาง'}
            </Button>
          </div>

          <Separator />

          {/* Existing Passports */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <User className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>หนังสือเดินทางที่อัปโหลดแล้ว</span>
              </div>
            </h3>

            <div className="text-xs text-muted-foreground space-y-1 mt-1">
              {passports.length > 0 ? (
                <>
                  {passports.map((passport, index) => (
                    <div key={passport.id} className="border rounded p-2 space-y-2">
                      <div className="flex items-start space-x-3">
                        {passport.image ? (
                          <img
                            src={getStorageUrl(passport.image)}
                            alt="Passport thumbnail"
                            className="w-16 h-16 object-cover rounded border cursor-pointer hover:opacity-80 flex-shrink-0"
                            onClick={() => window.open(getStorageUrl(passport.image), '_blank')}
                          />
                        ) : (
                          <div className="w-16 h-16 bg-muted/50 rounded border flex items-center justify-center flex-shrink-0">
                            <FileText className="h-6 w-6 text-muted-foreground" />
                          </div>
                        )}
                        <div className="flex-1">
                          <p className="font-medium text-sm text-foreground">
                            หนังสือเดินทาง #{index + 1}
                            {passport.name_en && passport.surname_en && (
                              <span className="ml-2 text-muted-foreground">
                                ({passport.name_en} {passport.surname_en})
                              </span>
                            )}
                          </p>
                          {passport.passport_number && (
                            <p className="text-xs text-muted-foreground">
                              เลขที่: {passport.passport_number}
                            </p>
                          )}
                          <p className="text-xs text-muted-foreground">
                            อัปโหลดเมื่อ: {new Date(passport.created_at).toLocaleDateString('th-TH', {
                              year: 'numeric',
                              month: 'long',
                              day: 'numeric',
                              hour: '2-digit',
                              minute: '2-digit'
                            })}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <User className="mx-auto h-12 w-12 text-muted-foreground/30" />
                  <p className="mt-2">ยังไม่มีหนังสือเดินทาง</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}