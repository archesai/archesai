import type { LucideIcon } from '#components/custom/icons'

export interface SiteRoute {
  children?: SiteRoute[]
  description?: string
  href: string
  Icon: LucideIcon
  section: string
  showInTabs?: boolean
  title: string
}
