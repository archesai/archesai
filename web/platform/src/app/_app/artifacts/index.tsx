import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import ArtifactDataTable from "#components/datatables/artifact-datatable";

export const Route = createFileRoute("/_app/artifacts/")({
  component: ArtifactsPage,
});

function ArtifactsPage(): JSX.Element {
  return <ArtifactDataTable />;
}
