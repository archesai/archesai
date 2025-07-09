import { Suspense } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'

import type { ArtifactEntity } from '@archesai/schemas'

import { getFindManyArtifactsSuspenseQueryOptions } from '@archesai/client'
import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'
import { ArtifactViewer } from '@archesai/ui/components/custom/artifact-viewer'
import { ContentTypeToIcon } from '@archesai/ui/components/custom/content-type-to-icon'
import { File, ScanSearch } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger
} from '@archesai/ui/components/shadcn/hover-card'
import { Skeleton } from '@archesai/ui/components/shadcn/skeleton'

import ContentForm from '#components/forms/artifact-form'

export default function ArtifactDataTable({
  readonly = false
}: {
  readonly?: boolean
}) {
  const navigate = useNavigate()

  return (
    <DataTable<ArtifactEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <ContentTypeToIcon contentType={row.original.mimeType || ''} />
                <Link
                  className='shrink truncate text-wrap text-blue-600 underline md:text-sm'
                  search={{
                    artifactId: row.original.id
                  }}
                  to={`/artifacts/single`}
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
                    <Link
                      search={{
                        artifactId: row.original.id
                      }}
                      to={`/artifacts/single`}
                    >
                      <HoverCardTrigger asChild>
                        <ScanSearch />
                      </HoverCardTrigger>
                    </Link>

                    <HoverCardContent
                      className='h-[500px] w-[500px]'
                      side='right'
                    >
                      <Suspense fallback={<Skeleton />}>
                        <ArtifactViewer content={row.original} />
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
            return row.original.parentId ?
                <Link
                  className='max-w-lg truncate text-base text-wrap md:text-sm'
                  search={{
                    artifactId: row.original.parentId
                  }}
                  to={`/artifacts/single`}
                >
                  {row.original.parentId}
                </Link>
              : <div className='text-muted-foreground'>None</div>
          },
          enableSorting: false
        },
        {
          accessorKey: 'producer',
          cell: ({ row }) => {
            return row.original.producerId ?
                <Link
                  className='max-w-lg truncate text-base text-wrap md:text-sm'
                  search={{
                    selectedRunId: row.original.producerId
                  }}
                  to={`/playground`}
                >
                  {row.original.producerId}
                </Link>
              : <div className='text-muted-foreground'>None</div>
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
      entityType={ARTIFACT_ENTITY_KEY}
      getEditFormFromItem={(content) => <ContentForm artifactId={content.id} />}
      grid={(_item) => {
        return (
          <div className='flex h-full w-full items-center justify-center'>
            {/* <Image
              alt='source image'
              height={256}
              src={item.previewImage}
              width={256}
            /> */}
          </div>
        )
      }}
      handleSelect={async (content) => {
        await navigate({ to: `/artifacts/single?artifactId=${content.id}` })
      }}
      icon={<File size={24} />}
      readonly={readonly}
      useFindMany={getFindManyArtifactsSuspenseQueryOptions()}
    />
  )
}
