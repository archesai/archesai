"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import {
  useContentControllerCreate,
  useContentControllerFindOne,
  useContentControllerUpdate,
} from "@/generated/archesApiComponents";
import {
  CreateContentDto,
  UpdateContentDto,
} from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import * as z from "zod";

const formSchema = z.object({
  name: z.string(),
});

export default function ContentForm({ contentId }: { contentId?: string }) {
  const { defaultOrgname } = useAuth();
  const { data: content } = useContentControllerFindOne(
    {
      pathParams: {
        contentId: contentId as string,
        orgname: defaultOrgname,
      },
    },
    {
      enabled: !!defaultOrgname && !!contentId,
    }
  );
  const { mutateAsync: updateContent } = useContentControllerUpdate({});
  const { mutateAsync: createContent } = useContentControllerCreate({});

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: content?.name,
      description: "This is the name that will be used for this content.",
      label: "Name",
      name: "name",
      props: {
        placeholder: "Content name here...",
      },
      validationRule: formSchema.shape.name,
    },
  ];

  return (
    <GenericForm<CreateContentDto, UpdateContentDto>
      description={
        !contentId ? "Invite a new content" : "Update an existing content"
      }
      fields={formFields}
      isUpdateForm={!!contentId}
      itemType="content"
      onSubmitCreate={async (createContentDto, mutateOptions) => {
        await createContent(
          {
            body: createContentDto,
            pathParams: {
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      onSubmitUpdate={async (data, mutateOptions) => {
        await updateContent(
          {
            body: data as any,
            pathParams: {
              contentId: contentId as string,
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
