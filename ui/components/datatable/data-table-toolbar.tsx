"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { useToggleView } from "@/hooks/useToggleView";
import { Cross2Icon } from "@radix-ui/react-icons";
import { Table } from "@tanstack/react-table";

import { Checkbox } from "../ui/checkbox";
import { DataTableFacetedFilter } from "./data-table-faceted-filter";
import { DatePickerWithRange } from "./date-range-picker";

interface DataTableToolbarProps<TData> {
  data: TData[];
  itemType?: string;
  table: Table<TData>;
}

export function DataTableToolbar<TData>({
  data,
  itemType,
  table,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  const { query, setQuery } = useFilterItems();
  const { selectedAllItems, selectedSomeItems, toggleSelectAll } =
    useSelectItems({ items: data || [] });

  return (
    <div className="flex flex-wrap items-center space-x-2 gap-2 md:gap-0">
      <Checkbox
        aria-label="Select all"
        checked={selectedAllItems || (selectedSomeItems && "indeterminate")}
        className="translate-y-[2px]"
        onCheckedChange={() => toggleSelectAll()}
      />
      <Input
        className="h-8 flex-1"
        onChange={(event) => setQuery(event.target.value)}
        placeholder={`Search ${itemType}s...`}
        value={query}
      />
      {table.getColumn("llmBase") && (
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
      )}
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
  );
}

interface DataTableViewOptionsProps<TData> {
  table: Table<TData>;
}

export function DataTableViewOptions<
  TData,
>({}: DataTableViewOptionsProps<TData>) {
  const { toggleView } = useToggleView();
  return <Switch onCheckedChange={() => toggleView()} />;
}
