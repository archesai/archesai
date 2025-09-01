import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const ArtifactEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  credits: z.ZodNumber
  description: z.ZodNullable<z.ZodString>
  id: z.ZodUUID
  mimeType: z.ZodString
  name: z.ZodNullable<z.ZodString>
  organizationId: z.ZodString
  previewImage: z.ZodNullable<z.ZodString>
  producerId: z.ZodNullable<z.ZodString>
  text: z.ZodNullable<z.ZodString>
  updatedAt: z.ZodString
  url: z.ZodNullable<z.ZodString>
}> = BaseEntitySchema.extend({
  credits: z
    .number()
    .describe(
      'The number of credits required to access this artifact. This is used for metering and billing purposes.'
    ),
  description: z.string().nullable().describe("The artifact's description"),
  mimeType: z
    .string()
    .describe('The MIME type of the artifact, e.g. image/png'),
  name: z
    .string()
    .nullable()
    .describe('The name of the artifact, used for display purposes'),
  organizationId: z.string().describe('The organization name'),
  previewImage: z
    .string()
    .nullable()
    .describe(
      'The URL of the preview image for this artifact. This is used for displaying a thumbnail in the UI.'
    ),
  producerId: z
    .string()
    .nullable()
    .describe('The ID of the run that produced this artifact, if applicable'),
  text: z.string().nullable().describe('The artifact text'),
  url: z.string().nullable().describe('The artifact URL')
}).meta({
  description: 'Schema for Artifact entity',
  id: 'ArtifactEntity'
})

export type ArtifactEntity = z.infer<typeof ArtifactEntitySchema>

export const ARTIFACT_ENTITY_KEY = 'artifacts'
