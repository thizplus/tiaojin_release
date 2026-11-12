
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import {
  FileText,
  Download,
  ExternalLink,
  FileImage,
  File,
  AlertCircle,
} from 'lucide-react';
import { getStorageUrl } from '@/lib/utils';

interface PeriodFile {
  name: string;
  type: string;
  url: string;
}

interface FileDownloadSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  files: PeriodFile[];
  periodInfo?: {
    startDate: string;
    endDate: string;
    days: number;
  };
  tourInfo?: {
    code: string;
    title: string;
  };
  /** ถ้าส่ง isFullyPaid มา จะเช็คเงื่อนไขชำระเงินครบก่อนดาวน์โหลด (agent side) */
  isFullyPaid?: boolean;
}

export function FileDownloadSheet({
  isOpen,
  onOpenChange,
  files,
  periodInfo,
  tourInfo,
  isFullyPaid
}: FileDownloadSheetProps) {

  // ถ้า isFullyPaid เป็น undefined → ไม่ต้องเช็ค (admin/sale), ถ้าเป็น boolean → เช็คชำระครบ
  const canDownload = isFullyPaid === undefined || isFullyPaid;

  const getFileTypeColor = (type: string) => {
    if (type.includes('pdf')) return 'bg-destructive/10 text-destructive border-destructive/20';
    if (type.includes('document') || type.includes('word')) return 'bg-chart-1/10 text-chart-1 border-chart-1/20';
    if (type.includes('image')) return 'bg-chart-2/10 text-chart-2 border-chart-2/20';
    if (type.includes('video')) return 'bg-chart-3/10 text-chart-3 border-chart-3/20';
    return 'bg-muted text-muted-foreground border-border';
  };

  const getFileTypeIcon = (type: string) => {
    if (type.includes('image')) {
      return <FileImage className="h-5 w-5" />;
    }
    if (type.includes('pdf') || type.includes('document') || type.includes('word')) {
      return <FileText className="h-5 w-5" />;
    }
    return <File className="h-5 w-5" />;
  };

  const getFileTypeDisplay = (type: string) => {
    if (type.includes('pdf')) return 'PDF';
    if (type.includes('image/jpeg') || type.includes('image/jpg')) return 'JPEG';
    if (type.includes('image/png')) return 'PNG';
    if (type.includes('image/webp')) return 'WebP';
    if (type.includes('word') || type.includes('officedocument.wordprocessingml')) return 'Word';
    if (type.includes('sheet') || type.includes('officedocument.spreadsheetml')) return 'Excel';
    if (type.includes('presentation') || type.includes('officedocument.presentationml')) return 'PowerPoint';
    return 'File';
  };

  const handleDownload = (file: PeriodFile) => {
    if (!canDownload) return;
    // สร้างชื่อไฟล์จาก tour code + title แทน UUID
    const ext = file.name.includes('.') ? file.name.substring(file.name.lastIndexOf('.')) : '';
    const downloadName = tourInfo
      ? `${tourInfo.code} ${tourInfo.title}${ext}`
      : file.name;

    // ส่ง filename ผ่าน query param เพื่อให้ backend ตั้ง Content-Disposition header
    const baseUrl = getStorageUrl(file.url);
    const separator = baseUrl.includes('?') ? '&' : '?';
    const downloadUrl = `${baseUrl}${separator}filename=${encodeURIComponent(downloadName)}`;

    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = downloadName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handleView = (file: PeriodFile) => {
    if (!canDownload) return;
    window.open(getStorageUrl(file.url), '_blank', 'noopener,noreferrer');
  };

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-lg overflow-y-auto p-6">
        <SheetHeader className="p-0">
          <SheetTitle className="text-left pr-6">
            เอกสารประกอบทัวร์
          </SheetTitle>
          <SheetDescription className="text-left">
            {periodInfo ? (
              `ช่วงเวลา: ${new Date(periodInfo.startDate).toLocaleDateString('th-TH')} - ${new Date(periodInfo.endDate).toLocaleDateString('th-TH')} (${periodInfo.days} วัน)`
            ) : (
              'ไฟล์เอกสารและข้อมูลประกอบการเดินทาง'
            )}
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-4">
          {/* แสดงเงื่อนไขเมื่อยังชำระเงินไม่ครบ */}
          {isFullyPaid === false && (
            <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4">
              <div className="flex items-center gap-2 text-destructive">
                <AlertCircle className="h-5 w-5" />
                <p className="text-sm font-medium">กรุณาชำระเงินให้ครบถ้วนก่อนดาวน์โหลดเอกสาร</p>
              </div>
            </div>
          )}

          {files.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              <FileText className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>ไม่มีไฟล์เอกสาร</p>
            </div>
          ) : (
            files.map((file, index) => (
              <Card key={index} className={`p-0 transition-shadow ${canDownload ? 'hover:shadow-md' : 'opacity-50'}`}>
                <CardContent className="p-4">
                  <div className="flex items-start space-x-3">
                    <div className="flex-shrink-0 mt-1">
                      {getFileTypeIcon(file.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h4 className="font-medium text-sm break-all leading-5 mb-2">
                        {file.name}
                      </h4>
                      <div className="flex items-center space-x-2 mb-3">
                        <Badge
                          variant="outline"
                          className={`${getFileTypeColor(file.type)} text-xs`}
                        >
                          {getFileTypeDisplay(file.type)}
                        </Badge>
                      </div>
                      <div className="flex gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleView(file)}
                          className="text-lg font-bold"
                          disabled={!canDownload}
                        >
                          <ExternalLink className="h-3 w-3 mr-1" />
                          เปิดดู
                        </Button>
                        <Button
                          size="sm"
                          onClick={() => handleDownload(file)}
                          className="text-lg font-bold"
                          disabled={!canDownload}
                        >
                          <Download className="h-3 w-3 mr-1" />
                          ดาวน์โหลด
                        </Button>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))
          )}

          {files.length > 0 && (
            <div className="pt-4 border-t">
              <p className="text-xs text-muted-foreground">
                {canDownload
                  ? '⚠️ ไฟล์เหล่านี้เป็นเอกสารประกอบการเดินทาง กรุณาตรวจสอบข้อมูลให้ครบถ้วนก่อนเดินทาง'
                  : '🔒 กรุณาชำระเงินให้ครบถ้วนเพื่อดาวน์โหลดเอกสาร'
                }
              </p>
            </div>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}
