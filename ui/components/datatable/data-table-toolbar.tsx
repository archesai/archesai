"use client";

import { DataTableViewOptions } from "@/components/datatable/data-table-view-options";
import { DatePickerWithRange } from "@/components/datatable/date-range-picker";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { useToggleView } from "@/hooks/useToggleView";
import { Cross2Icon } from "@radix-ui/react-icons";
import { Table } from "@tanstack/react-table";
import { GridIcon, ListIcon } from "lucide-react";

interface DataTableToolbarProps<TData> {
  createForm?: React.ReactNode;
  data: TData[];
  itemType: string;
  setFormOpen: (open: boolean) => void;
  table: Table<TData>;
}

export function DataTableToolbar<TData>({
  createForm,
  data,
  itemType,
  setFormOpen,
  table,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  const { query, setQuery } = useFilterItems();
  const { selectedAllItems, selectedSomeItems, toggleSelectAll } =
    useSelectItems({ items: data || [] });

  return (
    <div className="flex flex-wrap items-center gap-2">
      <Checkbox
        aria-label="Select all"
        checked={selectedAllItems || (selectedSomeItems && "indeterminate")}
        onCheckedChange={() => toggleSelectAll()}
      />
      <Input
        className="h-8 flex-1"
        onChange={(event) => setQuery(event.target.value)}
        placeholder={`Search ${itemType}s...`}
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
          className="flex h-8 gap-2 p-2"
          onClick={() => table.resetColumnFilters()}
          variant="ghost"
        >
          <span>Reset</span>
          <Cross2Icon className="h-5 w-5" />
        </Button>
      )}

      <DatePickerWithRange />
      <ViewToggle />
      <DataTableViewOptions table={table} />
      {createForm ? (
        <Button
          className="capitalize"
          onClick={() => setFormOpen(true)}
          size="sm"
        >
          Create {itemType.toLowerCase()}
        </Button>
      ) : null}
    </div>
  );
}

export function ViewToggle() {
  const { setView, view } = useToggleView();
  return (
    <div className="hidden h-8 gap-2 md:flex">
      <Button
        className={`flex h-full items-center justify-center transition-colors ${
          view === "table"
            ? "bg-secondary text-primary"
            : "bg-transparent text-muted-foreground"
        }`}
        onClick={() => setView("table")}
        size="icon"
        variant={"secondary"}
      >
        <ListIcon className="h-5 w-5" />
      </Button>
      <Button
        className={`flex h-full items-center justify-center transition-colors ${
          view === "grid"
            ? "bg-secondary text-primary"
            : "bg-transparent text-muted-foreground"
        }`}
        onClick={() => setView("grid")}
        size="icon"
        variant={"secondary"}
      >
        <GridIcon className="h-5 w-5" />
      </Button>
    </div>
  );
}
