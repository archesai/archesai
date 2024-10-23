import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { JobEntity } from "@/generated/archesApiSchemas";
import { CheckCircle2, CircleX, Loader2Icon } from "lucide-react";
import { useState } from "react";

export const JobStatusButton = ({ job }: { job: JobEntity }) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const renderIcon = () => {
    switch (job.status) {
      case "COMPLETE":
        return <CheckCircle2 className="text-teal-600" />;
      case "ERROR":
        return <CircleX className="text-rose-600" />;
      case "PROCESSING":
        return <Loader2Icon className="animate-spin text-primary" />;
      default:
        return null;
    }
  };

  return (
    <Popover onOpenChange={setIsPopoverOpen} open={isPopoverOpen}>
      <PopoverTrigger asChild>
        <Button className="flex items-center" size="sm" variant="ghost">
          {renderIcon()}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-4 text-nowrap flex-1 w-auto text-sm">
        <div>
          <strong className="font-semibold">Status:</strong> {job.status}
        </div>
        <div>
          <strong className="font-semibold">Started:</strong>{" "}
          {new Date(job.startedAt).toLocaleString()}
        </div>
        <div>
          <strong className="font-semibold">Completed:</strong>{" "}
          {job.completedAt ? new Date(job.completedAt).toLocaleString() : "N/A"}
        </div>
        {job.completedAt && (
          <div>
            <strong className="font-semibold">Duration:</strong>{" "}
            {Math.round(
              (new Date(job.completedAt).getTime() -
                new Date(job.startedAt).getTime()) /
                1000
            ) + " seconds"}
          </div>
        )}

        <div>
          <strong className="font-semibold">Progress:</strong>{" "}
          {Math.round(job.progress * 100)}%
        </div>
        {job.error && (
          <div>
            <strong className="font-semibold">Error:</strong> {job.error}
          </div>
        )}
      </PopoverContent>
    </Popover>
  );
};
