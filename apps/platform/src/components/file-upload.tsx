import type { DragEvent } from 'react'

import { useCallback, useEffect, useRef, useState } from 'react'
import { AnimatePresence, motion } from 'motion/react'

import { UploadCloudIcon } from '@archesai/ui/components/custom/icons'
import { cn } from '@archesai/ui/lib/utils'

interface FileError {
  code: string
  message: string
}

type FileStatus = 'dragging' | 'error' | 'idle' | 'uploading'

interface FileUploadProps {
  acceptedFileTypes?: string[]
  className?: string
  currentFile?: File | null
  maxFileSize?: number
  onFileRemove?: () => void
  onUploadError?: (error: FileError) => void
  onUploadSuccess?: (file: File) => void
  /** Duration in milliseconds for the upload simulation. Defaults to 2000ms (2s), 0 for no simulation */
  uploadDelay?: number
  validateFile?: (file: File) => FileError | null
}

const DEFAULT_MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB
const UPLOAD_STEP_SIZE = 5
const FILE_SIZES = [
  'Bytes',
  'KB',
  'MB',
  'GB',
  'TB',
  'PB',
  'EB',
  'ZB',
  'YB'
] as const

const formatBytes = (bytes: number, decimals = 2): string => {
  if (!bytes) return '0 Bytes'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const unit = FILE_SIZES[i] ?? FILE_SIZES[FILE_SIZES.length - 1]
  return `${Number.parseFloat((bytes / k ** i).toFixed(dm)).toString()} ${unit?.toString() ?? ''}`
}

const UploadIllustration = () => (
  <div className='relative h-16 w-16'>
    <svg
      aria-label='Upload illustration'
      className='h-full w-full'
      fill='none'
      viewBox='0 0 100 100'
      xmlns='http://www.w3.org/2000/svg'
    >
      <title>Upload File Illustration</title>
      <circle
        className='stroke-gray-200 dark:stroke-gray-700'
        cx='50'
        cy='50'
        r='45'
        strokeDasharray='4 4'
        strokeWidth='2'
      >
        <animateTransform
          attributeName='transform'
          dur='60s'
          from='0 50 50'
          repeatCount='indefinite'
          to='360 50 50'
          type='rotate'
        />
      </circle>

      <path
        className='fill-blue-100 stroke-blue-500 dark:fill-blue-900/30 dark:stroke-blue-400'
        d='M30 35H70C75 35 75 40 75 40V65C75 70 70 70 70 70H30C25 70 25 65 25 65V40C25 35 30 35 30 35Z'
        strokeWidth='2'
      >
        <animate
          attributeName='d'
          dur='2s'
          repeatCount='indefinite'
          values='
                        M30 35H70C75 35 75 40 75 40V65C75 70 70 70 70 70H30C25 70 25 65 25 65V40C25 35 30 35 30 35Z;
                        M30 38H70C75 38 75 43 75 43V68C75 73 70 73 70 73H30C25 73 25 68 25 68V43C25 38 30 38 30 38Z;
                        M30 35H70C75 35 75 40 75 40V65C75 70 70 70 70 70H30C25 70 25 65 25 65V40C25 35 30 35 30 35Z'
        />
      </path>

      <path
        className='stroke-blue-500 dark:stroke-blue-400'
        d='M30 35C30 35 35 35 40 35C45 35 45 30 50 30C55 30 55 35 60 35C65 35 70 35 70 35'
        fill='none'
        strokeWidth='2'
      />

      <g className='translate-y-2 transform'>
        <line
          className='stroke-blue-500 dark:stroke-blue-400'
          strokeLinecap='round'
          strokeWidth='2'
          x1='50'
          x2='50'
          y1='45'
          y2='60'
        >
          <animate
            attributeName='y2'
            dur='2s'
            repeatCount='indefinite'
            values='60;55;60'
          />
        </line>
        <polyline
          className='stroke-blue-500 dark:stroke-blue-400'
          fill='none'
          points='42,52 50,45 58,52'
          strokeLinecap='round'
          strokeLinejoin='round'
          strokeWidth='2'
        >
          <animate
            attributeName='points'
            dur='2s'
            repeatCount='indefinite'
            values='42,52 50,45 58,52;42,47 50,40 58,47;42,52 50,45 58,52'
          />
        </polyline>
      </g>
    </svg>
  </div>
)

const UploadingAnimation = ({ progress }: { progress: number }) => (
  <div className='relative h-16 w-16'>
    <svg
      aria-label={`Upload progress: ${Math.round(progress).toString()}%`}
      className='h-full w-full'
      fill='none'
      viewBox='0 0 240 240'
      xmlns='http://www.w3.org/2000/svg'
    >
      <title>Upload Progress Indicator</title>

      <defs>
        <mask id='progress-mask'>
          <rect
            fill='black'
            height='240'
            width='240'
          />
          <circle
            cx='120'
            cy='120'
            fill='white'
            r='120'
            strokeDasharray={`${((progress / 100) * 754).toString()}, 754`}
            transform='rotate(-90 120 120)'
          />
        </mask>
      </defs>

      <style>
        {`
                    @keyframes rotate-cw {
                        from { transform: rotate(0deg); }
                        to { transform: rotate(360deg); }
                    }
                    @keyframes rotate-ccw {
                        from { transform: rotate(360deg); }
                        to { transform: rotate(0deg); }
                    }
                    .g-spin circle {
                        transform-origin: 120px 120px;
                    }
                    .g-spin circle:nth-child(1) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(2) { animation: rotate-ccw 8s linear infinite; }
                    .g-spin circle:nth-child(3) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(4) { animation: rotate-ccw 8s linear infinite; }
                    .g-spin circle:nth-child(5) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(6) { animation: rotate-ccw 8s linear infinite; }
                    .g-spin circle:nth-child(7) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(8) { animation: rotate-ccw 8s linear infinite; }
                    .g-spin circle:nth-child(9) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(10) { animation: rotate-ccw 8s linear infinite; }
                    .g-spin circle:nth-child(11) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(12) { animation: rotate-ccw 8s linear infinite; }
                    .g-spin circle:nth-child(13) { animation: rotate-cw 8s linear infinite; }
                    .g-spin circle:nth-child(14) { animation: rotate-ccw 8s linear infinite; }

                    .g-spin circle:nth-child(2n) { animation-delay: 0.2s; }
                    .g-spin circle:nth-child(3n) { animation-delay: 0.3s; }
                    .g-spin circle:nth-child(5n) { animation-delay: 0.5s; }
                    .g-spin circle:nth-child(7n) { animation-delay: 0.7s; }
                `}
      </style>

      <g
        className='g-spin'
        mask='url(#progress-mask)'
        strokeDasharray='18% 40%'
        strokeWidth='10'
      >
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='150'
          stroke='#FF2E7E'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='140'
          stroke='#FFD600'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='130'
          stroke='#00E5FF'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='120'
          stroke='#FF3D71'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='110'
          stroke='#4ADE80'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='100'
          stroke='#2196F3'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='90'
          stroke='#FFA726'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='80'
          stroke='#FF1493'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='70'
          stroke='#FFEB3B'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='60'
          stroke='#00BCD4'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='50'
          stroke='#FF4081'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='40'
          stroke='#76FF03'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='30'
          stroke='#448AFF'
        />
        <circle
          cx='120'
          cy='120'
          opacity='0.95'
          r='20'
          stroke='#FF3D00'
        />
      </g>
    </svg>
  </div>
)

export default function FileUpload({
  acceptedFileTypes = [],
  className,
  currentFile: initialFile = null,
  maxFileSize = DEFAULT_MAX_FILE_SIZE,
  onFileRemove = () => ({}),
  onUploadError = () => ({}),
  onUploadSuccess = () => ({}),
  uploadDelay = 2000,
  validateFile = () => null
}: FileUploadProps) {
  const [file, setFile] = useState<File | null>(initialFile)
  const [status, setStatus] = useState<FileStatus>('idle')
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState<FileError | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const uploadIntervalRef = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    return () => {
      if (uploadIntervalRef.current) {
        clearInterval(uploadIntervalRef.current)
      }
    }
  }, [])

  const validateFileSize = useCallback(
    (file: File): FileError | null => {
      if (file.size > maxFileSize) {
        return {
          code: 'FILE_TOO_LARGE',
          message: `File size exceeds ${formatBytes(maxFileSize)}`
        }
      }
      return null
    },
    [maxFileSize]
  )

  const validateFileType = useCallback(
    (file: File): FileError | null => {
      if (!acceptedFileTypes.length) return null

      const fileType = file.type.toLowerCase()
      if (
        !acceptedFileTypes.some((type) => fileType.match(type.toLowerCase()))
      ) {
        return {
          code: 'INVALID_FILE_TYPE',
          message: `File type must be ${acceptedFileTypes.join(', ')}`
        }
      }
      return null
    },
    [acceptedFileTypes]
  )

  const handleError = useCallback(
    (error: FileError) => {
      setError(error)
      setStatus('error')
      onUploadError(error)

      setTimeout(() => {
        setError(null)
        setStatus('idle')
      }, 3000)
    },
    [onUploadError]
  )

  const simulateUpload = useCallback(
    (uploadingFile: File) => {
      let currentProgress = 0

      if (uploadIntervalRef.current) {
        clearInterval(uploadIntervalRef.current)
      }

      uploadIntervalRef.current = setInterval(
        () => {
          currentProgress += UPLOAD_STEP_SIZE
          if (currentProgress >= 100) {
            if (uploadIntervalRef.current) {
              clearInterval(uploadIntervalRef.current)
            }
            setProgress(0)
            setStatus('idle')
            setFile(null)
            onUploadSuccess(uploadingFile)
          } else {
            setStatus((prevStatus) => {
              if (prevStatus === 'uploading') {
                setProgress(currentProgress)
                return 'uploading'
              }
              if (uploadIntervalRef.current) {
                clearInterval(uploadIntervalRef.current)
              }
              return prevStatus
            })
          }
        },
        uploadDelay / (100 / UPLOAD_STEP_SIZE)
      )
    },
    [onUploadSuccess, uploadDelay]
  )

  const handleFileSelect = useCallback(
    (selectedFile: File | null) => {
      if (!selectedFile) return

      // Reset error state
      setError(null)

      // Validate file
      const sizeError = validateFileSize(selectedFile)
      if (sizeError) {
        handleError(sizeError)
        return
      }

      const typeError = validateFileType(selectedFile)
      if (typeError) {
        handleError(typeError)
        return
      }

      const customError = validateFile(selectedFile)
      if (customError) {
        handleError(customError)
        return
      }

      setFile(selectedFile)
      setStatus('uploading')
      setProgress(0)
      simulateUpload(selectedFile)
    },
    [
      simulateUpload,
      validateFileSize,
      validateFileType,
      validateFile,
      handleError
    ]
  )

  const handleDragOver = useCallback((e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setStatus((prev) => (prev !== 'uploading' ? 'dragging' : prev))
  }, [])

  const handleDragLeave = useCallback((e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setStatus((prev) => (prev === 'dragging' ? 'idle' : prev))
  }, [])

  const handleDrop = useCallback(
    (e: DragEvent<HTMLDivElement>) => {
      e.preventDefault()
      e.stopPropagation()
      if (status === 'uploading') return
      setStatus('idle')
      const droppedFile = e.dataTransfer.files[0]
      if (droppedFile) handleFileSelect(droppedFile)
    },
    [status, handleFileSelect]
  )

  const handleFileInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selectedFile = e.target.files?.[0]
      handleFileSelect(selectedFile ?? null)
      e.target.value = ''
    },
    [handleFileSelect]
  )

  const triggerFileInput = useCallback(() => {
    if (status === 'uploading') return
    fileInputRef.current?.click()
  }, [status])

  const resetState = useCallback(() => {
    setFile(null)
    setStatus('idle')
    setProgress(0)
    onFileRemove()
  }, [onFileRemove])

  return (
    <div
      aria-label='File upload'
      className={cn('relative mx-auto w-full max-w-sm', className ?? '')}
      role='complementary'
    >
      <div className='group relative w-full rounded-xl bg-white p-0.5 ring-1 ring-gray-200 dark:bg-black dark:ring-white/10'>
        <div className='absolute inset-x-0 -top-px h-px w-full bg-gradient-to-r from-transparent via-blue-500/20 to-transparent' />

        <div className='relative w-full rounded-[10px] bg-gray-50/50 p-1.5 dark:bg-white/[0.02]'>
          <div
            className={cn(
              'relative mx-auto w-full overflow-hidden rounded-lg border border-gray-100 bg-white dark:border-white/[0.08] dark:bg-black/50',
              error ? 'border-red-500/50' : ''
            )}
          >
            <div
              className={cn(
                'absolute inset-0 transition-opacity duration-300',
                status === 'dragging' ? 'opacity-100' : 'opacity-0'
              )}
            >
              <div className='absolute inset-x-0 top-0 h-[20%] bg-gradient-to-b from-blue-500/10 to-transparent' />
              <div className='absolute inset-x-0 bottom-0 h-[20%] bg-gradient-to-t from-blue-500/10 to-transparent' />
              <div className='absolute inset-y-0 left-0 w-[20%] bg-gradient-to-r from-blue-500/10 to-transparent' />
              <div className='absolute inset-y-0 right-0 w-[20%] bg-gradient-to-l from-blue-500/10 to-transparent' />
              <div className='absolute inset-[20%] animate-pulse rounded-lg bg-blue-500/5 transition-all duration-300' />
            </div>

            <div className='absolute -right-4 -top-4 h-8 w-8 bg-gradient-to-br from-blue-500/20 to-transparent opacity-0 blur-md transition-opacity duration-500 group-hover:opacity-100' />

            <div className='relative h-[240px]'>
              <AnimatePresence mode='wait'>
                {status === 'idle' || status === 'dragging' ?
                  <motion.div
                    animate={{
                      opacity: status === 'dragging' ? 0.8 : 1,
                      scale: status === 'dragging' ? 0.98 : 1,
                      y: 0
                    }}
                    className='absolute inset-0 flex flex-col items-center justify-center p-6'
                    exit={{ opacity: 0, y: -10 }}
                    initial={{ opacity: 0, y: 10 }}
                    key='dropzone'
                    onDragLeave={handleDragLeave}
                    onDragOver={handleDragOver}
                    onDrop={handleDrop}
                    transition={{ duration: 0.2 }}
                  >
                    <div className='mb-4'>
                      <UploadIllustration />
                    </div>

                    <div className='mb-4 space-y-1.5 text-center'>
                      <h3 className='text-lg font-semibold tracking-tight text-gray-900 dark:text-white'>
                        Drag and drop or
                      </h3>
                      <p className='text-xs text-gray-500 dark:text-gray-400'>
                        {acceptedFileTypes.length ?
                          acceptedFileTypes
                            .map((t) => t.split('/')[1])
                            .join(', ')
                            .toUpperCase()
                        : 'SVG, PNG, JPG or GIF'}{' '}
                        {maxFileSize && `up to ${formatBytes(maxFileSize)}`}
                      </p>
                    </div>

                    <button
                      className='group flex w-4/5 items-center justify-center gap-2 rounded-lg bg-gray-100 px-4 py-2.5 text-sm font-semibold text-gray-900 transition-all duration-200 hover:bg-gray-200 dark:bg-white/10 dark:text-white dark:hover:bg-white/20'
                      onClick={triggerFileInput}
                      type='button'
                    >
                      <span>Upload File</span>
                      <UploadCloudIcon className='h-4 w-4 transition-transform duration-200 group-hover:scale-110' />
                    </button>

                    <p className='mt-3 text-xs text-gray-500 dark:text-gray-400'>
                      or drag and drop your file here
                    </p>

                    <input
                      accept={acceptedFileTypes.join(',')}
                      aria-label='File input'
                      className='sr-only'
                      onChange={handleFileInputChange}
                      ref={fileInputRef}
                      type='file'
                    />
                  </motion.div>
                : status === 'uploading' ?
                  <motion.div
                    animate={{ opacity: 1, scale: 1 }}
                    className='absolute inset-0 flex flex-col items-center justify-center p-6'
                    exit={{ opacity: 0, scale: 0.95 }}
                    initial={{ opacity: 0, scale: 0.95 }}
                    key='uploading'
                  >
                    <div className='mb-4'>
                      <UploadingAnimation progress={progress} />
                    </div>

                    <div className='mb-4 space-y-1.5 text-center'>
                      <h3 className='truncate text-sm font-semibold text-gray-900 dark:text-white'>
                        {file?.name}
                      </h3>
                      <div className='flex items-center justify-center gap-2 text-xs'>
                        <span className='text-gray-500 dark:text-gray-400'>
                          {formatBytes(file?.size ?? 0)}
                        </span>
                        <span className='font-medium text-blue-500'>
                          {Math.round(progress)}%
                        </span>
                      </div>
                    </div>

                    <button
                      className='flex w-4/5 items-center justify-center gap-2 rounded-lg bg-gray-100 px-4 py-2.5 text-sm font-semibold text-gray-900 transition-all duration-200 hover:bg-gray-200 dark:bg-white/10 dark:text-white dark:hover:bg-white/20'
                      onClick={resetState}
                      type='button'
                    >
                      Cancel
                    </button>
                  </motion.div>
                : null}
              </AnimatePresence>
            </div>

            <AnimatePresence>
              {error && (
                <motion.div
                  animate={{ opacity: 1, y: 0 }}
                  className='absolute bottom-4 left-1/2 -translate-x-1/2 transform rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-2'
                  exit={{ opacity: 0, y: -10 }}
                  initial={{ opacity: 0, y: 10 }}
                >
                  <p className='text-sm text-red-500 dark:text-red-400'>
                    {error.message}
                  </p>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </div>
      </div>
    </div>
  )
}

FileUpload.displayName = 'FileUpload'
