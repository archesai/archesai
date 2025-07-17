import { createFileRoute } from '@tanstack/react-router'

import FileUpload from '#components/file-upload'

export const Route = createFileRoute('/_app/')({
  component: AppIndex
})

export default function AppIndex() {
  return (
    <div>
      <FileUpload />
    </div>
  )
}
