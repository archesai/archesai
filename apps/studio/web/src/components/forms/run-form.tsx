import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Input } from "@archesai/ui";
import type { JSX } from "react";
import type { CreateRunBody, UpdateRunBody } from "#lib/index";
import { useCreateRun, useGetRun, useUpdateRun } from "#lib/index";

export default function RunForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateRun } = useUpdateRun();
  const { mutateAsync: createRun } = useCreateRun();
  const { data: existingRun } = useGetRun(id, { query: { enabled: !!id } });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingRun?.data as Record<string, unknown> | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.pipelineID as string) ?? "",
      description: "UUID identifier",
      label: "Pipeline ID",
      name: "pipelineID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter pipeline id..."
          type="text"
        />
      ),
    },
  ];
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.pipelineID as string) ?? "",
      description: "UUID identifier",
      label: "Pipeline ID",
      name: "pipelineID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter pipeline id..."
          type="text"
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateRunBody, UpdateRunBody>
      description={!id ? "Create a new run" : "Update an existing run"}
      entityKey="runs"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createRun({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateRun({ data: updateDto, id: id });
      }}
      title="Run"
    />
  );
}
