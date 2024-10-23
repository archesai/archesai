"use client";

import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useContentControllerFindOne } from "@/generated/archesApiComponents";
import {
  useVectorRecordControllerFindAll,
  VectorRecordControllerFindAllPathParams,
} from "@/generated/archesApiComponents";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { VectorRecordEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { format } from "date-fns";
import { File } from "lucide-react";
import dynamic from "next/dynamic";
import Image from "next/image";
import { useSearchParams } from "next/navigation";

const ReactPlayer = dynamic(() => import("react-player"), { ssr: false });

export default function ContentDetailsPage() {
  const searchParams = useSearchParams();
  const contentId = searchParams?.get("contentId");

  const { defaultOrgname } = useAuth();

  const { data: content, isLoading } = useContentControllerFindOne(
    {
      pathParams: {
        contentId: contentId as string,
        orgname: defaultOrgname,
      },
    },
    {
      enabled: !!defaultOrgname && !!contentId,
    }
  );

  if (isLoading || !content) {
    return <div>Loading...</div>;
  }

  return (
    <div className="flex w-full gap-3 flex-shrink-0 items-start">
      <div className="flex-initial w-1/2 gap-3 flex flex-col">
        <Card>
          <CardHeader>
            <CardTitle className="flex justify-between items-center">
              <div>{content.name}</div>
              <Button asChild size="sm" variant="secondary">
                <a href={content.url} rel="noopener noreferrer" target="_blank">
                  Download Content
                </a>
              </Button>
            </CardTitle>
            <CardDescription>
              {content.description || "No Description"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center space-x-2">
              <Badge variant="secondary">{content.mimeType}</Badge>
              <Badge variant="secondary">
                {format(new Date(content.createdAt), "PPP")}
              </Badge>
            </div>
          </CardContent>
        </Card>
        <DataTable<
          { name: string } & VectorRecordEntity,
          VectorRecordControllerFindAllPathParams,
          undefined
        >
          columns={[
            {
              accessorKey: "text",
              cell: ({ row }) => {
                return (
                  <div className="flex space-x-2">{row.original.text}</div>
                );
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
      </div>

      <div className="w-1/2 rounded-lg overflow-hidden shadow-lg self-stretch">
        {renderContent(content)}
      </div>
    </div>
  );
}

function renderContent(content: ContentEntity) {
  const { mimeType, url } = content;

  if (mimeType.startsWith("video/") || mimeType.startsWith("audio/")) {
    return (
      <ReactPlayer
        config={{
          file: {
            attributes: {
              controlsList: "nodownload",
            },
          },
        }}
        controls
        height="100%"
        url={url}
        width="100%"
      />
    );
  } else if (mimeType.startsWith("image/")) {
    return (
      <Image
        alt={content.description}
        className="w-full h-full object-contain"
        height={516}
        src={url}
        width={516}
      />
    );
  } else if (mimeType === "application/pdf") {
    return (
      <iframe className="w-full h-full" src={url} title="PDF Document"></iframe>
    );
  } else {
    return (
      <div className="flex items-center justify-center h-full">
        <p>Cannot preview this content type. Please download to view.</p>
      </div>
    );
  }
}
