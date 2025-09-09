import type { UpdateUserBody } from "@archesai/client";
import {
  useGetOneSessionSuspense,
  useGetOneUserSuspense,
  useUpdateUser,
} from "@archesai/client";
import type { FormFieldConfig } from "@archesai/ui/components/custom/generic-form";
import { GenericForm } from "@archesai/ui/components/custom/generic-form";
import { Input } from "@archesai/ui/components/shadcn/input";
import { USER_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { JSX } from "react";

export default function UserForm(): JSX.Element {
  const { mutateAsync: updateUser } = useUpdateUser();
  const { data: sessionData } = useGetOneSessionSuspense("current");
  const { data: userData } = useGetOneUserSuspense(
    sessionData.data.userId || "current",
  );

  const formFields: FormFieldConfig<UpdateUserBody>[] = [
    {
      defaultValue: userData.data.email || "",
      description: "Your email address",
      label: "Email",
      name: "email",
      renderControl: (field) => (
        <Input
          {...field}
          type="email"
        />
      ),
    },
    {
      defaultValue: userData.data.image ?? "",
      description: "Your profile image URL",
      label: "Image URL",
      name: "image",
      renderControl: (field) => (
        <Input
          {...field}
          type="url"
        />
      ),
    },
  ];

  return (
    <GenericForm<never, UpdateUserBody>
      description="View and update your user details"
      entityKey={USER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (data) => {
        await updateUser({
          data,
          id: userData.data.id || "",
        });
      }}
      showCard={true}
      title="Profile"
    />
  );
}
