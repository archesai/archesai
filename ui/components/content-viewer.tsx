import { ContentEntity } from "@/generated/archesApiSchemas";
import dynamic from "next/dynamic";
import Image from "next/image";

const ReactPlayer = dynamic(() => import("react-player"), { ssr: false });
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import Link from "next/link";

export function ContentViewer({
  content,
  size,
}: {
  content: ContentEntity;
  size: "lg" | "sm";
}) {
  const { mimeType, url } = content;

  let hoverContent = null;

  if (mimeType?.startsWith("video/") || mimeType?.startsWith("audio/")) {
    hoverContent = (
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
  } else if (mimeType?.startsWith("image/")) {
    hoverContent = (
      <Image
        alt={content.description || ""}
        className="h-full w-full object-contain"
        height={516}
        src={url || ""}
        width={516}
      />
    );
  } else if (mimeType === "application/pdf") {
    hoverContent = (
      <iframe className="h-full w-full" src={url} title="PDF Document"></iframe>
    );
  } else {
    hoverContent = (
      <div className="flex h-full items-center justify-center">
        <p>Cannot preview this content type. Please download to view.</p>
      </div>
    );
  }

  if (size === "lg") {
    return hoverContent;
  } else {
    return (
      <HoverCard openDelay={0}>
        <HoverCardTrigger asChild>
          <Link
            className="underline underline-offset-4"
            href={`/content/single?contentId=${content.id}`}
          >
            View Content
          </Link>
        </HoverCardTrigger>
        <HoverCardContent>{hoverContent}</HoverCardContent>
      </HoverCard>
    );
  }
}
