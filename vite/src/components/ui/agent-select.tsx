import * as React from "react";
import { Check, ChevronsUpDown, Search, User } from "lucide-react";
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
import type { User as UserType } from "@/types";

interface AgentSelectProps {
  value?: string;
  onValueChange: (value: string | undefined) => void;
  placeholder?: string;
  disabled?: boolean;
  showUserId?: boolean;
}

export function AgentSelect({
  value,
  onValueChange,
  placeholder = "เลือก Agent...",
  disabled = false,
  showUserId = true,
}: AgentSelectProps) {
  const [open, setOpen] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState("");

  // Filter for agents only (role: 'agent' or 'super')
  const { usersList, getUsersList, isUsersListLoading } = useUser(
    {
      search: searchQuery,
      limit: 50,
      role: searchQuery.trim() ? undefined : 'agent' // Show agents by default, all when searching
    },
    undefined,
    { enableList: true, enableById: false, enableStats: false }
  );

  // Fetch users when search query changes
  React.useEffect(() => {
    if (searchQuery.trim()) {
      getUsersList({ search: searchQuery, limit: 50 });
    } else {
      getUsersList({ role: 'agent', limit: 50 });
    }
  }, [searchQuery, getUsersList]);

  // Filter to only show agent and super roles
  const filteredUsers = React.useMemo(() => {
    return usersList.filter(user =>
      user.role === 'agent' || user.role === 'super'
    );
  }, [usersList]);

  const selectedUser = filteredUsers.find(user => user.id === value);

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

  const getUserDisplayText = (user: UserType) => {
    if (showUserId) {
      return `${user.email} (${user.id})`;
    }
    return user.email;
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between font-mono text-sm"
          disabled={disabled}
        >
          {selectedUser ? (
            <div className="flex items-center gap-2 truncate">
              <User className="h-3 w-3 text-muted-foreground" />
              <span className="truncate">{getUserDisplayText(selectedUser)}</span>
            </div>
          ) : (
            <span className="text-muted-foreground">{placeholder}</span>
          )}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[500px] p-0">
        <Command>
          <div className="flex items-center border-b px-3">
            <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
            <CommandInput
              placeholder="ค้นหาอีเมลหรือ Agent ID..."
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
                  {searchQuery ? "ไม่พบ Agent ที่ค้นหา" : "ไม่มี Agent ในระบบ"}
                </CommandEmpty>
                <CommandGroup>
                  {filteredUsers.map((user) => (
                    <CommandItem
                      key={user.id}
                      value={user.email}
                      onSelect={() => handleSelect(user.id)}
                      className="flex items-center justify-between"
                    >
                      <div className="flex flex-col min-w-0 flex-1">
                        <div className="flex items-center gap-2">
                          <User className="h-3 w-3 text-muted-foreground" />
                          <span className="truncate font-mono text-sm">{user.email}</span>
                          <Badge variant={user.role === 'super' ? 'default' : 'secondary'} className="text-xs">
                            {user.role}
                          </Badge>
                        </div>
                        {showUserId && (
                          <div className="flex items-center gap-2 mt-1">
                            <span className="text-xs text-muted-foreground font-mono">
                              ID: {user.id}
                            </span>
                            <Badge variant="outline" className="text-xs">
                              {user.status}
                            </Badge>
                          </div>
                        )}
                      </div>
                      <Check
                        className={cn(
                          "ml-2 h-4 w-4 shrink-0",
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