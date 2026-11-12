import React, { useState, useCallback, useRef } from 'react';
import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';

// Global state เพื่อเก็บค่า search และ callback
let globalSearchValue = '';
let globalOnSearch: ((value: string) => void) | null = null;
let debounceTimeout: NodeJS.Timeout | null = null;
let shouldMaintainFocus = false;

interface IsolatedSearchInputProps {
  placeholder?: string;
  onSearchChange?: (value: string) => void;
}

// Component ที่แยกออกมาโดยสมบูรณ์
export const IsolatedSearchInput = React.memo(function IsolatedSearchInput({
  placeholder = "ค้นหา...",
  onSearchChange
}: IsolatedSearchInputProps) {
  const [localValue, setLocalValue] = useState(globalSearchValue);
  const inputRef = useRef<HTMLInputElement>(null);

  // ตั้งค่า global callback
  React.useEffect(() => {
    if (onSearchChange) {
      globalOnSearch = onSearchChange;
    }
  }, [onSearchChange]);

  // กู้คืน focus หลังจาก re-render
  React.useEffect(() => {
    if (shouldMaintainFocus && inputRef.current) {
      const input = inputRef.current;
      const length = input.value.length;
      input.focus();
      input.setSelectionRange(length, length); // วาง cursor ที่ท้ายข้อความ
      shouldMaintainFocus = false;
    }
  });

  const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setLocalValue(newValue);
    globalSearchValue = newValue;

    // Clear previous timeout
    if (debounceTimeout) {
      clearTimeout(debounceTimeout);
    }

    // Set new timeout for debounce
    debounceTimeout = setTimeout(() => {
      shouldMaintainFocus = true; // บอกให้กู้คืน focus หลัง search
      if (globalOnSearch) {
        globalOnSearch(newValue);
      }
    }, 300);
  }, []);

  const handleFocus = useCallback(() => {
    shouldMaintainFocus = false; // ยกเลิกการกู้คืน focus เมื่อ user focus เอง
  }, []);

  const handleBlur = useCallback(() => {
    shouldMaintainFocus = false; // ยกเลิกการกู้คืน focus เมื่อ user blur
  }, []);

  return (
    <div className="relative">
      <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
      <Input
        ref={inputRef}
        placeholder={placeholder}
        value={localValue}
        onChange={handleChange}
        onFocus={handleFocus}
        onBlur={handleBlur}
        className="pl-10"
      />
    </div>
  );
});