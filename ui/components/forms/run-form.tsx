'use client'
import { FormFieldConfig, GenericForm } from '@/components/forms/generic-form/generic-form'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { siteConfig } from '@/config/site'
import {
  useContentControllerFindAll,
  useRunsControllerCreate,
  useToolsControllerFindAll
} from '@/generated/archesApiComponents'
import { ContentEntity, CreateRunDto, ToolEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { usePlayground } from '@/hooks/use-playground'
import { useState } from 'react'
import * as z from 'zod'

import { DataSelector } from '../data-selector'
import ImportCard from '../import-card'
import { Textarea } from '../ui/textarea'

const formSchema = z.object({
  contentIds: z.array(z.string()).optional(),
  text: z.string().optional(),
  toolId: z.string().min(1, 'Tool selection is required'),
  url: z.string().optional()
})

export default function RunForm() {
  const { defaultOrgname } = useAuth()
  const { mutateAsync: runTool } = useRunsControllerCreate()
  const { selectedContent, selectedTool, setSelectedContent, setSelectedRunId, setSelectedTool } = usePlayground()

  const [tab, setTab] = useState<'contentIds' | 'file' | 'text' | 'url'>('contentIds')

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: selectedTool?.id,
      description: 'Select the tool you would like to run. Different tools have different inputs and outputs.',
      label: 'Tool',
      name: 'toolId',
      renderControl: (field) => (
        <DataSelector<ToolEntity>
          getItemDetails={(tool) => {
            return (
              <div className='grid gap-2'>
                <h4 className='flex items-center gap-1 font-medium leading-none'>{tool?.name}</h4>
                <div className='text-sm text-muted-foreground'>{tool?.description}</div>
              </div>
            )
          }}
          icons={[
            {
              Icon: siteConfig.toolBaseIcons['extract-text'],
              name: 'Extract Text'
            },
            {
              Icon: siteConfig.toolBaseIcons['create-embeddings'],
              name: 'Create Embeddings'
            },
            {
              Icon: siteConfig.toolBaseIcons['summarize'],
              name: 'Summarize'
            },
            {
              Icon: siteConfig.toolBaseIcons['text-to-image'],
              name: 'Text to Image'
            },
            {
              Icon: siteConfig.toolBaseIcons['text-to-speech'],
              name: 'Text to Speech'
            }
          ]}
          isMultiSelect={false}
          itemType='tool'
          selectedData={selectedTool as ToolEntity}
          setSelectedData={(tool: any) => {
            setSelectedTool(tool)
            field.onChange(tool.id)
          }}
          useFindAll={() =>
            useToolsControllerFindAll({
              pathParams: {
                orgname: defaultOrgname
              }
            })
          }
        />
      ),
      validationRule: formSchema.shape.toolId
    },
    {
      component: Input,
      description: 'Select the content you would like to run the tool on. You can select multiple content items.',
      label: 'Input',
      name: tab === 'file' ? 'contentIds' : tab,
      renderControl: (field) => (
        <Tabs value={tab}>
          <TabsList className='grid w-full grid-cols-4 px-1'>
            <TabsTrigger onClick={() => setTab('contentIds')} value='contentIds'>
              Content
            </TabsTrigger>
            <TabsTrigger onClick={() => setTab('text')} value='text'>
              Text
            </TabsTrigger>
            <TabsTrigger onClick={() => setTab('file')} value='file'>
              File
            </TabsTrigger>
            <TabsTrigger onClick={() => setTab('url')} value='url'>
              URL
            </TabsTrigger>
          </TabsList>
          <TabsContent value='contentIds'>
            <DataSelector<ContentEntity>
              isMultiSelect={true}
              itemType='content'
              selectedData={selectedContent as ContentEntity[]}
              setSelectedData={(content: any) => {
                setSelectedContent(content)
                field.onChange(content === null ? [] : content.map((c: any) => c.id))
              }}
              useFindAll={() =>
                useContentControllerFindAll({
                  pathParams: {
                    orgname: defaultOrgname
                  }
                })
              }
            />
          </TabsContent>
          <TabsContent value='text'>
            <Textarea {...field} placeholder='Enter text here' />
          </TabsContent>
          <TabsContent value='url'>
            <Textarea {...field} placeholder='Enter url here' rows={5} />
          </TabsContent>
          <TabsContent value='file'>
            <ImportCard
              cb={(content) => {
                setSelectedContent((old) => (old || [])?.concat(content))
                setTab('contentIds')
                field.onChange(content.map((c: any) => c.id))
              }}
            />
          </TabsContent>
        </Tabs>
      ),
      validationRule: (formSchema.shape as any)[tab]
    }
  ]

  return (
    <GenericForm<CreateRunDto, any>
      description={'Run a tool on a piece of content. You can select multiple content items.'}
      fields={formFields}
      isUpdateForm={false}
      itemType='tool run'
      onSubmitCreate={async (createToolRunDto, mutateOptions) => {
        const run = await runTool(
          {
            body: {
              contentIds: createToolRunDto.contentIds,
              runType: 'TOOL_RUN',
              text: createToolRunDto.text,
              toolId: createToolRunDto.toolId
            },
            pathParams: {
              orgname: defaultOrgname
            }
          },
          mutateOptions
        )
        setSelectedRunId(run.id)
      }}
      showCard={true}
      title='Try a Tool'
    />
  )
}
