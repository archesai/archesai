import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import TodoDataTable from "#components/datatables/todo-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/todos");

export const Route = createFileRoute("/_app/todos/")({
  component: TodosPage,
});

function TodosPage(): JSX.Element {
  return <TodoDataTable />;
}
