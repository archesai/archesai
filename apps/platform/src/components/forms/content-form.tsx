'use client'

import type { TSchema } from '@sinclair/typebox'

import { useState } from 'react'
import { Type } from '@sinclair/typebox'

import type { CreateContentBody, UpdateContentBody } from '@archesai/client'
import type { ContentEntity } from '@archesai/domain'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  createContent,
  updateContent,
  useGetOneContent
} from '@archesai/client'
import { CONTENT_ENTITY_KEY, ContentEntitySchema } from '@archesai/domain'
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

export default function ContentForm({ contentId }: { contentId?: string }) {
  const [tab, setTab] = useState<'file' | 'text' | 'url'>('file')

  const { data: existingContentResponse } = useGetOneContent(contentId!, {
    query: {
      enabled: !!contentId
    }
  })

  if (existingContentResponse?.status !== 200) {
    return <div>Content not found</div>
  }
  const content = existingContentResponse.data.data

  const formFields: FormFieldConfig<ContentEntity>[] = [
    {
      component: Input,
      defaultValue: content.attributes.name,
      description: 'This is the name that will be used for this content.',
      label: 'Name',
      name: 'name',
      props: {
        placeholder: 'Content name here...'
      },
      validationRule: Type.String({
        maxLength: 128,
        minLength: 1
      })
    },
    {
      component: Input,
      description:
        'Select the content you would like to run the tool on. You can select multiple content items.',
      label: 'Input',
      name: tab === 'file' ? 'text' : tab,
      renderControl: (field) => (
        <Tabs value={tab}>
          <TabsList className='grid w-full grid-cols-3 px-1'>
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
          <TabsContent value='text'>
            <Textarea
              {...field}
              placeholder='Enter text here'
              value={field.value as string}
            />
          </TabsContent>
          <TabsContent value='url'>
            <Textarea
              {...field}
              placeholder='Enter url here'
              rows={5}
              value={field.value as string}
            />
          </TabsContent>
          <TabsContent value='file'>
            <ImportCard
              cb={(content) => {
                field.onChange(content.map((c) => c.id))
              }}
            />
          </TabsContent>
        </Tabs>
      ),
      validationRule: ContentEntitySchema.properties[
        tab == 'file' ? 'text' : tab
      ] as unknown as TSchema
    }
  ]

  return (
    <GenericForm<ContentEntity, CreateContentBody, UpdateContentBody>
      description={
        !contentId ? 'Invite a new content' : 'Update an existing content'
      }
      entityKey={CONTENT_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!contentId}
      onSubmitCreate={async (createContentDto, mutateOptions) => {
        await createContent(createContentDto, mutateOptions)
      }}
      onSubmitUpdate={async (updateContentDto, mutateOptions) => {
        await updateContent(contentId!, updateContentDto, mutateOptions)
      }}
      title='Configuration'
    />
  )
}
