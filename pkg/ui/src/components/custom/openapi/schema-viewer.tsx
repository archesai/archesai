import type { JSX } from "react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "../../shadcn/accordion";
import { Badge } from "../../shadcn/badge";
import type { Reference, SchemaObject } from "./types";
import { getInternalValue, getRefName, isReference } from "./types";

interface SchemaViewerProps {
  name: string;
  schema: SchemaObject;
  schemas?: Record<string, SchemaObject> | undefined;
}

export function SchemaViewer({
  name,
  schema,
  schemas,
}: SchemaViewerProps): JSX.Element {
  const internalValue = getInternalValue(schema);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
        {schema.type && <Badge variant="secondary">{schema.type}</Badge>}
        {schema["x-codegen-schema-type"] && (
          <Badge variant="outline">{schema["x-codegen-schema-type"]}</Badge>
        )}
        {internalValue && (
          <>
            <Badge className="bg-sky-500/10 text-sky-500 text-xs">
              internal
            </Badge>
            <Badge className="bg-sky-500/10 text-sky-500 text-xs">
              {internalValue}
            </Badge>
          </>
        )}
        {schema.deprecated && <Badge variant="destructive">deprecated</Badge>}
      </div>
      {schema.description && (
        <p className="mb-3 text-muted-foreground text-sm">
          {schema.description}
        </p>
      )}
      <SchemaProperties
        schema={schema}
        schemas={schemas}
      />
    </div>
  );
}

interface SchemaPropertiesProps {
  schema: SchemaObject;
  schemas?: Record<string, SchemaObject> | undefined;
  depth?: number;
}

function SchemaProperties({
  schema,
  schemas,
  depth = 0,
}: SchemaPropertiesProps): JSX.Element | null {
  if (!schema.properties) {
    if (schema.type === "array" && schema.items) {
      return (
        <div className="border-l pl-4">
          <span className="text-muted-foreground text-sm">Array items: </span>
          {isReference(schema.items) ? (
            <Badge variant="outline">{getRefName(schema.items.$ref)}</Badge>
          ) : (
            <Badge variant="secondary">{schema.items.type ?? "unknown"}</Badge>
          )}
        </div>
      );
    }
    return null;
  }

  const requiredFields = new Set(schema.required ?? []);

  return (
    <div className={depth > 0 ? "border-l pl-4" : ""}>
      <div className="space-y-2">
        {Object.entries(schema.properties).map(([propName, propSchema]) => (
          <PropertyItem
            depth={depth}
            isRequired={requiredFields.has(propName)}
            key={propName}
            name={propName}
            schema={propSchema}
            schemas={schemas}
          />
        ))}
      </div>
    </div>
  );
}

interface PropertyItemProps {
  name: string;
  schema: SchemaObject | Reference;
  schemas?: Record<string, SchemaObject> | undefined;
  isRequired: boolean;
  depth: number;
}

function PropertyItem({
  name,
  schema,
  schemas,
  isRequired,
  depth,
}: PropertyItemProps): JSX.Element {
  if (isReference(schema)) {
    const refName = getRefName(schema.$ref);
    const referencedSchema = schemas?.[refName];

    return (
      <div className="py-1">
        <div className="flex items-center gap-2">
          <code className="font-mono text-sm">{name}</code>
          <Badge variant="outline">{refName}</Badge>
          {isRequired && (
            <Badge
              className="text-xs"
              variant="default"
            >
              required
            </Badge>
          )}
        </div>
        {referencedSchema && depth < 2 && (
          <Accordion
            className="mt-2"
            collapsible
            type="single"
          >
            <AccordionItem value="expand">
              <AccordionTrigger className="py-1 text-muted-foreground text-xs">
                Expand {refName}
              </AccordionTrigger>
              <AccordionContent>
                <SchemaProperties
                  depth={depth + 1}
                  schema={referencedSchema}
                  schemas={schemas}
                />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        )}
      </div>
    );
  }

  const hasNestedProperties =
    schema.properties ?? (schema.type === "array" && schema.items);

  return (
    <div className="py-1">
      <div className="flex flex-wrap items-center gap-2">
        <code className="font-mono text-sm">{name}</code>
        <Badge variant="secondary">{getTypeDisplay(schema)}</Badge>
        {isRequired && (
          <Badge
            className="text-xs"
            variant="default"
          >
            required
          </Badge>
        )}
        {schema.deprecated && (
          <Badge
            className="text-xs"
            variant="destructive"
          >
            deprecated
          </Badge>
        )}
        {schema.readOnly && (
          <Badge
            className="text-xs"
            variant="outline"
          >
            readOnly
          </Badge>
        )}
        {schema.format && (
          <Badge
            className="text-xs"
            variant="outline"
          >
            {schema.format}
          </Badge>
        )}
      </div>
      {schema.description && (
        <p className="mt-1 text-muted-foreground text-xs">
          {schema.description}
        </p>
      )}
      {schema.enum && (
        <div className="mt-1 text-xs">
          <span className="text-muted-foreground">Enum: </span>
          {schema.enum.map((v) => (
            <code
              className="mx-1 rounded bg-muted px-1"
              key={String(v)}
            >
              {String(v)}
            </code>
          ))}
        </div>
      )}
      {hasNestedProperties && depth < 2 && (
        <div className="mt-2">
          <SchemaProperties
            depth={depth + 1}
            schema={schema}
            schemas={schemas}
          />
        </div>
      )}
    </div>
  );
}

function getTypeDisplay(schema: SchemaObject): string {
  if (schema.type === "array" && schema.items) {
    if (isReference(schema.items)) {
      return `${getRefName(schema.items.$ref)}[]`;
    }
    return `${schema.items.type ?? "unknown"}[]`;
  }
  return schema.type ?? "object";
}
