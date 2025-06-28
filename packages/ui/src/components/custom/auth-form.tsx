import type { TSchema } from '@sinclair/typebox'

import { useMemo } from 'react'
import { typeboxResolver } from '@hookform/resolvers/typebox'
import { FormatRegistry, Type } from '@sinclair/typebox'
import { useForm } from 'react-hook-form'

import { FloatingLabelInput } from '#components/custom/floating-label'
import { Button } from '#components/shadcn/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage
} from '#components/shadcn/form'

FormatRegistry.Set('email', (value: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
)

export function AuthForm({
  description,
  fields,
  onSubmit,
  title
}: {
  description?: string
  fields: {
    defaultValue: boolean | number | string
    label: string
    name: string
    type?: 'email' | 'password' | 'text'
    validationRule?: TSchema
  }[]
  onSubmit: (data: Record<string, unknown>) => void
  title?: string
}) {
  const defaultValues = useMemo(
    () =>
      fields.reduce<Record<string, unknown>>((acc, field) => {
        acc[field.name] = field.defaultValue
        return acc
      }, {}),
    [fields]
  )

  const schema = useMemo(
    () =>
      Type.Object(
        fields.reduce<Record<string, TSchema>>((acc, field) => {
          if (field.validationRule) {
            acc[field.name] = field.validationRule
          }
          return acc
        }, {})
      ),
    [fields]
  )

  const form = useForm({
    defaultValues,
    resolver: typeboxResolver(schema)
  })

  return (
    <div className='flex w-full flex-col gap-4'>
      <div className='text-center'>
        <h1 className='text-2xl font-semibold tracking-tight'>{title}</h1>
        <p className='text-sm text-muted-foreground'>{description}</p>
      </div>
      <Form {...form}>
        <form
          className='flex flex-col gap-2'
          noValidate
          onSubmit={form.handleSubmit(onSubmit)}
        >
          {fields.map((field) => (
            <FormField
              control={form.control}
              key={String(field.name)}
              name={String(field.name)}
              render={({ field: f }) => (
                <FormItem>
                  <FormControl>
                    <FloatingLabelInput
                      id={field.label}
                      label={field.label}
                      type={field.type ?? 'text'}
                      {...f}
                      value={
                        f.value as
                          | number
                          | readonly string[]
                          | string
                          | undefined
                      }
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          ))}
          {form.formState.errors.root && (
            <div className='text-center text-red-600'>
              {form.formState.errors.root.message}
            </div>
          )}
          <Button
            disabled={form.formState.isSubmitting}
            type='submit'
          >
            Submit
          </Button>
        </form>
      </Form>
    </div>
  )
}
