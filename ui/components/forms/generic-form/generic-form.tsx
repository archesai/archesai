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
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { ReloadIcon } from "@radix-ui/react-icons";
import { useEffect } from "react";
import { ControllerRenderProps, useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "../../ui/button";
import { Separator } from "../../ui/separator";
import { useToast } from "../../ui/use-toast";

export interface FormFieldConfig {
  component: React.ComponentType<any>;
  defaultValue?: any;
  description: string;
  ignoreOnCreate?: boolean;
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
  showCard?: boolean;
  title?: string;
}

export function GenericForm<TCreateVariables, TUpdateVariables>({
  description,
  fields,
  isUpdateForm,
  itemType,
  onSubmitCreate,
  onSubmitUpdate,
  showCard = false,
  title,
}: GenericFormProps<TCreateVariables, TUpdateVariables>) {
  const { toast } = useToast();
  const defaultValues = fields.reduce<Record<string, any>>((acc, field) => {
    if (field.defaultValue !== undefined) {
      acc[field.name] = field.defaultValue;
    }
    return acc;
  }, {});

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
  }, [fields.map((f) => f.defaultValue).join()]);

  console.log(form.formState.errors);
  console.log(form.getValues()); // FIXME
  return (
    <Card
      className={cn(
        "flex flex-1 flex-col",
        showCard ? "" : "border-none shadow-none"
      )}
    >
      <CardHeader className="rounded-t-lg bg-sidebar">
        <CardTitle className="text-lg">{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <Separator />
      <Form {...form}>
        <form
          className="flex flex-1 flex-col"
          noValidate
          onSubmit={form.handleSubmit(
            isUpdateForm
              ? (data) => {
                  onSubmitUpdate?.(data as any, {
                    onError: (error: any) => {
                      toast({
                        description: (error as any)?.stack.message,
                        title: `Update failed`,
                        variant: "destructive",
                      });
                    },
                    onSuccess: () => {
                      toast({
                        description: `Your ${itemType} has been updated`,
                        title: `Update successful`,
                      });
                    },
                  });
                }
              : (data) => {
                  onSubmitCreate?.(data as any, {
                    onError: (error: any) => {
                      toast({
                        description: (error as any)?.stack.message,
                        title: `Create failed`,
                        variant: "destructive",
                      });
                    },
                    onSuccess: () => {
                      toast({
                        description: `Your ${itemType} has been created`,
                        title: `Creation successful`,
                      });
                    },
                  });
                }
          )}
        >
          <CardContent className="flex flex-1 flex-col gap-4 bg-muted/30 p-4">
            {fields
              .filter((f) => isUpdateForm || !f.ignoreOnCreate)
              .map((fieldConfig) => (
                <FormField
                  control={form.control}
                  key={fieldConfig.name}
                  name={fieldConfig.name}
                  render={({ field }) => (
                    <FormItem className="col-span-1 flex flex-col">
                      <FormLabel>{fieldConfig.label}</FormLabel>
                      <FormControl>
                        {fieldConfig.renderControl ? (
                          fieldConfig.renderControl(field)
                        ) : (
                          <fieldConfig.component
                            {...field}
                            {...fieldConfig.props}
                            value={field.value || ""}
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
          </CardContent>
          <Separator />
          <div className="flex rounded-b-lg bg-sidebar p-4">
            {(onSubmitCreate || onSubmitUpdate) && (
              <div className="flex w-full items-center justify-end gap-2">
                <Button
                  className="flex flex-1 gap-2"
                  disabled={
                    form.formState.isSubmitting || !form.formState.isDirty
                  }
                  size="sm"
                  type="submit"
                >
                  {form.formState.isSubmitting && (
                    <ReloadIcon className="h-5 w-5 animate-spin" />
                  )}
                  <span className="capitalize">
                    {isUpdateForm ? "Update" : "Create"} {itemType}
                  </span>
                </Button>
                <Button
                  className="flex-1"
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
          </div>
        </form>
      </Form>
    </Card>
  );
}
