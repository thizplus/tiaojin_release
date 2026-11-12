import { useState, useCallback } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { FileUpload } from '@/components/ui/file-upload';
import {
  FileText,
  Plus,
  Trash2,
  ExternalLink,
  Loader2,
  FileImage,
  File,
} from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { useUpdatePeriodMutation } from '@/stores/periodStore';
import { getStorageUrl } from '@/lib/utils';
import type { TourPeriod, PeriodFile } from '@/types/models/period';

interface PeriodFilesDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  period: TourPeriod | null;
}

export function PeriodFilesDialog({ open, onOpenChange, period }: PeriodFilesDialogProps) {
  const { toast } = useToast();
  const updateMutation = useUpdatePeriodMutation();
  const [files, setFiles] = useState<PeriodFile[]>([]);
  const [showUpload, setShowUpload] = useState(false);
  const [pendingFile, setPendingFile] = useState<PeriodFile | null>(null);

  // Sync files from period when dialog opens
  const handleOpenChange = useCallback((isOpen: boolean) => {
    if (isOpen && period) {
      setFiles(period.files || []);
      setShowUpload(false);
      setPendingFile(null);
    }
    onOpenChange(isOpen);
  }, [period, onOpenChange]);

  const getFileTypeColor = (type: string) => {
    if (type.includes('pdf')) return 'bg-destructive/10 text-destructive border-destructive/20';
    if (type.includes('document') || type.includes('word')) return 'bg-chart-1/10 text-chart-1 border-chart-1/20';
    if (type.includes('image')) return 'bg-chart-2/10 text-chart-2 border-chart-2/20';
    return 'bg-muted text-muted-foreground border-border';
  };

  const getFileTypeDisplay = (type: string) => {
    if (type.includes('pdf')) return 'PDF';
    if (type.includes('image/jpeg') || type.includes('image/jpg')) return 'JPEG';
    if (type.includes('image/png')) return 'PNG';
    if (type.includes('word') || type.includes('officedocument.wordprocessingml')) return 'Word';
    if (type.includes('sheet') || type.includes('officedocument.spreadsheetml')) return 'Excel';
    return 'File';
  };

  const getFileTypeIcon = (type: string) => {
    if (type.includes('image')) return <FileImage className="h-5 w-5" />;
    if (type.includes('pdf') || type.includes('document') || type.includes('word'))
      return <FileText className="h-5 w-5" />;
    return <File className="h-5 w-5" />;
  };

  const handleUploadSuccess = useCallback((result: { success: true; data: { url: string; fileName: string; mimeType: string } } | { success: false }) => {
    if (result.success) {
      setPendingFile({
        name: result.data.fileName || '',
        type: result.data.mimeType || '',
        url: result.data.url,
      });
    }
  }, []);

  const handleConfirmAdd = useCallback(async () => {
    if (!pendingFile || !period) return;

    const newFiles = [...files, pendingFile];
    try {
      await updateMutation.mutateAsync({
        id: period.id,
        data: { files: newFiles },
      });
      setFiles(newFiles);
      setPendingFile(null);
      setShowUpload(false);
      toast({ title: 'เพิ่มใบนัดหมายสำเร็จ' });
    } catch {
      toast({ title: 'เกิดข้อผิดพลาด', description: 'ไม่สามารถเพิ่มไฟล์ได้', variant: 'destructive' });
    }
  }, [pendingFile, period, files, updateMutation, toast]);

  const handleRemoveFile = useCallback(async (index: number) => {
    if (!period) return;

    const newFiles = files.filter((_, i) => i !== index);
    try {
      await updateMutation.mutateAsync({
        id: period.id,
        data: { files: newFiles },
      });
      setFiles(newFiles);
      toast({ title: 'ลบไฟล์สำเร็จ' });
    } catch {
      toast({ title: 'เกิดข้อผิดพลาด', description: 'ไม่สามารถลบไฟล์ได้', variant: 'destructive' });
    }
  }, [period, files, updateMutation, toast]);

  if (!period) return null;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-lg max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5" />
            จัดการใบนัดหมาย
          </DialogTitle>
          <DialogDescription>
            เพิ่มหรือลบไฟล์ใบนัดหมายสำหรับ period นี้
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          {/* รายการไฟล์ที่มีอยู่ */}
          {files.length === 0 && !showUpload ? (
            <div className="text-center py-8 text-muted-foreground">
              <FileText className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>ยังไม่มีใบนัดหมาย</p>
              <p className="text-sm">คลิกปุ่ม "เพิ่มไฟล์" เพื่อเพิ่มใบนัดหมาย</p>
            </div>
          ) : (
            <div className="space-y-3">
              {files.map((file, index) => (
                <div
                  key={`${file.url}-${index}`}
                  className="flex items-center justify-between p-3 border rounded-lg"
                >
                  <div className="flex items-center gap-3 min-w-0 flex-1">
                    <div className="flex-shrink-0">
                      {getFileTypeIcon(file.type)}
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="font-medium text-sm truncate">{file.name}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge variant="outline" className={`text-xs ${getFileTypeColor(file.type)}`}>
                          {getFileTypeDisplay(file.type)}
                        </Badge>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          className="h-6 px-2 text-xs"
                          onClick={() => window.open(getStorageUrl(file.url), '_blank')}
                        >
                          <ExternalLink className="h-3 w-3 mr-1" />
                          เปิดดู
                        </Button>
                      </div>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => handleRemoveFile(index)}
                    disabled={updateMutation.isPending}
                    className="text-destructive hover:text-destructive flex-shrink-0"
                  >
                    {updateMutation.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Trash2 className="h-4 w-4" />
                    )}
                  </Button>
                </div>
              ))}
            </div>
          )}

          {/* Upload area */}
          {showUpload && (
            <div className="border rounded-lg p-4 space-y-3">
              <h4 className="font-medium text-sm">อัปโหลดไฟล์ใหม่</h4>
              {!pendingFile ? (
                <FileUpload
                  structuredOptions={{
                    category: 'periods',
                    file_type: 'document',
                  }}
                  accept="*/*"
                  maxSize={10 * 1024 * 1024}
                  value=""
                  onChange={() => {}}
                  onUploadSuccess={handleUploadSuccess}
                  onUploadError={(error) => {
                    toast({ title: 'อัปโหลดไม่สำเร็จ', description: error.message, variant: 'destructive' });
                  }}
                  variant="dropzone"
                  size="md"
                  placeholder="เลือกไฟล์หรือลากมาวาง (สูงสุด 10MB)"
                />
              ) : (
                <div className="space-y-3">
                  <div className="flex items-center gap-3 p-3 border rounded-lg bg-muted/50">
                    {getFileTypeIcon(pendingFile.type)}
                    <div className="min-w-0 flex-1">
                      <p className="font-medium text-sm truncate">{pendingFile.name}</p>
                      <Badge variant="outline" className={`text-xs mt-1 ${getFileTypeColor(pendingFile.type)}`}>
                        {getFileTypeDisplay(pendingFile.type)}
                      </Badge>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      onClick={handleConfirmAdd}
                      disabled={updateMutation.isPending}
                      className="flex-1"
                    >
                      {updateMutation.isPending ? (
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      ) : null}
                      บันทึก
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() => {
                        setPendingFile(null);
                        setShowUpload(false);
                      }}
                      disabled={updateMutation.isPending}
                    >
                      ยกเลิก
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Add file button */}
          {!showUpload && (
            <Button
              variant="outline"
              onClick={() => setShowUpload(true)}
              className="w-full"
            >
              <Plus className="h-4 w-4 mr-2" />
              เพิ่มไฟล์
            </Button>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
