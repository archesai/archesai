"use no memo";

import type { UseSuspenseQueryOptions } from "@tanstack/react-query";
import type { AccessorKeyColumnDef, RowData } from "@tanstack/react-table";
import type { JSX } from "react";

import { useState } from "react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { VisuallyHidden } from "radix-ui";

import type { BaseEntity, SearchQuery } from "#types/entities";
import type { DataTableRowAction } from "#types/simple-data-table";

import { DataTablePagination } from "#components/datatable/components/data-table-pagination";
import { DataTableViewOptions } from "#components/datatable/components/data-table-view-options";
import { TasksTableActionBar } from "#components/datatable/components/tasks-table-action-bar";
import { DataTableFilterMenu } from "#components/datatable/components/toolbar/data-table-filter-menu";
import { DataTableSortList } from "#components/datatable/components/toolbar/data-table-sort-list";
import { ViewToggle } from "#components/datatable/components/view-toggle";
// import { ViewToggle } from '#components/datatable/components/view-toggle'
import { GridView } from "#components/datatable/components/views/grid-view";
import { TableView } from "#components/datatable/components/views/table-view";
// import { DataTableToolbar } from '#components/datatable/components/toolbar/data-table-toolbar'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "#components/shadcn/dialog";
import { useDataTable } from "#hooks/use-data-table";
import { useFilterState } from "#hooks/use-filter-state";
import { useToggleView } from "#hooks/use-toggle-view";

declare module "@tanstack/table-core" {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface TableMeta<TData extends RowData> {
    entityKey: string;
    label: string;
  }
}

export interface DataTableProps<TEntity extends BaseEntity> {
  actionBar?: React.ReactNode;
  columns: AccessorKeyColumnDef<TEntity>[];
  createForm?: React.ComponentType;
  deleteItem?: (id: string) => Promise<void>;
  entityKey?: string;
  getQueryOptions: (query: SearchQuery) => UseSuspenseQueryOptions<{
    data: TEntity[];
    meta: {
      total: number;
    };
  }>;
  grid?: (item: TEntity) => React.ReactNode;
  gridHover?: (item: TEntity) => React.ReactNode;
  handleSelect: (item: TEntity) => void;
  icon: React.ReactNode;
  minimal?: boolean;
  updateForm?: React.ComponentType<{ id: string }>;
}

export function DataTable<TEntity extends BaseEntity>(
  props: DataTableProps<TEntity>,
): JSX.Element {
  const [rowAction, setRowAction] =
    useState<DataTableRowAction<TEntity> | null>(null);

  const { view } = useToggleView();

  const { searchQuery } = useFilterState<TEntity>();
  const queryOptions = props.getQueryOptions(searchQuery);
  const { data: queryData } = useSuspenseQuery(queryOptions);

  const { table } = useDataTable<TEntity>({
    clearOnDefault: true,
    columns: props.columns,
    data: queryData.data,
    getRowId: (originalRow) => originalRow.id,
    initialState: {
      columnPinning: { right: ["actions"] },
      sorting: [
        { desc: true, id: "createdAt" as Extract<keyof TEntity, string> },
      ],
    },
    pageCount:
      Math.ceil(queryData.meta.total / (searchQuery.page?.size ?? 10)) || 1,
    shallow: false,
  });

  return (
    <div className="flex flex-1 flex-col gap-4">
      {/* FILTER TOOLBAR */}
      <div
        aria-orientation="horizontal"
        className="flex gap-2"
        role="toolbar"
      >
        <DataTableSortList
          align="start"
          table={table}
        />
        <DataTableFilterMenu table={table} />
        <DataTableViewOptions table={table} />
        <ViewToggle />
      </div>

      {/* DATA TABLE */}
      <div className="flex-1 overflow-auto">
        {view === "grid" ? (
          <GridView<TEntity>
            icon={props.icon}
            table={table}
          />
        ) : (
          <TableView<TEntity> table={table} />
        )}
      </div>

      {/* PAGINATION - Now uses filterState directly */}
      {!props.minimal && <DataTablePagination table={table} />}

      {/* DIALOG AND ACTION BAR remain the same */}
      <Dialog
        onOpenChange={() => {
          setRowAction(null);
        }}
        open={
          rowAction?.variant === "update" || rowAction?.variant === "custom"
        }
      >
        <VisuallyHidden.Root>
          <DialogDescription />
          <DialogTitle>
            {rowAction?.variant === "update" ? "Edit" : "Create"}{" "}
            {table.options.meta?.entityKey ?? "Entity"}
          </DialogTitle>
        </VisuallyHidden.Root>
        <DialogContent
          aria-description="Create/Edit"
          className="p-0"
          title="Create/Edit"
        >
          {rowAction?.variant === "update" && props.updateForm && (
            <props.updateForm id={rowAction.row.original.id} />
          )}

          {rowAction?.variant === "create" && props.createForm && (
            <props.createForm />
          )}
        </DialogContent>
      </Dialog>
      <TasksTableActionBar table={table} />
    </div>
  );
}
