import { CheckIcon, Loader2Icon, XIcon } from "lucide-react";

export const StatusToIcon = ({
  status,
}: {
  status: "COMPLETE" | "ERROR" | "PROCESSING" | "QUEUED";
}) => {
  switch (status) {
    case "COMPLETE":
      return <CheckIcon className="text-green-700" />;
    case "ERROR":
      return <XIcon className="text-red-500" />;
    case "PROCESSING":
      return <Loader2Icon className="animate-spin text-primary" />;
  }
};
