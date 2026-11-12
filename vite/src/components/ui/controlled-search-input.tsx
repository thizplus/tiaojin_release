import * as React from "react";
import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

interface ControlledSearchInputProps {
  onSearchChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  defaultValue?: string;
  value?: string; // Controlled mode - sync with URL
}

export const ControlledSearchInput = React.memo(
  React.forwardRef<HTMLInputElement, ControlledSearchInputProps>(
    ({ onSearchChange, placeholder = "ค้นหา...", className, defaultValue = "", value: controlledValue, ...props }, ref) => {
      // ใช้ controlled value ถ้ามี, ไม่งั้นใช้ internal state
      const isControlled = controlledValue !== undefined;
      const [internalValue, setInternalValue] = React.useState(defaultValue);
      const value = isControlled ? controlledValue : internalValue;

      const inputRef = React.useRef<HTMLInputElement>(null);
      const timeoutRef = React.useRef<NodeJS.Timeout>(null);
      const prevDefaultValueRef = React.useRef(defaultValue);

      // ใช้ ref ที่ส่งมาจากภายนอก หรือใช้ internal ref
      const finalRef = (ref as React.RefObject<HTMLInputElement>) || inputRef;

      // Sync internal value when defaultValue changes (for initial load from URL)
      // ใช้ ref เพื่อ track ค่าเดิม ป้องกัน reset เมื่อ user พิมพ์
      React.useEffect(() => {
        if (!isControlled && defaultValue !== prevDefaultValueRef.current) {
          setInternalValue(defaultValue);
          prevDefaultValueRef.current = defaultValue;
        }
      }, [defaultValue, isControlled]);

      const handleChange = React.useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        const newValue = e.target.value;

        if (!isControlled) {
          setInternalValue(newValue);
        }

        // Clear previous timeout
        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }

        // Set new timeout for debounce
        timeoutRef.current = setTimeout(() => {
          onSearchChange(newValue);
        }, 500);
      }, [onSearchChange, isControlled]);

      // Cleanup timeout on unmount
      React.useEffect(() => {
        return () => {
          if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
          }
        };
      }, []);

      return (
        <div className={cn("relative w-full max-w-sm", className)}>
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            ref={finalRef}
            placeholder={placeholder}
            value={value}
            onChange={handleChange}
            className="w-full pl-9"
            {...props}
          />
        </div>
      );
    }
  )
);

ControlledSearchInput.displayName = "ControlledSearchInput";