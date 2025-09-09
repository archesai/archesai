import type { Column } from "@tanstack/react-table";
import type { JSX } from "react";

import { useCallback, useMemo } from "react";
import { Input } from "#components/shadcn/input";
import { cn } from "#lib/utils";
import type { BaseEntity, FilterCondition } from "#types/entities";

interface DataTableRangeFilterProps<TData extends BaseEntity>
  extends React.ComponentProps<"div"> {
  column: Column<TData>;
  filter: FilterCondition & { id: string };
  inputId: string;
  onFilterUpdate: (
    filterId: string,
    updates: Partial<Omit<FilterCondition, "type">>,
  ) => void;
}

export function DataTableRangeFilter<TData extends BaseEntity>({
  className,
  column,
  filter,
  inputId,
  onFilterUpdate,
  ...props
}: DataTableRangeFilterProps<TData>): JSX.Element {
  const meta = column.columnDef.meta;

  const [min, max] = useMemo(() => {
    const range = column.columnDef.meta?.range;
    if (range) return range;

    const values = column.getFacetedMinMaxValues();
    if (!values) return [0, 100];

    return [values[0], values[1]];
  }, [column]);

  const formatValue = useCallback(
    (value: boolean | null | number | string | undefined) => {
      if (value === undefined || value === null || value === "") return "";
      const numValue = Number(value);
      return Number.isNaN(numValue)
        ? ""
        : numValue.toLocaleString(undefined, {
            maximumFractionDigits: 0,
          });
    },
    [],
  );

  const value = useMemo(() => {
    if (Array.isArray(filter.value)) return filter.value.map(formatValue);
    return [formatValue(filter.value as number | string | undefined), ""];
  }, [filter.value, formatValue]);

  const onRangeValueChange = useCallback(
    (value: string, isMin?: boolean) => {
      const numValue = Number(value);
      const currentValues = Array.isArray(filter.value)
        ? filter.value
        : ["", ""];
      const otherValue = isMin
        ? String(currentValues[1] ?? "")
        : String(currentValues[0] ?? "");

      if (
        value === "" ||
        (!Number.isNaN(numValue) &&
          (isMin
            ? numValue >= min && numValue <= (Number(otherValue) || max)
            : numValue <= max && numValue >= (Number(otherValue) || min)))
      ) {
        onFilterUpdate(filter.id, {
          value: isMin ? [value, otherValue] : [otherValue, value],
        });
      }
    },
    [filter.id, filter.value, min, max, onFilterUpdate],
  );

  return (
    <div
      className={cn("flex w-full items-center gap-2", className)}
      data-slot="range"
      {...props}
    >
      <Input
        aria-label={`${meta?.label ?? ""} minimum value`}
        aria-valuemax={max}
        aria-valuemin={min}
        className="h-8 w-full rounded"
        data-slot="range-min"
        defaultValue={value[0]}
        id={`${inputId}-min`}
        inputMode="numeric"
        max={max}
        min={min}
        onChange={(event) => {
          onRangeValueChange(event.target.value, true);
        }}
        placeholder={min.toString()}
        type="number"
      />
      <span className="sr-only shrink-0 text-muted-foreground">to</span>
      <Input
        aria-label={`${meta?.label ?? ""} maximum value`}
        aria-valuemax={max}
        aria-valuemin={min}
        className="h-8 w-full rounded"
        data-slot="range-max"
        defaultValue={value[1]}
        id={`${inputId}-max`}
        inputMode="numeric"
        max={max}
        min={min}
        onChange={(event) => {
          onRangeValueChange(event.target.value);
        }}
        placeholder={max.toString()}
        type="number"
      />
    </div>
  );
}
