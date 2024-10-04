"use client";
import { CustomCardForm, FormFieldConfig } from "@/components/custom-card-form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useChatbotsControllerFindOne } from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { useParams } from "next/navigation";
import * as z from "zod";
const formSchema = z.object({
  access_scope: z.string(),
  description: z.string().min(1).max(255),
  documents: z.string(),
  llm_base: z.string(),
  name: z.string().min(1).max(255),
});

export default function ChatbotConfigurationPage() {
  const { defaultOrgname } = useAuth();
  const { chatbotId } = useParams();
  const { data: agent, isLoading } = useChatbotsControllerFindOne({
    pathParams: {
      chatbotId: chatbotId as string,
      orgname: defaultOrgname,
    },
  });

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: agent?.name,
      description:
        "This is the name that will be used for this agent. If you do not specify one, it will automatically be set to your first message.",
      label: "Name",
      name: "name",
      props: {
        placeholder: "Agent name here...",
      },
      validationRule: formSchema.shape.name,
    },
    {
      component: Input,
      defaultValue: agent?.llmBase,
      description:
        "This is the LLM base that will be used for this agent. Note that different models have different credit usages.",
      label: "Language Model",
      name: "llm_base",
      props: {
        placeholder: "Search llms...",
      },
      validationRule: formSchema.shape.llm_base,
    },
    {
      component: Textarea,
      defaultValue: agent?.description,
      description:
        "This is the description that will be used for your agent. For example, 'You are a chatbot that will help people with their questions about Arches AI.'",
      label: "Description",
      name: "description",
      props: {
        placeholder: "Description here...",
      },
      validationRule: formSchema.shape.description,
    },
    {
      component: Input,
      description: "Description",
      label: "Documents",
      name: "documents",
      props: {
        placeholder: "Search documents...",
      },
      validationRule: formSchema.shape.description,
    },
  ];

  if (isLoading) {
    return null;
  }

  return (
    <CustomCardForm
      // description={"Configure your chatbot's settings"}
      fields={formFields}
      title="Configuration"
    />
  );
}
