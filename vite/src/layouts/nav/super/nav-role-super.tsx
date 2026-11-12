"use client"

import { NavSuperDashboard } from "@/layouts/nav/super/nav-super-dashboard"
import { NavSuperSystem } from "@/layouts/nav/super/nav-super-system"
import { NavSuperMain } from "@/layouts/nav/super/nav-super-main"
import { NavSuperData } from "@/layouts/nav/super/nav-super-data"

import { NavSuperReports } from "@/layouts/nav/super/nav-super-reports"

export function NavRoleSuper() {
  return (
    <>
      <NavSuperDashboard />
      <NavSuperSystem />
      <NavSuperMain />
      <NavSuperData />
    
      <NavSuperReports />
    </>
  )
}