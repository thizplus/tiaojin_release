"use client"

import { NavAdminDashboard } from "@/layouts/nav/admin/nav-admin-dashboard"
import { NavAdminMain } from "@/layouts/nav/admin/nav-admin-main"
import { NavAdminData } from "@/layouts/nav/admin/nav-admin-data"

import { NavAdminReports } from "@/layouts/nav/admin/nav-admin-reports"

export function NavRoleAdmin() {
  return (
    <>
      <NavAdminDashboard />
      <NavAdminMain />
      <NavAdminData />
      <NavAdminReports />
    </>
  )
}