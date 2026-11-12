import { useCallback, useState } from 'react';
import type { ApiKey } from '@/types/models/apiKey';
import { ApiKeySearchAndTable } from '@/components/manage/api-key/api-key-search-and-table';
import { CreateApiKeyDialog } from '@/components/manage/api-key/CreateApiKeyDialog';
import { EditApiKeyDialog } from '@/components/manage/api-key/EditApiKeyDialog';
import { DeleteApiKeyDialog } from '@/components/manage/api-key/DeleteApiKeyDialog';

export default function ApiKeysPage() {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedApiKey, setSelectedApiKey] = useState<ApiKey | null>(null);

  const handleEdit = useCallback((apiKey: ApiKey) => {
    setSelectedApiKey(apiKey);
    setIsEditDialogOpen(true);
  }, []);

  const handleDelete = useCallback((apiKey: ApiKey) => {
    setSelectedApiKey(apiKey);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleAdd = useCallback(() => {
    setIsCreateDialogOpen(true);
  }, []);

  const handleRowSelectionChange = useCallback((selectedIds: string[]) => {
    setSelectedRows(selectedIds);
  }, []);

  const handleCreateSuccess = useCallback(() => {
    console.log('API Key created successfully');
  }, []);

  const handleEditSuccess = useCallback(() => {
    console.log('API Key updated successfully');
    setSelectedApiKey(null);
  }, []);

  const handleDeleteSuccess = useCallback(() => {
    console.log('API Key deleted successfully');
    setSelectedApiKey(null);
  }, []);

  return (
    <div className="p-6">
      <ApiKeySearchAndTable
        title="จัดการ API Keys"
        description="จัดการ API Keys สำหรับเชื่อมต่อกับระบบ"
        onEdit={handleEdit}
        onDelete={handleDelete}
        onAdd={handleAdd}
        selectedRows={selectedRows}
        onRowSelectionChange={handleRowSelectionChange}
      />

      <CreateApiKeyDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={handleCreateSuccess}
      />

      <EditApiKeyDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        apiKey={selectedApiKey}
        onSuccess={handleEditSuccess}
      />

      <DeleteApiKeyDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        apiKey={selectedApiKey}
        onSuccess={handleDeleteSuccess}
      />
    </div>
  );
}