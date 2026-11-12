import { useCallback, useState } from 'react';
import type { City } from '@/types';
import { CitySearchAndTable } from '@/components/manage/city/city-search-and-table';
import { CreateCityDialog } from '@/components/manage/city/CreateCityDialog';
import { EditCityDialog } from '@/components/manage/city/EditCityDialog';
import { DeleteCityDialog } from '@/components/manage/city/DeleteCityDialog';

export default function CitiesPage() {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedCity, setSelectedCity] = useState<City | null>(null);

  const handleEdit = useCallback((city: City) => {
    setSelectedCity(city);
    setIsEditDialogOpen(true);
  }, []);

  const handleDelete = useCallback((city: City) => {
    setSelectedCity(city);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleAdd = useCallback(() => {
    setIsCreateDialogOpen(true);
  }, []);

  const handleBulkDelete = useCallback((selectedIds: string[]) => {
    console.log('Bulk delete requested for:', selectedIds);
    setSelectedRows(selectedIds);
    // TODO: Implement bulk delete dialog
  }, []);

  const handleRowSelectionChange = useCallback((selectedIds: string[]) => {
    setSelectedRows(selectedIds);
  }, []);

  const handleCreateSuccess = useCallback(() => {
    console.log('City created successfully');
  }, []);

  const handleEditSuccess = useCallback(() => {
    console.log('City updated successfully');
    setSelectedCity(null);
  }, []);

  const handleDeleteSuccess = useCallback(() => {
    console.log('City deleted successfully');
    setSelectedCity(null);
  }, []);

  return (
    <div className="p-6">
      <CitySearchAndTable
        title="จัดการเมือง"
        description="จัดการข้อมูลเมืองต่างๆ สำหรับระบบ"
        onEdit={handleEdit}
        onDelete={handleDelete}
        onBulkDelete={handleBulkDelete}
        onAdd={handleAdd}
        selectedRows={selectedRows}
        onRowSelectionChange={handleRowSelectionChange}
      />

      <CreateCityDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={handleCreateSuccess}
      />

      <EditCityDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        city={selectedCity}
        onSuccess={handleEditSuccess}
      />

      <DeleteCityDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        city={selectedCity}
        onSuccess={handleDeleteSuccess}
      />
    </div>
  );
}