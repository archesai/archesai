import { useGetArtifactSuspense } from "@archesai/client";
import {
  Badge,
  Button,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Timestamp,
} from "@archesai/ui";
import type { JSX } from "react";

export const ArtifactDetailsHeader = ({
  artifactID,
}: {
  artifactID: string;
}): JSX.Element => {
  const {
    data: { data: artifact },
  } = useGetArtifactSuspense(artifactID);

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
  artifactID,
}: {
  artifactID: string;
}): JSX.Element => {
  const {
    data: { data: artifact },
  } = useGetArtifactSuspense(artifactID);

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
