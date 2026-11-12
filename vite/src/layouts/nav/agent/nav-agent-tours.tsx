"use client"

import { Link, useLocation } from "react-router-dom"
import { agentNavigation } from "@/utils/navigation-agent"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavAgentTours() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>ข้อมูลทัวร์</SidebarGroupLabel>
      <SidebarMenu>
        {agentNavigation.tourData.map((item) => {
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