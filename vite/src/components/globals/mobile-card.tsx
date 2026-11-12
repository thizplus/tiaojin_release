import { type ColumnDef, flexRender } from "@tanstack/react-table";
import { cn } from "@/lib/utils";

interface MobileCardProps<TData, TValue = unknown> {
  data: TData;
  columns: ColumnDef<TData, TValue>[];
  columnLabels?: Record<string, string>;
  onClick?: () => void;
}

export function MobileCard<TData, TValue = unknown>({
  data,
  columns,
  columnLabels = {},
  onClick
}: MobileCardProps<TData, TValue>) {
  // แยก columns ตาม type
  const actionsColumn = columns.find(col => col.id === "actions");
  const dataColumns = columns.filter(col =>
    col.id !== "actions" &&
    'accessorKey' in col && col.accessorKey
  );

  // สร้าง mock cell context สำหรับ flexRender
  const createCellContext = (column: ColumnDef<TData, TValue>) => {
    const accessorKey = 'accessorKey' in column ? column.accessorKey as string : undefined;

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const mockContext: any = {
      getValue: () => {
        if (accessorKey && data) {
          return (data as Record<string, unknown>)[accessorKey];
        }
        return undefined;
      },
      row: {
        original: data,
        getValue: (key: string) => {
          if (data) {
            return (data as Record<string, unknown>)[key];
          }
          return undefined;
        },
      },
      column,
      table: {},
      cell: {
        getValue: () => {
          if (accessorKey && data) {
            return (data as Record<string, unknown>)[accessorKey];
          }
          return undefined;
        },
      },
      renderValue: () => {
        if (accessorKey && data) {
          return (data as Record<string, unknown>)[accessorKey];
        }
        return undefined;
      },
    };

    return mockContext;
  };

  return (
    <div
      className={cn(
        "rounded-lg border p-4 space-y-3 transition-colors",
        "bg-background hover:bg-muted/50",
        onClick && "cursor-pointer"
      )}
      onClick={onClick}
    >
      {/* Header with actions */}
      {actionsColumn && (
        <div className="flex items-center justify-end">
          <div onClick={(e) => e.stopPropagation()}>
            {flexRender(actionsColumn.cell, createCellContext(actionsColumn))}
          </div>
        </div>
      )}

      {/* Data fields */}
      <div className="space-y-2">
        {dataColumns.map((column) => {
          const accessorKey = 'accessorKey' in column ? column.accessorKey as string : undefined;
          const label = columnLabels[column.id as string] ||
                      (typeof column.header === 'string' ? column.header : column.id) ||
                      accessorKey;

          const cellContext = createCellContext(column);
          const value = cellContext.getValue();

          // ข้าม field ที่ว่าง
          if (value === undefined || value === null || value === '') {
            return null;
          }

          return (
            <div key={column.id || ('accessorKey' in column ? column.accessorKey as string : undefined)} className="flex flex-col space-y-1">
              <div className="text-sm font-medium text-muted-foreground">
                {label}
              </div>
              <div className="text-sm">
                {column.cell ?
                  flexRender(column.cell, cellContext) :
                  String(value)
                }
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default MobileCard;