import { useCallback } from 'react'
import {
  useApiKeyStore,
  useApiKeysListQuery,
  useApiKeyByIdQuery,
  useCreateApiKeyMutation,
  useUpdateApiKeyMutation,
  useDeleteApiKeyMutation
} from '@/stores/apiKeyStore'
import type { ApiKey } from '@/types/models/apiKey'
import type {
  CreateApiKeyRequest,
  UpdateApiKeyRequest,
  ApiKeyListRequest
} from '@/types/request/apiKey'

interface UseApiKeyReturn {
  // State
  apiKeys: ApiKey[]
  selectedApiKey: ApiKey | null
  isLoading: boolean

  // Actions
  createApiKey: (data: CreateApiKeyRequest) => Promise<any>
  updateApiKey: (id: string, data: UpdateApiKeyRequest) => Promise<ApiKey>
  deleteApiKey: (id: string) => Promise<null>

  // State management
  setApiKeys: (apiKeys: ApiKey[]) => void
  setSelectedApiKey: (apiKey: ApiKey | null) => void
  clearApiKeyData: () => void

  // Queries
  refetchApiKeysList: () => void
  refetchApiKeyById: () => void

  // Loading states
  isCreateLoading: boolean
  isUpdateLoading: boolean
  isDeleteLoading: boolean
  isApiKeysListLoading: boolean
  isApiKeyByIdLoading: boolean
}

export const useApiKey = (
  listParams?: ApiKeyListRequest,
  apiKeyId?: string,
  options: { enableList?: boolean; enableById?: boolean } = {}
): UseApiKeyReturn => {
  // Store state
  const {
    apiKeys,
    selectedApiKey,
    setApiKeys,
    setSelectedApiKey,
    clearApiKeyData
  } = useApiKeyStore()

  // Mutations
  const createApiKeyMutation = useCreateApiKeyMutation()
  const updateApiKeyMutation = useUpdateApiKeyMutation()
  const deleteApiKeyMutation = useDeleteApiKeyMutation()

  // Queries
  const {
    refetch: refetchApiKeysList,
    isLoading: isApiKeysListLoading
  } = useApiKeysListQuery(listParams || {}, {
    enabled: options.enableList === true
  })

  const {
    refetch: refetchApiKeyById,
    isLoading: isApiKeyByIdLoading
  } = useApiKeyByIdQuery(apiKeyId || '', {
    enabled: options.enableById !== false && !!apiKeyId
  })

  // Actions
  const createApiKey = useCallback(async (data: CreateApiKeyRequest) => {
    const result = await createApiKeyMutation.mutateAsync(data)
    return result
  }, [createApiKeyMutation])

  const updateApiKey = useCallback(async (id: string, data: UpdateApiKeyRequest) => {
    const result = await updateApiKeyMutation.mutateAsync({ id, data })
    return result
  }, [updateApiKeyMutation])

  const deleteApiKey = useCallback(async (id: string) => {
    const result = await deleteApiKeyMutation.mutateAsync(id)
    return result
  }, [deleteApiKeyMutation])

  // Calculate overall loading state
  const isLoading = createApiKeyMutation.isPending ||
                   updateApiKeyMutation.isPending ||
                   deleteApiKeyMutation.isPending ||
                   isApiKeysListLoading ||
                   isApiKeyByIdLoading

  return {
    // State
    apiKeys,
    selectedApiKey,
    isLoading,

    // Actions
    createApiKey,
    updateApiKey,
    deleteApiKey,

    // State management
    setApiKeys,
    setSelectedApiKey,
    clearApiKeyData,

    // Queries
    refetchApiKeysList,
    refetchApiKeyById,

    // Loading states
    isCreateLoading: createApiKeyMutation.isPending,
    isUpdateLoading: updateApiKeyMutation.isPending,
    isDeleteLoading: deleteApiKeyMutation.isPending,
    isApiKeysListLoading,
    isApiKeyByIdLoading
  }
}

export default useApiKey