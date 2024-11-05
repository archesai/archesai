"use client";

import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useToolsControllerFindAll } from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { cn } from "@/lib/utils";
import { CaretSortIcon, CheckIcon } from "@radix-ui/react-icons";
import { PopoverProps } from "@radix-ui/react-popover";
import * as React from "react";

import { Model, ModelType } from "../data/models";

const useMutationObserver = (
  ref: React.MutableRefObject<HTMLElement | null>,
  callback: MutationCallback,
  options = {
    attributes: true,
    characterData: true,
    childList: true,
    subtree: true,
  }
) => {
  React.useEffect(() => {
    if (ref.current) {
      const observer = new MutationObserver(callback);
      observer.observe(ref.current, options);
      return () => observer.disconnect();
    }
  }, [ref, callback, options]);
};

interface ModelSelectorProps extends PopoverProps {
  models: Model[];
  types: readonly ModelType[];
}

export function ModelSelector({ ...props }: ModelSelectorProps) {
  const [open, setOpen] = React.useState(false);
  const [selectedModel, setSelectedModel] = React.useState<ToolEntity>();
  const [peekedModel, setPeekedModel] = React.useState<ToolEntity>();
  const { defaultOrgname } = useAuth();
  const { data: tools } = useToolsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
  });
  const toolBases = tools?.results?.map((tool) => tool.toolBase) ?? [];

  return (
    <div className="grid gap-2">
      <HoverCard openDelay={200}>
        <HoverCardTrigger asChild>
          <Label htmlFor="model">Model</Label>
        </HoverCardTrigger>
        <HoverCardContent
          align="start"
          className="w-[260px] text-sm"
          side="left"
        >
          The model which will generate the completion. Some models are suitable
          for natural language tasks, others specialize in code. Learn more.
        </HoverCardContent>
      </HoverCard>
      <Popover onOpenChange={setOpen} open={open} {...props}>
        <PopoverTrigger asChild>
          <Button
            aria-expanded={open}
            aria-label="Select a tool"
            className="w-full justify-between"
            role="combobox"
            variant="outline"
          >
            {selectedModel ? selectedModel.name : "Select a tool..."}
            <CaretSortIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent align="end" className="w-[250px] p-0">
          <HoverCard>
            <HoverCardContent align="start" forceMount side="left">
              <div className="grid gap-2">
                <h4 className="font-medium leading-none">
                  {peekedModel?.name}
                </h4>
                <div className="text-sm text-muted-foreground">
                  {peekedModel?.description}
                </div>
                {/* {peekedModel?.strengths ? (
                  <div className="mt-4 grid gap-2">
                    <h5 className="text-sm font-medium leading-none">
                      Strengths
                    </h5>
                    <ul className="text-sm text-muted-foreground">
                      {peekedModel?.strengths}
                    </ul>
                  </div>
                ) : null} */}
              </div>
            </HoverCardContent>
            <Command loop>
              <CommandList className="h-[var(--cmdk-list-height)] max-h-[400px]">
                <CommandInput placeholder="Search Models..." />
                <CommandEmpty>No Models found.</CommandEmpty>
                <HoverCardTrigger />
                {toolBases.map((toolBase) => (
                  <CommandGroup heading={toolBase} key={toolBase}>
                    {tools?.results
                      ?.filter((tool) => tool.toolBase === toolBase)
                      .map((tool) => (
                        <ModelItem
                          isSelected={selectedModel?.id === tool.id}
                          key={tool.id}
                          onPeek={(tool) => setPeekedModel(tool)}
                          onSelect={() => {
                            setSelectedModel(tool);
                            setOpen(false);
                          }}
                          tool={tool}
                        />
                      ))}
                  </CommandGroup>
                ))}
              </CommandList>
            </Command>
          </HoverCard>
        </PopoverContent>
      </Popover>
    </div>
  );
}

interface ModelItemProps {
  isSelected: boolean;
  onPeek: (model: ToolEntity) => void;
  onSelect: () => void;
  tool: ToolEntity;
}

function ModelItem({ isSelected, onPeek, onSelect, tool }: ModelItemProps) {
  const ref = React.useRef<HTMLDivElement>(null);

  useMutationObserver(ref, (mutations) => {
    mutations.forEach((mutation) => {
      if (
        mutation.attributeName === "aria-selected" &&
        ref.current?.getAttribute("aria-selected") === "true"
      ) {
        onPeek(tool);
      }
    });
  });

  return (
    <CommandItem
      className="data-[selected=true]:bg-primary data-[selected=true]:text-primary-foreground"
      key={tool.id}
      onSelect={onSelect}
      ref={ref}
    >
      {tool.name}
      <CheckIcon
        className={cn(
          "ml-auto h-4 w-4",
          isSelected ? "opacity-100" : "opacity-0"
        )}
      />
    </CommandItem>
  );
}
