import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { CreateArtifactBody, UpdateArtifactBody } from "#lib/index";
import {
  useCreateArtifact,
  useGetArtifact,
  useUpdateArtifact,
} from "#lib/index";

export default function ArtifactForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateArtifact } = useUpdateArtifact();
  const { mutateAsync: createArtifact } = useCreateArtifact();
  const { data: existingArtifact } = useGetArtifact(id, {
    query: { enabled: !!id },
  });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingArtifact?.data as Record<string, unknown> | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the artifact",
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
      defaultValue: (data?.text as string) ?? "",
      description: "The artifact text",
      label: "Text",
      name: "text",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter text..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the artifact",
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
      defaultValue: (data?.text as string) ?? "",
      description: "The artifact text",
      label: "Text",
      name: "text",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter text..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.url as string) ?? "",
      description: "The artifact URL",
      label: "URL",
      name: "url",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter url..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateArtifactBody, UpdateArtifactBody>
      description={
        !id ? "Create a new artifact" : "Update an existing artifact"
      }
      entityKey="artifacts"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createArtifact({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateArtifact({ data: updateDto, id: id });
      }}
      title="Artifact"
    />
  );
}
