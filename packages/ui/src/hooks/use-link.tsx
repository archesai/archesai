/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import type { LinkProps as TanStackLinkProps } from '@tanstack/react-router'
import type { ComponentPropsWithoutRef, ReactNode } from 'react'

import { createContext, useCallback, useContext } from 'react'

export interface AnalyticsLinkProps {
  // Analytics tracking
  trackClick?: boolean
  trackingAction?: string
  trackingCategory?: string
  trackingLabel?: string
  trackingValue?: number
}

export interface BaseLinkProps
  extends Omit<ComponentPropsWithoutRef<'a'>, 'href' | 'onClick'> {
  children: ReactNode
  href: string
  onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void
}

export interface ExternalLinkProps {
  download?: boolean | string
  external?: boolean
  newTab?: boolean
  noOpener?: boolean
  noReferrer?: boolean
}

export type FilterUndefined<T> = {
  [K in keyof T as T[K] extends undefined ? never : K]: T[K]
}

export type LinkComponent = React.ComponentType<SmartLinkProps>

export type SmartLinkProps = AnalyticsLinkProps &
  BaseLinkProps &
  ExternalLinkProps &
  TanStackRouterProps

export interface TanStackRouterProps {
  disabled?: TanStackLinkProps['disabled']
  hash?: TanStackLinkProps['hash']
  mask?: TanStackLinkProps['mask']
  params?: TanStackLinkProps['params']
  preload?: TanStackLinkProps['preload']
  preloadDelay?: TanStackLinkProps['preloadDelay']
  preserveSearch?: boolean
  replace?: TanStackLinkProps['replace']
  resetScroll?: TanStackLinkProps['resetScroll']
  resetSearch?: boolean
  search?: TanStackLinkProps['search']
  state?: TanStackLinkProps['state']
}

// Context definition
interface LinkContextType {
  isExternalUrl?: (url: string) => boolean
  Link: LinkComponent
  trackEvent?: (
    category: string,
    action: string,
    label?: string,
    value?: number
  ) => void
}

export function filterUndefined<T extends Record<string, unknown>>(
  obj: T
): FilterUndefined<T> {
  return Object.fromEntries(
    Object.entries(obj).filter(([_, value]) => value !== undefined)
  ) as FilterUndefined<T>
}

const LinkContext = createContext<LinkContextType | undefined>(undefined)

export function buildTanStackProps(
  props: TanStackRouterProps & {
    className?: string
    href: string
    onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void
  }
) {
  const {
    className,
    disabled,
    hash,
    href,
    mask,
    onClick,
    params,
    preload,
    preloadDelay,
    replace,
    resetScroll,
    search,
    state
  } = props

  // Build props object with only defined values
  const tanstackProps: Record<string, unknown> = { to: href }

  if (search !== undefined) tanstackProps.search = search
  if (params !== undefined) tanstackProps.params = params
  if (hash !== undefined) tanstackProps.hash = hash
  if (state !== undefined) tanstackProps.state = state
  if (mask !== undefined) tanstackProps.mask = mask
  if (replace !== undefined) tanstackProps.replace = replace
  if (resetScroll !== undefined) tanstackProps.resetScroll = resetScroll
  if (preload !== undefined) tanstackProps.preload = preload
  if (preloadDelay !== undefined) tanstackProps.preloadDelay = preloadDelay
  if (disabled !== undefined) tanstackProps.disabled = disabled
  if (className !== undefined) tanstackProps.className = className
  if (onClick !== undefined) tanstackProps.onClick = onClick

  return tanstackProps as TanStackLinkProps
}

// Default implementations
const defaultIsExternalUrl = (url: string): boolean => {
  try {
    return (
      url.startsWith('http://') ||
      url.startsWith('https://') ||
      url.startsWith('mailto:') ||
      url.startsWith('tel:') ||
      url.startsWith('ftp://') ||
      url.startsWith('//')
    )
  } catch {
    return false
  }
}

const defaultTrackEvent = (
  category: string,
  action: string,
  label?: string,
  value?: number
) => {
  // Default analytics implementation
  if (typeof window !== 'undefined') {
    // Google Analytics 4
    // if ('gtag' in window && typeof window.gtag === 'function') {
    //   window.gtag('event', action, {
    //     event_category: category,
    //     event_label: label,
    //     value: value
    //   })
    // }

    // // Google Analytics Universal
    // if ('ga' in window && typeof window.ga === 'function') {
    //   window.ga('send', 'event', category, action, label, value)
    // }

    // Console log for development
    if (process.env.NODE_ENV === 'development') {
      console.log('📊 Analytics Event:', { action, category, label, value })
    }
  }
}

// Default anchor link component
const DefaultLink: LinkComponent = ({
  children,
  external,
  href,
  newTab = true,
  noOpener = true,
  noReferrer = false,
  onClick,
  trackClick,
  trackingAction = 'click',
  trackingCategory = 'Link',
  trackingLabel,
  ...props
}) => {
  const context = useContext(LinkContext)
  const isExternal =
    context?.isExternalUrl?.(href) ?? defaultIsExternalUrl(href)
  const trackEvent = context?.trackEvent ?? defaultTrackEvent

  const handleClick = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      // Track analytics if enabled
      if (trackClick) {
        trackEvent(
          trackingCategory,
          trackingAction,
          trackingLabel ?? href,
          undefined
        )
      }

      // Call original onClick
      onClick?.(e)
    },
    [
      trackClick,
      trackingCategory,
      trackingAction,
      trackingLabel,
      href,
      trackEvent,
      onClick
    ]
  )

  // Build props for external links
  const externalProps =
    isExternal || external ?
      {
        rel:
          [noOpener && 'noopener', noReferrer && 'noreferrer']
            .filter(Boolean)
            .join(' ') || undefined,
        target: newTab ? '_blank' : undefined
      }
    : {}

  return (
    <a
      href={href}
      onClick={handleClick}
      {...externalProps}
      {...props}
    >
      {children}
    </a>
  )
}

// Hook to access the Link component
export const useLinkComponent = (): LinkComponent => {
  const context = useContext(LinkContext)
  return context?.Link ?? DefaultLink
}

// Hook to access full context
export const useLinkContext = () => {
  const context = useContext(LinkContext)
  return {
    hasProvider: !!context,
    isExternalUrl: context?.isExternalUrl ?? defaultIsExternalUrl,
    Link: context?.Link ?? DefaultLink,
    trackEvent: context?.trackEvent ?? defaultTrackEvent
  }
}

// Provider component
interface LinkProviderProps {
  children: ReactNode
  isExternalUrl?: (url: string) => boolean
  Link?: LinkComponent
  trackEvent?: (
    category: string,
    action: string,
    label?: string,
    value?: number
  ) => void
}

export const LinkProvider: React.FC<LinkProviderProps> = ({
  children,
  isExternalUrl = defaultIsExternalUrl,
  Link = DefaultLink,
  trackEvent = defaultTrackEvent
}) => {
  return (
    <LinkContext.Provider value={{ isExternalUrl, Link, trackEvent }}>
      {children}
    </LinkContext.Provider>
  )
}
