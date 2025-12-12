import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Input, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type {
  CreateOrganizationBody,
  UpdateOrganizationBody,
} from "#lib/index";
import {
  useCreateOrganization,
  useGetOrganization,
  useUpdateOrganization,
} from "#lib/index";

export default function OrganizationForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateOrganization } = useUpdateOrganization();
  const { mutateAsync: createOrganization } = useCreateOrganization();
  const { data: existingOrganization } = useGetOrganization(id, {
    query: { enabled: !!id },
  });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingOrganization?.data as
    | Record<string, unknown>
    | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.billingEmail as string) ?? "",
      description: "The billing email to use for the organization",
      label: "Billing Email",
      name: "billingEmail",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter billing email..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.organizationID as string) ?? "",
      description: "UUID identifier",
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
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.billingEmail as string) ?? "",
      description: "The billing email to use for the organization",
      label: "Billing Email",
      name: "billingEmail",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter billing email..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.organizationID as string) ?? "",
      description: "UUID identifier",
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
    <GenericForm<CreateOrganizationBody, UpdateOrganizationBody>
      description={
        !id ? "Create a new organization" : "Update an existing organization"
      }
      entityKey="organizations"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createOrganization({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateOrganization({ data: updateDto, id: id });
      }}
      title="Organization"
    />
  );
}
