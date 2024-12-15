'use client'
import { ContentEntity, ToolEntity } from '@/generated/archesApiSchemas'
import { parseAsArrayOf, parseAsJson, parseAsString, useQueryState } from 'nuqs'
import { useCallback } from 'react'

export const usePlayground = () => {
  const [selectedTool, setSelectedTool] = useQueryState(
    'selectedTool',
    parseAsJson<any | ToolEntity>((tool) => tool as ToolEntity)
      .withOptions({
        clearOnDefault: true
      })
      .withDefault(undefined)
  )
  const [selectedContent, setSelectedContent] = useQueryState(
    'selectedContent',
    parseAsArrayOf(parseAsJson<ContentEntity>((tool) => tool as ContentEntity))
      .withOptions({
        clearOnDefault: true
      })
      .withDefault([])
  )
  const [selectedRunId, setSelectedRunId] = useQueryState(
    'selectedRunId',
    parseAsString
      .withOptions({
        clearOnDefault: true
      })
      .withDefault('')
  )

  const clearParams = useCallback(() => {
    setSelectedContent(null)
    setSelectedRunId(null)
    setSelectedTool(null)
  }, [])

  return {
    clearParams,
    selectedContent,
    selectedRunId,
    selectedTool,
    setSelectedContent,
    setSelectedRunId,
    setSelectedTool
  }
}
