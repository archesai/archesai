"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useToast } from "@/components/ui/use-toast";
import {
  usePasswordResetControllerRequest,
  useUserControllerDeactivate,
  useUserControllerFindOne,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { ReloadIcon } from "@radix-ui/react-icons";

export default function ProfileSecuritySettingsPage() {
  const { data: user } = useUserControllerFindOne({});
  const { toast } = useToast();
  const { logout } = useAuth();

  const { isPending: deactivatePending, mutateAsync: deactivateAccount } =
    useUserControllerDeactivate();
  const {
    isPending: requestPasswordResetPending,
    mutateAsync: requestPasswordReset,
  } = usePasswordResetControllerRequest();

  return (
    <div className="stack gap-3">
      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Reset Password</CardTitle>
          <CardDescription>
            If you would like to change your password, please click the button
            below. It will send you an email with instructions on how to reset.
          </CardDescription>
        </CardHeader>
        <CardFooter className="flex justify-between">
          <Button
            className="w-full h-8"
            disabled={requestPasswordResetPending}
            onClick={async () =>
              await requestPasswordReset(
                {
                  body: {
                    email: user?.email as string,
                  },
                },
                {
                  onError: (err) => {
                    toast({
                      description:
                        err?.stack.msg ||
                        "An error occurred while trying to reset your password.",
                      title: "Error",
                    });
                  },
                  onSuccess: () => {
                    toast({
                      description:
                        "We have sent you an email with instructions on how to reset your password.",
                      title: "Email Sent",
                    });
                  },
                }
              )
            }
          >
            {requestPasswordResetPending && (
              <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
            )}
            Reset Password
          </Button>
        </CardFooter>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Deactivate Account</CardTitle>
          <CardDescription>
            If you would like to deactivate your account, please click the
            button below. This action is irreversible.
          </CardDescription>
        </CardHeader>
        <CardFooter className="flex justify-between">
          <Button
            className="w-full h-8"
            disabled={deactivatePending}
            onClick={async () =>
              await deactivateAccount(
                {},
                {
                  onError: (err) => {
                    toast({
                      description:
                        err?.stack.msg ||
                        "An error occurred while trying to deactivate your account.",
                      title: "Error",
                    });
                  },
                  onSuccess: async () => {
                    await logout();
                  },
                }
              )
            }
            variant={"destructive"}
          >
            {deactivatePending && (
              <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
            )}
            Delete Account
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
