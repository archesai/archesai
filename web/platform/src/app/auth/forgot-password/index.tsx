import type { RequestPasswordResetBody } from "@archesai/client";
import { useRequestPasswordReset } from "@archesai/client";
import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Input } from "@archesai/ui";
import { createFileRoute, Link } from "@tanstack/react-router";
import type { JSX } from "react";

export const Route = createFileRoute("/auth/forgot-password/")({
  component: ForgotPasswordPage,
});

function ForgotPasswordPage(): JSX.Element {
  const { mutateAsync: requestPasswordReset } = useRequestPasswordReset();

  const formFields: FormFieldConfig<RequestPasswordResetBody>[] = [
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
  ];

  return (
    <>
      <GenericForm<RequestPasswordResetBody, never>
        description="Enter your email address to receive a password reset link"
        entityKey="auth"
        fields={formFields}
        isUpdateForm={false}
        onSubmitCreate={async (data) => {
          await requestPasswordReset({
            data: {
              email: data.email,
            },
          });
        }}
        postContent={
          <div className="text-center text-sm">
            Remembered your password?{" "}
            <Link
              className="underline"
              to="/auth/login"
            >
              Login
            </Link>
          </div>
        }
        showCard={true}
        title="Forgot Password"
      />
    </>
  );
}
