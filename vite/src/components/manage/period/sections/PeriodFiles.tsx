
import { useFieldArray, type UseFormReturn } from 'react-hook-form';
import { Button } from '@/components/ui/button';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { FileUpload } from '@/components/ui/file-upload';
import { FileText, Plus, Trash2, ExternalLink } from 'lucide-react';
import type { PeriodFormData } from '../schemas/period-form-schemas';

interface PeriodFilesProps {
  form: UseFormReturn<PeriodFormData>;
}

export function PeriodFiles({ form }: PeriodFilesProps) {
  const { control, formState: { errors } } = form;

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'files',
  });

  const addFile = () => {
    append({
      name: '',
      type: '',
      url: '',
    });
  };

  const handleFileUpload = (index: number, url: string | null, fileName?: string, fileType?: string) => {
    if (url) {
      // Auto-fill all file data from upload result
      console.log('File upload data:', { url, fileName, fileType });

      form.setValue(`files.${index}.url`, url);

      if (fileName) {
        form.setValue(`files.${index}.name`, fileName);
      }

      if (fileType) {
        // Truncate long MIME types to fit API constraints (max 100 chars)
        const shortType = fileType.length > 100 ? fileType.substring(0, 100) : fileType;
        form.setValue(`files.${index}.type`, shortType);
      }

      // Trigger validation
      form.trigger(`files.${index}`);
    }
  };

  const getFileTypeColor = (type: string) => {
    if (type.includes('pdf')) return 'bg-destructive/10 text-destructive border-destructive/20';
    if (type.includes('document') || type.includes('word')) return 'bg-chart-1/10 text-chart-1 border-chart-1/20';
    if (type.includes('image')) return 'bg-chart-2/10 text-chart-2 border-chart-2/20';
    if (type.includes('video')) return 'bg-chart-3/10 text-chart-3 border-chart-3/20';
    return 'bg-muted text-muted-foreground border-border';
  };

  const getFileTypeIcon = () => {
    return <FileText className="h-4 w-4" />;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileText className="h-5 w-5" />
          ไฟล์แนบ
        </CardTitle>
        <p className="text-sm text-muted-foreground">
          จัดการไฟล์แนบสำหรับช่วงเวลาทัวร์ เช่น รูปภาพ, เอกสาร, ใบปลิว
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        {fields.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            <FileText className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>ยังไม่มีไฟล์แนบ</p>
            <p className="text-sm">คลิกปุ่ม "เพิ่มไฟล์" เพื่อเพิ่มไฟล์แนบ</p>
          </div>
        ) : (
          <div className="space-y-3">
            {fields.map((field, index) => {
              const fileData = form.watch(`files.${index}`);
              const hasErrors = errors.files?.[index];

              return (
                <Card key={field.id} className={`${hasErrors ? 'border-destructive' : ''}`}>
                  <CardContent className="p-4">
                    {fileData?.name && fileData?.url ? (
                      // Show uploaded file info
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="flex-shrink-0">
                            {getFileTypeIcon()}
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="font-medium truncate">{fileData.name}</p>
                            <div className="flex items-center space-x-2 mt-1">
                              <Badge variant="outline" className={getFileTypeColor(fileData.type || '')}>
                                {fileData.type}
                              </Badge>
                              <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                className="h-6 px-2 text-xs"
                                onClick={() => window.open(fileData.url, '_blank')}
                              >
                                <ExternalLink className="h-3 w-3 mr-1" />
                                เปิดไฟล์
                              </Button>
                            </div>
                          </div>
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => remove(index)}
                          className="text-destructive hover:text-destructive"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    ) : (
                      // Show file upload area
                      <div className="space-y-4">
                        <div className="flex items-center justify-between">
                          <h4 className="font-medium">ไฟล์ที่ {index + 1}</h4>
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => remove(index)}
                            className="text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>

                        <FileUpload
                          structuredOptions={{
                            category: 'periods',
                            file_type: 'document'
                          }}
                          accept="*/*"
                          maxSize={10 * 1024 * 1024} // 10MB
                          value={fileData?.url || ''}
                          onChange={(url) => handleFileUpload(index, url)}
                          onUploadSuccess={(result) => {
                            console.log('Upload result:', result);
                            if (result.success) {
                              const fileName = result.data.fileName || '';
                              const fileType = result.data.mimeType || '';
                              console.log('Extracted from result:', { fileName, fileType, url: result.data.url });
                              handleFileUpload(index, result.data.url, fileName, fileType);
                            }
                          }}
                          onUploadError={(error) => {
                            console.error('File upload error:', error);
                          }}
                          variant="dropzone"
                          size="md"
                          placeholder="เลือกไฟล์หรือลากมาวาง"
                        />
                        {errors.files?.[index]?.url && (
                          <p className="text-xs text-destructive mt-1">
                            {errors.files[index]?.url?.message}
                          </p>
                        )}
                      </div>
                    )}
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}

        <Button
          type="button"
          variant="outline"
          onClick={addFile}
          className="w-full"
        >
          <Plus className="h-4 w-4 mr-2" />
          เพิ่มไฟล์
        </Button>

        {fields.length > 0 && (
          <div className="text-xs text-muted-foreground space-y-1">
            <p>• ไฟล์จะถูกอัปโหลดและจัดเก็บในระบบอัตโนมัติ</p>
            <p>• ขนาดไฟล์สูงสุด 10MB ต่อไฟล์</p>
            <p>• รองรับไฟล์ประเภท PDF, รูปภาพ, วิดีโอ, และเอกสารทุกประเภท</p>
            <p>• ข้อมูลชื่อไฟล์และประเภทจะถูกเติมอัตโนมัติหลังจากอัปโหลด</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}