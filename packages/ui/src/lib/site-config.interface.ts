import type { LucideIcon } from 'lucide-react'

export interface SiteRoute {
  children?: SiteRoute[]
  description?: string
  href: string
  Icon: LucideIcon
  section: string
  showInTabs?: boolean
  title: string
}
