export type BannerPosition = 'public' | 'internal';

export interface Banner {
  id: string;
  title: string;
  image: string;
  link_url: string | null;
  alt: string | null;
  display_order: number;
  position: BannerPosition;
  status: string;
  created_at: string;
  updated_at: string;
}