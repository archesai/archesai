import { createFileRoute, Link, useRouter } from '@tanstack/react-router'

import type { CreateAccountDto } from '@archesai/schemas'

import { useRegister } from '@archesai/client'
import { CreateAccountDtoSchema } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export const Route = createFileRoute('/auth/register/')({
  component: RegisterPage
})

export default function RegisterPage() {
  const router = useRouter()
  const { mutateAsync: register } = useRegister()

  return (
    <>
      <GenericForm<CreateAccountDto, CreateAccountDto>
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
            validationRule: CreateAccountDtoSchema.properties.email
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
            validationRule: CreateAccountDtoSchema.properties.password
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
            validationRule: CreateAccountDtoSchema.properties.password
          }
        ]}
        isUpdateForm={false}
        onSubmitCreate={async (data) => {
          await register({
            data: {
              email: data.email,
              password: data.password
            }
          })
          await router.navigate({
            to: '/'
          })
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
