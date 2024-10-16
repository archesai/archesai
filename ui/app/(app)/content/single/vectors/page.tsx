"use client";

import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import {
  useVectorRecordControllerFindAll,
  VectorRecordControllerFindAllPathParams,
} from "@/generated/archesApiComponents";
import { VectorRecordEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { File } from "lucide-react";
import { useSearchParams } from "next/navigation";

export default function ContentVectorsPage() {
  const searchParams = useSearchParams();
  const contentId = searchParams?.get("contentId");

  const { defaultOrgname } = useAuth();

  return (
    <DataTable<
      { name: string } & VectorRecordEntity,
      VectorRecordControllerFindAllPathParams,
      undefined
    >
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
      dataIcon={<File size={24} />}
      defaultView="table"
      findAllPathParams={{
        contentId: contentId as string,
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={() => {}}
      handleSelect={() => {}}
      itemType="vector"
      useFindAll={useVectorRecordControllerFindAll}
      useRemove={() => ({
        mutateAsync: async () => {},
      })}
    />
  );
}
