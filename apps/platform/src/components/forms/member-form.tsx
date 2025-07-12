import type { CreateMemberBody, UpdateMemberBody } from '@archesai/client'
import type { FormFieldConfig } from '@archesai/ui/components/custom/generic-form'

import {
  createMember,
  updateMember,
  useGetOneMemberSuspense
} from '@archesai/client'
import { MEMBER_ENTITY_KEY, Type } from '@archesai/schemas'
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

export default function MemberForm({ memberId }: { memberId?: string }) {
  const { data: existingMemberResponse } = useGetOneMemberSuspense(memberId)

  const member = existingMemberResponse.data

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: member.attributes.invitationId,
      description: 'This is the email that will be used for this member.',
      label: 'E-Mail',
      name: 'invitationId',
      props: {
        placeholder: 'Member email here...'
      },
      renderControl: (field) => (
        <Input
          {...field}
          type='text'
        />
      ),
      validationRule: Type.String({
        format: 'email',
        minLength: 1
      })
    },
    {
      defaultValue: member.attributes.role,
      description:
        'This is the role that will be used for this member. Note that different roles have different permissions.',
      label: 'RoleTypeEnum',
      name: 'role',
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
    <GenericForm<CreateMemberBody, UpdateMemberBody>
      description={
        !memberId ? 'Invite a new member' : 'Update an existing member'
      }
      entityKey={MEMBER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!memberId}
      onSubmitCreate={async (createMemberDto) => {
        await createMember(createMemberDto)
      }}
      onSubmitUpdate={async (data) => {
        await updateMember(memberId, data)
      }}
      title='Configuration'
    />
  )
}
