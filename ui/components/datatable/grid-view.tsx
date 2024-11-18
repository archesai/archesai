"use client";

import { DeleteItems } from "@/components/datatable/delete-items";
import { Card } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { FilePenLine } from "lucide-react";
import { useState } from "react";

import { BaseItem } from "./data-table";

interface GridViewProps<TItem extends BaseItem> {
  content?: (item: TItem) => JSX.Element;
  createForm?: React.ReactNode;
  data: TItem[];
  DataIcon: JSX.Element;
  deleteItem: (vars: any) => Promise<void>;
  getDeleteVariablesFromItem: (item: TItem) => any;
  getEditFormFromItem?: (item: TItem) => React.ReactNode;
  handleSelect: (item: TItem) => void;
  hoverContent?: (item: TItem) => JSX.Element;
  itemType: string;
  selectedItems: string[];
  setFinalForm: (form: React.ReactNode | undefined) => void;
  setFormOpen: (open: boolean) => void;
  toggleSelection: (id: string) => void;
}

export function GridView<TItem extends BaseItem>({
  content,
  data,
  DataIcon,
  deleteItem,
  getDeleteVariablesFromItem,
  getEditFormFromItem,
  handleSelect,
  hoverContent,
  itemType,
  selectedItems,
  setFinalForm,
  setFormOpen,
  toggleSelection,
}: GridViewProps<TItem>) {
  const [hover, setHover] = useState(-1);

  return (
    <div className="grid w-full grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-4">
      {/* Data Cards */}
      {data.map((item, i) => {
        const isItemSelected = selectedItems.includes(item.id);
        return (
          <Card
            className={`relative flex aspect-auto h-64 flex-col shadow-sm transition-all hover:bg-muted ${
              isItemSelected ? "ring-4 ring-blue-500" : ""
            } after:border-radius-inherit overflow-visible after:pointer-events-none after:absolute after:left-0 after:top-0 after:z-10 after:h-full after:w-full after:transition-shadow after:content-['']`}
            key={item.id}
          >
            {/* Top Content */}
            <div
              className="group relative grow cursor-pointer overflow-auto rounded-t-sm transition-all"
              onClick={() => handleSelect(item)}
              onMouseEnter={() => setHover(i)}
              onMouseLeave={() => setHover(-1)}
            >
              {content ? (
                content(item)
              ) : (
                <div className="flex h-full w-full items-center justify-center">
                  {DataIcon}
                </div>
              )}
            </div>
            <hr />

            {/* Footer */}
            <div className="mt-auto flex items-center justify-between p-2">
              <div className="flex min-w-0 items-center gap-2">
                <Checkbox
                  aria-label={`Select ${item.name}`}
                  checked={isItemSelected}
                  className="rounded text-blue-600 focus:ring-blue-500"
                  onCheckedChange={() => toggleSelection(item.id)}
                />
                <span className="overflow-hidden text-ellipsis whitespace-nowrap text-base leading-tight">
                  {item.name}
                </span>
              </div>
              <div className="flex flex-shrink-0 items-center gap-2">
                {getEditFormFromItem && (
                  <FilePenLine
                    className="h-5 w-5 cursor-pointer text-primary"
                    onClick={() => {
                      setFinalForm(getEditFormFromItem(item));
                      setFormOpen(true);
                    }}
                  />
                )}
                <DeleteItems
                  deleteFunction={async (vars) => {
                    await deleteItem(vars);
                  }}
                  deleteVariables={[getDeleteVariablesFromItem(item)]}
                  items={[
                    {
                      id: item.id,
                      name: item.name || item.id,
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
}
