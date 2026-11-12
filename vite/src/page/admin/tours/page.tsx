import { useState, useCallback, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { TourSearchAndTable } from "@/components/manage/tour/tour-search-and-table";
import { CreateTourDialog } from "@/components/manage/tour/CreateTourDialog";
import { EditTourDialog } from "@/components/manage/tour/EditTourDialog";
import { DeleteTourDialog } from "@/components/manage/tour/DeleteTourDialog";
import { useTours } from "@/hooks/useTours";
import type { Tour } from "@/types";

export default function TourPage() {
  const navigate = useNavigate();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedTour, setSelectedTour] = useState<Tour | null>(null);

  const { updateTourStatus, refetchToursList } = useTours();

  // Use ref to avoid re-renders from function references
  const refetchRef = useRef(refetchToursList);
  const updateStatusRef = useRef(updateTourStatus);

  useEffect(() => {
    refetchRef.current = refetchToursList;
    updateStatusRef.current = updateTourStatus;
  }, [refetchToursList, updateTourStatus]);

  const handleAdd = useCallback(() => {
    setCreateDialogOpen(true);
  }, []);

  const handleEdit = useCallback((tour: Tour) => {
    setSelectedTour(tour);
    setEditDialogOpen(true);
  }, []);

  const handleDelete = useCallback((tour: Tour) => {
    setSelectedTour(tour);
    setDeleteDialogOpen(true);
  }, []);

  const handleUpdateStatus = useCallback(async (tour: Tour) => {
    try {
      const newStatus = tour.status === "publish" ? "draft" : "publish";
      await updateStatusRef.current(tour.id, { status: newStatus });
      refetchRef.current();
    } catch (error) {
      console.error("Error updating tour status:", error);
    }
  }, []);

  const handleManagePeriods = useCallback((tour: Tour) => {
    navigate(`/admin/tours/${tour.id}/periods`);
  }, [navigate]);

  const handleDialogSuccess = useCallback(() => {
    refetchRef.current();
  }, []);

  return (
    <div className="p-6">
      <TourSearchAndTable
        onAdd={handleAdd}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onUpdateStatus={handleUpdateStatus}
        onManagePeriods={handleManagePeriods}
      />

      <CreateTourDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onSuccess={handleDialogSuccess}
      />

      <EditTourDialog
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        tour={selectedTour}
        onSuccess={handleDialogSuccess}
      />

      <DeleteTourDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        tour={selectedTour}
        onSuccess={handleDialogSuccess}
      />
    </div>
  );
}