"use client";

import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import { siteConfig } from "@/config/site";
import { cn } from "@/lib/utils";
import { type DialogProps } from "@radix-ui/react-dialog";
import { LaptopIcon, MoonIcon, SunIcon } from "@radix-ui/react-icons";
import * as VisuallyHidden from "@radix-ui/react-visually-hidden";
import { useRouter } from "next/navigation";
import { useTheme } from "next-themes";
import * as React from "react";

import { DialogDescription, DialogTitle } from "./ui/dialog";

export function CommandMenu({ ...props }: DialogProps) {
  const router = useRouter();
  const [open, setOpen] = React.useState(false);
  const { setTheme } = useTheme();

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || e.key === "/") {
        if (
          (e.target instanceof HTMLElement && e.target.isContentEditable) ||
          e.target instanceof HTMLInputElement ||
          e.target instanceof HTMLTextAreaElement ||
          e.target instanceof HTMLSelectElement
        ) {
          return;
        }

        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const runCommand = React.useCallback((command: () => unknown) => {
    setOpen(false);
    command();
  }, []);

  return (
    <div>
      <Button
        className={cn(
          "relative h-8 w-full justify-start rounded-[0.5rem] bg-muted/50 text-sm font-normal text-muted-foreground shadow-none pr-12"
        )}
        onClick={() => setOpen(true)}
        variant="outline"
        {...props}
      >
        <span className="hidden lg:inline-flex">
          Type a command or search...
        </span>
        <span className="inline-flex lg:hidden">Search...</span>
        <kbd className="pointer-events-none absolute right-[0.3rem] top-[0.3rem] h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>
      <CommandDialog onOpenChange={setOpen} open={open}>
        <VisuallyHidden.Root>
          <DialogTitle />
        </VisuallyHidden.Root>
        <DialogDescription />
        <CommandInput placeholder="Type a command or search..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          {Object.entries(siteConfig.links).map(([heading, links]) => (
            <CommandGroup
              heading={heading.charAt(0).toUpperCase() + heading.slice(1)}
              key={heading}
            >
              {links.map((navItem) => (
                <CommandItem
                  key={navItem.href}
                  onClick={() => {
                    runCommand(() => router.push(navItem.href as string));
                  }}
                  onSelect={() => {
                    runCommand(() => router.push(navItem.href as string));
                  }}
                  value={navItem.title}
                >
                  <navItem.Icon className="mr-2 h-4 w-4" />
                  {navItem.title}
                </CommandItem>
              ))}
            </CommandGroup>
          ))}
          <CommandSeparator />
          <CommandGroup heading="Theme">
            <CommandItem onSelect={() => runCommand(() => setTheme("light"))}>
              <SunIcon className="mr-2 h-4 w-4" />
              Light
            </CommandItem>
            <CommandItem onSelect={() => runCommand(() => setTheme("dark"))}>
              <MoonIcon className="mr-2 h-4 w-4" />
              Dark
            </CommandItem>
            <CommandItem onSelect={() => runCommand(() => setTheme("system"))}>
              <LaptopIcon className="mr-2 h-4 w-4" />
              System
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </div>
  );
}
