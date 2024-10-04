"use client";

import { CustomCardForm, FormFieldConfig } from "@/components/custom-card-form";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { useUserControllerFindOne } from "@/generated/archesApiComponents";
import React from "react";
import { z } from "zod";

export default function ProfileSettingsPage() {
  const { data: user } = useUserControllerFindOne({}, { enabled: true });

  if (!user) {
    return <div>Loading...</div>;
  }

  const formFields: FormFieldConfig[] = [
    // {
    //   component: Input,
    //   defaultValue: user.firstName,
    //   description: "Your first name",
    //   label: "First Name",
    //   name: "firstName",
    //   validationRule: z.string().min(1, "First name is required"),
    // },
    // {
    //   component: Input,
    //   defaultValue: user.lastName,
    //   description: "Your last name",
    //   label: "Last Name",
    //   name: "lastName",
    //   validationRule: z.string().min(1, "Last name is required"),
    // },
    {
      component: Input,
      defaultValue: user.displayName,
      description: "Your display name",
      label: "Display Name",
      name: "displayName",
      validationRule: z.string().min(1, "Display name is required"),
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
    {
      component: Checkbox,
      description: "Email Verified",
      label: "Email Verified",
      name: "emailVerified",
      props: {
        checked: user.emailVerified,
        disabled: true,
      },
    },
    {
      component: Input,
      defaultValue: user.defaultOrg,
      description: "Your default organization",
      label: "Default Organization",
      name: "defaultOrg",
      props: {
        disabled: true,
      },
    },
  ];

  return (
    <CustomCardForm
      description="View and update your user details"
      fields={formFields}
      title="Profile"
    />
  );
}
