'use client'

import { Suspense } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

import type { ContentEntity } from '@archesai/domain'

import { useFindManyContents } from '@archesai/client'
import { CONTENT_ENTITY_KEY } from '@archesai/domain'
import { ContentTypeToIcon } from '@archesai/ui/components/custom/content-type-to-icon'
import { ContentViewer } from '@archesai/ui/components/custom/content-viewer'
import { File, ScanSearch } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger
} from '@archesai/ui/components/shadcn/hover-card'
import { Skeleton } from '@archesai/ui/components/shadcn/skeleton'

import ContentForm from '#components/forms/content-form'

export default function ContentDataTable({
  readonly = false
}: {
  readonly?: boolean
}) {
  const router = useRouter()

  return (
    <DataTable<ContentEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <ContentTypeToIcon contentType={row.original.mimeType || ''} />
                <Link
                  className='shrink truncate text-wrap text-blue-600 underline md:text-sm'
                  href={`/content/single?contentId=${row.original.id}`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          }
        },
        {
          accessorKey: 'value',
          cell: ({ row }) => {
            return (
              <div className='truncate text-base text-wrap md:text-sm'>
                {row.original.text ?? (
                  <HoverCard openDelay={200}>
                    <Link href={`/content/single?contentId=${row.original.id}`}>
                      <HoverCardTrigger asChild>
                        <ScanSearch />
                      </HoverCardTrigger>
                    </Link>

                    <HoverCardContent
                      className='h-[500px] w-[500px]'
                      side='right'
                    >
                      <Suspense fallback={<Skeleton />}>
                        <ContentViewer content={row.original} />
                      </Suspense>
                    </HoverCardContent>
                  </HoverCard>
                )}
              </div>
            )
          },
          enableHiding: false,
          enableSorting: false
        },
        {
          accessorKey: 'parent',
          cell: ({ row }) => {
            return row.original.parentId ? (
              <Link
                className='max-w-lg truncate text-base text-wrap md:text-sm'
                href={`/content/single?contentId=${row.original.parentId}`}
              >
                {row.original.parentId}
              </Link>
            ) : (
              <div className='text-muted-foreground'>None</div>
            )
          },
          enableSorting: false
        },
        {
          accessorKey: 'producer',
          cell: ({ row }) => {
            return row.original.producerId ? (
              <Link
                className='max-w-lg truncate text-base text-wrap md:text-sm'
                href={`/playground?selectedRunId=${row.original.producerId}`}
              >
                {row.original.producerId}
              </Link>
            ) : (
              <div className='text-muted-foreground'>None</div>
            )
          },
          enableSorting: false
        },
        // {
        //   accessorKey: 'labels',
        //   cell: ({ row }) => {
        //     return (
        //       <div className='flex gap-1'>
        //         {row.original.labels.length ? (
        //           row.original.labels.map((label, index) => (
        //             <Badge
        //               key={index}
        //               variant={'secondary'}
        //             >
        //               {label.name}
        //             </Badge>
        //           ))
        //         ) : (
        //           <div className='text-muted-foreground'>None</div>
        //         )}
        //       </div>
        //     )
        //   },
        //   enableSorting: false
        // },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          }
        }
      ]}
      createForm={<ContentForm />}
      defaultView='table'
      entityType={CONTENT_ENTITY_KEY}
      getEditFormFromItem={(content) => <ContentForm contentId={content.id} />}
      grid={(item) => {
        return (
          <div className='flex h-full w-full items-center justify-center'>
            <Image
              alt='source image'
              height={256}
              src={item.previewImage}
              width={256}
            />
          </div>
        )
      }}
      handleSelect={(content) => {
        router.push(`/content/single?contentId=${content.id}`)
      }}
      icon={<File size={24} />}
      readonly={readonly}
      useFindMany={useFindManyContents}
    />
  )
}
