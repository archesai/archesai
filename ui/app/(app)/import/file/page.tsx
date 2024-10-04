"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  useContentControllerCreate,
  useStorageControllerGetReadUrl,
  useStorageControllerGetWriteUrl,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import React, { useRef, useState } from "react";

export default function ImportPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState<boolean>(false);
  const [dragActive, setDragActive] = useState<boolean>(false);
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const { mutateAsync: getWriteUrl } = useStorageControllerGetWriteUrl();
  const { mutateAsync: getReadUrl } = useStorageControllerGetReadUrl();
  const { mutateAsync: indexDocument } = useContentControllerCreate();

  const handleFiles = (files: FileList | null) => {
    if (files) {
      setSelectedFiles(Array.from(files));
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

  const uploadFile = (file: File, writeUrl: string) => {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      xhr.upload.addEventListener("progress", (event) => {
        if (event.lengthComputable) {
          const percentCompleted = Math.round(
            (event.loaded * 100) / event.total
          );
          // Update progress for this file
          setUploadProgress(percentCompleted);
        }
      });

      xhr.onload = async () => {
        if (xhr.status === 200 || xhr.status === 201) {
          try {
            // File uploaded successfully
            console.log(`File ${file.name} uploaded successfully`);

            const readUrlResponse = await getReadUrl({
              body: {
                path: `uploads/${file.name}`, // Include the file name in the path
              },
              pathParams: {
                orgname: defaultOrgname,
              },
            });

            await indexDocument({
              body: {
                buildArgs: {},
                name: file.name,
                type: "DOCUMENT",
                url: readUrlResponse.read,
              },
              pathParams: {
                orgname: defaultOrgname,
              },
            });

            resolve(null);
          } catch (error) {
            console.error(`Error processing file ${file.name}:`, error);
            reject(error);
          }
        } else {
          // Handle server errors
          console.error(`Upload failed for ${file.name}:`, xhr.responseText);
          reject(new Error(`Upload failed with status ${xhr.status}`));
        }
      };

      xhr.onerror = () => {
        // Handle network errors
        console.error(`Network error while uploading file ${file.name}`);
        reject(new Error("Network error"));
      };

      xhr.open("PUT", writeUrl);
      xhr.setRequestHeader("Content-Type", file.type);
      xhr.send(file);
    });
  };

  const uploadFiles = async () => {
    if (selectedFiles.length === 0) return;
    setUploading(true);

    try {
      for (const file of selectedFiles) {
        // Get a unique write URL for each file
        const writeUrlResponse = await getWriteUrl({
          body: {
            path: `uploads/${file.name}`, // Include the file name in the path
          },
          pathParams: {
            orgname: defaultOrgname,
          },
        });

        const writeUrl = writeUrlResponse.write; // Adjust this based on your API response structure
        console.log("Write URL for", file.name, ":", writeUrl);

        // Upload the file using the write URL
        await uploadFile(file, writeUrl);
      }

      setUploading(false);
      setSelectedFiles([]);
      setUploadProgress(0);
      console.log("All files uploaded successfully");
      router.push("/content");
    } catch (error) {
      console.error("An error occurred during file upload:", error);
      alert("An error occurred during file upload");
      setUploading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto h-[80%] mt-8 stack items-center justify-center">
      <Card
        className={cn(
          "border-dashed border-2 p-6 text-center",
          dragActive ? "border-blue-500" : "border-gray-400"
        )}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
      >
        <CardContent>
          <p className="text-gray-600">
            Drag and drop files here, or click to select files
          </p>
          <Button
            className="mt-4"
            onClick={() => fileInputRef.current?.click()}
            variant="outline"
          >
            Select Files
          </Button>
          <input
            className="hidden"
            multiple
            onChange={handleFileInputChange}
            ref={fileInputRef}
            type="file"
          />
        </CardContent>
      </Card>

      {selectedFiles.length > 0 && (
        <div className="mt-6">
          <h3 className="text-lg font-semibold">Selected Files:</h3>
          <ul className="list-disc list-inside">
            {selectedFiles.map((file, idx) => (
              <li key={idx}>{file.name}</li>
            ))}
          </ul>
          <Button className="mt-4" disabled={uploading} onClick={uploadFiles}>
            Upload
          </Button>
        </div>
      )}

      {uploading && (
        <div className="mt-6">
          <Progress value={uploadProgress} />
          <p className="text-center mt-2">{uploadProgress}%</p>
        </div>
      )}
    </div>
  );
}
