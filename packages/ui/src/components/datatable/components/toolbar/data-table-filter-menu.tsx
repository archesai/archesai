'use no memo'

import type { Column, Table } from '@tanstack/react-table'

import * as React from 'react'
import {
  BadgeCheck,
  CalendarIcon,
  Check,
  ListFilter,
  Text,
  X
} from 'lucide-react'

import type { BaseEntity, FilterCondition } from '@archesai/schemas'

import type { FilterOperator } from '#types/simple-data-table'

import { DataTableRangeFilter } from '#components/datatable/components/filters/data-table-range-filter'
import { Button } from '#components/shadcn/button'
import { Calendar } from '#components/shadcn/calendar'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '#components/shadcn/command'
import { Input } from '#components/shadcn/input'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '#components/shadcn/popover'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '#components/shadcn/select'
import { useFilterState } from '#hooks/use-filter-state'
import { formatDate } from '#lib/format'
import { cn } from '#lib/utils'
import {
  getDefaultFilterOperator,
  getFilterOperators
} from '#types/simple-data-table'

const OPEN_MENU_SHORTCUT = 'f'
const REMOVE_FILTER_SHORTCUTS = ['backspace', 'delete']

interface DataTableFilterItemProps<TData extends BaseEntity> {
  columns: Column<TData>[]
  filter: FilterCondition<TData> & { id: string }
  filterItemId: string
  onFilterRemove: (filterId: string) => void
  onFilterUpdate: (
    filterId: string,
    updates: Partial<Omit<FilterCondition<TData>, 'type'>>
  ) => void
}

interface DataTableFilterMenuProps<TData>
  extends React.ComponentProps<typeof PopoverContent> {
  debounceMs?: number
  shallow?: boolean
  table: Table<TData>
  throttleMs?: number
}

interface FilterValueSelectorProps<TData> {
  column: Column<TData>
  onSelect: (value: string) => void
  value: string
}

export function DataTableFilterMenu<TData extends BaseEntity>({
  align = 'start',
  table,
  ...props
}: DataTableFilterMenuProps<TData>) {
  const id = React.useId()

  const columns = React.useMemo(() => {
    return table
      .getAllColumns()
      .filter((column) => column.columnDef.enableColumnFilter)
  }, [table])

  const [open, setOpen] = React.useState(false)
  const [selectedColumn, setSelectedColumn] =
    React.useState<Column<TData> | null>(null)
  const [inputValue, setInputValue] = React.useState('')
  const triggerRef = React.useRef<HTMLButtonElement>(null)
  const inputRef = React.useRef<HTMLInputElement>(null)

  const onOpenChange = React.useCallback((open: boolean) => {
    setOpen(open)

    if (!open) {
      setTimeout(() => {
        setSelectedColumn(null)
        setInputValue('')
      }, 100)
    }
  }, [])

  const onInputKeyDown = React.useCallback(
    (event: React.KeyboardEvent<HTMLInputElement>) => {
      if (
        REMOVE_FILTER_SHORTCUTS.includes(event.key.toLowerCase()) &&
        !inputValue &&
        selectedColumn
      ) {
        event.preventDefault()
        setSelectedColumn(null)
      }
    },
    [inputValue, selectedColumn]
  )

  const { addCondition, filter, removeCondition, resetFilters, setCondition } =
    useFilterState<TData>()

  // Convert FilterNode to flat array for display
  const filters = React.useMemo(() => {
    const extractConditions = (
      node: typeof filter,
      id = 0
    ): (FilterCondition<TData> & { id: string })[] => {
      if (!node) return []
      if (node.type === 'condition') {
        return [{ ...node, id: `filter-${id.toString()}` }]
      } else {
        return node.children.flatMap((child, index) =>
          extractConditions(child, id * 100 + index)
        )
      }
    }

    return extractConditions(filter)
  }, [filter])

  const onFilterAdd = React.useCallback(
    (column: Column<TData>, value: string) => {
      if (!value.trim() && column.columnDef.meta?.filterVariant !== 'boolean') {
        return
      }

      const filterValue =
        column.columnDef.meta?.filterVariant === 'multiSelect' ? [value] : value

      // Add condition directly
      addCondition({
        field: column.id as keyof TData,
        operator: getDefaultFilterOperator(
          column.columnDef.meta?.filterVariant ?? 'text'
        ),
        type: 'condition',
        value: filterValue
      })

      setOpen(false)

      setTimeout(() => {
        setSelectedColumn(null)
        setInputValue('')
      }, 100)
    },
    [addCondition]
  )

  // Remove filter by field name
  const onFilterRemove = (filterId: string) => {
    const filterToRemove = filters.find((f) => f.id === filterId)
    if (filterToRemove) {
      removeCondition(filterToRemove.field)
    }
    requestAnimationFrame(() => {
      triggerRef.current?.focus()
    })
  }

  // Update filter condition
  const onFilterUpdate = (
    filterId: string,
    updates: Partial<Omit<FilterCondition<TData>, 'type'>>
  ) => {
    const filterToUpdate = filters.find((f) => f.id === filterId)
    if (
      filterToUpdate &&
      updates.field &&
      updates.operator &&
      updates.value !== undefined
    ) {
      setCondition(updates.field, updates.operator, updates.value)
    }
  }

  React.useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (
        event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement
      ) {
        return
      }

      if (
        event.key.toLowerCase() === OPEN_MENU_SHORTCUT &&
        !event.ctrlKey &&
        !event.metaKey &&
        !event.shiftKey
      ) {
        event.preventDefault()
        setOpen(true)
      }

      if (
        event.key.toLowerCase() === OPEN_MENU_SHORTCUT &&
        event.shiftKey &&
        !open &&
        filters.length > 0
      ) {
        event.preventDefault()
        onFilterRemove(filters[filters.length - 1]?.id ?? '')
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
    }
  }, [open, filters, onFilterRemove])

  const onTriggerKeyDown = React.useCallback(
    (event: React.KeyboardEvent<HTMLButtonElement>) => {
      if (
        REMOVE_FILTER_SHORTCUTS.includes(event.key.toLowerCase()) &&
        filters.length > 0
      ) {
        event.preventDefault()
        onFilterRemove(filters[filters.length - 1]?.id ?? '')
      }
    },
    [filters, onFilterRemove]
  )

  return (
    <div className='flex flex-wrap items-center gap-2'>
      {filters.map((filter) => (
        <DataTableFilterItem
          columns={columns}
          filter={filter}
          filterItemId={`${id}-filter-${filter.id}`}
          key={filter.id}
          onFilterRemove={onFilterRemove}
          onFilterUpdate={onFilterUpdate}
        />
      ))}
      {filters.length > 0 && (
        <Button
          aria-label='Reset all filters'
          className='size-8'
          onClick={resetFilters}
          size='icon'
          variant='outline'
        >
          <X />
        </Button>
      )}
      <Popover
        onOpenChange={onOpenChange}
        open={open}
      >
        <PopoverTrigger asChild>
          <Button
            aria-label='Open filter command menu'
            className={cn(filters.length > 0 && 'size-8', 'h-8')}
            onKeyDown={onTriggerKeyDown}
            ref={triggerRef}
            size={filters.length > 0 ? 'icon' : 'sm'}
            variant='outline'
          >
            <ListFilter />
            {filters.length > 0 ? null : 'Filter'}
          </Button>
        </PopoverTrigger>
        <PopoverContent
          align={align}
          className='w-full max-w-[var(--radix-popover-content-available-width)] origin-[var(--radix-popover-content-transform-origin)] p-0'
          {...props}
        >
          <Command
            className='[&_[cmdk-input-wrapper]_svg]:hidden'
            loop
          >
            <CommandInput
              onKeyDown={onInputKeyDown}
              onValueChange={setInputValue}
              placeholder={
                selectedColumn ?
                  (selectedColumn.columnDef.meta?.label ?? selectedColumn.id)
                : 'Search fields...'
              }
              ref={inputRef}
              value={inputValue}
            />
            <CommandList>
              {selectedColumn ?
                <>
                  {selectedColumn.columnDef.meta?.options && (
                    <CommandEmpty>No options found.</CommandEmpty>
                  )}
                  <FilterValueSelector
                    column={selectedColumn}
                    onSelect={(value) => {
                      onFilterAdd(selectedColumn, value)
                    }}
                    value={inputValue}
                  />
                </>
              : <>
                  <CommandEmpty>No fields found.</CommandEmpty>
                  <CommandGroup>
                    {columns.map((column) => (
                      <CommandItem
                        key={column.id}
                        onSelect={() => {
                          setSelectedColumn(column)
                          setInputValue('')
                          requestAnimationFrame(() => {
                            inputRef.current?.focus()
                          })
                        }}
                        value={column.id}
                      >
                        {column.columnDef.meta?.icon && (
                          <column.columnDef.meta.icon />
                        )}
                        <span className='truncate'>
                          {column.columnDef.meta?.label ?? column.id}
                        </span>
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </>
              }
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  )
}

function DataTableFilterItem<TData extends BaseEntity>({
  columns,
  filter,
  filterItemId,
  onFilterRemove,
  onFilterUpdate
}: DataTableFilterItemProps<TData>) {
  const [showFieldSelector, setShowFieldSelector] = React.useState(false)
  const [showOperatorSelector, setShowOperatorSelector] = React.useState(false)
  const [showValueSelector, setShowValueSelector] = React.useState(false)

  const column = columns.find((column) => column.id === filter.field)

  const operatorListboxId = `${filterItemId}-operator-listbox`
  const inputId = `${filterItemId}-input`

  const columnMeta = column?.columnDef.meta
  const filterOperators = getFilterOperators(
    columnMeta?.filterVariant ?? 'text'
  )

  const onItemKeyDown = React.useCallback(
    (event: React.KeyboardEvent<HTMLDivElement>) => {
      if (
        event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement
      ) {
        return
      }

      if (showFieldSelector || showOperatorSelector || showValueSelector) {
        return
      }

      if (REMOVE_FILTER_SHORTCUTS.includes(event.key.toLowerCase())) {
        event.preventDefault()
        onFilterRemove(filter.id)
      }
    },
    [
      filter.id,
      showFieldSelector,
      showOperatorSelector,
      showValueSelector,
      onFilterRemove
    ]
  )

  if (!column) return null

  return (
    <div
      className='flex h-8 items-center rounded-md bg-background'
      id={filterItemId}
      key={filter.id}
      onKeyDown={onItemKeyDown}
      role='listitem'
    >
      <Popover
        onOpenChange={setShowFieldSelector}
        open={showFieldSelector}
      >
        <PopoverTrigger asChild>
          <Button
            className='rounded-none rounded-l-md border border-r-0 font-normal dark:bg-input/30'
            size='sm'
            variant='ghost'
          >
            {columnMeta?.icon && (
              <columnMeta.icon className='text-muted-foreground' />
            )}
            {columnMeta?.label ?? column.id}
          </Button>
        </PopoverTrigger>
        <PopoverContent
          align='start'
          className='w-48 origin-[var(--radix-popover-content-transform-origin)] p-0'
        >
          <Command loop>
            <CommandInput placeholder='Search fields...' />
            <CommandList>
              <CommandEmpty>No fields found.</CommandEmpty>
              <CommandGroup>
                {columns.map((column) => (
                  <CommandItem
                    key={column.id}
                    onSelect={() => {
                      onFilterUpdate(filter.id, {
                        field: column.id as keyof TData,
                        operator: getDefaultFilterOperator(
                          column.columnDef.meta?.filterVariant ?? 'text'
                        ),
                        value: ''
                      })

                      setShowFieldSelector(false)
                    }}
                    value={column.id}
                  >
                    {column.columnDef.meta?.icon && (
                      <column.columnDef.meta.icon />
                    )}
                    <span className='truncate'>
                      {column.columnDef.meta?.label ?? column.id}
                    </span>
                    <Check
                      className={cn(
                        'ml-auto',
                        column.id === filter.field ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
      <Select
        onOpenChange={setShowOperatorSelector}
        onValueChange={(value: FilterOperator) => {
          onFilterUpdate(filter.id, {
            operator: value,
            value:
              value === 'isEmpty' || value === 'isNotEmpty' ? '' : filter.value
          })
        }}
        open={showOperatorSelector}
        value={filter.operator}
      >
        <SelectTrigger
          aria-controls={operatorListboxId}
          className='h-8 rounded-none border-r-0 px-2.5 lowercase [&_svg]:hidden [&[data-size]]:h-8'
        >
          <SelectValue placeholder={filter.operator} />
        </SelectTrigger>
        <SelectContent
          className='origin-[var(--radix-select-content-transform-origin)]'
          id={operatorListboxId}
        >
          {filterOperators.map((operator) => (
            <SelectItem
              className='lowercase'
              key={operator.value}
              value={operator.value}
            >
              {operator.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      {onFilterInputRender({
        column,
        filter,
        inputId,
        onFilterUpdate,
        setShowValueSelector,
        showValueSelector
      })}
      <Button
        aria-controls={filterItemId}
        className='h-full rounded-none rounded-r-md border border-l-0 px-1.5 font-normal dark:bg-input/30'
        onClick={() => {
          onFilterRemove(filter.id)
        }}
        size='sm'
        variant='ghost'
      >
        <X className='size-3.5' />
      </Button>
    </div>
  )
}

function FilterValueSelector<TData>({
  column,
  onSelect,
  value
}: FilterValueSelectorProps<TData>) {
  const variant = column.columnDef.meta?.filterVariant ?? 'text'

  switch (variant) {
    case 'boolean':
      return (
        <CommandGroup>
          <CommandItem
            onSelect={() => {
              onSelect('true')
            }}
            value='true'
          >
            True
          </CommandItem>
          <CommandItem
            onSelect={() => {
              onSelect('false')
            }}
            value='false'
          >
            False
          </CommandItem>
        </CommandGroup>
      )

    case 'date':
      return (
        <Calendar
          captionLayout='dropdown'
          mode='single'
          onSelect={(date) => {
            onSelect(date?.getTime().toString() ?? '')
          }}
          selected={value ? new Date(value) : undefined}
        />
      )

    case 'multiSelect':
    case 'select':
      return (
        <CommandGroup>
          {column.columnDef.meta?.options?.map((option) => (
            <CommandItem
              key={option.value}
              onSelect={() => {
                onSelect(option.value)
              }}
              value={option.value}
            >
              {option.icon && <option.icon />}
              <span className='truncate'>{option.label}</span>
              {option.count && (
                <span className='ml-auto font-mono text-xs'>
                  {option.count}
                </span>
              )}
            </CommandItem>
          ))}
        </CommandGroup>
      )

    default: {
      const isEmpty = !value.trim()

      return (
        <CommandGroup>
          <CommandItem
            disabled={isEmpty}
            onSelect={() => {
              onSelect(value)
            }}
            value={value}
          >
            {isEmpty ?
              <>
                <Text />
                <span>Type to add filter...</span>
              </>
            : <>
                <BadgeCheck />
                <span className='truncate'>Filter by &quot;{value}&quot;</span>
              </>
            }
          </CommandItem>
        </CommandGroup>
      )
    }
  }
}

function onFilterInputRender<TData extends BaseEntity>({
  column,
  filter,
  inputId,
  onFilterUpdate,
  setShowValueSelector,
  showValueSelector
}: {
  column: Column<TData>
  filter: FilterCondition<TData> & { id: string }
  inputId: string
  onFilterUpdate: (
    filterId: string,
    updates: Partial<Omit<FilterCondition<TData>, 'type'>>
  ) => void
  setShowValueSelector: (value: boolean) => void
  showValueSelector: boolean
}) {
  if (filter.operator === 'isEmpty' || filter.operator === 'isNotEmpty') {
    return (
      <div
        aria-label={`${column.columnDef.meta?.label ?? ''} filter is ${
          filter.operator === 'isEmpty' ? 'empty' : 'not empty'
        }`}
        aria-live='polite'
        className='h-full w-16 rounded-none border bg-transparent px-1.5 py-0.5 text-muted-foreground dark:bg-input/30'
        id={inputId}
        role='status'
      />
    )
  }

  const variant = column.columnDef.meta?.filterVariant ?? 'text'
  switch (variant) {
    case 'boolean': {
      const inputListboxId = `${inputId}-listbox`

      return (
        <Select
          onOpenChange={setShowValueSelector}
          onValueChange={(value: 'false' | 'true') => {
            onFilterUpdate(filter.id, { value })
          }}
          open={showValueSelector}
          value={typeof filter.value === 'string' ? filter.value : 'true'}
        >
          <SelectTrigger
            aria-controls={inputListboxId}
            className='rounded-none bg-transparent px-1.5 py-0.5 [&_svg]:hidden'
            id={inputId}
          >
            <SelectValue placeholder={filter.value ? 'True' : 'False'} />
          </SelectTrigger>
          <SelectContent id={inputListboxId}>
            <SelectItem value='true'>True</SelectItem>
            <SelectItem value='false'>False</SelectItem>
          </SelectContent>
        </Select>
      )
    }

    case 'date':
    case 'dateRange': {
      const inputListboxId = `${inputId}-listbox`

      const dateValue =
        Array.isArray(filter.value) ?
          filter.value.filter(Boolean)
        : [filter.value, filter.value].filter(Boolean)

      const displayValue =
        filter.operator === 'isBetween' && dateValue.length === 2 ?
          `${formatDate(new Date(Number(dateValue[0])))} - ${formatDate(
            new Date(Number(dateValue[1]))
          )}`
        : dateValue[0] ? formatDate(new Date(Number(dateValue[0])))
        : 'Pick date...'

      return (
        <Popover
          onOpenChange={setShowValueSelector}
          open={showValueSelector}
        >
          <PopoverTrigger asChild>
            <Button
              aria-controls={inputListboxId}
              className={cn(
                'h-full rounded-none border px-1.5 font-normal dark:bg-input/30',
                !filter.value && 'text-muted-foreground'
              )}
              id={inputId}
              size='sm'
              variant='ghost'
            >
              <CalendarIcon className='size-3.5' />
              <span className='truncate'>{displayValue}</span>
            </Button>
          </PopoverTrigger>
          <PopoverContent
            align='start'
            className='w-auto origin-[var(--radix-popover-content-transform-origin)] p-0'
            id={inputListboxId}
          >
            {filter.operator === 'isBetween' ?
              <Calendar
                captionLayout='dropdown'
                mode='range'
                onSelect={(date) => {
                  onFilterUpdate(filter.id, {
                    value:
                      date ?
                        [
                          (date.from?.getTime() ?? '').toString(),
                          (date.to?.getTime() ?? '').toString()
                        ]
                      : []
                  })
                }}
                selected={
                  dateValue.length === 2 ?
                    {
                      from: new Date(Number(dateValue[0])),
                      to: new Date(Number(dateValue[1]))
                    }
                  : {
                      from: new Date(),
                      to: new Date()
                    }
                }
              />
            : <Calendar
                captionLayout='dropdown'
                mode='single'
                onSelect={(date) => {
                  onFilterUpdate(filter.id, {
                    value: (date?.getTime() ?? '').toString()
                  })
                }}
                selected={
                  dateValue[0] ? new Date(Number(dateValue[0])) : undefined
                }
              />
            }
          </PopoverContent>
        </Popover>
      )
    }

    case 'multiSelect':
    case 'select': {
      const inputListboxId = `${inputId}-listbox`

      const options = column.columnDef.meta?.options ?? []
      const selectedValues =
        Array.isArray(filter.value) ? filter.value : [filter.value]

      const selectedOptions = options.filter((option) =>
        selectedValues.includes(option.value)
      )

      return (
        <Popover
          onOpenChange={setShowValueSelector}
          open={showValueSelector}
        >
          <PopoverTrigger asChild>
            <Button
              aria-controls={inputListboxId}
              className='h-full min-w-16 rounded-none border px-1.5 font-normal dark:bg-input/30'
              id={inputId}
              size='sm'
              variant='ghost'
            >
              {selectedOptions.length === 0 ?
                variant === 'multiSelect' ?
                  'Select options...'
                : 'Select option...'
              : <>
                  <div className='flex items-center -space-x-2 rtl:space-x-reverse'>
                    {selectedOptions.map((selectedOption) =>
                      selectedOption.icon ?
                        <div
                          className='rounded-full border bg-background p-0.5'
                          key={selectedOption.value}
                        >
                          <selectedOption.icon className='size-3.5' />
                        </div>
                      : null
                    )}
                  </div>
                  <span className='truncate'>
                    {selectedOptions.length > 1 ?
                      `${selectedOptions.length.toString()} selected`
                    : selectedOptions[0]?.label}
                  </span>
                </>
              }
            </Button>
          </PopoverTrigger>
          <PopoverContent
            align='start'
            className='w-48 origin-[var(--radix-popover-content-transform-origin)] p-0'
            id={inputListboxId}
          >
            <Command>
              <CommandInput placeholder='Search options...' />
              <CommandList>
                <CommandEmpty>No options found.</CommandEmpty>
                <CommandGroup>
                  {options.map((option) => (
                    <CommandItem
                      key={option.value}
                      onSelect={() => {
                        const value =
                          variant === 'multiSelect' ?
                            selectedValues.includes(option.value) ?
                              selectedValues.filter((v) => v !== option.value)
                            : [...selectedValues, option.value]
                          : option.value
                        onFilterUpdate(filter.id, { value })
                      }}
                      value={option.value}
                    >
                      {option.icon && <option.icon />}
                      <span className='truncate'>{option.label}</span>
                      {variant === 'multiSelect' && (
                        <Check
                          className={cn(
                            'ml-auto',
                            selectedValues.includes(option.value) ?
                              'opacity-100'
                            : 'opacity-0'
                          )}
                        />
                      )}
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>
      )
    }
    case 'number':
    case 'range':
    case 'text': {
      if (
        (variant === 'range' && filter.operator === 'isBetween') ||
        filter.operator === 'isBetween'
      ) {
        return (
          <DataTableRangeFilter
            className="size-full max-w-28 gap-0 [&_[data-slot='range-min']]:border-r-0 [&_input]:rounded-none [&_input]:px-1.5"
            column={column}
            filter={filter}
            inputId={inputId}
            onFilterUpdate={onFilterUpdate}
          />
        )
      }

      const isNumber = variant === 'number' || variant === 'range'

      return (
        <Input
          className='h-full w-24 rounded-none px-1.5'
          defaultValue={typeof filter.value === 'string' ? filter.value : ''}
          id={inputId}
          inputMode={isNumber ? 'numeric' : undefined}
          onChange={(event) => {
            onFilterUpdate(filter.id, { value: event.target.value })
          }}
          placeholder={column.columnDef.meta?.label ?? 'Enter value...'}
          type={isNumber ? 'number' : 'text'}
        />
      )
    }

    default:
      return null
  }
}
