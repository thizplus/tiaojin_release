import type { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, MoreHorizontal, Pencil, Trash, Shield, UserCheck, UserX, Crown, ShieldCheck, User as UserIcon, Briefcase, KeyRound } from "lucide-react";
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
import type { User } from "@/types";
import { getUserStatusConfig } from "@/utils/status-helpers";


export interface UserActionsHandlers {
  onEdit?: (user: User) => void;
  onDelete?: (user: User) => void;
  onChangeStatus?: (user: User) => void;
  onChangeRole?: (user: User) => void;
  onResetPassword?: (user: User) => void;
}


const getRoleLabel = (role: string) => {
  switch (role) {
    case "super":
      return "Super Admin";
    case "admin":
      return "Admin";
    case "agent":
      return "Agent";
    case "sales":
      return "Sales";
    default:
      return role;
  }
};


export const getUserColumns = (
  handlers: UserActionsHandlers
): ColumnDef<User>[] => {
  return [
    {
      accessorKey: "username",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          ชื่อผู้ใช้
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => <div className="font-medium">{row.getValue("username")}</div>,
    },
    {
      accessorKey: "email",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          อีเมล
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const email = row.getValue("email") as string | null;
        return email ? (
          <div>{email}</div>
        ) : (
          <span className="text-muted-foreground text-sm">-</span>
        );
      },
    },
    {
      accessorKey: "company_profile",
      header: "บริษัท",
      cell: ({ row }) => {
        const companyProfile = row.getValue("company_profile") as User["company_profile"];
        return companyProfile?.name ? (
          <div>
            <div className="font-medium">{companyProfile.name}</div>
            <div className="text-sm text-muted-foreground">{companyProfile.tax_id}</div>
          </div>
        ) : (
          <span className="text-muted-foreground text-sm">ไม่มีข้อมูลบริษัท</span>
        );
      },
    },
    {
      accessorKey: "role",
      header: "บทบาท",
      cell: ({ row }) => {
        const role = row.getValue("role") as string;

        const getRoleIcon = (role: string) => {
          switch (role) {
            case "super":
              return <Crown className="h-3.5 w-3.5 text-muted-foreground" />;
            case "admin":
              return <ShieldCheck className="h-3.5 w-3.5 text-muted-foreground" />;
            case "agent":
              return <UserIcon className="h-3.5 w-3.5 text-muted-foreground" />;
            case "sales":
              return <Briefcase className="h-3.5 w-3.5 text-muted-foreground" />;
            default:
              return <UserIcon className="h-3.5 w-3.5 text-muted-foreground" />;
          }
        };

      

        return (
          <Badge variant="outline" className="px-1.5">
            {getRoleIcon(role)}
            <span className={`ml-1 text-muted-foreground`}>{getRoleLabel(role)}</span>
          </Badge>
        );
      },
    },
    {
      accessorKey: "status",
      header: "สถานะ",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;

        const statusConfig = getUserStatusConfig(status);

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
        const user = row.original;
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
              <DropdownMenuItem onClick={() => handlers.onEdit?.(user)}>
                <Pencil className="mr-1 h-3.5 w-3.5" />
                แก้ไข
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handlers.onChangeRole?.(user)}>
                <Shield className="mr-1 h-3.5 w-3.5" />
                เปลี่ยนบทบาท
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handlers.onChangeStatus?.(user)}>
                {user.status === "active" ? (
                  <>
                    <UserX className="mr-1 h-3.5 w-3.5" />
                    ระงับการใช้งาน
                  </>
                ) : (
                  <>
                    <UserCheck className="mr-1 h-3.5 w-3.5" />
                    เปิดใช้งาน
                  </>
                )}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handlers.onResetPassword?.(user)}>
                <KeyRound className="mr-1 h-3.5 w-3.5" />
                รีเซ็ตรหัสผ่าน
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => handlers.onDelete?.(user)}
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
