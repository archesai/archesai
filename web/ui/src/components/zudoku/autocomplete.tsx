import { useCommandState } from "cmdk";
import type { KeyboardEvent, Ref } from "react";
import { useRef, useState } from "react";
import {
  Command,
  CommandInput,
  CommandItem,
  CommandList,
} from "#components/shadcn/command";
import {
  Popover,
  PopoverAnchor,
  PopoverContent,
} from "#components/shadcn/popover";
import { cn } from "./utils";

type AutocompleteProps = {
  value: string;
  options: readonly string[];
  onChange: (e: string) => void;
  className?: string;
  placeholder?: string;
  onEnterPress?: (e: KeyboardEvent<HTMLInputElement>) => void;
  ref?: Ref<HTMLInputElement>;
  shouldFilter?: boolean;
};

const AutocompletePopover = ({
  value,
  options,
  onChange,
  className,
  placeholder = "Value",
  onEnterPress,
  ref,
}: AutocompleteProps) => {
  const [open, setOpen] = useState(false);
  const [dontClose, setDontClose] = useState(false);
  const count = useCommandState((state) => state.filtered.count);
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <Popover open={open}>
      <PopoverAnchor>
        <CommandInput
          autoComplete="off"
          className={cn("h-9 bg-transparent", className)}
          onBlur={() => {
            if (dontClose) {
              return;
            }
            setOpen(false);
          }}
          onFocus={() => setOpen(true)}
          onKeyDown={(e) => {
            if (e.key === "Escape") {
              e.preventDefault();
              setOpen(false);
            }
            if (e.key === "Enter" && count === 0) {
              onEnterPress?.(e);
            }
          }}
          onValueChange={onChange}
          placeholder={placeholder}
          ref={ref || inputRef}
          value={value}
        />
      </PopoverAnchor>
      <PopoverContent
        align="start"
        className="w-[var(--radix-popper-anchor-width)] p-0"
        onInteractOutside={() => setOpen(false)}
        onOpenAutoFocus={(e) => e.preventDefault()}
        side="bottom"
      >
        <Command shouldFilter={false}>
          <CommandList>
            {options
              .filter((o) => o.toLowerCase().includes(value.toLowerCase()))
              .map((option) => (
                <CommandItem
                  key={option}
                  onMouseDown={() => setDontClose(true)}
                  onMouseUp={() => setDontClose(false)}
                  onSelect={() => {
                    onChange(option);
                    setOpen(false);
                  }}
                  value={option}
                >
                  {option}
                </CommandItem>
              ))}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
};

export const Autocomplete = AutocompletePopover;
