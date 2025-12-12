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

  const formFields: FormFieldConfig[] = [
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
      defaultValue: data?.keyHash ?? "",
      description: "Hashed version of the API key for secure storage",
      label: "Key Hash",
      name: "keyHash",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter key hash..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.lastUsedAt ?? "",
      description: "When this API key was last used",
      label: "Last Used At",
      name: "lastUsedAt",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter last used at..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.name ?? "",
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
      defaultValue: data?.prefix ?? "",
      label: "Prefix",
      name: "prefix",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter prefix..."
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
    {
      defaultValue: data?.userID ?? "",
      description: "The user who owns this API key",
      label: "User ID",
      name: "userID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter user id..."
          type="text"
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateAPIKeyBody, UpdateAPIKeyBody>
      description={!id ? "Create a new apiKey" : "Update an existing apiKey"}
      entityKey="api_keys"
      fields={formFields}
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
