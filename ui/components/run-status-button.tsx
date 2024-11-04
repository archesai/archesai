import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { RunEntity } from "@/generated/archesApiSchemas";
import { CheckCircle2, CircleX, Loader2Icon } from "lucide-react";
import { useState } from "react";

export const RunStatusButton = ({ run }: { run: RunEntity }) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const renderIcon = () => {
    switch (run.status) {
      case "COMPLETE":
        return <CheckCircle2 className="text-teal-600" />;
      case "ERROR":
        return <CircleX className="text-rose-600" />;
      case "PROCESSING":
        return (
          <div className="flex items-center gap-2">
            <Loader2Icon className="animate-spin text-primary" />
            <span>{(run.progress * 100).toFixed(0)}%</span>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <Popover onOpenChange={setIsPopoverOpen} open={isPopoverOpen}>
      <PopoverTrigger asChild>
        <Button className="flex items-center" size="icon" variant="link">
          {renderIcon()}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="overflow-auto p-4 text-sm">
        <div>
          <strong className="font-semibold">Tool:</strong> {run.toolId}
        </div>
        <div>
          <strong className="font-semibold">Status:</strong> {run.status}
        </div>
        <div>
          <strong className="font-semibold">Started:</strong>{" "}
          {run.startedAt && new Date(run.startedAt).toLocaleString()}
        </div>
        <div>
          <strong className="font-semibold">Completed:</strong>{" "}
          {run.completedAt ? new Date(run.completedAt).toLocaleString() : "N/A"}
        </div>
        {run.completedAt && (
          <div>
            <strong className="font-semibold">Duration:</strong>{" "}
            {run.startedAt &&
              run.completedAt &&
              new Date(run.completedAt).getTime() -
                new Date(run.startedAt).getTime()}
          </div>
        )}

        <div>
          <strong className="font-semibold">Progress:</strong>{" "}
          {Math.round(run.progress * 100)}%
        </div>
        {run.error && (
          <div>
            <strong className="font-semibold">Error:</strong> {run.error}
          </div>
        )}
      </PopoverContent>
    </Popover>
  );
};
