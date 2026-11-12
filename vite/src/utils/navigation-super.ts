import {
  Home,
  Settings,

  ClipboardList,
  Plane,
  MapPin,
  Users,
  UserCog,
  UserCheck,
  Image,
  Bell,

  BarChart3,
  Calendar,
  type LucideIcon,
  Flame,
} from "lucide-react";

export interface NavItem {
  title: string;
  url: string;
  icon?: LucideIcon;
  isActive?: boolean;
}

export const superNavigation = {
  dashboard: [
    {
      title: "แดชบอร์ด",
      url: "/admin/dashboard",
      icon: Home,
    }
  ],
  systemSettings: [
    {
      title: "ตั้งค่าระบบ",
      url: "/super/system-settings",
      icon: Settings,
    }
  ],
  mainMenu: [
    {
      title: "ทัวร์ทั้งหมด",
      url: "/admin/tours-regular",
      icon: Plane,
    },
    
    {
      title: "ทัวร์ไฟไหม้",
      url: "/admin/tours-special",
      icon: Flame,
    },
    {
      title: "รายการจองทัวร์ทั้งหมด",
      url: "/admin/bookings",
      icon: ClipboardList,
    }
  ],
  dataManagement: [
    {
      title: "โปรแกรมทัวร์",
      url: "/admin/tours",
      icon: Plane,
    },
    {
      title: "ข้อมูลเมือง",
      url: "/admin/cities",
      icon: MapPin,
    },
    {
      title: "ข้อมูลผู้ใช้งาน",
      url: "/admin/users",
      icon: Users,
    },
    {
      title: "ข้อมูลโปรไฟล์",
      url: "/admin/profiles",
      icon: UserCheck,
    },
    {
      title: "ข้อมูลพนักงานขาย",
      url: "/admin/sales",
      icon: UserCog,
    },
  
    {
      title: "ข้อมูลแบนเนอร์",
      url: "/admin/banners",
      icon: Image,
    },
    {
      title: "ข้อมูลการแจ้งเตือน",
      url: "/admin/notifications",
      icon: Bell,
    }
  ],

  reports: [
    {
      title: "รายงานตามทัวร์",
      url: "/admin/reports-tours",
      icon: BarChart3,
    },
    {
      title: "รายงานเลือกช่วงวันที่",
      url: "/admin/reports-daterange",
      icon: Calendar,
    }
  ]
};