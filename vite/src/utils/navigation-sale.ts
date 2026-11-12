import {
  Home,
  ClipboardList,

  type LucideIcon,

} from "lucide-react";

export interface NavItem {
  title: string;
  url: string;
  icon?: LucideIcon;
  isActive?: boolean;
}

export const saleNavigation = {
  dashboard: [
    {
      title: "แดชบอร์ด",
      url: "/sale/dashboard",
      icon: Home,
    }
  ],

  tourData: [
 
    {
      title: "รายการจองทัวร์ทั้งหมด",
      url: "/sale/bookings",
      icon: ClipboardList,
    }
  ]
};