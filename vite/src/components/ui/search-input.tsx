import React, { useState, useCallback } from 'react';
import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { useDebounce } from '@/hooks/useDebounce';

interface SearchInputProps {
  value: string;
  onSearch: (value: string) => void;
  placeholder?: string;
  debounceMs?: number;
  className?: string;
}

const SearchInput = React.memo(function SearchInput({
  value,
  onSearch,
  placeholder = "ค้นหา...",
  debounceMs = 500,
  className
}: SearchInputProps) {
  const [inputValue, setInputValue] = useState(value);
  const inputRef = React.useRef<HTMLInputElement>(null);
  const wasFocused = React.useRef(false);

  // Debounce the search
  const debouncedSearchTerm = useDebounce(inputValue, debounceMs);

  // Call onSearch when debounced value changes
  React.useEffect(() => {
    if (debouncedSearchTerm !== value) {
      onSearch(debouncedSearchTerm);
    }
  }, [debouncedSearchTerm, onSearch, value]);

  // Update local state when external value changes (only if input is not focused)
  React.useEffect(() => {
    if (value !== inputValue && !wasFocused.current) {
      setInputValue(value);
    }
  }, [value]);

  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  }, []);

  const handleFocus = useCallback(() => {
    wasFocused.current = true;
  }, []);

  const handleBlur = useCallback(() => {
    wasFocused.current = false;
  }, []);

  return (
    <div className={`relative ${className || ''}`}>
      <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
      <Input
        ref={inputRef}
        placeholder={placeholder}
        value={inputValue}
        onChange={handleInputChange}
        onFocus={handleFocus}
        onBlur={handleBlur}
        className="pl-10"
      />
    </div>
  );
});

export { SearchInput };