import type { AccessorKeyColumnDef, Table } from "@tanstack/react-table";
import { VisuallyHidden } from "radix-ui";
import React, { type JSX } from "react";
import { DataTablePagination } from "#components/datatable/components/data-table-pagination";
import { DataTableViewOptions } from "#components/datatable/components/data-table-view-options";
import { TasksTableActionBar } from "#components/datatable/components/tasks-table-action-bar";
import { DataTableFilterMenu } from "#components/datatable/components/toolbar/data-table-filter-menu";
import { DataTableSortList } from "#components/datatable/components/toolbar/data-table-sort-list";
import { ViewToggle } from "#components/datatable/components/view-toggle";
import { GridView } from "#components/datatable/components/views/grid-view";
import { TableView } from "#components/datatable/components/views/table-view";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "#components/shadcn/dialog";
import type {
  BaseEntity,
  FilterCondition,
  FilterNode,
  FilterValue,
} from "#types/entities";
import type { FilterOperator } from "#types/simple-data-table";

export interface PureDataTableProps<TEntity extends BaseEntity> {
  // Data
  columns: AccessorKeyColumnDef<TEntity>[];

  // State
  viewMode: "table" | "grid";

  // Filter state
  filter?: FilterNode;
  addCondition?: (condition: FilterCondition) => void;
  removeCondition?: (field: keyof TEntity) => void;
  resetFilters?: () => void;
  setCondition?: (
    field: keyof TEntity,
    operator: FilterOperator,
    value: FilterValue,
  ) => void;

  // Callbacks
  onViewModeChange: (mode: "table" | "grid") => void;

  // UI Configuration
  icon?: React.ReactNode;
  minimal?: boolean;
  showFilters?: boolean;
  showViewOptions?: boolean;
  showPagination?: boolean;

  // Forms (optional)
  createForm?: React.ComponentType;
  updateForm?: React.ComponentType<{ id: string }>;

  // Grid view customization
  gridRenderer?: (item: TEntity) => React.ReactNode;
  gridHoverRenderer?: (item: TEntity) => React.ReactNode;

  // Dialog state
  dialogOpen?: boolean;
  dialogVariant?: "create" | "update";
  selectedRow?: TEntity;
  onDialogChange?: (open: boolean) => void;

  // Table instance (passed from container)
  table: Table<TEntity>;

  // Optional action bar
  actionBar?: React.ReactNode;
}

/**
 * Pure presentational DataTable component.
 * All state and data fetching is handled by the container.
 */
export function PureDataTable<TEntity extends BaseEntity>({
  columns: _columns,
  viewMode,
  filter,
  addCondition,
  removeCondition,
  resetFilters,
  setCondition,
  onViewModeChange,
  icon,
  minimal = false,
  showFilters = true,
  showViewOptions = true,
  showPagination = true,
  createForm,
  updateForm,
  gridRenderer,
  gridHoverRenderer,
  dialogOpen = false,
  dialogVariant,
  selectedRow,
  onDialogChange,
  table,
  actionBar,
}: PureDataTableProps<TEntity>): JSX.Element {
  const handleViewToggle = () => {
    onViewModeChange(viewMode === "table" ? "grid" : "table");
  };

  return (
    <div className="flex flex-1 flex-col gap-4">
      {/* FILTER TOOLBAR */}
      {(showFilters || showViewOptions) && (
        <div
          aria-orientation="horizontal"
          className="flex gap-2"
          role="toolbar"
        >
          {showFilters && (
            <>
              <DataTableSortList
                align="start"
                table={table}
              />
              <DataTableFilterMenu
                addCondition={
                  addCondition as (condition: FilterCondition) => void
                }
                filter={filter as FilterNode}
                removeCondition={
                  removeCondition as (field: keyof TEntity) => void
                }
                resetFilters={resetFilters as () => void}
                setCondition={setCondition as () => void}
                table={table}
              />
            </>
          )}
          {showViewOptions && (
            <>
              <DataTableViewOptions table={table} />
              <ViewToggle
                onToggle={handleViewToggle}
                view={viewMode}
              />
            </>
          )}
        </div>
      )}

      {/* DATA TABLE */}
      <div className="flex-1 overflow-auto">
        {viewMode === "grid" ? (
          <GridView<TEntity>
            grid={gridRenderer as (item: TEntity) => React.ReactNode}
            gridHover={gridHoverRenderer as (item: TEntity) => React.ReactNode}
            icon={icon}
            table={table}
          />
        ) : (
          <TableView<TEntity> table={table} />
        )}
      </div>

      {/* PAGINATION */}
      {!minimal && showPagination && <DataTablePagination table={table} />}

      {/* DIALOG FOR CREATE/UPDATE */}
      <Dialog
        onOpenChange={onDialogChange || (() => {})}
        open={dialogOpen}
      >
        <VisuallyHidden.Root>
          <DialogDescription />
          <DialogTitle>
            {dialogVariant === "update" ? "Edit" : "Create"}{" "}
            {(
              table.options.meta as {
                entityKey: string;
              }
            ).entityKey ?? "Entity"}
          </DialogTitle>
        </VisuallyHidden.Root>
        <DialogContent
          aria-description="Create/Edit"
          className="p-0"
          title="Create/Edit"
        >
          {dialogVariant === "update" &&
            updateForm &&
            selectedRow &&
            React.createElement(updateForm, { id: selectedRow.id })}

          {dialogVariant === "create" &&
            createForm &&
            React.createElement(createForm)}
        </DialogContent>
      </Dialog>

      {/* ACTION BAR */}
      {actionBar || <TasksTableActionBar table={table} />}
    </div>
  );
}
