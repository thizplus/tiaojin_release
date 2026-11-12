import type { Banner } from '../models/banner';
import type { ApiResponse } from '../common';

export interface PublicBannersResponse extends ApiResponse<Banner[]> {}