import { SearchIcon } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { Button } from "#components/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "#components/shadcn/dialog";
import { Input } from "#components/shadcn/input";
import { ClientOnly } from "./client-only";

const detectOS = () => {
  if (typeof window === "undefined") return "unknown";
  const userAgent = window.navigator.userAgent.toLowerCase();
  if (userAgent.includes("mac")) return "mac";
  if (userAgent.includes("win")) return "windows";
  if (userAgent.includes("linux")) return "linux";
  return "unknown";
};

export const Search = ({ className }: { className?: string }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const onClose = useCallback(() => {
    setIsOpen(false);
    setSearchQuery("");
  }, []);

  useEffect(() => {
    if (isOpen) {
      return;
    }

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "k" && (event.metaKey || event.ctrlKey)) {
        event.preventDefault();
        setIsOpen(true);
      }
    }

    window.addEventListener("keydown", onKeyDown);

    return () => {
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [isOpen]);

  const os = detectOS();
  const cmdKey = os === "mac" ? "⌘" : "Ctrl";

  return (
    <div className={className}>
      <Button
        className="relative h-8 w-full justify-start rounded-lg text-muted-foreground text-sm shadow-none sm:w-72"
        onClick={() => setIsOpen(true)}
        variant="outline"
      >
        <div className="flex grow items-center gap-2">
          <SearchIcon size={14} />
          Search
        </div>
        <ClientOnly>
          <kbd className="pointer-events-none ml-auto inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-medium font-mono text-[10px] text-muted-foreground opacity-100">
            <span className="text-xs">{cmdKey}</span>K
          </kbd>
        </ClientOnly>
      </Button>

      <Dialog
        onOpenChange={setIsOpen}
        open={isOpen}
      >
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle className="sr-only">Search</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <SearchIcon
                className="text-muted-foreground"
                size={18}
              />
              <Input
                autoFocus
                className="flex-1 border-0 focus-visible:ring-0 focus-visible:ring-offset-0"
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search documentation..."
                type="text"
                value={searchQuery}
              />
              <Button
                onClick={onClose}
                size="sm"
                variant="ghost"
              >
                Cancel
              </Button>
            </div>
            <div className="min-h-[200px] text-center text-muted-foreground text-sm">
              {searchQuery
                ? `Searching for "${searchQuery}"...`
                : "Type to search"}
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};
