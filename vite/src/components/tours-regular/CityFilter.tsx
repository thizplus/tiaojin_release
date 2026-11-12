import { useRef, useCallback } from 'react';
import { useCities } from '@/hooks/useCities';
import { cn, getStorageUrl } from '@/lib/utils';
import { X } from 'lucide-react';

interface CityFilterProps {
  selectedCityId?: string;
  onCityChange: (cityId: string | undefined) => void;
}

export function CityFilter({ selectedCityId, onCityChange }: CityFilterProps) {
  const { cities } = useCities({ status: 'publish', stat: true });
  const scrollRef = useRef<HTMLDivElement>(null);

  // scroll ตัวที่กดไปอยู่ตรงกลาง
  const scrollToCenter = useCallback((el: HTMLButtonElement | null) => {
    if (!el || !scrollRef.current) return;
    const container = scrollRef.current;
    const scrollLeft = el.offsetLeft - container.offsetWidth / 2 + el.offsetWidth / 2;
    container.scrollTo({ left: scrollLeft, behavior: 'smooth' });
  }, []);

  if (cities.length === 0) return null;

  const handleClick = (cityId: string | undefined, e: React.MouseEvent<HTMLButtonElement>) => {
    onCityChange(cityId);
    scrollToCenter(e.currentTarget);
  };

  return (
    <div ref={scrollRef} className="flex items-center gap-2 overflow-x-auto pb-1 scrollbar-thin">
      {/* ปุ่ม "ทั้งหมด" */}
      <button
        onClick={(e) => handleClick(undefined, e)}
        className={cn(
          "flex-shrink-0 flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-medium transition-all border",
          !selectedCityId
            ? "bg-primary text-primary-foreground border-primary"
            : "bg-background text-muted-foreground border-border hover:border-foreground/30"
        )}
      >
        ทั้งหมด
      </button>

      {cities.map((city) => (
        <button
          key={city.id}
          onClick={(e) => handleClick(selectedCityId === city.id ? undefined : city.id, e)}
          className={cn(
            "flex-shrink-0 flex items-center gap-1.5 rounded-full pl-1 pr-3 py-1 text-xs font-medium transition-all border",
            selectedCityId === city.id
              ? "bg-primary text-primary-foreground border-primary"
              : "bg-background text-foreground border-border hover:border-foreground/30"
          )}
        >
          <img
            src={getStorageUrl(city.image)}
            alt={city.title}
            className="w-6 h-6 rounded-full object-cover"
            onError={(e) => { e.currentTarget.src = '/no-image.png' }}
          />
          <span className="whitespace-nowrap">{city.title}</span>
          {selectedCityId === city.id && (
            <X className="h-3 w-3 ml-0.5" />
          )}
        </button>
      ))}
    </div>
  );
}
