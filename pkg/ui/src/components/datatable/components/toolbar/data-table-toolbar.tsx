"use no memo";

import type { Column, Table } from "@tanstack/react-table";
import type { JSX } from "react";

import { useCallback, useMemo } from "react";
import { XCircleIcon } from "#components/custom/icons";
import { DataTableViewOptions } from "#components/datatable/components/data-table-view-options";
import { DataTableDateFilter } from "#components/datatable/components/filters/data-table-date-filter";
import { DataTableFacetedFilter } from "#components/datatable/components/filters/data-table-faceted-filter";
import { DataTableSliderFilter } from "#components/datatable/components/filters/data-table-slider-filter";
import { Button } from "#components/shadcn/button";
import { Input } from "#components/shadcn/input";
import { cn } from "#lib/utils";
import type { BaseEntity } from "#types/entities";

export interface DataTableToolbarProps<TEntity extends BaseEntity> {
  table: Table<TEntity>;
}

interface DataTableToolbarFilterProps<TData> {
  column: Column<TData>;
}

export function DataTableToolbar<TEntity extends BaseEntity>(
  props: DataTableToolbarProps<TEntity>,
): JSX.Element {
  const isFiltered = props.table.getState().columnFilters.length > 0;

  const columns = useMemo(
    () => props.table.getAllColumns().filter((column) => column.getCanFilter()),
    [props.table],
  );

  const onReset = useCallback(() => {
    props.table.resetColumnFilters();
  }, [props.table]);

  return (
    <div
      aria-orientation="horizontal"
      className="flex w-full items-start justify-between gap-2 p-1"
      role="toolbar"
    >
      <div className="flex flex-1 flex-wrap items-center gap-2">
        {columns.map((column) => (
          <DataTableToolbarFilter
            column={column}
            key={column.id}
          />
        ))}
        {isFiltered && (
          <Button
            aria-label="Reset filters"
            className="border-dashed"
            onClick={onReset}
            size="sm"
            variant="outline"
          >
            <XCircleIcon />
            Reset
          </Button>
        )}
      </div>
      <div className="flex items-center gap-2">
        {/* {children} */}
        <DataTableViewOptions table={props.table} />
      </div>
    </div>
  );
}

function DataTableToolbarFilter<TData>({
  column,
}: DataTableToolbarFilterProps<TData>) {
  {
    const columnMeta = column.columnDef.meta;

    const onFilterRender = useCallback(() => {
      if (!columnMeta?.filterVariant) return null;

      switch (columnMeta.filterVariant) {
        case "date":
        case "dateRange":
          return (
            <DataTableDateFilter
              column={column}
              multiple={columnMeta.filterVariant === "dateRange"}
              title={columnMeta.label ?? column.id}
            />
          );

        case "multiSelect":
        case "select":
          return (
            <DataTableFacetedFilter
              column={column}
              multiple={columnMeta.filterVariant === "multiSelect"}
              options={columnMeta.options ?? []}
              title={columnMeta.label ?? column.id}
            />
          );

        case "number":
          return (
            <div className="relative">
              <Input
                className={cn("h-8 w-[120px]", columnMeta.unit && "pr-8")}
                inputMode="numeric"
                onChange={(event) => {
                  column.setFilterValue(event.target.value);
                }}
                placeholder={columnMeta.label}
                type="number"
                value={(column.getFilterValue() as string | undefined) ?? ""}
              />
              {columnMeta.unit && (
                <span className="absolute top-0 right-0 bottom-0 flex items-center rounded-r-md bg-accent px-2 text-muted-foreground text-sm">
                  {columnMeta.unit}
                </span>
              )}
            </div>
          );

        case "range":
          return (
            <DataTableSliderFilter
              column={column}
              title={columnMeta.label ?? column.id}
            />
          );

        case "text":
          return (
            <Input
              className="h-8 w-40 lg:w-56"
              onChange={(event) => {
                column.setFilterValue(event.target.value);
              }}
              placeholder={columnMeta.label}
              value={((column.getFilterValue() as string) || undefined) ?? ""}
            />
          );

        default:
          return null;
      }
    }, [column, columnMeta]);

    return onFilterRender();
  }
}
