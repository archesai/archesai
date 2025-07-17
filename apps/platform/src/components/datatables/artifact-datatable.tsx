import { Suspense } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'

import type { ArtifactEntity } from '@archesai/schemas'

import { getFindManyArtifactsQueryOptions } from '@archesai/client'
import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'
import { ArtifactViewer } from '@archesai/ui/components/custom/artifact-viewer'
import { ContentTypeToIcon } from '@archesai/ui/components/custom/content-type-to-icon'
import {
  Calendar,
  File,
  LetterText,
  ScanSearch
} from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { TasksTableActionBar } from '@archesai/ui/components/datatable/components/tasks-table-action-bar'
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
                <ContentTypeToIcon contentType={row.original.mimeType ?? ''} />
                <Link
                  className='text-primary hover:underline'
                  params={{
                    artifactId: row.original.id
                  }}
                  to={`/artifacts/$artifactId`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          },
          enableColumnFilter: true,
          id: 'name',
          meta: {
            filterVariant: 'text',
            icon: LetterText,
            label: 'Name'
          }
        },
        {
          accessorKey: 'mimeType',
          cell: ({ row }) => {
            return (
              <div className='flex items-center gap-2'>
                <HoverCard openDelay={200}>
                  <Link
                    params={{
                      artifactId: row.original.id
                    }}
                    to={`/artifacts/$artifactId`}
                  >
                    <HoverCardTrigger asChild>
                      <ScanSearch />
                    </HoverCardTrigger>
                  </Link>

                  <HoverCardContent
                    className='h-min-[200] w-min-[200]'
                    side='right'
                  >
                    <Suspense fallback={<Skeleton />}>
                      <ArtifactViewer artifactId={row.original.id} />
                    </Suspense>
                  </HoverCardContent>
                </HoverCard>
                {row.original.mimeType}
              </div>
            )
          },
          enableColumnFilter: true,
          enableHiding: false,
          id: 'mimeType',
          meta: {
            filterVariant: 'multiSelect',
            icon: LetterText,
            label: 'Artifact Type',
            options: [
              { label: 'Text', value: 'text' },
              { label: 'Image', value: 'image' },
              { label: 'Audio', value: 'audio' },
              { label: 'Video', value: 'video' }
            ]
          }
        },
        {
          accessorKey: 'parent',
          cell: ({ row }) => {
            return row.original.parentId ?
                <Link
                  className='text-primary hover:underline'
                  params={{
                    artifactId: row.original.parentId
                  }}
                  to={`/artifacts/$artifactId`}
                >
                  {row.original.parentId}
                </Link>
              : <div className='text-muted-foreground'>None</div>
          },
          enableColumnFilter: true,
          enableSorting: true,
          id: 'parent',
          meta: {
            filterVariant: 'text',
            icon: LetterText,
            label: 'Parent'
          }
        },
        {
          accessorKey: 'producer',
          cell: ({ row }) => {
            return row.original.producerId ?
                <Link
                  className='text-primary hover:underline'
                  params={{
                    artifactId: row.original.id
                  }}
                  search={{
                    selectedRunId: row.original.producerId
                  }}
                  to={`/artifacts/$artifactId`}
                >
                  {row.original.producerId}
                </Link>
              : <div className='text-muted-foreground'>None</div>
          },
          enableColumnFilter: true,
          enableSorting: true,
          id: 'producer',
          meta: {
            filterVariant: 'text',
            icon: LetterText,
            label: 'Producer'
          }
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          enableColumnFilter: true,
          enableSorting: true,
          id: 'createdAt',
          meta: {
            filterVariant: 'date',
            icon: Calendar,
            label: 'Created'
          }
        }
      ]}
      createForm={<ContentForm />}
      defaultView='table'
      entityKey={ARTIFACT_ENTITY_KEY}
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
        await navigate({
          params: { artifactId: content.id },
          to: `/artifacts/$artifactId`
        })
      }}
      icon={<File size={24} />}
      readonly={readonly}
      useFindMany={getFindManyArtifactsQueryOptions}
    />
  )
}
