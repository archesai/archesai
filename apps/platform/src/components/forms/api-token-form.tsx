import type { CreateApiTokenBody, UpdateApiTokenBody } from '@archesai/client'
import type { ApiTokenEntity } from '@archesai/schemas'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  useCreateApiToken,
  useGetOneApiTokenSuspense,
  useUpdateApiToken
} from '@archesai/client'
import { API_TOKEN_ENTITY_KEY, Type } from '@archesai/schemas'
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
  const { data: existingApiTokenResponse, error } =
    useGetOneApiTokenSuspense(apiTokenId)

  if (error) {
    return <div>API Token not found</div>
  }
  const apiToken = existingApiTokenResponse.data

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: apiToken.attributes.name || '',
      description: 'This is the name that will be used for this API token.',
      label: 'Name',
      name: 'name' as keyof ApiTokenEntity,
      props: {
        placeholder: 'API Token name here...'
      },
      renderControl: (field) => (
        <Input
          {...field}
          type='text'
        />
      ),
      validationRule: Type.String({
        minLength: 1
      })
    },
    {
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
      )
      // validationRule: Type.Union([Type.Literal('ADMIN'), Type.Literal('USER')])
    }
  ]

  return (
    <GenericForm<CreateApiTokenBody, UpdateApiTokenBody>
      description={'API Tokens are used to authenticate with the API.'}
      entityKey={API_TOKEN_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!apiTokenId}
      onSubmitCreate={async (createApiTokenDto) => {
        await createApiToken({
          data: {
            ...createApiTokenDto
          },
          id: apiTokenId!
        })
      }}
      onSubmitUpdate={async (updateApiTokenDto) => {
        await updateApiToken({
          data: {
            ...updateApiTokenDto,
            name: apiToken.attributes.name,
            organizationId: apiToken.attributes.organizationId,
            role: apiToken.attributes.role
          }
        })
      }}
      title={
        !apiTokenId ? 'Create API Token' : (
          `Update API Token: ${apiToken.attributes.name}`
        )
      }
    />
  )
}
