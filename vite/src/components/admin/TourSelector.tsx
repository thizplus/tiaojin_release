import { useState, useEffect, useRef } from 'react';
import { Search, MapPin } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Loader2 } from 'lucide-react';
import { useTourSearch } from '@/hooks/useTourSearch';
import type { Tour } from '@/types';

interface TourSelectorProps {
  selectedTour: Tour | null;
  onTourSelect: (tour: Tour | null) => void;
}

export function TourSelector({ selectedTour, onTourSelect }: TourSelectorProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const {
    tours,
    searchTerm,
    setSearchTerm,
    isLoading,
    clearSearch,
  } = useTourSearch({
    limit: 10,
    debounceMs: 300,
    autoSearch: true,
  });

  // จัดการ click outside เพื่อปิด dropdown
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }
  }, [isOpen]);

  const handleTourSelect = (tour: Tour) => {
    onTourSelect(tour);
    setIsOpen(false);
    clearSearch();
  };

  const handleClear = () => {
    onTourSelect(null);
    clearSearch();
  };

  return (
    <div className="space-y-2">
      <label className="text-sm font-medium">เลือกทัวร์</label>
      <div className="relative" ref={dropdownRef}>
        <div className="flex gap-2">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="ค้นหาทัวร์ด้วยชื่อหรือรหัส..."
              value={selectedTour ? `${selectedTour.code} - ${selectedTour.title}` : searchTerm}
              onChange={(e) => {
                if (!selectedTour) {
                  setSearchTerm(e.target.value);
                  setIsOpen(true);
                }
              }}
              onFocus={() => setIsOpen(true)}
              className="pl-10"
              readOnly={!!selectedTour}
            />
          </div>
          {selectedTour && (
            <Button variant="outline" onClick={handleClear}>
              เคลียร์
            </Button>
          )}
        </div>

        {/* Dropdown รายการทัวร์ */}
        {isOpen && !selectedTour && (
          <div className="absolute top-full left-0 right-0 z-10 mt-1 bg-background border rounded-md shadow-lg">
            <ScrollArea className="max-h-64">
              {isLoading ? (
                <div className="flex items-center justify-center p-4">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span className="ml-2 text-sm text-muted-foreground">กำลังโหลด...</span>
                </div>
              ) : tours.length === 0 ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  {searchTerm ? 'ไม่พบทัวร์ที่ค้นหา' : 'พิมพ์เพื่อค้นหาทัวร์'}
                </div>
              ) : (
                <div className="p-2">
                  {tours.map((tour) => (
                    <div
                      key={tour.id}
                      className="flex items-center gap-3 p-3 hover:bg-muted cursor-pointer rounded"
                      onClick={() => handleTourSelect(tour)}
                    >
                      <MapPin className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs">
                            {tour.code}
                          </Badge>
                        </div>
                        <div className="text-sm font-medium truncate">
                          {tour.title}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </ScrollArea>
          </div>
        )}
      </div>
    </div>
  );
}