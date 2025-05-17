'use client'

import { useCallback } from 'react'
import { parseAsArrayOf, parseAsJson, parseAsString, useQueryState } from 'nuqs'

import type { ContentEntity, ToolEntity } from '@archesai/domain'

export const usePlayground = () => {
  const [selectedTool, setSelectedTool] = useQueryState(
    'selectedTool',
    parseAsJson<ToolEntity>((tool) => tool as ToolEntity).withOptions({
      clearOnDefault: true
    })
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

  const clearParams = useCallback(async () => {
    await setSelectedContent(null)
    await setSelectedRunId(null)
    await setSelectedTool(null)
  }, [setSelectedContent, setSelectedRunId, setSelectedTool])

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
