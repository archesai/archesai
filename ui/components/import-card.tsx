'use client'

import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import {
  useContentControllerCreate,
  useStorageControllerGetReadUrl,
  useStorageControllerGetWriteUrl
} from '@/generated/archesApiComponents'
import { ContentEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { useToast } from '@/hooks/use-toast'
import { cn } from '@/lib/utils'
import { CloudUpload, Loader2, Trash, Upload } from 'lucide-react'
import React, { useRef, useState } from 'react'

import { Badge } from './ui/badge'

export default function ImportCard({
  cb
}: {
  cb?: (content: ContentEntity[]) => void
}) {
  const { defaultOrgname } = useAuth()
  const [selectedFiles, setSelectedFiles] = useState<File[]>([])
  const [uploading, setUploading] = useState<boolean>(false)
  const [dragActive, setDragActive] = useState<boolean>(false)
  const [uploadProgress, setUploadProgress] = useState<number>(0)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const { toast } = useToast()

  const { mutateAsync: getWriteUrl } = useStorageControllerGetWriteUrl()
  const { mutateAsync: getReadUrl } = useStorageControllerGetReadUrl()
  const { mutateAsync: createContent } = useContentControllerCreate()

  const handleFiles = (files: FileList | null) => {
    if (files) {
      const newFiles = Array.from(files).filter(
        (file) =>
          !selectedFiles.some(
            (f) => f.name === file.name && f.size === file.size
          )
      )
      setSelectedFiles((prev) => [...prev, ...newFiles])
    }
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setDragActive(false)
    handleFiles(e.dataTransfer.files)
  }

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setDragActive(true)
  }

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setDragActive(false)
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFiles(e.target.files)
  }

  const removeFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index))
  }

  const uploadFile = (file: File, writeUrl: string): Promise<ContentEntity> => {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()

      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const percentCompleted = Math.round(
            (event.loaded * 100) / event.total
          )
          setUploadProgress((prev) => Math.max(prev, percentCompleted))
        }
      })

      xhr.onload = async () => {
        if (xhr.status === 200 || xhr.status === 201) {
          try {
            // File uploaded successfully
            console.log(`File ${file.name} uploaded successfully`)

            const readUrlResponse = await getReadUrl({
              body: {
                path: `uploads/${file.name}`
              },
              pathParams: {
                orgname: defaultOrgname
              }
            })

            const content = await createContent({
              body: {
                name: file.name,
                url: readUrlResponse.read
              },
              pathParams: {
                orgname: defaultOrgname
              }
            })

            resolve(content)
          } catch (error) {
            console.error(`Error processing file ${file.name}:`, error)
            reject(error)
          }
        } else {
          // Handle server errors
          console.error(`Upload failed for ${file.name}:`, xhr.responseText)
          reject(new Error(`Upload failed with status ${xhr.status}`))
        }
      }

      xhr.onerror = () => {
        // Handle network errors
        console.error(`Network error while uploading file ${file.name}`)
        reject(new Error('Network error'))
      }

      xhr.open('PUT', writeUrl)
      xhr.setRequestHeader('Content-Type', file.type)
      xhr.send(file)
    })
  }

  const uploadFiles = async () => {
    if (selectedFiles.length === 0) return
    setUploading(true)
    setUploadProgress(0)

    try {
      const urls = await Promise.all(
        selectedFiles.map(async (file) => {
          // Get a unique write URL for each file
          const writeUrlResponse = await getWriteUrl({
            body: {
              path: `uploads/${file.name}`
            },
            pathParams: {
              orgname: defaultOrgname
            }
          })

          const writeUrl = writeUrlResponse.write
          console.log('Write URL for', file.name, ':', writeUrl)

          // Upload the file using the write URL
          const readUrl = await uploadFile(file, writeUrl)

          return readUrl
        })
      )

      setUploading(false)
      setSelectedFiles([])
      setUploadProgress(100)
      toast({
        description: 'Files uploaded successfully.',
        title: 'Upload Complete'
      })
      if (cb) {
        cb(urls)
      }
    } catch (error) {
      console.error('An error occurred during file upload:', error)
      toast({
        description:
          (error as any).message || 'An error occurred while uploading files.',
        title: 'Upload Failed',
        variant: 'destructive'
      })
      setUploading(false)
    }
  }

  return (
    <div
      className={
        'flex flex-col items-center gap-2 rounded-lg transition-all duration-300'
      }
    >
      {/* Drop Area */}
      <Card
        className={cn(
          'flex w-full cursor-pointer flex-col items-center justify-center gap-2 border border-dashed p-8 transition-colors duration-300',
          dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-400'
        )}
        onClick={() => fileInputRef.current?.click()}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
      >
        <CloudUpload className='h-5 w-5 text-muted-foreground' />
        <p className='text-sm text-muted-foreground'>
          Drag and drop files here, or click to select files
        </p>

        <input
          className='hidden'
          multiple
          onChange={handleFileInputChange}
          ref={fileInputRef}
          type='file'
        />
      </Card>

      {/* Sidebar */}
      {selectedFiles.length > 0 && (
        <div className='flex w-full flex-col gap-2'>
          <ul className='flex max-h-52 grow flex-col gap-2 overflow-y-scroll'>
            {selectedFiles.map((file, idx) => (
              <li
                className='flex items-center justify-between rounded border bg-muted/50 p-2'
                key={idx}
              >
                <span className='flex w-4/5 items-center gap-2 truncate text-foreground'>
                  <span>{file.name}</span>
                  <Badge>{file.type}</Badge>
                </span>
                <Badge
                  className='text-nowrap text-primary'
                  variant='outline'
                >
                  {`${(file.size / 1024).toFixed(2)} KB`}
                </Badge>
                <button
                  aria-label={`Remove ${file.name}`}
                  className='text-red-500 hover:text-red-700 focus:outline-none'
                  onClick={() => removeFile(idx)}
                >
                  <Trash className='h-5 w-5' />
                </button>
              </li>
            ))}
          </ul>
          <Button
            className='flex items-center justify-center border'
            disabled={uploading || selectedFiles.length === 0}
            onClick={uploadFiles}
            size='sm'
            variant={'secondary'}
          >
            {uploading ? (
              <div className='flex gap-2'>
                <Loader2 className='h-5 w-5 animate-spin text-white' />
                <span>Uploading...</span>
              </div>
            ) : (
              <div className='flex gap-2'>
                <Upload className='h-5 w-5' />
                <span>Upload</span>
              </div>
            )}
          </Button>

          {uploading && (
            <div className='flex flex-col gap-2'>
              <Progress value={uploadProgress} />
              <p className='text-center text-sm text-gray-600'>
                {uploadProgress}% Uploaded
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
