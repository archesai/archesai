"use client";
import { DataTablePagination } from "@/components/datatable/data-table-pagination";
import { DataTableToolbar } from "@/components/datatable/data-table-toolbar";
import { DeleteItems } from "@/components/datatable/delete-items";
import { Card } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
} from "@/components/ui/drawer";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
import { DialogTitle } from "@radix-ui/react-dialog";
import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import * as VisuallyHidden from "@radix-ui/react-visually-hidden";
import {
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from "@tanstack/react-table";
import { FilePenLine, PlusSquare } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "../ui/button";
import { DropdownMenuShortcut } from "../ui/dropdown-menu";

export interface BaseItem {
  id: string;
  name: string;
}

interface DataTableProps<TItem extends BaseItem, TMutationVariables> {
  columns: ColumnDef<TItem, TMutationVariables>[];
  content: (item: TItem) => JSX.Element;
  createForm?: React.ReactNode;

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
  getEditFormFromItem?: (item: TItem) => React.ReactNode;
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
  createForm,
  data,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  dataIcon,
  defaultView,
  deleteItem,
  getDeleteVariablesFromItem,
  getEditFormFromItem,
  handleSelect,
  hoverContent,
  itemType,
}: DataTableProps<TItem, TMutationVariables>) {
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [sorting, setSorting] = useState<SortingState>([]);
  const [formOpen, setFormOpen] = useState(false);
  const [finalForm, setFinalForm] = useState<React.ReactNode | undefined>(
    createForm
  );

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
        cell: ({ row }) => (
          <DropdownMenu>
            <DropdownMenuTrigger asChild className="text-center">
              <Button
                className="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
                variant="ghost"
              >
                <DotsHorizontalIcon className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[160px]">
              {getEditFormFromItem ? (
                <>
                  <DropdownMenuItem
                    onClick={() => {
                      setFinalForm(getEditFormFromItem?.(row.original));
                      setFormOpen(true);
                    }}
                  >
                    Edit
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                </>
              ) : null}
              <DropdownMenuItem
                onSelect={(e) => e.preventDefault()} // Prevent closing on select
              >
                <DeleteItems
                  items={[
                    {
                      id: row.original.id,
                      name: row.original.name,
                    },
                  ]}
                  itemType={itemType}
                  mutationFunction={async (vars) => {
                    await deleteItem(vars);
                    setSelectedItems([]);
                  }}
                  mutationVariables={getDeleteVariablesFromItem(row.original)}
                  variant="md"
                />
                <DropdownMenuShortcut>⌘⌫</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ),
        header: () =>
          createForm ? (
            <Button onClick={() => setFormOpen(true)} size="sm">
              New {itemType}
            </Button>
          ) : null,
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
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 w-full">
      {createForm ? (
        <Card
          className={`shadow-sm relative w-full cursor-pointer overflow-visible hover:bg-gray-200 hover:dark:bg-gray-900 border-2 border-dashed border-gray-400 after:content-[''] after:absolute after:w-full after:h-full after:top-0 after:left-0 after:border-radius-inherit after:z-10 after:transition-shadow after:pointer-events-none flex items-center justify-center`}
          onClick={async () => {
            setFormOpen(true);
          }}
        >
          <div className="h-48 relative overflow-hidden group transition-all flex flex-col items-center justify-center">
            <PlusSquare size={30} />
            <span className="mt-2 text-lg">New {itemType}</span>
          </div>
        </Card>
      ) : null}
      {data?.results.map((item, i) => {
        const isItemSelected = selectedItems.includes(item.id);
        return (
          <Card
            className={`shadow-sm relative w-full ${isItemSelected ? "ring-4 ring-blue-500" : ""} overflow-visible after:content-[''] after:absolute after:w-full after:h-full after:top-0 after:left-0 after:border-radius-inherit after:z-10 after:transition-shadow after:pointer-events-none`}
            key={i}
          >
            <div
              className="h-48 cursor-pointer relative overflow-hidden group hover:bg-gray-200 hover:dark:bg-gray-900 transition-all"
              onClick={async () => handleSelect(item)}
              onMouseEnter={() => setHover(i)}
              onMouseLeave={() => setHover(-1)}
            >
              {content(item)}
            </div>
            <hr />

            <div className="flex justify-between items-center mt-auto p-2">
              <div className="flex items-center min-w-0">
                <Checkbox
                  aria-label={`Select ${item.name}`}
                  checked={isItemSelected}
                  className="h-5 w-5 text-blue-600 rounded focus:ring-blue-500"
                  onCheckedChange={() => toggleSelection(item.id)}
                />
                <h5 className="text-base leading-tight overflow-hidden text-ellipsis whitespace-nowrap pl-2">
                  {item.name}
                </h5>
              </div>
              <div className="flex items-center gap-2 p-2 flex-shrink-0">
                {getEditFormFromItem ? (
                  <FilePenLine
                    className="text-primary cursor-pointer"
                    onClick={() => {
                      setFinalForm(getEditFormFromItem?.(item));
                      setFormOpen(true);
                    }}
                  />
                ) : null}
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
      <Drawer
        onOpenChange={(o) => {
          setFormOpen(o);
          if (!o) {
            setFinalForm(createForm);
          }
        }}
        open={formOpen}
      >
        <VisuallyHidden.Root>
          <DrawerDescription />
          <DialogTitle>
            {finalForm ? "Edit" : "Create"} {itemType}
          </DialogTitle>
        </VisuallyHidden.Root>
        <DrawerContent
          aria-description="Create/Edit"
          className="p-3 bg-transparent border-none shadow-none"
          title="Create/Edit"
        >
          {finalForm}
        </DrawerContent>
      </Drawer>
    </div>
  );
}
