import {
  Home,
  Plane,
  Flame,
  ClipboardList,
  UserCheck,
  Bell,
  BarChart3,
  type LucideIcon,
} from "lucide-react";

export interface NavItem {
  title: string;
  url: string;
  icon?: LucideIcon;
  isActive?: boolean;
}

export const agentNavigation = {
  dashboard: [
    {
      title: "แดชบอร์ด",
      url: "/agent/dashboard",
      icon: Home,
    }
  ],
  tourData: [
    {
      title: "ทัวร์ทั้งหมด",
      url: "/agent/tours",
      icon: Plane,
    },
    {
      title: "ทัวร์ไฟไหม้",
      url: "/agent/tours-hot",
      icon: Flame,
    },
    {
      title: "รายการจองทัวร์",
      url: "/agent/bookings",
      icon: ClipboardList,
    }
  ],
  mainMenu: [
    {
      title: "โปรไฟล์",
      url: "/agent/profile",
      icon: UserCheck,
    },
    {
      title: "การแจ้งเตือน",
      url: "/agent/notifications",
      icon: Bell,
    },
    {
      title: "รายงาน",
      url: "/agent/reports",
      icon: BarChart3,
    }
  ]
};