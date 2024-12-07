'use client'

import { FormFieldConfig, GenericForm } from '@/components/forms/generic-form/generic-form'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import {
  useContentControllerCreate,
  useContentControllerFindOne,
  useContentControllerUpdate
} from '@/generated/archesApiComponents'
import { CreateContentDto, UpdateContentDto } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { useState } from 'react'
import * as z from 'zod'

import ImportCard from '../import-card'

const formSchema = z.object({
  name: z.string(),
  text: z.string(),
  type: z.string()
})

export default function ContentForm({ contentId }: { contentId?: string }) {
  const { defaultOrgname } = useAuth()
  const [tab, setTab] = useState<'file' | 'text' | 'url'>('file')
  const { data: content } = useContentControllerFindOne(
    {
      pathParams: {
        contentId: contentId as string,
        orgname: defaultOrgname
      }
    },
    {
      enabled: !!defaultOrgname && !!contentId
    }
  )
  const { mutateAsync: updateContent } = useContentControllerUpdate({})
  const { mutateAsync: createContent } = useContentControllerCreate({})

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: content?.name,
      description: 'This is the name that will be used for this content.',
      label: 'Name',
      name: 'name',
      props: {
        placeholder: 'Content name here...'
      },
      validationRule: formSchema.shape.name
    },
    {
      component: Input,
      description: 'Select the content you would like to run the tool on. You can select multiple content items.',
      label: 'Input',
      name: tab === 'file' ? 'contentIds' : tab,
      renderControl: (field) => (
        <Tabs value={tab}>
          <TabsList className='grid w-full grid-cols-3 px-1'>
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
          <TabsContent value='text'>
            <Textarea {...field} placeholder='Enter text here' />
          </TabsContent>
          <TabsContent value='url'>
            <Textarea {...field} placeholder='Enter url here' rows={5} />
          </TabsContent>
          <TabsContent value='file'>
            <ImportCard
              cb={(content) => {
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
    <GenericForm<CreateContentDto, UpdateContentDto>
      description={!contentId ? 'Invite a new content' : 'Update an existing content'}
      fields={formFields}
      isUpdateForm={!!contentId}
      itemType='content'
      onSubmitCreate={async (createContentDto, mutateOptions) => {
        await createContent(
          {
            body: createContentDto,
            pathParams: {
              orgname: defaultOrgname
            }
          },
          mutateOptions
        )
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateContent(
          {
            body: data as any,
            pathParams: {
              contentId: contentId as string,
              orgname: defaultOrgname
            }
          },
          mutateOptions
        )
      }}
      title='Configuration'
    />
  )
}
