'use no memo'

import type { Table } from '@tanstack/react-table'

import { useState } from 'react'
import { flexRender } from '@tanstack/react-table'

import type { BaseEntity } from '@archesai/schemas'

import { Card } from '#components/shadcn/card'
import { cn } from '#lib/utils'

export interface GridViewProps<TEntity extends BaseEntity> {
  grid?: (item: TEntity) => React.ReactNode
  gridHover?: (item: TEntity) => React.ReactNode
  icon: React.ReactNode
  table: Table<TEntity>
}

export function GridView<TEntity extends BaseEntity>({
  grid,
  gridHover,
  icon,
  table
}: GridViewProps<TEntity>) {
  const data = table.getRowModel().rows
  const [hover, setHover] = useState(-1)

  return (
    <div className='grid w-full grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-4'>
      {/* Data Cards */}
      {data.length > 0 ?
        data.map((item, i) => {
          const isItemSelected = item.getIsSelected()
          return (
            <div
              className={cn(
                `relative flex h-64 flex-col overflow-hidden rounded-xl bg-card/100 shadow-xs transition-all`,
                isItemSelected && 'bg-secondary/50'
              )}
              key={item.id}
            >
              {/* Top Content */}
              <div
                className='h-full cursor-pointer transition-all hover:bg-secondary'
                onClick={item.getToggleSelectedHandler()}
                onMouseEnter={() => {
                  setHover(i)
                }}
                onMouseLeave={() => {
                  setHover(-1)
                }}
              >
                {grid ?
                  grid(item.original)
                : <div className='flex h-full items-center justify-center'>
                    {icon}
                  </div>
                }
              </div>
              <hr />

              {/* Footer */}
              <div className='flex items-center justify-between bg-card p-4'>
                <div className='flex min-w-0 items-center gap-2'>
                  {flexRender(
                    item.getAllCells().at(0)?.column.columnDef.cell,
                    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion, @typescript-eslint/no-non-null-asserted-optional-chain
                    item.getAllCells().at(0)?.getContext()!
                  )}
                  {item.original.id}
                </div>
                {flexRender(
                  item.getAllCells().at(-1)?.column.columnDef.cell,
                  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion, @typescript-eslint/no-non-null-asserted-optional-chain
                  item.getAllCells().at(-1)?.getContext()!
                )}
              </div>
              {gridHover && hover === i && gridHover(item.original)}
            </div>
          )
        })
      : <div className='col-span-4 row-span-4 flex items-center justify-center pt-20 text-sm'>
          No items found
        </div>
      }
    </div>
  )
}
