import { useState, useCallback } from 'react';
import { fileService } from '@/services/admin/fileService';
import type {
  FileUploadOptions,
  FileUploadResult,
  FileUploadProgressEvent,
  StructuredPathUploadOptions
} from '@/types/file';

interface UseFileUploadOptions {
  onSuccess?: (result: FileUploadResult) => void;
  onError?: (error: Error) => void;
  onProgress?: (event: FileUploadProgressEvent) => void;
  validateFile?: boolean;
  maxSize?: number; // in bytes
  allowedTypes?: string[]; // mime types
}

interface UseFileUploadReturn {
  // States
  isUploading: boolean;
  progress: number;
  error: string | null;
  result: FileUploadResult | null;

  // Actions
  uploadFile: (options: FileUploadOptions) => Promise<FileUploadResult | null>;
  uploadWithCustomPath: (file: File, customPath: string) => Promise<FileUploadResult | null>;
  uploadWithStructuredPath: (
    file: File,
    options?: Omit<StructuredPathUploadOptions, 'file'>
  ) => Promise<FileUploadResult | null>;
  uploadCityImage: (file: File, cityId?: string) => Promise<FileUploadResult | null>;
  uploadBannerImage: (file: File, entityId?: string) => Promise<FileUploadResult | null>;
  uploadAvatar: (file: File, userId?: string) => Promise<FileUploadResult | null>;

  // Utilities
  reset: () => void;
  validateFile: (file: File) => { isValid: boolean; error?: string };
}

export const useFileUpload = (options: UseFileUploadOptions = {}): UseFileUploadReturn => {
  const [isUploading, setIsUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<FileUploadResult | null>(null);

  // Reset all states
  const reset = useCallback(() => {
    setIsUploading(false);
    setProgress(0);
    setError(null);
    setResult(null);
  }, []);

  // File validation
  const validateFile = useCallback((file: File) => {
    return fileService.validateFile(file, {
      maxSize: options.maxSize,
      allowedTypes: options.allowedTypes
    });
  }, [options.maxSize, options.allowedTypes]);

  // Progress handler
  const handleProgress = useCallback((event: FileUploadProgressEvent) => {
    setProgress(event.progress);
    options.onProgress?.(event);
  }, [options.onProgress]);

  // Generic upload function
  const uploadFile = useCallback(async (uploadOptions: FileUploadOptions): Promise<FileUploadResult | null> => {
    try {
      setIsUploading(true);
      setError(null);
      setProgress(0);

      // Validate file if enabled
      if (options.validateFile !== false) {
        const validation = validateFile(uploadOptions.file);
        if (!validation.isValid) {
          throw new Error(validation.error);
        }
      }

      // Upload file
      const uploadResult = await fileService.uploadFile(uploadOptions, handleProgress);

      setResult(uploadResult);

      if (uploadResult.success) {
        options.onSuccess?.(uploadResult);
      } else {
        const errorMessage = uploadResult.message || 'Upload failed';
        setError(errorMessage);
        options.onError?.(new Error(errorMessage));
      }

      return uploadResult;

    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Upload failed';
      setError(errorMessage);
      options.onError?.(err instanceof Error ? err : new Error(errorMessage));
      return null;
    } finally {
      setIsUploading(false);
    }
  }, [options, validateFile, handleProgress]);

  // Upload with custom path
  const uploadWithCustomPath = useCallback(async (
    file: File,
    customPath: string
  ): Promise<FileUploadResult | null> => {
    return uploadFile({
      file,
      custom_path: customPath
    });
  }, [uploadFile]);

  // Upload with structured path
  const uploadWithStructuredPath = useCallback(async (
    file: File,
    structuredOptions: Omit<StructuredPathUploadOptions, 'file'> = {}
  ): Promise<FileUploadResult | null> => {
    return uploadFile({
      file,
      ...structuredOptions
    });
  }, [uploadFile]);

  // Upload city image
  const uploadCityImage = useCallback(async (
    file: File,
    cityId?: string
  ): Promise<FileUploadResult | null> => {
    return uploadWithStructuredPath(file, {
      category: 'cities',
      entity_id: cityId,
      file_type: 'image'
    });
  }, [uploadWithStructuredPath]);

  // Upload banner image
  const uploadBannerImage = useCallback(async (
    file: File,
    entityId?: string
  ): Promise<FileUploadResult | null> => {
    return uploadWithStructuredPath(file, {
      category: 'banners',
      entity_id: entityId,
      file_type: 'image'
    });
  }, [uploadWithStructuredPath]);

  // Upload avatar
  const uploadAvatar = useCallback(async (
    file: File,
    userId?: string
  ): Promise<FileUploadResult | null> => {
    return uploadWithStructuredPath(file, {
      category: 'profiles',
      entity_id: userId,
      file_type: 'avatar'
    });
  }, [uploadWithStructuredPath]);

  return {
    // States
    isUploading,
    progress,
    error,
    result,

    // Actions
    uploadFile,
    uploadWithCustomPath,
    uploadWithStructuredPath,
    uploadCityImage,
    uploadBannerImage,
    uploadAvatar,

    // Utilities
    reset,
    validateFile,
  };
};

export default useFileUpload;