import type { JSX } from "react";
import { useEffect } from "react";
import type { ControllerRenderProps, FieldValues } from "react-hook-form";
import { useForm } from "react-hook-form";

import { Loader2Icon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "#components/shadcn/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "#components/shadcn/form";
import { cn } from "#lib/utils";

export interface FormFieldConfig<T extends FieldValues = FieldValues> {
  defaultValue?: boolean | number | string | undefined;
  description?: string;
  ignoreOnCreate?: boolean;
  label: string;
  name: keyof T;
  renderControl: (field: ControllerRenderProps) => React.ReactNode;
}

type GenericFormProps<
  CreateDto extends FieldValues,
  UpdateDto extends FieldValues,
> = {
  description?: string;
  entityKey: string;
  fields: FormFieldConfig<CreateDto | UpdateDto>[];
  /**
   * When supplied, `mutateOptions` is passed straight through.
   * Use it to wire in TanStack Query’s `useMutation` options and keep side-effects outside.
   */
  mutateOptions?: Record<string, unknown>;
  postContent?: React.ReactNode;
  preContent?: React.ReactNode;
  showCard?: boolean;
  title?: string;
} & (
  | {
      isUpdateForm: false;
      onSubmitCreate: (d: CreateDto) => Promise<void>;
      onSubmitUpdate?: (d: UpdateDto) => Promise<void>;
    }
  | {
      isUpdateForm: true;
      onSubmitCreate?: (d: CreateDto) => Promise<void>;
      onSubmitUpdate: (d: UpdateDto) => Promise<void>;
    }
);

export function GenericForm<
  CreateDto extends FieldValues,
  UpdateDto extends FieldValues,
>(props: GenericFormProps<CreateDto, UpdateDto>): JSX.Element {
  const {
    description,
    fields,
    isUpdateForm,
    // mutateOptions,
    onSubmitCreate,
    onSubmitUpdate,
    showCard = false,
    title,
  } = props;

  /* ---------- memoised defaults & schema ---------- */
  const defaultValues = fields.reduce<Record<string, unknown>>((acc, f) => {
    if (f.defaultValue !== undefined) acc[String(f.name)] = f.defaultValue;
    return acc;
  }, {});

  /* ---------- form instance ---------- */
  const form = useForm({
    defaultValues,
    mode: "onChange",
  });

  /* ---------- keep external defaults in sync ---------- */
  useEffect(() => {
    form.reset(defaultValues);
  }, [defaultValues, form]);

  /* ---------- submit helpers ---------- */
  async function handleSubmit(values: FieldValues) {
    if (isUpdateForm) {
      await onSubmitUpdate(values as UpdateDto);
    } else {
      await onSubmitCreate(values as CreateDto);
    }
  }

  return (
    <Form {...form}>
      <form
        noValidate
        onSubmit={form.handleSubmit(handleSubmit)}
      >
        <Card
          className={cn(!showCard && "border-none shadow-none", "min-w-sm")}
        >
          <CardHeader className="p-4 pb-2">
            {title && <CardTitle className="text-base">{title}</CardTitle>}
            {description && (
              <CardDescription className="text-sm">
                {description}
              </CardDescription>
            )}
          </CardHeader>

          <CardContent className="flex max-h-[60vh] flex-col gap-3 overflow-y-auto p-4 pt-2">
            {props.preContent}
            {fields
              .filter((f) => isUpdateForm || !f.ignoreOnCreate)
              .map((fieldConfig) => (
                <FormField
                  control={form.control}
                  key={String(fieldConfig.name)}
                  name={String(fieldConfig.name)}
                  render={({ field, fieldState }) => (
                    <FormItem className="space-y-1">
                      <FormLabel className="text-sm">
                        {fieldConfig.label}
                      </FormLabel>
                      <FormControl>
                        {fieldConfig.renderControl(field)}
                      </FormControl>
                      {!fieldState.error && fieldConfig.description && (
                        <FormDescription className="text-xs">
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

          <CardFooter className="border-t p-4 pt-3">
            <Button
              // disabled={
              //   !!(
              //     form.formState.isSubmitting ||
              //     !form.formState.isDirty ||
              //     !form.formState.isValid
              //   )
              // }
              size="sm"
              type="submit"
            >
              {form.formState.isSubmitting && (
                <Loader2Icon className="animate-spin" />
              )}
              <span className="capitalize">Submit</span>
            </Button>
            <Button
              // disabled={
              //   !!(form.formState.isSubmitting || !form.formState.isDirty)
              // }
              onClick={() => {
                form.reset();
              }}
              size="sm"
              type="button"
              variant="secondary"
            >
              Reset
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
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
