import { useState } from 'react';
import { ToursRegularList } from '@/components/tours-regular/ToursRegularList';
import { TourSearchForm } from '@/components/tours-regular/TourSearchForm';
import { CityFilter } from '@/components/tours-regular/CityFilter';
import { BannerCarousel } from '@/components/agent/dashboard/BannerCarousel';
import { usePublicBanners } from '@/hooks/usePublicBanners';
import type { Tour, Period } from '@/types/tours-regular';

export default function AdminToursRegularPage() {
  const [selectedCityId, setSelectedCityId] = useState<string | undefined>();
  const { banners, isLoading: bannersLoading } = usePublicBanners();

  const handleBookPeriod = (period: Period, tour: Tour) => {
    // TODO: Implement booking logic
    console.log('Book period:', period.id, 'for tour:', tour.id);
  };

  return (
    <div className="space-y-6 container max-w-6xl">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-3xl font-bold">ทัวร์ทั้งหมด</h1>
        <TourSearchForm variant="compact" />
      </div>

      {/* Banners */}
      {!bannersLoading && banners.length > 0 && (
        <BannerCarousel banners={banners} />
      )}

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
        userRole="admin" // Admin role for user selection in booking
        initialParams={{
          status: 'publish',
          available_only: true,
          ...(selectedCityId ? { city_id: selectedCityId } : {}),
        }}
      />
    </div>
  );
}