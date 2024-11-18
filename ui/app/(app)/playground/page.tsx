"use client";

import RunForm from "@/components/forms/run-form";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { siteConfig } from "@/config/site";
import { useToolsControllerFindAll } from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import { ArrowLeft } from "lucide-react";
import { useState } from "react";

export default function PlaygroundPage() {
  const { defaultOrgname } = useAuth();
  const [selectedTool, setSelectedTool] = useState<ToolEntity>();
  const { data: tools } = useToolsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
  });

  if (!selectedTool) {
    return (
      <div className="flex h-full flex-col items-center justify-start">
        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          {tools?.results?.map((tool, index) => {
            const Icon = siteConfig.toolBaseIcons[tool.toolBase];
            return (
              <Card
                className="flex flex-col justify-between bg-sidebar text-center transition-shadow hover:shadow-lg"
                key={index}
              >
                <CardHeader className="pt-6">
                  <Icon className="mx-auto mb-2 h-8 w-8 text-foreground" />
                  <CardTitle className="text-xl font-semibold">
                    {tool.name}
                  </CardTitle>
                </CardHeader>
                <CardContent>{tool.description}</CardContent>
                <CardFooter className="justify-center">
                  <Button
                    className="h-8"
                    onClick={() => setSelectedTool(tool)}
                    variant={"secondary"}
                  >
                    Select Tool
                  </Button>
                </CardFooter>
              </Card>
            );
          })}
        </div>
      </div>
    );
  }
  return (
    <div className="relative grid h-full gap-3 md:grid-cols-3">
      {
        // Back Button
      }
      <Button
        className="absolute left-0 top-0"
        onClick={() => setSelectedTool(undefined)}
        size="sm"
        variant={"secondary"}
      >
        <ArrowLeft className="h-4 w-4" />
      </Button>
      {/* OUTPUT */}
      <div className="col-span-2 flex flex-1 flex-col gap-2">
        <div className="flex-1">{}</div>
        <div className="flex-1">{}</div>
      </div>

      {/* SIDEBAR */}
      <div className="col-span-1 flex flex-col gap-4">
        <RunForm preSelectedTool={selectedTool} />
      </div>
    </div>
  );
}
