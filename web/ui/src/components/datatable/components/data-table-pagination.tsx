"use no memo"

import type { Table } from "@tanstack/react-table"
import type { JSX } from "react"

import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronsLeftIcon,
  ChevronsRightIcon
} from "#components/custom/icons"
import { Button } from "#components/shadcn/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "#components/shadcn/select"
import { cn } from "#lib/utils"

interface DataTablePaginationProps<TData> extends React.ComponentProps<"div"> {
  pageSizeOptions?: number[]
  table: Table<TData>
}

export function DataTablePagination<TData>({
  className,
  pageSizeOptions = [10, 20, 30, 40, 50],
  table,
  ...props
}: DataTablePaginationProps<TData>): JSX.Element {
  return (
    <div
      className={cn(
        "flex w-full flex-col-reverse items-center justify-between gap-4 overflow-auto p-1 sm:flex-row sm:gap-8",
        className
      )}
      {...props}
    >
      <div className="flex-1 text-sm whitespace-nowrap text-muted-foreground">
        {table.getFilteredSelectedRowModel().rows.length} of{" "}
        {table.getFilteredRowModel().rows.length} row(s) selected.
      </div>
      <div className="flex flex-col-reverse items-center gap-4 sm:flex-row sm:gap-6 lg:gap-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium whitespace-nowrap">Rows per page</p>
          <Select
            onValueChange={(value) => {
              table.setPageSize(Number(value))
            }}
            value={table.getState().pagination.pageSize.toString()}
          >
            <SelectTrigger className="h-8 w-[4.5rem] dark:border-none [&[data-size]]:h-8">
              <SelectValue placeholder={table.getState().pagination.pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {pageSizeOptions.map((pageSize) => (
                <SelectItem
                  key={pageSize}
                  value={pageSize.toString()}
                >
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-center justify-center text-sm font-medium">
          Page {table.getState().pagination.pageIndex + 1} of{" "}
          {table.getPageCount()}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            aria-label="Go to first page"
            className="hidden size-8 lg:flex"
            disabled={!table.getCanPreviousPage()}
            onClick={() => {
              table.setPageIndex(0)
            }}
            size="icon"
            variant="ghost"
          >
            <ChevronsLeftIcon />
          </Button>
          <Button
            aria-label="Go to previous page"
            className="size-8"
            disabled={!table.getCanPreviousPage()}
            onClick={() => {
              table.previousPage()
            }}
            size="icon"
            variant="ghost"
          >
            <ChevronLeftIcon />
          </Button>
          <Button
            aria-label="Go to next page"
            className="size-8"
            disabled={!table.getCanNextPage()}
            onClick={() => {
              table.nextPage()
            }}
            size="icon"
            variant="ghost"
          >
            <ChevronRightIcon />
          </Button>
          <Button
            aria-label="Go to last page"
            className="hidden size-8 lg:flex"
            disabled={!table.getCanNextPage()}
            onClick={() => {
              table.setPageIndex(table.getPageCount() - 1)
            }}
            size="icon"
            variant="ghost"
          >
            <ChevronsRightIcon />
          </Button>
        </div>
      </div>
    </div>
  )
}
