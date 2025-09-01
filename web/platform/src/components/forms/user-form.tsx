import type { UpdateUserBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import { useGetSessionSuspense, useUpdateUser } from '@archesai/client'
import { USER_ENTITY_KEY } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function UserForm() {
  const { mutateAsync: updateUser } = useUpdateUser()
  const { data: sessionData } = useGetSessionSuspense()

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: sessionData.user.name,
      description: 'Your username',
      label: 'Username',
      name: 'username',
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type='text'
        />
      )
    },
    {
      defaultValue: sessionData.user.name,
      description: 'Your full name',
      label: 'Name',
      name: 'name',
      renderControl: (field) => (
        <Input
          {...field}
          type='text'
        />
      )
    },
    {
      defaultValue: sessionData.user.email,
      description: 'Your email address',
      label: 'Email',
      name: 'email',
      renderControl: (field) => (
        <Input
          {...field}
          type='text'
        />
      )
    }
  ]

  return (
    <GenericForm<never, UpdateUserBody>
      description='View and update your user details'
      entityKey={USER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (data) => {
        await updateUser({
          data,
          id: sessionData.user.id
        })
      }}
      showCard={true}
      title='Profile'
    />
  )
}
