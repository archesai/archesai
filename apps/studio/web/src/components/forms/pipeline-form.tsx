import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Input, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { CreatePipelineBody, UpdatePipelineBody } from "#lib/index";
import {
  useCreatePipeline,
  useGetPipelineSuspense,
  useUpdatePipeline,
} from "#lib/index";

export default function PipelineForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updatePipeline } = useUpdatePipeline();
  const { mutateAsync: createPipeline } = useCreatePipeline();
  const { data: existingPipeline } = useGetPipelineSuspense(id);
  const data = existingPipeline?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.description ?? "",
      description: "Detailed description of the pipeline's purpose",
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
      defaultValue: data?.name ?? "",
      description: "The pipeline's display name",
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
    {
      defaultValue: data?.organizationID ?? "",
      description: "The organization identifier",
      label: "Organization ID",
      name: "organizationID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter organization id..."
          type="text"
        />
      ),
    },
  ];
  return (
    <GenericForm<CreatePipelineBody, UpdatePipelineBody>
      description={
        !id ? "Create a new pipeline" : "Update an existing pipeline"
      }
      entityKey="pipelines"
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createPipeline({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updatePipeline({ data: updateDto, id: id });
      }}
      title="Pipeline"
    />
  );
}
