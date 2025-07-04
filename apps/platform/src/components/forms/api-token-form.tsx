import { Type } from '@sinclair/typebox'

import type { CreateApiTokenBody, UpdateApiTokenBody } from '@archesai/client'
import type { ApiTokenEntity } from '@archesai/domain'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateApiToken,
  useGetOneApiToken,
  useUpdateApiToken
} from '@archesai/client'
import { API_TOKEN_ENTITY_KEY } from '@archesai/domain'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { FormControl } from '@archesai/ui/components/shadcn/form'
import { Input } from '@archesai/ui/components/shadcn/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@archesai/ui/components/shadcn/select'

export default function APITokenForm({ apiTokenId }: { apiTokenId?: string }) {
  const { mutateAsync: createApiToken } = useUpdateApiToken({})
  const { mutateAsync: updateApiToken } = useCreateApiToken({})
  const { data: existingApiTokenResponse, error } = useGetOneApiToken(
    apiTokenId!,
    {
      query: {
        enabled: !!apiTokenId
      }
    }
  )

  if (error) {
    return <div>API Token not found</div>
  }
  const apiToken = existingApiTokenResponse!.data

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: apiToken.attributes.name,
      description: 'This is the name that will be used for this API token.',
      label: 'Name',
      name: 'name' as keyof ApiTokenEntity,
      props: {
        placeholder: 'API Token name here...'
      },
      validationRule: Type.String({
        maxLength: 128,
        minLength: 1
      })
    },
    {
      component: Input,
      defaultValue: apiToken.attributes.role,
      description: 'This is the role that will be used for this API token.',
      label: 'RoleTypeEnum',
      name: 'role' as keyof ApiTokenEntity,
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={(value) => {
            field.onChange(value)
          }}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder={'Choose your role...'} />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            {[
              {
                id: 'ADMIN',
                name: 'Admin'
              },
              {
                id: 'USER',
                name: 'User'
              }
            ].map((option) => (
              <SelectItem
                key={option.id}
                value={option.id.toString()}
              >
                {option.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ),
      validationRule: Type.Union([Type.Literal('ADMIN'), Type.Literal('USER')])
    }
  ]

  return (
    <GenericForm<CreateApiTokenBody, UpdateApiTokenBody>
      description={'API Tokens are used to authenticate with the API.'}
      entityKey={API_TOKEN_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!apiTokenId}
      onSubmitCreate={async (createApiTokenDto, mutateOptions) => {
        await createApiToken(
          {
            data: {
              ...createApiTokenDto
            },
            id: apiTokenId!
          },
          mutateOptions
        )
      }}
      onSubmitUpdate={async (updateApiTokenDto, mutateOptions) => {
        await updateApiToken(
          {
            data: {
              ...updateApiTokenDto,
              name: apiToken.attributes.name,
              orgname: apiToken.attributes.orgname,
              role: apiToken.attributes.role
            }
          },
          mutateOptions
        )
      }}
      title={
        !apiTokenId ? 'Create API Token' : (
          `Update API Token: ${apiToken.attributes.name}`
        )
      }
    />
  )
}
