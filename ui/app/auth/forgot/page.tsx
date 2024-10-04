// src/pages/auth/forgot-password.tsx

"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { usePasswordResetControllerRequest } from "@/generated/archesApiComponents"; // Adjust the import path as needed
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
    usePasswordResetControllerRequest();
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
    <Card className="mx-auto max-w-sm mt-10">
      <CardHeader>
        <CardTitle className="text-2xl">Forgot Password</CardTitle>
        <CardDescription>
          Enter your email address to receive a password reset link.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {/* Display Success Message */}
        {message && (
          <div className="text-green-600 mb-4" role="alert">
            {message}
          </div>
        )}

        {/* Display Error Message */}
        {error && (
          <div className="text-red-600 mb-4" role="alert">
            {error}
          </div>
        )}

        <Form {...form}>
          <form
            className="space-y-4"
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
              className="w-full"
              disabled={form.formState.isSubmitting}
              type="submit"
            >
              {form.formState.isSubmitting ? "Sending..." : "Send Reset Link"}
            </Button>
          </form>
        </Form>

        {/* Redirect to Login */}
        <div className="mt-4 text-center text-sm">
          Remembered your password?{" "}
          <Link className="underline" href="/auth/login">
            Login
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}
