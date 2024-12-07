'use client'

import { useSidebar } from '@/components/ui/sidebar'
import { useOrganizationsControllerFindOne } from '@/generated/archesApiComponents'

export const CreditQuota = () => {
  const { open } = useSidebar()
  const { defaultOrgname } = useAuth()
  const { data: organization } = useOrganizationsControllerFindOne({
    pathParams: {
      orgname: defaultOrgname
    }
  })

  if (!open) {
    return <CreditCircularChart remaining={organization?.credits || 0} total={60000} />
  }
  return (
    <div className='inter flex w-full flex-col gap-2 rounded-lg bg-muted p-2 text-xs'>
      <div className='text-gray-alpha-500 flex items-center justify-between'>
        <div className='font-semibold'>Credit Usage</div>
        <div>
          <Link href='/organization/billing'>
            <Badge className='text-emerald-700 outline-black'>Upgrade</Badge>
          </Link>
        </div>
      </div>
      <div className='flex items-center gap-2'>
        <div className='flex flex-grow flex-col gap-2'>
          <div className='inter flex justify-between'>
            <div>Total</div>
            <div className='tabular-nums'>{organization?.credits}</div>
          </div>
          <div className='inter flex justify-between'>
            <div>Remaining</div>
            <div className='tabular-nums'>{organization?.credits}</div>
          </div>
        </div>
        <CreditCircularChart remaining={organization?.credits || 0} total={60000} />
      </div>
    </div>
  )
}

import { ChartConfig, ChartContainer } from '@/components/ui/chart'
import { useAuth } from '@/hooks/use-auth'
import Link from 'next/link'
import { PolarGrid, PolarRadiusAxis, RadialBar, RadialBarChart } from 'recharts'

import { Badge } from '../../ui/badge'

export const description = 'A radial chart with a custom shape'

const chartData = [{ browser: 'safari', fill: 'var(--color-safari)', visitors: 1260 }]

const chartConfig = {
  safari: {
    color: 'hsl(var(--chart-2))',
    label: 'Safari'
  },
  visitors: {
    label: 'Visitors'
  }
} satisfies ChartConfig

export function CreditCircularChart({ remaining, total }: { remaining: number; total: number }) {
  return (
    <ChartContainer className='mx-auto aspect-square h-[25px]' config={chartConfig}>
      <RadialBarChart
        data={chartData}
        endAngle={(remaining / total) * 360}
        height={200}
        innerRadius={10}
        outerRadius={20}
      >
        <PolarGrid gridType='circle' polarRadius={[86, 74]} radialLines={false} stroke='none' />
        <RadialBar background dataKey='visitors' />
        <PolarRadiusAxis axisLine={false} tick={false} tickLine={false}></PolarRadiusAxis>
      </RadialBarChart>
    </ChartContainer>
  )
}
