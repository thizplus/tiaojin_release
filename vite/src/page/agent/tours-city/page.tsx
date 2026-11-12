import { useParams } from 'react-router-dom';
import { ToursRegularList } from '@/components/tours-regular/ToursRegularList';
import { TourSearchForm } from '@/components/tours-regular/TourSearchForm';
import type { Tour, Period } from '@/types/tours-regular';

export default function ToursCityPage() {
  const { cityId } = useParams<{ cityId: string }>();

  const handleBookPeriod = (period: Period, tour: Tour) => {
    // TODO: Implement booking logic
    console.log('Book period:', period.id, 'for tour:', tour.id);
    // This could open a booking dialog or navigate to booking page
  };

  return (
    <div className="space-y-6 container max-w-6xl">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-3xl font-bold">ทัวร์ทั้งหมด</h1>
        <TourSearchForm variant="compact" />
      </div>

      {/* Tours List */}
      <ToursRegularList
        onBookPeriod={handleBookPeriod}
        canBook={true}
        userStatus="approve" // TODO: Get from user context/auth
        initialParams={{
          status: 'publish',
          city_id: cityId,
          available_only: true,
        }}
      />
    </div>
  );
}