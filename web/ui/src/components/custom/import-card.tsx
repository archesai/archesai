import type { JSX } from "react";

import { useRef, useState } from "react";

import {
  Loader2Icon,
  TrashIcon,
  UploadCloudIcon,
} from "#components/custom/icons";
// import { toast } from 'sonner'

// import type { ArtifactEntity } from '@archesai/schemas'

import { Badge } from "#components/shadcn/badge";
import { Button } from "#components/shadcn/button";
import { Card } from "#components/shadcn/card";
import { Progress } from "#components/shadcn/progress";
import { cn } from "#lib/utils";

export default function ImportCard(): JSX.Element {
  //   {
  //   cb
  // }: {
  //   cb?: (content: ArtifactEntity[]) => void

  // }
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState<boolean>(false);
  const [dragActive, setDragActive] = useState<boolean>(false);
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFiles = (files: FileList | null) => {
    if (files) {
      const newFiles = Array.from(files).filter(
        (file) =>
          !selectedFiles.some(
            (f) => f.name === file.name && f.size === file.size,
          ),
      );
      setSelectedFiles((prev) => [...prev, ...newFiles]);
    }
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setDragActive(false);
    handleFiles(e.dataTransfer.files);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setDragActive(true);
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setDragActive(false);
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFiles(e.target.files);
  };

  const removeFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  // const uploadFile = (
  //   file: File,
  //   writeUrl: string
  // ): Promise<ArtifactEntity> => {
  //   return new Promise((resolve, reject) => {
  //     // Create a XMLHttpRequest
  //     const xhr = new XMLHttpRequest()

  //     // Add progress event listener
  //     xhr.upload.addEventListener('progress', (event) => {
  //       if (event.lengthComputable) {
  //         const percentCompleted = Math.round(
  //           (event.loaded * 100) / event.total
  //         )
  //         setUploadProgress((prev) => Math.max(prev, percentCompleted))
  //       }
  //     })

  //     // Add onload event listener
  //     xhr.onload = async () => {
  //       if (xhr.status === 200 || xhr.status === 201) {
  //         try {
  //           // const readUrlResponse = await createFile({
  //           //   action: 'read',
  //           //   path: `uploads/${file.name}`
  //           // })
  //           resolve('' as unknown as ArtifactEntity) // FIXME
  //         } catch (error) {
  //           console.error(error)
  //           reject(new Error(`Failed to create content`))
  //         }
  //       } else {
  //         reject(new Error(`Upload failed: ${xhr.responseText}`))
  //       }
  //     }

  //     // Add onerror event listener
  //     xhr.onerror = () => {
  //       reject(new Error('Network error'))
  //     }

  //     // Open the request and send the file
  //     xhr.open('PUT', writeUrl)
  //     xhr.setRequestHeader('Content-Type', file.type)
  //     xhr.send(file)
  //   })
  // }

  const uploadFiles = () => {
    if (selectedFiles.length === 0) return;
    setUploading(true);
    setUploadProgress(0);

    // try {
    //   const urls = await Promise.all(
    //     selectedFiles.map(async (file) => {
    //       const response = await createFile({
    //         action: 'write',
    //         path: `uploads/${file.name}`
    //       })
    //       const writeUrl = response.data.write
    //       // if (response.status !== 201) {
    //       //   throw new Error('Failed to create file.')
    //       // }
    //       if (!writeUrl) {
    //         throw new Error('Failed to get write URL for file upload.')
    //       }

    //       return uploadFile(file, writeUrl)
    //     })
    //   )
    //   setUploading(false)
    //   setSelectedFiles([])
    //   setUploadProgress(100)
    //   toast('Upload Complete', {})
    //   if (cb) {
    //     cb(urls)
    //   }
    // } catch (error) {
    //   console.error(error)
    //   toast('Upload Failed', {
    //     description: 'An error occurred while uploading files.'
    //   })
    //   setUploading(false)
    // }
  };

  return (
    <div
      className={
        "flex flex-col items-center gap-2 rounded-lg transition-all duration-300"
      }
    >
      {/* Drop Area */}
      <Card
        className={cn(
          "flex w-full cursor-pointer flex-col items-center justify-center gap-2 border border-dashed p-8 transition-colors duration-300",
          dragActive ? "border-blue-500 bg-blue-50" : "border-gray-400",
        )}
        onClick={() => fileInputRef.current?.click()}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
      >
        <UploadCloudIcon className="h-5 w-5 text-muted-foreground" />
        <p className="text-muted-foreground text-sm">
          Drag and drop files here, or click to select files
        </p>

        <input
          className="hidden"
          multiple
          onChange={handleFileInputChange}
          ref={fileInputRef}
          type="file"
        />
      </Card>

      {/* Sidebar */}
      {selectedFiles.length > 0 && (
        <div className="flex w-full flex-col gap-2">
          <ul className="flex max-h-52 grow flex-col gap-2 overflow-y-scroll">
            {selectedFiles.map((file, idx) => (
              <li
                className="flex items-center justify-between rounded-xs border bg-muted/50 p-2"
                key={file.name}
              >
                <span className="flex w-4/5 items-center gap-2 truncate text-foreground">
                  <span>{file.name}</span>
                  <Badge>{file.type}</Badge>
                </span>
                <Badge
                  className="text-nowrap text-primary"
                  variant="outline"
                >
                  {`${(file.size / 1024).toFixed(2)} KB`}
                </Badge>
                <button
                  aria-label={`Remove ${file.name}`}
                  className="text-red-500 hover:text-red-700 focus:outline-hidden"
                  onClick={() => {
                    removeFile(idx);
                  }}
                  type="button"
                >
                  <TrashIcon className="h-5 w-5" />
                </button>
              </li>
            ))}
          </ul>
          <Button
            className="flex items-center justify-center border"
            disabled={uploading || selectedFiles.length === 0}
            onClick={uploadFiles}
            size="sm"
            variant={"secondary"}
          >
            {uploading ? (
              <div className="flex gap-2">
                <Loader2Icon className="h-5 w-5 animate-spin text-white" />
                <span>Uploading...</span>
              </div>
            ) : (
              <div className="flex gap-2">
                <UploadCloudIcon className="h-5 w-5" />
                <span>Upload</span>
              </div>
            )}
          </Button>

          {uploading && (
            <div className="flex flex-col gap-2">
              <Progress value={uploadProgress} />
              <p className="text-center text-gray-600 text-sm">
                {uploadProgress}% Uploaded
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
