// components/DataSelector.tsx

'use client'

import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList
} from '@/components/ui/command'
import { HoverCard, HoverCardContent } from '@/components/ui/hover-card'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover'
import { cn } from '@/lib/utils'
import { CaretSortIcon, CheckIcon, PlusIcon } from '@radix-ui/react-icons'
import React, { useCallback, useState } from 'react'

import { BaseItem } from './datatables/datatable/data-table'
import { RunsControllerFindAllQueryParams } from '@/generated/archesApiComponents'

interface DataSelectorProps<TItem extends BaseItem, TFindAllPathParams> {
  getItemDetails?: (item: TItem) => React.ReactNode
  //   iconMap?: { [key: string]: IconType };
  icons?: { Icon: any; name: string }[]
  isMultiSelect?: boolean
  itemType: string
  selectedData: TItem | TItem[] | undefined
  orgname?: string
  setSelectedData: (data: TItem | TItem[] | undefined) => void
  useFindAll: (s: {
    pathParams: TFindAllPathParams
    queryParams?: RunsControllerFindAllQueryParams
  }) => {
    data:
      | undefined
      | {
          results: TItem[]
        }
    isLoading: boolean
  }
  findAllParams: TFindAllPathParams
}

export function DataSelector<TItem extends BaseItem, TFindAllPathParams>({
  getItemDetails,
  icons,
  //   iconMap,
  isMultiSelect = false,
  itemType,
  selectedData,
  setSelectedData,
  useFindAll,
  findAllParams
}: DataSelectorProps<TItem, TFindAllPathParams>) {
  const { data } = useFindAll({
    pathParams: findAllParams,
    queryParams: { limit: 10 }
  })
  const [open, setOpen] = useState(false)
  const [hoveredItem, setHoveredItem] = useState<TItem | undefined>()
  const [searchTerm, setSearchTerm] = useState('')

  // Filter data based on search term
  const filteredData = React.useMemo(() => {
    if (!data) return []
    return (
      data.results.filter((item: any) =>
        item.name.toLowerCase().includes(searchTerm.toLowerCase())
      ) || []
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
        const isSelected = selectedArray.some(
          (i) => (i as any).id === (item as any).id
        )
        if (isSelected) {
          setSelectedData(
            selectedArray.filter((i) => (i as any).id !== (item as any).id)
          )
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
        setSelectedData(
          selectedData.filter((i) => (i as any).id !== (item as any).id)
        )
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
                  (selectedData as TItem[])?.length > 0 || !isMultiSelect
                    ? ''
                    : 'text-muted-foreground'
                )}
              >
                {isMultiSelect ? (
                  (selectedData as TItem[])?.length > 0 ? (
                    (selectedData as TItem[])?.length + ' selected'
                  ) : (
                    'Select ' + itemType.toLowerCase() + 's...'
                  )
                ) : (
                  <div className='flex items-center gap-1'>
                    {icons &&
                      icons
                        .filter((x) => x.name === (selectedData as TItem)?.name)
                        .map((x: any, i) => (
                          <x.Icon
                            className='h-4 w-4 text-muted-foreground'
                            key={i}
                          />
                        ))}
                    {(selectedData as TItem)?.name}
                  </div>
                )}
              </div>
            </div>
            {isMultiSelect ? (
              <PlusIcon className='ml-2 h-4 w-4 shrink-0 opacity-50' />
            ) : (
              <CaretSortIcon className='ml-2 h-4 w-4 shrink-0 opacity-50' />
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className='w-[300px] p-0'>
          <Command>
            <CommandInput
              onChangeCapture={(e: any) => setSearchTerm(e.target.value)}
              placeholder={`Search ${itemType.toLowerCase()}...`}
              value={searchTerm}
            />
            <CommandList className='max-h-[400px] overflow-y-auto'>
              <CommandEmpty>No {itemType.toLowerCase()} found.</CommandEmpty>
              {filteredData.map((item) => (
                <CommandItem
                  className={cn('flex items-center justify-between')}
                  key={(item as any).id}
                  onMouseEnter={() => setHoveredItem(item)}
                  onMouseLeave={() => setHoveredItem(undefined)}
                  onSelect={() => handleSelect(item)}
                >
                  <div className='flex items-center gap-1'>
                    {/* Icon Rendering */}
                    {icons &&
                      icons
                        .filter((x) => x.name === item.name)
                        .map((x: any, i) => (
                          <x.Icon
                            className='h-4 w-4 text-muted-foreground'
                            key={i}
                          />
                        ))}
                    <p>{item?.name}</p>
                  </div>
                  {/* Check Icon for Selected Items */}
                  {isMultiSelect &&
                    Array.isArray(selectedData) &&
                    selectedData.some(
                      (i) => (i as any).id === (item as any).id
                    ) && <CheckIcon className='h-4 w-4 text-green-500' />}
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
                className='inline-flex items-center rounded bg-blue-100 px-2 py-1 text-sm text-blue-700'
                key={(item as any).id}
              >
                {item.name}
                <button
                  className='ml-1 text-red-500 hover:text-red-700'
                  onClick={() => handleRemove(item)}
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
