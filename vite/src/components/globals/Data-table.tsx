"use client";

import * as React from "react";
import {
  type ColumnDef,
  type ColumnFiltersState,
  type SortingState,
  type VisibilityState,
  type PaginationState,
  type OnChangeFn, // เพิ่ม import OnChangeFn
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { ChevronDown, Plus, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-react";

import { Button } from "@/components/ui/button";

import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { MobileCard } from "./mobile-card";
import { useIsMobile } from "@/hooks/use-mobile";


// เพิ่ม interface สำหรับ pagination
interface PaginationProps {
  pageIndex: number;
  pageSize: number;
  pageCount: number;
  canPreviousPage: boolean;
  canNextPage: boolean;
  setPageIndex: (pageIndex: number) => void;
  setPageSize: (pageSize: number) => void;
  previousPage: () => void;
  nextPage: () => void;
  getPageCount: () => number;
}

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  searchKey?: string;
  searchPlaceholder?: string;
  showColumnToggle?: boolean;
  columnToggleText?: string;
  showAddButton?: boolean;
  addButtonText?: string;
  onAdd?: () => void;
  columnLabels?: Record<string, string>;
  // เพิ่ม props สำหรับการจัดการ pagination จากภายนอก
  pagination?: PaginationProps;
  // เพิ่ม getRowId สำหรับกำหนด ID ให้กับแถว
  getRowId?: (originalRow: TData) => string;
}

export function DataTable<TData, TValue>({
  columns,
  data,

  showColumnToggle = true,
  columnToggleText = "คอลัมน์",
  showAddButton = false,
  addButtonText = "เพิ่ม",
  onAdd,
  columnLabels = {},
  pagination: externalPagination,
  getRowId,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
  // ใช้ pagination state ภายในถ้าไม่มี external pagination
  const [pagination, setPagination] = React.useState<PaginationState>({
    pageIndex: 0,
    pageSize: 25,
  });

  // ตรวจสอบว่าเป็น mobile หรือไม่
  const isMobile = useIsMobile();



  // กำหนดค่า pagination สำหรับ table
  // เมื่อใช้ external pagination ต้องใช้ pageSize เท่ากับจำนวน data เพื่อแสดงทั้งหมด
  // เพราะ data จาก server มาเป็น page แล้ว ไม่ต้องการให้ table ทำ pagination ซ้ำ
  const paginationState = externalPagination
    ? {
      pageIndex: 0, // ใช้ 0 เสมอเพราะ data คือ current page แล้ว
      pageSize: data.length || externalPagination.pageSize, // ใช้ data.length เพื่อแสดงทั้งหมด
    }
    : pagination;

  // ฟังก์ชันจัดการเมื่อมีการเปลี่ยนแปลง pagination
  const handlePaginationChange: OnChangeFn<PaginationState> = React.useCallback(
    (updaterOrValue) => {
      // ตรวจสอบว่า updaterOrValue เป็นฟังก์ชันหรือค่า
      const updatedPagination = typeof updaterOrValue === 'function'
        ? updaterOrValue(paginationState)
        : updaterOrValue;
        
      if (externalPagination) {
        if (updatedPagination.pageIndex !== externalPagination.pageIndex) {
          externalPagination.setPageIndex(updatedPagination.pageIndex);
        }
        if (updatedPagination.pageSize !== externalPagination.pageSize) {
          externalPagination.setPageSize(updatedPagination.pageSize);
        }
      } else {
        setPagination(updatedPagination);
      }
    },
    [externalPagination, paginationState]
  );



  const table = useReactTable({
    data,
    columns,
    pageCount: externalPagination?.pageCount,
    getRowId: getRowId,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      pagination: paginationState,
    },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onPaginationChange: handlePaginationChange,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    manualPagination: !!externalPagination,
  });


  return (
    <div className="w-full">
      <div className="flex flex-wrap items-center justify-between py-4 gap-4 md:gap-0">
      

        <div className="flex justify-between w-full items-center space-x-2 md:w-auto">

          {showColumnToggle && (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline">
                  {columnToggleText} <ChevronDown className="ml-2 h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {table
                  .getAllColumns()
                  .filter((column) => column.getCanHide())
                  .map((column) => {
                    return (
                      <DropdownMenuCheckboxItem
                        key={column.id}
                        className="capitalize"
                        checked={column.getIsVisible()}
                        onCheckedChange={(value) =>
                          column.toggleVisibility(!!value)
                        }
                      >
                        {columnLabels[column.id] || column.id}
                      </DropdownMenuCheckboxItem>
                    );
                  })}
              </DropdownMenuContent>
            </DropdownMenu>
          )}

          {showAddButton && (
            <Button onClick={onAdd}>
              <Plus className="mr-2 h-4 w-4" />
              {addButtonText}
            </Button>
          )}
        </div>
      </div>



      {/* Desktop Table View */}
      {!isMobile ? (
        <div className="rounded-md border border-border">
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                    </TableHead>
                  ))}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    className="cursor-pointer hover:bg-muted/50"
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    ไม่พบข้อมูล
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      ) : (
        /* Mobile Card View */
        <div className="space-y-3">
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <MobileCard
                key={row.id}
                data={row.original}
                columns={columns}
                columnLabels={columnLabels}
              />
            ))
          ) : (
            <div className="rounded-lg border p-8 text-center text-muted-foreground">
              ไม่พบข้อมูล
            </div>
          )}
        </div>
      )}

      <div className="flex flex-col gap-4 py-4 md:flex-row md:items-center md:justify-between">

        <div className="flex flex-col gap-4 items-center md:flex-row md:gap-6">
          {/* ส่วนของ Rows per page */}
          <div className="flex items-center gap-2">
            <Label htmlFor="rows-per-page" className="text-sm font-medium">
              รายการต่อหน้า
            </Label>
            <Select
              value={`${externalPagination ? externalPagination.pageSize : table.getState().pagination.pageSize}`}
              onValueChange={(value) => {
                const newPageSize = Number(value);
                if (externalPagination) {
                  externalPagination.setPageSize(newPageSize);
                } else {
                  table.setPageSize(newPageSize);
                }
              }}
            >
              <SelectTrigger size="sm" className="w-20" id="rows-per-page">
                <SelectValue />
              </SelectTrigger>
              <SelectContent side="top">
                {[25, 50, 75, 100].map((pageSize) => (
                  <SelectItem key={pageSize} value={`${pageSize}`}>
                    {pageSize}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* แสดงหน้าปัจจุบัน/จำนวนหน้าทั้งหมด */}
          <div className="flex items-center justify-center text-sm font-medium whitespace-nowrap">
            หน้า {externalPagination ? externalPagination.pageIndex + 1 : table.getState().pagination.pageIndex + 1} จาก{" "}
            {externalPagination ? externalPagination.pageCount : table.getPageCount()}
          </div>

          {/* ส่วนของ Pagination */}
          <div className="flex items-center gap-1 md:gap-2">
            <Button
              variant="outline"
              className="h-8 w-8 p-0"
              onClick={() => {
                if (externalPagination) {
                  externalPagination.setPageIndex(0);
                } else {
                  table.setPageIndex(0);
                }
              }}
              disabled={!(externalPagination ? externalPagination.canPreviousPage : table.getCanPreviousPage())}
            >
              <span className="sr-only">ไปยังหน้าแรก</span>
              <ChevronsLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              className="h-8 w-8 p-0"
              onClick={() => externalPagination ? externalPagination.previousPage() : table.previousPage()}
              disabled={!(externalPagination ? externalPagination.canPreviousPage : table.getCanPreviousPage())}
            >
              <span className="sr-only">ไปยังหน้าก่อนหน้า</span>
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              className="h-8 w-8 p-0"
              onClick={() => externalPagination ? externalPagination.nextPage() : table.nextPage()}
              disabled={!(externalPagination ? externalPagination.canNextPage : table.getCanNextPage())}
            >
              <span className="sr-only">ไปยังหน้าถัดไป</span>
              <ChevronRight className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              className="h-8 w-8 p-0"
              onClick={() => {
                if (externalPagination) {
                  const lastPage = Math.max(0, externalPagination.pageCount - 1);
                  console.log(`DataTable: Last page button clicked, going to page ${lastPage}`);
                  externalPagination.setPageIndex(lastPage);
                } else {
                  const lastPage = Math.max(0, table.getPageCount() - 1);
                  console.log(`DataTable: Last page button clicked, going to page ${lastPage}`);
                  table.setPageIndex(lastPage);
                }
              }}
              disabled={!(externalPagination ? externalPagination.canNextPage : table.getCanNextPage())}
            >
              <span className="sr-only">ไปยังหน้าสุดท้าย</span>
              <ChevronsRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}