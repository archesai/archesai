import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetExecutorSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/executors/$executorID");

export const Route = createFileRoute("/_app/executors/$executorID/")({
  component: ExecutorDetailsPage,
});

function ExecutorDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const executorID = params.executorID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <ExecutorDetails executorID={executorID} />
        </Suspense>
      </Card>
    </div>
  );
}

function ExecutorDetails({ executorID }: { executorID: string }): JSX.Element {
  const {
    data: { data: executor },
  } = useGetExecutorSuspense(executorID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Executor Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(executor.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(executor.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            CPU Shares
          </dt>
          <dd className="mt-1 text-sm">{String(executor.cpuShares)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Dependencies
          </dt>
          <dd className="mt-1 text-sm">{String(executor.dependencies)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Description
          </dt>
          <dd className="mt-1 text-sm">{String(executor.description)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Env</dt>
          <dd className="mt-1 text-sm">{String(executor.env)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Execute Code
          </dt>
          <dd className="mt-1 text-sm">{String(executor.executeCode)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Extra Files
          </dt>
          <dd className="mt-1 text-sm">{String(executor.extraFiles)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Is Active
          </dt>
          <dd className="mt-1 text-sm">{String(executor.isActive)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Language
          </dt>
          <dd className="mt-1 text-sm">{String(executor.language)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Memory MB
          </dt>
          <dd className="mt-1 text-sm">{String(executor.memoryMB)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(executor.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(executor.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Schema In
          </dt>
          <dd className="mt-1 text-sm">{String(executor.schemaIn)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Schema Out
          </dt>
          <dd className="mt-1 text-sm">{String(executor.schemaOut)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Timeout</dt>
          <dd className="mt-1 text-sm">{String(executor.timeout)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Version</dt>
          <dd className="mt-1 text-sm">{String(executor.version)}</dd>
        </div>
      </dl>
    </div>
  );
}
