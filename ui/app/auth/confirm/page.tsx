"use client";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  useEmailChangeControllerConfirm,
  useEmailVerificationControllerConfirm,
  usePasswordResetControllerConfirm,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

// Define allowed action types
type ActionType = "email-change" | "email-verification" | "password-reset";

// Define schemas for different actions
const passwordResetSchema = z
  .object({
    confirmPassword: z
      .string()
      .min(8, { message: "Please confirm your password" }),
    password: z
      .string()
      .min(8, { message: "Password must be at least 8 characters" })
      .regex(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/, {
        message:
          "Password must contain at least one letter, one number, and one special character",
      }),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

type PasswordResetFormData = z.infer<typeof passwordResetSchema>;

export default function ConfirmPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const type = searchParams?.get("type") as ActionType;
  const token = searchParams?.get("token") as string;

  const { setAccessToken, setRefreshToken, user } = useAuth();

  const { mutateAsync: verifyEmail } = useEmailVerificationControllerConfirm();
  const { mutateAsync: resetPassword } = usePasswordResetControllerConfirm();
  const { mutateAsync: changeEmail } = useEmailChangeControllerConfirm();

  const [message, setMessage] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [operationSent, setOperationSent] = useState<boolean>(false);

  const form = useForm<PasswordResetFormData>({
    defaultValues: {
      confirmPassword: "",
      password: "",
    },
    resolver: zodResolver(passwordResetSchema),
  });

  useEffect(() => {
    const handleAction = async () => {
      if (!type || !token) {
        setError("Invalid request. Missing parameters.");
        router.push("/auth/login");
        return;
      }
      if (operationSent) {
        return;
      }
      setOperationSent(true);

      switch (type) {
        case "email-change":
          try {
            const { accessToken, refreshToken } = await changeEmail({
              body: {
                token,
              },
            });
            setMessage("Your email has been successfully updated!");
            setAccessToken(accessToken);
            setRefreshToken(refreshToken);
            router.push("/playground");
          } catch (err: any) {
            console.error(err);
            setError(
              err?.response?.data?.message ||
                "Email change failed. Please try again."
            );
          }
          break;
        case "email-verification":
          try {
            const { accessToken, refreshToken } = await verifyEmail({
              body: {
                token,
              },
            });
            setMessage("Your email has been successfully verified!");
            setAccessToken(accessToken);
            setRefreshToken(refreshToken);
            router.push("/playground");
          } catch (err: any) {
            console.error(err);
            // setError(
            //   err?.response?.data?.message || "Email verification failed."
            // );
          }
          break;
        case "password-reset":
          // Do nothing
          break;

        default:
          setError("Unsupported action type.");
          router.push("/auth/login");
          break;
      }
    };

    handleAction();
  }, [type, token]);

  const onSubmit = async (data: PasswordResetFormData) => {
    try {
      const { accessToken, refreshToken } = await resetPassword({
        body: {
          newPassword: data.password,
          token,
        },
      });
      setMessage("Your password has been successfully reset!");
      setAccessToken(accessToken);
      setRefreshToken(refreshToken);
      router.push("/playground");
    } catch (err: any) {
      console.error("Password reset error:", err);
      setError(
        err?.response?.data?.message ||
          "Password reset failed. Please try again."
      );
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-col gap-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">
          {type.split("-").join(" ")}
        </h1>
        <p className="text-sm text-muted-foreground">
          {message ||
            (error
              ? ""
              : type === "password-reset"
                ? "Please follow the instructions below."
                : "Verifying...")}
        </p>
      </div>
      <div className="flex flex-col gap-2">
        {/* Display Error Message */}
        {error && (
          <p className="text-red-500" role="alert">
            {error}
          </p>
        )}

        {/* Handle Password Reset Form */}
        {type === "password-reset" ? (
          <Form {...form}>
            <form
              className="flex flex-col gap-2"
              noValidate
              onSubmit={form.handleSubmit(onSubmit)}
            >
              {/* New Password Field */}
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel htmlFor="password">New Password</FormLabel>
                    <FormControl>
                      <Input
                        autoComplete="new-password"
                        id="password"
                        placeholder="Enter your new password"
                        type="password"
                        {...field}
                        aria-invalid={
                          form.formState.errors.password ? "true" : "false"
                        }
                      />
                    </FormControl>
                    <FormMessage>
                      {form.formState.errors.password?.message}
                    </FormMessage>
                  </FormItem>
                )}
              />

              {/* Confirm New Password Field */}
              <FormField
                control={form.control}
                name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel htmlFor="confirmPassword">
                      Confirm New Password
                    </FormLabel>
                    <FormControl>
                      <Input
                        autoComplete="new-password"
                        id="confirmPassword"
                        placeholder="Confirm your new password"
                        type="password"
                        {...field}
                        aria-invalid={
                          form.formState.errors.confirmPassword
                            ? "true"
                            : "false"
                        }
                      />
                    </FormControl>
                    <FormMessage>
                      {form.formState.errors.confirmPassword?.message}
                    </FormMessage>
                  </FormItem>
                )}
              />

              {/* Submit Button */}
              <Button
                className="w-full"
                disabled={form.formState.isSubmitting}
                type="submit"
              >
                {form.formState.isSubmitting
                  ? "Resetting Password..."
                  : "Reset Password"}
              </Button>
            </form>
          </Form>
        ) : (
          // Handle Email Verification & Email Change
          <div className="text-center">
            {/* Navigate to Home if Email is Verified */}
            {user?.emailVerified && (
              <Button asChild>
                <Link href="/playground">Go to Home</Link>
              </Button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
