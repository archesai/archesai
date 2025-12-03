/// <reference path="../../../types/table-meta.d.ts" />
import type { Table } from "@tanstack/react-table";
import type { JSX } from "react";
import { DeleteItems } from "#components/custom/delete-items";
import { MoreHorizontalIcon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "#components/shadcn/dropdown-menu";
import type { BaseEntity } from "#types/entities";

export interface DataTableRowActionsProps<TEntity extends BaseEntity> {
  deleteItem?: (id: string) => Promise<void>;
  getEditFormFromItem?: (item: TEntity) => React.ReactNode;
  row: {
    original: TEntity;
  };
  setFinalForm?: (form: React.ReactNode) => void;
  setFormOpen?: (open: boolean) => void;
  table: Table<TEntity>;
}

export function DataTableRowActions<TEntity extends BaseEntity>(
  props: DataTableRowActionsProps<TEntity>,
): JSX.Element {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label="Expand row options"
          className="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
          variant="ghost"
        >
          <MoreHorizontalIcon className="h-5 w-5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="end"
        className="w-[160px]"
      >
        {props.getEditFormFromItem ? (
          <>
            <DropdownMenuItem
              onClick={() => {
                if (!props.getEditFormFromItem) {
                  throw new Error(
                    "getEditFormFromItem function is not defined",
                  );
                }
                if (props.setFinalForm && props.setFormOpen) {
                  props.setFinalForm(
                    props.getEditFormFromItem(props.row.original),
                  );
                  props.setFormOpen(true);
                }
              }}
            >
              Edit
            </DropdownMenuItem>
            <DropdownMenuSeparator />
          </>
        ) : null}
        {props.deleteItem ? (
          <DropdownMenuItem
            onSelect={(e) => {
              e.preventDefault();
            }}
          >
            <DeleteItems
              deleteItem={props.deleteItem}
              entityKey={props.table.options.meta?.entityKey ?? "Entity"}
              items={[props.row.original]}
            />
          </DropdownMenuItem>
        ) : null}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
