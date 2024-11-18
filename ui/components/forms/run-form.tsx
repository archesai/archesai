"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import { siteConfig } from "@/config/site";
import {
  useContentControllerFindAll,
  useRunsControllerCreate,
  useToolsControllerFindAll,
} from "@/generated/archesApiComponents";
import {
  ContentEntity,
  CreateRunDto,
  ToolEntity,
} from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import { useState } from "react";
import * as z from "zod";

import { DataSelector } from "../data-selector";

const formSchema = z.object({
  contentIds: z.array(z.string()).optional(),
  text: z.string().optional(),
  toolId: z.string().min(1, "Tool selection is required"),
});

export default function RunForm({
  preSelectedTool,
}: {
  preSelectedTool?: ToolEntity;
}) {
  const { defaultOrgname } = useAuth();
  const { mutateAsync: runTool } = useRunsControllerCreate();
  const [selectedTool, setSelectedTool] = useState<ToolEntity | undefined>(
    preSelectedTool
  );
  const [selectedContent, setSelectedContent] = useState<ContentEntity>();

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: preSelectedTool?.id,
      description:
        "This is the role that will be used for this member. Note that different roles have different permissions.",
      label: "Tool",
      name: "toolId",
      renderControl: (field) => (
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
          itemType="tool"
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
      ),
      validationRule: formSchema.shape.toolId,
    },
    {
      component: Input,
      description:
        "This is the role that will be used for this member. Note that different roles have different permissions.",
      label: "Content",
      name: "contentIds",
      renderControl: (field) => (
        <DataSelector<ContentEntity>
          isMultiSelect={true}
          itemType="content"
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
      ),
      validationRule: formSchema.shape.contentIds,
    },
  ];

  return (
    <GenericForm<CreateRunDto, any>
      description={
        "Run a tool on a piece of content. This will create a new run of the tool on the selected content."
      }
      fields={formFields}
      isUpdateForm={false}
      itemType="tool run"
      onSubmitCreate={async (createTooRunDto, mutateOptions) => {
        await runTool(
          {
            body: {
              contentIds: createTooRunDto.contentIds,
              runType: "TOOL_RUN",
              text: createTooRunDto.text,
              toolId: createTooRunDto.toolId,
            },
            pathParams: {
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      showCard={true}
      title="Try a Tool"
    />
  );
}
