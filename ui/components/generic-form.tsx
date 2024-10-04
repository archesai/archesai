import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ReloadIcon } from "@radix-ui/react-icons";
import { useEffect } from "react";
import { ControllerRenderProps, useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "./ui/button";
import { useToast } from "./ui/use-toast";

export interface FormFieldConfig {
  component: React.ComponentType<any>;
  defaultValue?: any;
  description: string;
  label: string;
  name: string;
  props?: any;
  renderControl?: (
    field: ControllerRenderProps<Record<string, any>>
  ) => React.ReactNode;
  validationRule?: z.ZodType<any, any>;
}

interface GenericFormProps<TCreateVariables, TUpdateVariables> {
  description?: string;
  fields: FormFieldConfig[];
  isUpdateForm: boolean;
  itemType: string;
  onSubmitCreate?: (data: TCreateVariables, mutateOptions: any) => void;
  onSubmitUpdate?: (data: TUpdateVariables, mutateOptions: any) => void;
  title?: string;
}

export function GenericForm<TCreateVariables, TUpdateVariables>({
  description,
  fields,
  isUpdateForm,
  itemType,
  onSubmitCreate,
  onSubmitUpdate,
  title,
}: GenericFormProps<TCreateVariables, TUpdateVariables>) {
  const { toast } = useToast();
  const defaultValues = isUpdateForm
    ? fields.reduce<Record<string, any>>((acc, field) => {
        if (field.defaultValue !== undefined) {
          acc[field.name] = field.defaultValue;
        }
        return acc;
      }, {})
    : {};

  const schema = z.object(
    fields.reduce(
      (acc, field) => {
        if (field.validationRule) {
          acc[field.name] = field.validationRule;
        }
        return acc;
      },
      {} as { [key: string]: z.ZodType<any, any> }
    )
  );

  const form = useForm({
    defaultValues: defaultValues,
    resolver: zodResolver(schema),
  });

  useEffect(() => {
    form.reset(defaultValues);
  }, [fields]);

  return (
    <Card>
      {true ? (
        <CardHeader>
          <CardTitle className="text-xl">{title}</CardTitle>
          <CardDescription>{description}</CardDescription>
        </CardHeader>
      ) : null}

      <CardContent>
        <Form {...form}>
          <form
            className="flex flex-col gap-4"
            noValidate
            onSubmit={form.handleSubmit(
              isUpdateForm
                ? (data) => {
                    onSubmitUpdate?.(data as any, {
                      onError: (error: any) => {
                        toast({
                          description: (error as any)?.stack.msg,
                          title: `Error updating ${itemType}`,
                          variant: "destructive",
                        });
                      },
                      onSuccess: () => {
                        toast({
                          description: `Your ${itemType} has been updated`,
                          title: `${itemType} updated`,
                        });
                        form.reset();
                      },
                    });
                  }
                : (data) => {
                    onSubmitCreate?.(data as any, {
                      onError: (error: any) => {
                        toast({
                          description: (error as any)?.stack.msg,
                          title: `Error creating ${itemType}`,
                          variant: "destructive",
                        });
                      },
                      onSuccess: () => {
                        toast({
                          description: `Your ${itemType} has been created`,
                          title: `${itemType} created`,
                        });
                        form.reset();
                      },
                    });
                  }
            )}
          >
            {fields.map((fieldConfig) => (
              <FormField
                control={form.control}
                key={fieldConfig.name}
                name={fieldConfig.name}
                render={({ field }) => (
                  <FormItem className="flex flex-col col-span-1">
                    <FormLabel>{fieldConfig.label}</FormLabel>
                    <FormControl>
                      {fieldConfig.renderControl ? (
                        fieldConfig.renderControl(field)
                      ) : (
                        <fieldConfig.component
                          {...field}
                          {...fieldConfig.props}
                        />
                      )}
                    </FormControl>
                    {!form.formState.errors[
                      fieldConfig.name
                    ]?.message?.toString() && (
                      <FormDescription>
                        {fieldConfig.description}
                      </FormDescription>
                    )}
                    <FormMessage>
                      {form.formState.errors[
                        fieldConfig.name
                      ]?.message?.toString()}
                    </FormMessage>
                  </FormItem>
                )}
              />
            ))}

            {(onSubmitCreate || onSubmitUpdate) && (
              <div className="flex gap-4 w-full">
                <Button
                  className="w-full"
                  disabled={
                    form.formState.isSubmitting || !form.formState.isDirty
                  }
                  size="sm"
                  type="submit"
                  variant={"secondary"}
                >
                  {form.formState.isSubmitting && (
                    <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Submit
                </Button>
                <Button
                  className="w-full"
                  disabled={
                    form.formState.isSubmitting || !form.formState.isDirty
                  }
                  onClick={() => {
                    form.reset();
                  }}
                  size="sm"
                  variant={"secondary"}
                >
                  Clear
                </Button>
              </div>
            )}
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
