import { Link, useNavigate } from '@tanstack/react-router'

import type { ArtifactEntity } from '@archesai/schemas'

import { getFindManyArtifactsQueryOptions } from '@archesai/client'
import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'
import {
  CalendarIcon,
  FileIcon,
  TextIcon
} from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'

import ArtifactForm from '#components/forms/artifact-form'

export default function ArtifactDataTable() {
  const navigate = useNavigate()

  return (
    <DataTable<ArtifactEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
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
            icon: TextIcon,
            label: 'Name'
          }
        },
        {
          accessorKey: 'mimeType',
          cell: ({ row }) => {
            return <Badge variant={'secondary'}>{row.original.mimeType}</Badge>
          },
          enableColumnFilter: true,
          enableHiding: false,
          id: 'mimeType',
          meta: {
            filterVariant: 'multiSelect',
            icon: TextIcon,
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
            icon: TextIcon,
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
            icon: CalendarIcon,
            label: 'Created'
          }
        }
      ]}
      createForm={ArtifactForm}
      defaultView='table'
      entityKey={ARTIFACT_ENTITY_KEY}
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
      handleSelect={async (artifact) => {
        await navigate({
          params: { artifactId: artifact.id },
          to: `/artifacts/$artifactId`
        })
      }}
      icon={<FileIcon size={24} />}
      updateForm={ArtifactForm}
      useFindMany={getFindManyArtifactsQueryOptions}
    />
  )
}
