import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { FileUpload } from '@/components/ui/file-upload';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Upload,  FileText, Receipt, Loader2 } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { getApprovalStatusConfig } from '@/utils/bookingStatus';
import { getStorageUrl } from '@/lib/utils';
import agentBookingService from '@/services/agentBookingService';
import type { BookingResponse, BookingDocumentResponse } from '@/types';

interface PaymentSlipUploadSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  booking: BookingResponse | null;
  onUpload: (id: string, file: File, amount: number) => Promise<void>;
  isLoading?: boolean;
}

export function PaymentSlipUploadSheet({
  isOpen,
  onOpenChange,
  booking,
  onUpload,
  isLoading = false
}: PaymentSlipUploadSheetProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [amount, setAmount] = useState<string>('');
  const [documents, setDocuments] = useState<BookingDocumentResponse[]>([]);
  const [documentsLoading, setDocumentsLoading] = useState(false);
  const [documentsLoaded, setDocumentsLoaded] = useState(false);
  const { toast } = useToast();

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
    if (!booking || !selectedFile || !amount) {
      toast({
        title: 'ข้อมูลไม่ครบ',
        description: 'กรุณาเลือกไฟล์และกรอกจำนวนเงิน',
        variant: 'destructive',
      });
      return;
    }

    const amountNum = parseFloat(amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      toast({
        title: 'จำนวนเงินไม่ถูกต้อง',
        description: 'กรุณากรอกจำนวนเงินที่ถูกต้อง',
        variant: 'destructive',
      });
      return;
    }

    try {
      await onUpload(booking.id, selectedFile, amountNum);
      toast({
        title: 'อัปโหลดสำเร็จ',
        description: 'อัปโหลดสลิปการโอนเงินเรียบร้อยแล้ว',
      });
      // Reload documents silently after successful upload
      await loadDocuments(true);
      // Clear form but keep sheet open
      setSelectedFile(null);
      setAmount('');
    } catch (error) {
      toast({
        title: 'เกิดข้อผิดพลาด',
        description: error instanceof Error ? error.message : 'ไม่สามารถอัปโหลดสลิปได้',
        variant: 'destructive',
      });
    }
  };

  const handleClose = () => {
    setSelectedFile(null);
    setAmount('');
    setDocumentsLoaded(false);
    setDocuments([]);
    onOpenChange(false);
  };

  // Load documents when booking changes or upload succeeds
  const loadDocuments = async (silent = false) => {
    if (!booking?.id) return;

    if (!silent) {
      setDocumentsLoading(true);
    }
    try {
      const response = await agentBookingService.getDocuments(booking.id);
      if (response.success) {
        setDocuments(response.data.documents.filter(doc => doc.type === 'payment_slip'));
        if (!silent) {
          setDocumentsLoaded(true);
        }
      }
    } catch (error) {
      console.error('Failed to load documents:', error);
    } finally {
      if (!silent) {
        setDocumentsLoading(false);
      }
    }
  };

  // Reset documents when booking changes
  useEffect(() => {
    setDocumentsLoaded(false);
    setDocuments([]);
  }, [booking?.id]);

  useEffect(() => {
    if (isOpen && booking?.id && !documentsLoaded) {
      loadDocuments();
    }
  }, [isOpen, booking?.id, documentsLoaded]);



  return (
    <Sheet open={isOpen} onOpenChange={handleClose}>
      <SheetContent className="w-full sm:max-w-lg overflow-y-auto p-6">
        <SheetHeader className="p-0">
          <SheetTitle className="text-left pr-6">แจ้งชำระเงิน</SheetTitle>
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
                <span>อัปโหลดสลิปใหม่</span>
              </div>
            </h3>

            {/* Amount Input */}
            <div className="space-y-2">
              <Label htmlFor="amount">จำนวนเงินที่โอน (บาท)</Label>
              <Input
                id="amount"
                type="number"
                placeholder="0.00"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                min="0"
                step="0.01"
              />
            </div>

            {/* File Upload */}
            <div className="space-y-2">
              <Label>สลิปการโอนเงิน</Label>
              <FileUpload
                accept=".jpg,.jpeg,.png,.pdf"
                maxSize={5 * 1024 * 1024}
                allowedTypes={['image/jpeg', 'image/jpg', 'image/png', 'application/pdf']}
                placeholder="เลือกไฟล์สลิปการโอนเงิน"
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

            {/* Submit Button */}
            <Button
              onClick={handleSubmit}
              disabled={!selectedFile || !amount || isLoading}
              className="w-full"
            >
              {isLoading ? 'กำลังอัปโหลด...' : 'อัปโหลดสลิป'}
            </Button>
          </div>

          <Separator />

          {/* Existing Payment Slips */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium flex items-center gap-2">
              <div className="flex-shrink-0 p-2 bg-muted/50 rounded-md">
                <Receipt className="h-10 w-10 text-muted-foreground" />
              </div>
              <div className="leading-none">
                <span>สลิปที่อัปโหลดแล้ว</span>
              </div>
            </h3>

            <div className="text-xs text-muted-foreground space-y-1 mt-1">
              {documentsLoading && !documentsLoaded ? (
                <div className="text-center py-8">
                  <Loader2 className="mx-auto h-8 w-8 animate-spin text-muted-foreground" />
                  <p className="mt-2 text-muted-foreground">กำลังโหลดเอกสาร...</p>
                </div>
              ) : documents.length > 0 ? (
                <>
                  {documents.map((doc, index) => (
                    <div key={doc.id} className="border rounded p-2 space-y-2">
                      <div className="flex items-start justify-between">
                        <div className="flex items-start space-x-3">
                          {doc.file_url ? (
                            <img
                              src={getStorageUrl(doc.file_url)}
                              alt="Payment slip thumbnail"
                              className="w-16 h-16 object-cover rounded border cursor-pointer hover:opacity-80 flex-shrink-0"
                              onClick={() => window.open(getStorageUrl(doc.file_url), '_blank')}
                            />
                          ) : (
                            <div className="w-16 h-16 bg-muted/50 rounded border flex items-center justify-center flex-shrink-0">
                              <FileText className="h-6 w-6 text-muted-foreground" />
                            </div>
                          )}
                          <div className="flex-1">
                            <p className="font-medium text-sm text-foreground">สลิป #{index + 1}</p>
                            <p className="text-xs text-muted-foreground">
                              จำนวน: {doc.amount?.toLocaleString()} บาท
                            </p>
                            <p className="text-xs text-muted-foreground">
                              ไฟล์: {doc.filename}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              อัปโหลดเมื่อ: {new Date(doc.uploaded_at).toLocaleDateString('th-TH', {
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit'
                              })}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          {(() => {
                            const statusConfig = getApprovalStatusConfig(doc.status);
                            return (
                              <>
                                {statusConfig.icon}
                                <Badge className={`${statusConfig.bgColor} ${statusConfig.textColor}`}>
                                  {statusConfig.label}
                                </Badge>
                              </>
                            );
                          })()}
                        </div>
                      </div>
                    </div>
                  ))}
                </>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <FileText className="mx-auto h-12 w-12 text-muted-foreground/30" />
                  <p className="mt-2">ยังไม่มีสลิปการโอนเงิน</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}