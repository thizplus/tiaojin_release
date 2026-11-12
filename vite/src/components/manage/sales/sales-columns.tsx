import type { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, Edit, UserCheck, UserX, Trash } from "lucide-react";
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
import type { Sale } from "@/types";
import { getActiveInactiveStatusConfig } from "@/utils/status-helpers";

interface ActionsProps {
  onEdit: (sale: Sale) => void;
  onDelete: (sale: Sale) => void;
  onActivate: (sale: Sale) => void;
  onDeactivate: (sale: Sale) => void;
}

export const getSalesColumns = (actions: ActionsProps): ColumnDef<Sale>[] => [
  {
    accessorKey: "code",
    header: "รหัส",
    cell: ({ row }) => (
      <div className="font-mono text-sm">{row.getValue("code")}</div>
    ),
  },
  {
    accessorKey: "first_name",
    header: "ชื่อ",
    cell: ({ row }) => (
      <div className="font-medium">
        {row.original.first_name} {row.original.last_name}
      </div>
    ),
  },
  {
    accessorKey: "nick_name",
    header: "ชื่อเล่น",
    cell: ({ row }) => (
      <div className="text-muted-foreground">{row.getValue("nick_name")}</div>
    ),
  },
  {
    accessorKey: "user_email",
    header: "อีเมล",
    cell: ({ row }) => (
      <div className="text-sm">{row.getValue("user_email") || "-"}</div>
    ),
  },
  {
    accessorKey: "active",
    header: "สถานะ",
    cell: ({ row }) => {
      const active = row.getValue("active") as boolean;
      const statusConfig = getActiveInactiveStatusConfig(active);

      return (
        <Badge variant="outline" className="px-1.5">
          {statusConfig.icon}
          <span className={`ml-1 ${statusConfig.textColor}`}>{statusConfig.label}</span>
        </Badge>
      );
    },
  },
  {
    accessorKey: "created_at",
    header: "วันที่สร้าง",
    cell: ({ row }) => (
      <div className="text-sm text-muted-foreground">
        {format(new Date(row.getValue("created_at")), "dd MMM yyyy", { locale: th })}
      </div>
    ),
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const sale = row.original;
      const { onEdit, onDelete, onActivate, onDeactivate } = actions;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">เปิดเมนู</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onEdit(sale)}>
              <Edit className="mr-2 h-4 w-4" />
              แก้ไข
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            {sale.active ? (
              <DropdownMenuItem onClick={() => onDeactivate(sale)}>
                <UserX className="mr-2 h-4 w-4" />
                ปิดใช้งาน
              </DropdownMenuItem>
            ) : (
              <DropdownMenuItem onClick={() => onActivate(sale)}>
                <UserCheck className="mr-2 h-4 w-4" />
                เปิดใช้งาน
              </DropdownMenuItem>
            )}

            <DropdownMenuSeparator />

            <DropdownMenuItem
              onClick={() => onDelete(sale)}
              className="text-destructive focus:text-destructive"
            >
              <Trash className="mr-1 h-3.5 w-3.5 text-destructive" />
              ลบ
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];