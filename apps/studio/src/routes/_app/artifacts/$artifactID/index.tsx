import { ArtifactViewer, Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetArtifactSuspense } from "#lib/index";

import {
  ArtifactDetailsBody,
  ArtifactDetailsHeader,
} from "#routes/_app/artifacts/$artifactID/-details";

export const Route = createFileRoute("/_app/artifacts/$artifactID/")({
  component: ArtifactDetailsPage,
});

function ArtifactDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const artifactID = params.artifactID;

  return (
    <div className="flex h-full w-full gap-4">
      {/*LEFT SIDE*/}
      <Card>
        <Suspense>
          <ArtifactDetailsHeader artifactID={artifactID} />
        </Suspense>
        <Suspense>
          <ArtifactDetailsBody artifactID={artifactID} />
        </Suspense>
      </Card>

      {/*RIGHT SIDE*/}
      <Card className="w-1/2 overflow-hidden">
        <Suspense>
          <ArtifactViewerWrapper artifactID={artifactID} />
        </Suspense>
      </Card>
    </div>
  );
}

function ArtifactViewerWrapper({
  artifactID,
}: {
  artifactID: string;
}): JSX.Element {
  const {
    data: { data: artifact },
  } = useGetArtifactSuspense(artifactID);

  return <ArtifactViewer artifact={artifact} />;
}
