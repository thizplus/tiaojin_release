import type { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, Eye, Edit, Trash2, Calendar, MapPin, Plane, Users, Lock, Unlock, Flame, FileText } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import type { TourPeriod } from "@/types";
import { getPublishDraftStatusConfig } from "@/utils/status-helpers";
import { AllCommissionDisplay, AllRatesDisplay } from "./utils/period-display-components";

interface ActionsProps {
  onEdit: (period: TourPeriod) => void;
  onView?: (period: TourPeriod) => void;
  onDelete?: (period: TourPeriod) => void;
  onManageSeatHold?: (period: TourPeriod) => void;
  onManageFiles?: (period: TourPeriod) => void;
}


export const getTourPeriodsColumns = (actions: ActionsProps): ColumnDef<TourPeriod>[] => [
  {
    accessorKey: "start_date",
    header: "ช่วงเวลา",
    cell: ({ row }) => {
      const period = row.original;
      const startDate = new Date(period.start_date);
      const endDate = new Date(period.end_date);
      const diffTime = Math.abs(endDate.getTime() - startDate.getTime());
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;

      return (
        <div className="space-y-1">
          <div className="flex items-center space-x-2">
            <Calendar className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium">
              {format(startDate, "dd MMM", { locale: th })} - {format(endDate, "dd MMM yyyy", { locale: th })}
            </span>
          </div>
          <div className="text-xs text-muted-foreground">
            {diffDays} วัน
          </div>
          {/* Badge สำหรับ Expiration Override - ทัวร์ไฟไหม้ */}
          {period.expiration_override && (
            <Badge variant="destructive" className="text-xs h-5">
              <Flame className="h-3 w-3" />
              {period.expiration_display || 'ด่วน'}
            </Badge>
          )}
        </div>
      );
    },
  },
  {
    accessorKey: "seat",
    header: "ที่นั่ง",
    cell: ({ row }) => {
      const period = row.original;
      return (
        <div className="space-y-1">
          <div className="flex items-center space-x-2">
            <Users className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium">{period.seat}</span>
          </div>
          <div className="text-xs text-muted-foreground space-y-1">
            <div>จอง: {period.booking}</div>
            <div>จำนอง: {period.held_seats}</div>
            <div>ว่าง: {period.available_seats}</div>
            {period.has_waiting_list && (
              <div className="mt-1">
                <span className="inline-flex items-center rounded-full bg-orange-100 px-2 py-0.5 text-xs font-medium text-orange-800 dark:bg-orange-900/50 dark:text-orange-300">
                  WL
                </span>
              </div>
            )}
          </div>
        </div>
      );
    },
  },
  {
    accessorKey: "airline",
    header: "สายการบิน",
    cell: ({ row }) => {
      const period = row.original;
      return (
        <div className="space-y-1">
          <div className="flex items-center space-x-2">
            <Plane className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium">{period.airline}</span>
          </div>
          <div className="flex items-center space-x-2">
            <MapPin className="h-3 w-3 text-muted-foreground" />
            <span className="text-xs text-muted-foreground">{period.airport}</span>
          </div>
        </div>
      );
    },
  },
  {
    accessorKey: "rate_items",
    header: "ราคา",
    cell: ({ row }) => {
      const period = row.original;
      return <AllRatesDisplay period={period} />;
    },
  },
  {
    accessorKey: "deposit",
    header: "มัดจำ",
    cell: ({ row }) => {
      const period = row.original;
      const formatDeposit = () => {
        if (period.deposit_type === 'percent') {
          return `${period.deposit}%`;
        }
        return `฿${period.deposit.toLocaleString()}`;
      };

      const getDepositLabel = () => {
        switch (period.deposit_type) {
          case 'per_unit': return 'ต่อคน';
          case 'percent': return 'เปอร์เซ็นต์';
          default: return period.deposit_type;
        }
      };

      return (
        <div className="text-sm">
          <div className="font-medium">{formatDeposit()}</div>
          <div className="text-xs text-muted-foreground">
            {getDepositLabel()}
          </div>
        </div>
      );
    },
  },
  {
    accessorKey: "commission_items",
    header: "คอมมิชชั่น",
    cell: ({ row }) => {
      const period = row.original;
      return <AllCommissionDisplay period={period} />;
    },
  },
  {
    accessorKey: "status",
    header: "สถานะ",
    cell: ({ row }) => {
      const status = row.getValue("status") as string;
      const statusConfig = getPublishDraftStatusConfig(status);

      return (
        <Badge variant="outline" className="px-1.5">
          {statusConfig.icon}
          <span className={`ml-1 ${statusConfig.textColor}`}>{statusConfig.label}</span>
        </Badge>
      );
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const period = row.original;
      const { onEdit, onView, onDelete, onManageSeatHold, onManageFiles } = actions;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">เปิดเมนู</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {onView && (
              <DropdownMenuItem onClick={() => onView(period)}>
                <Eye className="mr-2 h-4 w-4" />
                ดูรายละเอียดเต็ม
              </DropdownMenuItem>
            )}
            <DropdownMenuItem onClick={() => onEdit(period)}>
              <Edit className="mr-2 h-4 w-4" />
              แก้ไข
            </DropdownMenuItem>
            {onManageSeatHold && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => onManageSeatHold(period)}>
                  {period.held_seats > 0 ? (
                    <Unlock className="mr-2 h-4 w-4" />
                  ) : (
                    <Lock className="mr-2 h-4 w-4" />
                  )}
                  จัดการที่นั่งจำนอง
                  {period.held_seats > 0 && (
                    <span className="ml-2 text-xs bg-muted text-muted-foreground px-1.5 py-0.5 rounded">
                      {period.held_seats}
                    </span>
                  )}
                </DropdownMenuItem>
              </>
            )}
            {onManageFiles && (
              <DropdownMenuItem onClick={() => onManageFiles(period)}>
                <FileText className="mr-2 h-4 w-4" />
                จัดการใบนัดหมาย
                {period.files?.length > 0 && (
                  <span className="ml-2 text-xs bg-muted text-muted-foreground px-1.5 py-0.5 rounded">
                    {period.files.length}
                  </span>
                )}
              </DropdownMenuItem>
            )}
            {onDelete && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={() => onDelete(period)}
                  className="text-destructive focus:text-destructive"
                >
                  <Trash2 className="mr-2 h-4 w-4 text-destructive" />
                  ลบ
                </DropdownMenuItem>
              </>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];