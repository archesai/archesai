import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetRunSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/runs/$runID");

export const Route = createFileRoute("/_app/runs/$runID/")({
  component: RunDetailsPage,
});

function RunDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const runID = params.runID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <RunDetails runID={runID} />
        </Suspense>
      </Card>
    </div>
  );
}

function RunDetails({ runID }: { runID: string }): JSX.Element {
  const {
    data: { data: run },
  } = useGetRunSuspense(runID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Run Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(run.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(run.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Completed At
          </dt>
          <dd className="mt-1 text-sm">{String(run.completedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Error</dt>
          <dd className="mt-1 text-sm">{String(run.error)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(run.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Pipeline ID
          </dt>
          <dd className="mt-1 text-sm">{String(run.pipelineID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Progress
          </dt>
          <dd className="mt-1 text-sm">{String(run.progress)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Started At
          </dt>
          <dd className="mt-1 text-sm">{String(run.startedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Status</dt>
          <dd className="mt-1 text-sm">{String(run.status)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Tool ID</dt>
          <dd className="mt-1 text-sm">{String(run.toolID)}</dd>
        </div>
      </dl>
    </div>
  );
}
