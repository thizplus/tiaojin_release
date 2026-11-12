// City Request Types
export interface CreateCityRequest {
    title: string;
    detail?: string;
    image: string;
    status: 'publish' | 'draft';
  }
  
  export interface UpdateCityRequest {
    title?: string;
    detail?: string;
    image?: string;
    status?: 'publish' | 'draft';
  }
  
  export interface UpdateCityStatusRequest {
    city_id: string;
    status: 'publish' | 'draft';
  }