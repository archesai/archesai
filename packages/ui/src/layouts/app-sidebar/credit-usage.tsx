import type { OrganizationEntity } from '@archesai/domain'

import { Badge } from '#components/shadcn/badge'
import { useSidebar } from '#components/shadcn/sidebar'
import { Skeleton } from '#components/shadcn/skeleton'

export interface CreditQuotaProps {
  children?: never
  organization?: OrganizationEntity
}

export const CreditQuota = ({ organization }: CreditQuotaProps) => {
  const { open } = useSidebar()

  if (!open) {
    return <></>
  }

  return (
    <div className='inter flex w-full flex-col gap-2 rounded-lg bg-sidebar-accent p-2 text-xs'>
      <div className='flex items-center justify-between'>
        <div className='font-semibold'>Credit Usage</div>
        <div>
          <a
            className='flex'
            href='/organization/billing'
          >
            <Badge variant='secondary'>Upgrade</Badge>
          </a>
        </div>
      </div>
      <div className='flex items-center gap-2'>
        <div className='flex flex-grow flex-col gap-2'>
          <div className='inter flex justify-between'>
            <div>Total</div>
            <div className='tabular-nums'>
              {organization ?
                organization.credits
              : <Skeleton className='h-2 w-20 bg-slate-600' />}
            </div>
          </div>
          <div className='inter flex justify-between'>
            <div>Remaining</div>
            <div className='tabular-nums'>
              {organization ?
                organization.credits
              : <Skeleton className='h-2 w-20 bg-slate-600' />}
            </div>
          </div>
        </div>
        <></>
      </div>
    </div>
  )
}
