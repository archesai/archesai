import type { SiteRoute } from '#lib/site-config.interface'

export const TitleAndDescription = ({
  siteRoute
}: {
  siteRoute: SiteRoute
}) => {
  if (!siteRoute.title) return null
  return (
    <div className='container flex items-center gap-3 border-b px-4 py-3'>
      <siteRoute.Icon className='h-8 w-8' />
      <div>
        <p className='text-xl font-semibold text-foreground/85'>
          {siteRoute.title}
        </p>
        <p className='text-sm text-muted-foreground'>{siteRoute.description}</p>
      </div>
    </div>
  )
}
