import type { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, MoreHorizontal, Eye, CheckCircle, XCircle, Pencil } from "lucide-react";
import { getStorageUrl } from "@/lib/utils";
import { format } from "date-fns";
import { th } from "date-fns/locale";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type { CompanyProfile } from "@/types";


export interface CompanyProfileActionsHandlers {
  onView?: (profile: CompanyProfile) => void;
  onEdit?: (profile: CompanyProfile) => void;
  onApprove?: (profile: CompanyProfile) => void;
}

export const getCompanyProfileColumns = (
  handlers: CompanyProfileActionsHandlers,
): ColumnDef<CompanyProfile>[] => {
  return [
    {
      accessorKey: "logo",
      header: "โลโก้",
      cell: ({ row }) => {
        const logo = row.getValue("logo") as string;
        return logo ? (
          <div className="flex items-center">
            <img
              src={getStorageUrl(logo)}
              alt="Company Logo"
              className="h-10 w-10 rounded-md object-cover"
            />
          </div>
        ) : (
          <span className="text-muted-foreground text-sm">ไม่มีโลโก้</span>
        );
      },
    },
    {
      accessorKey: "name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          ชื่อบริษัท
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => <div className="font-medium">{row.getValue("name")}</div>,
    },
    {
      accessorKey: "tax_id",
      header: "เลขประจำตัวผู้เสียภาษี",
      cell: ({ row }) => {
        const taxId = row.getValue("tax_id") as string;
        return (
          <div className="font-mono text-sm">
            {taxId.replace(/(\d{1})(\d{4})(\d{5})(\d{2})(\d{1})/, "$1-$2-$3-$4-$5")}
          </div>
        );
      },
    },
    {
      accessorKey: "phone",
      header: "เบอร์โทรศัพท์",
      cell: ({ row }) => {
        const phone = row.getValue("phone") as string;
        return <div className="font-mono text-sm">{phone}</div>;
      },
    },
    {
      accessorKey: "province",
      header: "จังหวัด",
      cell: ({ row }) => <div>{row.getValue("province")}</div>,
    },
    {
      accessorKey: "atta",
      header: "ATTA",
      cell: ({ row }) => {
        const value = row.getValue("atta") as string;
        if (!value) return <XCircle className="h-4 w-4 text-destructive/50" />;
        return (
          <button onClick={() => window.open(getStorageUrl(value), '_blank', 'noopener,noreferrer')} className="cursor-pointer" title="เปิดไฟล์ ATTA">
            <CheckCircle className="h-4 w-4 text-green-600" />
          </button>
        );
      },
    },
    {
      accessorKey: "dbd",
      header: "DBD",
      cell: ({ row }) => {
        const value = row.getValue("dbd") as string;
        if (!value) return <XCircle className="h-4 w-4 text-destructive/50" />;
        return (
          <button onClick={() => window.open(getStorageUrl(value), '_blank', 'noopener,noreferrer')} className="cursor-pointer" title="เปิดไฟล์ DBD">
            <CheckCircle className="h-4 w-4 text-green-600" />
          </button>
        );
      },
    },
    {
      accessorKey: "pp20",
      header: "PP20",
      cell: ({ row }) => {
        const value = row.getValue("pp20") as string;
        if (!value) return <XCircle className="h-4 w-4 text-destructive/50" />;
        return (
          <button onClick={() => window.open(getStorageUrl(value), '_blank', 'noopener,noreferrer')} className="cursor-pointer" title="เปิดไฟล์ PP20">
            <CheckCircle className="h-4 w-4 text-green-600" />
          </button>
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
          วันที่ส่งคำขอ
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
        const profile = row.original;
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
              <DropdownMenuItem onClick={() => handlers.onView?.(profile)}>
                <Eye className="mr-1 h-3.5 w-3.5" />
                ดูรายละเอียด
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handlers.onEdit?.(profile)}>
                <Pencil className="mr-1 h-3.5 w-3.5" />
                แก้ไข
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => handlers.onApprove?.(profile)}
                className="text-green-600 focus:text-green-600"
              >
                <CheckCircle className="mr-1 h-3.5 w-3.5 text-green-600" />
                อนุมัติ
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];
};