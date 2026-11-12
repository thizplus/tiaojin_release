export type TourStatus = 'draft' | 'publish';

export interface Tour {
  id: string;
  title: string;
  code: string;
  city_id: string;
  city_title?: string;
  detail: string;
  thumbnail: string;
  file_pdf?: string | null;
  file_word?: string | null;
  status: TourStatus;
  created_at: string;
  updated_at: string;
}

export interface TourFilter {
  page?: number;
  limit?: number;
  search?: string;
  status?: TourStatus;
  city_id?: string;
  sort?: string;
  order?: 'asc' | 'desc';
  with_periods?: boolean;
}

export interface CreateTourRequest {
  title: string;
  code: string;
  city_id: string;
  detail: string;
  thumbnail: string;
  file_pdf?: string | null;
  file_word?: string | null;
  status?: TourStatus;
}

export interface UpdateTourRequest {
  title?: string;
  city_id?: string;
  detail?: string;
  thumbnail?: string;
  file_pdf?: string | null;
  file_word?: string | null;
  status?: TourStatus;
}

export interface UpdateTourStatusRequest {
  status: TourStatus;
}