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
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingAccount?.data as Record<string, unknown> | undefined;
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.provider as string) ?? "",
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
      defaultValue: (data?.providerAccountIdentifier as string) ?? "",
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
      defaultValue: (data?.type as string) ?? "",
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
