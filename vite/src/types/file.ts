// File Upload Types
export interface FileUploadResponse {
  fileId: string;
  fileName: string;
  url: string;
  cdnPath: string;
  fileSize: number;
  mimeType: string;
  pathType: 'custom' | 'structured';
}

export interface CustomPathUploadOptions {
  file: File;
  custom_path: string;
}

export interface StructuredPathUploadOptions {
  file: File;
  category?: 'banners' | 'cities' | 'tours' | 'periods' | 'profiles';
  entity_id?: string; // UUID
  file_type?: 'image' | 'document' | 'avatar';
}

export type FileUploadOptions = CustomPathUploadOptions | StructuredPathUploadOptions;

export interface FileUploadProgressEvent {
  loaded: number;
  total: number;
  progress: number; // 0-100
}

export interface FileUploadError {
  success: false;
  message: string;
  errors?: string[];
}

export interface FileUploadSuccess {
  success: true;
  message: string;
  data: FileUploadResponse;
}

export type FileUploadResult = FileUploadSuccess | FileUploadError;