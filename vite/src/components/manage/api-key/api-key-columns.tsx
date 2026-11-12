import type { ColumnDef } from "@tanstack/react-table"
import { ArrowUpDown, MoreHorizontal, Pencil, Trash, CheckCircle2, XCircle } from "lucide-react"
import { format } from "date-fns"
import { th } from "date-fns/locale"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { ApiKey } from "@/types/models/apiKey"


export interface ApiKeyActionsHandlers {
  onEdit?: (apiKey: ApiKey) => void
  onDelete?: (apiKey: ApiKey) => void
}

export const getApiKeyColumns = (
  handlers: ApiKeyActionsHandlers
): ColumnDef<ApiKey>[] => {
  return [
    {
      accessorKey: "name",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          ชื่อ API Key
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      ),
      cell: ({ row }) => {
        const name = row.getValue("name") as string
        const description = row.original.description
        return (
          <div>
            <div className="font-medium">{name}</div>
            {description && (
              <div className="text-sm text-muted-foreground">{description}</div>
            )}
          </div>
        )
      },
    },
    {
      accessorKey: "key_preview",
      header: "API Key",
      cell: ({ row }) => {
        const keyPreview = row.getValue("key_preview") as string
        return (
          <div className="flex items-center gap-2">
            <code className="rounded bg-muted px-2 py-1 text-sm font-mono">
              {keyPreview}
            </code>
            <span className="text-xs text-muted-foreground ml-2">(preview)</span>
          </div>
        )
      },
    },
    {
      accessorKey: "is_active",
      header: "สถานะ",
      cell: ({ row }) => {
        const isActive = row.getValue("is_active") as boolean
        return (
          <Badge variant={isActive ? "default" : "secondary"} className="gap-1">
            {isActive ? (
              <>
                <CheckCircle2 className="h-3 w-3" />
                เปิดใช้งาน
              </>
            ) : (
              <>
                <XCircle className="h-3 w-3" />
                ปิดใช้งาน
              </>
            )}
          </Badge>
        )
      },
    },
    {
      accessorKey: "permissions",
      header: "สิทธิ์",
      cell: ({ row }) => {
        const permissions = row.getValue("permissions") as string[]
        return (
          <div className="flex flex-wrap gap-1">
            {permissions.slice(0, 3).map((permission, index) => (
              <Badge key={index} variant="outline" className="text-xs">
                {permission}
              </Badge>
            ))}
            {permissions.length > 3 && (
              <Badge variant="outline" className="text-xs">
                +{permissions.length - 3}
              </Badge>
            )}
          </div>
        )
      },
    },
    {
      accessorKey: "last_used_at",
      header: "ใช้งานล่าสุด",
      cell: ({ row }) => {
        const lastUsedAt = row.getValue("last_used_at") as string | null
        if (!lastUsedAt) {
          return <span className="text-sm text-muted-foreground">ยังไม่เคยใช้</span>
        }
        return (
          <span className="text-sm">
            {format(new Date(lastUsedAt), "d MMM yyyy HH:mm", { locale: th })}
          </span>
        )
      },
    },
    {
      accessorKey: "expires_at",
      header: "วันหมดอายุ",
      cell: ({ row }) => {
        const expiresAt = row.getValue("expires_at") as string | null
        if (!expiresAt) {
          return <span className="text-sm text-muted-foreground">ไม่มีกำหนด</span>
        }
        const isExpired = new Date(expiresAt) < new Date()
        return (
          <span className={`text-sm ${isExpired ? "text-red-500" : ""}`}>
            {format(new Date(expiresAt), "d MMM yyyy", { locale: th })}
          </span>
        )
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
        )
      },
    },
    {
      id: "actions",
      header: "จัดการ",
      cell: ({ row }) => {
        const apiKey = row.original
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
              <DropdownMenuItem onClick={() => handlers.onEdit?.(apiKey)}>
                <Pencil className="mr-2 h-3.5 w-3.5" />
                แก้ไข
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => handlers.onDelete?.(apiKey)}
                className="text-destructive focus:text-destructive"
              >
                <Trash className="mr-2 h-3.5 w-3.5 text-destructive" />
                ลบ
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}