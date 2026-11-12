import { useState } from 'react';
import { ToursRegularList } from '@/components/tours-regular/ToursRegularList';
import { TourSearchForm } from '@/components/tours-regular/TourSearchForm';
import { CityFilter } from '@/components/tours-regular/CityFilter';
import type { Tour, Period } from '@/types/tours-regular';

export default function ToursHotPage() {
  const [selectedCityId, setSelectedCityId] = useState<string | undefined>();

  const handleBookPeriod = (period: Period, tour: Tour) => {
    // TODO: Implement booking logic
    console.log('Book period:', period.id, 'for tour:', tour.id);
  };

  return (
    <div className="space-y-6 container max-w-6xl">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-3xl font-bold">ทัวร์ไฟไหม้ พิเศษ/โปรโมชั่น</h1>
        <TourSearchForm variant="compact" />
      </div>

      {/* City Filter */}
      <CityFilter
        selectedCityId={selectedCityId}
        onCityChange={setSelectedCityId}
      />

      {/* Tours List */}
      <ToursRegularList
        key={selectedCityId || 'all'}
        onBookPeriod={handleBookPeriod}
        canBook={true}
        userStatus="approve" // TODO: Get from user context/auth
        initialParams={{
          status: 'publish',
          available_only: true,
          ...(selectedCityId ? { city_id: selectedCityId } : {}),
        }}
        useToursHook="hot"
      />
    </div>
  );
}