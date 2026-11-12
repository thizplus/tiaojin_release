import type { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, MoreHorizontal, Pencil, Trash } from "lucide-react";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type { City } from "@/types";
import { getPublishDraftStatusConfigForBanner } from "@/utils/status-helpers";
import { getStorageUrl } from "@/lib/utils";

export interface CityActionsHandlers {
  onEdit?: (city: City) => void;
  onDelete?: (city: City) => void;
}

export const getCityColumns = (
  handlers: CityActionsHandlers
): ColumnDef<City>[] => {
  return [
    {
      accessorKey: "image",
      header: "รูปภาพ",
      cell: ({ row }) => {
        const image = row.getValue("image") as string;
        return image ? (
          <div className="flex items-center">
            <img
              src={getStorageUrl(image)}
              alt="City"
              className="h-10 w-10 rounded-md object-cover"
            />
          </div>
        ) : (
          <span className="text-muted-foreground text-sm">ไม่มีรูปภาพ</span>
        );
      },
    },
    {
      accessorKey: "title",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          ชื่อเมือง
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => <div className="font-medium">{row.getValue("title")}</div>,
    },
    
    {
      accessorKey: "detail",
      header: "รายละเอียด",
      cell: ({ row }) => {
        const detail = row.getValue("detail") as string | undefined;
        return detail ? (
          <div className="max-w-[300px] truncate">{detail}</div>
        ) : (
          <span className="text-muted-foreground text-sm">ไม่มีรายละเอียด</span>
        );
      },
    },
    {
      accessorKey: "status",
      header: "สถานะ",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        const statusConfig = getPublishDraftStatusConfigForBanner(status);

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
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          วันที่สร้าง
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        return format(
          new Date(row.getValue("created_at")),
          "d MMM yyyy HH:mm",
          { locale: th }
        );
      },
    },
    {
      accessorKey: "updated_at",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          วันที่อัปเดต
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        return format(
          new Date(row.getValue("updated_at")),
          "d MMM yyyy HH:mm",
          { locale: th }
        );
      },
    },
    {
      id: "actions",
      header: "จัดการ",
      cell: ({ row }) => {
        const city = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="h-8 w-8 p-0"
                onClick={(e) => e.stopPropagation()}
              >
                <span className="sr-only">เปิดเมนู</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => handlers.onEdit?.(city)}>
                <Pencil className="mr-1 h-3.5 w-3.5" />
                แก้ไข
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => handlers.onDelete?.(city)}
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
};