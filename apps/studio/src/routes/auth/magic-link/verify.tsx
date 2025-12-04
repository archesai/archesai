import { Button, Input, toast } from "@archesai/ui";
import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import type { JSX } from "react";
import { useCallback, useEffect, useState } from "react";

export const Route = createFileRoute("/auth/magic-link/verify")({
  component: MagicLinkVerifyPage,
});

function MagicLinkVerifyPage(): JSX.Element {
  const router = useRouter();
  const queryClient = useQueryClient();
  const searchParams = Route.useSearch() as { token?: string };
  const [verifying, setVerifying] = useState(false);
  const [otpCode, setOtpCode] = useState("");
  const [email, setEmail] = useState("");
  const [useOtp, setUseOtp] = useState(false);

  const verifyMagicLink = useCallback(
    async (token: string) => {
      setVerifying(true);
      try {
        const response = await fetch("/api/v1/auth/magic-links/verify", {
          body: JSON.stringify({ token }),
          headers: {
            "Content-Type": "application/json",
          },
          method: "POST",
        });

        if (!response.ok) {
          throw new Error("Invalid or expired magic link");
        }

        // Backend sets httpOnly cookie for session
        await response.json();

        toast.success("Successfully logged in!");

        // Invalidate queries and navigate
        await router.invalidate();
        await queryClient.clear();
        await router.navigate({ to: "/" });
      } catch (error) {
        console.error("Verification failed:", error);
        toast.error("Invalid or expired magic link");
        setVerifying(false);
      }
    },
    [router, queryClient],
  );

  // Auto-verify if token is in URL
  useEffect(() => {
    if (searchParams.token) {
      verifyMagicLink(searchParams.token);
    }
  }, [searchParams.token, verifyMagicLink]);

  const verifyOtp = async () => {
    if (!email || !otpCode) {
      toast.error("Please enter both email and OTP code");
      return;
    }

    setVerifying(true);
    try {
      const response = await fetch("/api/v1/auth/magic-links/verify", {
        body: JSON.stringify({
          code: otpCode,
          identifier: email,
        }),
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
      });

      if (!response.ok) {
        throw new Error("Invalid OTP code");
      }

      // Backend sets httpOnly cookie for session
      await response.json();

      toast.success("Successfully logged in!");

      // Invalidate queries and navigate
      await router.invalidate();
      await queryClient.clear();
      await router.navigate({ to: "/" });
    } catch (error) {
      console.error("OTP verification failed:", error);
      toast.error("Invalid OTP code. Please try again.");
      setVerifying(false);
    }
  };

  // If we have a token, show verifying state
  if (searchParams.token && verifying) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="w-full max-w-md space-y-4 p-8">
          <div className="text-center">
            <h1 className="font-semibold text-2xl">Verifying Magic Link</h1>
            <p className="mt-2 text-muted-foreground">
              Please wait while we verify your magic link...
            </p>
            <div className="mt-4">
              <div className="mx-auto h-8 w-8 animate-spin rounded-full border-primary border-b-2"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Manual OTP entry form
  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="w-full max-w-md space-y-4 p-8">
        <div className="rounded-lg border bg-card p-6 shadow-sm">
          <div className="space-y-4">
            <div className="text-center">
              <h1 className="font-semibold text-2xl">Verify Magic Link</h1>
              <p className="mt-2 text-muted-foreground text-sm">
                {useOtp
                  ? "Enter the OTP code from your console or email"
                  : "Click the link in your email or enter the OTP code"}
              </p>
            </div>

            {!searchParams.token && (
              <div className="space-y-4">
                {useOtp ? (
                  <>
                    <Input
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="Enter your email..."
                      type="email"
                      value={email}
                    />
                    <Input
                      maxLength={6}
                      onChange={(e) => setOtpCode(e.target.value)}
                      placeholder="Enter 6-digit OTP code..."
                      type="text"
                      value={otpCode}
                    />
                    <Button
                      className="w-full"
                      disabled={verifying || !email || !otpCode}
                      onClick={verifyOtp}
                    >
                      {verifying ? "Verifying..." : "Verify OTP"}
                    </Button>
                  </>
                ) : (
                  <div className="space-y-4 text-center">
                    <p className="text-muted-foreground text-sm">
                      Check your console (for development) or email for the
                      magic link
                    </p>
                    <Button
                      className="w-full"
                      onClick={() => setUseOtp(true)}
                      variant="outline"
                    >
                      Enter OTP Code Instead
                    </Button>
                  </div>
                )}

                <div className="text-center text-sm">
                  <a
                    className="underline underline-offset-4"
                    href="/auth/login"
                  >
                    Back to login
                  </a>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
