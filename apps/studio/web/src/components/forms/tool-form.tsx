import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { CreateToolBody, UpdateToolBody } from "#lib/index";
import { useCreateTool, useGetTool, useUpdateTool } from "#lib/index";

export default function ToolForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateTool } = useUpdateTool();
  const { mutateAsync: createTool } = useCreateTool();
  const { data: existingTool } = useGetTool(id, { query: { enabled: !!id } });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingTool?.data as Record<string, unknown> | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.description as string) ?? "",
      description: "The tool description",
      label: "Description",
      name: "description",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter description..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the tool",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter name..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.description as string) ?? "",
      description: "The tool description",
      label: "Description",
      name: "description",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter description..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the tool",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter name..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateToolBody, UpdateToolBody>
      description={!id ? "Create a new tool" : "Update an existing tool"}
      entityKey="tools"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createTool({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateTool({ data: updateDto, id: id });
      }}
      title="Tool"
    />
  );
}
