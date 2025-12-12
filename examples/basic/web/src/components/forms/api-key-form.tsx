import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Input, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { CreateAPIKeyBody, UpdateAPIKeyBody } from "#lib/index";
import { useCreateAPIKey, useGetAPIKey, useUpdateAPIKey } from "#lib/index";

export default function APIKeyForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateAPIKey } = useUpdateAPIKey();
  const { mutateAsync: createAPIKey } = useCreateAPIKey();
  const { data: existingAPIKey } = useGetAPIKey(id, {
    query: { enabled: !!id },
  });
  const data = existingAPIKey?.data;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: data?.expiresAt ?? "",
      description: "When this API key expires",
      label: "Expires At",
      name: "expiresAt",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter expires at..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.name ?? "",
      description: "A friendly name for the API key",
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
      description: "The organization this API key belongs to",
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
    {
      defaultValue: data?.rateLimit ?? "",
      description: "Requests per minute allowed for this API key",
      label: "Rate Limit",
      name: "rateLimit",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter rate limit..."
          type="text"
        />
      ),
    },
  ];
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: data?.expiresAt ?? "",
      description: "When this API key expires",
      label: "Expires At",
      name: "expiresAt",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter expires at..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.name ?? "",
      description: "A friendly name for the API key",
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
      defaultValue: data?.rateLimit ?? "",
      description: "Requests per minute allowed for this API key",
      label: "Rate Limit",
      name: "rateLimit",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter rate limit..."
          type="text"
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateAPIKeyBody, UpdateAPIKeyBody>
      description={!id ? "Create a new apiKey" : "Update an existing apiKey"}
      entityKey="api_keys"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createAPIKey({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateAPIKey({ data: updateDto, id: id });
      }}
      title="APIKey"
    />
  );
}
