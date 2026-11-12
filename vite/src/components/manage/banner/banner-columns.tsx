import type { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, MoreHorizontal, Pencil, Trash, ExternalLink } from "lucide-react";
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
import type { Banner } from "@/types";
import { getPublishDraftStatusConfigForBanner } from "@/utils/status-helpers";
import { getStorageUrl } from "@/lib/utils";

export interface BannerActionsHandlers {
  onEdit?: (banner: Banner) => void;
  onDelete?: (banner: Banner) => void;
}

export const getBannerColumns = (
  handlers: BannerActionsHandlers
): ColumnDef<Banner>[] => {
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
              alt="Banner"
              className="h-10 w-16 rounded-md object-cover"
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
          ชื่อแบนเนอร์
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => <div className="font-medium">{row.getValue("title")}</div>,
    },
    {
      accessorKey: "link_url",
      header: "ลิงก์",
      cell: ({ row }) => {
        const link_url = row.getValue("link_url") as string | null;
        return link_url ? (
          <div className="flex items-center gap-1">
            <a
              href={link_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:text-blue-800 underline max-w-[200px] truncate"
            >
              {link_url}
            </a>
            <ExternalLink className="h-3 w-3 text-muted-foreground" />
          </div>
        ) : (
          <span className="text-muted-foreground text-sm">ไม่มีลิงก์</span>
        );
      },
    },
    {
      accessorKey: "alt",
      header: "Alt Text",
      cell: ({ row }) => {
        const alt = row.getValue("alt") as string | null;
        return alt ? (
          <div className="max-w-[200px] truncate">{alt}</div>
        ) : (
          <span className="text-muted-foreground text-sm">ไม่มี Alt Text</span>
        );
      },
    },
    {
      accessorKey: "display_order",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          ลำดับ
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const order = row.getValue("display_order") as number;
        return (
          <div className="text-center font-medium">
            {order}
          </div>
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
        const banner = row.original;
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
              <DropdownMenuItem onClick={() => handlers.onEdit?.(banner)}>
                <Pencil className="mr-1 h-3.5 w-3.5" />
                แก้ไข
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => handlers.onDelete?.(banner)}
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