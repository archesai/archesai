import type { ControllerRenderProps, FieldValues } from 'react-hook-form'

import { useEffect } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

import { Loader2Icon } from '#components/custom/icons'
import { Button } from '#components/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
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

export interface FormFieldConfig {
  defaultValue?: boolean | number | string | undefined
  description?: string
  ignoreOnCreate?: boolean
  label: string
  name: string
  renderControl: (field: ControllerRenderProps) => React.ReactNode
  validationRule?: z.ZodType
}

type GenericFormProps<
  CreateDto extends FieldValues,
  UpdateDto extends FieldValues
> = {
  description?: string
  entityKey: string
  fields: FormFieldConfig[]
  /**
   * When supplied, `mutateOptions` is passed straight through.
   * Use it to wire in TanStack Query’s `useMutation` options and keep side-effects outside.
   */
  mutateOptions?: Record<string, unknown>
  postContent?: React.ReactNode
  preContent?: React.ReactNode
  showCard?: boolean
  title?: string
} & (
  | {
      isUpdateForm: false
      onSubmitCreate: (d: CreateDto) => Promise<void>
      onSubmitUpdate?: (d: UpdateDto) => Promise<void>
    }
  | {
      isUpdateForm: true
      onSubmitCreate?: (d: CreateDto) => Promise<void>
      onSubmitUpdate: (d: UpdateDto) => Promise<void>
    }
)

export function GenericForm<
  CreateDto extends FieldValues,
  UpdateDto extends FieldValues
>(props: GenericFormProps<CreateDto, UpdateDto>) {
  const {
    description,
    entityKey,
    fields,
    isUpdateForm,
    // mutateOptions,
    onSubmitCreate,
    onSubmitUpdate,
    showCard = false,
    title
  } = props

  /* ---------- memoised defaults & schema ---------- */
  const defaultValues = fields.reduce<Record<string, unknown>>((acc, f) => {
    if (f.defaultValue !== undefined) acc[f.name] = f.defaultValue
    return acc
  }, {})

  const schema = z.object(
    fields.reduce<Record<string, z.ZodType>>((acc, f) => {
      if (f.validationRule) acc[f.name] = f.validationRule
      return acc
    }, {})
  )

  /* ---------- form instance ---------- */
  const form = useForm({
    defaultValues,
    mode: 'onChange',
    resolver: zodResolver(schema)
  })

  /* ---------- keep external defaults in sync ---------- */
  useEffect(() => {
    form.reset(defaultValues)
  }, [defaultValues, form])

  /* ---------- submit helpers ---------- */
  async function handleSubmit(values: FieldValues) {
    const run =
      isUpdateForm ?
        async () => {
          await onSubmitUpdate(values as UpdateDto)
        }
      : async () => {
          await onSubmitCreate(values as CreateDto)
        }

    await run()
  }

  return (
    <Form {...form}>
      <form
        noValidate
        onSubmit={form.handleSubmit(handleSubmit)}
      >
        <Card className={cn(!showCard && 'border-none shadow-none')}>
          <CardHeader>
            {title && <CardTitle>{title}</CardTitle>}
            {description && <CardDescription>{description}</CardDescription>}
          </CardHeader>

          <Separator />

          <CardContent className='flex flex-col gap-6 p-4'>
            {props.preContent}
            {fields
              .filter((f) => isUpdateForm || !f.ignoreOnCreate)
              .map((fieldConfig) => (
                <FormField
                  control={form.control}
                  key={fieldConfig.name}
                  name={fieldConfig.name}
                  render={({ field, fieldState }) => (
                    <FormItem>
                      <FormLabel>{fieldConfig.label}</FormLabel>
                      <FormControl>
                        {fieldConfig.renderControl(field)}
                      </FormControl>
                      {!fieldState.error && (
                        <FormDescription>
                          {fieldConfig.description}
                        </FormDescription>
                      )}
                      <FormMessage />
                    </FormItem>
                  )}
                />
              ))}
            {props.postContent}
          </CardContent>

          <Separator />

          <CardFooter>
            <Button
              // disabled={
              //   !!(
              //     form.formState.isSubmitting ||
              //     !form.formState.isDirty ||
              //     !form.formState.isValid
              //   )
              // }
              size='sm'
              type='submit'
            >
              {form.formState.isSubmitting && (
                <Loader2Icon className='animate-spin' />
              )}
              <span className='capitalize'>
                {isUpdateForm ? 'Update' : 'Create'} {entityKey}
              </span>
            </Button>
            <Button
              // disabled={
              //   !!(form.formState.isSubmitting || !form.formState.isDirty)
              // }
              onClick={() => {
                form.reset()
              }}
              size='sm'
              type='button'
              variant='secondary'
            >
              Reset
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  )
}

// <FormItem>
//                 <FormControl>
//                   <FloatingLabelInput
//                     id={field.label}
//                     label={field.label}
//                     type={field.type ?? 'text'}
//                     {...f}
//                     value={
//                       f.value as
//                         | number
//                         | readonly string[]
//                         | string
//                         | undefined
//                     }
//                   />
//                 </FormControl>
//                 <FormMessage />
//               </FormItem>
