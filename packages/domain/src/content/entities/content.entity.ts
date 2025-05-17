import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

export const ContentEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    credits: Type.Number({
      description:
        'The number of credits required to access this content. This is used for metering and billing purposes.'
    }),
    description: Type.String({ description: "The content's description" }),
    embedding: Type.Array(Type.Number(), {
      description:
        "The content's embedding, used for semantic search and other ML tasks"
    }),
    mimeType: Type.String({
      description: 'The MIME type of the content, e.g. image/png'
    }),
    orgname: Type.String({ description: 'The organization name' }),
    parentId: Type.String({
      description:
        'The ID of the parent content, if this content is a child of another content'
    }),
    previewImage: Type.String({
      description:
        'The URL of the preview image for this content. This is used for displaying a thumbnail in the UI.'
    }),
    producerId: Type.String({
      description: 'The ID of the run that produced this content, if applicable'
    }),
    text: Type.Optional(Type.String({ description: 'The content text' })),
    url: Type.Optional(Type.String({ description: 'The content URL' }))
  },
  {
    $id: 'ContentEntity',
    description: 'The content entity',
    title: 'Content Entity'
  }
)

export class ContentEntity
  extends BaseEntity
  implements Static<typeof ContentEntitySchema>
{
  public credits: number
  public description: string
  public embedding: number[]
  public mimeType: string
  public orgname: string
  public parentId: string
  public previewImage: string
  public producerId: string
  public text?: string
  public type = CONTENT_ENTITY_KEY
  public url?: string

  constructor(props: ContentEntity) {
    super(props)
    this.credits = props.credits
    this.description = props.description
    this.embedding = props.embedding
    this.mimeType = props.mimeType
    this.orgname = props.orgname
    this.parentId = props.parentId
    this.previewImage = props.previewImage
    this.producerId = props.producerId
    if (props.text) {
      this.text = props.text
    }
    if (props.url) {
      this.url = props.url
    }
  }
}

// export class ContentRelationships {
//   children: BaseEntity[]
//   consumers: BaseEntity[]
//   labels: LabelEntity[]
//   parent?: BaseEntity | null
//   producer?: BaseEntity | null

//   constructor(props: ContentRelationships) {
//     this.children = props.children
//     this.consumers = props.consumers
//     this.labels = props.labels
//     this.parent = props.parent
//     this.producer = props.producer
//   }

//   static schema() {
//     return Type.Object({
//       children: Type.Array(BaseEntity.schema(), {
//         description: 'The children of the content'
//       }),
//       consumers: Type.Array(BaseEntity.schema(), {
//         description: 'The consumers of the content'
//       }),
//       labels: Type.Array(LabelEntity.schema(), {
//         description: 'The labels of the content'
//       }),
//       parent: Type.Optional(BaseEntity.schema({ nullable: true }), {
//         description: 'The parent of the content'
//       }),
//       producer: Type.Optional(BaseEntity.schema({ nullable: true }), {
//         description: 'The producer of the content'
//       })
//     })
//   }
// }

export const CONTENT_ENTITY_KEY = 'contents'
