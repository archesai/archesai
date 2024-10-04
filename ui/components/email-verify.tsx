import { Alert, AlertTitle } from "@/components/ui/alert";
import { useEmailVerificationControllerRequest } from "@/generated/archesApiComponents";
import { RocketIcon } from "@radix-ui/react-icons";

import { useToast } from "./ui/use-toast";

export function VerifyEmailAlert() {
  const { mutateAsync: requestEmailVerification } =
    useEmailVerificationControllerRequest();
  const { toast } = useToast();
  return (
    <Alert className="bg-primary border-none rounded-none">
      <RocketIcon className="h-4 w-4" color="white" />
      <AlertTitle className="text-white font-normal">
        <span className="flex gap-1">
          Please
          <div
            className="font-semibold cursor-pointer"
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
