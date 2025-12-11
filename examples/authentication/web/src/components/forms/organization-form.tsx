import type { FormFieldConfig } from "@archesai/ui";
import {
  FormControl,
  GenericForm,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Textarea,
} from "@archesai/ui";
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
  const data = existingOrganization?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.billingEmail ?? "",
      description: "Email address for billing communications",
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
      defaultValue: data?.credits ?? "",
      description: "Available credits for this organization",
      label: "Credits",
      name: "credits",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter credits..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.logo ?? "",
      description: "The organization's logo URL",
      label: "Logo",
      name: "logo",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter logo..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.name ?? "",
      description: "The organization's display name",
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
      defaultValue: data?.plan ?? "",
      description: "The current subscription plan",
      label: "Plan",
      name: "plan",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select plan..." />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="FREE">FREE</SelectItem>
            <SelectItem value="BASIC">BASIC</SelectItem>
            <SelectItem value="STANDARD">STANDARD</SelectItem>
            <SelectItem value="PREMIUM">PREMIUM</SelectItem>
            <SelectItem value="UNLIMITED">UNLIMITED</SelectItem>
          </SelectContent>
        </Select>
      ),
    },
    {
      defaultValue: data?.slug ?? "",
      description: "URL-friendly unique identifier for the organization",
      label: "Slug",
      name: "slug",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter slug..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.stripeCustomerIdentifier ?? "",
      description: "Stripe customer identifier",
      label: "Stripe Customer Identifier",
      name: "stripeCustomerIdentifier",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter stripe customer identifier..."
          rows={5}
          value={field.value as string}
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
      fields={formFields}
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
