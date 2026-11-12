import type { ChartConfig } from "@/components/ui/chart"

export const bookingStatusChartConfig = {
  pending: {
    label: "รอดำเนินการ",
    color: "hsl(48, 96%, 53%)", // yellow
  },
  processing: {
    label: "กำลังดำเนินการ",
    color: "hsl(217, 91%, 60%)", // blue
  },
  confirm: {
    label: "ยืนยันแล้ว",
    color: "hsl(142, 71%, 45%)", // green
  },
  cancel: {
    label: "ยกเลิก",
    color: "hsl(0, 84%, 60%)", // red
  },
  completed: {
    label: "เสร็จสิ้น",
    color: "hsl(215, 16%, 47%)", // gray
  },
  expired: {
    label: "หมดอายุ",
    color: "hsl(25, 95%, 53%)", // orange
  },
} satisfies ChartConfig