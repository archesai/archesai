import type { JSX } from "react";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";

import { CheckCircle2Icon, ChevronsUpDownIcon } from "#components/custom/icons";
import { Badge } from "#components/shadcn/badge";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "#components/shadcn/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "#components/shadcn/popover";
import { cn } from "#lib/utils";

interface FacetedContextValue<Multiple extends boolean = boolean> {
  multiple?: Multiple;
  onItemSelect?: (value: string) => void;
  value?: FacetedValue<Multiple> | undefined;
}

type FacetedValue<Multiple extends boolean> = Multiple extends true
  ? string[]
  : string;

const FacetedContext = createContext<FacetedContextValue | null>(null);

interface FacetedBadgeListProps extends React.ComponentProps<"div"> {
  badgeClassName?: string;
  max?: number;
  options?: { label: string; value: string }[];
  placeholder?: string;
}

interface FacetedProps<Multiple extends boolean = false>
  extends React.ComponentProps<typeof Popover> {
  children?: React.ReactNode;
  multiple?: Multiple;
  onValueChange?: (value: FacetedValue<Multiple> | undefined) => void;
  value?: FacetedValue<Multiple>;
}

function Faceted<Multiple extends boolean = false>(
  props: FacetedProps<Multiple>,
): JSX.Element {
  const {
    children,
    multiple = false,
    onOpenChange: onOpenChangeProp,
    onValueChange,
    open: openProp,
    value,
    ...facetedProps
  } = props;

  const [uncontrolledOpen, setUncontrolledOpen] = useState(false);
  const isControlled = openProp !== undefined;
  const open = isControlled ? openProp : uncontrolledOpen;

  const onOpenChange = useCallback(
    (newOpen: boolean) => {
      if (!isControlled) {
        setUncontrolledOpen(newOpen);
      }
      onOpenChangeProp?.(newOpen);
    },
    [isControlled, onOpenChangeProp],
  );

  const onItemSelect = useCallback(
    (selectedValue: string) => {
      if (!onValueChange) return;

      if (multiple) {
        const currentValue: unknown[] = Array.isArray(value) ? value : [];
        const newValue = currentValue.includes(selectedValue)
          ? currentValue.filter((v) => v !== selectedValue)
          : [...currentValue, selectedValue];
        onValueChange(newValue as FacetedValue<Multiple>);
      } else {
        if (value === selectedValue) {
          onValueChange(undefined);
        } else {
          onValueChange(selectedValue as FacetedValue<Multiple>);
        }

        requestAnimationFrame(() => {
          onOpenChange(false);
        });
      }
    },
    [multiple, value, onValueChange, onOpenChange],
  );

  const contextValue = useMemo<FacetedContextValue<typeof multiple>>(
    () => ({ multiple, onItemSelect, value }),
    [value, onItemSelect, multiple],
  );

  return (
    <FacetedContext.Provider value={contextValue}>
      <Popover
        onOpenChange={onOpenChange}
        open={open}
        {...facetedProps}
      >
        {children}
      </Popover>
    </FacetedContext.Provider>
  );
}

function FacetedBadgeList(props: FacetedBadgeListProps): JSX.Element {
  const {
    badgeClassName,
    className,
    max = 2,
    options = [],
    placeholder = "Select options...",
    ...badgeListProps
  } = props;

  const context = useFacetedContext("FacetedBadgeList");
  const values = Array.isArray(context.value)
    ? (context.value as string[] | undefined)
    : ([context.value].filter(Boolean) as string[]);

  const getLabel = useCallback(
    (value: string) => {
      const option = options.find((opt) => opt.value === value);
      return option?.label ?? value;
    },
    [options],
  );

  if (!values || values.length === 0) {
    return (
      <div
        {...badgeListProps}
        className="flex w-full items-center gap-1 text-muted-foreground"
      >
        {placeholder}
        <ChevronsUpDownIcon className="ml-auto size-4 shrink-0 opacity-50" />
      </div>
    );
  }

  return (
    <div
      {...badgeListProps}
      className={cn("flex flex-wrap items-center gap-1", className)}
    >
      {values.length > max ? (
        <Badge
          className={cn("rounded-sm px-1 font-normal", badgeClassName)}
          variant="secondary"
        >
          {values.length} selected
        </Badge>
      ) : (
        values.map((value) => (
          <Badge
            className={cn("rounded-sm px-1 font-normal", badgeClassName)}
            key={value}
            variant="secondary"
          >
            <span className="truncate">{getLabel(value)}</span>
          </Badge>
        ))
      )}
    </div>
  );
}

function FacetedContent(
  props: React.ComponentProps<typeof PopoverContent>,
): JSX.Element {
  const { children, className, ...contentProps } = props;

  return (
    <PopoverContent
      {...contentProps}
      align="start"
      className={cn(
        "w-[200px] origin-(--radix-popover-content-transform-origin) p-0",
        className,
      )}
    >
      <Command>{children}</Command>
    </PopoverContent>
  );
}

function FacetedTrigger(
  props: React.ComponentProps<typeof PopoverTrigger>,
): JSX.Element {
  const { children, className, ...triggerProps } = props;

  return (
    <PopoverTrigger
      {...triggerProps}
      className={cn("justify-between text-left", className)}
    >
      {children}
    </PopoverTrigger>
  );
}

function useFacetedContext(name: string) {
  const context = useContext(FacetedContext);
  if (!context) {
    throw new Error(`\`${name}\` must be within Faceted`);
  }
  return context;
}

const FacetedInput = CommandInput;

const FacetedList = CommandList;

const FacetedEmpty = CommandEmpty;

const FacetedGroup = CommandGroup;

interface FacetedItemProps extends React.ComponentProps<typeof CommandItem> {
  value: string;
}

function FacetedItem(props: FacetedItemProps): JSX.Element {
  const { children, className, onSelect, value, ...itemProps } = props;
  const context = useFacetedContext("FacetedItem");

  const isSelected = context.multiple
    ? Array.isArray(context.value) && context.value.includes(value)
    : context.value === value;

  const onItemSelect = useCallback(
    (currentValue: string) => {
      if (onSelect) {
        onSelect(currentValue);
      } else if (context.onItemSelect) {
        context.onItemSelect(currentValue);
      }
    },
    [onSelect, context],
  );

  return (
    <CommandItem
      aria-selected={isSelected}
      className={cn("gap-2", className)}
      data-selected={isSelected}
      onSelect={() => {
        onItemSelect(value);
      }}
      {...itemProps}
    >
      <span
        className={cn(
          "flex size-4 items-center justify-center rounded-sm border border-primary",
          isSelected
            ? "bg-primary text-primary-foreground"
            : "opacity-50 [&_svg]:invisible",
        )}
      >
        <CheckCircle2Icon className="size-4" />
      </span>
      {children}
    </CommandItem>
  );
}

const FacetedSeparator = CommandSeparator;

export {
  Faceted,
  FacetedBadgeList,
  FacetedContent,
  FacetedEmpty,
  FacetedGroup,
  FacetedInput,
  FacetedItem,
  FacetedList,
  FacetedSeparator,
  FacetedTrigger,
};
