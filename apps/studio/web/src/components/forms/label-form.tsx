import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { CreateLabelBody, UpdateLabelBody } from "#lib/index";
import { useCreateLabel, useGetLabel, useUpdateLabel } from "#lib/index";

export default function LabelForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateLabel } = useUpdateLabel();
  const { mutateAsync: createLabel } = useCreateLabel();
  const { data: existingLabel } = useGetLabel(id, { query: { enabled: !!id } });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingLabel?.data as Record<string, unknown> | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the label",
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
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the label",
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
    <GenericForm<CreateLabelBody, UpdateLabelBody>
      description={!id ? "Create a new label" : "Update an existing label"}
      entityKey="labels"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createLabel({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateLabel({ data: updateDto, id: id });
      }}
      title="Label"
    />
  );
}
