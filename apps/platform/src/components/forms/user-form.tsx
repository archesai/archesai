import { Type } from '@sinclair/typebox'

import type { CreateUserBody, UpdateUserBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import { useGetOneUser, useUpdateUser } from '@archesai/client'
import { USER_ENTITY_KEY } from '@archesai/domain'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export default function UserForm() {
  const { mutateAsync: updateUser } = useUpdateUser()
  const { data: userResponse, error } = useGetOneUser('user-form')

  if (error || !userResponse) {
    return <div>Run not found</div>
  }
  const user = userResponse.data

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: user.attributes.name,
      description: 'Your first name',
      label: 'Name',
      name: 'firstName',
      validationRule: Type.String({
        minLength: 1
      })
    },

    {
      component: Input,
      defaultValue: user.attributes.orgname,
      description: 'Your username',
      label: 'Username',
      name: 'username',
      props: {
        disabled: true
      }
    },
    {
      component: Input,
      defaultValue: user.attributes.email,
      description: 'Your email address',
      label: 'Email',
      name: 'email',
      props: {
        disabled: true
      }
    }
  ]

  return (
    <GenericForm<CreateUserBody, UpdateUserBody>
      description='View and update your user details'
      entityKey={USER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateUser(
          {
            data,
            id: user.id
          },
          mutateOptions
        )
      }}
      showCard={true}
      title='Profile'
    />
  )
}
