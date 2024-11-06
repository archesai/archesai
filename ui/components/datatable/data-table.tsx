"use client";
import { DataTablePagination } from "@/components/datatable/data-table-pagination";
import { DataTableToolbar } from "@/components/datatable/data-table-toolbar";
import { DeleteItems } from "@/components/datatable/delete-items";
import { Button } from "@/components/ui/button";
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
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { useToggleView } from "@/hooks/useToggleView";
import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import * as VisuallyHidden from "@radix-ui/react-visually-hidden";
import {
  ColumnDef,
  ColumnFiltersState,
  getCoreRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from "@tanstack/react-table";
import { endOfDay } from "date-fns";
import { useEffect, useState } from "react";

import { GridView } from "./grid-view";
import { TableView } from "./table-view";

export interface BaseItem {
  id: string;
  name?: string;
}

interface DataTableProps<
  TItem extends BaseItem,
  TFindAllPathParams,
  TDeleteVariables,
> {
  columns: ColumnDef<TItem, TDeleteVariables>[];
  content?: (item: TItem) => JSX.Element;
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
  dataIcon: DataIcon,
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
      ...(range?.to
        ? { endDate: range?.to && endOfDay(range.to).toISOString() }
        : {}),
      ...(range?.to
        ? { startDate: range?.from && range.from.toISOString() }
        : {}),
      limit,
      name: query,
      offset: page * limit,
      sortBy: sortBy as "createdAt",
      sortDirection: sortDirection,

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
                        name: row.original.name || row.original.id,
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

  return (
    <div className="flex h-full flex-col gap-3">
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

      <div className="flex-1 overflow-auto">
        {view === "grid" ? (
          <GridView
            content={content}
            createForm={createForm}
            data={data?.results || []}
            DataIcon={DataIcon}
            deleteItem={deleteItem}
            getDeleteVariablesFromItem={getDeleteVariablesFromItem}
            getEditFormFromItem={getEditFormFromItem}
            handleSelect={handleSelect}
            hoverContent={hoverContent}
            itemType={itemType}
            selectedItems={selectedItems}
            setFinalForm={setFinalForm}
            setFormOpen={setFormOpen}
            toggleSelection={toggleSelection}
          />
        ) : (
          <TableView columns={columns} itemType={itemType} table={table} />
        )}
      </div>

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
