"use client";
import { ContentTypeToIcon } from "@/components/content-type-to-icon";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import ImageForm from "@/components/forms/image-form";
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
                <ContentTypeToIcon contentType={row.original.mimeType} />
                <Link
                  className="max-w-[200px] truncate font-medium text-primary"
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
              <span className="font-light">
                {row.original.description || "No Description"}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Description" />
          ),
        },
        {
          accessorKey: "status",
          cell: ({ row }) => {
            return <JobStatusButton job={row.original.job} />;
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Status" />
          ),
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return (
              <span className="font-medium">
                {new Date(row.original.createdAt).toLocaleDateString()}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Created" />
          ),
        },
      ]}
      content={(item) => (
        <div className="flex h-full w-full items-center justify-center">
          {item.job.status !== "COMPLETE" ? (
            <JobStatusButton job={item.job} />
          ) : (
            <Image
              alt="source image"
              height={256}
              src={item.previewImage}
              width={256}
            />
          )}
        </div>
      )}
      createForm={<ImageForm />}
      dataIcon={<File size={24} />}
      defaultView="grid"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      findAllQueryParams={{
        type: "IMAGE",
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
      itemType="image"
      useFindAll={useContentControllerFindAll}
      useRemove={useContentControllerRemove}
    />
  );
}
