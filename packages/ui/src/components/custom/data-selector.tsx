'use client'

import type { LucideIcon } from 'lucide-react'

import { useCallback, useMemo, useState } from 'react'
import { CheckSquareIcon, PlusSquareIcon, SortAsc } from 'lucide-react'

import type { BaseEntity } from '@archesai/domain'

import type { TFindManyResponse } from '#components/datatable/data-table'

import { Button } from '#components/shadcn/button'
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList
} from '#components/shadcn/command'
import { HoverCard, HoverCardContent } from '#components/shadcn/hover-card'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '#components/shadcn/popover'
import { cn } from '#lib/utils'

interface DataSelectorProps<
  TItem extends BaseEntity,
  TFindAllPathParams,
  FindManyOptions extends {
    query?: {
      enabled?: boolean
    }
  } = {
    query?: {
      enabled?: boolean
    }
  }
> {
  findManyParams: TFindAllPathParams
  getItemDetails?: (item: TItem) => React.ReactNode
  //   iconMap?: { [key: string]: IconType };
  icons?: { color: string; Icon: LucideIcon; name: string }[]
  isMultiSelect?: boolean
  itemType: string
  orgname?: string
  selectedData: TItem | TItem[] | undefined
  setSelectedData: (data: TItem | TItem[] | undefined) => void
  useFindMany: (
    params: TFindAllPathParams,
    options: FindManyOptions
  ) => {
    data: TFindManyResponse<TItem> | undefined
    isFetched: boolean
  }
}

export function DataSelector<TItem extends BaseEntity, TFindAllPathParams>({
  findManyParams,
  getItemDetails,
  icons,
  isMultiSelect = false,
  itemType,
  selectedData,
  setSelectedData,
  useFindMany
}: DataSelectorProps<TItem, TFindAllPathParams>) {
  const { data } = useFindMany(findManyParams, {})
  const [open, setOpen] = useState(false)
  const [hoveredItem, setHoveredItem] = useState<TItem | undefined>()
  const [searchTerm, setSearchTerm] = useState('')
  if (!data) return null

  // Filter data based on search term
  const filteredData = useMemo(() => {
    return data.data.data.filter((item) =>
      item.attributes.name.toLowerCase().includes(searchTerm.toLowerCase())
    )
  }, [data, searchTerm])

  // Handler for selecting/deselecting items
  const handleSelect = useCallback(
    (item: TItem) => {
      if (isMultiSelect) {
        const selectedArray = Array.isArray(selectedData)
          ? selectedData
          : selectedData
            ? [selectedData]
            : []
        const isSelected = selectedArray.some((i) => i.id === item.id)
        if (isSelected) {
          setSelectedData(selectedArray.filter((i) => i.id !== item.id))
        } else {
          setSelectedData([...selectedArray, item])
        }
      } else {
        setSelectedData(item)
        setOpen(false) // Close popover on single select
      }
    },
    [isMultiSelect, selectedData, setSelectedData]
  )

  // Handler for removing selected item (for multi-select)
  const handleRemove = useCallback(
    (item: TItem) => {
      if (isMultiSelect && Array.isArray(selectedData)) {
        setSelectedData(selectedData.filter((i) => i.id !== item.id))
      } else {
        setSelectedData(undefined)
      }
    },
    [isMultiSelect, selectedData, setSelectedData]
  )

  return (
    <div className='flex flex-col gap-2'>
      {/* Popover for Selection */}
      <Popover
        onOpenChange={setOpen}
        open={open}
      >
        <PopoverTrigger asChild>
          <Button
            aria-expanded={open}
            aria-label={`Select ${itemType.toLowerCase()}`}
            className='w-full justify-between'
            role='combobox'
            variant='outline'
          >
            <div className='flex items-center gap-2'>
              <div
                className={cn(
                  'flex flex-wrap gap-2',
                  (Array.isArray(selectedData) && selectedData.length > 0) ||
                    !isMultiSelect
                    ? ''
                    : 'text-muted-foreground'
                )}
              >
                {isMultiSelect ? (
                  Array.isArray(selectedData) && selectedData.length > 0 ? (
                    selectedData.length.toString() + ' selected'
                  ) : (
                    'Select ' + itemType.toLowerCase() + 's...'
                  )
                ) : (
                  <div className='flex items-center gap-1'>
                    {icons
                      ?.filter((x) => x.name === (selectedData as TItem).name)
                      .map((x, i) => {
                        const iconColor = x.color
                        return (
                          <x.Icon
                            className={cn(
                              'mx-auto h-4 w-4',
                              iconColor.startsWith('text-') ? iconColor : ''
                            )}
                            key={i}
                            style={{
                              ...(iconColor.startsWith('#')
                                ? { color: iconColor }
                                : {})
                            }}
                          />
                        )
                      })}
                    {(selectedData as TItem).name}
                  </div>
                )}
              </div>
            </div>
            {isMultiSelect ? (
              <PlusSquareIcon className='ml-2 h-4 w-4 shrink-0 opacity-50' />
            ) : (
              <SortAsc className='ml-2 h-4 w-4 shrink-0 opacity-50' />
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className='w-[300px] p-0'>
          <Command>
            <CommandInput
              onInput={(e) => {
                setSearchTerm(e.currentTarget.value)
              }}
              placeholder={`Search ${itemType.toLowerCase()}...`}
              value={searchTerm}
            />
            <CommandList className='max-h-[400px] overflow-y-auto'>
              <CommandEmpty>No {itemType.toLowerCase()} found.</CommandEmpty>
              {filteredData.map((item) => (
                <CommandItem
                  className={cn('flex items-center justify-between')}
                  key={item.id}
                  onMouseEnter={() => {
                    setHoveredItem({
                      id: item.id,
                      type: itemType,
                      ...item.attributes
                    } as TItem)
                  }}
                  onMouseLeave={() => {
                    setHoveredItem(undefined)
                  }}
                  onSelect={() => {
                    handleSelect({
                      id: item.id,
                      type: itemType,
                      ...item.attributes
                    } as TItem)
                  }}
                >
                  <div className='flex items-center gap-1'>
                    {/* Icon Rendering */}
                    {icons
                      ?.filter((x) => x.name === item.attributes.name)
                      .map((x, i) => {
                        const iconColor = x.color
                        return (
                          <x.Icon
                            className={cn(
                              'mx-auto h-4 w-4',
                              iconColor.startsWith('text-') ? iconColor : ''
                            )}
                            key={i}
                            style={{
                              ...(iconColor.startsWith('#')
                                ? { color: iconColor }
                                : {})
                            }}
                          />
                        )
                      })}
                    <p>{item.attributes.name}</p>
                  </div>
                  {/* Check Icon for Selected Items */}
                  {isMultiSelect &&
                    Array.isArray(selectedData) &&
                    selectedData.some((i) => i.id === item.id) && (
                      <CheckSquareIcon className='h-4 w-4 text-green-500' />
                    )}
                </CommandItem>
              ))}
            </CommandList>
          </Command>
          {/* HoverCard for Item Details */}
          {hoveredItem && getItemDetails && (
            <HoverCard open={true}>
              <HoverCardContent
                align='end'
                className='w-[250px] p-4'
              >
                {getItemDetails(hoveredItem)}
              </HoverCardContent>
            </HoverCard>
          )}
        </PopoverContent>
      </Popover>

      {/* Display Selected Items (for Multi-Select) */}
      {isMultiSelect &&
        Array.isArray(selectedData) &&
        selectedData.length > 0 && (
          <div className='flex flex-wrap gap-2'>
            {selectedData.map((item) => (
              <span
                className='inline-flex items-center rounded-xs bg-blue-100 px-2 py-1 text-sm text-blue-700'
                key={item.id}
              >
                {item.name}
                <button
                  className='ml-1 text-red-500 hover:text-red-700'
                  onClick={() => {
                    handleRemove(item)
                  }}
                  type='button'
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
    </div>
  )
}
