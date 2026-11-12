"use client"

import { useNavigate } from "react-router-dom"
import {
  BadgeCheck,

  ChevronsUpDown,
  LogOut,
  Sparkles,
} from "lucide-react"

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { useAuth } from "@/hooks/useAuth"
import { toast } from "sonner"

export function NavUser() {
  const { isMobile } = useSidebar()
  const navigate = useNavigate()
  const { user, logout, isLogoutLoading } = useAuth()

  // If no user, don't render
  if (!user) {
    return null
  }

  const handleLogout = async () => {
    try {
      await logout()
      toast.success("ออกจากระบบเรียบร้อย")
      navigate('/login')
    } catch (error) {
      toast.error("เกิดข้อผิดพลาดในการออกจากระบบ", {
        description: error instanceof Error ? error.message : "เกิดข้อผิดพลาดในการออกจากระบบ"
      })
    }
  }

  const handleAccountClick = () => {
    navigate('/profile')
  }



  const handleApiClick = () => {
    navigate('/agent/api-keys')
  }

  // Generate user initials for avatar fallback
  const getInitials = (name?: string | null) => {
    if (!name) return 'U'
    return name.substring(0, 2).toUpperCase()
  }

  // Get display name - check status first, then use email or first_name + last_name
  const getDisplayName = () => {
    if (user.status === 'pending') {
      return getRoleDisplay()
    }
    if (user.first_name && user.last_name) {
      return `${user.first_name} ${user.last_name}`
    }
    return user.email || user.username || 'ผู้ใช้'
  }

  // Get email or username for display
  const getEmailDisplay = () => {
    return user.email || user.username || '-'
  }

  // Get role display text
  const getRoleDisplay = () => {
    switch (user.role) {
      case 'super':
        return 'ผู้ดูแลระบบสูงสุด'
      case 'admin':
        return 'ผู้ดูแลระบบ'
      case 'agent':
        return 'เอเจนต์'
      case 'sales':
        return 'เซลส์'
      default:
        return user.role
    }
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage src="" alt={getDisplayName()} />
                <AvatarFallback className="rounded-lg">
                  {getInitials(user.email || user.username)}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{getDisplayName()}</span>
                <span className="truncate text-xs">{getEmailDisplay()}</span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarImage src="" alt={getDisplayName()} />
                  <AvatarFallback className="rounded-lg">
                    {getInitials(user.email || user.username)}
                  </AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate text-xs">{getEmailDisplay()}</span>
                  <span className="truncate text-xs text-muted-foreground">
                    {getRoleDisplay()}
                  </span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {user.role === 'agent' && (
              <>
                <DropdownMenuGroup>
                  <DropdownMenuItem onClick={handleApiClick}>
                    <Sparkles />
                    API KEYS
                  </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
              </>
            )}
            <DropdownMenuGroup>
              <DropdownMenuItem onClick={handleAccountClick}>
                <BadgeCheck />
                บัญชี
              </DropdownMenuItem>

             
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={handleLogout}
              disabled={isLogoutLoading}
            >
              <LogOut />
              {isLogoutLoading ? "กำลังออกจากระบบ..." : "ออกจากระบบ"}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}