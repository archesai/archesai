import type {
  Announcements,
  DndContextProps,
  DragEndEvent,
  DraggableAttributes,
  DraggableSyntheticListeners,
  DragStartEvent,
  DropAnimation,
  ScreenReaderInstructions,
  UniqueIdentifier
} from '@dnd-kit/core'
import type { SortableContextProps } from '@dnd-kit/sortable'

import * as React from 'react'
import {
  closestCenter,
  closestCorners,
  defaultDropAnimationSideEffects,
  DndContext,
  DragOverlay,
  KeyboardSensor,
  MouseSensor,
  TouchSensor,
  useSensor,
  useSensors
} from '@dnd-kit/core'
import {
  restrictToHorizontalAxis,
  restrictToParentElement,
  restrictToVerticalAxis
} from '@dnd-kit/modifiers'
import {
  arrayMove,
  horizontalListSortingStrategy,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Slot } from '@radix-ui/react-slot'
import * as ReactDOM from 'react-dom'

import { useComposedRefs } from '#lib/compose-refs'
import { cn } from '#lib/utils'

const orientationConfig = {
  horizontal: {
    collisionDetection: closestCenter,
    modifiers: [restrictToHorizontalAxis, restrictToParentElement],
    strategy: horizontalListSortingStrategy
  },
  mixed: {
    collisionDetection: closestCorners,
    modifiers: [restrictToParentElement],
    strategy: undefined
  },
  vertical: {
    collisionDetection: closestCenter,
    modifiers: [restrictToVerticalAxis, restrictToParentElement],
    strategy: verticalListSortingStrategy
  }
}

const ROOT_NAME = 'Sortable'
const CONTENT_NAME = 'SortableContent'
const ITEM_NAME = 'SortableItem'
const ITEM_HANDLE_NAME = 'SortableItemHandle'
const OVERLAY_NAME = 'SortableOverlay'

interface SortableRootContextValue<T> {
  activeId: null | UniqueIdentifier
  flatCursor: boolean
  getItemValue: (item: T) => UniqueIdentifier
  id: string
  items: UniqueIdentifier[]
  modifiers: DndContextProps['modifiers']
  setActiveId: (id: null | UniqueIdentifier) => void
  strategy: SortableContextProps['strategy']
}

const SortableRootContext =
  React.createContext<null | SortableRootContextValue<unknown>>(null)
SortableRootContext.displayName = ROOT_NAME

interface GetItemValue<T> {
  /**
   * Callback that returns a unique identifier for each sortable item. Required for array of objects.
   * @example getItemValue={(item) => item.id}
   */
  getItemValue: (item: T) => UniqueIdentifier
}

type SortableRootProps<T> = DndContextProps &
  (T extends object ? GetItemValue<T> : Partial<GetItemValue<T>>) & {
    flatCursor?: boolean
    onMove?: (
      event: DragEndEvent & { activeIndex: number; overIndex: number }
    ) => void
    onValueChange?: (items: T[]) => void
    orientation?: 'horizontal' | 'mixed' | 'vertical'
    strategy?: SortableContextProps['strategy']
    value: T[]
  }

function SortableRoot<T>(props: SortableRootProps<T>) {
  const {
    accessibility,
    collisionDetection,
    flatCursor = false,
    getItemValue: getItemValueProp,
    modifiers,
    onMove,
    onValueChange,
    orientation = 'vertical',
    strategy,
    value,
    ...sortableProps
  } = props

  const id = React.useId()
  const [activeId, setActiveId] = React.useState<null | UniqueIdentifier>(null)

  const sensors = useSensors(
    useSensor(MouseSensor),
    useSensor(TouchSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates
    })
  )
  const config = React.useMemo(
    () => orientationConfig[orientation],
    [orientation]
  )

  const getItemValue = React.useCallback(
    (item: T): UniqueIdentifier => {
      if (typeof item === 'object' && !getItemValueProp) {
        throw new Error('getItemValue is required when using array of objects')
      }
      return getItemValueProp ?
          getItemValueProp(item)
        : (item as UniqueIdentifier)
    },
    [getItemValueProp]
  )

  const items = React.useMemo(() => {
    return value.map((item) => getItemValue(item))
  }, [value, getItemValue])

  const onDragStart = React.useCallback(
    (event: DragStartEvent) => {
      sortableProps.onDragStart?.(event)

      if (event.activatorEvent.defaultPrevented) return

      setActiveId(event.active.id)
    },
    [sortableProps.onDragStart]
  )

  const onDragEnd = React.useCallback(
    (event: DragEndEvent) => {
      sortableProps.onDragEnd?.(event)

      if (event.activatorEvent.defaultPrevented) return

      const { active, over } = event
      if (over && active.id !== over.id) {
        const activeIndex = value.findIndex(
          (item) => getItemValue(item) === active.id
        )
        const overIndex = value.findIndex(
          (item) => getItemValue(item) === over.id
        )

        if (onMove) {
          onMove({ ...event, activeIndex, overIndex })
        } else {
          onValueChange?.(arrayMove(value, activeIndex, overIndex))
        }
      }
      setActiveId(null)
    },
    [value, onValueChange, onMove, getItemValue, sortableProps.onDragEnd]
  )

  const onDragCancel = React.useCallback(
    (event: DragEndEvent) => {
      sortableProps.onDragCancel?.(event)

      if (event.activatorEvent.defaultPrevented) return

      setActiveId(null)
    },
    [sortableProps.onDragCancel]
  )

  const announcements: Announcements = React.useMemo(
    () => ({
      onDragCancel({ active }) {
        const activeIndex =
          (active.data.current as null | { sortable: { index: number } })
            ?.sortable.index ?? 0
        const activeValue = active.id.toString()
        return `Sorting cancelled. Sortable item "${activeValue}" returned to position ${(activeIndex + 1).toString()} of ${value.length.toString()}.`
      },
      onDragEnd({ active, over }) {
        const activeValue = active.id.toString()
        if (over) {
          const overIndex =
            (over.data.current as null | { sortable: { index: number } })
              ?.sortable.index ?? 0
          return `Sortable item "${activeValue}" dropped at position ${(overIndex + 1).toString()} of ${value.length.toString()}.`
        }
        return `Sortable item "${activeValue}" dropped. No changes were made.`
      },
      onDragMove({ active, over }) {
        if (over) {
          const overIndex =
            (over.data.current as null | { sortable?: { index: number } })
              ?.sortable?.index ?? 0
          const activeIndex =
            (active.data.current as null | { sortable?: { index: number } })
              ?.sortable?.index ?? 0
          const moveDirection = overIndex > activeIndex ? 'down' : 'up'
          const activeValue = active.id.toString()
          return `Sortable item "${activeValue}" is moving ${moveDirection} to position ${(overIndex + 1).toString()} of ${value.length.toString()}.`
        }
        return 'Sortable item is no longer over a droppable area. Press escape to cancel.'
      },
      onDragOver({ active, over }) {
        if (over) {
          const overIndex =
            (over.data.current as null | { sortable?: { index: number } })
              ?.sortable?.index ?? 0
          const activeIndex =
            (active.data.current as null | { sortable?: { index: number } })
              ?.sortable?.index ?? 0
          const moveDirection = overIndex > activeIndex ? 'down' : 'up'
          const activeValue = active.id.toString()
          return `Sortable item "${activeValue}" moved ${moveDirection} to position ${(overIndex + 1).toString()} of ${value.length.toString()}.`
        }
        return 'Sortable item is no longer over a droppable area. Press escape to cancel.'
      },
      onDragStart({ active }) {
        const activeValue = active.id.toString()
        const sortableData = active.data.current as null | {
          sortable?: { index: number }
        }
        const currentIndex = sortableData?.sortable?.index ?? 0
        return `Grabbed sortable item "${activeValue}". Current position is ${(currentIndex + 1).toString()} of ${value.length.toString()}. Use arrow keys to move, space to drop.`
      }
    }),
    [value]
  )

  const screenReaderInstructions: ScreenReaderInstructions = React.useMemo(
    () => ({
      draggable: `
        To pick up a sortable item, press space or enter.
        While dragging, use the ${
          orientation === 'vertical' ? 'up and down'
          : orientation === 'horizontal' ? 'left and right'
          : 'arrow'
        } keys to move the item.
        Press space or enter again to drop the item in its new position, or press escape to cancel.
      `
    }),
    [orientation]
  )

  const contextValue = React.useMemo(
    () => ({
      activeId,
      flatCursor,
      getItemValue,
      id,
      items,
      modifiers: modifiers ?? config.modifiers,
      setActiveId,
      strategy: strategy ?? config.strategy
    }),
    [
      id,
      items,
      modifiers,
      strategy,
      config.modifiers,
      config.strategy,
      activeId,
      getItemValue,
      flatCursor
    ]
  )

  return (
    <SortableRootContext.Provider
      value={contextValue as SortableRootContextValue<unknown>}
    >
      {/* @ts-ignore */}
      <DndContext
        collisionDetection={collisionDetection ?? config.collisionDetection}
        modifiers={modifiers ?? config.modifiers}
        sensors={sensors}
        {...sortableProps}
        accessibility={{
          announcements,
          screenReaderInstructions,
          ...accessibility
        }}
        id={id}
        onDragCancel={onDragCancel}
        onDragEnd={onDragEnd}
        onDragStart={onDragStart}
      />
    </SortableRootContext.Provider>
  )
}

function useSortableContext(consumerName: string) {
  const context = React.useContext(SortableRootContext)
  if (!context) {
    throw new Error(`\`${consumerName}\` must be used within \`${ROOT_NAME}\``)
  }
  return context
}

const SortableContentContext = React.createContext<boolean>(false)
SortableContentContext.displayName = CONTENT_NAME

interface SortableContentProps extends React.ComponentPropsWithoutRef<'div'> {
  asChild?: boolean
  children: React.ReactNode
  strategy?: SortableContextProps['strategy']
  withoutSlot?: boolean
}

const SortableContent = React.forwardRef<HTMLDivElement, SortableContentProps>(
  (props, forwardedRef) => {
    const {
      asChild,
      children,
      strategy: strategyProp,
      withoutSlot,
      ...contentProps
    } = props

    const context = useSortableContext(CONTENT_NAME)

    const ContentPrimitive = asChild ? Slot : 'div'

    return (
      <SortableContentContext.Provider value={true}>
        {/* @ts-ignore */}
        <SortableContext
          items={context.items}
          strategy={strategyProp ?? context.strategy}
        >
          {withoutSlot ?
            children
          : <ContentPrimitive
              data-slot='sortable-content'
              {...contentProps}
              ref={forwardedRef}
            >
              {children}
            </ContentPrimitive>
          }
        </SortableContext>
      </SortableContentContext.Provider>
    )
  }
)
SortableContent.displayName = CONTENT_NAME

interface SortableItemContextValue {
  attributes: DraggableAttributes
  disabled?: boolean
  id: string
  isDragging?: boolean
  listeners: DraggableSyntheticListeners | undefined
  setActivatorNodeRef: (node: HTMLElement | null) => void
}

const SortableItemContext =
  React.createContext<null | SortableItemContextValue>(null)
SortableItemContext.displayName = ITEM_NAME

interface SortableItemProps extends React.ComponentPropsWithoutRef<'div'> {
  asChild?: boolean
  asHandle?: boolean
  disabled?: boolean
  value: UniqueIdentifier
}

function useSortableItemContext(consumerName: string) {
  const context = React.useContext(SortableItemContext)
  if (!context) {
    throw new Error(`\`${consumerName}\` must be used within \`${ITEM_NAME}\``)
  }
  return context
}

const SortableItem = React.forwardRef<HTMLDivElement, SortableItemProps>(
  (props, forwardedRef) => {
    const {
      asChild,
      asHandle,
      className,
      disabled,
      style,
      value,
      ...itemProps
    } = props

    const inSortableContent = React.useContext(SortableContentContext)
    const inSortableOverlay = React.useContext(SortableOverlayContext)

    if (!inSortableContent && !inSortableOverlay) {
      throw new Error(
        `\`${ITEM_NAME}\` must be used within \`${CONTENT_NAME}\` or \`${OVERLAY_NAME}\``
      )
    }

    if (value === '') {
      throw new Error(`\`${ITEM_NAME}\` value cannot be an empty string`)
    }

    const context = useSortableContext(ITEM_NAME)
    const id = React.useId()
    const {
      attributes,
      isDragging,
      listeners,
      setActivatorNodeRef,
      setNodeRef,
      transform,
      transition
      // @ts-ignore
    } = useSortable({ disabled, id: value })

    const composedRef = useComposedRefs(forwardedRef, (node) => {
      if (disabled) return
      setNodeRef(node)
      if (asHandle) setActivatorNodeRef(node)
    })

    const composedStyle = React.useMemo<React.CSSProperties>(() => {
      return {
        transform: CSS.Translate.toString(transform),
        transition,
        ...style
      }
    }, [transform, transition, style])

    const itemContext = React.useMemo<SortableItemContextValue>(
      // @ts-ignore
      () => ({
        attributes,
        disabled,
        id,
        isDragging,
        listeners,
        setActivatorNodeRef
      }),
      [id, attributes, listeners, setActivatorNodeRef, isDragging, disabled]
    )

    const ItemPrimitive = asChild ? Slot : 'div'

    return (
      <SortableItemContext.Provider value={itemContext}>
        <ItemPrimitive
          data-disabled={disabled}
          data-dragging={isDragging ? '' : undefined}
          data-slot='sortable-item'
          id={id}
          {...itemProps}
          {...(asHandle && !disabled ? attributes : {})}
          {...(asHandle && !disabled ? listeners : {})}
          className={cn(
            'focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:outline-hidden',
            {
              'cursor-default': context.flatCursor,
              'cursor-grab': !isDragging && asHandle && !context.flatCursor,
              'data-dragging:cursor-grabbing': !context.flatCursor,
              'opacity-50': isDragging,
              'pointer-events-none opacity-50': disabled,
              'touch-none select-none': asHandle
            },
            className
          )}
          ref={composedRef}
          style={composedStyle}
        />
      </SortableItemContext.Provider>
    )
  }
)
SortableItem.displayName = ITEM_NAME

interface SortableItemHandleProps
  extends React.ComponentPropsWithoutRef<'button'> {
  asChild?: boolean
}

const SortableItemHandle = React.forwardRef<
  HTMLButtonElement,
  SortableItemHandleProps
>((props, forwardedRef) => {
  const { asChild, className, disabled, ...itemHandleProps } = props

  const context = useSortableContext(ITEM_HANDLE_NAME)
  const itemContext = useSortableItemContext(ITEM_HANDLE_NAME)

  const isDisabled = disabled ?? itemContext.disabled

  const composedRef = useComposedRefs(forwardedRef, (node) => {
    if (!isDisabled) return
    itemContext.setActivatorNodeRef(node)
  })

  const HandlePrimitive = asChild ? Slot : 'button'

  return (
    <HandlePrimitive
      aria-controls={itemContext.id}
      data-disabled={isDisabled}
      data-dragging={itemContext.isDragging ? '' : undefined}
      data-slot='sortable-item-handle'
      type='button'
      {...itemHandleProps}
      {...(isDisabled ? {} : itemContext.attributes)}
      {...(isDisabled ? {} : itemContext.listeners)}
      className={cn(
        'select-none disabled:pointer-events-none disabled:opacity-50',
        context.flatCursor ? 'cursor-default' : (
          'cursor-grab data-dragging:cursor-grabbing'
        ),
        className
      )}
      disabled={isDisabled}
      ref={composedRef}
    />
  )
})
SortableItemHandle.displayName = ITEM_HANDLE_NAME

const SortableOverlayContext = React.createContext(false)
SortableOverlayContext.displayName = OVERLAY_NAME

const dropAnimation: DropAnimation = {
  sideEffects: defaultDropAnimationSideEffects({
    styles: {
      active: {
        opacity: '0.4'
      }
    }
  })
}

interface SortableOverlayProps
  extends Omit<React.ComponentPropsWithoutRef<typeof DragOverlay>, 'children'> {
  children?:
    | ((params: { value: UniqueIdentifier }) => React.ReactNode)
    | React.ReactNode
  container?: DocumentFragment | Element | null
}

function SortableOverlay(props: SortableOverlayProps) {
  const { children, container: containerProp, ...overlayProps } = props

  const context = useSortableContext(OVERLAY_NAME)

  const [mounted, setMounted] = React.useState(false)
  React.useLayoutEffect(() => {
    setMounted(true)
  }, [])

  const container = containerProp ?? (mounted ? globalThis.document.body : null)

  if (!container) return null

  return ReactDOM.createPortal(
    // @ts-ignore
    <DragOverlay
      className={cn(!context.flatCursor && 'cursor-grabbing')}
      dropAnimation={dropAnimation}
      modifiers={context.modifiers}
      {...overlayProps}
    >
      <SortableOverlayContext.Provider value={true}>
        {context.activeId ?
          typeof children === 'function' ?
            children({ value: context.activeId })
          : children
        : null}
      </SortableOverlayContext.Provider>
    </DragOverlay>,
    container
  )
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
  SortableOverlay
}
