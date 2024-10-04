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
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

export interface FormFieldConfig {
  component: React.ComponentType<any>;
  defaultValue?: any;
  description: string;
  label: string;
  name: string;
  props?: any;
  renderControl?: (
    field: any,
    updateValue: (value: any) => void
  ) => React.ReactNode;
  validationRule?: z.ZodType<any, any>;
}

interface CustomCardFormProps {
  description?: string;
  fields: FormFieldConfig[];
  title?: string;
}

export const CustomCardForm: React.FC<CustomCardFormProps> = ({
  description,
  fields,
  title,
}) => {
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

  const onSubmit: SubmitHandler<Record<string, any>> = (data) => {
    console.log(data);
  };

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
            onSubmit={form.handleSubmit(onSubmit)}
          >
            {fields.map((fieldConfig) => (
              <FormField
                control={form.control}
                key={fieldConfig.name}
                name={fieldConfig.name}
                render={({ field }) =>
                  fieldConfig.renderControl ? (
                    // fieldConfig.renderControl(field, (value) =>
                    //   form.setValue(fieldConfig.name, value)
                    // )
                    <></>
                  ) : (
                    <FormItem className="flex flex-col col-span-1">
                      <FormLabel className="font-semibold text-sm">
                        {fieldConfig.label}
                      </FormLabel>
                      <FormControl>
                        <fieldConfig.component
                          {...field}
                          {...fieldConfig.props}
                        />
                      </FormControl>
                      <FormDescription className="text-sm">
                        {fieldConfig.description}
                      </FormDescription>
                    </FormItem>
                  )
                }
              />
            ))}
          </form>
        </Form>
      </CardContent>
    </Card>
  );
};
