import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { UpdateUserBody } from "#lib/index";
import { useGetUser, useUpdateUser } from "#lib/index";

export default function UserForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateUser } = useUpdateUser();
  const { data: existingUser } = useGetUser(id, { query: { enabled: !!id } });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingUser?.data as Record<string, unknown> | undefined;
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.email as string) ?? "",
      description: "The user's e-mail",
      label: "Email",
      name: "email",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter email..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.image as string) ?? "",
      description: "The user's avatar image URL",
      label: "Image",
      name: "image",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter image..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
  ];
  return (
    <GenericForm<UpdateUserBody, UpdateUserBody>
      description="Update an existing user"
      entityKey="users"
      fields={updateFormFields}
      isUpdateForm={true}
      onSubmitUpdate={async (updateDto) => {
        await updateUser({ data: updateDto, id: id });
      }}
      title="User"
    />
  );
}
