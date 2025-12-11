import type { FormFieldConfig } from "@archesai/ui";
import { GenericForm, Input } from "@archesai/ui";
import type { JSX } from "react";
import type { UpdateSessionBody } from "#lib/index";
import { useGetSession, useUpdateSession } from "#lib/index";

export default function SessionForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateSession } = useUpdateSession();
  const { data: existingSession } = useGetSession(id, {
    query: { enabled: !!id },
  });
  const data = existingSession?.data;
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: data?.organizationID ?? "",
      description: "The organization ID to set as active for this session",
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
  ];
  return (
    <GenericForm<UpdateSessionBody, UpdateSessionBody>
      description="Update an existing session"
      entityKey="sessions"
      fields={updateFormFields}
      isUpdateForm={true}
      onSubmitUpdate={async (updateDto) => {
        await updateSession({ data: updateDto, id: id });
      }}
      title="Session"
    />
  );
}
