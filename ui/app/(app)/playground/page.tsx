"use client";

import { DataSelector } from "@/components/data-selector";
import { RunStatusButton } from "@/components/run-status-button";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "@/components/ui/use-toast";
import { siteConfig } from "@/config/site";
import {
  useContentControllerFindAll,
  usePipelinesControllerCreatePipelineRun,
  useToolsControllerFindAll,
} from "@/generated/archesApiComponents";
import {
  ContentEntity,
  PipelineRunEntity,
  ToolEntity,
} from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { zodResolver } from "@hookform/resolvers/zod";
import { CounterClockwiseClockIcon } from "@radix-ui/react-icons";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z
  .object({
    runInputContentIds: z.array(z.string()).optional(),
    text: z.string().optional(),
    toolId: z.string().nonempty("Tool selection is required"),
  })
  .refine(
    (data) =>
      (data.runInputContentIds && data.runInputContentIds.length > 0) ||
      (data?.text?.trim()?.length || -1) > 0,
    {
      message: "Either content inputs or text must be provided.",
      path: ["runInputContentIds", "text"], // You can choose one or both paths for the error
    }
  );

type FormValues = z.infer<typeof formSchema>;

export default function PlaygroundPage() {
  const { defaultOrgname } = useAuth();
  const { mutateAsync: runPipeline } =
    usePipelinesControllerCreatePipelineRun();
  const [selectedTool, setSelectedTool] = useState<ToolEntity>();
  const [selectedContent, setSelectedContent] = useState<ContentEntity>();
  const [, setCurrentRun] = useState<PipelineRunEntity>();
  const form = useForm<FormValues>({
    defaultValues: {
      runInputContentIds: [],
      text: "",
      toolId: "",
    },
    resolver: zodResolver(formSchema),
  });
  // const { data: runs } = useRunsControllerFindAll(
  //   {
  //     pathParams: {
  //       orgname: defaultOrgname,
  //     },
  //     queryParams: {
  //       filters: JSON.stringify([
  //         {
  //           field: "type",
  //           operator: "equals",
  //           value: "TOOL_RUN",
  //         },
  //         {
  //           field: "toolId",
  //           operator: "equals",
  //           value: selectedTool?.id || "",
  //         },
  //       ]) as any,
  //     },
  //   },
  //   {
  //     enabled: !!selectedTool,
  //   }
  // );

  // const { data: runDetailed } = useRunsControllerFindOne(
  //   {
  //     pathParams: {
  //       id: currentRun?.id || "",
  //       orgname: defaultOrgname,
  //     },
  //   },
  //   {
  //     enabled: !!currentRun,
  //   }
  // );

  return (
    <Form {...form}>
      <form
        className="grid h-full gap-3 md:grid-cols-3"
        onSubmit={form.handleSubmit(async (values) => {
          const run = await runPipeline(
            {
              body: {
                runInputContentIds: values.runInputContentIds,
                text: values.text,
                url: "",
              },
              pathParams: {
                orgname: defaultOrgname,
                pipelineId: values.toolId,
              },
            },
            {
              onError: (error) => {
                toast({
                  description: error?.stack.message,
                  title: "Error",
                });
              },
              onSuccess: () => {
                toast({
                  description: "Tool run successful",
                  title: "Success",
                });
              },
            }
          );
          setCurrentRun(run);
        })}
      >
        {/* OUTPUT */}
        <div className="col-span-2 flex flex-1 flex-col gap-2">
          <div className="flex-1">{}</div>
          <div className="flex-1">{}</div>
        </div>

        {/* SIDEBAR */}
        <div className="col-span-1 flex flex-col gap-4">
          {/* Tool Selector */}
          <Controller
            control={form.control}
            name="toolId"
            render={({ field, fieldState }) => (
              <>
                <DataSelector<ToolEntity>
                  getItemDetails={(tool) => {
                    return (
                      <div className="grid gap-2">
                        <h4 className="flex items-center gap-1 font-medium leading-none">
                          {tool?.name}
                        </h4>
                        <div className="text-sm text-muted-foreground">
                          {tool?.description}
                        </div>
                      </div>
                    );
                  }}
                  icons={[
                    {
                      Icon: siteConfig.toolBaseIcons["extract-text"],
                      name: "Extract Text",
                    },
                    {
                      Icon: siteConfig.toolBaseIcons["create-embeddings"],
                      name: "Create Embeddings",
                    },
                    {
                      Icon: siteConfig.toolBaseIcons["summarize"],
                      name: "Summarize",
                    },
                    {
                      Icon: siteConfig.toolBaseIcons["text-to-image"],
                      name: "Text to Image",
                    },
                    {
                      Icon: siteConfig.toolBaseIcons["text-to-speech"],
                      name: "Text to Speech",
                    },
                  ]}
                  isMultiSelect={false}
                  label="Tool"
                  selectedData={selectedTool}
                  setSelectedData={(tool: any) => {
                    setSelectedTool(tool);
                    field.onChange(tool.id);
                  }}
                  useFindAll={() =>
                    useToolsControllerFindAll({
                      pathParams: {
                        orgname: defaultOrgname,
                      },
                    })
                  }
                />

                {fieldState.error && (
                  <span className="text-sm text-red-500">
                    {(fieldState.error as any)?.toolId?.message}
                  </span>
                )}
              </>
            )}
          />

          {/* Content Selector */}
          <Controller
            control={form.control}
            name="runInputContentIds"
            render={({ field, fieldState }) => (
              <>
                <DataSelector<ContentEntity>
                  isMultiSelect={true}
                  label="Content"
                  selectedData={selectedContent}
                  setSelectedData={(content: any) => {
                    setSelectedContent(content);
                    field.onChange(
                      content === null ? [] : content.map((c: any) => c.id)
                    );
                  }}
                  useFindAll={() =>
                    useContentControllerFindAll({
                      pathParams: {
                        orgname: defaultOrgname,
                      },
                    })
                  }
                />
                {fieldState.error && (
                  <span className="text-sm text-red-500">
                    {(fieldState.error as any)?.runInputContentIds?.message}
                  </span>
                )}
              </>
            )}
          />

          {/* Input Text */}
          <div className="flex flex-1 flex-col space-y-2">
            <Label htmlFor="text">Input Text</Label>
            <Textarea
              id="text"
              placeholder="Add text as input..."
              {...form.register("text")}
              className={form.formState.errors.text ? "border-red-500" : ""}
            />
            {form.formState.errors.text && (
              <span className="text-sm text-red-500">
                {form.formState.errors.text.message}
              </span>
            )}
          </div>
          {/* Submit Button */}
          <div className="flex items-center justify-end space-x-2">
            <div>{form.formState.errors.root?.message}</div>
            <Button disabled={!selectedTool} size="sm" type="submit">
              Submit
            </Button>
            <Button variant="secondary">
              <span className="sr-only">Show history</span>
              <CounterClockwiseClockIcon className="h-4 w-4" />
            </Button>
          </div>
          {/* Tool Runs */}
          <Label>Tool Runs</Label>
          <div>
            {[].map((run, i) => (
              <RunStatusButton
                key={i}
                onClick={() => setCurrentRun(run)}
                run={run}
              />
            ))}
          </div>
        </div>
      </form>
    </Form>
  );
}
