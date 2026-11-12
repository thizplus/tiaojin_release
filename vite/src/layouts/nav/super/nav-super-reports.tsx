"use client"

import { Link, useLocation } from "react-router-dom"
import { superNavigation } from "@/utils/navigation-super"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavSuperReports() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>รายงาน</SidebarGroupLabel>
      <SidebarMenu>
        {superNavigation.reports.map((item) => {
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