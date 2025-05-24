'use client'

import type { TSchema } from '@sinclair/typebox'

import { useState } from 'react'

import type {
  CreateRunBody,
  FindManyContentsParams,
  FindManyToolsParams,
  UpdateRunBody
} from '@archesai/client'
import type { ContentEntity, RunEntity, ToolEntity } from '@archesai/domain'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateRun,
  useFindManyContents,
  useFindManyTools,
  useGetOneRun,
  useUpdateRun
} from '@archesai/client'
import {
  CONTENT_ENTITY_KEY,
  RunEntitySchema,
  TOOL_ENTITY_KEY,
  ToolEntitySchema
} from '@archesai/domain'
import { DataSelector } from '@archesai/ui/components/custom/data-selector'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import ImportCard from '@archesai/ui/components/custom/import-card'
import { Input } from '@archesai/ui/components/shadcn/input'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger
} from '@archesai/ui/components/shadcn/tabs'
import { Textarea } from '@archesai/ui/components/shadcn/textarea'
import { stringToColor } from '@archesai/ui/lib/utils'

import { usePlayground } from '#hooks/use-playground'
import { toolBaseIcons } from '#lib/site-config'

export function RunForm({ runId }: { runId?: string }) {
  const [tab, setTab] = useState<'contentIds' | 'file' | 'text' | 'url'>(
    'contentIds'
  )

  const { mutateAsync: createRun } = useCreateRun()
  const { mutateAsync: _updateRun } = useUpdateRun()
  const {
    selectedContent,
    selectedTool,
    setSelectedContent,
    setSelectedRunId,
    setSelectedTool
  } = usePlayground()
  const { data: existingRunResponse } = useGetOneRun(runId!, {
    query: {
      enabled: !!runId
    }
  })

  if (existingRunResponse?.status !== 200) {
    return <div>Run not found</div>
  }

  const formFields: FormFieldConfig<RunEntity>[] = [
    {
      component: Input,
      defaultValue: selectedTool?.id,
      description:
        'Select the tool you would like to run. Different tools have different inputs and outputs.',
      label: 'Tool',
      name: 'toolId',
      renderControl: (field) => (
        <DataSelector<ToolEntity, FindManyToolsParams>
          findManyParams={{}}
          getItemDetails={(tool) => {
            return (
              <div className='grid gap-2'>
                <h4 className='flex items-center gap-1 leading-none font-medium'>
                  {tool.name}
                </h4>
                <div className='text-sm text-muted-foreground'>
                  {tool.description}
                </div>
              </div>
            )
          }}
          icons={[
            {
              color: stringToColor('extract-text'),
              Icon: toolBaseIcons['extract-text'],
              name: 'Extract Text'
            },
            {
              color: stringToColor('create-embeddings'),
              Icon: toolBaseIcons['create-embeddings'],
              name: 'Create Embeddings'
            },
            {
              color: stringToColor('summarize'),
              Icon: toolBaseIcons.summarize,
              name: 'Summarize'
            },
            {
              color: stringToColor('text-to-image'),
              Icon: toolBaseIcons['text-to-image'],
              name: 'Text to Image'
            },
            {
              color: stringToColor('text-to-speech'),
              Icon: toolBaseIcons['text-to-speech'],
              name: 'Text to Speech'
            }
          ]}
          isMultiSelect={false}
          itemType={TOOL_ENTITY_KEY}
          selectedData={selectedTool!}
          setSelectedData={async (tool) => {
            if (!tool) {
              await setSelectedTool(null)
              field.onChange('')
              return
            }
            if (!Array.isArray(tool)) {
              await setSelectedTool(tool)
              field.onChange(tool.id)
            }
          }}
          useFindMany={useFindManyTools}
        />
      ),
      validationRule: ToolEntitySchema.properties.id as unknown as TSchema
    },
    {
      component: Input,
      description:
        'Select the content you would like to run the tool on. You can select multiple content items.',
      label: 'Input',
      name: tab === 'file' ? 'contentIds' : tab,
      renderControl: (field) => (
        <Tabs value={tab}>
          <TabsList className='grid w-full grid-cols-4 px-1'>
            <TabsTrigger
              onClick={() => {
                setTab('contentIds')
              }}
              value='contentIds'
            >
              Content
            </TabsTrigger>
            <TabsTrigger
              onClick={() => {
                setTab('text')
              }}
              value='text'
            >
              Text
            </TabsTrigger>
            <TabsTrigger
              onClick={() => {
                setTab('file')
              }}
              value='file'
            >
              File
            </TabsTrigger>
            <TabsTrigger
              onClick={() => {
                setTab('url')
              }}
              value='url'
            >
              URL
            </TabsTrigger>
          </TabsList>
          <TabsContent value='contentIds'>
            <DataSelector<ContentEntity, FindManyContentsParams>
              findManyParams={{}}
              isMultiSelect={true}
              itemType={CONTENT_ENTITY_KEY}
              selectedData={selectedContent}
              setSelectedData={async (content) => {
                if (!content) {
                  await setSelectedContent([])
                  field.onChange([])
                  return
                }
                if (!Array.isArray(content)) {
                  await setSelectedContent([content])
                  field.onChange([content.id])
                  return
                }
                await setSelectedContent(content)
                field.onChange(content.map((c) => c.id))
              }}
              useFindMany={useFindManyContents}
            />
          </TabsContent>
          <TabsContent value='text'>
            <Textarea
              {...field}
              placeholder='Enter text here'
            />
          </TabsContent>
          <TabsContent value='url'>
            <Textarea
              {...field}
              placeholder='Enter url here'
              rows={5}
            />
          </TabsContent>
          <TabsContent value='file'>
            <ImportCard
              cb={async (content) => {
                await setSelectedContent((old) => old.concat(content))
                setTab('contentIds')
                field.onChange(content.map((c) => c.id))
              }}
            />
          </TabsContent>
        </Tabs>
      ),
      validationRule: RunEntitySchema.properties
        .completedAt as unknown as TSchema
    }
  ]

  return (
    <GenericForm<RunEntity, CreateRunBody, UpdateRunBody>
      description={
        'Run a tool on a piece of content. You can select multiple content items.'
      }
      entityKey={TOOL_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={false}
      onSubmitCreate={async (createToolRunDto, mutateOptions) => {
        const run = await createRun(
          {
            data: createToolRunDto
          },
          mutateOptions
        )
        await setSelectedRunId(run.data.data.id)
      }}
      showCard={true}
      title='Try a Tool'
    />
  )
}
