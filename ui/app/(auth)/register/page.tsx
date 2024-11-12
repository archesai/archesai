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
import { useAuth } from "@/hooks/use-auth";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

// Define schema using Zod for form validation
const schema = z
  .object({
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
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

type RegisterFormData = z.infer<typeof schema>;

export default function RegisterPage() {
  const router = useRouter();
  const { registerWithEmailAndPassword, user } = useAuth();
  const [error, setError] = useState<null | string>(null);

  useEffect(() => {
    if (user) {
      router.push("/playground");
    }
  }, [user, router]);

  const form = useForm<RegisterFormData>({
    defaultValues: {
      confirmPassword: "",
      email: "",
      password: "",
    },
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: RegisterFormData) => {
    try {
      await registerWithEmailAndPassword(data.email, data.password);
      // Optionally, you can set a success message or perform additional actions here
    } catch (err: any) {
      console.error("Registration error:", err);
      // Enhanced error handling to capture specific error messages
      if (err?.message) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred. Please try again.");
      }
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-col gap-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">Register</h1>
        <p className="text-sm text-muted-foreground">
          Create your account by entering your email and password
        </p>
      </div>
      <div className="flex flex-col gap-2">
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
                  <FormLabel htmlFor="email">Email</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete="email"
                      id="email"
                      placeholder="m@example.com"
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

            {/* Password Field */}
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor="password">Password</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete="new-password"
                      id="password"
                      placeholder="Enter your password"
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

            {/* Confirm Password Field */}
            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor="confirmPassword">
                    Confirm Password
                  </FormLabel>
                  <FormControl>
                    <Input
                      autoComplete="new-password"
                      id="confirmPassword"
                      placeholder="Confirm your password"
                      type="password"
                      {...field}
                      aria-invalid={
                        form.formState.errors.confirmPassword ? "true" : "false"
                      }
                    />
                  </FormControl>
                  <FormMessage>
                    {form.formState.errors.confirmPassword?.message}
                  </FormMessage>
                </FormItem>
              )}
            />

            {/* Display Error Message */}
            {error && (
              <div className="text-center text-red-600" role="alert">
                {error}
              </div>
            )}

            {/* Submit Button */}
            <Button
              className="mt-5 w-full"
              disabled={form.formState.isSubmitting}
              type="submit"
            >
              {form.formState.isSubmitting ? "Registering..." : "Register"}
            </Button>
          </form>
        </Form>

        {/* Redirect to Login */}
        <div className="text-center text-sm">
          Already have an account?{" "}
          <Link className="underline" href="/login">
            Login
          </Link>
        </div>
      </div>
    </div>
  );
}
