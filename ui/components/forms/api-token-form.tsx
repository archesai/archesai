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
import { useAuth } from "@/hooks/useAuth";
import * as z from "zod";

const formSchema = z.object({
  description: z.string().min(1),
  llmBase: z.enum(["gpt-4o", "gpt-4o-mini"], {
    message: "Invalid language model. Must be one of 'gpt-4o', 'gpt-4o-mini'.",
  }),
  name: z.string().min(1).max(255),
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
      description: "This is the name that will be used for this chatbot.",
      label: "Name",
      name: "name",
      props: {
        placeholder: "Chatbot name here...",
      },
      validationRule: formSchema.shape.name,
    },
    {
      component: Input,
      defaultValue: apiToken?.role,
      description:
        "This is the LLM base that will be used for this chatbot. Note that different models have different credit usages.",
      label: "Language Model",
      name: "llmBase",
      props: {
        placeholder: "Search llms...",
      },
      validationRule: formSchema.shape.llmBase,
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
