import {
  Button,
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Loader2Icon,
  Separator,
} from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import UserForm from "#components/forms/user-form";
import {
  useDeleteUser,
  useGetSessionSuspense,
  useGetUserSuspense,
  useRequestPasswordReset,
} from "#lib/index";

export const Route = createFileRoute("/_app/profile/")({
  component: ProfileSecuritySettingsPage,
});

function ProfileSecuritySettingsPage(): JSX.Element {
  const { data: sessionData } = useGetSessionSuspense("current");
  const { data: userData } = useGetUserSuspense(sessionData.data.userID);
  const { isPending: deactivatePending, mutateAsync: deactivateAccount } =
    useDeleteUser();
  const {
    isPending: requestPasswordResetPending,
    mutateAsync: requestPasswordReset,
  } = useRequestPasswordReset();

  return (
    <div className="flex flex-col gap-4">
      <UserForm />
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Reset Password</CardTitle>
            <CardDescription>
              If you would like to change your password, please click the button
              below. It will send you an email with instructions on how to
              reset.
            </CardDescription>
          </CardHeader>
          <Separator />
          <CardFooter>
            <Button
              disabled={requestPasswordResetPending}
              onClick={async () => {
                await requestPasswordReset({
                  data: {
                    email: userData.data.email,
                  },
                });
              }}
              size={"sm"}
              type="submit"
            >
              {requestPasswordResetPending && (
                <Loader2Icon className="animate-spin" />
              )}
              Reset Password
            </Button>
          </CardFooter>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Deactivate Account</CardTitle>
            <CardDescription>
              If you would like to deactivate your account, please click the
              button below. This action is irreversible.
            </CardDescription>
          </CardHeader>
          <Separator />
          <CardFooter>
            <Button
              disabled={deactivatePending}
              onClick={async () => {
                await deactivateAccount({
                  id: userData.data.id,
                });
              }}
              size="sm"
              variant={"destructive"}
            >
              {deactivatePending && <Loader2Icon className="animate-spin" />}
              Delete Account
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
