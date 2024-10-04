"use client";

import { CustomCardForm, FormFieldConfig } from "@/components/custom-card-form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import {
  useUserControllerFindOne,
  useUserControllerUpdate,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import React from "react";
import { z } from "zod";

export default function ProfileSettingsPage() {
  const { defaultOrgname } = useAuth();
  const { data: user } = useUserControllerFindOne(
    {},
    { enabled: !!defaultOrgname }
  );
  const { mutateAsync: updateUser } = useUserControllerUpdate();
  const { toast } = useToast();

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
    <CustomCardForm
      description="View and update your user details"
      fields={formFields}
      onSubmit={async (data) => {
        await updateUser(
          {
            body: {
              defaultOrg: data.defaultOrg,
              firstName: data.firstName,
              lastName: data.lastName,
            },
          },
          {
            onError: (error) => {
              toast({
                description: error?.stack?.msg,
                title: "Error updating profile",
              });
            },
            onSuccess: () => {
              toast({
                title: "Profile updated",
              });
            },
          }
        );
      }}
      title="Profile"
    />
  );
}
