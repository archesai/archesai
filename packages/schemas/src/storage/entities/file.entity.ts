import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const FileEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    isDir: Type.Boolean({
      description: 'Whether or not this is a directory'
    }),
    orgname: Type.String({ description: 'The original name of the file' }),
    path: Type.String({ description: 'The path to the item' }),
    read: Type.Optional(
      Type.String({
        description:
          'The read-only URL that you can use to download the file from secure storage',
        format: 'uri'
      })
    ),
    size: Type.Number({ description: 'The size of the item in bytes' }),
    write: Type.Optional(
      Type.String({
        description:
          'The write-only URL that you can use to upload the file to secure storage',
        format: 'uri'
      })
    )
  },
  {
    $id: 'FileEntity',
    description: 'The file entity',
    title: 'File Entity'
  }
)

export type FileEntity = Static<typeof FileEntitySchema>

export const FILE_ENTITY_KEY = 'files'
