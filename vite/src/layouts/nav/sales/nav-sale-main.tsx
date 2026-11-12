"use client"

import { Link, useLocation } from "react-router-dom"
import { saleNavigation, type NavItem } from "@/utils/navigation-sale"
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavSaleMain() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarMenu>
        {saleNavigation.dashboard.map((item: NavItem) => {
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