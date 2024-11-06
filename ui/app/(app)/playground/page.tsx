"use client";

import { RunStatusButton } from "@/components/run-status-button";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  useRunsControllerFindAll,
  useToolsControllerRun,
} from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { zodResolver } from "@hookform/resolvers/zod";
import { CounterClockwiseClockIcon } from "@radix-ui/react-icons";
import { useEffect, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";

import { ContentSelector } from "./components/content-selector"; // We'll create this component next
import { ToolSelector } from "./components/tool-selector";

const formSchema = z
  .object({
    runContentInputIds: z.array(z.string()).optional(),
    text: z.string().optional(),
    toolId: z.string().nonempty("Tool selection is required"),
  })
  .refine(
    (data) =>
      (data.runContentInputIds && data.runContentInputIds.length > 0) ||
      (data?.text?.trim()?.length || -1) > 0,
    {
      message: "Either content inputs or text must be provided.",
      path: ["runContentInputIds", "text"], // You can choose one or both paths for the error
    }
  );

type FormValues = z.infer<typeof formSchema>;

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

  const form = useForm<FormValues>({
    defaultValues: {
      runContentInputIds: [],
      text: "",
      toolId: selectedTool?.id || "",
    },
    resolver: zodResolver(formSchema),
  });

  // Update toolId in form when selectedTool changes
  useEffect(() => {
    form.setValue("toolId", selectedTool?.id || "");
  }, [selectedTool]);

  return (
    <Form {...form}>
      <form
        className="grid h-full gap-3 md:grid-cols-3"
        onSubmit={form.handleSubmit(async (values) => {
          console.log(values);
          await runTool({
            body: {
              runInputContentIds: values.runContentInputIds,
              text: values.text,
              url: "",
            },
            pathParams: {
              orgname: defaultOrgname,
              toolId: values.toolId,
            },
          });
        })}
      >
        {/* OUTPUT */}
        <div className="col-span-2 flex flex-1 flex-col gap-2">
          <Label htmlFor="output">Output</Label>
          <Textarea
            className="flex-1 rounded-md border bg-muted"
            disabled
            id="output"
            placeholder="Output will appear here."
          />
        </div>

        {/* SIDEBAR */}
        <div className="col-span-1 flex flex-col gap-4">
          {/* Tool Selector */}
          <Controller
            control={form.control}
            name="toolId"
            render={({ field, fieldState }) => (
              <>
                <ToolSelector
                  selectedTool={selectedTool}
                  setSelectedTool={(tool) => {
                    setSelectedTool(tool);
                    field.onChange(tool.id);
                  }}
                />
                {fieldState.error && (
                  <span className="text-sm text-red-500">
                    {fieldState.error.message}
                  </span>
                )}
              </>
            )}
          />

          {/* Content Selector */}
          <Controller
            control={form.control}
            name="runContentInputIds"
            render={({ field, fieldState }) => (
              <>
                <ContentSelector
                  selectedContentIds={field.value as string[]}
                  setSelectedContentIds={field.onChange}
                />
                {fieldState.error && (
                  <span className="text-sm text-red-500">
                    {fieldState.error.message}
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
            {runs?.results?.map((run) => (
              <RunStatusButton key={run.id} run={run} />
            ))}
          </div>
        </div>
      </form>
    </Form>
  );
}
