'use client'

import { Type } from '@sinclair/typebox'

import type { CreateLabelBody, UpdateLabelBody } from '@archesai/client'
import type { LabelEntity } from '@archesai/domain'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateLabel,
  useGetOneLabel,
  useUpdateLabel
} from '@archesai/client'
import { LABEL_ENTITY_KEY } from '@archesai/domain'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function LabelForm({ labelId }: { labelId?: string }) {
  const { mutateAsync: updateLabel } = useUpdateLabel({})
  const { mutateAsync: createLabel } = useCreateLabel({})
  const { data: existingLabelResponse } = useGetOneLabel(labelId!, {
    query: {
      enabled: !!labelId
    }
  })

  if (existingLabelResponse?.status !== 200) {
    return <div>Label not found</div>
  }
  const label = existingLabelResponse.data.data

  const formFields: FormFieldConfig<LabelEntity>[] = [
    {
      component: Input,
      defaultValue: label.attributes.name,
      description: 'This is the name that will be used for this label.',
      label: 'Name',
      name: 'name',
      props: {
        placeholder: 'Label name here...'
      },
      validationRule: Type.String({})
    }
  ]

  return (
    <GenericForm<LabelEntity, CreateLabelBody, UpdateLabelBody>
      description={!labelId ? 'Invite a new label' : 'Update an existing label'}
      entityKey={LABEL_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!labelId}
      onSubmitCreate={async (createLabelDto, mutateOptions) => {
        await createLabel(
          {
            data: createLabelDto
          },
          mutateOptions
        )
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateLabel(
          {
            data: data,
            id: labelId!
          },
          mutateOptions
        )
      }}
      title='Configuration'
    />
  )
}
