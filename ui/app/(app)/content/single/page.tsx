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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useContentControllerFindOne } from "@/generated/archesApiComponents";
import {
  TextChunksControllerFindAllPathParams,
  useTextChunksControllerFindAll,
} from "@/generated/archesApiComponents";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { TextChunkEntity } from "@/generated/archesApiSchemas";
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
    <div className="flex h-full w-full gap-3">
      {/*LEFT SIDE*/}
      <div className="flex w-1/2 flex-initial flex-col gap-3">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
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
            <div className="flex items-center gap-2">
              <Badge variant="secondary">{content.mimeType}</Badge>
              <Badge variant="secondary">
                {format(new Date(content.createdAt), "PPP")}
              </Badge>
              {content.jobs.map((job) => {
                return <Badge variant="secondary">{job.toolId}</Badge>;
              })}
            </div>
          </CardContent>
        </Card>

        <Tabs
          className="flex h-full flex-col gap-1"
          defaultValue={content.jobs[0].toolId}
        >
          <TabsList>
            {content.jobs.map((job) => {
              return <TabsTrigger value={job.toolId}>{job.toolId}</TabsTrigger>;
            })}
          </TabsList>
          <TabsContent className="grow" value="extract-text">
            <DataTable<
              { name: string } & TextChunkEntity,
              TextChunksControllerFindAllPathParams,
              undefined
            >
              columns={[
                {
                  accessorKey: "text",
                  cell: ({ row }) => {
                    return row.original.text;
                  },
                  header: ({ column }) => (
                    <DataTableColumnHeader column={column} title="Text" />
                  ),
                },
              ]}
              content={() => (
                <div className="flex h-full w-full items-center justify-center"></div>
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
              useFindAll={useTextChunksControllerFindAll}
              useRemove={() => ({
                mutateAsync: async () => {},
              })}
            />
          </TabsContent>
          <TabsContent value="image-to-text"></TabsContent>
          <TabsContent value="password">Change your password here.</TabsContent>
        </Tabs>
      </div>
      {/*RIGHT SIDE*/}
      <Card className="w-1/2 overflow-hidden">{renderContent(content)}</Card>
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
        className="h-full w-full object-contain"
        height={516}
        src={url}
        width={516}
      />
    );
  } else if (mimeType === "application/pdf") {
    return (
      <iframe className="h-full w-full" src={url} title="PDF Document"></iframe>
    );
  } else {
    return (
      <div className="flex h-full items-center justify-center">
        <p>Cannot preview this content type. Please download to view.</p>
      </div>
    );
  }
}
