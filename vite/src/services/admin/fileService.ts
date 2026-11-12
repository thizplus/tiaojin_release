import API from '@/constants/api';
import apiService from '@/services/apiService';
import type {
  FileUploadOptions,
  FileUploadResult,
  FileUploadProgressEvent,
  CustomPathUploadOptions,
  StructuredPathUploadOptions
} from '@/types/file';

class FileService {
  /**
   * Check if upload options are for custom path
   */
  private isCustomPathUpload(options: FileUploadOptions): options is CustomPathUploadOptions {
    return 'custom_path' in options;
  }

  /**
   * Upload file with progress tracking
   */
  async uploadFile(
    options: FileUploadOptions,
    onProgress?: (event: FileUploadProgressEvent) => void
  ): Promise<FileUploadResult> {
    try {
      // Create FormData
      const formData = new FormData();
      formData.append('file', options.file);

      if (this.isCustomPathUpload(options)) {
        // Custom Path Upload
        formData.append('custom_path', options.custom_path);
      } else {
        // Structured Path Upload
        if (options.category) {
          formData.append('category', options.category);
        }
        if (options.entity_id) {
          formData.append('entity_id', options.entity_id);
        }
        if (options.file_type) {
          formData.append('file_type', options.file_type);
        }
      }

      // Convert progress callback format
      const progressCallback = onProgress ? (percentage: number) => {
        // Calculate approximate loaded/total based on percentage
        const total = options.file.size || 1;
        const loaded = Math.round((percentage / 100) * total);

        onProgress({
          loaded,
          total,
          progress: percentage
        });
      } : undefined;

      // Use apiService.upload instead of XMLHttpRequest
      const result = await apiService.upload<FileUploadResult>(
        API.FILE.UPLOAD,
        formData,
        progressCallback
      );

      return result;

    } catch (error) {
      console.error('File upload error:', error);
      throw error;
    }
  }

  /**
   * Upload file using custom path
   */
  async uploadWithCustomPath(
    file: File,
    customPath: string,
    onProgress?: (event: FileUploadProgressEvent) => void
  ): Promise<FileUploadResult> {
    return this.uploadFile({
      file,
      custom_path: customPath
    }, onProgress);
  }

  /**
   * Upload file using structured path
   */
  async uploadWithStructuredPath(
    file: File,
    options: Omit<StructuredPathUploadOptions, 'file'> = {},
    onProgress?: (event: FileUploadProgressEvent) => void
  ): Promise<FileUploadResult> {
    return this.uploadFile({
      file,
      ...options
    }, onProgress);
  }

  /**
   * Upload image for city
   */
  async uploadCityImage(
    file: File,
    cityId?: string,
    onProgress?: (event: FileUploadProgressEvent) => void
  ): Promise<FileUploadResult> {
    return this.uploadWithStructuredPath(file, {
      category: 'cities',
      entity_id: cityId,
      file_type: 'image'
    }, onProgress);
  }

  /**
   * Upload banner image
   */
  async uploadBannerImage(
    file: File,
    entityId?: string,
    onProgress?: (event: FileUploadProgressEvent) => void
  ): Promise<FileUploadResult> {
    return this.uploadWithStructuredPath(file, {
      category: 'banners',
      entity_id: entityId,
      file_type: 'image'
    }, onProgress);
  }

  /**
   * Upload profile avatar
   */
  async uploadAvatar(
    file: File,
    userId?: string,
    onProgress?: (event: FileUploadProgressEvent) => void
  ): Promise<FileUploadResult> {
    return this.uploadWithStructuredPath(file, {
      category: 'profiles',
      entity_id: userId,
      file_type: 'avatar'
    }, onProgress);
  }

  /**
   * Validate file before upload
   */
  validateFile(file: File, options?: {
    maxSize?: number; // in bytes
    allowedTypes?: string[]; // mime types
  }): { isValid: boolean; error?: string } {
    // Check if file exists
    if (!file) {
      return { isValid: false, error: 'No file provided' };
    }

    // Check file size
    if (file.size === 0) {
      return { isValid: false, error: 'File is empty' };
    }

    // Check max size (default 10MB)
    const maxSize = options?.maxSize || 10 * 1024 * 1024;
    if (file.size > maxSize) {
      return {
        isValid: false,
        error: `File size exceeds ${Math.round(maxSize / 1024 / 1024)}MB limit`
      };
    }

    // Check allowed types
    if (options?.allowedTypes && !options.allowedTypes.includes(file.type)) {
      return {
        isValid: false,
        error: `File type ${file.type} is not allowed`
      };
    }

    return { isValid: true };
  }
}

export const fileService = new FileService();
export default fileService;