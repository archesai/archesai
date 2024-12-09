'use client'
import { ContentViewer } from '@/components/content-viewer'
import { DataTable } from '@/components/datatables/datatable/data-table'
import { DataTableColumnHeader } from '@/components/datatables/datatable/data-table-column-header'
import ContentForm from '@/components/forms/content-form'
import { Badge } from '@/components/ui/badge'
import {
  ContentControllerFindAllPathParams,
  ContentControllerRemoveVariables,
  useContentControllerFindAll,
  useContentControllerRemove
} from '@/generated/archesApiComponents'
import { ContentEntity, FieldFieldQuery } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger
} from '@/components/ui/hover-card'

import { format } from 'date-fns'
import { File } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { Suspense } from 'react'
import { Skeleton } from '../ui/skeleton'

export default function ContentDataTable({
  customFilters,
  readonly
}: {
  customFilters?: FieldFieldQuery[]
  readonly?: boolean
}) {
  const router = useRouter()
  const { defaultOrgname } = useAuth()

  return (
    <DataTable<
      ContentEntity,
      ContentControllerFindAllPathParams,
      ContentControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                {/* <ContentTypeToIcon contentType={row.original.mimeType} /> */}
                <Link
                  className='shrink truncate text-wrap md:text-sm'
                  href={`/content/single?contentId=${row.original.id}`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          },
          header: ({ column }) => (
            <DataTableColumnHeader
              column={column}
              title='Name'
            />
          )
        },
        {
          accessorKey: 'mimeType',
          cell: ({ row }) => {
            return <Badge>{row.original?.mimeType}</Badge>
          },
          header: ({ column }) => (
            <DataTableColumnHeader
              column={column}
              title='Type'
            />
          )
        },
        {
          accessorKey: 'value',
          cell: ({ row }) => {
            return (
              <div className='truncate text-wrap text-base md:text-sm'>
                {row.original.text || (
                  <HoverCard openDelay={0}>
                    <HoverCardTrigger asChild>
                      <Link
                        className='underline underline-offset-4'
                        href={`/content/single?contentId=${row.original?.id}`}
                      >
                        View Content
                      </Link>
                    </HoverCardTrigger>
                    <HoverCardContent>
                      <Suspense fallback={<Skeleton />}>
                        <ContentViewer id={row.original.id} />
                      </Suspense>
                    </HoverCardContent>
                  </HoverCard>
                )}
              </div>
            )
          },
          enableHiding: false,
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className='-ml-2 text-sm'
              column={column}
              title='Data'
            />
          )
        },
        {
          accessorKey: 'parent',
          cell: ({ row }) => {
            return row.original.parent?.name ? (
              <Link
                className='max-w-lg truncate text-wrap text-base md:text-sm'
                href={`/content/single?contentId=${row.original.parent?.id}`}
              >
                {row.original.parent?.name}
              </Link>
            ) : (
              <div className='text-muted-foreground'>None</div>
            )
          },
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className='-ml-2 text-sm'
              column={column}
              title='Parent'
            />
          )
        },
        {
          accessorKey: 'producedBy',
          cell: ({ row }) => {
            return row.original.producedBy ? (
              <Link
                className='max-w-lg truncate text-wrap text-base md:text-sm'
                href={`/playground?selectedRunId=${row.original.producedBy?.id}`}
              >
                {row.original.producedBy?.name}
              </Link>
            ) : (
              <div className='text-muted-foreground'>None</div>
            )
          },
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className='-ml-2 text-sm'
              column={column}
              title='Produced By'
            />
          )
        },
        {
          accessorKey: 'labels',
          cell: ({ row }) => {
            return (
              <div className='flex gap-1'>
                {row.original.labels?.length ? (
                  row.original.labels?.map((label, index) => (
                    <Badge key={index}>{label.name}</Badge>
                  ))
                ) : (
                  <div className='text-muted-foreground'>None</div>
                )}
              </div>
            )
          },
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className='-ml-2 text-sm'
              column={column}
              title='Labels'
            />
          )
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {format(new Date(row.original.createdAt), 'M/d/yy h:mm a')}
              </span>
            )
          },
          header: ({ column }) => (
            <DataTableColumnHeader
              column={column}
              title='Created'
            />
          )
        }
      ]}
      content={(item) => {
        return (
          <div className='flex h-full w-full items-center justify-center'>
            <Image
              alt='source image'
              height={256}
              src={item.previewImage || ''}
              width={256}
            />
          </div>
        )
      }}
      createForm={<ContentForm />}
      customFilters={customFilters}
      dataIcon={<File size={24} />}
      defaultView='table'
      findAllPathParams={{
        orgname: defaultOrgname
      }}
      getDeleteVariablesFromItem={(content) => ({
        pathParams: {
          id: content.id,
          orgname: defaultOrgname
        }
      })}
      getEditFormFromItem={(content) => <ContentForm contentId={content.id} />}
      handleSelect={(content) =>
        router.push(`/content/single?contentId=${content.id}`)
      }
      itemType='content'
      readonly={readonly}
      useFindAll={useContentControllerFindAll}
      useRemove={useContentControllerRemove}
    />
  )
}
