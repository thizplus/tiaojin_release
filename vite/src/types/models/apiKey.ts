// API Key Model Types

export interface ApiKey {
  id: string
  name: string
  description: string
  key_preview: string
  is_active: boolean
  last_used_at: string | null
  expires_at: string | null
  ip_whitelist: string[]
  permissions: string[]
  created_at: string
  updated_at: string
}

export interface ApiKeyWithFullKey extends Omit<ApiKey, 'key_preview'> {
  api_key: string
}