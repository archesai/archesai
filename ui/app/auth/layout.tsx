"use client";
import { LogoSVG } from "@/components/logo-svg";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import Link from "next/link";

export default function AuthenticationPage({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <>
      <div className="relative grid h-screen flex-col items-center justify-center lg:max-w-none lg:grid-cols-2">
        <Link
          className={cn(
            buttonVariants({ variant: "ghost" }),
            "absolute right-4 top-4 md:right-8 md:top-8"
          )}
          href="/"
        >
          Back
        </Link>
        {/* Left side of the screen */}
        <div className="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
          <div className="absolute inset-0 bg-zinc-900" /> {/* FIXME */}
          <div className="relative z-20 flex items-center text-lg font-medium">
            <LogoSVG />
          </div>
          <div className="relative z-20 mt-auto">
            <blockquote className="flex flex-col gap-2">
              <p className="text-lg">
                &ldquo;This library has saved me countless hours of work and
                helped me deliver stunning designs to my clients faster than
                ever before.&rdquo;
              </p>
              <footer className="text-sm">Sofia Davis</footer>
            </blockquote>
          </div>
        </div>

        {/* Right side of the screen or main*/}
        <div className="mx-auto flex w-[350px] flex-col items-center justify-center gap-3">
          {children}
          <p className="text-center text-sm text-muted-foreground">
            By clicking continue, you agree to our{" "}
            <Link
              className="underline underline-offset-4 hover:text-primary"
              href="/terms"
            >
              Terms of Service
            </Link>{" "}
            and{" "}
            <Link
              className="underline underline-offset-4 hover:text-primary"
              href="/privacy"
            >
              Privacy Policy
            </Link>
            .
          </p>
        </div>
      </div>
    </>
  );
}
