"use client";
import { RunStatusButton } from "@/components/run-status-button";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  useRunsControllerFindAll,
  useToolsControllerRun,
} from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { CounterClockwiseClockIcon } from "@radix-ui/react-icons";
import { useState } from "react";

import { CodeViewer } from "./components/code-viewer";
import { ModelSelector } from "./components/model-selector";
import { PresetActions } from "./components/preset-actions";
import { PresetSave } from "./components/preset-save";
import { PresetSelector } from "./components/preset-selector";
import { PresetShare } from "./components/preset-share";
import { presets } from "./data/presets";

export default function PlaygroundPage() {
  const { defaultOrgname } = useAuth();
  const [selectedTool, setSelectedTool] = useState<ToolEntity>();
  const { mutateAsync: runTool } = useToolsControllerRun();
  const { data: runs } = useRunsControllerFindAll(
    {
      pathParams: {
        orgname: defaultOrgname,
      },
      queryParams: {
        toolId: selectedTool?.id,
      },
    },
    {
      enabled: !!selectedTool,
    }
  );

  console.log(runs);
  return (
    <div
      className="flex flex-col gap-4"
      style={{
        height: "calc(100% - 64px)",
      }}
    >
      {
        // PRESET SELECTOR CODE VIEWER
      }
      <div className="-mt-0 flex justify-end gap-2 lg:-mt-12">
        <PresetSelector presets={presets} />
        <PresetSave />
        <div className="hidden space-x-2 md:flex">
          <CodeViewer />
          <PresetShare />
        </div>
        <PresetActions />
      </div>

      {
        // MAIN
      }
      <div className="grid h-full items-stretch gap-6 md:grid-cols-[1fr_200px]">
        {
          // FORM
        }
        <div className="hidden flex-col gap-3 sm:flex md:order-2">
          <ModelSelector
            selectedTool={selectedTool}
            setSelectedTool={setSelectedTool}
          />
          <Label>Runs</Label>
          {runs?.results?.map((run) => (
            <RunStatusButton key={run.id} run={run} />
          ))}
        </div>

        {
          // CONTENT
        }
        <div className="flex flex-1 md:order-1">
          <div className="flex flex-1 flex-col space-y-4">
            <div className="grid flex-1 gap-6 lg:grid-cols-2">
              <div className="flex flex-col space-y-4">
                <div className="flex flex-1 flex-col space-y-2">
                  <Label htmlFor="input">Input</Label>
                  <Textarea
                    className="flex-1"
                    id="input"
                    placeholder="We is going to the market."
                  />
                </div>
                <div className="flex flex-1 flex-col space-y-2">
                  <Label htmlFor="instructions">Instructions</Label>
                  <Textarea
                    className="flex-1"
                    id="instructions"
                    placeholder="Fix the grammar."
                  />
                </div>
              </div>
              <div className="rounded-md border bg-muted" />
            </div>
            <div className="flex items-center space-x-2">
              <Button
                disabled={!selectedTool}
                onClick={async () => {
                  await runTool({
                    body: {
                      text: "We is going to the market.",
                    },
                    pathParams: {
                      orgname: defaultOrgname,
                      toolId: selectedTool?.id || "",
                    },
                  });
                }}
              >
                Submit
              </Button>
              <Button variant="secondary">
                <span className="sr-only">Show history</span>
                <CounterClockwiseClockIcon className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
