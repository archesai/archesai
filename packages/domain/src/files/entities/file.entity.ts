import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

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

export class FileEntity
  extends BaseEntity
  implements Static<typeof FileEntitySchema>
{
  public isDir: boolean
  public orgname: string
  public path: string
  public read?: string
  public size: number
  public type = FILE_ENTITY_KEY
  public write?: string

  constructor(props: FileEntity) {
    super(props)
    this.isDir = props.isDir
    this.orgname = props.orgname
    this.path = props.path
    if (props.read) {
      this.read = props.read
    }
    if (props.write) {
      this.write = props.write
    }
    this.size = props.size
  }
}

export const FILE_ENTITY_KEY = 'files'
