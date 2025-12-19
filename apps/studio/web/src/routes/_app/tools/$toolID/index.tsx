import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetToolSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/tools/$toolID");

export const Route = createFileRoute("/_app/tools/$toolID/")({
  component: ToolDetailsPage,
});

function ToolDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const toolID = params.toolID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <ToolDetails toolID={toolID} />
        </Suspense>
      </Card>
    </div>
  );
}

function ToolDetails({ toolID }: { toolID: string }): JSX.Element {
  const {
    data: { data: tool },
  } = useGetToolSuspense(toolID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Tool Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(tool.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(tool.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Description
          </dt>
          <dd className="mt-1 text-sm">{String(tool.description)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Input Mime Type
          </dt>
          <dd className="mt-1 text-sm">{String(tool.inputMimeType)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(tool.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(tool.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Output Mime Type
          </dt>
          <dd className="mt-1 text-sm">{String(tool.outputMimeType)}</dd>
        </div>
      </dl>
    </div>
  );
}
