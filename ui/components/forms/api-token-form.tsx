"use client";
import { CustomCardForm, FormFieldConfig } from "@/components/custom-card-form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import {
  useApiTokensControllerCreate,
  useApiTokensControllerFindOne,
  useApiTokensControllerUpdate,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import * as z from "zod";

const formSchema = z.object({
  name: z.string().min(1).max(255),
  role: z.enum(["ADMIN", "USER"]),
});

export default function ApiTokenForm({ apiTokenId }: { apiTokenId?: string }) {
  const { toast } = useToast();
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

  const { mutateAsync: updateApiToken } = useApiTokensControllerUpdate({
    onError: (error) => {
      toast({
        description: error?.stack.msg,
        title: "Error updating API Token",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "Your API Token has been updated.",
        title: "API Token updated",
      });
    },
  });

  const { mutateAsync: createApiToken } = useApiTokensControllerCreate({
    onError: (error) => {
      toast({
        description: error?.stack.msg,
        title: "Error creating API Token",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "Your API Token has been created.",
        title: "API Token created",
      });
    },
  });

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: apiToken?.name,
      description: "This is the name that will be used for this API Token.",
      label: "Name",
      name: "name",
      props: {
        placeholder: "API Token name here...",
      },
      validationRule: formSchema.shape.name,
    },
    {
      component: Input,
      defaultValue: apiToken?.role,
      description:
        "This is the LLM base that will be used for this API Token. Note that different models have different credit usages.",
      label: "Language Model",
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
    <CustomCardForm
      description={"Configure your API Token's settings"}
      fields={formFields}
      onSubmit={
        apiTokenId
          ? async () => {
              await updateApiToken({
                pathParams: {
                  id: apiTokenId || "",
                  orgname: defaultOrgname,
                },
              });
            }
          : async (data) => {
              await createApiToken({
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
