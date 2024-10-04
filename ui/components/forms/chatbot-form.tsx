"use client";
import { CustomCardForm, FormFieldConfig } from "@/components/custom-card-form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import {
  useChatbotsControllerCreate,
  useChatbotsControllerFindOne,
  useChatbotsControllerUpdate,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import * as z from "zod";

const formSchema = z.object({
  access_scope: z.string(),
  description: z.string().min(1),
  documents: z.string(),
  llmBase: z.string(),
  name: z.string().min(1).max(255),
});

export default function ChatbotForm({ chatbotId }: { chatbotId?: string }) {
  const { toast } = useToast();
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

  const { mutateAsync: updateChatbot } = useChatbotsControllerUpdate({
    onError: (error) => {
      toast({
        description: error?.stack.msg,
        title: "Error updating chatbot",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "Your chatbot has been updated.",
        title: "Chatbot updated",
      });
    },
  });

  const { mutateAsync: createChatbot } = useChatbotsControllerCreate({
    onError: (error) => {
      toast({
        description: error?.stack.msg,
        title: "Error creating chatbot",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "Your chatbot has been created.",
        title: "Chatbot created",
      });
    },
  });

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
      defaultValue: chatbot?.llmBase,
      description:
        "This is the LLM base that will be used for this chatbot. Note that different models have different credit usages.",
      label: "Language Model",
      name: "llmBase",
      props: {
        placeholder: "Search llms...",
      },
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
    <CustomCardForm
      description={"Configure your chatbot's settings"}
      fields={formFields}
      onSubmit={
        chatbotId
          ? async (data) => {
              await updateChatbot({
                body: data,
                pathParams: {
                  chatbotId: chatbotId || "",
                  orgname: defaultOrgname,
                },
              });
            }
          : async (data) => {
              await createChatbot({
                body: data as any,
                pathParams: {
                  orgname: defaultOrgname,
                },
              });
            }
      }
      title="Configuration"
    />
  );
}
