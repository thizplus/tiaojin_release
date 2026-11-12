import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { MapPin, Calendar, Clock } from 'lucide-react'
import type { City } from '@/types/models/city'

interface CityCardProps {
  city: City
  onClick: (cityId: string) => void
}

export function CityCard({ city, onClick }: CityCardProps) {
  return (
    <Card
      className="cursor-pointer transition-all hover:shadow-lg hover:scale-105 pt-0"
      onClick={() => onClick(city.id)}
    >
      <div className="w-full overflow-hidden rounded-t-lg">
        <img
          src={city.image}
          alt={city.title}
          className="h-full w-full"
          onError={(e) => {
            e.currentTarget.src = '/no-image.png'
          }}
        />
      </div>
      <CardHeader className="px-4">
        <div className="flex items-center gap-2">
          <MapPin className="h-5 w-5 text-primary" />
          <CardTitle className="text-xl">{city.title}</CardTitle>
        </div>
        {city.detail && (
          <CardDescription className="line-clamp-2">
            {city.detail}
          </CardDescription>
        )}
        <div className="flex flex-wrap gap-2 pt-2">
          <Badge variant="secondary" className="flex items-center gap-1">
            <Calendar className="h-3 w-3" />
            {city.active_tours || 0} โปรแกรมทัวร์
          </Badge>
          <Badge variant="secondary" className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {city.active_periods || 0} รายการ
          </Badge>
        </div>
      </CardHeader>
    </Card>
  )
}