import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

export const ArtifactEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    credits: Type.Number({
      description:
        'The number of credits required to access this artifact. This is used for metering and billing purposes.'
    }),
    description: Type.String({ description: "The artifact's description" }),
    embedding: Type.Array(Type.Number(), {
      description:
        "The artifact's embedding, used for semantic search and other ML tasks"
    }),
    mimeType: Type.String({
      description: 'The MIME type of the artifact, e.g. image/png'
    }),
    orgname: Type.String({ description: 'The organization name' }),
    parentId: Type.String({
      description:
        'The ID of the parent artifact, if this artifact is a child of another artifact'
    }),
    previewImage: Type.String({
      description:
        'The URL of the preview image for this artifact. This is used for displaying a thumbnail in the UI.'
    }),
    producerId: Type.String({
      description:
        'The ID of the run that produced this artifact, if applicable'
    }),
    text: Type.Optional(Type.String({ description: 'The artifact text' })),
    url: Type.Optional(Type.String({ description: 'The artifact URL' }))
  },
  {
    $id: 'ArtifactEntity',
    description: 'The artifact entity',
    title: 'Artifact Entity'
  }
)

export class ArtifactEntity
  extends BaseEntity
  implements Static<typeof ArtifactEntitySchema>
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
  public type = ARTIFACT_ENTITY_KEY
  public url?: string

  constructor(props: ArtifactEntity) {
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

// export class artifactRelationships {
//   children: BaseEntity[]
//   consumers: BaseEntity[]
//   labels: LabelEntity[]
//   parent?: BaseEntity | null
//   producer?: BaseEntity | null

//   constructor(props: artifactRelationships) {
//     this.children = props.children
//     this.consumers = props.consumers
//     this.labels = props.labels
//     this.parent = props.parent
//     this.producer = props.producer
//   }

//   static schema() {
//     return Type.Object({
//       children: Type.Array(BaseEntity.schema(), {
//         description: 'The children of the artifact'
//       }),
//       consumers: Type.Array(BaseEntity.schema(), {
//         description: 'The consumers of the artifact'
//       }),
//       labels: Type.Array(LabelEntity.schema(), {
//         description: 'The labels of the artifact'
//       }),
//       parent: Type.Optional(BaseEntity.schema({ nullable: true }), {
//         description: 'The parent of the artifact'
//       }),
//       producer: Type.Optional(BaseEntity.schema({ nullable: true }), {
//         description: 'The producer of the artifact'
//       })
//     })
//   }
// }

export const ARTIFACT_ENTITY_KEY = 'artifacts'
