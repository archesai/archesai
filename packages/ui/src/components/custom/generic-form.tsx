import type { TSchema } from '@sinclair/typebox'
import type { ControllerRenderProps, FieldValues } from 'react-hook-form'

import { useEffect } from 'react'
import { typeboxResolver } from '@hookform/resolvers/typebox'
import { Type } from '@sinclair/typebox'
import { LoaderIcon } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'

import type { BaseEntity } from '@archesai/domain'

import { Button } from '#components/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '#components/shadcn/card'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '#components/shadcn/form'
import { Separator } from '#components/shadcn/separator'
import { cn } from '#lib/utils'

export interface FormFieldConfig<TEntity extends BaseEntity> {
  component: React.ComponentType
  defaultValue?: boolean | number | string | undefined
  description: string
  ignoreOnCreate?: boolean
  label: string
  name: string
  props?: Record<string, unknown>
  renderControl?: (field: ControllerRenderProps<TEntity>) => React.ReactNode
  validationRule?: TSchema
}

type GenericFormProps<
  TEntity extends BaseEntity,
  CreateDto extends FieldValues,
  UpdateDto extends FieldValues
> = {
  description?: string
  entityKey: string
  fields: FormFieldConfig<TEntity>[]
  onSubmitCreate?: (
    data: CreateDto,
    mutateOptions: Record<string, unknown>
  ) => void
  showCard?: boolean
  title?: string
} & (
  | {
      isUpdateForm: false
      onSubmitUpdate?: (
        data: UpdateDto,
        mutateOptions: Record<string, unknown>
      ) => void
    }
  | {
      isUpdateForm: true
      onSubmitUpdate: (
        data: UpdateDto,
        mutateOptions: Record<string, unknown>
      ) => void
    }
)

export function GenericForm<
  TEntity extends BaseEntity,
  CreateDto extends FieldValues,
  UpdateDto extends FieldValues
>({
  description,
  entityKey,
  fields,
  isUpdateForm,
  onSubmitCreate,
  onSubmitUpdate,
  showCard = false,
  title
}: GenericFormProps<TEntity, CreateDto, UpdateDto>) {
  const defaultValues = fields.reduce<Record<string, unknown>>((acc, field) => {
    if (field.defaultValue !== undefined) {
      acc[field.name] = field.defaultValue
    }
    return acc
  }, {})

  const schema = Type.Object(
    fields.reduce<Record<string, TSchema>>((acc, field) => {
      if (field.validationRule) {
        acc[field.name] = field.validationRule
      }
      return acc
    }, {})
  )

  const form = useForm({
    defaultValues: defaultValues,
    resolver: typeboxResolver(schema)
  })

  useEffect(() => {
    form.reset(defaultValues)
  }, [fields.map((f) => f.defaultValue).join()])

  return (
    <Card
      className={cn(
        'flex flex-1 flex-col',
        showCard ? '' : 'border-none shadow-none'
      )}
    >
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <Separator />
      <Form {...form}>
        <form
          className='flex flex-1 flex-col'
          noValidate
          onSubmit={form.handleSubmit(
            isUpdateForm ?
              (data) => {
                onSubmitUpdate(data as UpdateDto, {
                  onError: (error: Error) => {
                    toast('`Update failed`', {
                      description: error.message
                    })
                  },
                  onSuccess: () => {
                    toast(`Update successful`, {
                      description: `Your ${entityKey} has been updated`
                    })
                  }
                })
              }
            : onSubmitCreate ?
              (data) => {
                onSubmitCreate(data as CreateDto, {
                  onError: (error: Error) => {
                    toast(`Create failed`, {
                      description: error.message
                    })
                  },
                  onSuccess: () => {
                    toast(`Creation successful`, {
                      description: `Your ${entityKey} has been created`
                    })
                  }
                })
              }
            : () => {
                toast(`Error`, {
                  description: `No submit function provided`
                })
              }
          )}
        >
          <CardContent className='flex flex-1 flex-col gap-4 p-4'>
            {fields
              .filter((f) => isUpdateForm || !f.ignoreOnCreate)
              .map((fieldConfig) => (
                <FormField
                  control={form.control}
                  key={fieldConfig.name}
                  name={fieldConfig.name}
                  render={({ field, fieldState }) => {
                    return (
                      <FormItem className='col-span-1 flex flex-col'>
                        <FormLabel>{fieldConfig.label}</FormLabel>
                        <FormControl>
                          {
                            fieldConfig.renderControl ?
                              // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-explicit-any
                              fieldConfig.renderControl(field.value as any)
                              // <fieldConfig.component
                              //   {...field}
                              //   {...fieldConfig.props}
                              //   value={field.value}
                              // /> // FIXME
                            : <></>
                          }
                        </FormControl>
                        {!fieldState.error?.message && (
                          <FormDescription>
                            {fieldConfig.description}
                          </FormDescription>
                        )}
                        <FormMessage>{fieldState.error?.message}</FormMessage>
                      </FormItem>
                    )
                  }}
                />
              ))}
          </CardContent>
          <Separator />
          <div className='flex justify-end rounded-xl p-4 py-2'>
            {(onSubmitCreate ?? onSubmitUpdate) && (
              <div className='flex w-full items-center justify-end gap-2'>
                <Button
                  className='flex flex-1 gap-2'
                  disabled={
                    form.formState.isSubmitting || !form.formState.isDirty
                  }
                  size='sm'
                  type='submit'
                >
                  {form.formState.isSubmitting && (
                    <LoaderIcon className='h-5 w-5 animate-spin' />
                  )}
                  <span className='capitalize'>
                    {isUpdateForm ? 'Update' : 'Create'} {entityKey}
                  </span>
                </Button>
                <Button
                  className='flex-1'
                  disabled={
                    form.formState.isSubmitting || !form.formState.isDirty
                  }
                  onClick={() => {
                    form.reset()
                  }}
                  size='sm'
                  variant={'secondary'}
                >
                  Clear
                </Button>
              </div>
            )}
          </div>
        </form>
      </Form>
    </Card>
  )
}
