"use client";
import { DataTablePagination } from "@/components/datatable/data-table-pagination";
import { DataTableRowActions } from "@/components/datatable/data-table-row-actions";
import { DataTableToolbar } from "@/components/datatable/data-table-toolbar";
import { DeleteItems } from "@/components/datatable/delete-items";
import { Card } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useSelectItems } from "@/hooks/useSelectItems";
import { useToggleView } from "@/hooks/useToggleView";
import {
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from "@tanstack/react-table";
import { PlusSquare } from "lucide-react";
import { useEffect, useState } from "react";

export interface BaseItem {
  id: string;
  name: string;
}

interface DataTableProps<TItem extends BaseItem, TMutationVariables> {
  columns: ColumnDef<TItem, TMutationVariables>[];
  content: (item: TItem) => JSX.Element;
  data: {
    metadata: {
      limit: number;
      offset: number;
      totalResults: number;
    };
    results: TItem[];
  };

  dataIcon: JSX.Element;
  defaultView?: "grid" | "table";
  deleteItem: (vars: TMutationVariables) => Promise<void>;

  getDeleteVariablesFromItem: (item: TItem) => TMutationVariables[];
  handleSelect: (item: TItem) => void;

  hidePagination?: boolean;
  hideSearch?: boolean;

  hoverContent?: (item: TItem) => JSX.Element;
  itemType: string;
  loading: boolean;

  mutationVariables: TMutationVariables[];
}

export function DataTable<TItem extends BaseItem, TMutationVariables>({
  columns,
  content,
  data,
  dataIcon,
  defaultView,
  deleteItem,
  getDeleteVariablesFromItem,
  handleSelect,
  hoverContent,
  itemType,
}: DataTableProps<TItem, TMutationVariables>) {
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [sorting, setSorting] = useState<SortingState>([]);

  const { setView, view } = useToggleView();

  useEffect(() => {
    if (defaultView) {
      setView(defaultView);
    }
  }, [defaultView]);

  const {
    selectedAllItems,
    selectedItems,
    selectedSomeItems,
    setSelectedItems,
    toggleSelectAll,
    toggleSelection,
  } = useSelectItems({ items: data?.results || [] });

  const table = useReactTable({
    columns: [
      {
        cell: ({ row }) => (
          <Checkbox
            aria-label="Select row"
            checked={selectedItems.includes(row.original.id)}
            className="translate-y-[2px]"
            onCheckedChange={() => toggleSelection(row.original.id)}
          />
        ),
        enableHiding: false,
        enableSorting: false,
        header: () => (
          <Checkbox
            aria-label="Select all"
            checked={selectedAllItems || (selectedSomeItems && "indeterminate")}
            className="translate-y-[2px]"
            onCheckedChange={() => toggleSelectAll()}
          />
        ),
        id: "select",
      },
      ...columns,
      {
        cell: ({ row }) => <DataTableRowActions row={row} />,
        id: "actions",
      },
    ],
    data: data?.results || [],
    enableRowSelection: true,
    getCoreRowModel: getCoreRowModel(),

    manualFiltering: true,
    manualPagination: true,
    manualSorting: true,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onSortingChange: setSorting,
    state: {
      columnFilters,
      columnVisibility,
      sorting,
    },
  });
  const [hover, setHover] = useState(-1);

  const grid_view = (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 w-full">
      {data?.results.map((item, i) => {
        const isItemSelected = selectedItems.includes(item.id);
        return (
          <Card
            className={`shadow-sm relative w-full ${isItemSelected ? "ring-4 ring-blue-500" : ""} overflow-visible after:content-[''] after:absolute after:w-full after:h-full after:top-0 after:left-0 after:border-radius-inherit after:z-10 after:transition-shadow after:pointer-events-none`}
            key={i}
          >
            <div
              className="h-64 cursor-pointer relative overflow-hidden group hover:bg-primary-foreground transition-all"
              onClick={async () => handleSelect(item)}
              onMouseEnter={() => setHover(i)}
              onMouseLeave={() => setHover(-1)}
            >
              {content(item)}
            </div>
            <hr />

            <div className="flex justify-between items-center mt-auto p-2">
              <div className="flex items-center">
                {dataIcon}
                <h5 className="text-base leading-tight overflow-hidden text-ellipsis whitespace-nowrap pl-1 max-w-[8.5rem]">
                  {item.name}
                </h5>
              </div>
              <div className="flex items-center gap-1 p-2">
                <PlusSquare
                  className="text-primary cursor-pointer"
                  onClick={async () => handleSelect(item)}
                />
                <DeleteItems
                  items={[
                    {
                      id: item.id,
                      name: item.name,
                    },
                  ]}
                  itemType={itemType}
                  mutationFunction={async (vars) => {
                    await deleteItem(vars);
                    setSelectedItems([]);
                  }}
                  mutationVariables={getDeleteVariablesFromItem(item)}
                />
                <Checkbox
                  aria-label={`Select ${item.name}`}
                  checked={isItemSelected}
                  className="h-5 w-5 m-1 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  onCheckedChange={() => toggleSelection(item.id)}
                />
              </div>
            </div>
            {hoverContent && hover === i && hoverContent(item)}
          </Card>
        );
      })}
    </div>
  );

  const table_view = (
    <div className="rounded-md border bg-background shadow-sm">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead colSpan={header.colSpan} key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                data-state={row.getIsSelected() && "selected"}
                key={row.id}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                className="h-24 text-center"
                colSpan={columns.length + 2}
              >
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );

  return (
    <div className="flex space-y-4 flex-col justify-between backdrop-blur-sm opacity-70 h-full">
      <div className="space-y-4">
        <DataTableToolbar itemType={itemType} table={table} />
        {data?.results?.length > 10 && <DataTablePagination data={data} />}
        {view === "grid" ? grid_view : table_view}
      </div>
      <div>
        <DataTablePagination data={data} />
      </div>
    </div>
  );
}
