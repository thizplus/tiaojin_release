import * as React from "react";
import { Check, ChevronsUpDown, Search } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useCity } from "@/hooks/useCity";
import { Badge } from "@/components/ui/badge";


interface CitySelectProps {
  value?: string;
  onValueChange: (value: string | undefined) => void;
  placeholder?: string;
  disabled?: boolean;
}

export function CitySelect({
  value,
  onValueChange,
  placeholder = "เลือกจังหวัด...",
  disabled = false,
}: CitySelectProps) {
  const [open, setOpen] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState("");

  const { citiesList, getCitiesList, isCitiesListLoading } = useCity(
    { search: searchQuery, limit: 50, status: 'publish' },
    undefined,
    { enableList: true, enableById: false }
  );

  // Fetch cities when search query changes
  React.useEffect(() => {
    if (searchQuery.trim()) {
      getCitiesList();
    }
  }, [searchQuery]);

  // Fetch initial cities when component mounts
  React.useEffect(() => {
    getCitiesList();
  }, []);

  const selectedCity = citiesList.find(city => city.id === value);

  const handleSelect = (cityId: string) => {
    if (cityId === value) {
      onValueChange(undefined); // Clear selection
    } else {
      onValueChange(cityId);
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
          {selectedCity ? (
            <div className="flex items-center gap-2">
              <span>{selectedCity.title}</span>
              <Badge variant="secondary" className="text-xs">
                {selectedCity.status}
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
          <div className="flex items-center border-b px-3">
            <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
            <CommandInput
              placeholder="ค้นหาชื่อจังหวัด..."
              value={searchQuery}
              onValueChange={setSearchQuery}
              className="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
          <CommandList>
            {isCitiesListLoading ? (
              <div className="p-4 text-center text-sm text-muted-foreground">
                กำลังค้นหา...
              </div>
            ) : (
              <>
                <CommandEmpty>
                  {searchQuery ? "ไม่พบจังหวัดที่ค้นหา" : "พิมพ์เพื่อค้นหาจังหวัด"}
                </CommandEmpty>
                <CommandGroup>
                  {citiesList.map((city) => (
                    <CommandItem
                      key={city.id}
                      value={city.title}
                      onSelect={() => handleSelect(city.id)}
                      className="flex items-center justify-between"
                    >
                      <div className="flex flex-col">
                        <div className="flex items-center gap-2">
                          <span>{city.title}</span>
                          <Badge variant="outline" className="text-xs">
                            {city.status}
                          </Badge>
                        </div>
                        {city.detail && (
                          <span className="text-xs text-muted-foreground truncate max-w-[300px]">
                            {city.detail}
                          </span>
                        )}
                      </div>
                      <Check
                        className={cn(
                          "ml-2 h-4 w-4",
                          value === city.id ? "opacity-100" : "opacity-0"
                        )}
                      />
                    </CommandItem>
                  ))}
                </CommandGroup>
              </>
            )}
          </CommandList>
        </Command>
        {selectedCity && (
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