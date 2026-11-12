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

export function NavAdminMain() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>เมนูหลัก</SidebarGroupLabel>
      <SidebarMenu>
        {adminNavigation.mainMenu.map((item) => {
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