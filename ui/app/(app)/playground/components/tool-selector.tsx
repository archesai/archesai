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
import { siteConfig } from "@/config/site";
import { useToolsControllerFindAll } from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { cn } from "@/lib/utils";
import { CaretSortIcon, CheckIcon } from "@radix-ui/react-icons";
import * as React from "react";

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

const toolBaseIcons = siteConfig.toolBaseIcons;

export function ToolSelector({
  selectedTool,
  setSelectedTool,
}: {
  selectedTool?: ToolEntity;
  setSelectedTool: (model: ToolEntity) => void;
}) {
  const [open, setOpen] = React.useState(false);
  const [peekedTool, setPeekedTool] = React.useState<ToolEntity>();
  const { defaultOrgname } = useAuth();
  const { data: tools } = useToolsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
  });
  const toolBases = tools?.results?.map((tool) => tool.toolBase) ?? [];
  React.useEffect(() => {
    if (tools?.results?.length) {
      setSelectedTool(tools.results[0]);
    }
  }, [tools]);
  return (
    <div className="grid gap-2">
      <HoverCard openDelay={200}>
        <HoverCardTrigger asChild>
          <Label htmlFor="model">Tool</Label>
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
      <Popover onOpenChange={setOpen} open={open}>
        <PopoverTrigger asChild>
          <Button
            aria-expanded={open}
            aria-label="Select a tool"
            className="w-full justify-between"
            role="combobox"
            variant="outline"
          >
            <div className="flex items-center gap-1">
              {selectedTool?.toolBase &&
                Object.entries(toolBaseIcons)
                  .filter(([toolBase]) => toolBase === selectedTool.toolBase)
                  .map(([toolBase, Icon]) => (
                    <div className="flex items-center gap-2" key={toolBase}>
                      <Icon className="h-4 w-4 text-muted-foreground" />
                    </div>
                  ))}
              {selectedTool ? selectedTool.name : "Select a tool..."}
            </div>

            <CaretSortIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent align="end" className="w-[250px] p-0">
          <HoverCard>
            <HoverCardContent align="start" forceMount side="left">
              <div className="grid gap-2">
                <h4 className="flex items-center gap-1 font-medium leading-none">
                  {Object.entries(toolBaseIcons)
                    .filter(([toolBase]) => toolBase === peekedTool?.toolBase)
                    .map(([toolBase, Icon]) => (
                      <div className="flex items-center gap-2" key={toolBase}>
                        <Icon className="h-4 w-4 text-muted-foreground" />
                      </div>
                    ))}
                  {peekedTool?.name}
                </h4>
                <div className="text-sm text-muted-foreground">
                  {peekedTool?.description}
                </div>
              </div>
            </HoverCardContent>
            <Command loop>
              <CommandList className="h-[var(--cmdk-list-height)] max-h-[400px]">
                <CommandInput placeholder="Search tools..." />
                <CommandEmpty>No Tools found.</CommandEmpty>
                <HoverCardTrigger />
                {toolBases.map((toolBase) => (
                  <CommandGroup heading={toolBase} key={toolBase}>
                    {tools?.results
                      ?.filter((tool) => tool.toolBase === toolBase)
                      .map((tool) => (
                        <ToolItem
                          isSelected={selectedTool?.id === tool.id}
                          key={tool.id}
                          onPeek={(tool) => setPeekedTool(tool)}
                          onSelect={() => {
                            setSelectedTool(tool);
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

interface ToolItemProps {
  isSelected: boolean;
  onPeek: (model: ToolEntity) => void;
  onSelect: () => void;
  tool: ToolEntity;
}

function ToolItem({ isSelected, onPeek, onSelect, tool }: ToolItemProps) {
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
      className="gap-1 data-[selected=true]:bg-primary data-[selected=true]:text-primary-foreground"
      key={tool.id}
      onSelect={onSelect}
      ref={ref}
    >
      {Object.entries(toolBaseIcons)
        .filter(([toolBase]) => toolBase === tool?.toolBase)
        .map(([toolBase, Icon]) => (
          <div className="flex items-center gap-2" key={toolBase}>
            <Icon className="h-4 w-4 text-muted-foreground" />
          </div>
        ))}
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
