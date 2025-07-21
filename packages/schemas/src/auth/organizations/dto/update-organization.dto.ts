import type { z } from 'zod'

import { CreateOrganizationDtoSchema } from '#auth/organizations/dto/create-organization.dto'

export const UpdateOrganizationDtoSchema: z.ZodObject<{
  billingEmail: z.ZodOptional<z.ZodString>
  organizationId: z.ZodOptional<z.ZodString>
}> = CreateOrganizationDtoSchema.partial()

export type UpdateOrganizationDto = z.infer<typeof UpdateOrganizationDtoSchema>
