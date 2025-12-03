import type { FormFieldConfig } from "@archesai/ui";
// import ImportCard from '@archesai/ui
import {
  GenericForm,
  Input,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
} from "@archesai/ui";
import { ARTIFACT_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { JSX } from "react";
import { useState } from "react";
import type { CreateArtifactBody, UpdateArtifactBody } from "#lib/index";
import {
  createArtifact,
  updateArtifact,
  useGetArtifactSuspense,
} from "#lib/index";

export default function ArtifactForm({ id }: { id?: string }): JSX.Element {
  const [tab, setTab] = useState<"file" | "text" | "url">("file");

  const { data: existingContentResponse, error } = useGetArtifactSuspense(id);

  if (error) {
    return <div>Content not found</div>;
  }
  const content = existingContentResponse.data;

  const formFields: FormFieldConfig[] = [
    {
      defaultValue: content.name ?? "",
      description: "This is the name that will be used for this content.",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Content name here..."
          type="text"
        />
      ),
    },
    {
      description:
        "Select the content you would like to run the tool on. You can select multiple content items.",
      label: "Input",
      name: tab === "file" ? "text" : tab,
      renderControl: (field) => (
        <Tabs value={tab}>
          <TabsList className="grid w-full grid-cols-3 px-1">
            <TabsTrigger
              onClick={() => {
                setTab("text");
              }}
              value="text"
            >
              Text
            </TabsTrigger>
            <TabsTrigger
              onClick={() => {
                setTab("file");
              }}
              value="file"
            >
              File
            </TabsTrigger>
            <TabsTrigger
              onClick={() => {
                setTab("url");
              }}
              value="url"
            >
              URL
            </TabsTrigger>
          </TabsList>
          <TabsContent value="text">
            <Textarea
              {...field}
              placeholder="Enter text here"
              value={field.value as string}
            />
          </TabsContent>
          <TabsContent value="url">
            <Textarea
              {...field}
              placeholder="Enter url here"
              rows={5}
              value={field.value as string}
            />
          </TabsContent>
          <TabsContent value="file">
            {/* <ImportCard
              cb={(content) => {
                field.onChange(content.map((c) => c.id))
              }}
            /> */}
          </TabsContent>
        </Tabs>
      ),
    },
  ];

  return (
    <GenericForm<CreateArtifactBody, UpdateArtifactBody>
      description={!id ? "Invite a new content" : "Update an existing content"}
      entityKey={ARTIFACT_ENTITY_KEY}
      fields={formFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createArtifactDto) => {
        await createArtifact(createArtifactDto);
      }}
      onSubmitUpdate={async (updateContentDto) => {
        await updateArtifact(id, updateContentDto);
      }}
      title="Configuration"
    />
  );
}
