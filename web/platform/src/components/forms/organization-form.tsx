import type {
  CreateOrganizationBody,
  UpdateOrganizationMutationBody,
} from "@archesai/client";
import {
  useCreateOrganization,
  useGetOneOrganizationSuspense,
  useGetOneSessionSuspense,
  useUpdateOrganization,
} from "@archesai/client";
import type { FormFieldConfig } from "@archesai/ui/components/custom/generic-form";
import { GenericForm } from "@archesai/ui/components/custom/generic-form";
import { Input } from "@archesai/ui/components/shadcn/input";
import { ORGANIZATION_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { JSX } from "react";

export default function OrganizationForm(): JSX.Element {
  const {
    data: { data: session },
  } = useGetOneSessionSuspense("current");
  const { mutateAsync: createOrganization } = useCreateOrganization();
  const { mutateAsync: updateOrganization } = useUpdateOrganization();
  const {
    data: { data: organization },
  } = useGetOneOrganizationSuspense(session.activeOrganizationId);

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: organization.name,
      description: "Organization name",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type="text"
        />
      ),
    },
    {
      defaultValue: organization.billingEmail ?? "",
      description: "Email address for billing notifications",
      label: "Billing Email",
      name: "billingEmail",
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type="email"
        />
      ),
    },
  ];

  return (
    <GenericForm<CreateOrganizationBody, UpdateOrganizationMutationBody>
      description={"View your organization's details"}
      entityKey={ORGANIZATION_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitCreate={async (createOrganizationDto: CreateOrganizationBody) => {
        await createOrganization({
          data: createOrganizationDto,
        });
      }}
      onSubmitUpdate={async (
        updateOrganizationDto: UpdateOrganizationMutationBody,
      ) => {
        await updateOrganization({
          data: updateOrganizationDto,
          id: session.activeOrganizationId,
        });
      }}
      showCard={true}
      title={"Organiation"}
    />
  );
}
