import { CheckSquare2Icon, Loader2Icon, XIcon } from "lucide-react";

export const StatusToIcon = ({
  status,
}: {
  status: "COMPLETE" | "ERROR" | "PROCESSING" | "QUEUED";
}) => {
  switch (status) {
    case "COMPLETE":
      return <CheckSquare2Icon className="text-primary" />;
    case "ERROR":
      return <XIcon className="text-red-500" />;
    case "PROCESSING":
      return <Loader2Icon className="animate-spin text-primary" />;
  }
};
