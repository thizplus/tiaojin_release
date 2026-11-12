import { useState, useEffect } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn, getStorageUrl } from '@/lib/utils'
import type { Banner } from '@/types/models/banner'
import { useIsMobile } from '@/hooks/use-mobile'
import FsLightbox from 'fslightbox-react'

interface BannerCarouselProps {
  banners: Banner[]
}

export function BannerCarousel({ banners }: BannerCarouselProps) {
  const isMobile = useIsMobile()
  const [currentIndex, setCurrentIndex] = useState(0)
  const [lightboxController, setLightboxController] = useState({
    toggler: false,
    slide: 1
  })

  // Auto-play for mobile slider
  useEffect(() => {
    if (!isMobile || banners.length <= 1) return

    const interval = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % banners.length)
    }, 5000) // Change slide every 5 seconds

    return () => clearInterval(interval)
  }, [isMobile, banners.length])

  const handlePrevious = () => {
    setCurrentIndex((prev) => (prev - 1 + banners.length) % banners.length)
  }

  const handleNext = () => {
    setCurrentIndex((prev) => (prev + 1) % banners.length)
  }

  const handleBannerClick = (index: number) => {
    setLightboxController({
      toggler: !lightboxController.toggler,
      slide: index + 1
    })
  }

  if (banners.length === 0) return null

  // Mobile: Slider view
  if (isMobile) {
    return (
      <>
        <div className="relative w-full">
          <div className="relative aspect-[1/1] w-full overflow-hidden rounded-lg">
            {banners.map((banner, index) => (
              <div
                key={banner.id}
                className={cn(
                  'absolute inset-0 transition-opacity duration-500',
                  index === currentIndex ? 'opacity-100' : 'opacity-0 pointer-events-none'
                )}
              >
                <img
                  src={getStorageUrl(banner.image)}
                  alt={banner.alt || banner.title}
                  className="h-full w-full cursor-pointer object-cover"
                  onClick={() => handleBannerClick(index)}
                />
              </div>
            ))}

            {/* Navigation Buttons */}
            {banners.length > 1 && (
              <>
                <Button
                  variant="ghost"
                  size="icon"
                  className="absolute left-2 top-1/2 z-10 -translate-y-1/2 bg-background/80 text-foreground hover:bg-background/90 hover:text-foreground"
                  onClick={handlePrevious}
                >
                  <ChevronLeft className="h-6 w-6" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="absolute right-2 top-1/2 z-10 -translate-y-1/2 bg-background/80 text-foreground hover:bg-background/90 hover:text-foreground"
                  onClick={handleNext}
                >
                  <ChevronRight className="h-6 w-6" />
                </Button>
              </>
            )}
          </div>

          {/* Dots Indicator */}
          {banners.length > 1 && (
            <div className="mt-4 flex justify-center gap-2">
              {banners.map((_, index) => (
                <button
                  key={index}
                  aria-label={`Go to slide ${index + 1}`}
                  className={cn(
                    'h-2 w-2 rounded-full transition-all',
                    index === currentIndex
                      ? 'w-8 bg-primary'
                      : 'bg-muted-foreground/30'
                  )}
                  onClick={() => setCurrentIndex(index)}
                />
              ))}
            </div>
          )}
        </div>

        <FsLightbox
          toggler={lightboxController.toggler}
          sources={banners.map(b => getStorageUrl(b.image))}
          slide={lightboxController.slide}
        />
      </>
    )
  }

  // Desktop: Grid view (4 columns)
  return (
    <>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {banners.map((banner, index) => (
          <div
            key={banner.id}
            className="group relative aspect-[1/1] cursor-pointer overflow-hidden rounded-lg"
            onClick={() => handleBannerClick(index)}
          >
            <img
              src={getStorageUrl(banner.image)}
              alt={banner.alt || banner.title}
              className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
          </div>
        ))}
      </div>

      <FsLightbox
        toggler={lightboxController.toggler}
        sources={banners.map(b => getStorageUrl(b.image))}
        slide={lightboxController.slide}
      />
    </>
  )
}