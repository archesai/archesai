import type {
  Artifact,
  ArtifactsFilterParameter,
  ArtifactsSortParameter,
  PageQueryParameter,
} from "@archesai/client";
import { getListArtifactsSuspenseQueryOptions } from "@archesai/client";
import {
  CalendarIcon,
  FileIcon,
  TextIcon,
} from "@archesai/ui/components/custom/icons";
import { Timestamp } from "@archesai/ui/components/custom/timestamp";
import { DataTable } from "@archesai/ui/components/datatable/data-table";
import { Badge } from "@archesai/ui/components/shadcn/badge";
import { ARTIFACT_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { Link, useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";

import ArtifactForm from "#components/forms/artifact-form";

export default function ArtifactDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListArtifactsSuspenseQueryOptions({
      filter: query.filter as unknown as ArtifactsFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as ArtifactsSortParameter,
    });
  };

  return (
    <DataTable<Artifact>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex gap-2">
                <Link
                  className="text-blue-500 hover:underline"
                  params={{
                    artifactID: row.original.id,
                  }}
                  to={`/artifacts/$artifactID`}
                >
                  {row.original.name}
                </Link>
              </div>
            );
          },
          enableColumnFilter: true,
          id: "name",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Name",
          },
        },
        {
          accessorKey: "mimeType",
          cell: ({ row }) => {
            return <Badge variant={"secondary"}>{row.original.mimeType}</Badge>;
          },
          enableColumnFilter: true,
          enableHiding: false,
          id: "mimeType",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Artifact Type",
            options: [
              { label: "Text", value: "text" },
              { label: "Image", value: "image" },
              { label: "Audio", value: "audio" },
              { label: "Video", value: "video" },
            ],
          },
        },
        {
          accessorKey: "producer",
          cell: ({ row }) => {
            return row.original.producerID ? (
              <Link
                className="text-blue-500 hover:underline"
                params={{
                  artifactID: row.original.id,
                }}
                search={{
                  selectedRunID: row.original.producerID,
                }}
                to={`/artifacts/$artifactID`}
              >
                {row.original.producerID}
              </Link>
            ) : (
              <div className="text-muted-foreground">None</div>
            );
          },
          enableColumnFilter: true,
          enableSorting: true,
          id: "producer",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Producer",
          },
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />;
          },
          enableColumnFilter: true,
          enableSorting: true,
          id: "createdAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Created",
          },
        },
      ]}
      createForm={ArtifactForm}
      entityKey={ARTIFACT_ENTITY_KEY}
      // biome-ignore lint/suspicious/noExplicitAny: FIXME
      getQueryOptions={getQueryOptions as any}
      grid={(_item) => {
        return (
          <div className="flex h-full w-full items-center justify-center">
            {/* <Image
              alt='source image'
              height={256}
              src={item.previewImage}
              width={256}
            /> */}
          </div>
        );
      }}
      handleSelect={async (artifact) => {
        await navigate({
          params: { artifactID: artifact.id },
          to: `/artifacts/$artifactID`,
        });
      }}
      icon={<FileIcon size={24} />}
      updateForm={ArtifactForm}
    />
  );
}
