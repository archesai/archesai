// src/pages/auth/forgot-password.tsx

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
import { useAuthControllerPasswordResetRequest } from "@/generated/archesApiComponents"; // Adjust the import path as needed
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

export const forgotPasswordSchema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
});

export type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

export default function ForgotPasswordPage() {
  const { mutateAsync: requestPasswordReset } =
    useAuthControllerPasswordResetRequest();
  const [message, setMessage] = useState<string>("");
  const [error, setError] = useState<string>("");

  const form = useForm<ForgotPasswordFormData>({
    defaultValues: {
      email: "",
    },
    resolver: zodResolver(forgotPasswordSchema),
  });

  const onSubmit = async (data: ForgotPasswordFormData) => {
    setMessage("");
    setError("");
    try {
      await requestPasswordReset({
        body: {
          email: data.email,
        },
      });
      setMessage(
        "If an account with that email exists, a password reset link has been sent."
      );
    } catch (err: any) {
      console.error("Password reset request error:", err);
      setError(
        err?.response?.data?.message ||
          "An unexpected error occurred. Please try again."
      );
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-col gap-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">
          Forgot Password
        </h1>
        <p className="text-md text-muted-foreground">
          Enter your email address to receive a password reset link.
        </p>
      </div>
      <div className="flex flex-col gap-2">
        {/* Display Success Message */}
        {message && (
          <div className="text-green-600" role="alert">
            {message}
          </div>
        )}

        {/* Display Error Message */}
        {error && (
          <div className="text-red-600" role="alert">
            {error}
          </div>
        )}

        <Form {...form}>
          <form
            className="flex flex-col gap-2"
            noValidate
            onSubmit={form.handleSubmit(onSubmit)}
          >
            {/* Email Field */}
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor="email">Email Address</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete="email"
                      id="email"
                      placeholder="you@example.com"
                      type="email"
                      {...field}
                      aria-invalid={
                        form.formState.errors.email ? "true" : "false"
                      }
                    />
                  </FormControl>
                  <FormMessage>
                    {form.formState.errors.email?.message}
                  </FormMessage>
                </FormItem>
              )}
            />

            {/* Submit Button */}
            <Button
              className="mt-5 w-full"
              disabled={form.formState.isSubmitting}
              type="submit"
            >
              {form.formState.isSubmitting ? "Sending..." : "Send Reset Link"}
            </Button>
          </form>
        </Form>

        {/* Redirect to Login */}
        <div className="text-center text-sm">
          Remembered your password?{" "}
          <Link className="underline" href="/auth/login">
            Login
          </Link>
        </div>
      </div>
    </div>
  );
}
