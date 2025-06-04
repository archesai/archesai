/* eslint-disable @typescript-eslint/no-unnecessary-condition */
'use client'

import { useFindManyRuns } from '@archesai/client'
import { StatusTypeEnumButton } from '@archesai/ui/components/custom/run-status-button'
import { cn } from '@archesai/ui/lib/utils'

import ArtifactDataTable from '#components/datatables/artifact-datatable'
import { RunForm } from '#components/forms/run-form'
import { ToolCards } from '#components/tool-cards'
import { usePlayground } from '#hooks/use-playground'

export default function Playground() {
  const { selectedRunId, selectedTool } = usePlayground()

  const { data: runsResponse } = useFindManyRuns()

  const runs = runsResponse?.data.data ?? []
  const run = runs.find((r) => r.id === selectedRunId)

  const hasInputs = false
  const hasOutputs = false

  return selectedTool ?
      <div className='flex h-full min-h-0 gap-3'>
        <div
          className={cn(
            'flex w-2/3 flex-1 flex-col gap-4 overflow-auto',
            !hasInputs ? 'hidden' : ''
          )}
        >
          {hasInputs ?
            <div
              className={cn(
                'overflow-auto transition-all',
                hasOutputs ? 'h-1/2' : 'h-full'
              )}
            >
              <ArtifactDataTable readonly />
            </div>
          : null}

          {hasOutputs ?
            <div className='h-1/2 overflow-auto'>
              <ArtifactDataTable readonly />
            </div>
          : null}
        </div>

        <div
          className={cn(
            'flex flex-col gap-1 transition-all',
            !hasInputs ? 'h-auto max-h-0 w-full items-center' : 'w-1/3 gap-3'
          )}
        >
          {selectedRunId && run && (
            <StatusTypeEnumButton
              run={{
                ...run.attributes,
                id: run.id,
                type: run.type
              }}
            />
          )}
          <RunForm />
        </div>
      </div>
    : <ToolCards />
}
