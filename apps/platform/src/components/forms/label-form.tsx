import type { CreateLabelBody, UpdateLabelBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateLabel,
  useGetOneLabel,
  useUpdateLabel
} from '@archesai/client'
import { LABEL_ENTITY_KEY, Type } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function LabelForm({ labelId }: { labelId?: string }) {
  const { mutateAsync: updateLabel } = useUpdateLabel()
  const { mutateAsync: createLabel } = useCreateLabel()
  const { data: label } = useGetOneLabel(labelId, {
    query: {
      enabled: !!labelId,
      placeholderData: {
        data: {
          createdAt: '',
          id: labelId ?? '',
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
      props: {
        placeholder: 'Label name here...'
      },
      renderControl: (field) => (
        <Input
          {...field}
          type='text'
        />
      ),
      validationRule: Type.String({
        minLength: 1
      })
    }
  ]

  return (
    <GenericForm<CreateLabelBody, UpdateLabelBody>
      description={!labelId ? 'Invite a new label' : 'Update an existing label'}
      entityKey={LABEL_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!labelId}
      onSubmitCreate={async (createLabelDto) => {
        await createLabel({
          data: createLabelDto
        })
      }}
      onSubmitUpdate={async (data) => {
        await updateLabel({
          data: data,
          id: labelId!
        })
      }}
      title='Configuration'
    />
  )
}
