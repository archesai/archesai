// src/components/content-selector.tsx

"use client";

import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useContentControllerFindAll } from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { CheckIcon, PlusIcon } from "@radix-ui/react-icons";
import * as React from "react";

interface ContentSelectorProps {
  selectedContentIds: string[];
  setSelectedContentIds: (ids: string[]) => void;
}

export function ContentSelector({
  selectedContentIds,
  setSelectedContentIds,
}: ContentSelectorProps) {
  const { defaultOrgname } = useAuth();
  const { data: contents } = useContentControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
  });
  const [open, setOpen] = React.useState(false);
  const [searchTerm, setSearchTerm] = React.useState("");

  const filteredContents =
    contents?.results?.filter((content) =>
      content.name.toLowerCase().includes(searchTerm.toLowerCase())
    ) ?? [];

  const toggleSelection = (id: string) => {
    if (selectedContentIds.includes(id)) {
      setSelectedContentIds(selectedContentIds.filter((cid) => cid !== id));
    } else {
      setSelectedContentIds([...selectedContentIds, id]);
    }
  };

  return (
    <div className="grid gap-2">
      <Label>Content</Label>
      <Popover onOpenChange={setOpen} open={open}>
        <PopoverTrigger asChild>
          <Button
            aria-expanded={open}
            aria-label="Select contents"
            className="w-full justify-between"
            role="combobox"
            variant="outline"
          >
            {selectedContentIds.length > 0
              ? `${selectedContentIds.length} selected`
              : "Select contents..."}
            <PlusIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[300px] p-0">
          <Command>
            <CommandInput
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setSearchTerm(e.target.value)
              }
              placeholder="Search contents..."
            />
            <CommandList className="max-h-[400px] overflow-y-auto">
              <CommandEmpty>No contents found.</CommandEmpty>
              {filteredContents.map((content) => (
                <CommandItem
                  className="flex items-center justify-between"
                  key={content.id}
                  onSelect={() => toggleSelection(content.id)}
                >
                  <span>{content.name}</span>
                  {selectedContentIds.includes(content.id) && (
                    <CheckIcon className="h-4 w-4 text-green-500" />
                  )}
                </CommandItem>
              ))}
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {/* Display selected contents */}
      {selectedContentIds.length ? (
        <div className="mt-2 flex flex-wrap gap-2">
          {selectedContentIds.map((id) => {
            const content = contents?.results?.find((c) => c.id === id);
            return (
              <span
                className="inline-flex items-center rounded bg-blue-100 px-2 py-1 text-sm text-blue-700"
                key={id}
              >
                {content?.name || id}
                <button
                  className="ml-1 text-red-500 hover:text-red-700"
                  onClick={() => toggleSelection(id)}
                  type="button"
                >
                  ×
                </button>
              </span>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}
