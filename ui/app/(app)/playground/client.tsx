'use client'
import ContentDataTable from '@/components/datatables/content-datatable'
import RunForm from '@/components/forms/run-form'
import { RunStatusButton } from '@/components/run-status-button'
import { ToolCards } from '@/components/tool-cards'
import { useRunsControllerFindOne } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { usePlayground } from '@/hooks/use-playground'
import { cn } from '@/lib/utils'
import { Suspense } from 'react'

export function Client() {
  const { defaultOrgname } = useAuth()
  const { selectedContent, selectedRunId, selectedTool } = usePlayground()

  const { data: run } = useRunsControllerFindOne(
    {
      pathParams: {
        id: selectedRunId,
        orgname: defaultOrgname
      }
    },
    {
      enabled: !!selectedRunId
    }
  )

  const hasInputs = !!selectedContent.length || !!run?.inputs?.length
  const hasOutputs = !!run?.outputs?.length

  return selectedTool ? (
    <div className='flex h-full min-h-0 gap-3'>
      <div
        className={cn(
          'flex w-2/3 flex-1 flex-col gap-4 overflow-auto',
          !hasInputs ? 'hidden' : ''
        )}
      >
        {hasInputs ? (
          <div
            className={cn(
              'overflow-auto transition-all',
              hasOutputs ? 'h-1/2' : 'h-full'
            )}
          >
            <Suspense fallback={<p>Loading feed...</p>}>
              <ContentDataTable
                customFilters={[
                  {
                    field: 'id',
                    operator: 'in',
                    value: run
                      ? run.inputs.map((r) => r.id)
                      : selectedContent?.map((r) => r.id) || []
                  }
                ]}
                readonly
              />
            </Suspense>
          </div>
        ) : null}
        {/* <Separator /> */}

        {run?.outputs?.length ? (
          <div className='h-1/2 overflow-auto'>
            <ContentDataTable
              customFilters={[
                {
                  field: 'id',
                  operator: 'in',
                  value: run ? run.outputs.map((r) => r.id) : []
                }
              ]}
              readonly
            />
          </div>
        ) : null}
      </div>

      {/* SIDEBAR */}
      <div
        className={cn(
          'flex flex-col gap-1 transition-all',
          !hasInputs
            ? 'h-auto w-full items-center justify-center py-24'
            : 'w-1/3 gap-3'
        )}
      >
        {selectedRunId && run && <RunStatusButton run={run} />}
        <RunForm />
      </div>
    </div>
  ) : (
    <ToolCards />
  )
}
