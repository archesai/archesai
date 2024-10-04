"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useToggleView } from "@/hooks/useToggleView";
import { DropdownMenuTrigger } from "@radix-ui/react-dropdown-menu";
import { Cross2Icon } from "@radix-ui/react-icons";
import { MixerHorizontalIcon } from "@radix-ui/react-icons";
import { Table } from "@tanstack/react-table";

import { DatePickerWithRange } from "./date-range-picker";

interface DataTableToolbarProps<TData> {
  itemType?: string;
  table: Table<TData>;
}

export function DataTableToolbar<TData>({
  itemType,
  table,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  const { query, setQuery } = useFilterItems();

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 items-center space-x-2 flex-col md:flex-row gap-2 md:gap-0">
        <Input
          className="h-8"
          onChange={(event) => setQuery(event.target.value)}
          placeholder={`Filter ${itemType}s...`}
          value={query}
        />
        {/* {table.getColumn("llmBase") && (
          <DataTableFacetedFilter
            column={table.getColumn("llmBase")}
            options={[
              {
                label: "GPT-4",
                value: "GPT-4",
              },
            ]}
            title="Language Model"
          />
        )} */}
        {isFiltered && (
          <Button
            className="h-8 px-2 lg:px-3"
            onClick={() => table.resetColumnFilters()}
            variant="ghost"
          >
            Reset
            <Cross2Icon className="ml-2 h-4 w-4" />
          </Button>
        )}
        <DatePickerWithRange />
        <DataTableViewOptions table={table} />
      </div>
    </div>
  );
}

interface DataTableViewOptionsProps<TData> {
  table: Table<TData>;
}

export function DataTableViewOptions<TData>({
  table,
}: DataTableViewOptionsProps<TData>) {
  const { setView, view } = useToggleView();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button className="ml-auto h-8" size="sm" variant="outline">
          <MixerHorizontalIcon className="mr-2 h-4 w-4" />
          View
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-[150px]">
        <DropdownMenuLabel>Toggle columns</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {table
          .getAllColumns()
          .filter(
            (column) =>
              typeof column.accessorFn !== "undefined" && column.getCanHide()
          )
          .map((column) => {
            return (
              <DropdownMenuCheckboxItem
                checked={column.getIsVisible()}
                className="capitalize"
                key={column.id}
                onCheckedChange={(value) => column.toggleVisibility(!!value)}
              >
                {column.id}
              </DropdownMenuCheckboxItem>
            );
          })}
        <DropdownMenuSeparator />
        <DropdownMenuLabel>Toggle view</DropdownMenuLabel>
        <DropdownMenuCheckboxItem
          checked={view === "grid"}
          className="capitalize"
          disabled={view === "grid"}
          key={"grid"}
          onCheckedChange={(value) => {
            setView(value ? "grid" : "table");
          }}
        >
          {"Grid"}
        </DropdownMenuCheckboxItem>
        <DropdownMenuCheckboxItem
          checked={view === "table"}
          className="capitalize"
          disabled={view === "table"}
          key={"table"}
          onCheckedChange={(value) => {
            setView(value ? "table" : "grid");
          }}
        >
          {"Table"}
        </DropdownMenuCheckboxItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
