import type { JSX } from "react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "../../shadcn/accordion";
import { Badge } from "../../shadcn/badge";
import type { Operation, PathItem, SchemaObject } from "./types";
import { getInternalValue, getRefName, isReference } from "./types";

interface PathViewerProps {
  path: string;
  pathItem: PathItem;
  schemas?: Record<string, SchemaObject> | undefined;
}

const methodColors: Record<string, string> = {
  delete: "bg-red-500/10 text-red-500 border-red-500/20",
  get: "bg-green-500/10 text-green-500 border-green-500/20",
  head: "bg-purple-500/10 text-purple-500 border-purple-500/20",
  options: "bg-gray-500/10 text-gray-500 border-gray-500/20",
  patch: "bg-orange-500/10 text-orange-500 border-orange-500/20",
  post: "bg-blue-500/10 text-blue-500 border-blue-500/20",
  put: "bg-yellow-500/10 text-yellow-500 border-yellow-500/20",
};

const HTTP_METHODS = [
  "get",
  "post",
  "put",
  "patch",
  "delete",
  "options",
  "head",
] as const;

export function PathViewer({
  path,
  pathItem,
  schemas,
}: PathViewerProps): JSX.Element {
  const allOperations = HTTP_METHODS.filter(
    (method) => pathItem[method] !== undefined,
  ).map((method) => ({
    internalValue: getInternalValue(pathItem[method]),
    method,
    operation: pathItem[method] as Operation,
  }));

  // Sort so internal operations are at the bottom
  const operations = allOperations.sort((a, b) => {
    if (Boolean(a.internalValue) === Boolean(b.internalValue)) return 0;
    return a.internalValue ? 1 : -1;
  });

  return (
    <div className="overflow-hidden rounded-lg border">
      <Accordion
        collapsible
        type="single"
      >
        {operations.map(({ method, operation, internalValue }) => (
          <AccordionItem
            className={internalValue ? "bg-sky-500/5 opacity-60" : ""}
            key={method}
            value={method}
          >
            <AccordionTrigger className="px-4 hover:no-underline">
              <div className="flex w-full items-center gap-3">
                <Badge
                  className={`font-mono text-xs uppercase ${internalValue ? "border-sky-500/20 bg-sky-500/10 text-sky-500" : methodColors[method]}`}
                  variant="outline"
                >
                  {method}
                </Badge>
                <code className="flex-1 text-left font-mono text-sm">
                  {path}
                </code>
                {operation.operationId && (
                  <span className="text-muted-foreground text-sm">
                    {operation.operationId}
                  </span>
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
                {operation.deprecated && (
                  <Badge
                    className="text-xs"
                    variant="destructive"
                  >
                    deprecated
                  </Badge>
                )}
              </div>
            </AccordionTrigger>
            <AccordionContent className="px-4 pb-4">
              <OperationDetails
                operation={operation}
                schemas={schemas}
              />
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  );
}

interface OperationDetailsProps {
  operation: Operation;
  schemas?: Record<string, SchemaObject> | undefined;
}

function OperationDetails({
  operation,
  schemas: _schemas,
}: OperationDetailsProps): JSX.Element {
  return (
    <div className="space-y-4">
      {operation.summary && (
        <div>
          <h5 className="mb-1 font-medium text-sm">Summary</h5>
          <p className="text-muted-foreground text-sm">{operation.summary}</p>
        </div>
      )}

      {operation.description && (
        <div>
          <h5 className="mb-1 font-medium text-sm">Description</h5>
          <p className="text-muted-foreground text-sm">
            {operation.description}
          </p>
        </div>
      )}

      {operation.tags && operation.tags.length > 0 && (
        <div>
          <h5 className="mb-1 font-medium text-sm">Tags</h5>
          <div className="flex flex-wrap gap-1">
            {operation.tags.map((tag) => (
              <Badge
                key={tag}
                variant="secondary"
              >
                {tag}
              </Badge>
            ))}
          </div>
        </div>
      )}

      {operation.parameters && operation.parameters.length > 0 && (
        <div>
          <h5 className="mb-2 font-medium text-sm">Parameters</h5>
          <div className="space-y-2">
            {operation.parameters.map((param) => (
              <div
                className="flex items-start gap-2 text-sm"
                key={`${param.in}-${param.name}`}
              >
                <code className="rounded bg-muted px-1 font-mono">
                  {param.name}
                </code>
                <Badge
                  className="text-xs"
                  variant="outline"
                >
                  {param.in}
                </Badge>
                {param.required && (
                  <Badge
                    className="text-xs"
                    variant="default"
                  >
                    required
                  </Badge>
                )}
                {param.schema && (
                  <Badge
                    className="text-xs"
                    variant="secondary"
                  >
                    {param.schema.type ?? "object"}
                  </Badge>
                )}
                {param.description && (
                  <span className="text-muted-foreground">
                    {param.description}
                  </span>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {operation.requestBody && (
        <div>
          <h5 className="mb-2 font-medium text-sm">Request Body</h5>
          {operation.requestBody.required && (
            <Badge
              className="mb-2 text-xs"
              variant="default"
            >
              required
            </Badge>
          )}
          {operation.requestBody.description && (
            <p className="mb-2 text-muted-foreground text-sm">
              {operation.requestBody.description}
            </p>
          )}
          {operation.requestBody.content && (
            <div className="space-y-2">
              {Object.entries(operation.requestBody.content).map(
                ([contentType, mediaType]) => (
                  <div
                    className="text-sm"
                    key={contentType}
                  >
                    <Badge
                      className="mb-1"
                      variant="outline"
                    >
                      {contentType}
                    </Badge>
                    {mediaType.schema && (
                      <div className="mt-1 pl-2">
                        {isReference(mediaType.schema) ? (
                          <code className="rounded bg-muted px-1 text-xs">
                            {getRefName(mediaType.schema.$ref)}
                          </code>
                        ) : (
                          <code className="rounded bg-muted px-1 text-xs">
                            {mediaType.schema.type ?? "object"}
                          </code>
                        )}
                      </div>
                    )}
                  </div>
                ),
              )}
            </div>
          )}
        </div>
      )}

      {operation.responses && Object.keys(operation.responses).length > 0 && (
        <div>
          <h5 className="mb-2 font-medium text-sm">Responses</h5>
          <div className="space-y-2">
            {Object.entries(operation.responses).map(([code, response]) => (
              <div
                className="flex items-start gap-2 border-l-2 pl-2 text-sm"
                key={code}
                style={{
                  borderColor: code.startsWith("2")
                    ? "rgb(34, 197, 94)"
                    : code.startsWith("4")
                      ? "rgb(239, 68, 68)"
                      : code.startsWith("5")
                        ? "rgb(249, 115, 22)"
                        : "rgb(156, 163, 175)",
                }}
              >
                <Badge
                  className={
                    code.startsWith("2")
                      ? "bg-green-500/10 text-green-500"
                      : code.startsWith("4")
                        ? "bg-red-500/10 text-red-500"
                        : code.startsWith("5")
                          ? "bg-orange-500/10 text-orange-500"
                          : ""
                  }
                  variant="outline"
                >
                  {code}
                </Badge>
                <div className="flex-1">
                  {response.description && (
                    <span className="text-muted-foreground">
                      {response.description}
                    </span>
                  )}
                  {response.content &&
                    Object.entries(response.content).map(
                      ([contentType, mediaType]) => (
                        <div
                          className="mt-1"
                          key={contentType}
                        >
                          <Badge
                            className="text-xs"
                            variant="outline"
                          >
                            {contentType}
                          </Badge>
                          {mediaType.schema && (
                            <span className="ml-2">
                              {isReference(mediaType.schema) ? (
                                <code className="rounded bg-muted px-1 text-xs">
                                  {getRefName(mediaType.schema.$ref)}
                                </code>
                              ) : (
                                <code className="rounded bg-muted px-1 text-xs">
                                  {mediaType.schema.type ?? "object"}
                                </code>
                              )}
                            </span>
                          )}
                        </div>
                      ),
                    )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {operation.security && operation.security.length > 0 && (
        <div>
          <h5 className="mb-2 font-medium text-sm">Security</h5>
          <div className="flex flex-wrap gap-1">
            {operation.security.map((secReq, idx) =>
              Object.keys(secReq).map((name) => (
                <Badge
                  key={`${idx}-${name}`}
                  variant="outline"
                >
                  {name}
                </Badge>
              )),
            )}
          </div>
        </div>
      )}
    </div>
  );
}
