"use client";

import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { useVectorRecordControllerFindAll } from "@/generated/archesApiComponents";
import { VectorRecordEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { useSelectItems } from "@/hooks/useSelectItems";
import { File } from "lucide-react";
import { useSearchParams } from "next/navigation";

export default function ContentVectorsPage() {
  const searchParams = useSearchParams();
  const contentId = searchParams?.get("contentId");

  const { defaultOrgname } = useAuth();
  const { data: vectorRecords, isLoading: vectorRecordsIsLoading } =
    useVectorRecordControllerFindAll(
      {
        pathParams: {
          contentId: contentId as string,
          orgname: defaultOrgname,
        },
      },
      {
        enabled: !!defaultOrgname,
      }
    );

  const { selectedItems } = useSelectItems({
    items: vectorRecords?.results || [],
  });

  if (!vectorRecords) {
    return <div>Loading...</div>;
  }

  return (
    <DataTable<{ name: string } & VectorRecordEntity, undefined>
      columns={[
        {
          accessorKey: "text",
          cell: ({ row }) => {
            return <div className="flex space-x-2">{row.original.text}</div>;
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Text" />
          ),
        },
      ]}
      content={() => (
        <div className="flex w-full justify-center items-center h-full"></div>
      )}
      data={vectorRecords as any}
      dataIcon={<File size={24} />}
      defaultView="table"
      deleteItem={async () => {}}
      getDeleteVariablesFromItem={() => []}
      handleSelect={() => {}}
      itemType="vector"
      loading={vectorRecordsIsLoading}
      mutationVariables={selectedItems.map(() => undefined)}
    />
  );
}
