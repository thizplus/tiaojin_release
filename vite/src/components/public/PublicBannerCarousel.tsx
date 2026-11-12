import { useState, useEffect, useCallback } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { getStorageUrl } from '@/lib/utils';
import type { Banner } from '@/types';

interface PublicBannerCarouselProps {
  banners: Banner[];
  autoSlideInterval?: number;
}

/** Banner Carousel สำหรับหน้า public — landscape, auto-slide, รองรับหลาย banner */
export function PublicBannerCarousel({
  banners,
  autoSlideInterval = 5000,
}: PublicBannerCarouselProps) {
  const [current, setCurrent] = useState(0);

  const next = useCallback(() => {
    setCurrent((prev) => (prev + 1) % banners.length);
  }, [banners.length]);

  const prev = useCallback(() => {
    setCurrent((prev) => (prev - 1 + banners.length) % banners.length);
  }, [banners.length]);

  // Auto-slide
  useEffect(() => {
    if (banners.length <= 1) return;
    const timer = setInterval(next, autoSlideInterval);
    return () => clearInterval(timer);
  }, [banners.length, next, autoSlideInterval]);

  if (banners.length === 0) return null;

  const banner = banners[current];

  return (
    <div className="group relative overflow-hidden rounded-lg">
      <a
        href={banner.link_url || undefined}
        target={banner.link_url ? '_blank' : undefined}
        rel={banner.link_url ? 'noopener noreferrer' : undefined}
        className="block"
      >
        <img
          src={getStorageUrl(banner.image)}
          alt={banner.alt || banner.title}
          className="w-full h-auto object-cover"
        />
      </a>

      {/* Arrows — ซ่อนใน mobile, จางๆ แสดงตอน hover ใน desktop */}
      {banners.length > 1 && (
        <>
          <button
            onClick={prev}
            className="absolute left-3 top-1/2 -translate-y-1/2 rounded-full p-1.5 bg-black/20 text-white opacity-0 transition-opacity hover:bg-black/40 md:group-hover:opacity-100"
          >
            <ChevronLeft className="h-5 w-5" />
          </button>
          <button
            onClick={next}
            className="absolute right-3 top-1/2 -translate-y-1/2 rounded-full p-1.5 bg-black/20 text-white opacity-0 transition-opacity hover:bg-black/40 md:group-hover:opacity-100"
          >
            <ChevronRight className="h-5 w-5" />
          </button>

          {/* Dots */}
          <div className="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-2">
            {banners.map((_, i) => (
              <button
                key={i}
                onClick={() => setCurrent(i)}
                className={`h-2 rounded-full transition-all ${i === current ? 'w-6 bg-white' : 'w-2 bg-white/50'}`}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
