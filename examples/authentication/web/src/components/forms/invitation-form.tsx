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
  Textarea,
} from "@archesai/ui";
import type { JSX } from "react";
import type { CreateInvitationBody, UpdateInvitationBody } from "#lib/index";
import {
  useCreateInvitation,
  useGetInvitationSuspense,
  useUpdateInvitation,
} from "#lib/index";

export default function InvitationForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateInvitation } = useUpdateInvitation();
  const { mutateAsync: createInvitation } = useCreateInvitation();
  const { data: existingInvitation } = useGetInvitationSuspense(id);

  const data = existingInvitation?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.email ?? "",
      description: "The email of the invitated user",
      label: "Email",
      name: "email",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter email..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.expiresAt ?? "",
      description: "The date and time when the invitation expires",
      label: "Expires At",
      name: "expiresAt",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter expires at..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.inviterID ?? "",
      description: "The ID of the user who sent this invitation",
      label: "Inviter ID",
      name: "inviterID",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter inviter id..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.organizationID ?? "",
      description: "The organization the user is being invited to join",
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
      description: "The role of the invitation",
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
      defaultValue: data?.status ?? "",
      description:
        "The status of the invitation, e.g., pending, accepted, declined",
      label: "Status",
      name: "status",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select status..." />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="pending">Pending</SelectItem>
            <SelectItem value="accepted">Accepted</SelectItem>
            <SelectItem value="declined">Declined</SelectItem>
            <SelectItem value="expired">Expired</SelectItem>
          </SelectContent>
        </Select>
      ),
    },
  ];

  return (
    <GenericForm<CreateInvitationBody, UpdateInvitationBody>
      description={
        !id ? "Create a new invitation" : "Update an existing invitation"
      }
      entityKey="invitations"
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createInvitation({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateInvitation({ data: updateDto, id: id });
      }}
      title="Invitation"
    />
  );
}
