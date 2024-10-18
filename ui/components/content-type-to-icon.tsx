import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import { File, FilePlus, FileText, Image, Music, Video } from "lucide-react";

export const ContentTypeToIcon = ({ contentType }: { contentType: string }) => {
  const getLabel = (contentType: string) => {
    switch (contentType) {
      case "application/msword":
      case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
        return <FileText />;
      case "application/pdf":
        return <FileText />;
      case "application/vnd.ms-excel":
      case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
        return <FileText />;
      case "application/x-rar-compressed":
        return <FilePlus />;
      case "application/zip":
      case "audio/mp3":
        return <Music />;
      case "audio/mpeg":
      case "image/gif":
        return <Image />;
      case "image/jpeg":
      case "image/png":
      case "video/mp4":
      case "video/mpeg":
        return <Video />;
      default:
        return <File />;
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
