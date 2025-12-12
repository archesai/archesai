export interface OpenAPISpec {
  openapi: string;
  info: OpenAPIInfo;
  paths?: Record<string, PathItem>;
  components?: OpenAPIComponents;
  tags?: OpenAPITag[];
  servers?: Server[];
}

export interface OpenAPIInfo {
  title: string;
  version: string;
  description?: string;
}

export interface OpenAPITag {
  name: string;
  description?: string;
}

export interface PathItem {
  get?: Operation;
  post?: Operation;
  put?: Operation;
  patch?: Operation;
  delete?: Operation;
  options?: Operation;
  head?: Operation;
  summary?: string;
  description?: string;
  parameters?: Parameter[];
}

export interface Operation {
  operationId?: string;
  summary?: string;
  description?: string;
  tags?: string[];
  parameters?: Parameter[];
  requestBody?: RequestBody;
  responses?: Record<string, Response>;
  security?: SecurityRequirement[];
  deprecated?: boolean;
  "x-internal"?: boolean | string;
}

export interface Parameter {
  name: string;
  in: "query" | "header" | "path" | "cookie";
  description?: string;
  required?: boolean;
  schema?: SchemaObject;
  deprecated?: boolean;
  "x-internal"?: boolean | string;
}

export interface RequestBody {
  description?: string;
  required?: boolean;
  content?: Record<string, MediaType>;
  "x-internal"?: boolean | string;
}

export interface Response {
  description?: string;
  content?: Record<string, MediaType>;
  headers?: Record<string, Header>;
  "x-internal"?: boolean | string;
}

export interface MediaType {
  schema?: SchemaObject | Reference;
  example?: unknown;
  examples?: Record<string, Example>;
}

export interface Header {
  description?: string;
  required?: boolean;
  deprecated?: boolean;
  schema?: SchemaObject | Reference;
  example?: unknown;
  "x-internal"?: boolean | string;
}

export interface Example {
  summary?: string;
  description?: string;
  value?: unknown;
  externalValue?: string;
  "x-internal"?: boolean | string;
}

export interface OpenAPILink {
  operationRef?: string;
  operationId?: string;
  parameters?: Record<string, unknown>;
  requestBody?: unknown;
  description?: string;
  server?: Server;
  "x-internal"?: boolean | string;
}

export interface Server {
  url: string;
  description?: string;
  variables?: Record<string, ServerVariable>;
}

export interface ServerVariable {
  enum?: string[];
  default: string;
  description?: string;
}

export interface Callback {
  [expression: string]: PathItem;
}

export interface Reference {
  $ref: string;
}

export interface SecurityRequirement {
  [name: string]: string[];
}

export interface OpenAPIComponents {
  schemas?: Record<string, SchemaObject>;
  securitySchemes?: Record<string, SecurityScheme>;
  parameters?: Record<string, Parameter>;
  responses?: Record<string, Response>;
  requestBodies?: Record<string, RequestBody>;
  headers?: Record<string, Header>;
  examples?: Record<string, Example>;
  links?: Record<string, OpenAPILink>;
  callbacks?: Record<string, Callback>;
}

export interface SchemaObject {
  type?: string;
  format?: string;
  description?: string;
  properties?: Record<string, SchemaObject | Reference>;
  items?: SchemaObject | Reference;
  required?: string[];
  enum?: unknown[];
  $ref?: string;
  allOf?: (SchemaObject | Reference)[];
  oneOf?: (SchemaObject | Reference)[];
  anyOf?: (SchemaObject | Reference)[];
  nullable?: boolean;
  readOnly?: boolean;
  writeOnly?: boolean;
  deprecated?: boolean;
  example?: unknown;
  default?: unknown;
  additionalProperties?: boolean | SchemaObject | Reference;
  minLength?: number;
  maxLength?: number;
  minimum?: number;
  maximum?: number;
  pattern?: string;
  "x-codegen-schema-type"?: string;
  "x-internal"?: boolean | string;
}

export interface SecurityScheme {
  type: "apiKey" | "http" | "oauth2" | "openIdConnect";
  description?: string;
  name?: string;
  in?: "query" | "header" | "cookie";
  scheme?: string;
  bearerFormat?: string;
  flows?: OAuthFlows;
  openIdConnectUrl?: string;
}

export interface OAuthFlows {
  implicit?: OAuthFlow;
  password?: OAuthFlow;
  clientCredentials?: OAuthFlow;
  authorizationCode?: OAuthFlow;
}

export interface OAuthFlow {
  authorizationUrl?: string;
  tokenUrl?: string;
  refreshUrl?: string;
  scopes?: Record<string, string>;
}

export function isReference(obj: unknown): obj is Reference {
  return (
    typeof obj === "object" &&
    obj !== null &&
    "$ref" in obj &&
    typeof (obj as Reference).$ref === "string"
  );
}

export function getRefName(ref: string): string {
  const parts = ref.split("/");
  return parts[parts.length - 1] ?? ref;
}

export interface InternalMarker {
  "x-internal"?: boolean | string;
}

export function isInternal(obj: InternalMarker | undefined | null): boolean {
  const value = obj?.["x-internal"];
  // Handle both boolean true and string values like "server", "internal", etc.
  return value === true || (typeof value === "string" && value.length > 0);
}

export function getInternalValue(
  obj: InternalMarker | undefined | null,
): string | null {
  const value = obj?.["x-internal"];
  if (value === true) return "internal";
  if (typeof value === "string" && value.length > 0) return value;
  return null;
}
