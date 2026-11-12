import { useNavigate } from 'react-router-dom'
import { useCities } from '@/hooks/useCities'
import { usePublicBanners } from '@/hooks/usePublicBanners'
import { CityCard } from '@/components/agent/dashboard/CityCard'
import { BannerCarousel } from '@/components/agent/dashboard/BannerCarousel'
import { Card, CardContent } from '@/components/ui/card'
import { Loader2 } from 'lucide-react'

export default function AgentDashboardPage() {
  const navigate = useNavigate()

  const { cities, isLoading: citiesLoading, error: citiesError, refetch: refetchCities } = useCities({
    status: 'publish',
    stat: true
  })

  const { banners, isLoading: bannersLoading } = usePublicBanners()

  const handleCityClick = (cityId: string) => {
    navigate(`/agent/cities/${cityId}`)
  }

  return (
    <div className="space-y-8">
      {/* Banners Section */}
      {!bannersLoading && banners.length > 0 && (
        <section>
          <BannerCarousel banners={banners} />
        </section>
      )}

      {/* Cities Section */}
      <section className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">เลือกเมืองที่ต้องการ</h1>
        </div>

        {citiesLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin" />
          </div>
        ) : citiesError ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-red-500">เกิดข้อผิดพลาดในการโหลดข้อมูล</p>
              <button onClick={() => refetchCities()} className="mt-4 px-4 py-2 bg-primary text-white rounded-md">
                ลองอีกครั้ง
              </button>
            </CardContent>
          </Card>
        ) : cities.length > 0 ? (
          <div className="grid grid-cols-2 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
            {cities.map((city) => (
              <CityCard key={city.id} city={city} onClick={handleCityClick} />
            ))}
          </div>
        ) : (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">ไม่มีข้อมูลเมือง</p>
            </CardContent>
          </Card>
        )}
      </section>
    </div>
  )
}