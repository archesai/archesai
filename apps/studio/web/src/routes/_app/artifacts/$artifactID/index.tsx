import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetArtifactSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/artifacts/$artifactID");

export const Route = createFileRoute("/_app/artifacts/$artifactID/")({
  component: ArtifactDetailsPage,
});

function ArtifactDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const artifactID = params.artifactID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <ArtifactDetails artifactID={artifactID} />
        </Suspense>
      </Card>
    </div>
  );
}

function ArtifactDetails({ artifactID }: { artifactID: string }): JSX.Element {
  const {
    data: { data: artifact },
  } = useGetArtifactSuspense(artifactID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Artifact Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Credits</dt>
          <dd className="mt-1 text-sm">{String(artifact.credits)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Description
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.description)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Mime Type
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.mimeType)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(artifact.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Preview Image
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.previewImage)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Producer ID
          </dt>
          <dd className="mt-1 text-sm">{String(artifact.producerID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Text</dt>
          <dd className="mt-1 text-sm">{String(artifact.text)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">URL</dt>
          <dd className="mt-1 text-sm">{String(artifact.url)}</dd>
        </div>
      </dl>
    </div>
  );
}
