import { OpenAPIPage } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/openapi/")({
  component: () => <OpenAPIPage url="http://moose:3000/api/openapi.yaml" />,
});
