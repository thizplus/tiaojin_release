import type { ColumnDef } from "@tanstack/react-table";
import { MoreHorizontal, Edit, FileDown, FileText, CheckCircle, XCircle, Trash, Calendar } from "lucide-react";
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
import type { Tour } from "@/types";
import { getPublishDraftStatusConfig } from "@/utils/status-helpers";
import { getStorageUrl } from "@/lib/utils";

interface ActionsProps {
  onEdit: (tour: Tour) => void;
  onDelete: (tour: Tour) => void;
  onUpdateStatus: (tour: Tour) => void;
  onManagePeriods: (tour: Tour) => void;
}

export const getTourColumns = (actions: ActionsProps): ColumnDef<Tour>[] => [
  {
    accessorKey: "thumbnail",
    header: "รูปภาพ",
    cell: ({ row }) => {
      const thumbnail = row.getValue("thumbnail") as string;
      return thumbnail ? (
        <div className="flex items-center">
          <img
            src={getStorageUrl(thumbnail)}
            alt="Tour"
            className="h-10 w-10 rounded-md object-cover"
          />
        </div>
      ) : (
        <span className="text-muted-foreground text-sm">ไม่มีรูปภาพ</span>
      );
    },
  },
  {
    accessorKey: "code",
    header: "รหัสทัวร์",
    cell: ({ row }) => (
      <div className="font-mono text-sm">{row.getValue("code")}</div>
    ),
  },
  {
    accessorKey: "title",
    header: "ชื่อทัวร์",
    cell: ({ row }) => (
      <div className="font-medium max-w-[300px] truncate">
        {row.getValue("title")}
      </div>
    ),
  },
  {
    accessorKey: "city_title",
    header: "จังหวัด",
    cell: ({ row }) => (
      <div className="text-sm">{row.getValue("city_title") || "-"}</div>
    ),
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
    id: "files",
    header: "ไฟล์แนบ",
    cell: ({ row }) => {
      const tour = row.original;
      return (
        <div className="flex gap-1">
          {tour.file_pdf && (
            <Badge variant="outline" className="text-xs">
              <FileDown className="h-3 w-3 mr-1" />
              PDF
            </Badge>
          )}
          {tour.file_word && (
            <Badge variant="outline" className="text-xs">
              <FileText className="h-3 w-3 mr-1" />
              DOC
            </Badge>
          )}
        </div>
      );
    },
  },
  {
    id: "periods",
    header: "จัดการ Periods",
    cell: ({ row }) => {
      const tour = row.original;
      const { onManagePeriods } = actions;
      return (
        <Button
          variant="outline"
          size="sm"
          onClick={() => onManagePeriods(tour)}
          className="h-8"
        >
          <Calendar className="h-4 w-4 mr-2" />
          จัดการ Periods
        </Button>
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
      const tour = row.original;
      const { onEdit, onDelete, onUpdateStatus, onManagePeriods } = actions;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">เปิดเมนู</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onEdit(tour)}>
              <Edit className="mr-2 h-4 w-4" />
              แก้ไข
            </DropdownMenuItem>

            <DropdownMenuItem onClick={() => onManagePeriods(tour)}>
              <Calendar className="mr-2 h-4 w-4" />
              จัดการ Periods
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            {tour.status === "draft" ? (
              <DropdownMenuItem onClick={() => onUpdateStatus(tour)}>
                <CheckCircle className="mr-2 h-4 w-4" />
                เผยแพร่
              </DropdownMenuItem>
            ) : (
              <DropdownMenuItem onClick={() => onUpdateStatus(tour)}>
                <XCircle className="mr-2 h-4 w-4" />
                ยกเลิกเผยแพร่
              </DropdownMenuItem>
            )}

            <DropdownMenuSeparator />

            <DropdownMenuItem
              onClick={() => onDelete(tour)}
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