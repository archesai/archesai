"use client";
import { ContentTypeToIcon } from "@/components/content-type-to-icon";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import {
  ContentControllerRemoveVariables,
  useContentControllerFindAll,
  useContentControllerRemove,
} from "@/generated/archesApiComponents";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { File } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ContentPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  const { limit, page, query, range } = useFilterItems();

  const {
    data: content,
    isLoading,
    isPlaceholderData,
  } = useContentControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
    queryParams: {
      limit,
      offset: page * limit,
      searchTerm: query,
      sortBy: "createdAt",
      sortDirection: "asc" as const,
      startDate: range.from?.toISOString(),
    },
  });

  const loading = isPlaceholderData || isLoading;
  const { mutateAsync: deleteDocument } = useContentControllerRemove();

  const { selectedItems } = useSelectItems({ items: content?.results || [] });

  return (
    <DataTable<ContentEntity, ContentControllerRemoveVariables>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <div className="flex space-x-2 justify-center text-muted-foreground max-w-10">
                  <ContentTypeToIcon contentType={row.original.mimeType} />
                </div>
                <Link
                  className="max-w-[200px] truncate font-medium text-primary"
                  href={`/content/single/general?contentId=${row.original.id}`}
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
              <div className="flex space-x-2">
                <span className="font-light">
                  {row.original.description || "n/a"}
                </span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Description" />
          ),
        },
        {
          accessorKey: "status",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="font-light">{row.original.job.status}</span>
                {row.original.job.status === "PROCESSING" && (
                  <span className="text-priamry">
                    {(row.original.job.progress * 100).toFixed(0)}%
                  </span>
                )}
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Status" />
          ),
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="font-medium">
                  {new Date(row.original.createdAt).toLocaleDateString()}
                </span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Created" />
          ),
        },
      ]}
      content={(item) => (
        <div className="flex w-full justify-center items-center h-full">
          <Image
            alt="source image"
            height={256}
            src={item.previewImage}
            width={256}
          />
        </div>
      )}
      data={content as any}
      dataIcon={<File size={24} />}
      defaultView="table"
      deleteItem={deleteDocument}
      getDeleteVariablesFromItem={(content) => [
        {
          pathParams: {
            contentId: content.id,
            orgname: defaultOrgname,
          },
        },
      ]}
      handleSelect={(content) =>
        router.push(`/content/single/general?contentId=${content.id}`)
      }
      itemType="content"
      loading={loading}
      mutationVariables={selectedItems.map((id) => ({
        pathParams: {
          contentId: id,
          orgname: defaultOrgname,
        },
      }))}
    />
  );
}
