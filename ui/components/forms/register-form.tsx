"use client";
import {
  FormFieldConfig,
  GenericForm,
} from "@/components/forms/generic-form/generic-form";
import { Input } from "@/components/ui/input";
import { RegisterDto } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import * as z from "zod";

// Define schema using Zod for form validation
const formSchema = z.object({
  confirmPassword: z
    .string()
    .min(8, { message: "Please confirm your password" }),
  email: z.string().email({ message: "Invalid email address" }),
  password: z
    .string()
    .min(8, { message: "Password must be at least 8 characters" })
    .regex(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/, {
      message:
        "Password must contain at least one letter, one number, and one special character",
    }),
});
//   .refine((data) => data.password === data.confirmPassword, {
//     message: "Passwords do not match",
//     path: ["confirmPassword"],
//   });

export default function RegisterForm() {
  const { registerWithEmailAndPassword } = useAuth();

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: "",
      description: "This is the email that will be used for this user.",
      label: "Email",
      name: "email",
      props: {
        placeholder: "Enter your email here...",
      },
      validationRule: formSchema.shape.email,
    },
    {
      component: Input,
      defaultValue: "",
      description: "This is the password that will be used for your account.",
      label: "Password",
      name: "password",
      props: {
        placeholder: "Enter your password here...",
      },
      validationRule: formSchema.shape.password,
    },
    {
      component: Input,
      defaultValue: "",
      description:
        "This is the role that will be used for this member. Note that different roles have different permissions.",
      label: "Confirm Password",
      name: "confirmPassword",
      props: {
        placeholder: "Confirm your password here...",
      },
      validationRule: formSchema.shape.confirmPassword,
    },
  ];

  return (
    <GenericForm<RegisterDto, undefined>
      description={"Configure your member's settings"}
      fields={formFields}
      isUpdateForm={false}
      itemType="member"
      onSubmitCreate={async (registerDto) => {
        await registerWithEmailAndPassword(
          registerDto.email,
          registerDto.password
        );
      }}
      onSubmitUpdate={async () => {}}
      title="Configuration"
    />
  );
}
