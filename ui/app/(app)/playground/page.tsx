import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { CounterClockwiseClockIcon } from "@radix-ui/react-icons";
import { Metadata } from "next";

import { CodeViewer } from "./components/code-viewer";
import { MaxLengthSelector } from "./components/maxlength-selector";
import { ModelSelector } from "./components/model-selector";
import { PresetActions } from "./components/preset-actions";
import { PresetSave } from "./components/preset-save";
import { PresetSelector } from "./components/preset-selector";
import { PresetShare } from "./components/preset-share";
import { TemperatureSelector } from "./components/temperature-selector";
import { TopPSelector } from "./components/top-p-selector";
import { models, types } from "./data/models";
import { presets } from "./data/presets";

export const metadata: Metadata = {
  description: "The OpenAI Playground built using the components.",
  title: "Playground",
};

export default function PlaygroundPage() {
  return (
    <div className="flex h-full flex-col">
      {
        // TOP PART
      }

      <div className="flex justify-end gap-2">
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
      <div className="container h-full py-6">
        <div className="grid h-full items-stretch gap-6 md:grid-cols-[1fr_200px]">
          {
            // SIDEBAR
          }
          <div className="hidden flex-col space-y-4 sm:flex md:order-2">
            <ModelSelector models={models} types={types} />
            <TemperatureSelector defaultValue={[0.56]} />
            <MaxLengthSelector defaultValue={[256]} />
            <TopPSelector defaultValue={[0.9]} />
          </div>

          {
            // MAIN CONTENT
          }
          <div className="md:order-1">
            <div className="flex flex-col space-y-4">
              <div className="grid h-full gap-6 lg:grid-cols-2">
                <div className="flex flex-col space-y-4">
                  <div className="flex flex-1 flex-col space-y-2">
                    <Label htmlFor="input">Input</Label>
                    <Textarea
                      className="flex-1 lg:min-h-[580px]"
                      id="input"
                      placeholder="We is going to the market."
                    />
                  </div>
                  <div className="flex flex-col space-y-2">
                    <Label htmlFor="instructions">Instructions</Label>
                    <Textarea
                      id="instructions"
                      placeholder="Fix the grammar."
                    />
                  </div>
                </div>
                <div className="mt-[21px] min-h-[400px] rounded-md border bg-muted lg:min-h-[700px]" />
              </div>
              <div className="flex items-center space-x-2">
                <Button>Submit</Button>
                <Button variant="secondary">
                  <span className="sr-only">Show history</span>
                  <CounterClockwiseClockIcon className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
