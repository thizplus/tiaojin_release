import React, { memo, useCallback, useState, useRef } from 'react';
import { Button } from '@/components/ui/button';

import { Progress } from '@/components/ui/progress';
import { Card, CardContent } from '@/components/ui/card';
import { useFileUpload } from '@/hooks/useFileUpload';
import { cn, getStorageUrl } from '@/lib/utils';
import {
  Upload,
  X,
  File,

  CheckCircle,
  AlertCircle,
  Loader2,
  FileText,
  FileImage
} from 'lucide-react';
import type {
  FileUploadOptions,
  FileUploadResult,
  StructuredPathUploadOptions
} from '@/types/file';

interface FileUploadProps {
  // Upload options
  uploadType?: 'custom' | 'structured';
  customPath?: string;
  structuredOptions?: Omit<StructuredPathUploadOptions, 'file'>;

  // File constraints
  accept?: string;
  maxSize?: number; // in bytes
  allowedTypes?: string[];

  // UI options
  className?: string;
  placeholder?: string;
  showProgress?: boolean;
  showPreview?: boolean;
  multiple?: boolean;
  disabled?: boolean;

  // Event handlers
  onUploadSuccess?: (result: FileUploadResult) => void;
  onUploadError?: (error: Error) => void;
  onFileSelect?: (file: File | null) => void;

  // Controlled value
  value?: string; // URL
  onChange?: (url: string | null) => void;

  // Appearance
  variant?: 'default' | 'dropzone' | 'button';
  size?: 'sm' | 'md' | 'lg';
}

export const FileUpload = memo(function FileUpload({
  uploadType = 'structured',
  customPath,
  structuredOptions = {},
  accept = 'image/*',
  maxSize = 10 * 1024 * 1024, // 10MB
  allowedTypes,
  className,
  placeholder = 'เลือกไฟล์หรือลากมาวาง',
  showProgress = true,
  showPreview = true,
  multiple = false,
  disabled = false,
  onUploadSuccess,
  onUploadError,
  onFileSelect,
  value,
  onChange,
  variant = 'dropzone',
  size = 'md'
}: FileUploadProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isDragOver, setIsDragOver] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string | null>(value || null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // resolve URL: blob URL ปล่อยผ่าน, backend URL ผ่าน getStorageUrl()
  const resolveUrl = (url: string | null): string | null => {
    if (!url) return null;
    if (url.startsWith('blob:') || url.startsWith('data:')) return url;
    return getStorageUrl(url);
  };

  // Sync previewUrl with value prop
  React.useEffect(() => {
    setPreviewUrl(value || null);
  }, [value]);

  const {
    isUploading,
    progress,
    error,
    result,
    uploadFile,
    reset,
    validateFile
  } = useFileUpload({
    maxSize,
    allowedTypes,
    onSuccess: (uploadResult) => {
      if (uploadResult.success) {
        const url = uploadResult.data.url;
        setPreviewUrl(url);
        setSelectedFile(null); // Clear selected file to show success state
        onChange?.(url);
        onUploadSuccess?.(uploadResult);
      }
    },
    onError: onUploadError
  });

  // Handle file selection
  const handleFileSelect = useCallback((file: File | null) => {
    setSelectedFile(file);
    onFileSelect?.(file);

    if (file) {
      // Create preview URL for images
      if (file.type.startsWith('image/') && showPreview) {
        const url = URL.createObjectURL(file);
        setPreviewUrl(url);
      }

      // Auto upload if configured
      handleUpload(file);
    } else {
      setPreviewUrl(value || null);
    }
  }, [onFileSelect, showPreview, value]);

  // Handle file upload
  const handleUpload = useCallback(async (file: File) => {
    if (!file) return;

    // Validate file
    const validation = validateFile(file);
    if (!validation.isValid) {
      onUploadError?.(new Error(validation.error));
      return;
    }

    // Prepare upload options
    const uploadOptions: FileUploadOptions = uploadType === 'custom' && customPath
      ? { file, custom_path: customPath }
      : { file, ...structuredOptions };

    // Upload file
    await uploadFile(uploadOptions);
  }, [uploadType, customPath, structuredOptions, uploadFile, validateFile, onUploadError]);

  // Handle input change
  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    handleFileSelect(file);
  }, [handleFileSelect]);

  // Handle drag events
  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);

    const files = Array.from(e.dataTransfer.files);
    const file = files[0] || null;
    handleFileSelect(file);
  }, [handleFileSelect]);

  // Handle remove file
  const handleRemove = useCallback(() => {
    setSelectedFile(null);
    setPreviewUrl(null);
    onChange?.(null);
    reset();

    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  }, [onChange, reset]);

  // Handle click to select file
  const handleClick = useCallback(() => {
    if (!disabled) {
      fileInputRef.current?.click();
    }
  }, [disabled]);

  // Get file icon
  const getFileIcon = useCallback((file: File | null, size: 'sm' | 'lg' = 'sm') => {
    const iconSize = size === 'lg' ? 'h-12 w-12' : 'h-4 w-4';
    const iconColor = size === 'lg' ? 'text-muted-foreground' : '';

    if (!file) {
      return <File className={`${iconSize} ${iconColor}`} />;
    }

    if (file.type.startsWith('image/')) {
      return <FileImage className={`${iconSize} ${iconColor}`} />;
    }
    if (file.type === 'application/pdf') {
      return <FileText className={`${iconSize} ${iconColor} text-destructive`} />;
    }
    if (file.type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document' ||
        file.type === 'application/msword') {
      return <FileText className={`${iconSize} ${iconColor} text-primary`} />;
    }
    return <File className={`${iconSize} ${iconColor}`} />;
  }, []);

  // Size classes
  const sizeClasses = {
    sm: 'p-4 text-sm',
    md: 'p-6 text-base',
    lg: 'p-8 text-lg'
  };

  // Render button variant
  if (variant === 'button') {
    return (
      <div className={className}>
        <input
          ref={fileInputRef}
          type="file"
          accept={accept}
          multiple={multiple}
          onChange={handleInputChange}
          className="hidden"
        />
        <Button
          type="button"
          variant="outline"
          onClick={handleClick}
          disabled={disabled || isUploading}
          className="w-full"
        >
          {isUploading ? (
            <Loader2 className="h-4 w-4 animate-spin mr-2" />
          ) : (
            <Upload className="h-4 w-4 mr-2" />
          )}
          {isUploading ? 'กำลังอัปโหลด...' : 'เลือกไฟล์'}
        </Button>

        {showProgress && isUploading && (
          <div className="mt-2">
            <Progress value={progress} className="h-2" />
            <p className="text-sm text-muted-foreground mt-1">{progress}%</p>
          </div>
        )}

        {error && (
          <div className="mt-2 flex items-center text-destructive text-sm">
            <AlertCircle className="h-4 w-4 mr-1" />
            {error}
          </div>
        )}
      </div>
    );
  }

  // Render dropzone variant
  return (
    <div className={className}>
      <input
        ref={fileInputRef}
        type="file"
        accept={accept}
        multiple={multiple}
        onChange={handleInputChange}
        className="hidden"
      />

      <Card
        className={cn(
          'border-2 border-dashed transition-colors',
          !disabled && 'cursor-pointer',
          sizeClasses[size],
          isDragOver && !disabled && 'border-primary bg-primary/5',
          isUploading && 'border-muted-foreground/50',
          error && 'border-destructive',
          disabled && 'opacity-50 cursor-not-allowed'
        )}
        onDragOver={disabled ? undefined : handleDragOver}
        onDragLeave={disabled ? undefined : handleDragLeave}
        onDrop={disabled ? undefined : handleDrop}
        onClick={disabled ? undefined : handleClick}
      >
        <CardContent className="p-0">
          {/* Preview */}
          {(selectedFile || (showPreview && previewUrl)) && (
            <div className="relative mb-4">
              {/* Show image preview for image files (both during selection and after upload) */}
              {showPreview && previewUrl && (previewUrl.includes('.jpg') || previewUrl.includes('.jpeg') || previewUrl.includes('.png') || previewUrl.includes('.webp') || (selectedFile && selectedFile.type.startsWith('image/'))) ? (
                // Image preview with square aspect ratio
                <div className="aspect-square w-full max-w-32 mx-auto">
                  <img
                    src={resolveUrl(previewUrl) || ''}
                    alt="Preview"
                    className="w-full h-full rounded-lg object-cover border-2 border-border"
                  />
                </div>
              ) : previewUrl && !selectedFile ? (
                // Uploaded file preview for non-images (after upload success, no selectedFile)
                <div className="flex items-center gap-3 p-3 border-2 border-dashed border-border rounded-lg bg-muted/20">
                  <CheckCircle className="h-6 w-6 text-green-500" />
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm break-all">
                      {previewUrl.includes('.pdf') ? 'PDF Document.pdf' :
                       previewUrl.includes('.doc') ? 'Word Document.docx' :
                       'Uploaded File'}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {previewUrl.includes('.pdf') ? 'PDF Document' :
                       previewUrl.includes('.doc') ? 'Word Document' :
                       'File'} • อัพโหลดสำเร็จ
                    </p>
                  </div>
                </div>
              ) : selectedFile ? (
                // File icon + name for non-images during selection/upload
                <div className="flex items-center gap-3 p-3 border-2 border-dashed border-border rounded-lg bg-muted/20">
                  {getFileIcon(selectedFile, 'lg')}
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm break-all">{selectedFile.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {selectedFile.type === 'application/pdf' ? 'PDF Document' :
                       selectedFile.type.includes('word') ? 'Word Document' :
                       selectedFile.type || 'Unknown File'} • {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                    </p>
                  </div>
                </div>
              ) : null}
              <Button
                type="button"
                variant="destructive"
                size="sm"
                className="absolute -top-2 -right-2 h-6 w-6 p-0 rounded-full"
                onClick={(e) => {
                  e.stopPropagation();
                  handleRemove();
                }}
              >
                <X className="h-3 w-3" />
              </Button>
            </div>
          )}

          {/* Upload area */}
          <div className="text-center">
            {isUploading ? (
              <div className="space-y-2">
                <Loader2 className="h-8 w-8 animate-spin mx-auto text-primary" />
                <p className="text-muted-foreground">กำลังอัปโหลด...</p>
                {showProgress && (
                  <div className="space-y-1">
                    <Progress value={progress} className="h-2" />
                    <p className="text-sm text-muted-foreground">{progress}%</p>
                  </div>
                )}
              </div>
            ) : result?.success ? (
              <div className="space-y-2">
                <CheckCircle className="h-8 w-8 mx-auto text-green-500" />
                <p className="text-green-500">อัปโหลดสำเร็จ</p>
                <p className="text-sm text-muted-foreground">{result.data.fileName}</p>
              </div>
            ) : (selectedFile || previewUrl) ? (
              <div className="mt-3 p-2 bg-muted/50 rounded-md">
                <p className="text-xs text-muted-foreground">
                  คลิกเพื่อเปลี่ยนไฟล์
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                <Upload className="h-8 w-8 mx-auto text-muted-foreground" />
                <p className="text-muted-foreground">{placeholder}</p>
                <p className="text-xs text-muted-foreground">
                  ขนาดไฟล์สูงสุด: {(maxSize / 1024 / 1024).toFixed(0)} MB
                </p>
              </div>
            )}

            {error && (
              <div className="mt-4 flex items-center justify-center text-destructive text-sm">
                <AlertCircle className="h-4 w-4 mr-1" />
                {error}
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
});