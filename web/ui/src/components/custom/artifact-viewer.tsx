import type { JSX } from "react"

interface ArtifactViewerProps {
  artifact: {
    mimeType: string
    text?: null | string
  }
}

export function ArtifactViewer({
  artifact
}: ArtifactViewerProps): JSX.Element | null {
  let hoverContent: React.ReactNode = null
  if (
    artifact.mimeType.startsWith("video/") ||
    artifact.mimeType.startsWith("audio/")
  ) {
    hoverContent = (
      <video
        className="h-full w-full object-contain"
        controls
        src={artifact.text ?? ""}
      />
    )
  } else if (artifact.mimeType.startsWith("image/") && artifact.text) {
    hoverContent = (
      <image
        // className='h-full w-full object-contain'
        height={516}
        href={artifact.text}
        width={516}
      />
    )
  } else if (artifact.mimeType === "application/pdf" && artifact.text) {
    hoverContent = (
      <iframe
        className="h-full w-full"
        src={artifact.text}
        title="PDF Document"
      ></iframe>
    )
  } else if (artifact.mimeType.startsWith("text/") && artifact.text) {
    hoverContent = (
      <div className="flex h-full items-center justify-center p-4 text-center">
        <p>{artifact.text}</p>
      </div>
    )
  } else {
    hoverContent = (
      <div className="flex h-full items-center justify-center">
        <p>Cannot preview this content type. Please download to view.</p>
      </div>
    )
  }

  return hoverContent
}
