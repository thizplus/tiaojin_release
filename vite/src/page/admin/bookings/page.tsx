import { useCallback, useState } from "react";
import type { BookingResponse } from "@/types";
import { BookingCardList } from "@/components/manage/booking/BookingCardList";
import { DeleteBookingDialog } from "@/components/manage/booking/DeleteBookingDialog";
import { CreatePaymentDialog } from "@/components/manage/booking/CreatePaymentDialog";

export default function BookingsPage() {
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isCreatePaymentDialogOpen, setIsCreatePaymentDialogOpen] = useState(false);
  const [selectedBooking, setSelectedBooking] = useState<BookingResponse | null>(null);

  const handleView = useCallback(() => {
    // BookingCardList จะจัดการ view เอง ผ่าน BookingDetailSheet
  }, []);

  const handleDelete = useCallback((booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleCreatePayment = useCallback((booking: BookingResponse) => {
    setSelectedBooking(booking);
    setIsCreatePaymentDialogOpen(true);
  }, []);

  const handleDeleteSuccess = useCallback(() => {
    console.log("Booking deleted successfully");
    setSelectedBooking(null);
  }, []);

  const handleCreatePaymentSuccess = useCallback(() => {
    console.log("Payment created successfully");
    setSelectedBooking(null);
  }, []);

  return (
    <div className="p-6 space-y-6">
      <BookingCardList
        title="จัดการการจองทัวร์"
        description="จัดการข้อมูลการจองทัวร์และการชำระเงิน"
        onView={handleView}
        onDelete={handleDelete}
        onCreatePayment={handleCreatePayment}
      />


      <DeleteBookingDialog
        booking={selectedBooking}
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onSuccess={handleDeleteSuccess}
      />

      <CreatePaymentDialog
        booking={selectedBooking}
        open={isCreatePaymentDialogOpen}
        onOpenChange={setIsCreatePaymentDialogOpen}
        onSuccess={handleCreatePaymentSuccess}
      />
    </div>
  );
}