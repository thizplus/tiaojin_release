export interface City {
  id: string;
  title: string;
  detail: string | null;
  image: string;
  status: string;
  created_at: string;
  updated_at: string;
  active_tours?: number;
  active_periods?: number;
}