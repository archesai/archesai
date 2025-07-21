import type { CreateLabelBody, UpdateLabelBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateLabel,
  useGetOneLabel,
  useUpdateLabel
} from '@archesai/client'
import { LABEL_ENTITY_KEY, StringSchema } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function LabelForm({ id }: { id?: string }) {
  const { mutateAsync: updateLabel } = useUpdateLabel()
  const { mutateAsync: createLabel } = useCreateLabel()
  const { data: label } = useGetOneLabel(id, {
    query: {
      enabled: !!id,
      placeholderData: {
        data: {
          createdAt: '',
          id: id ?? '',
          name: '',
          organizationId: '',
          updatedAt: ''
        }
      }
    }
  })
  if (!label) {
    return null
  }

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: label.data.name,
      description: 'This is the name that will be used for this label.',
      label: 'Name',
      name: 'name',
      renderControl: (field) => (
        <Input
          placeholder='Label name here...'
          {...field}
          type='text'
        />
      ),
      validationRule: StringSchema
    }
  ]

  return (
    <GenericForm<CreateLabelBody, UpdateLabelBody>
      description={!id ? 'Invite a new label' : 'Update an existing label'}
      entityKey={LABEL_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createLabelDto) => {
        await createLabel({
          data: createLabelDto
        })
      }}
      onSubmitUpdate={async (data) => {
        await updateLabel({
          data: data,
          id: id
        })
      }}
      title='Configuration'
    />
  )
}
