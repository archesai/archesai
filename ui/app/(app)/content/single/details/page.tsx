"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useContentControllerFindOne } from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
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
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{content.name}</CardTitle>
          <CardDescription>{content.description}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center space-x-2">
            <Badge variant="outline">{content.type}</Badge>
            <Badge variant="outline">{content.mimeType}</Badge>
            <Badge variant="outline">
              {format(new Date(content.createdAt), "PPP")}
            </Badge>
          </div>
          <Separator className="my-4" />

          {/* Additional Information */}
          <div className="mt-6 space-y-2">
            <h3 className="text-lg font-semibold">Details</h3>
            <p>
              <strong>Credits Used:</strong> {content.credits}
            </p>
          </div>

          {/* Annotations and Build Args */}
          <div className="mt-6 space-y-2">
            <h3 className="text-lg font-semibold">Build Arguments</h3>
            <pre className="bg-gray-100 p-4 rounded-md">
              {JSON.stringify(content.buildArgs, null, 2)}
            </pre>
          </div>

          {/* Text Content */}
          {content.text && (
            <div className="mt-6 space-y-2">
              <h3 className="text-lg font-semibold">Text</h3>
              <p>{content.text}</p>
            </div>
          )}

          {/* Download Button */}
          <div className="mt-6">
            <Button asChild>
              <a href={content.url} rel="noopener noreferrer" target="_blank">
                Download Content
              </a>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
