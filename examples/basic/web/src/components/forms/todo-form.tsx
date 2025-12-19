import type { FormFieldConfig } from "@archesai/ui";
import { Checkbox, GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { CreateTodoBody, UpdateTodoBody } from "#lib/index";
import { useCreateTodo, useGetTodo, useUpdateTodo } from "#lib/index";

export default function TodoForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateTodo } = useUpdateTodo();
  const { mutateAsync: createTodo } = useCreateTodo();
  const { data: existingTodo } = useGetTodo(id, { query: { enabled: !!id } });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingTodo?.data as Record<string, unknown> | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.completed as boolean) ?? false,
      label: "Completed",
      name: "completed",
      renderControl: (field) => (
        <Checkbox
          checked={field.value as boolean}
          onCheckedChange={field.onChange}
        />
      ),
    },
    {
      defaultValue: (data?.title as string) ?? "",
      label: "Title",
      name: "title",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter title..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.completed as boolean) ?? false,
      label: "Completed",
      name: "completed",
      renderControl: (field) => (
        <Checkbox
          checked={field.value as boolean}
          onCheckedChange={field.onChange}
        />
      ),
    },
    {
      defaultValue: (data?.title as string) ?? "",
      label: "Title",
      name: "title",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter title..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateTodoBody, UpdateTodoBody>
      description={!id ? "Create a new todo" : "Update an existing todo"}
      entityKey="todos"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createTodo({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateTodo({ data: updateDto, id: id });
      }}
      title="Todo"
    />
  );
}
