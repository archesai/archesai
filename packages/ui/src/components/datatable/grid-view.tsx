import { useState } from 'react'
import { FilePenLine } from 'lucide-react'

import type { BaseEntity } from '@archesai/domain'

import type { DataTableContainerProps } from '#components/datatable/data-table'

import { DeleteItems } from '#components/datatable/delete-items'
import { Card } from '#components/shadcn/card'
import { Checkbox } from '#components/shadcn/checkbox'

export type GridViewProps<TEntity extends BaseEntity> = Pick<
  DataTableContainerProps<TEntity>,
  | 'createForm'
  | 'data'
  | 'defaultView'
  | 'deleteItem'
  | 'entityType'
  | 'getEditFormFromItem'
  | 'grid'
  | 'gridHover'
  | 'handleSelect'
  | 'icon'
  | 'isFetched'
  | 'readonly'
  | 'selectedItems'
  | 'setFinalForm'
  | 'setFormOpen'
  | 'toggleSelection'
>

export function GridView<TEntity extends BaseEntity>({
  data,
  deleteItem,
  entityType,
  getEditFormFromItem,
  grid,
  gridHover,
  handleSelect,
  icon,
  readonly,
  selectedItems,
  setFinalForm,
  setFormOpen,
  toggleSelection
}: GridViewProps<TEntity>) {
  const [hover, setHover] = useState(-1)

  return (
    <div className='grid w-full grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-4'>
      {/* Data Cards */}
      {data.map((item, i) => {
        const isItemSelected = selectedItems.includes(item.id)
        return (
          <Card
            className={`relative flex aspect-auto h-64 flex-col shadow-xs transition-all hover:bg-muted ${
              isItemSelected ? 'ring-1 ring-blue-500' : ''
            } after:border-radius-inherit overflow-visible after:pointer-events-none after:absolute after:top-0 after:left-0 after:z-10 after:h-full after:w-full after:transition-shadow after:content-['']`}
            key={item.id}
          >
            {/* Top Content */}
            <div
              className='group relative grow cursor-pointer overflow-hidden rounded-t-xl transition-all'
              onClick={() => {
                handleSelect(item)
              }}
              onMouseEnter={() => {
                setHover(i)
              }}
              onMouseLeave={() => {
                setHover(-1)
              }}
            >
              {grid ?
                grid(item)
              : <div className='flex h-full w-full items-center justify-center'>
                  {icon}
                </div>
              }
            </div>
            <hr />

            {/* Footer */}
            <div className='mt-auto flex items-center justify-between p-2'>
              <div className='flex min-w-0 items-center gap-2'>
                {!readonly && (
                  <Checkbox
                    aria-label={`Select ${item.name}`}
                    checked={isItemSelected}
                    className='rounded-xs text-blue-600 focus:ring-blue-500'
                    onCheckedChange={() => {
                      toggleSelection(item.id)
                    }}
                  />
                )}
                <span className='overflow-hidden text-base leading-tight text-ellipsis whitespace-nowrap'>
                  {item.name}
                </span>
              </div>
              <div className='flex shrink-0 items-center gap-2'>
                {!readonly && getEditFormFromItem && (
                  <FilePenLine
                    className='h-5 w-5 cursor-pointer text-primary'
                    onClick={() => {
                      setFinalForm(getEditFormFromItem(item))
                      setFormOpen(true)
                    }}
                  />
                )}
                {!readonly && deleteItem ?
                  <DeleteItems
                    deleteItem={deleteItem}
                    entityType={entityType}
                    items={[item]}
                  />
                : null}
              </div>
            </div>
            {gridHover && hover === i && gridHover(item)}
          </Card>
        )
      })}
    </div>
  )
}
