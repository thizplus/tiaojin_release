// API Key Request Types

export interface CreateApiKeyRequest {
  name: string
  description: string
  expires_at?: string | null
  ip_whitelist?: string[]
  permissions?: string[]
}

export interface UpdateApiKeyRequest {
  name?: string
  description?: string
  is_active?: boolean
  expires_at?: string | null
  ip_whitelist?: string[]
  permissions?: string[]
}

export interface ApiKeyListRequest {
  page?: number
  limit?: number
  is_active?: boolean
  search?: string
  show_expired?: boolean
}

export interface ValidateApiKeyRequest {
  api_key: string
}