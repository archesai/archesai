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
import { useAuth } from "@/hooks/useAuth";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

// Define schema using Zod for form validation
const schema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
  password: z.string().min(1, { message: "Password is required" }),
});

type LoginFormData = z.infer<typeof schema>;

export default function LoginPage() {
  const router = useRouter();
  const { signInWithEmailAndPassword, signInWithGoogle, user } = useAuth();
  const [formError, setFormError] = useState<null | string>(null);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  useEffect(() => {
    if (user) {
      router.push("/home");
    }
  }, [user, router]);

  const form = useForm<LoginFormData>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: LoginFormData) => {
    setIsSubmitting(true);
    setFormError(null);
    try {
      await signInWithEmailAndPassword(data.email, data.password);
      // Redirect handled by useEffect
    } catch (error: any) {
      console.error("Login error", error);
      // Enhanced error handling to capture specific error messages
      if (error?.message) {
        setFormError(error.message);
      } else {
        setFormError("An unexpected error occurred. Please try again.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleGoogleSignIn = async () => {
    setIsSubmitting(true);
    try {
      await signInWithGoogle();
      // Redirect handled by useEffect
    } catch (error: any) {
      console.error("Google Sign-In error", error);
      if (error?.message) {
        setFormError(error.message);
      } else {
        setFormError("An unexpected error occurred. Please try again.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="stack gap-2">
      <div className="flex flex-col space-y-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">Login</h1>
        <p className="text-sm text-muted-foreground">
          Enter your email and password to login to your account
        </p>
      </div>
      <div>
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
                      autoComplete="current-password"
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
                  {/* Forgot Password Link */}
                  <div className="text-right">
                    <Link
                      className="inline-block text-sm underline"
                      href="/auth/forgot"
                    >
                      Forgot your password?
                    </Link>
                  </div>
                </FormItem>
              )}
            />

            {/* Display Form Error */}
            {formError && (
              <div className="text-red-600 text-center" role="alert">
                {formError}
              </div>
            )}

            {/* Submit Button */}
            <Button className="w-full" disabled={isSubmitting} type="submit">
              {isSubmitting ? "Logging in..." : "Login"}
            </Button>
          </form>
        </Form>

        {/* Conditional Firebase Login Button */}
        {process.env.NEXT_PUBLIC_USE_FIREBASE === "true" && (
          <Button
            className="w-full mt-2"
            disabled={isSubmitting}
            onClick={handleGoogleSignIn}
            variant="outline"
          >
            {isSubmitting ? "Processing..." : "Login with Google"}
          </Button>
        )}

        {/* Redirect to Register */}
        <div className="mt-4 text-center text-sm">
          Don&apos;t have an account?{" "}
          <Link className="underline" href="/auth/register">
            Sign up
          </Link>
        </div>
      </div>
    </div>
  );
}
