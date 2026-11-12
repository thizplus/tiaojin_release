export interface CityFilter {
    search?: string;
    status?: string;
    stat?: boolean;
    page?: number;
    limit?: number;
    sort?: string;
    order?: string;
}

export interface BannerFilter {
    search?: string;
    status?: string;
    position?: string;
    is_active?: boolean;
    page?: number;
    limit?: number;
    sort?: string;
    order?: string;
}

export interface UserFilter {
    search?: string;
    role?: string;
    status?: string;
    page?: number;
    limit?: number;
    sort?: string;
    order?: string;
}