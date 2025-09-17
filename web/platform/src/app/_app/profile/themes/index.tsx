import { ThemeEditorPage } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/profile/themes/")({
  component: RouteComponent,
});

function RouteComponent() {
  return <ThemeEditorPage />;
}
