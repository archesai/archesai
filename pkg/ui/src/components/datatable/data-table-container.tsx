import type { UseSuspenseQueryOptions } from "@tanstack/react-query";
import { useSuspenseQuery } from "@tanstack/react-query";
import type { AccessorKeyColumnDef } from "@tanstack/react-table";
import type { ComponentType, JSX } from "react";
import { useState } from "react";
import { useDataTable } from "../../hooks/use-data-table";
import { useFilterState } from "../../hooks/use-filter-state";
import { useToggleView } from "../../hooks/use-toggle-view";
import type { BaseEntity, FilterNode, SearchQuery } from "../../types/entities";
import { PureDataTable } from "./pure-data-table";

export interface DataTableContainerProps<TEntity extends BaseEntity> {
  columns: AccessorKeyColumnDef<TEntity>[];
  getQueryOptions: (query: SearchQuery) => UseSuspenseQueryOptions<{
    data: TEntity[];
    meta: {
      total: number;
    };
  }>;
  entityKey?: string;
  handleSelect?: (item: TEntity) => void;
  deleteItem?: (id: string) => Promise<void>;
  createForm?: React.ComponentType;
  updateForm?: React.ComponentType<{
    id: string;
  }>;
  icon?: React.ReactNode;
  minimal?: boolean;
  grid?: (item: TEntity) => React.ReactNode;
  gridHover?: (item: TEntity) => React.ReactNode;
  actionBar?: React.ReactNode;
}

/**
 * Container component that handles all business logic for DataTable.
 * Fetches data, manages state, and passes everything to PureDataTable.
 */
export function DataTableContainer<TEntity extends BaseEntity>({
  columns,
  getQueryOptions,
  entityKey: _entityKey,
  handleSelect: _handleSelect,
  deleteItem: _deleteItem,
  createForm,
  updateForm,
  icon,
  minimal = false,
  grid,
  gridHover,
  actionBar,
}: DataTableContainerProps<TEntity>): JSX.Element {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogVariant, setDialogVariant] = useState<"create" | "update">(
    "create",
  );
  const [selectedRow, setSelectedRow] = useState<TEntity | undefined>();

  // Use hooks from platform
  const { view, toggleView } = useToggleView();
  const {
    // filter,
    searchQuery,
    setPage,
    setPageSize,
    filter,
    addCondition,
    removeCondition,
    resetFilters,
    setCondition,
  } = useFilterState<TEntity>();

  // Fetch data
  const queryOptions = getQueryOptions(searchQuery);
  const { data: queryData } = useSuspenseQuery(queryOptions);

  // Setup table with controlled pagination
  const { table } = useDataTable<TEntity>({
    clearOnDefault: true,
    columns,
    data: queryData.data,
    getRowId: (originalRow: TEntity) => originalRow.id,
    initialState: {
      columnPinning: {
        right: ["actions"],
      },
      sorting: [
        {
          desc: true,
          id: "createdAt" as Extract<keyof TEntity, string>,
        },
      ],
    },
    onPaginationChange: (updater) => {
      const state = table.getState().pagination;
      const newState = typeof updater === "function" ? updater(state) : updater;

      if (newState.pageIndex !== state.pageIndex) {
        setPage(newState.pageIndex + 1);
      }
      if (newState.pageSize !== state.pageSize) {
        setPageSize(newState.pageSize);
      }
    },
    pageCount:
      Math.ceil(queryData.meta.total / (searchQuery.page?.size ?? 10)) || 1,
    shallow: false,
  });

  const handleDialogChange = (open: boolean) => {
    setDialogOpen(open);
    if (!open) {
      setSelectedRow(undefined);
    }
  };

  const handleOpenCreateDialog = () => {
    setDialogVariant("create");
    setSelectedRow(undefined);
    setDialogOpen(true);
  };

  const handleViewModeChange = (_mode: "table" | "grid") => {
    toggleView();
  };

  return (
    <PureDataTable<TEntity>
      actionBar={actionBar}
      addCondition={addCondition}
      columns={columns}
      createForm={createForm as ComponentType}
      dialogOpen={dialogOpen}
      dialogVariant={dialogVariant}
      filter={filter as FilterNode}
      gridHoverRenderer={gridHover as (item: TEntity) => React.ReactNode}
      gridRenderer={grid as (item: TEntity) => React.ReactNode}
      icon={icon}
      minimal={minimal}
      onDialogChange={handleDialogChange}
      onOpenCreateDialog={handleOpenCreateDialog}
      onViewModeChange={handleViewModeChange}
      removeCondition={removeCondition}
      resetFilters={resetFilters}
      selectedRow={selectedRow as TEntity}
      setCondition={setCondition}
      showFilters={true}
      showPagination={true}
      showViewOptions={true}
      table={table}
      updateForm={
        updateForm as ComponentType<{
          id: string;
        }>
      }
      viewMode={view}
    />
  );
}
