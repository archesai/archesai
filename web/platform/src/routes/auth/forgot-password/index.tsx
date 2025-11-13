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
  const { mutateAsync: requestPasswordReset, isSuccess } =
    useRequestPasswordReset();

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

  if (isSuccess) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="w-full max-w-md space-y-4 rounded-lg border bg-card p-6">
          <div className="text-center">
            <h2 className="font-bold text-2xl">Check Your Email</h2>
            <p className="mt-2 text-muted-foreground">
              We've sent a password reset link to your email address. Please
              check your inbox and follow the instructions.
            </p>
          </div>
          <div className="text-center text-sm">
            <Link
              className="underline"
              to="/auth/login"
            >
              Return to Login
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
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
  );
}
