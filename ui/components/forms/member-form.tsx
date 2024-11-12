"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import {
  useMembersControllerCreate,
  useMembersControllerUpdate,
} from "@/generated/archesApiComponents";
import { CreateMemberDto, UpdateMemberDto } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import * as z from "zod";

import { FormControl } from "../ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

const formSchema = z.object({
  description: z.string().min(1),
  email: z.string().email(),
  name: z.string().min(1).max(255),
  role: z.enum(["ADMIN", "USER"], {
    message: "Invalid role. Must be one of 'ADMIN', 'USER'.",
  }),
});

export default function MemberForm({ memberId }: { memberId?: string }) {
  const { defaultOrgname } = useAuth();
  const { mutateAsync: updateMember } = useMembersControllerUpdate({});
  const { mutateAsync: createMember } = useMembersControllerCreate({});

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: "",
      description: "This is the name that will be used for this member.",
      label: "Name",
      name: "name",
      props: {
        placeholder: "Member name here...",
      },
      validationRule: formSchema.shape.name,
    },
    {
      component: Input,
      defaultValue: "",
      description: "This is the email that will be used for this member.",
      label: "E-Mail",
      name: "inviteEmail",
      props: {
        placeholder: "Member email here...",
      },
      validationRule: formSchema.shape.email,
    },
    {
      component: Input,
      defaultValue: "admin",
      description:
        "This is the role that will be used for this member. Note that different roles have different permissions.",
      label: "Role",
      name: "role",
      props: {
        placeholder: "Search roles...",
      },
      renderControl: (field) => (
        <Select
          defaultValue="ADMIN"
          onValueChange={(value) => field.onChange(value)}
          value={String(field.value) || undefined}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder={"Choose your model"} />
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
              <SelectItem key={option.id} value={option.id.toString()}>
                {option.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ),
      validationRule: formSchema.shape.role,
    },
  ];

  return (
    <GenericForm<CreateMemberDto, UpdateMemberDto>
      description={"Configure your member's settings"}
      fields={formFields}
      isUpdateForm={!!memberId}
      itemType="member"
      onSubmitCreate={async (createMemberDto, mutateOptions) => {
        await createMember(
          {
            body: createMemberDto,
            pathParams: {
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateMember(
          {
            body: data as any,
            pathParams: {
              memberId: memberId as string,
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      title="Configuration"
    />
  );
}
