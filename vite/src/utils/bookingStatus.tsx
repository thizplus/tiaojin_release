import { Clock, AlertCircle, CheckCircle, XCircle } from 'lucide-react';

export const getBookingStatusConfig = (status: string) => {
  switch (status) {
    case "pending":
      return {
        label: "รอดำเนินการ",
        variant: "outline" as const,
        icon: <Clock className="h-4 w-4 text-chart-4" />,
        iconLarge: <Clock className="h-10 w-10 text-chart-4" />,
        textColor: "text-chart-4",
        bgColor: "bg-chart-4/10"
      };
    case "expired":
      return {
        label: "หมดอายุ",
        variant: "outline" as const,
        icon: <Clock className="h-4 w-4 text-muted-foreground" />,
        iconLarge: <Clock className="h-10 w-10 text-muted-foreground" />,
        textColor: "text-muted-foreground",
        bgColor: "bg-muted/10"
      };
    case "processing":
      return {
        label: "กำลังดำเนินการ",
        variant: "outline" as const,
        icon: <AlertCircle className="h-4 w-4 text-chart-1" />,
        iconLarge: <AlertCircle className="h-10 w-10 text-chart-1" />,
        textColor: "text-chart-1",
        bgColor: "bg-chart-1/10"
      };
    case "confirm":
      return {
        label: "ยืนยันแล้ว",
        variant: "outline" as const,
        icon: <CheckCircle className="h-4 w-4 text-chart-2" />,
        iconLarge: <CheckCircle className="h-10 w-10 text-chart-2" />,
        textColor: "text-chart-2",
        bgColor: "bg-chart-2/10"
      };
    case "cancel":
      return {
        label: "ยกเลิก",
        variant: "outline" as const,
        icon: <XCircle className="h-4 w-4 text-destructive" />,
        iconLarge: <XCircle className="h-10 w-10 text-destructive" />,
        textColor: "text-destructive",
        bgColor: "bg-destructive/10"
      };
    case "completed":
      return {
        label: "เสร็จสิ้น",
        variant: "outline" as const,
        icon: <CheckCircle className="h-4 w-4 text-chart-2" />,
        iconLarge: <CheckCircle className="h-10 w-10 text-chart-2" />,
        textColor: "text-chart-2",
        bgColor: "bg-chart-2/10"
      };
    default:
      return {
        label: status,
        variant: "outline" as const,
        icon: <Clock className="h-4 w-4 text-muted-foreground" />,
        iconLarge: <Clock className="h-10 w-10 text-muted-foreground" />,
        textColor: "text-muted-foreground",
        bgColor: "bg-muted/10"
      };
  }
};

// สำหรับ payment slip และ passport status
export const getApprovalStatusConfig = (status: string) => {
  switch (status) {
    case 'approved':
      return {
        label: 'อนุมัติแล้ว',
        icon: <CheckCircle className="h-4 w-4 text-chart-2" />,
        textColor: 'text-chart-2',
        bgColor: 'bg-chart-2/10'
      };
    case 'pending':
      return {
        label: 'รออนุมัติ',
        icon: <Clock className="h-4 w-4 text-chart-4" />,
        textColor: 'text-chart-4',
        bgColor: 'bg-chart-4/10'
      };
    case 'rejected':
      return {
        label: 'ไม่อนุมัติ',
        icon: <XCircle className="h-4 w-4 text-destructive" />,
        textColor: 'text-destructive',
        bgColor: 'bg-destructive/10'
      };
   
      

// slip status
case 'confirm':
  return {
    label: 'ยืนยัน',
    icon: <CheckCircle className="h-4 w-4 text-chart-2" />,
    textColor: 'text-chart-2',
    bgColor: 'bg-chart-2/10'
  };

  case 'reject':
    return {
      label: 'ปฏิเสธ',
      icon: <XCircle className="h-4 w-4 text-destructive" />,
      textColor: 'text-destructive',
      bgColor: 'bg-destructive/10'
    };
    case 'reversed':
    return {
      label: 'ยกเลิก',
      icon: <XCircle className="h-4 w-4 text-destructive" />,
      textColor: 'text-destructive',
      bgColor: 'bg-destructive/10'
    };

 
    default:
      return {
        label: 'ไม่ทราบสถานะ',
        icon: <Clock className="h-4 w-4 text-muted-foreground" />,
        textColor: 'text-muted-foreground',
        bgColor: 'bg-muted/10'
      };
  }
};