import type {
  Announcements,
  DndContextProps,
  DragEndEvent,
  DraggableAttributes,
  DraggableSyntheticListeners,
  DragStartEvent,
  DropAnimation,
  ScreenReaderInstructions,
  UniqueIdentifier,
} from "@dnd-kit/core";
import {
  closestCenter,
  closestCorners,
  DndContext,
  DragOverlay,
  defaultDropAnimationSideEffects,
  KeyboardSensor,
  MouseSensor,
  TouchSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  restrictToHorizontalAxis,
  restrictToParentElement,
  restrictToVerticalAxis,
} from "@dnd-kit/modifiers";
import type { SortableContextProps } from "@dnd-kit/sortable";
import {
  arrayMove,
  horizontalListSortingStrategy,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Slot } from "@radix-ui/react-slot";
import type { JSX } from "react";
import {
  createContext,
  forwardRef,
  useCallback,
  useContext,
  useId,
  useLayoutEffect,
  useMemo,
  useState,
} from "react";
import * as ReactDOM from "react-dom";

import { useComposedRefs } from "#lib/compose-refs";
import { cn } from "#lib/utils";

const orientationConfig = {
  horizontal: {
    collisionDetection: closestCenter,
    modifiers: [restrictToHorizontalAxis, restrictToParentElement],
    strategy: horizontalListSortingStrategy,
  },
  mixed: {
    collisionDetection: closestCorners,
    modifiers: [restrictToParentElement],
    strategy: undefined,
  },
  vertical: {
    collisionDetection: closestCenter,
    modifiers: [restrictToVerticalAxis, restrictToParentElement],
    strategy: verticalListSortingStrategy,
  },
};

const ROOT_NAME = "Sortable";
const CONTENT_NAME = "SortableContent";
const ITEM_NAME = "SortableItem";
const ITEM_HANDLE_NAME = "SortableItemHandle";
const OVERLAY_NAME = "SortableOverlay";

interface SortableRootContextValue<T> {
  activeID: null | UniqueIdentifier;
  flatCursor: boolean;
  getItemValue: (item: T) => UniqueIdentifier;
  id: string;
  items: UniqueIdentifier[];
  modifiers: DndContextProps["modifiers"];
  setActiveID: (id: null | UniqueIdentifier) => void;
  strategy: SortableContextProps["strategy"];
}

const SortableRootContext =
  createContext<null | SortableRootContextValue<unknown>>(null);
SortableRootContext.displayName = ROOT_NAME;

interface GetItemValue<T> {
  /**
   * Callback that returns a unique identifier for each sortable item. Required for array of objects.
   * @example getItemValue={(item) => item.id}
   */
  getItemValue: (item: T) => UniqueIdentifier;
}

type SortableRootProps<T> = DndContextProps &
  (T extends object ? GetItemValue<T> : Partial<GetItemValue<T>>) & {
    flatCursor?: boolean;
    onMove?: (
      event: DragEndEvent & {
        activeIndex: number;
        overIndex: number;
      },
    ) => void;
    onValueChange?: (items: T[]) => void;
    orientation?: "horizontal" | "mixed" | "vertical";
    strategy?: SortableContextProps["strategy"];
    value: T[];
  };

function SortableRoot<T>(props: SortableRootProps<T>): JSX.Element {
  const {
    accessibility,
    collisionDetection,
    flatCursor = false,
    getItemValue: getItemValueProp,
    modifiers,
    onMove,
    onValueChange,
    orientation = "vertical",
    strategy,
    value,
    ...sortableProps
  } = props;

  const id = useId();
  const [activeID, setActiveID] = useState<null | UniqueIdentifier>(null);

  const sensors = useSensors(
    useSensor(MouseSensor),
    useSensor(TouchSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );
  const config = useMemo(() => orientationConfig[orientation], [orientation]);

  const getItemValue = (item: T): UniqueIdentifier => {
    if (typeof item === "object" && !getItemValueProp) {
      throw new Error("getItemValue is required when using array of objects");
    }
    return getItemValueProp
      ? getItemValueProp(item)
      : (item as UniqueIdentifier);
  };

  const items = value.map((item) => getItemValue(item));

  const onDragStart = (event: DragStartEvent) => {
    sortableProps.onDragStart?.(event);

    if (event.activatorEvent.defaultPrevented) return;

    setActiveID(event.active.id);
  };

  const onDragEnd = (event: DragEndEvent) => {
    sortableProps.onDragEnd?.(event);

    if (event.activatorEvent.defaultPrevented) return;

    const { active, over } = event;
    if (over && active.id !== over.id) {
      const activeIndex = value.findIndex(
        (item) => getItemValue(item) === active.id,
      );
      const overIndex = value.findIndex(
        (item) => getItemValue(item) === over.id,
      );

      if (onMove) {
        onMove({
          ...event,
          activeIndex,
          overIndex,
        });
      } else {
        onValueChange?.(arrayMove(value, activeIndex, overIndex));
      }
    }
    setActiveID(null);
  };

  const onDragCancel = useCallback(
    (event: DragEndEvent) => {
      sortableProps.onDragCancel?.(event);

      if (event.activatorEvent.defaultPrevented) return;

      setActiveID(null);
    },
    [sortableProps.onDragCancel],
  );

  const announcements: Announcements = useMemo(
    () => ({
      onDragCancel({ active }) {
        const activeIndex =
          (
            active.data.current as null | {
              sortable: {
                index: number;
              };
            }
          )?.sortable.index ?? 0;
        const activeValue = active.id.toString();
        return `Sorting cancelled. Sortable item "${activeValue}" returned to position ${(activeIndex + 1).toString()} of ${value.length.toString()}.`;
      },
      onDragEnd({ active, over }) {
        const activeValue = active.id.toString();
        if (over) {
          const overIndex =
            (
              over.data.current as null | {
                sortable: {
                  index: number;
                };
              }
            )?.sortable.index ?? 0;
          return `Sortable item "${activeValue}" dropped at position ${(overIndex + 1).toString()} of ${value.length.toString()}.`;
        }
        return `Sortable item "${activeValue}" dropped. No changes were made.`;
      },
      onDragMove({ active, over }) {
        if (over) {
          const overIndex =
            (
              over.data.current as null | {
                sortable?: {
                  index: number;
                };
              }
            )?.sortable?.index ?? 0;
          const activeIndex =
            (
              active.data.current as null | {
                sortable?: {
                  index: number;
                };
              }
            )?.sortable?.index ?? 0;
          const moveDirection = overIndex > activeIndex ? "down" : "up";
          const activeValue = active.id.toString();
          return `Sortable item "${activeValue}" is moving ${moveDirection} to position ${(overIndex + 1).toString()} of ${value.length.toString()}.`;
        }
        return "Sortable item is no longer over a droppable area. Press escape to cancel.";
      },
      onDragOver({ active, over }) {
        if (over) {
          const overIndex =
            (
              over.data.current as null | {
                sortable?: {
                  index: number;
                };
              }
            )?.sortable?.index ?? 0;
          const activeIndex =
            (
              active.data.current as null | {
                sortable?: {
                  index: number;
                };
              }
            )?.sortable?.index ?? 0;
          const moveDirection = overIndex > activeIndex ? "down" : "up";
          const activeValue = active.id.toString();
          return `Sortable item "${activeValue}" moved ${moveDirection} to position ${(overIndex + 1).toString()} of ${value.length.toString()}.`;
        }
        return "Sortable item is no longer over a droppable area. Press escape to cancel.";
      },
      onDragStart({ active }) {
        const activeValue = active.id.toString();
        const sortableData = active.data.current as null | {
          sortable?: {
            index: number;
          };
        };
        const currentIndex = sortableData?.sortable?.index ?? 0;
        return `Grabbed sortable item "${activeValue}". Current position is ${(currentIndex + 1).toString()} of ${value.length.toString()}. Use arrow keys to move, space to drop.`;
      },
    }),
    [value],
  );

  const screenReaderInstructions: ScreenReaderInstructions = useMemo(
    () => ({
      draggable: `
        To pick up a sortable item, press space or enter.
        While dragging, use the ${
          orientation === "vertical"
            ? "up and down"
            : orientation === "horizontal"
              ? "left and right"
              : "arrow"
        } keys to move the item.
        Press space or enter again to drop the item in its new position, or press escape to cancel.
      `,
    }),
    [orientation],
  );

  const contextValue = {
    activeID,
    flatCursor,
    getItemValue,
    id,
    items,
    modifiers: modifiers ?? config.modifiers,
    setActiveID,
    strategy: strategy ?? config.strategy,
  };

  return (
    <SortableRootContext.Provider
      value={contextValue as SortableRootContextValue<unknown> | null}
    >
      {/* @ts-ignore FIXME */}
      <DndContext
        collisionDetection={collisionDetection ?? config.collisionDetection}
        modifiers={modifiers ?? config.modifiers}
        sensors={sensors}
        {...sortableProps}
        accessibility={{
          announcements,
          screenReaderInstructions,
          ...accessibility,
        }}
        id={id}
        onDragCancel={onDragCancel}
        onDragEnd={onDragEnd}
        onDragStart={onDragStart}
      />
    </SortableRootContext.Provider>
  );
}

function useSortableContext(consumerName: string) {
  const context = useContext(SortableRootContext);
  if (!context) {
    throw new Error(`\`${consumerName}\` must be used within \`${ROOT_NAME}\``);
  }
  return context;
}

const SortableContentContext = createContext<boolean>(false);
SortableContentContext.displayName = CONTENT_NAME;

interface SortableContentProps extends React.ComponentPropsWithoutRef<"div"> {
  asChild?: boolean;
  children: React.ReactNode;
  strategy?: SortableContextProps["strategy"];
  withoutSlot?: boolean;
}

const SortableContent = forwardRef<HTMLDivElement, SortableContentProps>(
  (props, forwardedRef) => {
    const {
      asChild,
      children,
      strategy: strategyProp,
      withoutSlot,
      ...contentProps
    } = props;

    const context = useSortableContext(CONTENT_NAME);

    const ContentPrimitive = asChild ? Slot : "div";

    return (
      <SortableContentContext.Provider value={true}>
        {/* @ts-ignore FIXME */}
        <SortableContext
          items={context.items}
          strategy={strategyProp ?? context.strategy}
        >
          {withoutSlot ? (
            children
          ) : (
            <ContentPrimitive
              data-slot="sortable-content"
              {...contentProps}
              ref={forwardedRef}
            >
              {children}
            </ContentPrimitive>
          )}
        </SortableContext>
      </SortableContentContext.Provider>
    );
  },
);
SortableContent.displayName = CONTENT_NAME;

interface SortableItemContextValue {
  attributes: DraggableAttributes;
  disabled?: boolean;
  id: string;
  isDragging?: boolean;
  listeners: DraggableSyntheticListeners | undefined;
  setActivatorNodeRef: (node: HTMLElement | null) => void;
}

const SortableItemContext = createContext<null | SortableItemContextValue>(
  null,
);
SortableItemContext.displayName = ITEM_NAME;

interface SortableItemProps extends React.ComponentPropsWithoutRef<"div"> {
  asChild?: boolean;
  asHandle?: boolean;
  disabled?: boolean;
  value: UniqueIdentifier;
}

function useSortableItemContext(consumerName: string) {
  const context = useContext(SortableItemContext);
  if (!context) {
    throw new Error(`\`${consumerName}\` must be used within \`${ITEM_NAME}\``);
  }
  return context;
}

const SortableItem = forwardRef<HTMLDivElement, SortableItemProps>(
  (props, forwardedRef) => {
    const {
      asChild,
      asHandle,
      className,
      disabled,
      style,
      value,
      ...itemProps
    } = props;

    const inSortableContent = useContext(SortableContentContext);
    const inSortableOverlay = useContext(SortableOverlayContext);

    if (!inSortableContent && !inSortableOverlay) {
      throw new Error(
        `\`${ITEM_NAME}\` must be used within \`${CONTENT_NAME}\` or \`${OVERLAY_NAME}\``,
      );
    }

    if (value === "") {
      throw new Error(`\`${ITEM_NAME}\` value cannot be an empty string`);
    }

    const context = useSortableContext(ITEM_NAME);
    const id = useId();
    const {
      attributes,
      isDragging,
      listeners,
      setActivatorNodeRef,
      setNodeRef,
      transform,
      transition,
      // @ts-expect-error FIXME
    } = useSortable({
      disabled,
      id: value,
    });

    const composedRef = useComposedRefs(forwardedRef, (node) => {
      if (disabled) return;
      setNodeRef(node);
      if (asHandle) setActivatorNodeRef(node);
    });

    const composedStyle = useMemo<React.CSSProperties>(() => {
      return {
        transform: CSS.Translate.toString(transform),
        transition,
        ...style,
      };
    }, [transform, transition, style]);

    const itemContext = useMemo<SortableItemContextValue>(
      // @ts-expect-error FIXME
      () => ({
        attributes,
        disabled,
        id,
        isDragging,
        listeners,
        setActivatorNodeRef,
      }),
      [id, attributes, listeners, setActivatorNodeRef, isDragging, disabled],
    );

    const ItemPrimitive = asChild ? Slot : "div";

    return (
      <SortableItemContext.Provider value={itemContext}>
        <ItemPrimitive
          data-disabled={disabled}
          data-dragging={isDragging ? "" : undefined}
          data-slot="sortable-item"
          id={id}
          {...itemProps}
          {...(asHandle && !disabled ? attributes : {})}
          {...(asHandle && !disabled ? listeners : {})}
          className={cn(
            "focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1",
            {
              "cursor-default": context.flatCursor,
              "cursor-grab": !isDragging && asHandle && !context.flatCursor,
              "data-dragging:cursor-grabbing": !context.flatCursor,
              "opacity-50": isDragging,
              "pointer-events-none opacity-50": disabled,
              "touch-none select-none": asHandle,
            },
            className,
          )}
          ref={composedRef}
          style={composedStyle}
        />
      </SortableItemContext.Provider>
    );
  },
);
SortableItem.displayName = ITEM_NAME;

interface SortableItemHandleProps
  extends React.ComponentPropsWithoutRef<"button"> {
  asChild?: boolean;
}

const SortableItemHandle = forwardRef<
  HTMLButtonElement,
  SortableItemHandleProps
>((props, forwardedRef) => {
  const { asChild, className, disabled, ...itemHandleProps } = props;

  const context = useSortableContext(ITEM_HANDLE_NAME);
  const itemContext = useSortableItemContext(ITEM_HANDLE_NAME);

  const isDisabled = disabled ?? itemContext.disabled;

  const composedRef = useComposedRefs(forwardedRef, (node) => {
    if (!isDisabled) return;
    itemContext.setActivatorNodeRef(node);
  });

  const HandlePrimitive = asChild ? Slot : "button";

  return (
    <HandlePrimitive
      aria-controls={itemContext.id}
      data-disabled={isDisabled}
      data-dragging={itemContext.isDragging ? "" : undefined}
      data-slot="sortable-item-handle"
      type="button"
      {...itemHandleProps}
      {...(isDisabled ? {} : itemContext.attributes)}
      {...(isDisabled ? {} : itemContext.listeners)}
      className={cn(
        "select-none disabled:pointer-events-none disabled:opacity-50",
        context.flatCursor
          ? "cursor-default"
          : "cursor-grab data-dragging:cursor-grabbing",
        className,
      )}
      disabled={isDisabled}
      ref={composedRef}
    />
  );
});
SortableItemHandle.displayName = ITEM_HANDLE_NAME;

const SortableOverlayContext = createContext(false);
SortableOverlayContext.displayName = OVERLAY_NAME;

const dropAnimation: DropAnimation = {
  sideEffects: defaultDropAnimationSideEffects({
    styles: {
      active: {
        opacity: "0.4",
      },
    },
  }),
};

interface SortableOverlayProps
  extends Omit<React.ComponentPropsWithoutRef<typeof DragOverlay>, "children"> {
  children?:
    | ((params: { value: UniqueIdentifier }) => React.ReactNode)
    | React.ReactNode;
  container?: DocumentFragment | Element | null;
}

function SortableOverlay(props: SortableOverlayProps): JSX.Element | null {
  const { children, container: containerProp, ...overlayProps } = props;

  const context = useSortableContext(OVERLAY_NAME);

  const [mounted, setMounted] = useState(false);
  useLayoutEffect(() => {
    setMounted(true);
  }, []);

  const container =
    containerProp ?? (mounted ? globalThis.document.body : null);

  if (!container) return null;

  return ReactDOM.createPortal(
    // @ts-expect-error FIXME
    <DragOverlay
      className={cn(!context.flatCursor && "cursor-grabbing")}
      dropAnimation={dropAnimation}
      modifiers={context.modifiers}
      {...overlayProps}
    >
      <SortableOverlayContext.Provider value={true}>
        {context.activeID
          ? typeof children === "function"
            ? children({
                value: context.activeID,
              })
            : children
          : null}
      </SortableOverlayContext.Provider>
    </DragOverlay>,
    container,
  );
}

export {
  SortableContent as Content,
  SortableItem as Item,
  SortableItemHandle as ItemHandle,
  SortableOverlay as Overlay,
  //
  SortableRoot as Root,
  SortableRoot as Sortable,
  SortableContent,
  SortableItem,
  SortableItemHandle,
  SortableOverlay,
};
