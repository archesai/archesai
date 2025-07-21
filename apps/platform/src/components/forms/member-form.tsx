import type { CreateMemberBody, UpdateMemberBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  createMember,
  updateMember,
  useGetOneMemberSuspense
} from '@archesai/client'
import { MEMBER_ENTITY_KEY, StringSchema } from '@archesai/schemas'
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

export default function MemberForm({ id }: { id?: string }) {
  const { data: existingMemberResponse } = useGetOneMemberSuspense(id)

  const member = existingMemberResponse.data

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: member.userId,
      description: 'This is the email that will be used for this member.',
      label: 'User ID',
      name: 'userId',
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type='text'
        />
      ),
      validationRule: StringSchema
    },
    {
      defaultValue: member.role,
      description:
        'This is the role that will be used for this member. Note that different roles have different permissions.',
      label: 'RoleTypeEnum',
      name: 'role',
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
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
    <GenericForm<CreateMemberBody, UpdateMemberBody>
      description={!id ? 'Invite a new member' : 'Update an existing member'}
      entityKey={MEMBER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createMemberDto) => {
        await createMember(createMemberDto)
      }}
      onSubmitUpdate={async (data) => {
        await updateMember(id, data)
      }}
      title='Configuration'
    />
  )
}
