import type { LucideIcon } from "@archesai/ui";
import { PureDataSelector } from "@archesai/ui";
import type { BaseEntity } from "@archesai/ui/types/entities";
import type { UseSuspenseQueryOptions } from "@tanstack/react-query";
import { useSuspenseQuery } from "@tanstack/react-query";
import type { JSX } from "react";

export interface DataSelectorContainerProps<TItem extends BaseEntity> {
  // Query configuration
  queryOptions: UseSuspenseQueryOptions<{
    data: TItem[];
  }>;

  // Selection state
  selectedData: TItem | TItem[] | undefined;
  setSelectedData: (data: TItem | TItem[] | undefined) => void;

  // Config
  itemType: string;
  isMultiSelect?: boolean;
  icons?: {
    color: string;
    Icon: LucideIcon;
    name: string;
  }[];
  getItemDetails?: (item: TItem) => React.ReactNode;

  // Optional
  placeholder?: string;
  className?: string;
}

/**
 * Container component that handles data fetching for DataSelector.
 * Fetches data and passes it to PureDataSelector.
 */
export function DataSelectorContainer<TItem extends BaseEntity>({
  queryOptions,
  selectedData,
  setSelectedData,
  itemType,
  isMultiSelect = false,
  icons,
  getItemDetails,
  placeholder,
  className,
}: DataSelectorContainerProps<TItem>): JSX.Element {
  // Fetch data using the provided query options
  const {
    data: { data: items },
  } = useSuspenseQuery(queryOptions);

  return (
    <PureDataSelector<TItem>
      className={className || ""}
      getItemDetails={getItemDetails || ((item: TItem) => item.id)}
      icons={icons || []}
      isMultiSelect={isMultiSelect}
      items={items}
      itemType={itemType}
      onSelect={setSelectedData}
      placeholder={placeholder || "Select an item"}
      selectedData={selectedData}
    />
  );
}
