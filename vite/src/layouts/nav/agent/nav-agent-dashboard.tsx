"use client"

import { Link, useLocation } from "react-router-dom"
import { agentNavigation } from "@/utils/navigation-agent"
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export function NavAgentDashboard() {
  const location = useLocation()

  return (
    <SidebarGroup>
      <SidebarMenu>
        {agentNavigation.dashboard.map((item) => {
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