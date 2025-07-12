import type { TSchema } from '@sinclair/typebox'

import { useState } from 'react'
import { Type } from '@sinclair/typebox'

import type { CreateArtifactBody, UpdateArtifactBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  createArtifact,
  updateArtifact,
  useGetOneArtifactSuspense
} from '@archesai/client'
import { ARTIFACT_ENTITY_KEY, ArtifactEntitySchema } from '@archesai/schemas'
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

export default function ContentForm({ artifactId }: { artifactId?: string }) {
  const [tab, setTab] = useState<'file' | 'text' | 'url'>('file')

  const { data: existingContentResponse, error } =
    useGetOneArtifactSuspense(artifactId)

  if (error) {
    return <div>Content not found</div>
  }
  const content = existingContentResponse.data

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: content.attributes.name,
      description: 'This is the name that will be used for this content.',
      label: 'Name',
      name: 'name',
      props: {
        placeholder: 'Content name here...'
      },
      renderControl: (field) => (
        <Input
          {...field}
          type='text'
        />
      ),
      validationRule: Type.String({
        maxLength: 128,
        minLength: 1
      })
    },
    {
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
      validationRule: ArtifactEntitySchema.properties[
        tab == 'file' ? 'text' : tab
      ] as unknown as TSchema
    }
  ]

  return (
    <GenericForm<CreateArtifactBody, UpdateArtifactBody>
      description={
        !artifactId ? 'Invite a new content' : 'Update an existing content'
      }
      entityKey={ARTIFACT_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!artifactId}
      onSubmitCreate={async (createArtifactDto) => {
        await createArtifact(createArtifactDto)
      }}
      onSubmitUpdate={async (updateContentDto) => {
        await updateArtifact(artifactId, updateContentDto)
      }}
      title='Configuration'
    />
  )
}
