import { useState, useCallback, useMemo, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { useDebounce } from '@/hooks/useDebounce';
import {
  Search,
  X,
  Filter
} from 'lucide-react';

interface TourSearchFormProps {
  onSearchChange?: (searchText: string) => void;
  variant?: 'compact' | 'full';
}

export function TourSearchForm({ onSearchChange, variant = 'full' }: TourSearchFormProps) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [searchInput, setSearchInput] = useState(searchParams.get('search') || '');

  // Get current search value from URL
  const currentSearch = useMemo(() => searchParams.get('search') || '', [searchParams]);

  // Debounce search input for auto search
  const debouncedSearchInput = useDebounce(searchInput, 500);

  // Auto search when debounced input changes
  useEffect(() => {
    const newSearchParams = new URLSearchParams(searchParams);

    if (debouncedSearchInput.trim()) {
      newSearchParams.set('search', debouncedSearchInput.trim());
      newSearchParams.delete('page'); // Reset page when searching
    } else {
      newSearchParams.delete('search');
      newSearchParams.delete('page');
    }

    // Only navigate if the URL actually needs to change
    const newSearchString = newSearchParams.toString();
    if (newSearchString !== searchParams.toString()) {
      navigate({ search: newSearchString });
      onSearchChange?.(debouncedSearchInput.trim());
    }
  }, [debouncedSearchInput, searchParams, navigate, onSearchChange]);

  const handleSearch = useCallback(() => {
    const newSearchParams = new URLSearchParams(searchParams);

    if (searchInput.trim()) {
      newSearchParams.set('search', searchInput.trim());
      newSearchParams.delete('page'); // Reset page when searching
    } else {
      newSearchParams.delete('search');
      newSearchParams.delete('page');
    }

    navigate({ search: newSearchParams.toString() });
    onSearchChange?.(searchInput.trim());
  }, [searchInput, searchParams, navigate, onSearchChange]);

  const handleClearSearch = useCallback(() => {
    const newSearchParams = new URLSearchParams(searchParams);
    newSearchParams.delete('search');
    newSearchParams.delete('page');

    setSearchInput('');
    navigate({ search: newSearchParams.toString() });
    onSearchChange?.('');
  }, [searchParams, navigate, onSearchChange]);

  // Sync input with URL when URL changes (e.g., browser back/forward)
  useEffect(() => {
    const urlSearch = searchParams.get('search') || '';
    if (urlSearch !== searchInput) {
      setSearchInput(urlSearch);
    }
  }, [searchParams]); // Remove searchInput from dependencies to prevent infinite loop

  if (variant === 'compact') {
    return (
      <div className="relative w-full sm:w-auto sm:max-w-md">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="ค้นหาทัวร์..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="pl-10 pr-10 h-9 w-full"
        />
        {searchInput && (
          <button
            onClick={() => setSearchInput('')}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
    );
  }

  return (
    <Card>
      <CardContent className="p-6">
        <div className="space-y-4">
          {/* Search Input */}
          <div className="flex gap-2">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="ค้นหาทัวร์ด้วยชื่อ, รหัสทัวร์, จังหวัด หรือรายละเอียด..."
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="pl-10 pr-10"
              />
              {searchInput && (
                <button
                  onClick={() => setSearchInput('')}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>
            <Button onClick={handleSearch} className="px-6">
              <Search className="h-4 w-4 mr-2" />
              ค้นหา
            </Button>
          </div>

          {/* Current Search Display */}
          {currentSearch && (
            <div className="flex items-center gap-2">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">กำลังค้นหา:</span>
              <Badge variant="secondary" className="gap-1">
                {currentSearch}
                <button
                  onClick={handleClearSearch}
                  className="ml-1 hover:bg-muted rounded-full p-0.5"
                >
                  <X className="h-3 w-3" />
                </button>
              </Badge>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleClearSearch}
                className="text-muted-foreground hover:text-foreground"
              >
                ล้างการค้นหา
              </Button>
            </div>
          )}

          {/* Search Instructions */}
          <div className="text-xs text-muted-foreground">
            <p>💡 ค้นหาอัตโนมัติขณะพิมพ์ - รองรับชื่อทัวร์, รหัสทัวร์, ชื่อจังหวัด หรือรายละเอียดต่างๆ</p>
            {debouncedSearchInput !== searchInput && searchInput && (
              <p className="text-primary mt-1">🔍 กำลังค้นหา...</p>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}