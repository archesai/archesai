import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { UpdateAccountBody } from "#lib/index";
import { useGetAccount, useUpdateAccount } from "#lib/index";

export default function AccountForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateAccount } = useUpdateAccount();
  const { data: existingAccount } = useGetAccount(id, {
    query: { enabled: !!id },
  });
  const data = existingAccount?.data;
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: data?.provider ?? "",
      description: "The account provider",
      label: "Provider",
      name: "provider",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter provider..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.providerAccountIdentifier ?? "",
      description: "The provider account ID",
      label: "Provider Account Identifier",
      name: "providerAccountIdentifier",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter provider account identifier..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.type ?? "",
      description: "The account type",
      label: "Type",
      name: "type",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter type..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  return (
    <GenericForm<UpdateAccountBody, UpdateAccountBody>
      description="Update an existing account"
      entityKey="accounts"
      fields={updateFormFields}
      isUpdateForm={true}
      onSubmitUpdate={async (updateDto) => {
        await updateAccount({ data: updateDto, id: id });
      }}
      title="Account"
    />
  );
}
