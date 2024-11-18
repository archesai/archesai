"use client";

import { ContentViewer } from "@/components/content-viewer";
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
import { useAuth } from "@/hooks/use-auth";
import { format } from "date-fns";
import { useSearchParams } from "next/navigation";

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
              <Badge>{content.mimeType}</Badge>
              <Badge>{format(new Date(content.createdAt), "PPP")}</Badge>
            </div>
          </CardContent>
        </Card>
      </div>
      {/*RIGHT SIDE*/}
      <Card className="w-1/2 overflow-hidden">
        {<ContentViewer content={content} size="lg" />}
      </Card>
    </div>
  );
}
