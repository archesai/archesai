import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'
import { createFileRoute, Link, useRouter } from '@tanstack/react-router'

import { useRegister } from '@archesai/client'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

const RegisterSchema = Type.Object({
  confirmPassword: Type.String({ minLength: 8 }),
  email: Type.String({ format: 'email' }),
  password: Type.String({ minLength: 8 })
})

export const Route = createFileRoute('/auth/register/')({
  component: RegisterPage
})

export default function RegisterPage() {
  const router = useRouter()
  const { mutateAsync: register } = useRegister()

  const onSubmit = async (data: Static<typeof RegisterSchema>) => {
    await register({
      data: {
        email: data.email,
        password: data.password
      }
    })
    await router.navigate({
      to: '/chat'
    })
  }

  return (
    <>
      <GenericForm
        description='Create your account by entering your email and password'
        entityKey='auth'
        fields={[
          {
            defaultValue: '',
            label: 'Email',
            name: 'email',
            renderControl: (field) => (
              <Input
                {...field}
                type='email'
              />
            ),
            validationRule: RegisterSchema.properties.email
          },
          {
            defaultValue: '',
            label: 'Password',
            name: 'password',
            renderControl: (field) => (
              <Input
                {...field}
                type='password'
              />
            ),
            validationRule: RegisterSchema.properties.password
          },
          {
            defaultValue: '',
            label: 'Confirm Password',
            name: 'confirmPassword',
            renderControl: (field) => (
              <Input
                {...field}
                type='password'
              />
            ),
            validationRule: Type.String({
              errorMessage: 'Passwords must match',
              minLength: 8
            })
          }
        ]}
        isUpdateForm={false}
        onSubmitCreate={async (data) => {
          await onSubmit(data as Static<typeof RegisterSchema>)
        }}
        postContent={
          <div className='text-center text-sm'>
            Already have an account?{' '}
            <Link
              className='underline'
              to='/auth/login'
            >
              Login
            </Link>
          </div>
        }
        showCard={true}
        title='Register'
      />
    </>
  )
}
