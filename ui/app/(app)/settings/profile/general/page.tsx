"use client";

import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import {
  useUserControllerFindOne,
  useUserControllerUpdate,
} from "@/generated/archesApiComponents";
import { UpdateUserDto } from "@/generated/archesApiSchemas";
import React from "react";
import { z } from "zod";

export default function ProfileSettingsPage() {
  const { data: user } = useUserControllerFindOne({});
  const { mutateAsync: updateUser } = useUserControllerUpdate();

  if (!user) {
    return <div>Loading...</div>;
  }

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: user.firstName,
      description: "Your first name",
      label: "First Name",
      name: "firstName",
      validationRule: z.string().min(1, "First name is required"),
    },
    {
      component: Input,
      defaultValue: user.lastName,
      description: "Your last name",
      label: "Last Name",
      name: "lastName",
      validationRule: z.string().min(1, "Last name is required"),
    },
    {
      component: Input,
      defaultValue: user.username,
      description: "Your username",
      label: "Username",
      name: "username",
      props: {
        disabled: true,
      },
    },
    {
      component: Input,
      defaultValue: user.email,
      description: "Your email address",
      label: "Email",
      name: "email",
      props: {
        disabled: true,
      },
    },
  ];

  return (
    <GenericForm<any, UpdateUserDto>
      description="View and update your user details"
      fields={formFields}
      isUpdateForm={true}
      itemType="user"
      onSubmitCreate={() => {}}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateUser(
          {
            body: {
              defaultOrg: data.defaultOrg,
              firstName: data.firstName,
              lastName: data.lastName,
            },
          },
          mutateOptions
        );
      }}
      title="Profile"
    />
  );
}
