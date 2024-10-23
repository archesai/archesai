"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/hooks/useAuth";

export default function OrganizationSettingsPage() {
  const { defaultOrgname } = useAuth();
  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: defaultOrgname,
      description: "The name of the organization. This cannot be changed.",
      label: "Name",
      name: "name",
      props: {
        disabled: true,
      },
    },
  ];

  return (
    <GenericForm
      description={"View your organization's details"}
      fields={formFields}
      isUpdateForm={true}
      itemType="organization"
      showCard={true}
      title="Organization"
    />
  );
}
