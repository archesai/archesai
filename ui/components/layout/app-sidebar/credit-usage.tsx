'use client'

import { useSidebar } from '@/components/ui/sidebar'
import { useOrganizationsControllerFindOne } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { Badge } from '@/components/ui/badge'
import Link from 'next/link'
import { Skeleton } from '@/components/ui/skeleton'

export const CreditQuota = () => {
  const { open } = useSidebar()
  const { defaultOrgname } = useAuth()
  const { data: organization, isFetched } = useOrganizationsControllerFindOne({
    pathParams: {
      orgname: defaultOrgname
    }
  })

  if (!open) {
    return <></>
  }
  return (
    <div className='inter flex w-full flex-col gap-2 rounded-lg bg-sidebar-accent p-2 text-xs'>
      <div className='flex items-center justify-between'>
        <div className='font-semibold'>Credit Usage</div>
        <div>
          <Link href='/organization/billing'>
            <Badge>Upgrade</Badge>
          </Link>
        </div>
      </div>
      <div className='flex items-center gap-2'>
        <div className='flex flex-grow flex-col gap-2'>
          <div className='inter flex justify-between'>
            <div>Total</div>
            <div className='tabular-nums'>
              {isFetched ? (
                organization?.credits
              ) : (
                <Skeleton className='h-2 w-20 bg-slate-600' />
              )}
            </div>
          </div>
          <div className='inter flex justify-between'>
            <div>Remaining</div>
            <div className='tabular-nums'>
              {isFetched ? (
                organization?.credits
              ) : (
                <Skeleton className='h-2 w-20 bg-slate-600' />
              )}
            </div>
          </div>
        </div>
        <></>
      </div>
    </div>
  )
}
