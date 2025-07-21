import { z } from 'zod'

import { ArtifactEntitySchema } from '#orchestration/artifacts/entities/artifact.entity'

export const CreateArtifactDtoSchema: z.ZodObject<{
  name: z.ZodString
  text: z.ZodNullable<z.ZodString>
  url: z.ZodNullable<z.ZodString>
}> = z.object({
  name: z.string().describe('The name of the artifact'),
  text: ArtifactEntitySchema.shape.text,
  url: ArtifactEntitySchema.shape.url
})

export type CreateArtifactDto = z.infer<typeof CreateArtifactDtoSchema>
