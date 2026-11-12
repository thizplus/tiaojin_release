import type { BannerPosition } from '../models/banner';

// Banner Request Types
export interface CreateBannerRequest {
  title: string;
  image: string;
  link_url?: string;
  alt?: string;
  display_order?: number;
  position?: BannerPosition;
  status: 'publish' | 'draft';
}

export interface UpdateBannerRequest {
  title?: string;
  image?: string;
  link_url?: string;
  alt?: string;
  display_order?: number;
  position?: BannerPosition;
  status?: 'publish' | 'draft';
}

export interface UpdateBannerStatusRequest {
  banner_id: string;
  status: 'publish' | 'draft';
}

export interface UpdateBannerOrderRequest {
  display_order: number;
  banner_orders?: Array<{
    banner_id: string;
    order: number;
  }>;
}