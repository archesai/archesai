import type { FormFieldConfig } from "@archesai/ui";
import {
  FormControl,
  GenericForm,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@archesai/ui";
import type { JSX } from "react";
import type { CreateMemberBody, UpdateMemberBody } from "#lib/index";
import {
  useCreateMember,
  useGetMemberSuspense,
  useUpdateMember,
} from "#lib/index";

export default function MemberForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateMember } = useUpdateMember();
  const { mutateAsync: createMember } = useCreateMember();
  const { data: existingMember } = useGetMemberSuspense(id);

  const data = existingMember?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.organizationID ?? "",
      description: "The organization this member belongs to",
      label: "Organization ID",
      name: "organizationID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter organization id..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.role ?? "",
      description: "The role of the member",
      label: "Role",
      name: "role",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select role..." />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="admin">Admin</SelectItem>
            <SelectItem value="owner">Owner</SelectItem>
            <SelectItem value="basic">Basic</SelectItem>
          </SelectContent>
        </Select>
      ),
    },
    {
      defaultValue: data?.userID ?? "",
      description: "The user who is a member of the organization",
      label: "User ID",
      name: "userID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter user id..."
          type="text"
        />
      ),
    },
  ];

  return (
    <GenericForm<CreateMemberBody, UpdateMemberBody>
      description={!id ? "Create a new member" : "Update an existing member"}
      entityKey="members"
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createMember({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateMember({ data: updateDto, id: id });
      }}
      title="Member"
    />
  );
}
