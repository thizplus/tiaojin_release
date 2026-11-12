import type { City } from '../models/city';
import type { ApiResponse } from '../common';

export interface CitiesListResponse extends ApiResponse<{
  cities: City[];
  pagination: {
    total: number;
    has_more: boolean;
    current_page: number;
    total_pages: number;
    page_size: number;
    limit: number;
    offset: number;
  };
}> {}

export interface CityResponse extends ApiResponse<City> {}