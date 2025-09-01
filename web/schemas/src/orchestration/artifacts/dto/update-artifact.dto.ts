import type { z } from 'zod'

import { CreateArtifactDtoSchema } from '#orchestration/artifacts/dto/create-artifact.dto'

export const UpdateArtifactDtoSchema: z.ZodObject<{
  name: z.ZodOptional<z.ZodString>
  text: z.ZodOptional<z.ZodNullable<z.ZodString>>
  url: z.ZodOptional<z.ZodNullable<z.ZodString>>
}> = CreateArtifactDtoSchema.partial()

export type UpdateArtifactDto = z.infer<typeof UpdateArtifactDtoSchema>
