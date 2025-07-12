import type { CreateUserBody, UpdateUserBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useGetOneUser,
  useGetSessionSuspense,
  useUpdateUser
} from '@archesai/client'
import { USER_ENTITY_KEY } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function UserForm() {
  const { mutateAsync: updateUser } = useUpdateUser()
  const { data: session } = useGetSessionSuspense()
  const { data: user } = useGetOneUser(session.id)

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: user?.data.attributes.name,
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
      defaultValue: user?.data.attributes.email,
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
    <GenericForm<CreateUserBody, UpdateUserBody>
      description='View and update your user details'
      entityKey={USER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (data) => {
        await updateUser({
          data,
          id: user?.data.id
        })
      }}
      showCard={true}
      title='Profile'
    />
  )
}
