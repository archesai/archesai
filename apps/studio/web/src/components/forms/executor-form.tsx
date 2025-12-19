import type { FormFieldConfig } from "@archesai/ui";
import {
  Checkbox,
  FormControl,
  GenericForm,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Textarea,
} from "@archesai/ui";
import type { JSX } from "react";
import type { CreateExecutorBody, UpdateExecutorBody } from "#lib/index";
import {
  useCreateExecutor,
  useGetExecutor,
  useUpdateExecutor,
} from "#lib/index";

export default function ExecutorForm({ id }: { id?: string }): JSX.Element {
  const { mutateAsync: updateExecutor } = useUpdateExecutor();
  const { mutateAsync: createExecutor } = useCreateExecutor();
  const { data: existingExecutor } = useGetExecutor(id, {
    query: { enabled: !!id },
  });
  // Cast to Record to allow accessing request body fields that may differ from entity fields
  const data = existingExecutor?.data as Record<string, unknown> | undefined;
  const createFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.cpuShares as string) ?? "",
      description: "CPU shares (relative weight)",
      label: "CPU Shares",
      name: "cpuShares",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter cpu shares..."
          type="text"
        />
      ),
    },
    {
      defaultValue: (data?.dependencies as string) ?? "",
      description: "Dependencies configuration",
      label: "Dependencies",
      name: "dependencies",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter dependencies..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.description as string) ?? "",
      description: "The executor description",
      label: "Description",
      name: "description",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter description..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.env as string) ?? "",
      description: "Environment variables as JSON array",
      label: "Env",
      name: "env",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter env..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.executeCode as string) ?? "",
      description: "The custom execute function code",
      label: "Execute Code",
      name: "executeCode",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter execute code..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.extraFiles as string) ?? "",
      description: "Additional files to mount",
      label: "Extra Files",
      name: "extraFiles",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter extra files..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.language as string) ?? "",
      description: "The programming language for the executor",
      label: "Language",
      name: "language",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select language..." />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="nodejs">Nodejs</SelectItem>
            <SelectItem value="python">Python</SelectItem>
            <SelectItem value="go">Go</SelectItem>
          </SelectContent>
        </Select>
      ),
    },
    {
      defaultValue: (data?.memoryMB as string) ?? "",
      description: "Memory limit in megabytes",
      label: "Memory MB",
      name: "memoryMB",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter memory mb..."
          type="text"
        />
      ),
    },
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the executor",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter name..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.schemaIn as string) ?? "",
      description: "Input JSON Schema",
      label: "Schema In",
      name: "schemaIn",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter schema in..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.schemaOut as string) ?? "",
      description: "Output JSON Schema",
      label: "Schema Out",
      name: "schemaOut",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter schema out..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.timeout as string) ?? "",
      description: "Execution timeout in seconds",
      label: "Timeout",
      name: "timeout",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter timeout..."
          type="text"
        />
      ),
    },
  ];
  const updateFormFields: FormFieldConfig[] = [
    {
      defaultValue: (data?.cpuShares as string) ?? "",
      description: "CPU shares (relative weight)",
      label: "CPU Shares",
      name: "cpuShares",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter cpu shares..."
          type="text"
        />
      ),
    },
    {
      defaultValue: (data?.dependencies as string) ?? "",
      description: "Dependencies configuration",
      label: "Dependencies",
      name: "dependencies",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter dependencies..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.description as string) ?? "",
      description: "The executor description",
      label: "Description",
      name: "description",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter description..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.env as string) ?? "",
      description: "Environment variables as JSON array",
      label: "Env",
      name: "env",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter env..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.executeCode as string) ?? "",
      description: "The custom execute function code",
      label: "Execute Code",
      name: "executeCode",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter execute code..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.extraFiles as string) ?? "",
      description: "Additional files to mount",
      label: "Extra Files",
      name: "extraFiles",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter extra files..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.isActive as boolean) ?? false,
      description: "Whether the executor is active",
      label: "Is Active",
      name: "isActive",
      renderControl: (field) => (
        <Checkbox
          checked={field.value as boolean}
          onCheckedChange={field.onChange}
        />
      ),
    },
    {
      defaultValue: (data?.language as string) ?? "",
      description: "The programming language for the executor",
      label: "Language",
      name: "language",
      renderControl: (field) => (
        <Select
          defaultValue={field.value as string}
          onValueChange={field.onChange}
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue placeholder="Select language..." />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="nodejs">Nodejs</SelectItem>
            <SelectItem value="python">Python</SelectItem>
            <SelectItem value="go">Go</SelectItem>
          </SelectContent>
        </Select>
      ),
    },
    {
      defaultValue: (data?.memoryMB as string) ?? "",
      description: "Memory limit in megabytes",
      label: "Memory MB",
      name: "memoryMB",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter memory mb..."
          type="text"
        />
      ),
    },
    {
      defaultValue: (data?.name as string) ?? "",
      description: "The name of the executor",
      label: "Name",
      name: "name",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter name..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.schemaIn as string) ?? "",
      description: "Input JSON Schema",
      label: "Schema In",
      name: "schemaIn",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter schema in..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.schemaOut as string) ?? "",
      description: "Output JSON Schema",
      label: "Schema Out",
      name: "schemaOut",
      renderControl: (field) => (
        <Textarea
          {...field}
          placeholder="Enter schema out..."
          rows={5}
          value={field.value as string}
        />
      ),
    },
    {
      defaultValue: (data?.timeout as string) ?? "",
      description: "Execution timeout in seconds",
      label: "Timeout",
      name: "timeout",
      renderControl: (field) => (
        <Input
          {...field}
          placeholder="Enter timeout..."
          type="text"
        />
      ),
    },
  ];
  return (
    <GenericForm<CreateExecutorBody, UpdateExecutorBody>
      description={
        !id ? "Create a new executor" : "Update an existing executor"
      }
      entityKey="executors"
      fields={!id ? createFormFields : updateFormFields}
      isUpdateForm={!!id}
      onSubmitCreate={async (createDto) => {
        await createExecutor({ data: createDto });
      }}
      onSubmitUpdate={async (updateDto) => {
        await updateExecutor({ data: updateDto, id: id });
      }}
      title="Executor"
    />
  );
}
