import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const ToolEntitySchema: TObject<{
  createdAt: TString
  description: TString
  id: TString
  inputMimeType: TString
  name: TString
  organizationId: TString
  outputMimeType: TString
  toolBase: TString
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    description: Type.String({ description: 'The tool description' }),
    inputMimeType: Type.String({
      description: 'The MIME type of the input for the tool, e.g. text/plain'
    }),
    name: Type.String({
      description: 'The name of the tool'
    }),
    organizationId: Type.String({ description: 'The organization name' }),
    outputMimeType: Type.String({
      description: 'The MIME type of the output for the tool, e.g. text/plain'
    }),
    toolBase: Type.String({ description: 'The base of the tool' })
  },
  {
    $id: 'ToolEntity',
    description: 'The tool entity',
    title: 'Tool Entity'
  }
)

export type ToolEntity = Static<typeof ToolEntitySchema>

export const TOOL_ENTITY_KEY = 'tools'
