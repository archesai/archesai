'use client'
import { FormFieldConfig, GenericForm } from '@/components/forms/generic-form/generic-form'
import { Input } from '@/components/ui/input'
import {
  useLabelsControllerCreate,
  useLabelsControllerFindOne,
  useLabelsControllerUpdate
} from '@/generated/archesApiComponents'
import { CreateLabelDto, UpdateLabelDto } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import * as z from 'zod'

const formSchema = z.object({
  name: z.string()
})

export default function LabelForm({ labelId }: { labelId?: string }) {
  const { defaultOrgname } = useAuth()
  const { data: label } = useLabelsControllerFindOne(
    {
      pathParams: {
        id: labelId as string,
        orgname: defaultOrgname
      }
    },
    {
      enabled: !!defaultOrgname && !!labelId
    }
  )
  const { mutateAsync: updateLabel } = useLabelsControllerUpdate({})
  const { mutateAsync: createLabel } = useLabelsControllerCreate({})

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: label?.name,
      description: 'This is the name that will be used for this label.',
      label: 'Name',
      name: 'name',
      props: {
        placeholder: 'Label name here...'
      },
      validationRule: formSchema.shape.name
    }
  ]

  return (
    <GenericForm<CreateLabelDto, UpdateLabelDto>
      description={!labelId ? 'Invite a new label' : 'Update an existing label'}
      fields={formFields}
      isUpdateForm={!!labelId}
      itemType='label'
      onSubmitCreate={async (createLabelDto, mutateOptions) => {
        await createLabel(
          {
            body: createLabelDto,
            pathParams: {
              orgname: defaultOrgname
            }
          },
          mutateOptions
        )
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateLabel(
          {
            body: data as any,
            pathParams: {
              id: labelId as string,
              orgname: defaultOrgname
            }
          },
          mutateOptions
        )
      }}
      title='Configuration'
    />
  )
}
