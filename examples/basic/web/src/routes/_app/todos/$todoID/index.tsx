import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetTodoSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/todos/$todoID");

export const Route = createFileRoute("/_app/todos/$todoID/")({
  component: TodoDetailsPage,
});

function TodoDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const todoID = params.todoID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <TodoDetails todoID={todoID} />
        </Suspense>
      </Card>
    </div>
  );
}

function TodoDetails({ todoID }: { todoID: string }): JSX.Element {
  const {
    data: { data: todo },
  } = useGetTodoSuspense(todoID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Todo Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(todo.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(todo.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Completed
          </dt>
          <dd className="mt-1 text-sm">{String(todo.completed)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Title</dt>
          <dd className="mt-1 text-sm">{String(todo.title)}</dd>
        </div>
      </dl>
    </div>
  );
}
