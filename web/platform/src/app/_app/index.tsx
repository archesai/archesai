import type { JSX } from "react"

import { createFileRoute } from "@tanstack/react-router"

import FileUpload from "#components/file-upload"

export const Route = createFileRoute("/_app/")({
  component: AppIndex
})

export default function AppIndex(): JSX.Element {
  return (
    <div className="-mt-16 flex h-full flex-col items-center justify-center">
      <FileUpload />
    </div>
  )
}
