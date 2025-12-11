import {
  Badge,
  CalendarIcon,
  DataTableContainer,
  ListIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import ArtifactForm from "#components/forms/artifact-form";
import type { Artifact, PageQueryParameter } from "#lib/index";
import {
  deleteArtifact,
  getListArtifactsSuspenseQueryOptions,
} from "#lib/index";

export default function ArtifactDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListArtifactsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Artifact>
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
          accessorKey: "credits",
          cell: ({ row }) => {
            return row.original.credits;
          },
          id: "credits",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Credits",
          },
        },
        {
          accessorKey: "description",
          cell: ({ row }) => {
            return row.original.description;
          },
          id: "description",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Description",
          },
        },
        {
          accessorKey: "mimeType",
          cell: ({ row }) => {
            return row.original.mimeType;
          },
          id: "mimeType",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Mime Type",
          },
        },
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.name}</Badge>;
          },
          id: "name",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Name",
          },
        },
        {
          accessorKey: "organizationID",
          cell: ({ row }) => {
            return row.original.organizationID;
          },
          id: "organizationID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Organization ID",
          },
        },
        {
          accessorKey: "previewImage",
          cell: ({ row }) => {
            return row.original.previewImage;
          },
          id: "previewImage",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Preview Image",
          },
        },
        {
          accessorKey: "producerID",
          cell: ({ row }) => {
            return row.original.producerID;
          },
          id: "producerID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Producer ID",
          },
        },
        {
          accessorKey: "text",
          cell: ({ row }) => {
            return row.original.text;
          },
          id: "text",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Text",
          },
        },
        {
          accessorKey: "url",
          cell: ({ row }) => {
            return row.original.url;
          },
          id: "url",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "URL",
          },
        },
      ]}
      createForm={ArtifactForm}
      deleteItem={async (id) => {
        await deleteArtifact(id);
      }}
      entityKey="artifacts"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (artifact) => {
        await navigate({
          params: {
            artifactID: artifact.id,
          },
          to: `/artifacts/$artifactID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={ArtifactForm}
    />
  );
}
