"use client";
import { DataTablePagination } from "@/components/datatable/data-table-pagination";
import { DataTableToolbar } from "@/components/datatable/data-table-toolbar";
import { DeleteItems } from "@/components/datatable/delete-items";
import { Card } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { useToggleView } from "@/hooks/useToggleView";
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
import { endOfDay } from "date-fns";
import { FilePenLine, PlusSquare } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "../ui/button";

export interface BaseItem {
  id: string;
  name: string;
}

interface DataTableProps<
  TItem extends BaseItem,
  TFindAllPathParams,
  TDeleteVariables,
> {
  columns: ColumnDef<TItem, TDeleteVariables>[];
  content: (item: TItem) => JSX.Element;
  createForm?: React.ReactNode;

  dataIcon: JSX.Element;
  defaultView?: "grid" | "table";

  findAllPathParams: TFindAllPathParams;
  findAllQueryParams?: object;
  getDeleteVariablesFromItem: (item: TItem) => TDeleteVariables;
  getEditFormFromItem?: (item: TItem) => React.ReactNode;
  handleSelect: (item: TItem) => void;

  hoverContent?: (item: TItem) => JSX.Element;
  itemType: string;
  useFindAll: (s: any) => {
    data:
      | {
          metadata: {
            limit: number;
            offset: number;
            totalResults: number;
          };
          results: TItem[];
        }
      | undefined;
    isLoading: boolean;
    isPlaceholderData: boolean;
  };
  useRemove: () => {
    mutateAsync: (vars: TDeleteVariables) => Promise<void>;
  };
}

export function DataTable<
  TItem extends BaseItem,
  TFindAllPathParams,
  TDeleteVariables,
>({
  columns,
  content,
  createForm,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  dataIcon,
  defaultView,
  findAllPathParams,
  findAllQueryParams,
  getDeleteVariablesFromItem,
  getEditFormFromItem,
  handleSelect,
  hoverContent,
  itemType,
  useFindAll,
  useRemove,
}: DataTableProps<TItem, TFindAllPathParams, TDeleteVariables>) {
  const { limit, page, query, range, sortBy, sortDirection } = useFilterItems();
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [sorting, setSorting] = useState<SortingState>([
    {
      desc: true,
      id: "createdAt",
    },
  ]);
  const [formOpen, setFormOpen] = useState(false);
  const [finalForm, setFinalForm] = useState<React.ReactNode | undefined>(
    createForm
  );
  const { setSortBy, setSortDirection } = useFilterItems();

  useEffect(() => {
    setSortDirection(sorting[0]?.desc ? "desc" : "asc");
    setSortBy(sorting[0]?.id);
  }, [sorting]);

  const { setView, view } = useToggleView();

  useEffect(() => {
    if (defaultView) {
      setView(defaultView);
    }
  }, [defaultView]);

  const { data } = useFindAll({
    pathParams: findAllPathParams,
    queryParams: {
      endDate: endOfDay(range.to).toISOString(),
      limit,
      name: query,
      offset: page * limit,
      sortBy: sortBy as "createdAt",
      sortDirection: sortDirection,
      startDate: range.from.toISOString(),
      ...findAllQueryParams,
    },
  });
  const { mutateAsync: deleteItem } = useRemove();

  const { selectedItems, setSelectedItems, toggleSelection } = useSelectItems({
    items: data?.results || [],
  });

  const table = useReactTable({
    columns: [
      {
        cell: ({ row }) => (
          <Checkbox
            aria-label="Select row"
            checked={selectedItems.includes(row.original.id)}
            className=""
            onCheckedChange={() => toggleSelection(row.original.id)}
          />
        ),
        enableHiding: false,
        enableSorting: false,

        id: "select",
      },
      ...columns,
      {
        cell: ({ row }) => (
          <div className="flex justify-end">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
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
                    deleteFunction={async (vars) => {
                      await deleteItem(vars);
                      setSelectedItems([]);
                    }}
                    deleteVariables={[getDeleteVariablesFromItem(row.original)]}
                    items={[
                      {
                        id: row.original.id,
                        name: row.original.name,
                      },
                    ]}
                    itemType={itemType}
                    variant="md"
                  />
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        ),
        header: () =>
          createForm ? (
            <div className="text-right">
              <Button
                onClick={() => setFormOpen(true)}
                size="sm"
                variant={"secondary"}
              >
                Create {itemType}
              </Button>
            </div>
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
          className={`m-30 bg-transparent shadow-sm relative w-full cursor-pointer overflow-visible transition-all hover:bg-gray-200 hover:dark:bg-gray-900 border-2 border-dashed border-gray-400 after:content-[''] after:absolute after:w-full after:h-full after:top-0 after:left-0 after:border-radius-inherit after:z-10 after:transition-shadow after:pointer-events-none flex items-center justify-center`}
          onClick={async () => {
            setFormOpen(true);
          }}
        >
          <div className="h-48 relative overflow-hidden group  flex flex-col items-center justify-center">
            <PlusSquare className="h-6 w-h-6 text-primary" />
            <span className="mt-2 text-md">Create {itemType}</span>
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
              className="h-48 rounded-t-sm cursor-pointer relative overflow-hidden group hover:bg-gray-200 hover:dark:bg-gray-900 transition-all"
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
                  className=" text-blue-600 rounded focus:ring-blue-500"
                  onCheckedChange={() => toggleSelection(item.id)}
                />
                <h5 className="text-base leading-tight overflow-hidden text-ellipsis whitespace-nowrap pl-2">
                  {item.name}
                </h5>
              </div>
              <div className="flex items-center gap-2 flex-shrink-0">
                {getEditFormFromItem ? (
                  <FilePenLine
                    className=" text-primary cursor-pointer h-5 w-5"
                    onClick={() => {
                      setFinalForm(getEditFormFromItem?.(item));
                      setFormOpen(true);
                    }}
                  />
                ) : null}
                <DeleteItems
                  deleteFunction={async (vars) => {
                    await deleteItem(vars);
                    setSelectedItems([]);
                  }}
                  deleteVariables={[getDeleteVariablesFromItem(item)]}
                  items={[
                    {
                      id: item.id,
                      name: item.name,
                    },
                  ]}
                  itemType={itemType}
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
                  <TableCell className="py-2" key={cell.id}>
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
    <div className="flex space-y-4 flex-col justify-between backdrop-blur-sm h-full">
      <div className="space-y-4">
        <DataTableToolbar
          data={data?.results || []}
          itemType={itemType}
          table={table}
        />
        {selectedItems.length > 0 && (
          <DeleteItems
            deleteFunction={async (vars) => {
              await deleteItem(vars);
              setSelectedItems([]);
            }}
            deleteVariables={selectedItems.map((id) =>
              getDeleteVariablesFromItem(
                data?.results.find((i) => i.id === id) as TItem
              )
            )}
            items={selectedItems.map((id) => {
              const item = data?.results.find((i) => i.id === id);
              return {
                id: item?.id || "",
                name: item?.name || "",
              };
            })}
            itemType={itemType}
            variant="lg"
          />
        )}

        {view === "grid" ? grid_view : table_view}
      </div>
      <div>
        <DataTablePagination data={data as any} />
      </div>
      {/* THIS IS THE FORM DIALOG */}
      <Dialog
        onOpenChange={(o) => {
          setFormOpen(o);
          if (!o) {
            setFinalForm(createForm);
          }
        }}
        open={formOpen}
      >
        <VisuallyHidden.Root>
          <DialogDescription />
          <DialogTitle>
            {finalForm ? "Edit" : "Create"} {itemType}
          </DialogTitle>
        </VisuallyHidden.Root>
        <DialogContent
          aria-description="Create/Edit"
          className="p-0"
          title="Create/Edit"
        >
          {finalForm}
        </DialogContent>
      </Dialog>
    </div>
  );
}
