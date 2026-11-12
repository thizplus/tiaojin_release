// src/stores/apiKeyStore.ts
import { create } from 'zustand'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiKeyService from '@/services/apiKeyService'
import type { ApiKey } from '@/types/models/apiKey'
import type {
  CreateApiKeyRequest,
  UpdateApiKeyRequest,
  ApiKeyListRequest
} from '@/types/request/apiKey'

// Zustand Store State
interface ApiKeyState {
  apiKeys: ApiKey[]
  selectedApiKey: ApiKey | null
  setApiKeys: (apiKeys: ApiKey[]) => void
  setSelectedApiKey: (apiKey: ApiKey | null) => void
  clearApiKeyData: () => void
}

// Zustand Store
export const useApiKeyStore = create<ApiKeyState>((set) => ({
  apiKeys: [],
  selectedApiKey: null,
  setApiKeys: (apiKeys) => set({ apiKeys }),
  setSelectedApiKey: (apiKey) => set({ selectedApiKey: apiKey }),
  clearApiKeyData: () => set({ apiKeys: [], selectedApiKey: null })
}))

// React Query Hooks

// Query: Get all API Keys
export const useApiKeysListQuery = (
  params: ApiKeyListRequest = {},
  options: { enabled?: boolean } = {}
) => {
  const { setApiKeys } = useApiKeyStore()

  return useQuery({
    queryKey: ['api-keys', params],
    queryFn: async () => {
      const response = await apiKeyService.getAll(params)
      if (response.success && response.data?.api_keys) {
        setApiKeys(response.data.api_keys)
      }
      return response
    },
    enabled: options.enabled !== false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Query: Get API Key by ID
export const useApiKeyByIdQuery = (
  id: string,
  options: { enabled?: boolean } = {}
) => {
  const { setSelectedApiKey } = useApiKeyStore()

  return useQuery({
    queryKey: ['api-key', id],
    queryFn: async () => {
      const response = await apiKeyService.getById(id)
      if (response.success && response.data) {
        setSelectedApiKey(response.data)
      }
      return response
    },
    enabled: options.enabled !== false && !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Mutation: Create API Key
export const useCreateApiKeyMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: CreateApiKeyRequest) => {
      const response = await apiKeyService.create(data)
      if (!response.success) {
        throw new Error(response.message || 'Failed to create API key')
      }
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['api-keys'] })
    }
  })
}

// Mutation: Update API Key
export const useUpdateApiKeyMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateApiKeyRequest }) => {
      const response = await apiKeyService.update(id, data)
      if (!response.success) {
        throw new Error(response.message || 'Failed to update API key')
      }
      return response.data
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['api-keys'] })
      queryClient.invalidateQueries({ queryKey: ['api-key', variables.id] })
    }
  })
}

// Mutation: Delete API Key
export const useDeleteApiKeyMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await apiKeyService.delete(id)
      if (!response.success) {
        throw new Error(response.message || 'Failed to delete API key')
      }
      return null
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['api-keys'] })
    }
  })
}