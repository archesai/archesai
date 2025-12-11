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
import type { UpdateSessionBody } from "#lib/index";
import { useGetSession, useUpdateSession } from "#lib/index";

export default function SessionForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateSession } = useUpdateSession();
  const { data: existingSession } = useGetSession(id, {
    query: { enabled: !!id },
  });
  const data = existingSession?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.authMethod ?? "",
      description:
        "The authentication method used (magic_link, oauth_google, oauth_github, etc.)",
      label: "Auth Method",
      name: "authMethod",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter auth method..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.authProvider ?? "",
      description:
        "The authentication provider (google, github, microsoft, local)",
      label: "Auth Provider",
      name: "authProvider",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select auth provider..." />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="local">Local</SelectItem>
            <SelectItem value="google">Google</SelectItem>
            <SelectItem value="github">Github</SelectItem>
            <SelectItem value="microsoft">Microsoft</SelectItem>
            <SelectItem value="apple">Apple</SelectItem>
          </SelectContent>
        </Select>
      ),
    },
    {
      defaultValue: data?.expiresAt ?? "",
      description: "The expiration date of the session",
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
      defaultValue: data?.ipAddress ?? "",
      description: "The IP address of the session",
      label: "IP Address",
      name: "ipAddress",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter ip address..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.organizationID ?? "",
      description:
        "The organization ID for this session (nullable for users without org)",
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
      defaultValue: data?.token ?? "",
      description: "The session token",
      label: "Token",
      name: "token",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter token..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.userAgent ?? "",
      description: "The user agent of the session",
      label: "User Agent",
      name: "userAgent",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter user agent..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.userID ?? "",
      description: "The user who owns this session",
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
    <GenericForm<UpdateSessionBody, UpdateSessionBody>
      description="Update an existing session"
      entityKey="sessions"
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (updateDto) => {
        await updateSession({ data: updateDto, id: id });
      }}
      title="Session"
    />
  );
}
