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
import type { UpdateAccountBody } from "#lib/index";
import { useGetAccount, useUpdateAccount } from "#lib/index";

export default function AccountForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateAccount } = useUpdateAccount();
  const { data: existingAccount } = useGetAccount(id, {
    query: { enabled: !!id },
  });
  const data = existingAccount?.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: data?.accessToken ?? "",
      description: "The OAuth access token",
      label: "Access Token",
      name: "accessToken",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter access token..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.accessTokenExpiresAt ?? "",
      description: "The access token expiration timestamp",
      label: "Access Token Expires At",
      name: "accessTokenExpiresAt",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter access token expires at..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.accountIdentifier ?? "",
      description: "The unique identifier for the account from the provider",
      label: "Account Identifier",
      name: "accountIdentifier",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter account identifier..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.idToken ?? "",
      description: "The OpenID Connect ID token",
      label: "ID Token",
      name: "idToken",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter id token..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.provider ?? "",
      description: "The authentication provider identifier",
      label: "Provider",
      name: "provider",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select provider..." />
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
      defaultValue: data?.refreshToken ?? "",
      description: "The OAuth refresh token",
      label: "Refresh Token",
      name: "refreshToken",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter refresh token..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.refreshTokenExpiresAt ?? "",
      description: "The refresh token expiration timestamp",
      label: "Refresh Token Expires At",
      name: "refreshTokenExpiresAt",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter refresh token expires at..."
          type="text"
        />
      ),
    },
    {
      defaultValue: data?.scope ?? "",
      description: "The OAuth scope granted",
      label: "Scope",
      name: "scope",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter scope..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: data?.userID ?? "",
      description: "The user ID this account belongs to",
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
    <GenericForm<UpdateAccountBody, UpdateAccountBody>
      description="Update an existing account"
      entityKey="accounts"
      fields={formFields}
      isUpdateForm={true}
      onSubmitUpdate={async (updateDto) => {
        await updateAccount({ data: updateDto, id: id });
      }}
      title="Account"
    />
  );
}
