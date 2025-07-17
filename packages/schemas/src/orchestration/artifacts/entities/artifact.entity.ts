import type {
  Static,
  TNull,
  TNumber,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const ArtifactEntitySchema: TObject<{
  createdAt: TString
  credits: TNumber
  description: TUnion<[TNull, TString]>
  id: TString
  mimeType: TString
  name: TUnion<[TNull, TString]>
  organizationId: TString
  previewImage: TUnion<[TNull, TString]>
  producerId: TUnion<[TNull, TString]>
  text: TUnion<[TNull, TString]>
  updatedAt: TString
  url: TUnion<[TNull, TString]>
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    credits: Type.Number({
      description:
        'The number of credits required to access this artifact. This is used for metering and billing purposes.'
    }),
    description: Type.Union([Type.Null(), Type.String()], {
      description: "The artifact's description"
    }),
    // embedding: Type.Optional(
    //   Type.Array(Type.Number(), {
    //     description:
    //       "The artifact's embedding, used for semantic search and other ML tasks"
    //   })
    // ),
    mimeType: Type.String({
      description: 'The MIME type of the artifact, e.g. image/png'
    }),
    name: Type.Union([Type.Null(), Type.String()], {
      description: 'The name of the artifact, used for display purposes'
    }),
    organizationId: Type.String({ description: 'The organization name' }),
    previewImage: Type.Union([Type.Null(), Type.String()], {
      description:
        'The URL of the preview image for this artifact. This is used for displaying a thumbnail in the UI.'
    }),
    producerId: Type.Union([Type.Null(), Type.String()], {
      description:
        'The ID of the run that produced this artifact, if applicable'
    }),
    text: Type.Union([Type.Null(), Type.String()], {
      description: 'The artifact text'
    }),
    url: Type.Union([Type.Null(), Type.String()], {
      description: 'The artifact URL'
    })
  },
  {
    $id: 'ArtifactEntity',
    description: 'The artifact entity',
    title: 'Artifact Entity'
  }
)

export type ArtifactEntity = Static<typeof ArtifactEntitySchema>

export const ARTIFACT_ENTITY_KEY = 'artifacts'
