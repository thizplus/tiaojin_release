"use client"

import { NavSaleDashboard } from "@/layouts/nav/sales/nav-sale-dashboard"

import { NavSaleTours } from "@/layouts/nav/sales/nav-sale-tours"

export function NavRoleSale() {
  return (
    <>
      <NavSaleDashboard />
    
      <NavSaleTours />
    </>
  )
}