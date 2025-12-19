import type {
  AccessorKeyColumnDef,
  ColumnFiltersState,
  ColumnSort,
  PaginationState,
  RowSelectionState,
  SortingState,
  TableOptions,
  TableState,
  Updater,
  VisibilityState,
} from "@tanstack/react-table";
import { getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useCallback, useMemo, useState } from "react";
import { DataTableColumnHeader } from "../components/datatable/components/data-table-column-header";
import { Checkbox } from "../components/shadcn/checkbox";
import { toSentenceCase } from "../lib/utils";
import type { BaseEntity } from "../types/entities";
import { useFilterState } from "./use-filter-state";

const DEBOUNCE_MS = 300;
const THROTTLE_MS = 50;

export interface ExtendedColumnSort<TData extends BaseEntity>
  extends Omit<ColumnSort, "id"> {
  id: Extract<keyof TData, string>;
}

interface useDataTableProps<TEntity extends BaseEntity>
  extends Omit<
      TableOptions<TEntity>,
      | "getCoreRowModel"
      | "manualFiltering"
      | "manualPagination"
      | "manualSorting"
      | "pageCount"
      | "state"
    >,
    Required<Pick<TableOptions<TEntity>, "pageCount">> {
  clearOnDefault?: boolean;
  debounceMs?: number;
  enableAdvancedFilter?: boolean;
  history?: "push" | "replace";
  initialState?: Omit<Partial<TableState>, "sorting"> & {
    sorting?: ExtendedColumnSort<TEntity>[];
  };
  scroll?: boolean;
  shallow?: boolean;
  startTransition?: React.TransitionStartFunction;
  throttleMs?: number;
}

export function useDataTable<TData extends BaseEntity>(
  props: useDataTableProps<TData>,
): {
  debounceMs: number;
  shallow: boolean;
  table: ReturnType<typeof useReactTable<TData>>;
  throttleMs: number;
} {
  const {
    // clearOnDefault = false,
    columns,
    debounceMs = DEBOUNCE_MS,
    enableAdvancedFilter = false,
    // history = 'replace',
    initialState,
    pageCount = -1,
    // scroll = false,
    shallow = true,
    // startTransition,
    throttleMs = THROTTLE_MS,
    ...tableProps
  } = props;

  const [rowSelection, setRowSelection] = useState<RowSelectionState>(
    initialState?.rowSelection ?? {},
  );
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>(
    initialState?.columnVisibility ?? {},
  );

  const {
    // filter,
    pageNumber,
    pageSize,
    // setFilter,
    setPage,
    setPageSize,
    setSorting,
    sorting,
  } = useFilterState<TData>();

  const pagination: PaginationState = useMemo(() => {
    return {
      pageIndex: pageNumber - 1, // zero-based index -> one-based index
      pageSize: pageSize,
    };
  }, [pageNumber, pageSize]);

  const onPaginationChange = useCallback(
    (updaterOrValue: Updater<PaginationState>) => {
      if (typeof updaterOrValue === "function") {
        const newPagination = updaterOrValue(pagination);
        setPage(newPagination.pageIndex + 1);
        setPageSize(newPagination.pageSize);
      } else {
        setPage(updaterOrValue.pageIndex + 1);
        setPageSize(updaterOrValue.pageSize);
      }
    },
    [pagination, setPage, setPageSize],
  );

  const onSortingChange = useCallback(
    (updaterOrValue: Updater<SortingState>) => {
      if (typeof updaterOrValue === "function") {
        const newSorting = updaterOrValue(sorting);
        setSorting(newSorting as ExtendedColumnSort<TData>[]);
      } else {
        setSorting(updaterOrValue as ExtendedColumnSort<TData>[]);
      }
    },
    [sorting, setSorting],
  );

  const initialColumnFilters: ColumnFiltersState = useMemo(() => {
    if (enableAdvancedFilter) return [];
    return [];
  }, [enableAdvancedFilter]);

  const [columnFilters, setColumnFilters] =
    useState<ColumnFiltersState>(initialColumnFilters);

  const onColumnFiltersChange = useCallback(
    (updaterOrValue: Updater<ColumnFiltersState>) => {
      if (enableAdvancedFilter) return;

      setColumnFilters((prev) => {
        const next =
          typeof updaterOrValue === "function"
            ? updaterOrValue(prev)
            : updaterOrValue;

        return next;
      });
    },
    [enableAdvancedFilter],
  );

  const enhancedColumns = useMemo<AccessorKeyColumnDef<TData>[]>(
    () => [
      // Checkbox column
      {
        accessorKey: "select",
        cell: ({ row }) => (
          <div className="flex w-4">
            <Checkbox
              aria-label="Select row"
              checked={row.getIsSelected()}
              onCheckedChange={(value) => {
                row.toggleSelected(!!value);
              }}
            />
          </div>
        ),
        enableHiding: false,
        enableSorting: false,
        header: ({ table }) => (
          <div className="flex">
            <Checkbox
              aria-label="Select all"
              checked={
                table.getIsAllPageRowsSelected() ||
                (table.getIsSomePageRowsSelected() && "indeterminate")
              }
              className="translate-y-0.5"
              onCheckedChange={(value) => {
                table.toggleAllPageRowsSelected(!!value);
              }}
            />
          </div>
        ),
        id: "select",
      },
      ...columns.map((column) => ({
        ...column,
        accessorKey: (column as AccessorKeyColumnDef<TData>).accessorKey,
        header:
          column.header ??
          (({ column: col }) => (
            <DataTableColumnHeader
              column={col}
              title={toSentenceCase(column.meta?.label?.toString() ?? "")}
            />
          )),
      })),
    ],
    [columns],
  );

  // Create table with minimal state - let filterState handle pagination/sorting
  const table = useReactTable({
    ...tableProps,
    columns: enhancedColumns,
    defaultColumn: {
      ...tableProps.defaultColumn,
      enableColumnFilter: false,
    },
    enableRowSelection: true,
    getCoreRowModel: getCoreRowModel(),
    // initialState,
    manualFiltering: true,
    manualPagination: true,
    manualSorting: true,
    onColumnFiltersChange,
    onColumnVisibilityChange: setColumnVisibility,
    onPaginationChange,
    onRowSelectionChange: setRowSelection,
    onSortingChange,
    pageCount,
    state: {
      columnFilters,
      columnVisibility,
      pagination,
      rowSelection,
      sorting,
    },
  });

  return {
    debounceMs,
    shallow,
    table,
    throttleMs,
  };
}
