import { useState } from "react";

import { PlusSquareIcon, TrashIcon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import { Input } from "#components/shadcn/input";
import { Label } from "#components/shadcn/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#components/shadcn/select";

interface FieldDefinition {
  constraints?: {
    elementType?: string;
    fieldName?: string;
  };
  fieldName: string;
  fieldType: string;
  isOptional: boolean;
  subFields?: FieldDefinition[];
}

const fieldTypes = [
  { label: "Text", value: "string" },
  { label: "Number", value: "number" },
  { label: "True/False", value: "boolean" },
  { label: "List", value: "array" },
  { label: "Sub-Item", value: "object" },
];

const SchemaBuilder: React.FC = () => {
  const [fields, setFields] = useState<FieldDefinition[]>([]);

  const handleSubmit = () => {
    generateJsonSchema(fields);
    // Send schemaString to backend
    // For example, use fetch to send schemaString to your backend
  };

  return (
    <div className="flex flex-wrap">
      <div className="w-full">
        <FieldList
          fields={fields}
          setFields={setFields}
        />
        <div className="flex w-full justify-end">
          <Button
            className="mt-4"
            onClick={handleSubmit}
          >
            Submit Schema
          </Button>
        </div>
      </div>
    </div>
  );
};

interface FieldListProps {
  fields: FieldDefinition[];
  setFields: React.Dispatch<React.SetStateAction<FieldDefinition[]>>;
}

const FieldList: React.FC<FieldListProps> = ({ fields, setFields }) => {
  const addField = () => {
    setFields([
      ...fields,
      {
        constraints: {},
        fieldName: "",
        fieldType: "string",
        isOptional: false,
        subFields: [],
      },
    ]);
  };

  return (
    <div className="flex flex-col gap-4">
      {fields.map((field, index) => (
        <FieldEditor
          field={field}
          fields={fields}
          index={index}
          key={field.fieldName}
          setFields={setFields}
        />
      ))}
      <Button
        onClick={addField}
        variant="outline"
      >
        <PlusSquareIcon className="mr-2 h-4 w-4" />
        Add Field
      </Button>
    </div>
  );
};

interface FieldEditorProps {
  field: FieldDefinition;
  fields: FieldDefinition[];
  index: number;
  setFields: React.Dispatch<React.SetStateAction<FieldDefinition[]>>;
}

const FieldEditor: React.FC<FieldEditorProps> = ({
  field,
  fields,
  index,
  setFields,
}) => {
  const handleFieldChange = (newField: FieldDefinition) => {
    const newFields = [...fields];
    newFields[index] = newField;
    setFields(newFields);
  };

  const removeField = () => {
    const newFields = [...fields];
    newFields.splice(index, 1);
    setFields(newFields);
  };

  return (
    <div>
      <div className="flex items-center gap-2">
        <Input
          onChange={(e) => {
            handleFieldChange({ ...field, fieldName: e.target.value });
          }}
          placeholder="Field Name"
          value={field.fieldName}
        />
        <Select
          onValueChange={(value) => {
            handleFieldChange({ ...field, fieldType: value });
          }}
          value={field.fieldType}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select type" />
          </SelectTrigger>
          <SelectContent>
            {fieldTypes.map((type) => (
              <SelectItem
                key={type.value}
                value={type.value}
              >
                {type.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {/* Optional Checkbox if needed */}
        {/* Uncomment if you decide to use the optional field
        <div className="flex items-center gap-1">
          <Checkbox
            checked={field.isOptional}
            onCheckedChange={(checked) =>
              handleFieldChange({ ...field, isOptional: !!checked })
            }
          />
          <span>Optional</span>
        </div>
        */}

        <Button
          className="h-8"
          onClick={removeField}
          variant="ghost"
        >
          <TrashIcon className="h-4 w-4 text-destructive" />
        </Button>
      </div>

      {/* Additional constraints based on type */}
      {field.fieldType === "object" && (
        <div className="mt-4 ml-4">
          <Label>Sub Item Fields</Label>
          <FieldList
            fields={field.subFields ?? []}
            setFields={(subFields) => {
              handleFieldChange({
                ...field,
                subFields: subFields as FieldDefinition[],
              });
            }}
          />
        </div>
      )}

      {field.fieldType === "array" && (
        <div className="mt-4 ml-4">
          <Label>List Item Type</Label>
          <Select
            onValueChange={(value) => {
              handleFieldChange({
                ...field,
                constraints: { ...field.constraints, elementType: value },
              });
            }}
            value={field.constraints?.elementType ?? "string"}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select element type" />
            </SelectTrigger>
            <SelectContent>
              {fieldTypes
                .filter((type) => type.value !== "array") // Avoid nesting arrays
                .map((type) => (
                  <SelectItem
                    key={type.value}
                    value={type.value}
                  >
                    {type.label}
                  </SelectItem>
                ))}
            </SelectContent>
          </Select>
        </div>
      )}
    </div>
  );
};

const generateJsonSchema = (fields: FieldDefinition[]): string => {
  let schemaString = "z.object({\n";
  fields.forEach((field) => {
    schemaString += generateFieldSchema(field, 1);
  });
  schemaString += "})";
  return schemaString;
};

const generateFieldSchema = (
  field: FieldDefinition,
  indentLevel: number,
): string => {
  const indent = "  ".repeat(indentLevel);
  let fieldString = `${indent}${field.fieldName}: `;
  let fieldSchema = "";

  switch (field.fieldType) {
    case "array":
      fieldSchema = `z.array(${field.constraints?.elementType ? `z.${field.constraints.elementType}()` : "z.any()"})`;
      break;
    case "boolean":
      fieldSchema = "z.boolean()";
      break;
    case "number":
      fieldSchema = "z.number()";
      break;
    case "object":
      fieldSchema = "z.object({\n";
      field.subFields?.forEach((subField) => {
        fieldSchema += generateFieldSchema(subField, indentLevel + 1);
      });
      fieldSchema += `${indent}})`;
      break;
    case "string":
      fieldSchema = "z.string()";
      break;
    default:
      fieldSchema = "z.any()";
  }

  if (field.isOptional) {
    fieldSchema += ".optional()";
  }

  fieldString += `${fieldSchema},\n`;
  return fieldString;
};

// Function to generate example JSON
const generateExampleJSON = (fields: FieldDefinition[]): unknown => {
  const obj: Record<string, unknown> = {};
  fields.forEach((field) => {
    if (field.fieldName) {
      if (!(field.fieldName in obj)) {
        return;
      }
      obj[field.fieldName] = generateFieldExample(field);
    }
  });
  return obj;
};

const generateFieldExample = (field: FieldDefinition): unknown => {
  if (field.isOptional) return undefined;

  switch (field.fieldType) {
    case "array":
      return [
        generateFieldExample({
          fieldName: "",
          fieldType: field.constraints?.elementType ?? "string",
          isOptional: false,
        }),
      ];
    case "boolean":
      return true;
    case "number":
      return 123;
    case "object":
      return generateExampleJSON(field.subFields ?? []);
    case "string":
      return "example text";
    default:
      return null;
  }
};

export default SchemaBuilder;
