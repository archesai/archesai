import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetPipelineSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/pipelines/$pipelineID");

export const Route = createFileRoute("/_app/pipelines/$pipelineID/")({
  component: PipelineDetailsPage,
});

function PipelineDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const pipelineID = params.pipelineID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <PipelineDetails pipelineID={pipelineID} />
        </Suspense>
      </Card>
    </div>
  );
}

function PipelineDetails({ pipelineID }: { pipelineID: string }): JSX.Element {
  const {
    data: { data: pipeline },
  } = useGetPipelineSuspense(pipelineID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Pipeline Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(pipeline.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(pipeline.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Description
          </dt>
          <dd className="mt-1 text-sm">{String(pipeline.description)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(pipeline.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(pipeline.organizationID)}</dd>
        </div>
      </dl>
    </div>
  );
}
