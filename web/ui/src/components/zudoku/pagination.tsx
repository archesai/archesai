import { ArrowLeftIcon, ArrowRightIcon } from "lucide-react";
import { Button } from "#components/shadcn/button";
import { Link } from "./Link";
import { cn } from "./utils";

export const Pagination = ({
  prev,
  next,
  className,
}: {
  prev: { to: string; label: string } | undefined;
  next: { to: string; label: string } | undefined;
  className?: string;
}) => {
  const linkClass =
    "group transition-all p-5 space-x-1 rtl:space-x-reverse transition-all hover:text-foreground";

  return (
    <div
      className={cn(
        "-mx-4 flex font-semibold text-muted-foreground",
        prev ? "justify-between" : "justify-end",
        className,
      )}
      data-pagefind-ignore="all"
    >
      {prev && (
        <Button
          asChild
          variant="ghost"
        >
          <Link
            className={linkClass}
            to={prev.to}
          >
            <ArrowLeftIcon
              size={14}
              strokeWidth={2.5}
            />
            <span className="truncate text-lg">{prev.label}</span>
          </Link>
        </Button>
      )}
      {next && (
        <Button
          asChild
          variant="ghost"
        >
          <Link
            className={linkClass}
            to={next.to}
          >
            <span className="truncate text-lg">{next.label}</span>
            <ArrowRightIcon
              size={14}
              strokeWidth={2.5}
            />
          </Link>
        </Button>
      )}
    </div>
  );
};
