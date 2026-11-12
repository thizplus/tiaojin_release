"use client"

import { NavAgentDashboard } from "@/layouts/nav/agent/nav-agent-dashboard"
import { NavAgentTours } from "@/layouts/nav/agent/nav-agent-tours"
import { NavAgentMain } from "@/layouts/nav/agent/nav-agent-main"

export function NavRoleAgent() {
  return (
    <>
      <NavAgentDashboard />
      <NavAgentTours />
      <NavAgentMain />
    </>
  )
}