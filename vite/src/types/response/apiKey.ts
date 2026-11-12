// API Key Response Types

import type { ApiKey, ApiKeyWithFullKey } from '../models/apiKey'
import type { ApiResponse } from '../common'

export interface ApiKeyListResponse extends ApiResponse<{
  api_keys: ApiKey[]
  pagination: {
    total: number
    offset: number
    limit: number
  }
}> {}

export interface ApiKeyResponse extends ApiResponse<ApiKey> {}

export interface CreateApiKeyResponse extends ApiResponse<ApiKeyWithFullKey> {}

export interface UpdateApiKeyResponse extends ApiResponse<ApiKey> {}

export interface DeleteApiKeyResponse extends ApiResponse<null> {}

export interface ValidateApiKeyResponse extends ApiResponse<{
  valid: boolean
  user_id?: string
  key_id?: string
  permissions?: string[]
  expires_at?: string | null
}> {}