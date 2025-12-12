import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetSessionSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/sessions/$sessionID");

export const Route = createFileRoute("/_app/sessions/$sessionID/")({
  component: SessionDetailsPage,
});

function SessionDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const sessionID = params.sessionID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <SessionDetails sessionID={sessionID} />
        </Suspense>
      </Card>
    </div>
  );
}

function SessionDetails({ sessionID }: { sessionID: string }): JSX.Element {
  const {
    data: { data: session },
  } = useGetSessionSuspense(sessionID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Session Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(session.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(session.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Auth Method
          </dt>
          <dd className="mt-1 text-sm">{String(session.authMethod)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Auth Provider
          </dt>
          <dd className="mt-1 text-sm">{String(session.authProvider)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Expires At
          </dt>
          <dd className="mt-1 text-sm">{String(session.expiresAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            IP Address
          </dt>
          <dd className="mt-1 text-sm">{String(session.ipAddress)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(session.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Token</dt>
          <dd className="mt-1 text-sm">{String(session.token)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            User Agent
          </dt>
          <dd className="mt-1 text-sm">{String(session.userAgent)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">User ID</dt>
          <dd className="mt-1 text-sm">{String(session.userID)}</dd>
        </div>
      </dl>
    </div>
  );
}
