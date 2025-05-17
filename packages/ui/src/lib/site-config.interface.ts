import type { Workflow } from 'lucide-react'

export type RouteIcon = typeof Workflow

export interface SiteRoute {
  children?: SiteRoute[]
  description?: string
  href: string
  Icon: RouteIcon
  section: string
  showInTabs?: boolean
  title: string
}
