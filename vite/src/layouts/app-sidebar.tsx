"use client"

import * as React from "react"

import { NavRoleSuper } from "@/layouts/nav/super/nav-role-super"
import { NavRoleAdmin } from "@/layouts/nav/admin/nav-role-admin"
import { NavRoleAgent } from "@/layouts/nav/agent/nav-role-agent"
import { NavRoleSale } from "@/layouts/nav/sales/nav-role-sale"
import { NavUser } from "@/layouts/nav-user"
import { useAuth } from "@/hooks/useAuth"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth()

  // Don't render sidebar if no user
  if (!user) {
    return null
  }

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="/dashboard">
                <img src="/images/id/logo.png" alt="TIAOJIN" className="size-8" />
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-semibold">TIAOJIN</span>
                  <span className="text-xs text-muted-foreground">ทัวร์จีนต้องเที่ยวจีน</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        {user?.role === "super" && <NavRoleSuper />}
        {user?.role === "admin" && <NavRoleAdmin />}
        {user?.role === "agent" && <NavRoleAgent />}
        {user?.role === "sales" && <NavRoleSale />}
      </SidebarContent>

      <SidebarFooter>
        <NavUser />
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>
  )
}