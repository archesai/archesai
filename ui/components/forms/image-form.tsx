"use client";
import { FormFieldConfig, GenericForm } from "@/components/generic-form";
import { Input } from "@/components/ui/input";
import {
  useContentControllerFindOne,
  useImagesControllerCreate,
} from "@/generated/archesApiComponents";
import { CreateImageDto } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import * as z from "zod";

import { Textarea } from "../ui/textarea";

const formSchema = z.object({
  name: z.optional(z.string().min(1).max(255)),
  prompt: z.string(),
});

export default function ImageForm({ imageId }: { imageId?: string }) {
  const { defaultOrgname } = useAuth();
  const { data: image, isLoading } = useContentControllerFindOne(
    {
      pathParams: {
        contentId: imageId as string,
        orgname: defaultOrgname,
      },
    },
    {
      enabled: !!defaultOrgname && !!imageId,
    }
  );
  const { mutateAsync: createChatbot } = useImagesControllerCreate({});

  const formFields: FormFieldConfig[] = [
    {
      component: Input,
      defaultValue: image?.name,
      description: "This is the name that will be used for this image.",
      ignoreOnCreate: true,
      label: "Name",
      name: "name",
      props: {
        placeholder: "Image name here...",
      },
      validationRule: formSchema.shape.name,
    },
    {
      component: Textarea,
      defaultValue: image?.buildArgs.prompt,
      description: "This is the prompt that will be used for this image.",
      label: "Prompt",
      name: "prompt",
      props: {
        placeholder: "Image prompt here...",
      },
      validationRule: formSchema.shape.prompt,
    },
  ];

  if (isLoading) {
    return null;
  }

  return (
    <GenericForm<CreateImageDto, undefined>
      description={"Create a new image."}
      fields={formFields}
      isUpdateForm={!!imageId}
      itemType="image"
      onSubmitCreate={async (createImageDto, mutateOptions) => {
        await createChatbot(
          {
            body: {
              height: 1024,
              name: createImageDto.prompt,
              prompt: createImageDto.prompt,
              width: 1024,
            },
            pathParams: {
              orgname: defaultOrgname,
            },
          },
          mutateOptions
        );
      }}
      title="Create Image"
    />
  );
}
