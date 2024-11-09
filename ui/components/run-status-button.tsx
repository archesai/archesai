import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { siteConfig } from "@/config/site";
import { PipelineRunEntity } from "@/generated/archesApiSchemas";
import { CounterClockwiseClockIcon } from "@radix-ui/react-icons";
import { Ban, CheckCircle2, Loader2Icon } from "lucide-react";
import { useState } from "react";

export const RunStatusButton = ({
  onClick,
  run,
}: {
  onClick?: () => void;
  run: PipelineRunEntity;
}) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const renderIcon = () => {
    switch (run.status) {
      case "QUEUED":
        return <CounterClockwiseClockIcon className="text-primary" />;
      case "COMPLETE":
        return <CheckCircle2 className="text-green-600" />;
      case "ERROR":
        return <Ban className="text-red-600" />;
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

  const Icon = siteConfig.toolBaseIcons["text-to-image"];

  return (
    <Popover onOpenChange={setIsPopoverOpen} open={isPopoverOpen}>
      <PopoverTrigger asChild>
        <Button
          className="my-1 flex w-full items-center justify-between"
          onClick={onClick}
          size="sm"
          variant="secondary"
        >
          <div className="flex flex-1 items-center justify-start gap-1 overflow-hidden truncate">
            <Icon className="h-4 w-4 shrink-0" />
            {run.name}
          </div>
          <div className="ml-2 flex-shrink-0">{renderIcon()}</div>
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
