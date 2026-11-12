"use client"

import { Link, useLocation } from "react-router-dom"
import { adminNavigation } from "@/utils/navigation-admin"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavAdminData() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>จัดการข้อมูล</SidebarGroupLabel>
      <SidebarMenu>
        {adminNavigation.dataManagement.map((item) => {
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