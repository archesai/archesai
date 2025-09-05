import type { JSX } from "react";

import { createFileRoute } from "@tanstack/react-router";

import ArtifactDataTable from "#components/datatables/artifact-datatable";

export const Route = createFileRoute("/_app/artifacts/")({
  component: ArtifactsPage,
});

export default function ArtifactsPage(): JSX.Element {
  return <ArtifactDataTable />;
}
