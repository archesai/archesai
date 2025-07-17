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
import { ORGANIZATION_ENTITY_KEY } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function OrganizationForm() {
  const {
    data: { session }
  } = useGetSessionSuspense()
  const { mutateAsync: createOrganization } = useCreateOrganization({})
  const { mutateAsync: updateOrganization } = useUpdateOrganization()
  const { error } = useGetOneOrganizationSuspense(session.activeOrganizationId)
  if (error) {
    return <div>Organization not found</div>
  }

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: session.activeOrganizationId,
      description: 'The name of the organization. This cannot be changed.',
      label: 'Name',
      name: 'name',
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type='text'
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
      title={ORGANIZATION_ENTITY_KEY}
    />
  )
}
