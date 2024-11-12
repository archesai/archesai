import { Alert, AlertTitle } from "@/components/ui/alert";
import { useAuthControllerEmailVerificationRequest } from "@/generated/archesApiComponents";
import { RocketIcon } from "@radix-ui/react-icons";

import { useToast } from "../../ui/use-toast";

export function VerifyEmailAlert() {
  const { mutateAsync: requestEmailVerification } =
    useAuthControllerEmailVerificationRequest();
  const { toast } = useToast();
  return (
    <Alert className="rounded-none border-none bg-primary">
      <RocketIcon className="h-5 w-5" color="white" />
      <AlertTitle className="font-normal text-white">
        <span className="flex gap-1">
          Please
          <div
            className="cursor-pointer font-semibold"
            onClick={async () => {
              try {
                await requestEmailVerification({});
                toast({
                  description:
                    "Please check your inbox for the verification email",
                  title: "Email verification sent",
                });
              } catch (error) {
                toast({
                  description: error as any,
                  title: "Error sending verification email",
                });
              }
            }}
          >
            {" "}
            verify your email address{" "}
          </div>{" "}
          to continue using the app.
        </span>
      </AlertTitle>
    </Alert>
  );
}
