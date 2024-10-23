import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import { FileText, Image, Music, Video } from "lucide-react";

export const ContentTypeToIcon = ({ contentType }: { contentType: string }) => {
  const getLabel = (contentType: string) => {
    switch (contentType) {
      case "application/msword":
      case "application/pdf":
      case "application/vnd.ms-excel":
      case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
      case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
        return <FileText className="h-5 w-5" />;
      case "audio/mp3":
      case "audio/mpeg":
        return <Music className="h-5 w-5" />;
      case "image/gif":
      case "image/jpeg":
      case "image/png":
      case "image/svg+xml":
        return <Image className="h-5 w-5" />;

      case "video/mp4":
      case "video/mpeg":
        return <Video className="h-5 w-5" />;
      default:
        return <FileText className="h-5 w-5" />;
    }
  };

  return (
    <HoverCard>
      <HoverCardTrigger asChild>{getLabel(contentType)}</HoverCardTrigger>
      <HoverCardContent>
        <div className="flex items-center space-x-4">
          {getLabel(contentType)}
          <div className="space-y-1">
            <p className="text-sm">{contentType}</p>
          </div>
        </div>
      </HoverCardContent>
    </HoverCard>
  );
};
