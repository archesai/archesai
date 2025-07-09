// import { useCallback } from 'react'
// import { useRouter, useSearch } from '@tanstack/react-router'

// import type { ArtifactEntity, ToolEntity } from '@archesai/schemas'

// export const usePlayground = () => {
//   const router = useRouter()
//   const search = useSearch({ strict: false })

//   const selectedTool = search.selectedTool as ToolEntity | undefined
//   const selectedContent = search.selectedContent as ArtifactEntity[] | undefined
//   const selectedRunId = search.selectedRunId as string | undefined

//   const setSelectedTool = useCallback(
//     (tool: null | ToolEntity) => {
//       router.setSearch((prev) => ({
//         ...prev,
//         selectedTool: tool ?? undefined
//       }))
//     },
//     [router]
//   )

//   const setSelectedContent = useCallback(
//     (content: ArtifactEntity[] | null) => {
//       router.setSearch((prev) => ({
//         ...prev,
//         selectedContent: content ?? undefined
//       }))
//     },
//     [router]
//   )

//   const setSelectedRunId = useCallback(
//     (runId: null | string) => {
//       router.setSearch((prev) => ({
//         ...prev,
//         selectedRunId: runId ?? undefined
//       }))
//     },
//     [router]
//   )

//   const clearParams = useCallback(() => {
//     router.setSearch((prev) => ({
//       ...prev,
//       selectedContent: undefined,
//       selectedRunId: undefined,
//       selectedTool: undefined
//     }))
//   }, [router])

//   return {
//     clearParams,
//     selectedContent: selectedContent ?? [],
//     selectedRunId: selectedRunId ?? '',
//     selectedTool,
//     setSelectedContent,
//     setSelectedRunId,
//     setSelectedTool
//   }
// }
