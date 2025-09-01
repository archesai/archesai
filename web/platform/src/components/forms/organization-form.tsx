import type {
  CreateOrganizationBody,
  UpdateOrganizationMutationBody
} from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateOrganization,
  useGetOneOrganizationSuspense,
  useGetSessionSuspense,
  useUpdateOrganization
} from '@archesai/client'
import {
  ORGANIZATION_ENTITY_KEY,
  OrganizationEntitySchema
} from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function OrganizationForm() {
  const {
    data: { session }
  } = useGetSessionSuspense()
  const { mutateAsync: createOrganization } = useCreateOrganization()
  const { mutateAsync: updateOrganization } = useUpdateOrganization()
  const {
    data: { data: organization }
  } = useGetOneOrganizationSuspense(session.activeOrganizationId)

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: organization.name,
      description: OrganizationEntitySchema.shape.name.description ?? '',
      label: 'Name',
      name: 'name',
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type='text'
        />
      )
    },
    {
      defaultValue: organization.billingEmail ?? '',
      description:
        OrganizationEntitySchema.shape.billingEmail.description ?? '',
      label: 'Billing Email',
      name: 'billingEmail',
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type='email'
        />
      )
    }
  ]

  return (
    <GenericForm<CreateOrganizationBody, UpdateOrganizationMutationBody>
      description={"View your organization's details"}
      entityKey={ORGANIZATION_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitCreate={async (createOrganizationDto) => {
        await createOrganization({
          data: createOrganizationDto
        })
      }}
      onSubmitUpdate={async (updateOrganizationDto) => {
        await updateOrganization({
          data: updateOrganizationDto,
          id: session.activeOrganizationId
        })
      }}
      showCard={true}
      title={'Organiation'}
    />
  )
}
