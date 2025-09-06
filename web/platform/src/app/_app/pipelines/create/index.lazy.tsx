import { createLazyFileRoute } from "@tanstack/react-router";

import { CreatePipelineContent } from "#components/create-pipeline";

// import RunForm from "#components/forms/run-form";

export const Route = createLazyFileRoute("/_app/pipelines/create/")({
  component: CreatePipelineContent, // default export
});
