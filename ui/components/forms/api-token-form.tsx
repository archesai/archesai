"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import {
  useApiTokensControllerCreate,
  useApiTokensControllerFindOne,
  useApiTokensControllerUpdate,
} from "@/generated/archesApiComponents";
import {
  CreateApiTokenDto,
  UpdateApiTokenDto,
} from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import * as z from "zod";

const formSchema = z.object({
  domains: z.string(),
  name: z.string(),
  role: z.enum(["USER", "ADMIN"]),
});

export default function APITokenForm({ apiTokenId }: { apiTokenId?: string }) {
  const { defaultOrgname } = useAuth();
  const { data: apiToken, isLoading } = useApiTokensControllerFindOne(
    {
      pathParams: {
        id: apiTokenId as string,
        orgname: defaultOrgname,
      },
    },
    {
      enabled: !!defaultOrgname && !!apiTokenId,
    }
  );
  const { mutateAsync: updateApiToken } = useApiTokensControllerUpdate({});
  const { mutateAsync: createChatbot } = useApiTokensControllerCreate({});

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: apiToken?.name,
      description: "This is the name that will be used for this API token.",
      label: "Name",
      name: "name",
      props: {
        placeholder: "API Token name here...",
      },
      validationRule: formSchema.shape.name,
    },
    {
      component: Input,
      defaultValue: apiToken?.domains || "",
      description:
        "These are the domains that will be used for this API token.",
      label: "Domains",
      name: "domains",
      props: {
        placeholder: "http://example.com, https://example.com",
      },
      validationRule: formSchema.shape.domains,
    },
    {
      component: Input,
      defaultValue: apiToken?.role,
      description: "This is the role that will be used for this API token.",
      label: "Role",
      name: "role",
      props: {
        placeholder: "Search llms...",
      },
      validationRule: formSchema.shape.role,
    },
  ];

  if (isLoading) {
    return null;
  }

  return (
    <GenericForm<CreateApiTokenDto, UpdateApiTokenDto>
      description={"Configure your API Tokens's settings"}
      fields={formFields}
      isUpdateForm={!!apiTokenId}
      itemType="API Token"
      onSubmitCreate={async (createApiTokenDto, mutateOptions) => {
        await createChatbot(
          {
            body: createApiTokenDto,
            pathParams: {
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateApiToken(
          {
            body: data as any,
            pathParams: {
              id: apiTokenId as string,
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
