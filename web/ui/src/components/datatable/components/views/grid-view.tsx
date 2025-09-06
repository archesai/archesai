"use no memo";

import type { Table } from "@tanstack/react-table";
import type { JSX } from "react";

import { useState } from "react";
import { flexRender } from "@tanstack/react-table";

import type { BaseEntity } from "#types/entities";

import { Card, CardContent, CardFooter } from "#components/shadcn/card";
import { cn } from "#lib/utils";

export interface GridViewProps<TEntity extends BaseEntity> {
  grid?: (item: TEntity) => React.ReactNode;
  gridHover?: (item: TEntity) => React.ReactNode;
  icon: React.ReactNode;
  table: Table<TEntity>;
}

export function GridView<TEntity extends BaseEntity>({
  grid,
  gridHover,
  icon,
  table,
}: GridViewProps<TEntity>): JSX.Element {
  const data = table.getRowModel().rows;
  const [hover, setHover] = useState(-1);

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
      {/* Data Cards */}
      {data.length > 0 ? (
        data.map((item, i) => {
          const isItemSelected = item.getIsSelected();
          const checkbox = item.getAllCells().at(0)?.column.columnDef.cell;
          const context = item.getAllCells().at(0)?.getContext();

          return (
            <Card
              className={cn(
                `h-64 transition-all duration-200 hover:bg-accent hover:shadow-lg`,
                isItemSelected && "bg-accent/80",
              )}
              key={item.id}
            >
              {/* Top Content */}
              <CardContent
                className="h-full cursor-pointer"
                onClick={item.getToggleSelectedHandler()}
                onMouseEnter={() => {
                  setHover(i);
                }}
                onMouseLeave={() => {
                  setHover(-1);
                }}
              >
                {grid ? (
                  grid(item.original)
                ) : (
                  <div className="flex h-full items-center justify-center">
                    {icon}
                  </div>
                )}
              </CardContent>

              <hr />

              {/* Footer */}
              <CardFooter className="justify-start">
                {context && flexRender(checkbox, context)}
                <span className="truncate">{item.original.id}</span>
              </CardFooter>
              {gridHover && hover === i && gridHover(item.original)}
            </Card>
          );
        })
      ) : (
        <div className="col-span-4 row-span-4 flex items-center justify-center pt-20 text-sm">
          No items found
        </div>
      )}
    </div>
  );
}
