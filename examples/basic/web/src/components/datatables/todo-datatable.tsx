import {
  Badge,
  CalendarIcon,
  CheckIcon,
  DataTableContainer,
  ListIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import TodoForm from "#components/forms/todo-form";
import type { PageQueryParameter, Todo } from "#lib/index";
import { deleteTodo, getListTodosSuspenseQueryOptions } from "#lib/index";

export default function TodoDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListTodosSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Todo>
      columns={[
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            const val = row.original.createdAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "createdAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Created At",
          },
        },
        {
          accessorKey: "updatedAt",
          cell: ({ row }) => {
            const val = row.original.updatedAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "updatedAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Updated At",
          },
        },
        {
          accessorKey: "completed",
          cell: ({ row }) => {
            return (
              <Badge variant={row.original.completed ? "default" : "secondary"}>
                {row.original.completed ? "Yes" : "No"}
              </Badge>
            );
          },
          id: "completed",
          meta: {
            filterVariant: "boolean",
            icon: CheckIcon,
            label: "Completed",
          },
        },
        {
          accessorKey: "title",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.title}</Badge>;
          },
          id: "title",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Title",
          },
        },
      ]}
      createForm={TodoForm}
      deleteItem={async (id) => {
        await deleteTodo(id);
      }}
      entityKey="todos"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (todo) => {
        await navigate({
          params: {
            todoID: todo.id,
          },
          to: `/todos/$todoID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={TodoForm}
    />
  );
}
