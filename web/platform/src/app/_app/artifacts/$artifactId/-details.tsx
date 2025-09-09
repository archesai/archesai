import { useGetOneArtifactSuspense } from "@archesai/client";
import { Timestamp } from "@archesai/ui/components/custom/timestamp";
import { Badge } from "@archesai/ui/components/shadcn/badge";
import { Button } from "@archesai/ui/components/shadcn/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@archesai/ui/components/shadcn/card";
import type { JSX } from "react";

export const ArtifactDetailsHeader = ({
  artifactId,
}: {
  artifactId: string;
}): JSX.Element => {
  const {
    data: { data: artifact },
  } = useGetOneArtifactSuspense(artifactId);

  return (
    <CardHeader>
      <CardTitle className="flex items-center justify-between">
        <div>{artifact.name}</div>
        <Button
          asChild
          size="sm"
          variant="outline"
        >
          <a
            href={artifact.text} // FIXME - not a link
            rel="noopener noreferrer"
            target="_blank"
          >
            Download Artifact
          </a>
        </Button>
      </CardTitle>
      <CardDescription>{artifact.description}</CardDescription>
    </CardHeader>
  );
};

export const ArtifactDetailsBody = ({
  artifactId,
}: {
  artifactId: string;
}): JSX.Element => {
  const {
    data: { data: artifact },
  } = useGetOneArtifactSuspense(artifactId);

  return (
    <CardContent>
      <div className="flex items-center gap-2">
        <Badge>{artifact.mimeType}</Badge>
        {artifact.createdAt && (
          <Badge>
            <Timestamp date={artifact.createdAt} />
          </Badge>
        )}
      </div>
    </CardContent>
  );
};
