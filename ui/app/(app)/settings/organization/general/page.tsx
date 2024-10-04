"use client";
import { CustomCardForm, FormFieldConfig } from "@/components/custom-card-form";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/hooks/useAuth";

export default function OrganizationSettingsPage() {
  const { defaultOrgname } = useAuth();
  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      description: "The name of the organization. This cannot be changed.",
      label: "Name",
      name: "name",
      props: {
        disabled: true,
        value: defaultOrgname,
      },
    },
  ];

  return (
    <CustomCardForm
      description={"View your organization's details"}
      fields={formFields}
      isUpdateForm={true}
      title="Organization"
    />
  );
}
