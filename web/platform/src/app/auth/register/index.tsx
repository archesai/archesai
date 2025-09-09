import type { RegisterBody } from "@archesai/client";
import { useRegister } from "@archesai/client";
import type { FormFieldConfig } from "@archesai/ui/components/custom/generic-form";
import { GenericForm } from "@archesai/ui/components/custom/generic-form";
import { Input } from "@archesai/ui/components/shadcn/input";
import { createFileRoute, Link, useRouter } from "@tanstack/react-router";
import type { JSX } from "react";

import { TermsIndicator } from "#components/terms-indicator";

type RegisterFormData = RegisterBody & {
  confirmPassword: string;
};

export const Route = createFileRoute("/auth/register/")({
  component: RegisterPage,
});

export default function RegisterPage(): JSX.Element {
  const router = useRouter();
  const { mutateAsync: register } = useRegister();

  const formFields: FormFieldConfig<RegisterFormData>[] = [
    {
      defaultValue: "",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Input
          {...field}
          type="text"
        />
      ),
    },
    {
      defaultValue: "",
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
      defaultValue: "",
      label: "Password",
      name: "password",
      renderControl: (field) => (
        <Input
          {...field}
          type="password"
        />
      ),
    },
    {
      defaultValue: "",
      label: "Confirm Password",
      name: "confirmPassword",
      renderControl: (field) => (
        <Input
          {...field}
          type="password"
        />
      ),
    },
  ];

  return (
    <>
      <GenericForm<RegisterFormData, never>
        description="Create your account by entering your email and password"
        entityKey="auth"
        fields={formFields}
        isUpdateForm={false}
        onSubmitCreate={async (data) => {
          await register({
            data: {
              email: data.email,
              name: data.name,
              password: data.password,
            },
          });
          await router.navigate({
            to: "/",
          });
        }}
        postContent={
          <div className="text-center text-sm">
            Already have an account?{" "}
            <Link
              className="underline"
              to="/auth/login"
            >
              Login
            </Link>
          </div>
        }
        showCard={true}
        title="Register"
      />
      <TermsIndicator />
    </>
  );
}
