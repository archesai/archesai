import type { UseSuspenseQueryOptions } from "@tanstack/react-query";
import { useSuspenseQuery } from "@tanstack/react-query";
import type { JSX } from "react";
import { useState } from "react";

import type { LucideIcon } from "#components/custom/icons";
import {
  CheckCircle2Icon,
  PlusSquareIcon,
  SortAscIcon,
} from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList,
} from "#components/shadcn/command";
import { HoverCard, HoverCardContent } from "#components/shadcn/hover-card";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "#components/shadcn/popover";
import { cn } from "#lib/utils";
import type { BaseEntity } from "#types/entities";

interface DataSelectorProps<TItem extends BaseEntity> {
  getItemDetails?: (item: TItem) => React.ReactNode;
  icons?: { color: string; Icon: LucideIcon; name: string }[];
  isMultiSelect?: boolean;
  itemType: string;
  organizationId?: string;
  selectedData: TItem | TItem[] | undefined;
  setSelectedData: (data: TItem | TItem[] | undefined) => void;
  useFindMany: UseSuspenseQueryOptions<{
    data: TItem[];
  }>;
}

export function DataSelector<TItem extends BaseEntity>({
  getItemDetails,
  icons,
  isMultiSelect = false,
  itemType,
  selectedData,
  setSelectedData,
  useFindMany,
}: DataSelectorProps<TItem>): JSX.Element {
  const {
    data: { data },
  } = useSuspenseQuery(useFindMany);
  const [open, setOpen] = useState(false);
  const [hoveredItem, setHoveredItem] = useState<TItem | undefined>();
  const [searchTerm, setSearchTerm] = useState("");

  // item.name.toLowerCase().includes(searchTerm.toLowerCase())

  // Handler for selecting/deselecting items
  const handleSelect = (item: TItem) => {
    if (isMultiSelect) {
      const selectedArray = Array.isArray(selectedData)
        ? selectedData
        : selectedData
          ? [selectedData]
          : [];
      const isSelected = selectedArray.some((i) => i.id === item.id);
      if (isSelected) {
        setSelectedData(selectedArray.filter((i) => i.id !== item.id));
      } else {
        setSelectedData([...selectedArray, item]);
      }
    } else {
      setSelectedData(item);
      setOpen(false); // Close popover on single select
    }
  };

  // Handler for removing selected item (for multi-select)
  const handleRemove = (item: TItem) => {
    if (isMultiSelect && Array.isArray(selectedData)) {
      setSelectedData(selectedData.filter((i) => i.id !== item.id));
    } else {
      setSelectedData(undefined);
    }
  };

  return (
    <div className="flex flex-col gap-2">
      {/* Popover for Selection */}
      <Popover
        onOpenChange={setOpen}
        open={open}
      >
        <PopoverTrigger asChild>
          <Button
            aria-expanded={open}
            aria-label={`Select ${itemType.toLowerCase()}`}
            className="w-full justify-between"
            role="combobox"
            variant="outline"
          >
            <div className="flex items-center gap-2">
              <div
                className={cn(
                  "flex flex-wrap gap-2",
                  (Array.isArray(selectedData) && selectedData.length > 0) ||
                    !isMultiSelect
                    ? ""
                    : "text-muted-foreground",
                )}
              >
                {isMultiSelect ? (
                  Array.isArray(selectedData) && selectedData.length > 0 ? (
                    `${selectedData.length.toString()} selected`
                  ) : (
                    `Select ${itemType.toLowerCase()}s...`
                  )
                ) : (
                  <div className="flex items-center gap-1">
                    {icons
                      ?.filter((x) => x.name === (selectedData as TItem).id)
                      .map((x) => {
                        const iconColor = x.color;
                        return (
                          <x.Icon
                            className={cn(
                              "mx-auto h-4 w-4",
                              iconColor.startsWith("text-") ? iconColor : "",
                            )}
                            key={`${x.name}-${x.color}`}
                            style={{
                              ...(iconColor.startsWith("#")
                                ? { color: iconColor }
                                : {}),
                            }}
                          />
                        );
                      })}
                    {(selectedData as TItem).id}
                  </div>
                )}
              </div>
            </div>
            {isMultiSelect ? (
              <PlusSquareIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            ) : (
              <SortAscIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[300px] p-0">
          <Command>
            <CommandInput
              onInput={(e) => {
                setSearchTerm(e.currentTarget.value);
              }}
              placeholder={`Search ${itemType.toLowerCase()}...`}
              value={searchTerm}
            />
            <CommandList className="max-h-[400px] overflow-y-auto">
              <CommandEmpty>No {itemType.toLowerCase()} found.</CommandEmpty>
              {data.map((item) => (
                <CommandItem
                  className={cn("flex items-center justify-between")}
                  key={item.id}
                  onMouseEnter={() => {
                    setHoveredItem(item);
                  }}
                  onMouseLeave={() => {
                    setHoveredItem(undefined);
                  }}
                  onSelect={() => {
                    handleSelect(item);
                  }}
                >
                  <div className="flex items-center gap-1">
                    {/* Icon Rendering */}
                    {icons
                      ?.filter((x) => x.name === item.id)
                      .map((x) => {
                        const iconColor = x.color;
                        return (
                          <x.Icon
                            className={cn(
                              "mx-auto h-4 w-4",
                              iconColor.startsWith("text-") ? iconColor : "",
                            )}
                            key={`${item.id}-${x.name}-${x.color}`}
                            style={{
                              ...(iconColor.startsWith("#")
                                ? { color: iconColor }
                                : {}),
                            }}
                          />
                        );
                      })}
                    <p>{item.id}</p>
                  </div>
                  {/* Check Icon for Selected Items */}
                  {isMultiSelect &&
                    Array.isArray(selectedData) &&
                    selectedData.some((i) => i.id === item.id) && (
                      <CheckCircle2Icon className="h-4 w-4 text-green-500" />
                    )}
                </CommandItem>
              ))}
            </CommandList>
          </Command>
          {/* HoverCard for Item Details */}
          {hoveredItem && getItemDetails && (
            <HoverCard open={true}>
              <HoverCardContent
                align="end"
                className="w-[250px] p-4"
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
          <div className="flex flex-wrap gap-2">
            {selectedData.map((item) => (
              <span
                className="inline-flex items-center rounded-xs bg-blue-100 px-2 py-1 text-sm text-blue-700"
                key={item.id}
              >
                {item.id}
                <button
                  className="ml-1 text-red-500 hover:text-red-700"
                  onClick={() => {
                    handleRemove(item);
                  }}
                  type="button"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
    </div>
  );
}
