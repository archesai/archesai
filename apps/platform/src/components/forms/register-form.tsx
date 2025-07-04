import { Type } from '@sinclair/typebox'
import { FormatRegistry } from '@sinclair/typebox/type'

import type { RegisterBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import { register } from '@archesai/client'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

FormatRegistry.Set('email', (value: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
)

export const RegisterSchema = Type.Object({
  confirmPassword: Type.String({
    maxLength: 128,
    minLength: 8
  }),
  email: Type.String({
    format: 'email'
  }),
  password: Type.String({
    maxLength: 128,
    minLength: 8
  })
})

export default function RegisterForm() {
  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: '',
      description: 'This is the email that will be used for this user.',
      label: 'Email',
      name: 'email',
      props: {
        placeholder: 'Enter your email here...'
      },
      validationRule: RegisterSchema.properties.email
    },
    {
      component: Input,
      defaultValue: '',
      description: 'This is the password that will be used for your account.',
      label: 'Password',
      name: 'password',
      props: {
        placeholder: 'Enter your password here...'
      },
      validationRule: RegisterSchema.properties.password
    },
    {
      component: Input,
      defaultValue: '',
      description:
        'This is the role that will be used for this member. Note that different roles have different permissions.',
      label: 'Confirm Password',
      name: 'confirmPassword',
      props: {
        placeholder: 'Confirm your password here...'
      },
      validationRule: RegisterSchema.properties.confirmPassword
    }
  ]

  return (
    <GenericForm<RegisterBody, never>
      description={"Configure your member's settings"}
      entityKey='register'
      fields={formFields}
      isUpdateForm={false}
      onSubmitCreate={async (registerDto) => {
        await register({
          email: registerDto.email,
          password: registerDto.password
        })
      }}
      title='Configuration'
    />
  )
}
