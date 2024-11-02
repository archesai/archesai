"use client";
// import { ContentTypeToIcon } from "@/components/content-type-to-icon";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import ImportCard from "@/components/import-card";
import { JobStatusButton } from "@/components/job-status-button";
import {
  ContentControllerFindAllPathParams,
  ContentControllerRemoveVariables,
  useContentControllerFindAll,
  useContentControllerRemove,
} from "@/generated/archesApiComponents";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { File } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ContentPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  return (
    <div className="flex h-full flex-col gap-3">
      <ImportCard />
      <DataTable<
        ContentEntity,
        ContentControllerFindAllPathParams,
        ContentControllerRemoveVariables
      >
        columns={[
          {
            accessorKey: "name",
            cell: ({ row }) => {
              return (
                <div className="flex gap-2">
                  {/* <ContentTypeToIcon contentType={row.original.mimeType} /> */}
                  <Link
                    className="max-w-[200px] shrink truncate text-base font-medium text-primary md:text-sm"
                    href={`/content/single?contentId=${row.original.id}`}
                  >
                    {row.original.name}
                  </Link>
                </div>
              );
            },
            header: ({ column }) => (
              <DataTableColumnHeader column={column} title="Name" />
            ),
          },
          {
            accessorKey: "description",
            cell: ({ row }) => {
              return (
                <span className="text-base md:text-sm">
                  {row.original.description || "No Description"}
                </span>
              );
            },
            enableHiding: false,
            enableSorting: false,
            header: ({ column }) => (
              <DataTableColumnHeader
                className="-ml-2 text-sm"
                column={column}
                title="Description"
              />
            ),
          },

          {
            accessorKey: "createdAt",
            cell: ({ row }) => {
              return (
                <span className="font-light">
                  {new Date(row.original.createdAt).toLocaleDateString()}
                </span>
              );
            },
            header: ({ column }) => (
              <DataTableColumnHeader column={column} title="Created" />
            ),
          },
          {
            accessorKey: "tools",
            cell: ({ row }) => {
              return (
                <>
                  {row.original.jobs.map((job) => {
                    return <JobStatusButton job={job} />;
                  })}
                </>
              );
            },
            header: ({ column }) => (
              <DataTableColumnHeader column={column} title="Tools" />
            ),
          },
        ]}
        content={(item) => {
          const pendingJobs = item.jobs.filter(
            (job) => job.status !== "COMPLETE"
          );
          return (
            <div className="flex h-full w-full items-center justify-center">
              {pendingJobs.length && pendingJobs[0].status !== "COMPLETE" ? (
                <JobStatusButton job={pendingJobs[0]} />
              ) : (
                <Image
                  alt="source image"
                  height={256}
                  src={item.previewImage}
                  width={256}
                />
              )}
            </div>
          );
        }}
        dataIcon={<File size={24} />}
        defaultView="table"
        findAllPathParams={{
          orgname: defaultOrgname,
        }}
        getDeleteVariablesFromItem={(content) => ({
          pathParams: {
            contentId: content.id,
            orgname: defaultOrgname,
          },
        })}
        handleSelect={(content) =>
          router.push(`/content/single?contentId=${content.id}`)
        }
        itemType="content"
        useFindAll={useContentControllerFindAll}
        useRemove={useContentControllerRemove}
      />
    </div>
  );
}
