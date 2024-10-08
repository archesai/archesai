import { CheckSquare2Icon, Loader2Icon } from "lucide-react";

export const StatusToIcon = ({
  status,
}: {
  status: "COMPLETE" | "ERROR" | "PROCESSING" | "QUEUED";
}) => {
  switch (status) {
    case "COMPLETE":
      return <CheckSquare2Icon className="text-primary" />;
    case "ERROR":
      return (
        <svg
          className="text-destructive"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            d="M6 18L18 6M6 6l12 12"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
          />
        </svg>
      );
    case "PROCESSING":
      return <Loader2Icon className="animate-spin text-primary" />;
  }
};
