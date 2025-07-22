import { z } from 'zod'

import { OrganizationEntitySchema } from '#auth/organizations/entities/organization.entity'

export const CreateOrganizationDtoSchema: z.ZodObject<{
  billingEmail: z.ZodNullable<z.ZodString>
  organizationId: z.ZodUUID
}> = z.object({
  billingEmail: OrganizationEntitySchema.shape.billingEmail,
  organizationId: OrganizationEntitySchema.shape.id
})

export type CreateOrganizationDto = z.infer<typeof CreateOrganizationDtoSchema>
