import type { CreateMemberBody, UpdateMemberBody } from "@archesai/client";
import {
  createMember,
  updateMember,
  useGetMemberSuspense,
  useGetSessionSuspense,
} from "@archesai/client";
import type { FormFieldConfig } from "@archesai/ui/components/custom/generic-form";
import { GenericForm } from "@archesai/ui/components/custom/generic-form";
import { FormControl } from "@archesai/ui/components/shadcn/form";
import { Input } from "@archesai/ui/components/shadcn/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@archesai/ui/components/shadcn/select";
import { MEMBER_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { JSX } from "react";

export default function MemberForm({ id }: { id?: string }): JSX.Element {
  const { data: sessionData } = useGetSessionSuspense("current");
  const memberQuery = useGetMemberSuspense(
    sessionData.data.activeOrganizationID,
    id,
  );

  const member = id ? memberQuery.data.data : null;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: member?.userID,
      description: "This is the email that will be used for this member.",
      label: "User ID",
      name: "userID",
      renderControl: (field) => (
        <Input
          {...field}
          disabled={true}
          type="text"
        />
      ),
    },
    {
      defaultValue: member?.role,
      description:
        "This is the role that will be used for this member. Note that different roles have different permissions.",
      label: "RoleTypeEnum",
      name: "role",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder={"Choose your role..."} />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            {[
              {
                id: "ADMIN",
                name: "Admin",
              },
              {
                id: "USER",
                name: "User",
              },
            ].map((option) => (
              <SelectItem
                key={option.id}
                value={option.id}
              >
                {option.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ),
      // validationRule: Type.Union([Type.Literal('ADMIN'), Type.Literal('USER')])
    },
  ];

  return (
    <GenericForm<CreateMemberBody, UpdateMemberBody>
      description={!id ? "Invite a new member" : "Update an existing member"}
      entityKey={MEMBER_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createMemberDto) => {
        await createMember(
          sessionData.data.activeOrganizationID,
          createMemberDto,
        );
      }}
      onSubmitUpdate={async (data) => {
        if (id) {
          await updateMember(sessionData.data.activeOrganizationID, id, data);
        }
      }}
      title="Configuration"
    />
  );
}
