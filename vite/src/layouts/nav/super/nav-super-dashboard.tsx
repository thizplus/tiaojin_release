"use client"

import { Link, useLocation } from "react-router-dom"
import { superNavigation } from "@/utils/navigation-super"
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavSuperDashboard() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarMenu>
        {superNavigation.dashboard.map((item) => {
          const isCurrentPage = location.pathname === item.url

          return (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild tooltip={item.title} isActive={isCurrentPage}>
                <Link to={item.url}>
                  {item.icon && <item.icon />}
                  <span>{item.title}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          )
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}