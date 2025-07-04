import type {
  CreateOrganizationBody,
  UpdateOrganizationBody
} from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateOrganization,
  useGetOneOrganization,
  useUpdateOrganization
} from '@archesai/client'
import { ORGANIZATION_ENTITY_KEY } from '@archesai/domain'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function OrganizationForm({ orgname }: { orgname?: string }) {
  const { mutateAsync: createOrganization } = useCreateOrganization({})
  const { mutateAsync: updateOrganization } = useUpdateOrganization()
  const { error } = useGetOneOrganization(orgname!, {
    query: {
      enabled: !!orgname
    }
  })
  if (error) {
    return <div>Organization not found</div>
  }

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: orgname,
      description: 'The name of the organization. This cannot be changed.',
      label: 'Name',
      name: 'name',
      props: {
        disabled: true
      }
    }
  ]

  return (
    <GenericForm<CreateOrganizationBody, UpdateOrganizationBody>
      description={"View your organization's details"}
      entityKey={ORGANIZATION_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitCreate={async (createOrganizationDto, mutateOptions) => {
        await createOrganization(
          {
            data: createOrganizationDto
          },
          mutateOptions
        )
      }}
      onSubmitUpdate={async (updateOrganizationDto, mutateOptions) => {
        await updateOrganization(
          {
            data: updateOrganizationDto,
            id: orgname!
          },
          mutateOptions
        )
      }}
      showCard={true}
      title={ORGANIZATION_ENTITY_KEY}
    />
  )
}
