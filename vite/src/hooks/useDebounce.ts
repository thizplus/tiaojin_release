import { useState, useEffect } from 'react';

/**
 * Custom hook สำหรับ debounce value
 * @param value - ค่าที่ต้องการ debounce
 * @param delay - เวลาหน่วงในมิลลิวินาที (default: 300ms)
 * @returns debounced value
 */
export function useDebounce<T>(value: T, delay: number = 300): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    // ตั้งค่า timeout สำหรับการอัปเดต debounced value
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    // ล้าง timeout หากมีการเปลี่ยนแปลง value ก่อนหมดเวลา
    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

export default useDebounce;