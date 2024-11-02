"use client";
import { DataTablePagination } from "@/components/datatable/data-table-pagination";
import { DataTableToolbar } from "@/components/datatable/data-table-toolbar";
import { DeleteItems } from "@/components/datatable/delete-items";
import { Button } from "@/components/ui/button";
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
  minimal?: boolean;
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
  minimal,
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
                  <DotsHorizontalIcon className="h-5 w-5" />
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
    <div className="grid w-full grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-4">
      {createForm ? (
        <Card
          className={`flex h-48 cursor-pointer flex-col items-center justify-center border-2 border-dashed border-gray-400 shadow-sm transition-all hover:bg-secondary`}
          onClick={async () => {
            setFormOpen(true);
          }}
        >
          <PlusSquare className="w-h-6 h-6 text-primary" />
          <span className="text-md">Create {itemType}</span>
        </Card>
      ) : null}
      {data?.results.map((item, i) => {
        const isItemSelected = selectedItems.includes(item.id);
        return (
          <Card
            className={`relative w-full shadow-sm ${isItemSelected ? "ring-4 ring-blue-500" : ""} after:border-radius-inherit overflow-visible after:pointer-events-none after:absolute after:left-0 after:top-0 after:z-10 after:h-full after:w-full after:transition-shadow after:content-['']`}
            key={i}
          >
            <div
              className="group relative h-40 cursor-pointer overflow-hidden rounded-t-sm transition-all hover:bg-gray-200 hover:dark:bg-gray-900"
              onClick={async () => handleSelect(item)}
              onMouseEnter={() => setHover(i)}
              onMouseLeave={() => setHover(-1)}
            >
              {content(item)}
            </div>
            <hr />

            <div className="mt-auto flex items-center justify-between p-2">
              <div className="flex min-w-0 items-center gap-2">
                <Checkbox
                  aria-label={`Select ${item.name}`}
                  checked={isItemSelected}
                  className="rounded text-blue-600 focus:ring-blue-500"
                  onCheckedChange={() => toggleSelection(item.id)}
                />
                <h5 className="overflow-hidden text-ellipsis whitespace-nowrap text-base leading-tight">
                  {item.name}
                </h5>
              </div>
              <div className="flex flex-shrink-0 items-center gap-2">
                {getEditFormFromItem ? (
                  <FilePenLine
                    className="h-5 w-5 cursor-pointer text-primary"
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
                  <TableHead
                    className="text-base"
                    colSpan={header.colSpan}
                    key={header.id}
                  >
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
                className="transition-all hover:bg-slate-200 hover:dark:bg-slate-900"
                data-state={row.getIsSelected() && "selected"}
                key={row.id}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell className="p-2" key={cell.id}>
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
                No {itemType}s found
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );

  return (
    <div className="flex h-full flex-grow flex-col gap-3">
      {!minimal && (
        <DataTableToolbar
          data={data?.results || []}
          itemType={itemType}
          table={table}
        />
      )}
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

      <div className="grow"> {view === "grid" ? grid_view : table_view}</div>

      {!minimal && (
        <div className="self-auto">
          <DataTablePagination data={data as any} />
        </div>
      )}

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
