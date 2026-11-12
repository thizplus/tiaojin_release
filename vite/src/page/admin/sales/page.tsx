import { useCallback, useState } from "react";
import type { Sale } from "@/types";
import { SalesSearchAndTable } from "@/components/manage/sales/sales-search-and-table";
import { CreateSaleDialog } from "@/components/manage/sales/CreateSaleDialog";
import { EditSaleDialog } from "@/components/manage/sales/EditSaleDialog";
import { DeleteSaleDialog } from "@/components/manage/sales/DeleteSaleDialog";
import { useSales } from "@/hooks/useSales";

export default function SalesPage() {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedSale, setSelectedSale] = useState<Sale | null>(null);

  const { activateSale, deactivateSale, refetchSalesList } = useSales();

  const handleEdit = useCallback((sale: Sale) => {
    setSelectedSale(sale);
    setIsEditDialogOpen(true);
  }, []);

  const handleDelete = useCallback((sale: Sale) => {
    setSelectedSale(sale);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleActivate = useCallback(async (sale: Sale) => {
    try {
      await activateSale(sale.id);
      refetchSalesList();
      console.log("Sale activated successfully");
    } catch (error) {
      console.error("Failed to activate sale:", error);
    }
  }, [activateSale, refetchSalesList]);

  const handleDeactivate = useCallback(async (sale: Sale) => {
    try {
      await deactivateSale(sale.id);
      refetchSalesList();
      console.log("Sale deactivated successfully");
    } catch (error) {
      console.error("Failed to deactivate sale:", error);
    }
  }, [deactivateSale, refetchSalesList]);

  const handleAdd = useCallback(() => {
    setIsCreateDialogOpen(true);
  }, []);

  const handleCreateSuccess = useCallback(() => {
    console.log("Sale created successfully");
  }, []);

  const handleEditSuccess = useCallback(() => {
    console.log("Sale updated successfully");
    setSelectedSale(null);
  }, []);

  const handleDeleteSuccess = useCallback(() => {
    console.log("Sale deleted successfully");
    setSelectedSale(null);
  }, []);

  return (
    <div className="p-6">
      <SalesSearchAndTable
        title="จัดการพนักงานขาย"
        description="จัดการข้อมูลพนักงานขายและสถานะการใช้งาน"
        onEdit={handleEdit}
        onDelete={handleDelete}
        onActivate={handleActivate}
        onDeactivate={handleDeactivate}
        onAdd={handleAdd}
      />

      <CreateSaleDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={handleCreateSuccess}
      />

      <EditSaleDialog
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        sale={selectedSale}
        onSuccess={handleEditSuccess}
      />

      <DeleteSaleDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        sale={selectedSale}
        onSuccess={handleDeleteSuccess}
      />
    </div>
  );
}