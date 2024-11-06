"use client";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useRunsControllerFindOne } from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { useSearchParams } from "next/navigation";
import React from "react";

const RunDetailsPage: React.FC = () => {
  const searchParams = useSearchParams();
  const runId = searchParams?.get("runId");

  const { defaultOrgname } = useAuth();

  const { data: detailedRun, error } = useRunsControllerFindOne({
    pathParams: {
      id: runId as string,
      orgname: defaultOrgname,
    },
  });

  if (error || !detailedRun) {
    return (
      <div className="p-4">
        <Alert variant="destructive">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            {error?.stack?.message || "Failed to load run data."}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  const isToolRun = detailedRun.type === "TOOL_RUN";
  const isPipelineRun = detailedRun.type === "PIPELINE_RUN";

  return (
    <div className="container mx-auto p-4">
      <Tabs className="w-full" defaultValue="overview">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          {isPipelineRun && (
            <TabsTrigger value="pipeline">Pipeline</TabsTrigger>
          )}
          {isToolRun && <TabsTrigger value="tool">Tool</TabsTrigger>}
        </TabsList>
        <TabsContent value="overview">
          <RunInfo run={detailedRun} />
        </TabsContent>
        {isPipelineRun && (
          <TabsContent value="pipeline">
            <PipelineDetails
              childRuns={detailedRun.childRuns}
              pipeline={detailedRun.pipeline}
            />
          </TabsContent>
        )}
        {isToolRun && (
          <TabsContent value="tool">
            <ToolDetails tool={detailedRun.tool} />
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
};

interface RunInfoProps {
  run: any;
}

const RunInfo: React.FC<RunInfoProps> = ({ run }) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Field</TableHead>
          <TableHead>Value</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {Object.entries(run).map(([key, value]) => {
          if (
            key === "tool" ||
            key === "pipeline" ||
            key === "childRuns" ||
            key === "parentRun"
          ) {
            return null;
          }
          return (
            <TableRow key={key}>
              <TableCell className="font-medium capitalize">{key}</TableCell>
              <TableCell>{String(value)}</TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};

interface ToolDetailsProps {
  tool: any;
}

const ToolDetails: React.FC<ToolDetailsProps> = ({ tool }) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Tool Information</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Field</TableHead>
              <TableHead>Value</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {Object.entries(tool).map(([key, value]) => (
              <TableRow key={key}>
                <TableCell className="font-medium capitalize">{key}</TableCell>
                <TableCell>{String(value)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};

interface PipelineDetailsProps {
  childRuns: any[];
  pipeline: any;
}

const PipelineDetails: React.FC<PipelineDetailsProps> = ({
  childRuns,
  pipeline,
}) => {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Pipeline Information</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Field</TableHead>
                <TableHead>Value</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {Object.entries(pipeline).map(([key, value]) => {
                if (key === "pipelineTools") return null;
                return (
                  <TableRow key={key}>
                    <TableCell className="font-medium capitalize">
                      {key}
                    </TableCell>
                    <TableCell>{String(value)}</TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Pipeline Tools</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Tool Name</TableHead>
                <TableHead>Depends On</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {pipeline.pipelineTools.map((tool: any) => (
                <TableRow key={tool.id}>
                  <TableCell>{tool.id}</TableCell>
                  <TableCell>{tool.tool}</TableCell>
                  <TableCell>{tool.dependsOnId}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {childRuns && childRuns.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Child Runs</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Progress</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {childRuns.map((childRun) => (
                  <TableRow key={childRun.id}>
                    <TableCell>{childRun.id}</TableCell>
                    <TableCell>{childRun.name}</TableCell>
                    <TableCell>{childRun.status}</TableCell>
                    <TableCell>{childRun.progress}%</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default RunDetailsPage;
