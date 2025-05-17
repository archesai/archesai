'use client'

import { Type } from '@sinclair/typebox'

import type { RegisterBody } from '@archesai/client'
import type { AccountEntity } from '@archesai/domain'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'
import { useAuth } from '@archesai/ui/hooks/use-auth'

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
  const { registerWithEmailAndPassword } = useAuth()

  const formFields: FormFieldConfig<AccountEntity>[] = [
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
    <GenericForm<AccountEntity, RegisterBody, never>
      description={"Configure your member's settings"}
      entityKey='register'
      fields={formFields}
      isUpdateForm={false}
      onSubmitCreate={async (registerDto) => {
        await registerWithEmailAndPassword(
          registerDto.email,
          registerDto.password
        )
      }}
      title='Configuration'
    />
  )
}
