import {
  Badge,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  CheckCircle2Icon,
  CircleIcon,
  CpuIcon,
  FileIcon,
  LayersIcon,
  RocketIcon,
  Separator,
  ServerIcon,
  ShieldIcon,
  UploadCloudIcon,
  ZapIcon,
} from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { useState } from "react";

export const Route = createFileRoute("/_app/configuration/")({
  component: ConfigurationPage,
});

// Mock current configuration - in real app this would come from API
const currentConfig = {
  api: {
    cors: {
      enabled: true,
      origins: ["http://localhost:3000"],
    },
    environment: "development",
    host: "0.0.0.0",
    port: 3001,
  },
  auth: {
    jwtExpiresIn: "7d",
    mfaEnabled: false,
    providers: ["github", "google", "email"],
    sessionTimeout: 3600,
  },
  database: {
    host: "localhost",
    maxConnections: 20,
    name: "archesai",
    port: 5432,
    ssl: false,
    type: "postgresql",
  },
  intelligence: {
    defaultModel: "gpt-4",
    embeddingModel: "text-embedding-3-small",
    maxTokens: 4096,
    temperature: 0.7,
  },
  logging: {
    format: "json",
    level: "info",
    outputs: ["stdout", "file"],
  },
  redis: {
    db: 0,
    host: "localhost",
    port: 6379,
    ttl: 86400,
  },
  storage: {
    bucket: "archesai-storage",
    provider: "s3",
    publicUrl: "https://cdn.archesai.com",
    region: "us-east-1",
  },
};

const configSections = [
  {
    bgColor: "bg-blue-50",
    color: "text-blue-600",
    description: "Core API server configuration",
    icon: ServerIcon,
    id: "api",
    items: [
      { label: "Host", type: "primary", value: currentConfig.api.host },
      { label: "Port", type: "primary", value: currentConfig.api.port },
      {
        label: "Environment",
        type: "badge",
        value: currentConfig.api.environment,
      },
      {
        label: "CORS",
        type: "status",
        value: currentConfig.api.cors.enabled ? "Enabled" : "Disabled",
      },
    ],
    title: "API Server",
  },
  {
    bgColor: "bg-green-50",
    color: "text-green-600",
    description: "Security and authentication settings",
    icon: ShieldIcon,
    id: "auth",
    items: [
      {
        label: "Providers",
        type: "primary",
        value: currentConfig.auth.providers.join(", "),
      },
      {
        label: "JWT Expiry",
        type: "secondary",
        value: currentConfig.auth.jwtExpiresIn,
      },
      {
        label: "Session Timeout",
        type: "secondary",
        value: `${currentConfig.auth.sessionTimeout}s`,
      },
      {
        label: "MFA",
        type: "status",
        value: currentConfig.auth.mfaEnabled ? "Enabled" : "Disabled",
      },
    ],
    title: "Authentication",
  },
  {
    bgColor: "bg-purple-50",
    color: "text-purple-600",
    description: "Database connection and pooling",
    icon: LayersIcon,
    id: "database",
    items: [
      { label: "Type", type: "badge", value: currentConfig.database.type },
      {
        label: "Host",
        type: "primary",
        value: `${currentConfig.database.host}:${currentConfig.database.port}`,
      },
      {
        label: "Database",
        type: "primary",
        value: currentConfig.database.name,
      },
      {
        label: "Max Connections",
        type: "secondary",
        value: currentConfig.database.maxConnections,
      },
      {
        label: "SSL",
        type: "status",
        value: currentConfig.database.ssl ? "Enabled" : "Disabled",
      },
    ],
    title: "Database",
  },
  {
    bgColor: "bg-red-50",
    color: "text-red-600",
    description: "Cache and session storage",
    icon: ZapIcon,
    id: "redis",
    items: [
      {
        label: "Host",
        type: "primary",
        value: `${currentConfig.redis.host}:${currentConfig.redis.port}`,
      },
      { label: "Database", type: "secondary", value: currentConfig.redis.db },
      { label: "TTL", type: "secondary", value: `${currentConfig.redis.ttl}s` },
    ],
    title: "Redis Cache",
  },
  {
    bgColor: "bg-indigo-50",
    color: "text-indigo-600",
    description: "AI models and settings",
    icon: CpuIcon,
    id: "intelligence",
    items: [
      {
        label: "Model",
        type: "badge",
        value: currentConfig.intelligence.defaultModel,
      },
      {
        label: "Embedding",
        type: "secondary",
        value: currentConfig.intelligence.embeddingModel,
      },
      {
        label: "Max Tokens",
        type: "secondary",
        value: currentConfig.intelligence.maxTokens,
      },
      {
        label: "Temperature",
        type: "secondary",
        value: currentConfig.intelligence.temperature,
      },
    ],
    title: "AI Intelligence",
  },
  {
    bgColor: "bg-orange-50",
    color: "text-orange-600",
    description: "File storage and CDN",
    icon: UploadCloudIcon,
    id: "storage",
    items: [
      {
        label: "Provider",
        type: "badge",
        value: currentConfig.storage.provider.toUpperCase(),
      },
      { label: "Bucket", type: "primary", value: currentConfig.storage.bucket },
      {
        label: "Region",
        type: "secondary",
        value: currentConfig.storage.region,
      },
      {
        label: "CDN URL",
        type: "secondary",
        value: currentConfig.storage.publicUrl,
      },
    ],
    title: "Storage",
  },
  {
    bgColor: "bg-gray-50",
    color: "text-gray-600",
    description: "Logging and monitoring",
    icon: FileIcon,
    id: "logging",
    items: [
      {
        label: "Level",
        type: "badge",
        value: currentConfig.logging.level.toUpperCase(),
      },
      {
        label: "Format",
        type: "secondary",
        value: currentConfig.logging.format,
      },
      {
        label: "Outputs",
        type: "secondary",
        value: currentConfig.logging.outputs.join(", "),
      },
    ],
    title: "Logging",
  },
];

function ConfigurationPage(): JSX.Element {
  const [selectedSection, setSelectedSection] = useState<string | null>(null);

  return (
    <div className="flex h-full gap-4">
      {/* Left sidebar with sections */}
      <div className="w-72 shrink-0">
        <div className="mb-3">
          <h1 className="font-bold text-xl">Configuration</h1>
          <p className="text-muted-foreground text-xs">
            Current system configuration
          </p>
        </div>

        <div
          className="h-[calc(100vh-10rem)] overflow-y-auto pr-1"
          style={{
            scrollbarColor: "#d1d5db transparent",
            scrollbarWidth: "thin",
          }}
        >
          <div className="space-y-1.5">
            {configSections.map((section) => {
              const Icon = section.icon;
              const isSelected = selectedSection === section.id;

              return (
                <Card
                  className={`cursor-pointer transition-all hover:shadow-sm ${
                    isSelected ? "border-primary ring-1 ring-primary" : ""
                  }`}
                  key={section.id}
                  onClick={() => setSelectedSection(section.id)}
                >
                  <CardHeader className="p-3">
                    <div className="flex items-start gap-2">
                      <div className={`rounded p-1.5 ${section.bgColor}`}>
                        <Icon className={`h-4 w-4 ${section.color}`} />
                      </div>
                      <div className="min-w-0 flex-1">
                        <CardTitle className="font-medium text-sm leading-tight">
                          {section.title}
                        </CardTitle>
                        <CardDescription className="mt-0.5 line-clamp-2 text-xs">
                          {section.description}
                        </CardDescription>
                      </div>
                    </div>
                  </CardHeader>
                </Card>
              );
            })}
          </div>
        </div>
      </div>

      {/* Right content area */}
      <div className="flex-1">
        {selectedSection ? (
          <div className="space-y-3">
            {(() => {
              const section = configSections.find(
                (s) => s.id === selectedSection,
              );
              if (!section) return null;
              const Icon = section.icon;

              return (
                <>
                  <Card>
                    <CardHeader className="pt-4 pb-3">
                      <div className="flex items-center gap-2">
                        <div className={`rounded p-2 ${section.bgColor}`}>
                          <Icon className={`h-5 w-5 ${section.color}`} />
                        </div>
                        <div>
                          <CardTitle className="text-base">
                            {section.title}
                          </CardTitle>
                          <CardDescription className="text-xs">
                            {section.description}
                          </CardDescription>
                        </div>
                      </div>
                    </CardHeader>
                    <Separator />
                    <CardContent className="pt-4">
                      <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
                        {section.items.map((item, itemIndex) => (
                          <div
                            className="space-y-1"
                            key={`${section.id}-item-${itemIndex}`}
                          >
                            <p className="font-medium text-muted-foreground text-xs">
                              {item.label}
                            </p>
                            <div className="flex items-center gap-1.5">
                              {item.type === "badge" ? (
                                <Badge
                                  className="h-5 font-mono text-xs"
                                  variant="secondary"
                                >
                                  {item.value}
                                </Badge>
                              ) : item.type === "status" ? (
                                <div className="flex items-center gap-1">
                                  {item.value === "Enabled" ? (
                                    <CheckCircle2Icon className="h-3.5 w-3.5 text-green-600" />
                                  ) : (
                                    <CircleIcon className="h-3.5 w-3.5 text-gray-400" />
                                  )}
                                  <span className="font-medium text-xs">
                                    {item.value}
                                  </span>
                                </div>
                              ) : item.type === "primary" ? (
                                <p className="font-mono font-semibold text-sm">
                                  {item.value}
                                </p>
                              ) : (
                                <p className="font-mono text-xs">
                                  {item.value}
                                </p>
                              )}
                            </div>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>

                  {/* Additional details card */}
                  <Card>
                    <CardHeader className="pt-4 pb-3">
                      <CardTitle className="text-sm">
                        Configuration Source
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2">
                      <div className="flex items-center justify-between rounded border p-2">
                        <div className="flex items-center gap-1.5">
                          <Badge
                            className="h-5 text-xs"
                            variant="outline"
                          >
                            YAML
                          </Badge>
                          <code className="text-xs">arches.yaml</code>
                        </div>
                        <span className="text-muted-foreground text-xs">
                          Primary
                        </span>
                      </div>
                      <div className="flex items-center justify-between rounded border p-2">
                        <div className="flex items-center gap-1.5">
                          <Badge
                            className="h-5 text-xs"
                            variant="outline"
                          >
                            ENV
                          </Badge>
                          <code className="text-xs">Environment vars</code>
                        </div>
                        <span className="text-muted-foreground text-xs">
                          Override
                        </span>
                      </div>
                      <div className="flex items-center justify-between rounded border p-2">
                        <div className="flex items-center gap-1.5">
                          <Badge
                            className="h-5 text-xs"
                            variant="outline"
                          >
                            CLI
                          </Badge>
                          <code className="text-xs">Command flags</code>
                        </div>
                        <span className="text-muted-foreground text-xs">
                          Runtime
                        </span>
                      </div>
                    </CardContent>
                  </Card>
                </>
              );
            })()}
          </div>
        ) : (
          <div className="flex h-full items-center justify-center">
            <div className="text-center">
              <RocketIcon className="mx-auto h-10 w-10 text-muted-foreground opacity-50" />
              <h3 className="mt-3 font-medium text-sm">
                Select a configuration section
              </h3>
              <p className="mt-1 text-muted-foreground text-xs">
                Choose a section from the left to view its settings
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
