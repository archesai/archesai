"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  useChatbotsControllerCreate,
  useChatbotsControllerFindOne,
  useChatbotsControllerUpdate,
} from "@/generated/archesApiComponents";
import {
  CreateChatbotDto,
  UpdateChatbotDto,
} from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
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
  llmBase: z.enum(["gpt-4o", "gpt-4o-mini"], {
    message: "Invalid language model. Must be one of 'gpt-4o', 'gpt-4o-mini'.",
  }),
  name: z.string().min(1).max(255),
});

export default function ChatbotForm({ chatbotId }: { chatbotId?: string }) {
  const { defaultOrgname } = useAuth();
  const { data: chatbot, isLoading } = useChatbotsControllerFindOne(
    {
      pathParams: {
        chatbotId: chatbotId as string,
        orgname: defaultOrgname,
      },
    },
    {
      enabled: !!defaultOrgname && !!chatbotId,
    }
  );
  const { mutateAsync: updateChatbot } = useChatbotsControllerUpdate({});
  const { mutateAsync: createChatbot } = useChatbotsControllerCreate({});

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: chatbot?.name,
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
      defaultValue: chatbot?.llmBase || "gpt-4o",
      description:
        "This is the LLM base that will be used for this chatbot. Note that different models have different credit usages.",
      label: "Language Model",
      name: "llmBase",
      props: {
        placeholder: "Search llms...",
      },
      renderControl: (field) => (
        <Select
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
                id: "gpt-4o",
                name: "GPT-4o",
              },
              {
                id: "gpt-4o-mini",
                name: "GPT-4o Mini",
              },
            ].map((option) => (
              <SelectItem key={option.id} value={option.id.toString()}>
                {option.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ),
      validationRule: formSchema.shape.llmBase,
    },
    {
      component: Textarea,
      defaultValue: chatbot?.description,
      description:
        "This is the description that will be used for your chatbot. For example, 'You are a chatbot that will help people with their questions about Arches AI.'",
      label: "Description",
      name: "description",
      props: {
        className: "h-40",
        placeholder: "Description here...",
      },
      validationRule: formSchema.shape.description,
    },
  ];

  if (isLoading) {
    return null;
  }

  return (
    <GenericForm<CreateChatbotDto, UpdateChatbotDto>
      description={"Configure your chatbot's settings"}
      fields={formFields}
      isUpdateForm={!!chatbotId}
      itemType="chatbot"
      onSubmitCreate={async (createChatbotDto, mutateOptions) => {
        await createChatbot(
          {
            body: createChatbotDto,
            pathParams: {
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateChatbot(
          {
            body: data as any,
            pathParams: {
              chatbotId: chatbotId as string,
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
