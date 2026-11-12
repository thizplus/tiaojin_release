import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useSales } from "@/hooks/useSales";
import { Badge } from "@/components/ui/badge";

interface SalesSelectProps {
  value?: string;
  onValueChange: (value: string | undefined) => void;
  placeholder?: string;
  disabled?: boolean;
}

export function SalesSelect({
  value,
  onValueChange,
  placeholder = "เลือกเซลส์...",
  disabled = false,
}: SalesSelectProps) {
  const [open, setOpen] = React.useState(false);

  // ดึงรายการ sales โดยใช้ public endpoint สำหรับ active sales
  const { salesList, isSalesActiveLoading } = useSales(
    undefined,
    undefined,
    { enableList: false, enableById: false, enableActive: true }
  );

  const selectedSales = salesList.find(sales => sales.id === value);

  const handleSelect = (salesId: string) => {
    if (salesId === value) {
      onValueChange(undefined); // Clear selection
    } else {
      onValueChange(salesId);
    }
    setOpen(false);
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onValueChange(undefined);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
          disabled={disabled}
        >
          {selectedSales ? (
            <div className="flex items-center gap-2">
              <span>{selectedSales.code} - {selectedSales.first_name} {selectedSales.last_name}</span>
              <Badge variant="outline" className="text-xs">
                {selectedSales.nick_name}
              </Badge>
            </div>
          ) : (
            <span className="text-muted-foreground">{placeholder}</span>
          )}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[400px] p-0">
        <Command>
          <CommandList>
            {isSalesActiveLoading ? (
              <div className="p-4 text-center text-sm text-muted-foreground">
                กำลังโหลด...
              </div>
            ) : (
              <>
                <CommandEmpty>
                  ไม่พบข้อมูลเซลส์
                </CommandEmpty>
                <CommandGroup>
                  {salesList.map((sales) => (
                    <CommandItem
                      key={sales.id}
                      value={`${sales.code}-${sales.first_name}-${sales.last_name}`}
                      onSelect={() => handleSelect(sales.id)}
                      className="flex items-center justify-between"
                    >
                      <div className="flex flex-col">
                        <div className="flex items-center gap-2">
                          <span>{sales.code} - {sales.first_name} {sales.last_name}</span>
                          <Badge variant="outline" className="text-xs">
                            {sales.nick_name}
                          </Badge>
                        </div>
                        {sales.user_email && (
                          <span className="text-xs text-muted-foreground">
                            {sales.user_email}
                          </span>
                        )}
                      </div>
                      <Check
                        className={cn(
                          "ml-2 h-4 w-4",
                          value === sales.id ? "opacity-100" : "opacity-0"
                        )}
                      />
                    </CommandItem>
                  ))}
                </CommandGroup>
              </>
            )}
          </CommandList>
        </Command>
        {selectedSales && (
          <div className="border-t p-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleClear}
              className="w-full text-xs"
            >
              ล้างการเลือก
            </Button>
          </div>
        )}
      </PopoverContent>
    </Popover>
  );
}