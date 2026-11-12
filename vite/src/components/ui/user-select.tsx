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
import { useUser } from "@/hooks/useUser";
import { Badge } from "@/components/ui/badge";


interface UserSelectProps {
  value?: string;
  onValueChange: (value: string | undefined) => void;
  placeholder?: string;
  disabled?: boolean;
}

export function UserSelect({
  value,
  onValueChange,
  placeholder = "เลือกผู้ใช้...",
  disabled = false,
}: UserSelectProps) {
  const [open, setOpen] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState("");

  const { usersList, getUsersList, isUsersListLoading } = useUser(
    { search: searchQuery, limit: 50 },
    undefined,
    { enableList: true, enableById: false, enableStats: false }
  );

  // Fetch users when search query changes
  React.useEffect(() => {
    if (searchQuery.trim()) {
      getUsersList({ search: searchQuery, limit: 50 });
    }
  }, [searchQuery, getUsersList]);

  // Debug: log usersList
  React.useEffect(() => {
    console.log("UserSelect - usersList:", usersList);
    console.log("UserSelect - searchQuery:", searchQuery);
    console.log("UserSelect - isLoading:", isUsersListLoading);
  }, [usersList, searchQuery, isUsersListLoading]);

  const selectedUser = usersList.find(user => user.id === value);

  const handleSelect = (userId: string) => {
    if (userId === value) {
      onValueChange(undefined); // Clear selection
    } else {
      onValueChange(userId);
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
          {selectedUser ? (
            <div className="flex items-center gap-2">
              <span>{selectedUser.email}</span>
             
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
              placeholder="ค้นหาอีเมลผู้ใช้..."
              value={searchQuery}
              onValueChange={setSearchQuery}
              className="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
          <CommandList>
            {isUsersListLoading ? (
              <div className="p-4 text-center text-sm text-muted-foreground">
                กำลังค้นหา...
              </div>
            ) : (
              <>
                <CommandEmpty>
                  {searchQuery ? "ไม่พบผู้ใช้ที่ค้นหา" : "พิมพ์เพื่อค้นหาผู้ใช้"}
                </CommandEmpty>
                <CommandGroup>
                  {usersList.map((user) => (
                    <CommandItem
                      key={user.id}
                      value={user.email}
                      onSelect={() => handleSelect(user.id)}
                      className="flex items-center justify-between"
                    >
                      <div className="flex flex-col">
                        <div className="flex items-center gap-2">
                          <span>{user.email}</span>
                          <Badge variant="outline" className="text-xs">
                            {user.role}
                          </Badge>
                        </div>
                     
                      </div>
                      <Check
                        className={cn(
                          "ml-2 h-4 w-4",
                          value === user.id ? "opacity-100" : "opacity-0"
                        )}
                      />
                    </CommandItem>
                  ))}
                </CommandGroup>
              </>
            )}
          </CommandList>
        </Command>
        {selectedUser && (
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