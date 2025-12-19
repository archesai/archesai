import { useQuery } from "@tanstack/react-query";
import yaml from "js-yaml";
import type { JSX } from "react";
import { Badge } from "../../shadcn/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../../shadcn/card";
import { ScrollArea } from "../../shadcn/scroll-area";
import { Skeleton } from "../../shadcn/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../shadcn/tabs";
import { PathViewer } from "./path-viewer";
import { SchemaViewer } from "./schema-viewer";
import type {
  Callback,
  Example,
  Header,
  InternalMarker,
  OpenAPILink,
  OpenAPISpec,
  Parameter,
  PathItem,
  RequestBody,
  Response,
  SchemaObject,
} from "./types";
import { getInternalValue, getRefName, isInternal, isReference } from "./types";

interface OpenAPIPageProps {
  url: string;
}

async function fetchOpenAPISpec(url: string): Promise<OpenAPISpec> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to fetch OpenAPI spec: ${response.statusText}`);
  }
  const text = await response.text();
  const spec = yaml.load(text) as OpenAPISpec;
  return spec;
}

export function OpenAPIPage({ url }: OpenAPIPageProps): JSX.Element {
  const {
    data: spec,
    error,
    isLoading,
  } = useQuery({
    queryFn: () => fetchOpenAPISpec(url),
    queryKey: ["openapi-spec", url],
    staleTime: 5 * 60 * 1000,
  });

  if (isLoading) {
    return (
      <div className="space-y-4 p-4">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-4 w-96" />
        <div className="space-y-2">
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-32 w-full" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <Card className="m-4 border-destructive">
        <CardHeader>
          <CardTitle className="text-destructive">
            Error Loading OpenAPI Spec
          </CardTitle>
          <CardDescription>
            {error instanceof Error ? error.message : "Unknown error"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm">
            Make sure the OpenAPI spec is available at {url}
          </p>
        </CardContent>
      </Card>
    );
  }

  if (!spec) {
    return (
      <Card className="m-4 border-destructive">
        <CardHeader>
          <CardTitle className="text-destructive">
            Error Loading OpenAPI Spec
          </CardTitle>
          <CardDescription>No spec data received</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm">
            Make sure the OpenAPI spec is available at {url}
          </p>
        </CardContent>
      </Card>
    );
  }

  return <OpenAPIViewer spec={spec} />;
}

interface OpenAPIViewerProps {
  spec: OpenAPISpec;
}

const HTTP_METHODS = [
  "get",
  "post",
  "put",
  "patch",
  "delete",
  "options",
  "head",
] as const;

function sortInternalLast<T>(
  entries: [string, T][],
  isInternalFn: (entry: [string, T]) => boolean,
): [string, T][] {
  return [...entries].sort((a, b) => {
    const aInternal = isInternalFn(a);
    const bInternal = isInternalFn(b);
    if (aInternal === bInternal) return 0;
    return aInternal ? 1 : -1;
  });
}

function sortInternalMarkerLast<T extends InternalMarker>(
  entries: [string, T][],
): [string, T][] {
  return sortInternalLast(entries, ([, item]) => isInternal(item));
}

function isPathInternal(pathItem: PathItem): boolean {
  let hasAnyOperation = false;
  for (const method of HTTP_METHODS) {
    const operation = pathItem[method];
    if (operation) {
      hasAnyOperation = true;
      if (!isInternal(operation)) {
        return false;
      }
    }
  }
  return hasAnyOperation;
}

function groupPathsByTag(
  pathEntries: [string, NonNullable<OpenAPISpec["paths"]>[string]][],
  _tags?: OpenAPISpec["tags"],
): Record<string, [string, NonNullable<OpenAPISpec["paths"]>[string]][]> {
  const grouped: Record<
    string,
    [string, NonNullable<OpenAPISpec["paths"]>[string]][]
  > = {};

  for (const [path, pathItem] of pathEntries) {
    const methods = [
      "get",
      "post",
      "put",
      "patch",
      "delete",
      "options",
      "head",
    ];
    const operationTags = new Set<string>();

    for (const method of methods) {
      const op = pathItem[method as keyof typeof pathItem];
      if (op && typeof op === "object" && "tags" in op && op.tags) {
        for (const tag of op.tags) {
          operationTags.add(tag);
        }
      }
    }

    if (operationTags.size === 0) {
      operationTags.add("Untagged");
    }

    for (const tag of operationTags) {
      if (!grouped[tag]) {
        grouped[tag] = [];
      }
      grouped[tag].push([path, pathItem]);
    }
  }

  const sortedKeys = Object.keys(grouped).sort((a, b) => {
    if (a === "Untagged") return 1;
    if (b === "Untagged") return -1;
    return a.localeCompare(b);
  });

  const sorted: Record<
    string,
    [string, NonNullable<OpenAPISpec["paths"]>[string]][]
  > = {};
  for (const key of sortedKeys) {
    const group = grouped[key];
    if (group) {
      sorted[key] = group.sort((a, b) => {
        const aInternal = isPathInternal(a[1]);
        const bInternal = isPathInternal(b[1]);
        if (aInternal === bInternal) return 0;
        return aInternal ? 1 : -1;
      });
    }
  }

  return sorted;
}

export function OpenAPIViewer({ spec }: OpenAPIViewerProps): JSX.Element {
  const paths = spec.paths ?? {};
  const schemas = spec.components?.schemas ?? {};
  const securitySchemes = spec.components?.securitySchemes ?? {};
  const responses = spec.components?.responses ?? {};
  const requestBodies = spec.components?.requestBodies ?? {};
  const parameters = spec.components?.parameters ?? {};
  const headers = spec.components?.headers ?? {};
  const examples = spec.components?.examples ?? {};
  const links = spec.components?.links ?? {};
  const callbacks = spec.components?.callbacks ?? {};

  const pathEntries = sortInternalLast(Object.entries(paths), ([, item]) =>
    isPathInternal(item),
  );
  const schemaEntries = sortInternalMarkerLast(Object.entries(schemas));
  const securityEntries = Object.entries(securitySchemes);
  const responseEntries = sortInternalMarkerLast(Object.entries(responses));
  const requestBodyEntries = sortInternalMarkerLast(
    Object.entries(requestBodies),
  );
  const parameterEntries = sortInternalMarkerLast(Object.entries(parameters));
  const headerEntries = sortInternalMarkerLast(Object.entries(headers));
  const exampleEntries = sortInternalMarkerLast(Object.entries(examples));
  const linkEntries = sortInternalMarkerLast(Object.entries(links));
  const callbackEntries = Object.entries(callbacks);

  const pathsByTag = groupPathsByTag(pathEntries, spec.tags);

  return (
    <div className="flex h-full flex-col">
      <div className="border-b p-4">
        <div className="flex items-center gap-3">
          <h1 className="font-bold text-2xl">{spec.info.title}</h1>
          <Badge variant="secondary">v{spec.info.version}</Badge>
          <Badge variant="outline">OpenAPI {spec.openapi}</Badge>
        </div>
        {spec.info.description && (
          <p className="mt-2 text-muted-foreground">{spec.info.description}</p>
        )}
        {spec.servers && spec.servers.length > 0 && (
          <div className="mt-2 flex flex-wrap gap-2">
            {spec.servers.map((server) => (
              <Badge
                key={server.url}
                variant="outline"
              >
                {server.url}
              </Badge>
            ))}
          </div>
        )}
      </div>

      <Tabs
        className="flex flex-1 flex-col"
        defaultValue="paths"
      >
        <div className="overflow-x-auto border-b px-4">
          <TabsList className="w-max">
            <TabsTrigger value="paths">
              Paths ({pathEntries.length})
            </TabsTrigger>
            <TabsTrigger value="schemas">
              Schemas ({schemaEntries.length})
            </TabsTrigger>
            {responseEntries.length > 0 && (
              <TabsTrigger value="responses">
                Responses ({responseEntries.length})
              </TabsTrigger>
            )}
            {requestBodyEntries.length > 0 && (
              <TabsTrigger value="requestBodies">
                Request Bodies ({requestBodyEntries.length})
              </TabsTrigger>
            )}
            {parameterEntries.length > 0 && (
              <TabsTrigger value="parameters">
                Parameters ({parameterEntries.length})
              </TabsTrigger>
            )}
            {headerEntries.length > 0 && (
              <TabsTrigger value="headers">
                Headers ({headerEntries.length})
              </TabsTrigger>
            )}
            {exampleEntries.length > 0 && (
              <TabsTrigger value="examples">
                Examples ({exampleEntries.length})
              </TabsTrigger>
            )}
            {linkEntries.length > 0 && (
              <TabsTrigger value="links">
                Links ({linkEntries.length})
              </TabsTrigger>
            )}
            {callbackEntries.length > 0 && (
              <TabsTrigger value="callbacks">
                Callbacks ({callbackEntries.length})
              </TabsTrigger>
            )}
            {securityEntries.length > 0 && (
              <TabsTrigger value="security">
                Security ({securityEntries.length})
              </TabsTrigger>
            )}
          </TabsList>
        </div>

        <TabsContent
          className="m-0 flex-1"
          value="paths"
        >
          <ScrollArea className="h-[calc(100vh-220px)]">
            <div className="space-y-6 p-4">
              {Object.entries(pathsByTag).map(([tag, tagPaths]) => (
                <div key={tag}>
                  <h3 className="sticky top-0 mb-3 bg-background py-2 font-semibold text-lg">
                    {tag}
                    <Badge
                      className="ml-2"
                      variant="secondary"
                    >
                      {tagPaths.length}
                    </Badge>
                  </h3>
                  <div className="space-y-3">
                    {tagPaths.map(([path, pathItem]) => (
                      <PathViewer
                        key={path}
                        path={path}
                        pathItem={pathItem}
                        schemas={schemas}
                      />
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>
        </TabsContent>

        <TabsContent
          className="m-0 flex-1"
          value="schemas"
        >
          <ScrollArea className="h-[calc(100vh-220px)]">
            <div className="space-y-4 p-4">
              {schemaEntries.map(([name, schema]) => (
                <SchemaViewer
                  key={name}
                  name={name}
                  schema={schema}
                  schemas={schemas}
                />
              ))}
            </div>
          </ScrollArea>
        </TabsContent>

        {responseEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="responses"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {responseEntries.map(([name, response]) => (
                  <ResponseViewer
                    key={name}
                    name={name}
                    response={response}
                    schemas={schemas}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {requestBodyEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="requestBodies"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {requestBodyEntries.map(([name, requestBody]) => (
                  <RequestBodyViewer
                    key={name}
                    name={name}
                    requestBody={requestBody}
                    schemas={schemas}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {parameterEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="parameters"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {parameterEntries.map(([name, parameter]) => (
                  <ParameterViewer
                    key={name}
                    name={name}
                    parameter={parameter}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {headerEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="headers"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {headerEntries.map(([name, header]) => (
                  <HeaderViewer
                    header={header}
                    key={name}
                    name={name}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {exampleEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="examples"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {exampleEntries.map(([name, example]) => (
                  <ExampleViewer
                    example={example}
                    key={name}
                    name={name}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {linkEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="links"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {linkEntries.map(([name, link]) => (
                  <LinkViewer
                    key={name}
                    link={link}
                    name={name}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {callbackEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="callbacks"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {callbackEntries.map(([name, callback]) => (
                  <CallbackViewer
                    callback={callback}
                    key={name}
                    name={name}
                    schemas={schemas}
                  />
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}

        {securityEntries.length > 0 && (
          <TabsContent
            className="m-0 flex-1"
            value="security"
          >
            <ScrollArea className="h-[calc(100vh-220px)]">
              <div className="space-y-4 p-4">
                {securityEntries.map(([name, scheme]) => (
                  <div
                    className="rounded-lg border p-4"
                    key={name}
                  >
                    <div className="mb-2 flex items-center gap-2">
                      <h4 className="font-semibold text-base">{name}</h4>
                      <Badge variant="secondary">{scheme.type}</Badge>
                    </div>
                    <dl className="grid grid-cols-2 gap-2 text-sm">
                      {scheme.description && (
                        <>
                          <dt className="text-muted-foreground">Description</dt>
                          <dd>{scheme.description}</dd>
                        </>
                      )}
                      {scheme.scheme && (
                        <>
                          <dt className="text-muted-foreground">Scheme</dt>
                          <dd>{scheme.scheme}</dd>
                        </>
                      )}
                      {scheme.bearerFormat && (
                        <>
                          <dt className="text-muted-foreground">
                            Bearer Format
                          </dt>
                          <dd>{scheme.bearerFormat}</dd>
                        </>
                      )}
                      {scheme.in && (
                        <>
                          <dt className="text-muted-foreground">In</dt>
                          <dd>{scheme.in}</dd>
                        </>
                      )}
                      {scheme.name && (
                        <>
                          <dt className="text-muted-foreground">Name</dt>
                          <dd>{scheme.name}</dd>
                        </>
                      )}
                      {scheme.flows && (
                        <>
                          <dt className="text-muted-foreground">Flows</dt>
                          <dd>
                            {Object.keys(scheme.flows).map((flow) => (
                              <Badge
                                className="mr-1"
                                key={flow}
                                variant="outline"
                              >
                                {flow}
                              </Badge>
                            ))}
                          </dd>
                        </>
                      )}
                    </dl>
                  </div>
                ))}
              </div>
            </ScrollArea>
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}

// Response Viewer Component
interface ResponseViewerProps {
  name: string;
  response: Response;
  schemas: Record<string, SchemaObject>;
}

function ResponseViewer({ name, response }: ResponseViewerProps): JSX.Element {
  const internalValue = getInternalValue(response);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
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
      </div>
      {response.description && (
        <p className="mb-3 text-muted-foreground text-sm">
          {response.description}
        </p>
      )}
      <div className="space-y-4">
        {response.content && Object.keys(response.content).length > 0 && (
          <div>
            <h5 className="mb-2 font-medium text-sm">Content</h5>
            <div className="space-y-2">
              {Object.entries(response.content).map(
                ([contentType, mediaType]) => (
                  <div
                    className="flex items-center gap-2 text-sm"
                    key={contentType}
                  >
                    <Badge variant="outline">{contentType}</Badge>
                    {mediaType.schema &&
                      (isReference(mediaType.schema) ? (
                        <Badge variant="secondary">
                          {getRefName(mediaType.schema.$ref)}
                        </Badge>
                      ) : (
                        <Badge variant="secondary">
                          {mediaType.schema.type ?? "object"}
                        </Badge>
                      ))}
                  </div>
                ),
              )}
            </div>
          </div>
        )}
        {response.headers && Object.keys(response.headers).length > 0 && (
          <div>
            <h5 className="mb-2 font-medium text-sm">Headers</h5>
            <div className="space-y-2">
              {Object.entries(response.headers).map(([headerName, header]) => (
                <div
                  className="flex items-start gap-2 text-sm"
                  key={headerName}
                >
                  <code className="rounded bg-muted px-1 font-mono">
                    {headerName}
                  </code>
                  {header.required && (
                    <Badge
                      className="text-xs"
                      variant="default"
                    >
                      required
                    </Badge>
                  )}
                  {header.schema && (
                    <Badge
                      className="text-xs"
                      variant="secondary"
                    >
                      {isReference(header.schema)
                        ? getRefName(header.schema.$ref)
                        : (header.schema.type ?? "unknown")}
                    </Badge>
                  )}
                  {header.description && (
                    <span className="text-muted-foreground">
                      {header.description}
                    </span>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

// Request Body Viewer Component
interface RequestBodyViewerProps {
  name: string;
  requestBody: RequestBody;
  schemas: Record<string, SchemaObject>;
}

function RequestBodyViewer({
  name,
  requestBody,
}: RequestBodyViewerProps): JSX.Element {
  const internalValue = getInternalValue(requestBody);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
        {requestBody.required && (
          <Badge
            className="text-xs"
            variant="default"
          >
            required
          </Badge>
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
      </div>
      {requestBody.description && (
        <p className="mb-3 text-muted-foreground text-sm">
          {requestBody.description}
        </p>
      )}
      {requestBody.content && Object.keys(requestBody.content).length > 0 && (
        <div className="space-y-2">
          {Object.entries(requestBody.content).map(
            ([contentType, mediaType]) => (
              <div
                className="flex items-center gap-2 text-sm"
                key={contentType}
              >
                <Badge variant="outline">{contentType}</Badge>
                {mediaType.schema &&
                  (isReference(mediaType.schema) ? (
                    <Badge variant="secondary">
                      {getRefName(mediaType.schema.$ref)}
                    </Badge>
                  ) : (
                    <Badge variant="secondary">
                      {mediaType.schema.type ?? "object"}
                    </Badge>
                  ))}
              </div>
            ),
          )}
        </div>
      )}
    </div>
  );
}

// Parameter Viewer Component
interface ParameterViewerProps {
  name: string;
  parameter: Parameter;
}

function ParameterViewer({
  name,
  parameter,
}: ParameterViewerProps): JSX.Element {
  const internalValue = getInternalValue(parameter);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
        <Badge variant="outline">{parameter.in}</Badge>
        {parameter.required && (
          <Badge
            className="text-xs"
            variant="default"
          >
            required
          </Badge>
        )}
        {parameter.deprecated && (
          <Badge
            className="text-xs"
            variant="destructive"
          >
            deprecated
          </Badge>
        )}
        {parameter.schema &&
          (isReference(parameter.schema) ? (
            <Badge variant="secondary">
              {getRefName(parameter.schema.$ref)}
            </Badge>
          ) : (
            <>
              <Badge variant="secondary">
                {parameter.schema.type ?? "unknown"}
              </Badge>
              {parameter.schema.format && (
                <Badge variant="outline">{parameter.schema.format}</Badge>
              )}
            </>
          ))}
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
      </div>
      {parameter.description && (
        <p className="mb-3 text-muted-foreground text-sm">
          {parameter.description}
        </p>
      )}
      {parameter.schema &&
        !isReference(parameter.schema) &&
        parameter.schema.enum && (
          <div className="text-sm">
            <span className="text-muted-foreground">Enum: </span>
            {parameter.schema.enum.map((v) => (
              <code
                className="mx-1 rounded bg-muted px-1"
                key={String(v)}
              >
                {String(v)}
              </code>
            ))}
          </div>
        )}
    </div>
  );
}

// Header Viewer Component
interface HeaderViewerProps {
  name: string;
  header: Header;
}

function HeaderViewer({ name, header }: HeaderViewerProps): JSX.Element {
  const internalValue = getInternalValue(header);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
        {header.required && (
          <Badge
            className="text-xs"
            variant="default"
          >
            required
          </Badge>
        )}
        {header.deprecated && (
          <Badge
            className="text-xs"
            variant="destructive"
          >
            deprecated
          </Badge>
        )}
        {header.schema &&
          (isReference(header.schema) ? (
            <Badge variant="secondary">{getRefName(header.schema.$ref)}</Badge>
          ) : (
            <Badge variant="secondary">{header.schema.type ?? "unknown"}</Badge>
          ))}
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
      </div>
      {header.description && (
        <p className="mb-3 text-muted-foreground text-sm">
          {header.description}
        </p>
      )}
      {header.example !== undefined && (
        <div className="text-sm">
          <span className="text-muted-foreground">Example: </span>
          <code className="rounded bg-muted px-1">
            {JSON.stringify(header.example)}
          </code>
        </div>
      )}
    </div>
  );
}

// Example Viewer Component
interface ExampleViewerProps {
  name: string;
  example: Example;
}

function ExampleViewer({ name, example }: ExampleViewerProps): JSX.Element {
  const internalValue = getInternalValue(example);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
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
      </div>
      {example.summary && (
        <p className="mb-3 text-muted-foreground text-sm">{example.summary}</p>
      )}
      {example.description && (
        <p className="mb-3 text-muted-foreground text-sm">
          {example.description}
        </p>
      )}
      {example.value !== undefined && (
        <div>
          <h5 className="mb-1 font-medium text-sm">Value</h5>
          <pre className="overflow-auto rounded bg-muted p-2 text-xs">
            {JSON.stringify(example.value, null, 2)}
          </pre>
        </div>
      )}
      {example.externalValue && (
        <div className="text-sm">
          <span className="text-muted-foreground">External Value: </span>
          <code className="rounded bg-muted px-1">{example.externalValue}</code>
        </div>
      )}
    </div>
  );
}

// Link Viewer Component
interface LinkViewerProps {
  name: string;
  link: OpenAPILink;
}

function LinkViewer({ name, link }: LinkViewerProps): JSX.Element {
  const internalValue = getInternalValue(link);

  return (
    <div
      className={`rounded-lg border p-4 ${internalValue ? "bg-sky-500/5 opacity-60" : ""}`}
    >
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
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
      </div>
      {link.description && (
        <p className="mb-3 text-muted-foreground text-sm">{link.description}</p>
      )}
      <dl className="grid grid-cols-2 gap-2 text-sm">
        {link.operationId && (
          <>
            <dt className="text-muted-foreground">Operation ID</dt>
            <dd>
              <code className="rounded bg-muted px-1">{link.operationId}</code>
            </dd>
          </>
        )}
        {link.operationRef && (
          <>
            <dt className="text-muted-foreground">Operation Ref</dt>
            <dd>
              <code className="rounded bg-muted px-1">{link.operationRef}</code>
            </dd>
          </>
        )}
        {link.parameters && Object.keys(link.parameters).length > 0 && (
          <>
            <dt className="text-muted-foreground">Parameters</dt>
            <dd>
              <pre className="overflow-auto rounded bg-muted p-1 text-xs">
                {JSON.stringify(link.parameters, null, 2)}
              </pre>
            </dd>
          </>
        )}
        {link.server && (
          <>
            <dt className="text-muted-foreground">Server</dt>
            <dd>
              <code className="rounded bg-muted px-1">{link.server.url}</code>
            </dd>
          </>
        )}
      </dl>
    </div>
  );
}

// Callback Viewer Component
interface CallbackViewerProps {
  name: string;
  callback: Callback;
  schemas: Record<string, SchemaObject>;
}

function CallbackViewer({
  name,
  callback,
  schemas,
}: CallbackViewerProps): JSX.Element {
  return (
    <div className="rounded-lg border p-4">
      <div className="mb-2 flex items-center gap-2">
        <h4 className="font-semibold text-base">{name}</h4>
      </div>
      <div className="space-y-4">
        {Object.entries(callback).map(([expression, pathItem]) => (
          <div key={expression}>
            <h5 className="mb-2 font-mono text-sm">{expression}</h5>
            <PathViewer
              path={expression}
              pathItem={pathItem}
              schemas={schemas}
            />
          </div>
        ))}
      </div>
    </div>
  );
}
