import type {
  Static,
  TArray,
  TNull,
  TNumber,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const ArtifactEntitySchema: TObject<{
  createdAt: TString
  credits: TNumber
  description: TString
  embedding: TOptional<TArray<TNumber>>
  id: TString
  mimeType: TString
  name: TString
  organizationId: TString
  parentId: TString
  previewImage: TString
  producerId: TUnion<[TString, TNull]>
  text: TOptional<TString>
  updatedAt: TString
  url: TOptional<TString>
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    credits: Type.Number({
      description:
        'The number of credits required to access this artifact. This is used for metering and billing purposes.'
    }),
    description: Type.String({ description: "The artifact's description" }),
    embedding: Type.Optional(
      Type.Array(Type.Number(), {
        description:
          "The artifact's embedding, used for semantic search and other ML tasks"
      })
    ),
    mimeType: Type.String({
      description: 'The MIME type of the artifact, e.g. image/png'
    }),
    name: Type.String({
      description: 'The name of the artifact, used for display purposes'
    }),
    organizationId: Type.String({ description: 'The organization name' }),
    parentId: Type.String({
      description:
        'The ID of the parent artifact, if this artifact is a child of another artifact'
    }),
    previewImage: Type.String({
      description:
        'The URL of the preview image for this artifact. This is used for displaying a thumbnail in the UI.'
    }),
    producerId: Type.Union([Type.String(), Type.Null()], {
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

export type ArtifactEntity = Static<typeof ArtifactEntitySchema>

export const ARTIFACT_ENTITY_KEY = 'artifacts'
