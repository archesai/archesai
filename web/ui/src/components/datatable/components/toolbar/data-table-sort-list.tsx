"use no memo"

import type { ColumnSort, SortDirection, Table } from "@tanstack/react-table"
import type { JSX } from "react"

import { useCallback, useEffect, useId, useMemo, useRef, useState } from "react"

import {
  ArrowDownUpIcon,
  ChevronsUpDownIcon,
  GripVerticalIcon,
  TrashIcon
} from "#components/custom/icons"
import {
  Sortable,
  SortableContent,
  SortableItem,
  SortableItemHandle,
  SortableOverlay
} from "#components/custom/sortable"
import { Badge } from "#components/shadcn/badge"
import { Button } from "#components/shadcn/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from "#components/shadcn/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from "#components/shadcn/popover"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "#components/shadcn/select"
import { cn } from "#lib/utils"

const OPEN_MENU_SHORTCUT = "s"
const REMOVE_SORT_SHORTCUTS = ["backspace", "delete"]

interface DataTableSortItemProps {
  columnLabels: Map<string, string>
  columns: { id: string; label: string }[]
  onSortRemove: (sortId: string) => void
  onSortUpdate: (sortId: string, updates: Partial<ColumnSort>) => void
  sort: ColumnSort
  sortItemId: string
}

interface DataTableSortListProps<TData>
  extends React.ComponentProps<typeof PopoverContent> {
  table: Table<TData>
}

export function DataTableSortList<TData>({
  table,
  ...props
}: DataTableSortListProps<TData>): JSX.Element {
  const id = useId()
  const labelId = useId()
  const descriptionId = useId()
  const [open, setOpen] = useState(false)
  const addButtonRef = useRef<HTMLButtonElement>(null)

  const sorting = table.getState().sorting
  const onSortingChange = table.setSorting

  const { columnLabels, columns } = useMemo(() => {
    const labels = new Map<string, string>()
    const sortingIds = new Set(sorting.map((s) => s.id))
    const availableColumns: { id: string; label: string }[] = []

    for (const column of table.getAllColumns()) {
      if (!column.getCanSort()) continue

      const label = column.columnDef.meta?.label ?? column.id
      labels.set(column.id, label)

      if (!sortingIds.has(column.id)) {
        availableColumns.push({ id: column.id, label })
      }
    }

    return {
      columnLabels: labels,
      columns: availableColumns
    }
  }, [sorting, table])

  const onSortAdd = useCallback(() => {
    const firstColumn = columns[0]
    if (!firstColumn) return

    onSortingChange((prevSorting) => [
      ...prevSorting,
      { desc: false, id: firstColumn.id }
    ])
  }, [columns, onSortingChange])

  const onSortUpdate = useCallback(
    (sortId: string, updates: Partial<ColumnSort>) => {
      onSortingChange((prevSorting) => {
        if (!prevSorting[0]) return prevSorting
        return prevSorting.map((sort) =>
          sort.id === sortId ? { ...sort, ...updates } : sort
        )
      })
    },
    [onSortingChange]
  )

  const onSortRemove = useCallback(
    (sortId: string) => {
      onSortingChange((prevSorting) =>
        prevSorting.filter((item) => item.id !== sortId)
      )
    },
    [onSortingChange]
  )

  const onSortingReset = useCallback(() => {
    onSortingChange(table.initialState.sorting)
  }, [onSortingChange, table.initialState.sorting])

  useEffect(() => {
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
        sorting.length > 0
      ) {
        event.preventDefault()
        onSortingReset()
      }
    }

    window.addEventListener("keydown", onKeyDown)
    return () => {
      window.removeEventListener("keydown", onKeyDown)
    }
  }, [sorting.length, onSortingReset])

  const onTriggerKeyDown = useCallback(
    (event: React.KeyboardEvent<HTMLButtonElement>) => {
      if (
        REMOVE_SORT_SHORTCUTS.includes(event.key.toLowerCase()) &&
        sorting.length > 0
      ) {
        event.preventDefault()
        onSortingReset()
      }
    },
    [sorting.length, onSortingReset]
  )

  return (
    <Sortable
      getItemValue={(item) => item.id}
      onValueChange={onSortingChange}
      value={sorting}
    >
      <Popover
        onOpenChange={setOpen}
        open={open}
      >
        <PopoverTrigger asChild>
          <Button
            onKeyDown={onTriggerKeyDown}
            size="sm"
            variant="ghost"
          >
            <ArrowDownUpIcon />
            Sort
            {sorting.length > 0 && (
              <Badge
                className="h-[18.24px] rounded-[3.2px] px-[5.12px] font-mono text-[10.4px] font-normal"
                variant="secondary"
              >
                {sorting.length}
              </Badge>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent
          aria-describedby={descriptionId}
          aria-labelledby={labelId}
          className="flex w-full max-w-[var(--radix-popover-content-available-width)] origin-[var(--radix-popover-content-transform-origin)] flex-col gap-3.5 p-4 sm:min-w-[380px]"
          {...props}
        >
          <div className="flex flex-col gap-1">
            <h4
              className="leading-none font-medium"
              id={labelId}
            >
              {sorting.length > 0 ? "Sort by" : "No sorting applied"}
            </h4>
            <p
              className={cn(
                "text-sm text-muted-foreground",
                sorting.length > 0 && "sr-only"
              )}
              id={descriptionId}
            >
              {sorting.length > 0
                ? "Modify sorting to organize your rows."
                : "Add sorting to organize your rows."}
            </p>
          </div>
          {sorting.length > 0 && (
            <SortableContent asChild>
              <ul className="flex max-h-[300px] flex-col gap-2 overflow-y-auto p-1">
                {sorting.map((sort) => (
                  <DataTableSortItem
                    columnLabels={columnLabels}
                    columns={columns}
                    key={sort.id}
                    onSortRemove={onSortRemove}
                    onSortUpdate={onSortUpdate}
                    sort={sort}
                    sortItemId={`${id}-sort-${sort.id}`}
                  />
                ))}
              </ul>
            </SortableContent>
          )}
          <div className="flex w-full items-center gap-2">
            <Button
              className="rounded"
              disabled={columns.length === 0}
              onClick={onSortAdd}
              ref={addButtonRef}
              size="sm"
            >
              Add sort
            </Button>
            {sorting.length > 0 && (
              <Button
                className="rounded"
                onClick={onSortingReset}
                size="sm"
                variant="ghost"
              >
                Reset sorting
              </Button>
            )}
          </div>
        </PopoverContent>
      </Popover>
      <SortableOverlay>
        <div className="flex items-center gap-2">
          <div className="h-8 w-[180px] rounded-sm bg-primary/10" />
          <div className="h-8 w-24 rounded-sm bg-primary/10" />
          <div className="size-8 shrink-0 rounded-sm bg-primary/10" />
          <div className="size-8 shrink-0 rounded-sm bg-primary/10" />
        </div>
      </SortableOverlay>
    </Sortable>
  )
}

function DataTableSortItem({
  columnLabels,
  columns,
  onSortRemove,
  onSortUpdate,
  sort,
  sortItemId
}: DataTableSortItemProps) {
  const fieldListboxId = `${sortItemId}-field-listbox`
  const fieldTriggerId = `${sortItemId}-field-trigger`
  const directionListboxId = `${sortItemId}-direction-listbox`

  const [showFieldSelector, setShowFieldSelector] = useState(false)
  const [showDirectionSelector, setShowDirectionSelector] = useState(false)

  const onItemKeyDown = useCallback(
    (event: React.KeyboardEvent<HTMLLIElement>) => {
      if (
        event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement
      ) {
        return
      }

      if (showFieldSelector || showDirectionSelector) {
        return
      }

      if (REMOVE_SORT_SHORTCUTS.includes(event.key.toLowerCase())) {
        event.preventDefault()
        onSortRemove(sort.id)
      }
    },
    [sort.id, showFieldSelector, showDirectionSelector, onSortRemove]
  )

  return (
    <SortableItem
      asChild
      value={sort.id}
    >
      <li
        className="flex items-center gap-2"
        id={sortItemId}
        onKeyDown={onItemKeyDown}
        tabIndex={-1}
      >
        <Popover
          onOpenChange={setShowFieldSelector}
          open={showFieldSelector}
        >
          <PopoverTrigger asChild>
            <Button
              aria-controls={fieldListboxId}
              className="w-44 justify-between rounded font-normal"
              id={fieldTriggerId}
              size="sm"
              variant="ghost"
            >
              <span className="truncate">{columnLabels.get(sort.id)}</span>
              <ChevronsUpDownIcon className="opacity-50" />
            </Button>
          </PopoverTrigger>
          <PopoverContent
            className="w-[var(--radix-popover-trigger-width)] origin-[var(--radix-popover-content-transform-origin)] p-0"
            id={fieldListboxId}
          >
            <Command>
              <CommandInput placeholder="Search fields..." />
              <CommandList>
                <CommandEmpty>No fields found.</CommandEmpty>
                <CommandGroup>
                  {columns.map((column) => (
                    <CommandItem
                      key={column.id}
                      onSelect={(value) => {
                        onSortUpdate(sort.id, { id: value })
                      }}
                      value={column.id}
                    >
                      <span className="truncate">{column.label}</span>
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>
        <Select
          onOpenChange={setShowDirectionSelector}
          onValueChange={(value: SortDirection) => {
            onSortUpdate(sort.id, { desc: value === "desc" })
          }}
          open={showDirectionSelector}
          value={sort.desc ? "desc" : "asc"}
        >
          <SelectTrigger
            aria-controls={directionListboxId}
            className="h-8 w-24 rounded [&[data-size]]:h-8"
          >
            <SelectValue />
          </SelectTrigger>
          <SelectContent
            className="min-w-[var(--radix-select-trigger-width)] origin-[var(--radix-select-content-transform-origin)]"
            id={directionListboxId}
          >
            {[
              { label: "Asc", value: "asc" as const },
              { label: "Desc", value: "desc" as const }
            ].map((order) => (
              <SelectItem
                key={order.value}
                value={order.value}
              >
                {order.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Button
          aria-controls={sortItemId}
          className="size-8 shrink-0 rounded"
          onClick={() => {
            onSortRemove(sort.id)
          }}
          size="icon"
          variant="ghost"
        >
          <TrashIcon />
        </Button>
        <SortableItemHandle asChild>
          <Button
            className="size-8 shrink-0 rounded"
            size="icon"
            variant="ghost"
          >
            <GripVerticalIcon />
          </Button>
        </SortableItemHandle>
      </li>
    </SortableItem>
  )
}
