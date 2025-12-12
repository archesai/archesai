import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetAPIKeySuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/api-keys/$apiKeyID");

export const Route = createFileRoute("/_app/api-keys/$apiKeyID/")({
  component: APIKeyDetailsPage,
});

function APIKeyDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const apiKeyID = params.apiKeyID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <APIKeyDetails apiKeyID={apiKeyID} />
        </Suspense>
      </Card>
    </div>
  );
}

function APIKeyDetails({ apiKeyID }: { apiKeyID: string }): JSX.Element {
  const {
    data: { data: apiKey },
  } = useGetAPIKeySuspense(apiKeyID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">APIKey Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Expires At
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.expiresAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Key Hash
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.keyHash)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Last Used At
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.lastUsedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(apiKey.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Prefix</dt>
          <dd className="mt-1 text-sm">{String(apiKey.prefix)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Rate Limit
          </dt>
          <dd className="mt-1 text-sm">{String(apiKey.rateLimit)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Scopes</dt>
          <dd className="mt-1 text-sm">{String(apiKey.scopes)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">User ID</dt>
          <dd className="mt-1 text-sm">{String(apiKey.userID)}</dd>
        </div>
      </dl>
    </div>
  );
}
