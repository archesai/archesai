"use client";

import { useContentControllerFindOne } from "@/generated/archesApiComponents";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import dynamic from "next/dynamic";
import Image from "next/image";
import { useSearchParams } from "next/navigation";

// Dynamically import react-player to prevent SSR issues
const ReactPlayer = dynamic(() => import("react-player"), { ssr: false });

export default function ContentViewPage() {
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

  return <div className="w-full h-full">{renderContent(content)}</div>;
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
