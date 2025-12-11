import type { FormFieldConfig } from "@archesai/ui";
import { Checkbox, GenericForm, Textarea } from "@archesai/ui";
import type { JSX } from "react";
import type { UpdateUserBody } from "#lib/index";
import { useGetUser, useUpdateUser } from "#lib/index";

export default function UserForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateUser } = useUpdateUser();
  const { data: existingUser } = useGetUser(id, { query: { enabled: !!id } });
  const data = existingUser?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.email ?? "",
      description: "The user's email address",
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
      defaultValue: data?.emailVerified ?? false,
      description: "Whether the user's email has been verified",
      label: "Email Verified",
      name: "emailVerified",
      renderControl: (field) => (
        <Checkbox
          checked={field.value as boolean}
          onCheckedChange={field.onChange}
        />
      ),
    },
    {
      defaultValue: data?.image ?? "",
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
    {
      defaultValue: data?.name ?? "",
      description: "The user's display name",
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
  ];
  return (
    <GenericForm<UpdateUserBody, UpdateUserBody>
      description="Update an existing user"
      entityKey="users"
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (updateDto) => {
        await updateUser({ data: updateDto, id: id });
      }}
      title="User"
    />
  );
}
