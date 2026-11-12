import { memo, useState } from 'react';
import { useApiKeysData } from '@/hooks/useApiKeysData';
import { Button } from '@/components/ui/button';
import { ApiKeyTable } from './api-key-table';
import { ApiDocumentationDialog } from './ApiDocumentationDialog';
import type { ApiKey } from '@/types/models/apiKey';
import { BookOpen } from 'lucide-react';

interface ApiKeySearchAndTableProps {
  title?: string;
  description?: string;
  onEdit: (apiKey: ApiKey) => void;
  onDelete: (apiKey: ApiKey) => void;
  onAdd: () => void;
  selectedRows?: string[];
  onRowSelectionChange?: (selectedRows: string[]) => void;
}

export const ApiKeySearchAndTable = memo(function ApiKeySearchAndTable({
  title = "จัดการ API Keys",
  description = "จัดการ API Keys สำหรับเชื่อมต่อกับระบบ",
  onEdit,
  onDelete,
  onAdd,
  selectedRows = [],
  onRowSelectionChange
}: ApiKeySearchAndTableProps) {
  const [isDocDialogOpen, setIsDocDialogOpen] = useState(false);

  // ใช้ useApiKeysData hook ที่จัดการทุกอย่าง
  const {
    apiKeys,
    pagination,
    isLoading,
    currentPage,
    currentPageSize,
    handlePageChange,
    handlePageSizeChange,
  } = useApiKeysData();

  return (
    <div className="space-y-6">
     <div className="flex flex-col items-left gap-4 md:flex-row md:items-center justify-between mb-6 ">
         {/* Header Section */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">{title}</h1>
          <p className="text-muted-foreground mt-1">{description}</p>
        </div>
      </div>

      {/* API Documentation Button */}
      <div className="relative w-full max-w-sm">
        <Button
          variant="outline"
          onClick={() => setIsDocDialogOpen(true)}
          className="w-full"
        >
          <BookOpen className="mr-2 h-4 w-4" />
          รายละเอียด API
        </Button>
      </div>
      </div>

      {/* Loading State */}
      {isLoading ? (
        <div className="p-6 text-center">กำลังโหลดข้อมูล...</div>
      ) : (
        /* API Key Table */
        <ApiKeyTable
          data={apiKeys}
          onEdit={onEdit}
          onDelete={onDelete}
          onAdd={onAdd}
          page={currentPage - 1} // ปรับให้เริ่มจาก 0 สำหรับ DataTable
          pageSize={currentPageSize}
          totalPages={pagination?.totalPages || 1}
          setPage={handlePageChange}
          setPageSize={handlePageSizeChange}
          selectedRows={selectedRows}
          onRowSelectionChange={onRowSelectionChange}
          searchQuery=""
          onSearchChange={() => {}}
        />
      )}

      {/* API Documentation Dialog */}
      <ApiDocumentationDialog
        open={isDocDialogOpen}
        onOpenChange={setIsDocDialogOpen}
      />
    </div>
  );
});

export default ApiKeySearchAndTable;