"use client";

import { GridView } from "@/components/datatable/grid-view";
import RunForm from "@/components/forms/run-form";
import { RunStatusButton } from "@/components/run-status-button";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { siteConfig } from "@/config/site";
import {
  useContentControllerFindAll,
  useRunsControllerFindOne,
  useToolsControllerFindAll,
} from "@/generated/archesApiComponents";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import { usePlayground } from "@/hooks/use-playground";
import { useRouter } from "next/navigation";
// import { ArrowLeft } from "lucide-react";

export default function PlaygroundPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();
  const {
    // clearParams,
    selectedContent,
    selectedRunId,
    selectedTool,
    setSelectedTool,
  } = usePlayground();

  const { data: tools } = useToolsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
  });
  const { data: run } = useRunsControllerFindOne(
    {
      pathParams: {
        orgname: defaultOrgname,
        runId: selectedRunId || "",
      },
    },
    {
      enabled: !!selectedRunId,
    }
  );

  const { data: inputs } = useContentControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
    queryParams: {
      filters: JSON.stringify([
        {
          field: "id",
          operator: "in",
          value: run
            ? run.inputs.map((r) => r.id)
            : selectedContent?.map((r) => r.id) || [],
        },
      ]) as any,
    },
  });

  const { data: outputs } = useContentControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
    queryParams: {
      filters: JSON.stringify([
        {
          field: "id",
          operator: "in",
          value: run?.outputs?.map((r) => r.id) || [],
        },
      ]) as any,
    },
  });

  if (!selectedTool) {
    return (
      <div className="grid grid-cols-1 gap-6 p-0 md:grid-cols-3">
        {tools?.results?.map((tool, index) => {
          const Icon = siteConfig.toolBaseIcons[tool.toolBase];
          return (
            <Card
              className="flex flex-col justify-between gap-2 bg-sidebar p-4 text-center transition-shadow hover:shadow-lg"
              key={index}
            >
              <Icon className="mx-auto h-8 w-8 text-primary/80" />
              <div className="text-lg font-semibold">{tool.name}</div>
              <div className="text-sm font-normal"> {tool.description}</div>

              <Button
                className="mt-1 h-8"
                onClick={() => setSelectedTool(tool)}
                variant={"outline"}
              >
                Select Tool
              </Button>
            </Card>
          );
        })}
      </div>
    );
  }
  return (
    <div className="relative grid h-full gap-3 md:grid-cols-3">
      <div className="col-span-2 flex flex-1 flex-col gap-4">
        {inputs?.results?.length ? (
          <>
            <Label>Input Content</Label>
            <GridView<ContentEntity>
              content={(item) => (
                <div className="p-1">
                  <div className="font-mono">{item.name}</div>
                  <div className="text-sm text-muted-foreground">
                    {item.text || item.url}
                  </div>
                </div>
              )}
              data={inputs?.results || []}
              DataIcon={<div>Content</div>}
              deleteItem={async () => {}}
              getDeleteVariablesFromItem={() => {}}
              handleSelect={(content) => {
                router.push(`/content/single?contentId=${content.id}`);
              }}
              itemType={"content"}
              selectedItems={[]}
              setFinalForm={() => {}}
              setFormOpen={() => {}}
              toggleSelection={() => {}}
            />
          </>
        ) : null}

        {outputs?.results?.length ? (
          <>
            <Separator />
            <Label>Output Content</Label>
            <GridView<ContentEntity>
              content={(item) => (
                <div className="p-1">
                  <div className="font-mono">{item.name}</div>
                  <div className="text-sm text-muted-foreground">
                    {item.text || item.url}
                  </div>
                </div>
              )}
              data={outputs?.results || []}
              DataIcon={<div>Content</div>}
              deleteItem={async () => {}}
              getDeleteVariablesFromItem={() => {}}
              handleSelect={() => {}}
              itemType={"content"}
              selectedItems={[]}
              setFinalForm={() => {}}
              setFormOpen={() => {}}
              toggleSelection={() => {}}
            />
          </>
        ) : null}

        {!selectedContent && (
          <>
            <div className="flex flex-1 flex-col items-center justify-center text-muted-foreground">
              <p>Use the form to start using tools.</p>
            </div>
          </>
        )}
      </div>

      {/* SIDEBAR */}
      <div className="col-span-1 flex flex-col gap-1">
        {selectedRunId && run && <RunStatusButton run={run} />}
        <RunForm />
      </div>
    </div>
  );
}
