import type { JSX } from 'react'
import type { Column } from '@tanstack/react-table'

import { useCallback, useId, useMemo } from 'react'

import { PlusSquareIcon, XCircleIcon } from '#components/custom/icons'
import { Button } from '#components/shadcn/button'
import { Input } from '#components/shadcn/input'
import { Label } from '#components/shadcn/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '#components/shadcn/popover'
import { Separator } from '#components/shadcn/separator'
import { Slider } from '#components/shadcn/slider'
import { cn } from '#lib/utils'

interface DataTableSliderFilterProps<TData> {
  column: Column<TData>
  title?: string
}

interface Range {
  max: number
  min: number
}

type RangeValue = [number, number]

export function DataTableSliderFilter<TData>({
  column,
  title
}: DataTableSliderFilterProps<TData>): JSX.Element {
  const id = useId()

  const columnFilterValue =
    getIsValidRange(column.getFilterValue()) ?
      (column.getFilterValue() as RangeValue)
    : undefined

  const defaultRange = column.columnDef.meta?.range
  const unit = column.columnDef.meta?.unit

  const { max, min, step } = useMemo<Range & { step: number }>(() => {
    let minValue = 0
    let maxValue = 100

    if (defaultRange && getIsValidRange(defaultRange)) {
      ;[minValue, maxValue] = defaultRange
    } else {
      const values = column.getFacetedMinMaxValues()
      if (
        values &&
        Array.isArray(values) &&
        (values as number[]).length === 2
      ) {
        const [facetMinValue, facetMaxValue] = values
        if (
          typeof facetMinValue === 'number' &&
          typeof facetMaxValue === 'number'
        ) {
          minValue = facetMinValue
          maxValue = facetMaxValue
        }
      }
    }

    const rangeSize = maxValue - minValue
    const step =
      rangeSize <= 20 ? 1
      : rangeSize <= 100 ? Math.ceil(rangeSize / 20)
      : Math.ceil(rangeSize / 50)

    return { max: maxValue, min: minValue, step }
  }, [column, defaultRange])

  const range = useMemo((): RangeValue => {
    return columnFilterValue ?? [min, max]
  }, [columnFilterValue, min, max])

  const formatValue = useCallback((value: number) => {
    return value.toLocaleString(undefined, { maximumFractionDigits: 0 })
  }, [])

  const onFromInputChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const numValue = Number(event.target.value)
      if (!Number.isNaN(numValue) && numValue >= min && numValue <= range[1]) {
        column.setFilterValue([numValue, range[1]])
      }
    },
    [column, min, range]
  )

  const onToInputChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const numValue = Number(event.target.value)
      if (!Number.isNaN(numValue) && numValue <= max && numValue >= range[0]) {
        column.setFilterValue([range[0], numValue])
      }
    },
    [column, max, range]
  )

  const onSliderValueChange = useCallback(
    (value: RangeValue) => {
      if (Array.isArray(value) && (value as number[]).length === 2) {
        column.setFilterValue(value)
      }
    },
    [column]
  )

  const onReset = useCallback(
    (event: React.MouseEvent) => {
      if (event.target instanceof HTMLDivElement) {
        event.stopPropagation()
      }
      column.setFilterValue(undefined)
    },
    [column]
  )

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          className='border-dashed'
          size='sm'
          variant='outline'
        >
          {columnFilterValue ?
            <div
              aria-label={`Clear ${title ?? ''} filter`}
              className='rounded-sm opacity-70 transition-opacity hover:opacity-100 focus-visible:ring-1 focus-visible:ring-ring focus-visible:outline-none'
              onClick={onReset}
              role='button'
              tabIndex={0}
            >
              <XCircleIcon />
            </div>
          : <PlusSquareIcon />}
          <span>{title}</span>
          {columnFilterValue ?
            <>
              <Separator
                className='mx-0.5 data-[orientation=vertical]:h-4'
                orientation='vertical'
              />
              {formatValue(columnFilterValue[0])} -{' '}
              {formatValue(columnFilterValue[1])}
              {unit ? ` ${unit}` : ''}
            </>
          : null}
        </Button>
      </PopoverTrigger>
      <PopoverContent
        align='start'
        className='flex w-auto flex-col gap-4'
      >
        <div className='flex flex-col gap-3'>
          <p className='leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70'>
            {title}
          </p>
          <div className='flex items-center gap-4'>
            <Label
              className='sr-only'
              htmlFor={`${id}-from`}
            >
              From
            </Label>
            <div className='relative'>
              <Input
                aria-valuemax={max}
                aria-valuemin={min}
                className={cn('h-8 w-24', unit && 'pr-8')}
                id={`${id}-from`}
                inputMode='numeric'
                max={max}
                min={min}
                onChange={onFromInputChange}
                pattern='[0-9]*'
                placeholder={min.toString()}
                type='number'
                value={range[0].toString()}
              />
              {unit && (
                <span className='absolute top-0 right-0 bottom-0 flex items-center rounded-r-md bg-accent px-2 text-sm text-muted-foreground'>
                  {unit}
                </span>
              )}
            </div>
            <Label
              className='sr-only'
              htmlFor={`${id}-to`}
            >
              to
            </Label>
            <div className='relative'>
              <Input
                aria-valuemax={max}
                aria-valuemin={min}
                className={cn('h-8 w-24', unit && 'pr-8')}
                id={`${id}-to`}
                inputMode='numeric'
                max={max}
                min={min}
                onChange={onToInputChange}
                pattern='[0-9]*'
                placeholder={max.toString()}
                type='number'
                value={range[1].toString()}
              />
              {unit && (
                <span className='absolute top-0 right-0 bottom-0 flex items-center rounded-r-md bg-accent px-2 text-sm text-muted-foreground'>
                  {unit}
                </span>
              )}
            </div>
          </div>
          <Label
            className='sr-only'
            htmlFor={`${id}-slider`}
          >
            {title} slider
          </Label>
          <Slider
            id={`${id}-slider`}
            max={max}
            min={min}
            onValueChange={onSliderValueChange}
            step={step}
            value={range}
          />
        </div>
        <Button
          aria-label={`Clear ${title ?? ''} filter`}
          onClick={onReset}
          size='sm'
          variant='outline'
        >
          Clear
        </Button>
      </PopoverContent>
    </Popover>
  )
}

function getIsValidRange(value: unknown): value is RangeValue {
  return (
    Array.isArray(value) &&
    value.length === 2 &&
    typeof value[0] === 'number' &&
    typeof value[1] === 'number'
  )
}
