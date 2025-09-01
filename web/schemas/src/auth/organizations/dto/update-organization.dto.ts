import type { z } from 'zod'

import { CreateOrganizationDtoSchema } from '#auth/organizations/dto/create-organization.dto'

export const UpdateOrganizationDtoSchema: z.ZodObject<{
  billingEmail: z.ZodOptional<z.ZodNullable<z.ZodString>>
  organizationId: z.ZodOptional<z.ZodUUID>
}> = CreateOrganizationDtoSchema.partial()

export type UpdateOrganizationDto = z.infer<typeof UpdateOrganizationDtoSchema>
